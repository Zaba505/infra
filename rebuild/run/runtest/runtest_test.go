package runtest_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

var errBroken = errors.New("broken")

func TestPhase_FailsOnceThenPasses(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	r := &run.Run{Target: target}
	p := runtest.Phase("0",
		runtest.FailApplyOnce(errBroken),
		runtest.FailVerifyOnce(errBroken),
		runtest.FailTeardownOnce(errBroken),
	)
	ctx := t.Context()

	require.ErrorIs(t, p.Apply(ctx, r), errBroken)
	require.NoError(t, p.Apply(ctx, r))
	require.ErrorIs(t, p.Verify(ctx, r), errBroken)
	require.NoError(t, p.Verify(ctx, r))
	require.ErrorIs(t, p.Teardown(ctx, r), errBroken)
	require.Equal(t, []string{"0"}, target.Left(), "a teardown that failed removed nothing")
	require.NoError(t, p.Teardown(ctx, r))
	require.Empty(t, target.Left())
	require.Equal(t, []string{
		"0 apply", "0 apply", "0 verify", "0 verify", "0 teardown", "0 teardown",
	}, target.Calls())
}

func TestPhase_FailsEveryTime(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	r := &run.Run{Target: target}
	p := runtest.Phase("0", runtest.FailApply(errBroken), runtest.FailVerify(errBroken))

	for range 3 {
		require.ErrorIs(t, p.Apply(t.Context(), r), errBroken)
		require.ErrorIs(t, p.Verify(t.Context(), r), errBroken)
	}
	require.Equal(t, []string{"0"}, target.Left(), "an apply that failed may have made something")
}

func TestPhase_TakesTimeOnApplyAndVerifyButNotTeardown(t *testing.T) {
	start := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	clock := runtest.NewClock(start)
	r := &run.Run{Target: runtest.NewTarget(t, runtest.Drill)}
	p := runtest.Phase("0", runtest.Takes(clock, time.Minute))

	require.NoError(t, p.Apply(t.Context(), r))
	require.NoError(t, p.Verify(t.Context(), r))
	require.NoError(t, p.Teardown(t.Context(), r))
	require.Equal(t, start.Add(2*time.Minute), clock.Now())
}

func TestPhase_RefusesATargetNewTargetDidNotMake(t *testing.T) {
	r := &run.Run{Target: other{}}
	p := runtest.Phase("0")

	require.ErrorIs(t, p.Apply(t.Context(), r), runtest.ErrNotARuntestTarget)
	require.ErrorIs(t, p.Verify(t.Context(), r), runtest.ErrNotARuntestTarget)
	require.ErrorIs(t, p.Teardown(t.Context(), r), runtest.ErrNotARuntestTarget)
}

func TestNewTarget_IsNamedAfterItsKindAndHasCleanCommits(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Live)

	require.Equal(t, "live", target.Name())
	require.Equal(t, run.Live, target.Kind())
	src, err := target.Source(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, src.Repo.SHA)
	require.NotEmpty(t, src.Values.SHA)
	require.False(t, src.Uncommitted())
}

func TestStore_ReplacesARecordWithTheSameID(t *testing.T) {
	store := runtest.NewTarget(t, runtest.Drill).Store()
	ctx := t.Context()

	require.NoError(t, store.Write(ctx, run.Record{ID: "a"}))
	require.NoError(t, store.Write(ctx, run.Record{ID: "b"}))
	require.NoError(t, store.Write(ctx, run.Record{ID: "a", Finished: true}))

	recs, err := store.Records(ctx)
	require.NoError(t, err)
	require.Equal(t, []run.Record{{ID: "a", Finished: true}, {ID: "b"}}, recs)
}

func TestStore_SharesNothingWithItsCallers(t *testing.T) {
	store := runtest.NewTarget(t, runtest.Drill).Store()
	ctx := t.Context()
	rec := run.Record{
		ID:      "a",
		Phases:  []run.PhaseRecord{{Name: "0"}},
		Stopped: &run.Stop{Phase: "0"},
	}
	require.NoError(t, store.Write(ctx, rec))
	rec.Phases[0].Name = "changed"
	rec.Stopped.Phase = "changed"

	recs, err := store.Records(ctx)
	require.NoError(t, err)
	require.Equal(t, "0", recs[0].Phases[0].Name)
	require.Equal(t, "0", recs[0].Stopped.Phase)
}

func TestStore_WithNoBackendKeepsNothing(t *testing.T) {
	store := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend()).Store()

	require.ErrorIs(t, store.Write(t.Context(), run.Record{ID: "a"}), run.ErrNoBackend)
	_, err := store.Records(t.Context())
	require.ErrorIs(t, err, run.ErrNoBackend)
}

// other is a target runtest didn't make.
type other struct{}

func (other) Name() string                               { return "other" }
func (other) Kind() run.Kind                             { return run.Drill }
func (other) Source(context.Context) (run.Source, error) { return run.Source{}, nil }
func (other) Store() run.Store                           { return nil }
