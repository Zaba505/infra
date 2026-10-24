package run_test

import (
	"encoding/json"
	"testing"
	"time"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

func TestLastKnownGood_IsTheLastRunThatFinished(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))
	require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))
	good := lastRun(t, target)
	require.Len(t, records(t, target), 2)

	ps := phases(map[string][]runtest.PhaseOption{"1b": {runtest.FailApply(errNoEdge)}})
	require.ErrorIs(t, up(t.Context(), target, ps, run.NoPause()), errNoEdge)
	failed := lastRun(t, target)
	require.NotEqual(t, good.ID, failed.ID)

	lkg, ok, err := run.LastKnownGood(t.Context(), target.Store())
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, good, lkg)
}

func TestLastKnownGood_IsNothingUntilARunFinishes(t *testing.T) {
	for name, target := range map[string]*runtest.Target{
		"with no state backend": runtest.NewTarget(t, runtest.Drill, runtest.NoBackend()),
		"with no records":       runtest.NewTarget(t, runtest.Drill),
		"after a run stopped where Through said": func() *runtest.Target {
			target := runtest.NewTarget(t, runtest.Drill)
			require.NoError(t, up(t.Context(), target, phases(nil), run.Through("0")))
			return target
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			_, ok, err := run.LastKnownGood(t.Context(), target.Store())
			require.NoError(t, err)
			require.False(t, ok)
		})
	}
}

func TestLastKnownGood_FailsWhenTheRecordsCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailStore(errBroken))

	_, _, err := run.LastKnownGood(t.Context(), target.Store())
	require.ErrorIs(t, err, errBroken)
	_, _, err = run.LastRun(t.Context(), target.Store())
	require.ErrorIs(t, err, errBroken)
}

func TestRecord_SurvivesARoundTripThroughJSON(t *testing.T) {
	rec := run.Record{
		ID:      "r1",
		Target:  "scratch",
		Started: time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC),
		Phases: []run.PhaseRecord{{
			Name: "0",
			Source: run.Source{
				Repo:   run.Commit{SHA: "a", Uncommitted: true},
				Values: run.Commit{SHA: "b"},
			},
			Made:     []string{"the state backend"},
			Duration: 3 * time.Minute,
			Verified: true,
		}},
		Stopped: &run.Stop{Phase: "0", Step: run.StepCheckpoint},
		Total:   61 * time.Minute,
		Over:    true,
	}

	b, err := json.Marshal(rec)
	require.NoError(t, err)
	var got run.Record
	require.NoError(t, json.Unmarshal(b, &got))
	require.Equal(t, rec, got)
}

func TestKind_String(t *testing.T) {
	require.Equal(t, "live", run.Live.String())
	require.Equal(t, "drill", run.Drill.String())
	require.Equal(t, "Kind(0)", run.Kind(0).String())
}
