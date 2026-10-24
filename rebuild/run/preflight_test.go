package run_test

import (
	"testing"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

func TestPreflight_PassesFromNothing(t *testing.T) {
	for name, target := range map[string]*runtest.Target{
		"with no state backend":         runtest.NewTarget(t, runtest.Drill, runtest.NoBackend()),
		"with a backend and no records": runtest.NewTarget(t, runtest.Drill),
	} {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, run.Preflight(t.Context(), target, phases(nil)))
			require.Empty(t, target.Calls())
		})
	}
}

func TestPreflight_StopsOnAnUnfinishedRunAndSaysItCanBeResumedOrTornDown(t *testing.T) {
	for name, tc := range map[string]struct {
		stop *run.Stop
		want string
	}{
		"that failed": {
			stop: &run.Stop{Phase: "1a", Step: run.StepApply, Error: "boom"},
			want: "stopped at phase 1a's apply",
		},
		"that paused at a checkpoint": {
			stop: &run.Stop{Phase: "0", Step: run.StepCheckpoint},
			want: "paused at the checkpoint after phase 0",
		},
		"that stopped where Through said": {
			stop: &run.Stop{Phase: "0", Step: run.StepThrough},
			want: "stopped after phase 0, as it was asked to",
		},
		"that never said where it stopped": {
			want: "never recorded where it stopped",
		},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Drill)
			require.NoError(t, target.Store().Write(t.Context(), run.Record{
				ID:      "earlier",
				Phases:  []run.PhaseRecord{{Name: "0", Verified: true}},
				Stopped: tc.stop,
			}))

			err := run.Preflight(t.Context(), target, phases(nil))
			require.ErrorIs(t, err, run.ErrUnfinishedRun)
			require.ErrorContains(t, err, "run earlier "+tc.want)
			require.ErrorContains(t, err, "Resume it with up --resume, or tear the target down with down")
			require.Empty(t, target.Calls(), "nothing is verified on top of a partial run")

			err = up(t.Context(), target, phases(nil), run.NoPause())
			require.ErrorIs(t, err, run.ErrUnfinishedRun)
			require.Empty(t, target.Calls())
		})
	}
}

func TestPreflight_VerifiesEveryPhaseTheLastRunApplied(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))
	applied := len(target.Calls())

	require.NoError(t, run.Preflight(t.Context(), target, phases(nil)))
	require.Equal(t, []string{"0 verify", "1a verify", "1b verify"}, target.Calls()[applied:])
}

func TestPreflight_StopsOnAFailedVerify(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))
	applied := len(target.Calls())

	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailVerify(errNoEdge)}})
	err := run.Preflight(t.Context(), target, ps)
	require.ErrorIs(t, err, run.ErrDrift)
	require.ErrorIs(t, err, errNoEdge)
	require.ErrorContains(t, err, "phase 1a")
	require.Equal(t, []string{"0 verify", "1a verify"}, target.Calls()[applied:])
}

func TestPreflight_StopsWhenTheLastRunAppliedAPhaseNoLongerInTheOrder(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))

	err := run.Preflight(t.Context(), target, []run.Phase{runtest.Phase("0"), runtest.Phase("1a")})
	require.ErrorIs(t, err, run.ErrUnknownPhase)
	require.ErrorContains(t, err, "phase 1b")
}

func TestPreflight_ChangesNothingSoItRunsAgainstLiveFromUncommittedChanges(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Live, runtest.WithSource(run.Source{
		Repo: run.Commit{SHA: "a", Uncommitted: true},
	}))

	require.NoError(t, run.Preflight(t.Context(), target, phases(nil)))
}

func TestPreflight_StopsWhenTheTargetsRecordsCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailStore(errBroken))

	err := run.Preflight(t.Context(), target, phases(nil))
	require.ErrorIs(t, err, errBroken)
	require.ErrorContains(t, err, "reading the target's records")
}

func TestPreflight_StopsWhenTheTargetsSourceCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailSource(errBroken))

	require.ErrorIs(t, run.Preflight(t.Context(), target, phases(nil)), errBroken)
}
