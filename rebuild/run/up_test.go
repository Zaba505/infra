package run_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

var (
	errNoEdge = errors.New("the edge does not answer")
	errBroken = errors.New("broken")
)

// order is the phase order the tests run, in the shape the rebuild's own takes.
var order = []string{"0", "1a", "1b"}

func phases(opts map[string][]runtest.PhaseOption) []run.Phase {
	ps := make([]run.Phase, len(order))
	for i, name := range order {
		ps[i] = runtest.Phase(name, opts[name]...)
	}
	return ps
}

func lastRun(t *testing.T, target run.Target) run.Record {
	t.Helper()
	rec, ok, err := run.LastRun(t.Context(), target.Store())
	require.NoError(t, err)
	require.True(t, ok, "the target holds no record")
	return rec
}

func records(t *testing.T, target run.Target) []run.Record {
	t.Helper()
	recs, err := target.Store().Records(t.Context())
	require.NoError(t, err)
	return recs
}

func TestUp_RunsEachPhasesApplyThenVerifyInOrder(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	err := run.Up(t.Context(), target, phases(nil))
	require.NoError(t, err)
	require.Equal(t, []string{
		"0 apply", "0 verify", "1a apply", "1a verify", "1b apply", "1b verify",
	}, target.Calls())

	rec := lastRun(t, target)
	require.True(t, rec.Finished)
	require.Nil(t, rec.Stopped)
	require.Equal(t, "drill", rec.Target)
	require.NotEmpty(t, rec.ID)
	require.Len(t, rec.Phases, 3)
	for i, p := range rec.Phases {
		require.Equal(t, order[i], p.Name)
		require.True(t, p.Verified)
		require.Equal(t, "0000000000000000000000000000000000000001", p.Source.Repo.SHA)
		require.Equal(t, "0000000000000000000000000000000000000002", p.Source.Values.SHA)
		require.False(t, p.Source.Uncommitted())
	}
}

func TestUp_WithNoPhasesFinishesWithNothingToDo(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	require.NoError(t, run.Up(t.Context(), target, nil))
	require.Empty(t, target.Calls())
	rec := lastRun(t, target)
	require.True(t, rec.Finished)
	require.Empty(t, rec.Phases)
}

func TestUp_EachRunGetsItsOwnID(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	require.NoError(t, run.Up(t.Context(), target, phases(nil)))
	require.NoError(t, run.Up(t.Context(), target, phases(nil)))

	recs := records(t, target)
	require.Len(t, recs, 2)
	require.NotEqual(t, recs[0].ID, recs[1].ID)
}

func TestUp_StopsAfterThePhaseThroughNames(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	err := run.Up(t.Context(), target, phases(nil), run.Through("1a"))
	require.NoError(t, err)
	require.Equal(t, []string{"0 apply", "0 verify", "1a apply", "1a verify"}, target.Calls())

	rec := lastRun(t, target)
	require.False(t, rec.Finished, "a run that stopped short of the last phase hasn't finished")
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepThrough}, rec.Stopped)
	require.Len(t, rec.Phases, 2)
}

func TestUp_ThroughTheLastPhaseFinishes(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	require.NoError(t, run.Up(t.Context(), target, phases(nil), run.Through("1b")))
	require.True(t, lastRun(t, target).Finished)
}

func TestUp_ThroughAPhaseNotInTheOrderRunsNothing(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	err := run.Up(t.Context(), target, phases(nil), run.Through("9"))
	require.ErrorIs(t, err, run.ErrUnknownPhase)
	require.Empty(t, target.Calls())
	require.Empty(t, records(t, target))
}

func TestUp_AFailedApplyStopsTheRunAtOnce(t *testing.T) {
	for i, failing := range order {
		t.Run(failing, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Drill)

			err := run.Up(t.Context(), target, phases(map[string][]runtest.PhaseOption{
				failing: {runtest.FailApply(errNoEdge)},
			}))
			require.ErrorIs(t, err, errNoEdge)
			require.ErrorContains(t, err, "phase "+failing+": apply")

			var want []string
			for _, name := range order[:i] {
				want = append(want, name+" apply", name+" verify")
			}
			want = append(want, failing+" apply")
			require.Equal(t, want, target.Calls(), "no later phase runs")

			rec := lastRun(t, target)
			require.False(t, rec.Finished)
			require.Equal(t, &run.Stop{Phase: failing, Step: run.StepApply, Error: errNoEdge.Error()}, rec.Stopped)
			require.Len(t, rec.Phases, i+1)
			require.False(t, rec.Phases[i].Verified)
		})
	}
}

func TestUp_AFailedVerifyStopsTheRunAndLeavesItUnfinished(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailVerify(errNoEdge)}})

	err := run.Up(t.Context(), target, ps)
	require.ErrorIs(t, err, errNoEdge)
	require.ErrorContains(t, err, "phase 1a: verify")
	require.Equal(t, []string{"0 apply", "0 verify", "1a apply", "1a verify"}, target.Calls())

	rec := lastRun(t, target)
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepVerify, Error: errNoEdge.Error()}, rec.Stopped)

	err = run.Preflight(t.Context(), target, ps)
	require.ErrorIs(t, err, run.ErrUnfinishedRun)
}

func TestUp_RunsThePreflightCheckBeforeApplyingAnything(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, target.Store().Write(t.Context(), run.Record{
		ID:       "earlier",
		Finished: true,
		Phases:   []run.PhaseRecord{{Name: "0", Verified: true}, {Name: "1a", Verified: true}},
	}))

	require.NoError(t, run.Up(t.Context(), target, phases(nil)))
	require.Equal(t, []string{
		"0 verify", "1a verify",
		"0 apply", "0 verify", "1a apply", "1a verify", "1b apply", "1b verify",
	}, target.Calls())
}

func TestUp_StopsBeforeAnythingChangesWhenThePreflightCheckFails(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, run.Up(t.Context(), target, phases(nil)))
	before := target.Calls()

	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailVerify(errNoEdge)}})
	err := run.Up(t.Context(), target, ps)
	require.ErrorIs(t, err, run.ErrDrift)
	require.ErrorIs(t, err, errNoEdge)
	require.Equal(t, append(before, "0 verify", "1a verify"), target.Calls())
	require.Len(t, records(t, target), 1, "a run refused by the preflight check leaves no record")
}

func TestUp_ALiveRunRefusesUncommittedChanges(t *testing.T) {
	for name, tc := range map[string]struct {
		src  run.Source
		want string
	}{
		"in this repository": {
			src:  run.Source{Repo: run.Commit{SHA: "a", Uncommitted: true}, Values: run.Commit{SHA: "b"}},
			want: "this repository has uncommitted changes",
		},
		"in the target's values": {
			src:  run.Source{Repo: run.Commit{SHA: "a"}, Values: run.Commit{SHA: "b", Uncommitted: true}},
			want: "the target's values have uncommitted changes",
		},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Live, runtest.WithSource(tc.src))

			err := run.Up(t.Context(), target, phases(nil))
			require.ErrorIs(t, err, run.ErrUncommitted)
			require.ErrorContains(t, err, tc.want)
			require.Empty(t, target.Calls())
			require.Empty(t, records(t, target))
		})
	}
}

func TestUp_ALiveRunFromCommittedChangesRuns(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Live)

	require.NoError(t, run.Up(t.Context(), target, phases(nil)))
	require.True(t, lastRun(t, target).Finished)
}

func TestUp_ADrillMayRunFromUncommittedChangesAndItsRecordSaysSo(t *testing.T) {
	src := run.Source{
		Repo:   run.Commit{SHA: "a", Uncommitted: true},
		Values: run.Commit{SHA: "b", Uncommitted: true},
	}
	target := runtest.NewTarget(t, runtest.Drill, runtest.WithSource(src))

	require.NoError(t, run.Up(t.Context(), target, phases(nil)))
	for _, p := range lastRun(t, target).Phases {
		require.Equal(t, src, p.Source)
		require.True(t, p.Source.Uncommitted())
	}
}

func TestUp_KeepsItsRecordOncePhaseZeroHasMadeTheStateBackend(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())
	ps := phases(map[string][]runtest.PhaseOption{"0": {runtest.ProvidesBackend()}})

	require.NoError(t, run.Up(t.Context(), target, ps))
	require.Len(t, target.Calls(), 6)
	rec := lastRun(t, target)
	require.True(t, rec.Finished)
	require.Len(t, rec.Phases, 3)
}

func TestUp_WithNoStateBackendAndNoPhasesFinishesWithNothingToKeep(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())

	require.NoError(t, run.Up(t.Context(), target, nil))
	_, ok, err := run.LastRun(t.Context(), target.Store())
	require.NoError(t, err)
	require.False(t, ok)
}

func TestUp_AFailedFirstPhaseNeedsNoStateBackendToStop(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())
	ps := phases(map[string][]runtest.PhaseOption{"0": {runtest.ProvidesBackend(), runtest.FailApply(errNoEdge)}})

	err := run.Up(t.Context(), target, ps)
	require.ErrorIs(t, err, errNoEdge)
	require.NotErrorIs(t, err, run.ErrNoBackend)
	require.Equal(t, []string{"0 apply"}, target.Calls())
}

func TestUp_StopsWhenThereIsStillNoStateBackendAfterPhaseZero(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())

	err := run.Up(t.Context(), target, phases(nil))
	require.ErrorIs(t, err, run.ErrNoBackend)
	require.Equal(t, []string{"0 apply", "0 verify"}, target.Calls(), "no later phase runs unrecorded")
}

func TestUp_RecordsWhereItStoppedWhenItsContextIsCancelled(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ps := []run.Phase{runtest.Phase("0"), cancelling{name: "1a", cancel: cancel}, runtest.Phase("1b")}

	err := run.Up(ctx, target, ps)
	require.ErrorIs(t, err, context.Canceled)
	rec := lastRun(t, target)
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepApply, Error: context.Canceled.Error()}, rec.Stopped)
}

// cancelling is a phase whose Apply cancels the run's context, as an operator
// interrupting the run does, and fails with the context's error.
type cancelling struct {
	name   string
	cancel context.CancelFunc
}

func (c cancelling) Name() string { return c.name }

func (c cancelling) Apply(ctx context.Context, _ *run.Run) error {
	c.cancel()
	return ctx.Err()
}

func (c cancelling) Verify(context.Context, *run.Run) error   { return nil }
func (c cancelling) Teardown(context.Context, *run.Run) error { return nil }

func TestUp_StopsWhenItsRecordCannotBeWritten(t *testing.T) {
	for name, tc := range map[string]struct {
		writes int
		opts   map[string][]runtest.PhaseOption
		calls  []string
		also   error
	}{
		"when the run starts": {writes: 0},
		"after a phase": {
			writes: 1,
			calls:  []string{"0 apply", "0 verify"},
		},
		"after a phase failed": {
			writes: 1,
			opts:   map[string][]runtest.PhaseOption{"0": {runtest.FailApply(errNoEdge)}},
			calls:  []string{"0 apply"},
			also:   errNoEdge,
		},
		"when the run finishes": {
			writes: 4,
			calls:  []string{"0 apply", "0 verify", "1a apply", "1a verify", "1b apply", "1b verify"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Drill, runtest.FailStoreWrites(tc.writes, errBroken))

			err := run.Up(t.Context(), target, phases(tc.opts))
			require.ErrorIs(t, err, errBroken)
			require.ErrorContains(t, err, "writing run")
			if tc.also != nil {
				require.ErrorIs(t, err, tc.also)
			}
			require.Equal(t, tc.calls, target.Calls())
		})
	}
}

func TestUp_StopsWhenTheTargetsRecordsCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailStore(errBroken))

	err := run.Up(t.Context(), target, phases(nil))
	require.ErrorIs(t, err, errBroken)
	require.Empty(t, target.Calls())
}

func TestUp_StopsWhenTheTargetsSourceCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailSource(errBroken))

	err := run.Up(t.Context(), target, phases(nil))
	require.ErrorIs(t, err, errBroken)
	require.Empty(t, target.Calls())
}

func TestUp_MeasuresEachPhaseAndTheTotalFromTheStartOfUp(t *testing.T) {
	start := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	clock := runtest.NewClock(start)
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, target.Store().Write(t.Context(), run.Record{
		ID: "earlier", Finished: true, Phases: []run.PhaseRecord{{Name: "0", Verified: true}},
	}))
	// Every apply and every verify takes five minutes, so each phase takes
	// ten, and the preflight check's verify of phase 0 takes five.
	ps := phases(map[string][]runtest.PhaseOption{
		"0":  {runtest.Takes(clock, 5*time.Minute)},
		"1a": {runtest.Takes(clock, 5*time.Minute)},
		"1b": {runtest.Takes(clock, 5*time.Minute)},
	})

	require.NoError(t, run.Up(t.Context(), target, ps, run.WithClock(clock)))

	rec := lastRun(t, target)
	require.Equal(t, start, rec.Started)
	for _, p := range rec.Phases {
		require.Equal(t, 10*time.Minute, p.Duration, "phase %s", p.Name)
	}
	require.Equal(t, 35*time.Minute, rec.Total, "the total includes the preflight check")
	require.False(t, rec.Over)
}

func TestUp_AFailedRunsTimeRunsToWhereItStopped(t *testing.T) {
	clock := runtest.NewClock(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC))
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{
		"0":  {runtest.Takes(clock, 5*time.Minute)},
		"1a": {runtest.Takes(clock, 5*time.Minute), runtest.FailVerify(errNoEdge)},
	})

	require.ErrorIs(t, run.Up(t.Context(), target, ps, run.WithClock(clock)), errNoEdge)

	rec := lastRun(t, target)
	require.Equal(t, 20*time.Minute, rec.Total)
	require.Equal(t, 10*time.Minute, rec.Phases[1].Duration)
}

func TestUp_ARunOver60MinutesIsMarkedAndFinishes(t *testing.T) {
	for name, tc := range map[string]struct {
		each time.Duration
		over bool
	}{
		// Three phases, each an apply and a verify.
		"over":       {each: 15 * time.Minute, over: true},
		"at exactly": {each: 10 * time.Minute, over: false},
	} {
		t.Run(name, func(t *testing.T) {
			clock := runtest.NewClock(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC))
			target := runtest.NewTarget(t, runtest.Drill)
			ps := phases(map[string][]runtest.PhaseOption{
				"0":  {runtest.Takes(clock, tc.each)},
				"1a": {runtest.Takes(clock, tc.each)},
				"1b": {runtest.Takes(clock, tc.each)},
			})

			require.NoError(t, run.Up(t.Context(), target, ps, run.WithClock(clock)))
			require.Len(t, target.Calls(), 6, "a run running long is never stopped")

			rec := lastRun(t, target)
			require.True(t, rec.Finished)
			require.Equal(t, 6*tc.each, rec.Total)
			require.Equal(t, tc.over, rec.Over)
		})
	}
}

func TestUp_RefusesATargetItCannotBePointedAt(t *testing.T) {
	for name, tc := range map[string]struct {
		target run.Target
		want   error
	}{
		"none":            {target: nil, want: run.ErrNoTarget},
		"with no name":    {target: runtest.NewTarget(t, runtest.Drill, runtest.Named("")), want: run.ErrNoTarget},
		"of no kind":      {target: runtest.NewTarget(t, run.Kind(0)), want: run.ErrUnknownKind},
		"of another kind": {target: runtest.NewTarget(t, run.Kind(7)), want: run.ErrUnknownKind},
	} {
		t.Run(name, func(t *testing.T) {
			require.ErrorIs(t, run.Up(t.Context(), tc.target, phases(nil)), tc.want)
			require.ErrorIs(t, run.Preflight(t.Context(), tc.target, phases(nil)), tc.want)
			require.ErrorIs(t, run.Down(t.Context(), tc.target, phases(nil), run.Confirm("x")), tc.want)
		})
	}
}

func TestUp_RefusesAnInvalidPhaseOrder(t *testing.T) {
	for name, ps := range map[string][]run.Phase{
		"with a nil phase":      {runtest.Phase("0"), nil},
		"with an unnamed phase": {runtest.Phase("0"), runtest.Phase("")},
		"with a name twice":     {runtest.Phase("0"), runtest.Phase("1a"), runtest.Phase("0")},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Drill)
			require.ErrorIs(t, run.Up(t.Context(), target, ps), run.ErrInvalidPhases)
			require.ErrorIs(t, run.Preflight(t.Context(), target, ps), run.ErrInvalidPhases)
			require.ErrorIs(t, run.Down(t.Context(), target, ps), run.ErrInvalidPhases)
			require.Empty(t, target.Calls())
		})
	}
}

func TestUp_HandsEveryPhaseTheSameRun(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	var seen []*run.Run
	ps := []run.Phase{spy{name: "0", seen: &seen}, spy{name: "1a", seen: &seen}}

	require.NoError(t, run.Up(t.Context(), target, ps))
	require.Len(t, seen, 4)
	for _, r := range seen {
		require.Same(t, seen[0], r)
	}
	require.Same(t, target, seen[0].Target)
	require.Equal(t, lastRun(t, target).ID, seen[0].ID)
}

// spy is a phase that keeps every Run it's handed.
type spy struct {
	name string
	seen *[]*run.Run
}

func (s spy) Name() string { return s.name }

func (s spy) Apply(_ context.Context, r *run.Run) error {
	*s.seen = append(*s.seen, r)
	return nil
}

func (s spy) Verify(_ context.Context, r *run.Run) error {
	*s.seen = append(*s.seen, r)
	return nil
}

func (s spy) Teardown(_ context.Context, r *run.Run) error {
	*s.seen = append(*s.seen, r)
	return nil
}
