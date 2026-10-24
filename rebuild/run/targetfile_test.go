package run_test

import (
	"testing"

	"dagger/rebuild/run"

	"github.com/stretchr/testify/require"
)

func TestTargetFilePath_IsTheTargetsFileInTheTargetsDirectory(t *testing.T) {
	for _, name := range []string{"live", "drill", "drill-2", "scratch_1", "a.b"} {
		t.Run(name, func(t *testing.T) {
			got, err := run.TargetFilePath(name)
			require.NoError(t, err)
			require.Equal(t, "targets/"+name+".yaml", got)
		})
	}
}

func TestTargetFilePath_RefusesNoName(t *testing.T) {
	_, err := run.TargetFilePath("")
	require.ErrorIs(t, err, run.ErrNoTarget)
}

func TestTargetFilePath_RefusesANameThatReachesOutsideTheTargetsDirectory(t *testing.T) {
	for _, name := range []string{"..", ".hidden", "a/b", "../live", "a\\b", "with space", "-flag"} {
		t.Run(name, func(t *testing.T) {
			_, err := run.TargetFilePath(name)
			require.ErrorIs(t, err, run.ErrUnknownTarget)
			require.ErrorContains(t, err, "is not a target name")
		})
	}
}

func TestParseTargetFile_ReadsTheKind(t *testing.T) {
	for text, want := range map[string]run.Kind{
		"kind: live\n":                    run.Live,
		"kind: drill\n":                   run.Drill,
		"---\nkind: drill\n":              run.Drill,
		"# a scratch target\nkind: drill": run.Drill,
	} {
		t.Run(text, func(t *testing.T) {
			got, err := run.ParseTargetFile([]byte(text))
			require.NoError(t, err)
			require.Equal(t, run.TargetFile{Kind: want}, got)
		})
	}
}

func TestParseTargetFile_IsStrict(t *testing.T) {
	for name, tc := range map[string]struct {
		text string
		want string
	}{
		"an empty file":           {text: "", want: "it is empty"},
		"a file of comments only": {text: "# nothing here\n", want: "it is empty"},
		"no kind":                 {text: "{}\n", want: "it must say kind: live or kind: drill"},
		"a null kind":             {text: "kind:\n", want: "it must say kind: live or kind: drill"},
		"an unknown kind":         {text: "kind: staging\n", want: `kind "staging" is neither live nor drill`},
		"a kind in capitals":      {text: "kind: Live\n", want: `kind "Live" is neither live nor drill`},
		"a numeric kind":          {text: "kind: 1\n", want: `kind "1" is neither live nor drill`},
		"an unknown key":          {text: "kind: drill\nproject: p\n", want: "field project not found"},
		"two documents":           {text: "kind: drill\n---\nkind: live\n", want: "it must hold one YAML document"},
		"not a mapping":           {text: "- kind: drill\n", want: "cannot unmarshal"},
		"not YAML":                {text: "kind: [drill\n", want: "invalid target file"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := run.ParseTargetFile([]byte(tc.text))
			require.ErrorIs(t, err, run.ErrInvalidTargetFile)
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestTargetKind_ReadsTheTargetsKind(t *testing.T) {
	files := map[string][]byte{
		"live":  []byte("kind: live\n"),
		"drill": []byte("kind: drill\n"),
	}
	for name, want := range map[string]run.Kind{"live": run.Live, "drill": run.Drill} {
		t.Run(name, func(t *testing.T) {
			got, err := run.TargetKind(name, files, files)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}

func TestTargetKind_ATargetIsLiveWhenEitherVersionOfItsFileSaysSo(t *testing.T) {
	live := []byte("kind: live\n")
	drill := []byte("kind: drill\n")
	for name, tc := range map[string]struct {
		committed, working map[string][]byte
		want               run.Kind
	}{
		"an uncommitted edit to drill over a committed live": {
			committed: map[string][]byte{"t": live},
			working:   map[string][]byte{"t": drill},
			want:      run.Live,
		},
		"an uncommitted edit to live over a committed drill": {
			committed: map[string][]byte{"t": drill},
			working:   map[string][]byte{"t": live},
			want:      run.Live,
		},
		"a drill not committed yet": {
			committed: map[string][]byte{},
			working:   map[string][]byte{"t": drill},
			want:      run.Drill,
		},
		"a drill in both": {
			committed: map[string][]byte{"t": drill},
			working:   map[string][]byte{"t": drill},
			want:      run.Drill,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := run.TargetKind("t", tc.committed, tc.working)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestTargetKind_RefusesMoreThanOneLiveTarget(t *testing.T) {
	live := []byte("kind: live\n")
	drill := []byte("kind: drill\n")
	for name, tc := range map[string]struct {
		committed, working map[string][]byte
	}{
		"two live files": {
			committed: map[string][]byte{"live": live, "prod": live},
			working:   map[string][]byte{"live": live, "prod": live},
		},
		"a second live target not committed yet": {
			committed: map[string][]byte{"live": live, "prod": drill},
			working:   map[string][]byte{"live": live, "prod": live},
		},
		"a live target removed but not committed": {
			committed: map[string][]byte{"live": live, "prod": live},
			working:   map[string][]byte{"live": live},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, target := range []string{"live", "prod"} {
				_, err := run.TargetKind(target, tc.committed, tc.working)
				require.ErrorIs(t, err, run.ErrMoreThanOneLive)
				require.ErrorContains(t, err, "live, prod each say kind: live")
			}
		})
	}
}

func TestTargetKind_RefusesATargetWithNoFile(t *testing.T) {
	files := map[string][]byte{"drill": []byte("kind: drill\n")}
	for name, committed := range map[string]map[string][]byte{
		"that never had one":                 files,
		"whose file is removed, uncommitted": {"drill": files["drill"], "staging": []byte("kind: drill\n")},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := run.TargetKind("staging", committed, files)
			require.ErrorIs(t, err, run.ErrUnknownTarget)
			require.ErrorContains(t, err, "staging: the values checkout has no targets/staging.yaml")
		})
	}
}

func TestTargetKind_ReadsEveryFileStrictlyInBothVersions(t *testing.T) {
	good := map[string][]byte{"drill": []byte("kind: drill\n")}
	bad := map[string][]byte{"drill": []byte("kind: drill\n"), "other": []byte("kind: staging\n")}
	for name, tc := range map[string]struct {
		committed, working map[string][]byte
		want               string
	}{
		"another target's file as it stands": {committed: good, working: bad, want: "target other: targets/other.yaml, as it stands"},
		"another target's file as committed": {committed: bad, working: good, want: "target other: targets/other.yaml, as committed"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := run.TargetKind("drill", tc.committed, tc.working)
			require.ErrorIs(t, err, run.ErrInvalidTargetFile)
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestTargetKind_RefusesABadName(t *testing.T) {
	_, err := run.TargetKind("../live", nil, nil)
	require.ErrorIs(t, err, run.ErrUnknownTarget)
	_, err = run.TargetKind("", nil, nil)
	require.ErrorIs(t, err, run.ErrNoTarget)
}
