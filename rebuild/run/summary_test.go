package run_test

import (
	"testing"
	"time"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

func TestSummary_AtACheckpointSaysWhatThePhaseMadeItsTimeAndTheRunsTimeSoFar(t *testing.T) {
	clock := runtest.NewClock(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC))
	target := runtest.NewTarget(t, runtest.Drill, runtest.Named("scratch"))
	ps := phases(map[string][]runtest.PhaseOption{
		"0": {runtest.Takes(clock, 4*time.Minute), runtest.Makes("the state backend", "the run's key")},
	})

	sum, err := run.Up(t.Context(), target, ps, run.WithClock(clock))
	require.NoError(t, err)
	require.Equal(t, "run "+sum.ID+" against scratch\n"+
		"phase 0: verified, 8m0s\n"+
		"  made the state backend\n"+
		"  made the run's key\n"+
		"the run's time so far: 8m0s\n"+
		"the run paused at the checkpoint after phase 0. Continue with up --resume, or tear the target down with down.\n",
		sum.String())
}

func TestSummary_OfAFailureSaysWhatFailedAndHowToGoOn(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailVerify(errNoEdge)}})

	clock := runtest.NewClock(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC))

	sum, err := run.Up(t.Context(), target, ps, run.WithClock(clock), run.NoPause())
	require.ErrorIs(t, err, errNoEdge)
	got := sum.String()
	require.Contains(t, got, "phase 0: verified, 0s\n")
	require.Contains(t, got, "phase 1a: did not verify, 0s\n")
	require.Contains(t, got, "the run's time so far: 0s\n")
	require.Contains(t, got, "the run failed at phase 1a's verify: the edge does not answer\n")
	require.Contains(t, got, "continue with up --resume, which applies phase 1a again. Or tear the target down with down.\n")
}

func TestSummary_OfAFinishedRunSaysItFinishedAndWhenItWentOver(t *testing.T) {
	clock := runtest.NewClock(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC))
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{"1b": {runtest.Takes(clock, 31*time.Minute)}})

	sum, err := run.Up(t.Context(), target, ps, run.WithClock(clock), run.NoPause())
	require.NoError(t, err)
	got := sum.String()
	require.Contains(t, got, "phase 1b: verified, 1h2m0s\n")
	require.Contains(t, got, "the run took 1h2m0s, over the 1h0m0s target\n")
	require.Contains(t, got, "the run finished: every phase verified.\n")
}

func TestSummary_OfARunStoppedWhereThroughSaidSaysSo(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	sum, err := run.Up(t.Context(), target, phases(nil), run.Through("0"))
	require.NoError(t, err)
	require.Contains(t, sum.String(),
		"the run stopped after phase 0, as it was asked to. Continue with up --resume, or tear the target down with down.\n")
}

func TestSummary_OfARunStillGoingSaysSo(t *testing.T) {
	sum := run.Summary{ID: "r1", Target: "scratch", Total: time.Minute}

	require.Equal(t, "run r1 against scratch\nthe run's time so far: 1m0s\nthe run is still going.\n", sum.String())
}
