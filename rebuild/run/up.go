package run

import (
	"context"
	"errors"
	"fmt"
)

// An UpOption changes how Up runs.
type UpOption func(*upOptions)

type upOptions struct {
	through string
	noPause bool
	resume  bool
	clock   Clock
}

// Through names the last phase Up runs. Without it, Up runs through the last
// phase in the order.
//
// A run that stops where Through said has not finished, because it hasn't
// verified the whole rebuild, so it is never the target's last-known-good. It
// is unfinished, as a run paused at a checkpoint is: the next Up refuses to
// start a new run on top of it, and Up with Resume continues it.
func Through(phase string) UpOption {
	return func(o *upOptions) { o.through = phase }
}

// NoPause makes Up run on from one phase to the next without stopping at the
// checkpoint between them, as CI and drills do. Each phase's Verify is then
// the only gate. Through still names the last phase it runs.
func NoPause() UpOption {
	return func(o *upOptions) { o.noPause = true }
}

// Resume makes Up continue the target's unfinished run instead of starting a
// new one, whether that run paused at a checkpoint, stopped where Through
// said, or failed.
//
// Whether the target is fit to resume is the operator's call. A failure can
// leave a phase in a state its Apply can't fix; resuming then fails again, and
// Down followed by a fresh Up is the way back.
func Resume() UpOption {
	return func(o *upOptions) { o.resume = true }
}

// WithClock gives Up the clock it measures the run's time with. The tests give
// it a fake one.
func WithClock(c Clock) UpOption {
	return func(o *upOptions) { o.clock = c }
}

// Up runs the rebuild against target, one call at a time.
//
// Without Resume, Up starts a new run: the preflight check, then the first
// phase's Apply and then its Verify. With Resume, it continues the target's
// unfinished run:
//
//  1. it stops with ErrNoUnfinishedRun unless the target's last run is
//     unfinished;
//  2. it runs Verify for every phase the run has verified so far, and stops
//     with ErrDrift if one fails, changing nothing;
//  3. it continues from the first phase in the order that the run hasn't
//     verified, starting with that phase's Apply.
//
// After each phase verifies, Up writes the run's record, returns a Summary of
// what the phase made, how long it took and the run's time so far, and stops
// at the checkpoint after the phase. The run is then unfinished, and the
// operator continues it with Resume or tears the target down with Down. With
// NoPause, Up runs on to the next phase instead. Either way, it stops after
// the phase Through names, and the run finishes once the last phase in the
// order verifies.
//
// A failed Apply or Verify stops the run at once: Up writes the record and
// returns the error, with the Summary, and no later phase runs. The target is
// left as it is for the operator to look at. Resume applies the failed phase
// again, and Down is always the way back to nothing.
//
// Against the live target, every call refuses uncommitted changes, in this
// repository or in the target's values. A drill may run from them, and its
// record says so. A resumed run may run from a different commit than the one
// it started from; the record keeps the commits each phase was applied from.
//
// The run's time is measured from the start of its first Up to the last
// phase's verify. It includes the preflight check, the time between calls at
// every checkpoint, and the time between a failure and its resume. A run over
// TimeTarget is marked in its record, and never stopped.
func Up(ctx context.Context, target Target, phases []Phase, opts ...UpOption) (Summary, error) {
	o := upOptions{clock: systemClock{}}
	for _, opt := range opts {
		opt(&o)
	}
	started := o.clock.Now()

	if err := checkTarget(target); err != nil {
		return Summary{}, err
	}
	if err := checkPhases(phases); err != nil {
		return Summary{}, err
	}
	last := len(phases) - 1
	if o.through != "" {
		last = indexOf(phases, o.through)
		if last < 0 {
			return Summary{}, fmt.Errorf("through %s: %w", o.through, ErrUnknownPhase)
		}
	}

	r, err := newRun(ctx, target, true)
	if err != nil {
		return Summary{}, err
	}
	var (
		rec   Record
		start int
	)
	if o.resume {
		rec, start, err = resume(ctx, r, phases)
		if err != nil {
			return Summary{}, err
		}
		if o.through != "" && start > last {
			return Summary{}, fmt.Errorf("through %s: run %s has already verified phase %s", o.through, rec.ID, o.through)
		}
		rec.Stopped = nil
	} else {
		if err := preflight(ctx, r, phases); err != nil {
			return Summary{}, err
		}
		rec = Record{ID: r.ID, Target: target.Name(), Started: started}
	}

	// The record is written whatever became of the run's context, so a run
	// cancelled partway still records where it stopped.
	wctx := context.WithoutCancel(ctx)
	store := target.Store()
	rec.measure(o.clock.Now())
	// Until the first phase, phase 0, has applied and verified, the target
	// may have no state backend to keep the record in.
	if err := write(wctx, store, rec, start == 0); err != nil {
		return Summary{}, err
	}

	var ran []PhaseRecord
	for i := start; i <= last; i++ {
		p := phases[i]
		// What was reported before the phase began is not the phase's.
		r.takeMade()
		begun := o.clock.Now()
		step := StepApply
		err := p.Apply(ctx, r)
		if err == nil {
			step = StepVerify
			err = p.Verify(ctx, r)
		}
		now := o.clock.Now()
		pr := PhaseRecord{
			Name:     p.Name(),
			Source:   r.Source,
			Made:     r.takeMade(),
			Duration: now.Sub(begun),
			Verified: err == nil,
		}
		rec.put(pr)
		ran = append(ran, pr)
		rec.measure(now)
		if err != nil {
			rec.Stopped = &Stop{Phase: p.Name(), Step: step, Error: err.Error()}
			return summarize(rec, ran), errors.Join(
				fmt.Errorf("phase %s: %s: %w", p.Name(), step, err),
				write(wctx, store, rec, i == 0),
			)
		}

		switch {
		case i == len(phases)-1:
			rec.Finished = true
		case i == last:
			rec.Stopped = &Stop{Phase: p.Name(), Step: StepThrough}
		case !o.noPause:
			rec.Stopped = &Stop{Phase: p.Name(), Step: StepCheckpoint}
		}
		if err := write(wctx, store, rec, false); err != nil {
			return summarize(rec, ran), err
		}
		if rec.Finished || rec.Stopped != nil {
			return summarize(rec, ran), nil
		}
	}

	// Only a run with no phase left to run gets here: one with an empty
	// order, or a resumed one that had verified every phase in it.
	rec.Finished = true
	return summarize(rec, ran), write(wctx, store, rec, start == 0)
}
