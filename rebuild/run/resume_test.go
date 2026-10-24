package run_test

import (
	"context"
	"testing"
	"time"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

// The test from revision 1's Drill.
func TestUp_AResumedRunStartsAtThePhaseThatFailed(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	phases := []run.Phase{
		runtest.Phase("0"),
		runtest.Phase("1a", runtest.FailApplyOnce(errNoEdge)),
		runtest.Phase("1b"),
	}

	_, err := run.Up(t.Context(), target, phases, run.NoPause())
	require.ErrorIs(t, err, errNoEdge)

	_, err = run.Up(t.Context(), target, phases, run.NoPause(), run.Resume())
	require.NoError(t, err)
	require.Equal(t, []string{
		"0 apply", "0 verify", "1a apply",
		"0 verify", "1a apply", "1a verify", "1b apply", "1b verify",
	}, target.Calls())
}

func TestUp_StopsAtTheCheckpointAfterEachPhaseAndResumeRunsTheNextOne(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{
		"0":  {runtest.Makes("the state backend")},
		"1a": {runtest.Makes("the edge", "its address")},
	})

	sum, err := run.Up(t.Context(), target, ps)
	require.NoError(t, err)
	require.Equal(t, []string{"0 apply", "0 verify"}, target.Calls())
	rec := lastRun(t, target)
	require.False(t, rec.Finished, "a run paused at a checkpoint is unfinished")
	require.Equal(t, &run.Stop{Phase: "0", Step: run.StepCheckpoint}, rec.Stopped)
	require.Equal(t, rec.ID, sum.ID)
	require.Equal(t, "drill", sum.Target)
	require.Equal(t, rec.Stopped, sum.Stopped)
	require.False(t, sum.Finished)
	require.Len(t, sum.Phases, 1)
	require.Equal(t, "0", sum.Phases[0].Name)
	require.Equal(t, []string{"the state backend"}, sum.Phases[0].Made)
	id := rec.ID

	sum, err = run.Up(t.Context(), target, ps, run.Resume())
	require.NoError(t, err)
	require.Equal(t, []string{"0 verify", "1a apply", "1a verify"}, target.Calls()[2:])
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepCheckpoint}, sum.Stopped)
	require.Len(t, sum.Phases, 1, "the summary holds what this call ran")
	require.Equal(t, []string{"the edge", "its address"}, sum.Phases[0].Made)
	require.Equal(t, id, sum.ID, "a resumed run is the same run")

	sum, err = run.Up(t.Context(), target, ps, run.Resume())
	require.NoError(t, err)
	require.Equal(t, []string{"0 verify", "1a verify", "1b apply", "1b verify"}, target.Calls()[5:])
	require.True(t, sum.Finished)
	require.Nil(t, sum.Stopped)

	recs := records(t, target)
	require.Len(t, recs, 1, "the run kept one record across its three calls")
	require.Equal(t, id, recs[0].ID)
	require.True(t, recs[0].Finished)
	require.Len(t, recs[0].Phases, 3)
	require.Equal(t, []string{"the state backend"}, recs[0].Phases[0].Made)
	require.Equal(t, []string{"the edge", "its address"}, recs[0].Phases[1].Made)
	require.Empty(t, recs[0].Phases[2].Made)
}

func TestUp_WithNoPauseRunsOnWithoutStopping(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	sum, err := run.Up(t.Context(), target, phases(nil), run.NoPause())
	require.NoError(t, err)
	require.Len(t, target.Calls(), 6)
	require.True(t, sum.Finished)
	require.Len(t, sum.Phases, 3)
}

func TestUp_ThroughNamesTheLastPhaseWithOrWithoutPausing(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	sum, err := run.Up(t.Context(), target, phases(nil), run.Through("1a"))
	require.NoError(t, err)
	require.Equal(t, &run.Stop{Phase: "0", Step: run.StepCheckpoint}, sum.Stopped,
		"without NoPause, the run still stops at the checkpoint before Through's phase")

	sum, err = run.Up(t.Context(), target, phases(nil), run.Resume(), run.Through("1a"))
	require.NoError(t, err)
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepThrough}, sum.Stopped)

	sum, err = run.Up(t.Context(), target, phases(nil), run.Resume())
	require.NoError(t, err, "a run stopped where Through said is resumed like any other")
	require.True(t, sum.Finished)
	require.Equal(t, []string{
		"0 apply", "0 verify",
		"0 verify", "1a apply", "1a verify",
		"0 verify", "1a verify", "1b apply", "1b verify",
	}, target.Calls())
}

func TestUp_ResumeAfterAFailureAtEachPhaseInTurnAppliesThatPhaseAgainAndGoesOn(t *testing.T) {
	for i, failing := range order {
		for step, fail := range map[run.Step]func(error) runtest.PhaseOption{
			run.StepApply:  runtest.FailApplyOnce,
			run.StepVerify: runtest.FailVerifyOnce,
		} {
			t.Run(failing+" "+string(step), func(t *testing.T) {
				target := runtest.NewTarget(t, runtest.Drill)
				ps := phases(map[string][]runtest.PhaseOption{failing: {fail(errNoEdge)}})

				_, err := run.Up(t.Context(), target, ps, run.NoPause())
				require.ErrorIs(t, err, errNoEdge)
				failedAt := len(target.Calls())
				require.Equal(t, &run.Stop{Phase: failing, Step: step, Error: errNoEdge.Error()}, lastRun(t, target).Stopped)

				sum, err := run.Up(t.Context(), target, ps, run.NoPause(), run.Resume())
				require.NoError(t, err)

				var want []string
				for _, name := range order[:i] {
					want = append(want, name+" verify")
				}
				for _, name := range order[i:] {
					want = append(want, name+" apply", name+" verify")
				}
				require.Equal(t, want, target.Calls()[failedAt:])

				require.True(t, sum.Finished)
				require.Len(t, sum.Phases, len(order)-i)
				rec := lastRun(t, target)
				require.True(t, rec.Finished)
				require.Nil(t, rec.Stopped)
				require.Len(t, rec.Phases, len(order), "the phase applied again replaced its failed entry")
				for j, p := range rec.Phases {
					require.Equal(t, order[j], p.Name)
					require.True(t, p.Verified)
				}
				require.Len(t, records(t, target), 1)
			})
		}
	}
}

func TestUp_ResumeStopsWhenAnEarlierPhaseNoLongerVerifies(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	_, err := run.Up(t.Context(), target, phases(nil), run.Through("1a"), run.NoPause())
	require.NoError(t, err)
	before := lastRun(t, target)
	applied := len(target.Calls())

	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailVerify(errNoEdge)}})
	_, err = run.Up(t.Context(), target, ps, run.Resume())
	require.ErrorIs(t, err, run.ErrDrift)
	require.ErrorIs(t, err, errNoEdge)
	require.ErrorContains(t, err, "phase 1a no longer verifies")
	require.Equal(t, []string{"0 verify", "1a verify"}, target.Calls()[applied:], "nothing is applied on top of drift")
	require.Equal(t, before, lastRun(t, target), "a resume that found drift changed nothing")
}

func TestUp_WithoutResumeStopsOnAnUnfinishedRunAndSaysItCanBeResumed(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	_, err := run.Up(t.Context(), target, phases(nil))
	require.NoError(t, err)
	applied := len(target.Calls())

	_, err = run.Up(t.Context(), target, phases(nil))
	require.ErrorIs(t, err, run.ErrUnfinishedRun)
	require.ErrorContains(t, err, "paused at the checkpoint after phase 0")
	require.ErrorContains(t, err, "Resume it with up --resume, or tear the target down with down")
	require.Len(t, target.Calls(), applied)
	require.Len(t, records(t, target), 1, "no new run starts on top of an unfinished one")
}

func TestUp_WithResumeStopsWhenThereIsNoUnfinishedRun(t *testing.T) {
	for name, tc := range map[string]struct {
		target *runtest.Target
		want   string
	}{
		"with no state backend": {
			target: runtest.NewTarget(t, runtest.Drill, runtest.NoBackend()),
			want:   "drill holds no run",
		},
		"with no records": {
			target: runtest.NewTarget(t, runtest.Drill),
			want:   "drill holds no run",
		},
		"after a run that finished": {
			target: func() *runtest.Target {
				target := runtest.NewTarget(t, runtest.Drill)
				require.NoError(t, up(t.Context(), target, phases(nil), run.NoPause()))
				return target
			}(),
			want: "finished. Start a new one with up",
		},
	} {
		t.Run(name, func(t *testing.T) {
			before := len(tc.target.Calls())

			_, err := run.Up(t.Context(), tc.target, phases(nil), run.Resume())
			require.ErrorIs(t, err, run.ErrNoUnfinishedRun)
			require.ErrorContains(t, err, tc.want)
			require.Len(t, tc.target.Calls(), before)
		})
	}
}

func TestUp_ResumeOfARunThatVerifiedEveryPhaseFinishesIt(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, target.Store().Write(t.Context(), run.Record{
		ID: "earlier",
		Phases: []run.PhaseRecord{
			{Name: "0", Verified: true}, {Name: "1a", Verified: true}, {Name: "1b", Verified: true},
		},
	}))

	sum, err := run.Up(t.Context(), target, phases(nil), run.Resume())
	require.NoError(t, err)
	require.Equal(t, []string{"0 verify", "1a verify", "1b verify"}, target.Calls())
	require.True(t, sum.Finished)
	require.Empty(t, sum.Phases)
	require.True(t, lastRun(t, target).Finished)
}

func TestUp_ResumeThroughAPhaseTheRunHasVerifiedRunsNothing(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	_, err := run.Up(t.Context(), target, phases(nil), run.Through("1a"), run.NoPause())
	require.NoError(t, err)
	applied := len(target.Calls())

	_, err = run.Up(t.Context(), target, phases(nil), run.Resume(), run.Through("0"))
	require.ErrorContains(t, err, "has already verified phase 0")
	require.Equal(t, []string{"0 verify", "1a verify"}, target.Calls()[applied:])
	require.Equal(t, &run.Stop{Phase: "1a", Step: run.StepThrough}, lastRun(t, target).Stopped)
}

func TestUp_ResumeStopsWhenTheRunAppliedAPhaseNoLongerInTheOrder(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	require.NoError(t, target.Store().Write(t.Context(), run.Record{
		ID:     "earlier",
		Phases: []run.PhaseRecord{{Name: "0", Verified: true}, {Name: "1z"}},
	}))

	_, err := run.Up(t.Context(), target, phases(nil), run.Resume())
	require.ErrorIs(t, err, run.ErrUnknownPhase)
	require.ErrorContains(t, err, "phase 1z")
	require.Equal(t, []string{"0 verify"}, target.Calls())
}

func TestUp_ResumeStopsWhenTheTargetsRecordsCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailStore(errBroken))

	_, err := run.Up(t.Context(), target, phases(nil), run.Resume())
	require.ErrorIs(t, err, errBroken)
	require.Empty(t, target.Calls())
}

func TestUp_ResumeStopsWhenItsRecordCannotBeWritten(t *testing.T) {
	// The first call writes twice: when the run starts, and at the checkpoint.
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailStoreWrites(2, errBroken))
	_, err := run.Up(t.Context(), target, phases(nil))
	require.NoError(t, err)

	_, err = run.Up(t.Context(), target, phases(nil), run.Resume())
	require.ErrorIs(t, err, errBroken)
	require.ErrorContains(t, err, "writing run")
	require.Equal(t, []string{"0 apply", "0 verify", "0 verify"}, target.Calls(), "no phase is applied unrecorded")
}

func TestUp_ARunCanResumeFromADifferentCommitAndEachPhaseKeepsItsOwn(t *testing.T) {
	first := run.Source{Repo: run.Commit{SHA: "a"}, Values: run.Commit{SHA: "v"}}
	fix := run.Source{Repo: run.Commit{SHA: "b"}, Values: run.Commit{SHA: "v"}}
	target := runtest.NewTarget(t, runtest.Live, runtest.WithSource(first))
	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailApplyOnce(errNoEdge)}})

	_, err := run.Up(t.Context(), target, ps, run.NoPause())
	require.ErrorIs(t, err, errNoEdge)

	target.SetSource(fix)
	_, err = run.Up(t.Context(), target, ps, run.NoPause(), run.Resume())
	require.NoError(t, err)

	rec := lastRun(t, target)
	require.True(t, rec.Finished)
	require.Equal(t, first, rec.Phases[0].Source)
	require.Equal(t, fix, rec.Phases[1].Source)
	require.Equal(t, fix, rec.Phases[2].Source)
}

func TestUp_ALiveRunRefusesUncommittedChangesOnEveryCall(t *testing.T) {
	for name, src := range map[string]run.Source{
		"in this repository":     {Repo: run.Commit{SHA: "a", Uncommitted: true}, Values: run.Commit{SHA: "v"}},
		"in the target's values": {Repo: run.Commit{SHA: "a"}, Values: run.Commit{SHA: "v", Uncommitted: true}},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Live)
			_, err := run.Up(t.Context(), target, phases(nil))
			require.NoError(t, err)
			before := lastRun(t, target)
			applied := len(target.Calls())

			target.SetSource(src)
			_, err = run.Up(t.Context(), target, phases(nil), run.Resume())
			require.ErrorIs(t, err, run.ErrUncommitted)
			require.Len(t, target.Calls(), applied)
			require.Equal(t, before, lastRun(t, target))
		})
	}
}

func TestUp_ARunsTimeIncludesTheTimeBetweenCalls(t *testing.T) {
	start := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	clock := runtest.NewClock(start)
	target := runtest.NewTarget(t, runtest.Drill)
	// Every apply and every verify takes five minutes.
	ps := phases(map[string][]runtest.PhaseOption{
		"0":  {runtest.Takes(clock, 5*time.Minute)},
		"1a": {runtest.Takes(clock, 5*time.Minute), runtest.FailApplyOnce(errNoEdge)},
		"1b": {runtest.Takes(clock, 5*time.Minute)},
	})

	// 0 applies and verifies, and the run pauses at 10 minutes.
	sum, err := run.Up(t.Context(), target, ps, run.WithClock(clock))
	require.NoError(t, err)
	require.Equal(t, 10*time.Minute, sum.Total)
	require.Equal(t, 10*time.Minute, sum.Phases[0].Duration)

	// The operator looks for 20 minutes, then resumes: 0 verifies again, and
	// 1a's apply fails, at 40 minutes.
	clock.Advance(20 * time.Minute)
	sum, err = run.Up(t.Context(), target, ps, run.WithClock(clock), run.Resume())
	require.ErrorIs(t, err, errNoEdge)
	require.Equal(t, 40*time.Minute, sum.Total)
	require.Equal(t, 5*time.Minute, sum.Phases[0].Duration)
	require.False(t, sum.Over)

	// The operator takes 30 minutes over the failure, then resumes with no
	// more pauses: 0 verifies, then 1a and 1b each take ten minutes.
	clock.Advance(30 * time.Minute)
	sum, err = run.Up(t.Context(), target, ps, run.WithClock(clock), run.Resume(), run.NoPause())
	require.NoError(t, err)
	require.True(t, sum.Finished)
	require.Equal(t, 95*time.Minute, sum.Total)
	require.True(t, sum.Over, "the time between calls counts toward the 60 minutes")

	rec := lastRun(t, target)
	require.Equal(t, start, rec.Started, "a run's time starts at its first up")
	require.Equal(t, 95*time.Minute, rec.Total)
	require.True(t, rec.Over)
}

func TestUp_APhaseReportsWhatItMadeThroughRunAndOnlyFromItsOwnCall(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := []run.Phase{reporting{name: "0"}, reporting{name: "1a"}}

	_, err := run.Up(t.Context(), target, ps)
	require.NoError(t, err)
	sum, err := run.Up(t.Context(), target, ps, run.Resume())
	require.NoError(t, err)

	require.Len(t, sum.Phases, 1)
	require.Equal(t, []string{"1a applied", "1a verified"}, sum.Phases[0].Made,
		"what phase 0 reported while the resume verified it is not phase 1a's")
	rec := lastRun(t, target)
	require.Equal(t, []string{"0 applied", "0 verified"}, rec.Phases[0].Made,
		"the resume's verify of phase 0 did not change what phase 0 made")
}

// reporting is a phase that reports through the Run from its Apply and its
// Verify.
type reporting struct{ name string }

func (p reporting) Name() string { return p.name }

func (p reporting) Apply(_ context.Context, r *run.Run) error {
	r.Made(p.name + " applied")
	return nil
}

func (p reporting) Verify(_ context.Context, r *run.Run) error {
	r.Made(p.name + " verified")
	return nil
}

func (p reporting) Teardown(context.Context, *run.Run) error { return nil }

func TestUp_AResumedRunHandsItsPhasesTheRunsID(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	var seen []*run.Run
	ps := []run.Phase{spy{name: "0", seen: &seen}, spy{name: "1a", seen: &seen}}

	_, err := run.Up(t.Context(), target, ps)
	require.NoError(t, err)
	_, err = run.Up(t.Context(), target, ps, run.Resume())
	require.NoError(t, err)

	require.Len(t, seen, 5)
	for _, r := range seen {
		require.Equal(t, lastRun(t, target).ID, r.ID)
	}
}
