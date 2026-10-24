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
