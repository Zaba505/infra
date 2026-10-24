package run

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// An UpOption changes how Up runs.
type UpOption func(*upOptions)

type upOptions struct {
	through string
	clock   Clock
}

// Through names the last phase Up runs. Without it, Up runs through the last
// phase in the order.
//
// A run that stops where Through said has not finished, because it hasn't
// verified the whole rebuild, so it is never the target's last-known-good. That
// makes it unfinished, and the preflight check refuses to start a new run on
// top of it: the next Up stops until the target is torn down. Through is for a
// drill of the first phases, not for building a target up in steps.
func Through(phase string) UpOption {
	return func(o *upOptions) { o.through = phase }
}

// WithClock gives Up the clock it measures the run's time with. The tests give
// it a fake one.
func WithClock(c Clock) UpOption {
	return func(o *upOptions) { o.clock = c }
}

// Up runs the rebuild against target: the preflight check, then each phase's
// Apply and then its Verify, in the order given, through the last phase or the
// one Through names.
//
// A failed Apply or Verify stops the run at once: Up writes the record and
// returns the error, and no later phase runs. The target is left as it is for
// the operator to look at, and Down is always the way back to nothing.
//
// Against the live target, Up refuses uncommitted changes, in this
// repository or in the target's values. A drill may run from them, and its
// record says so.
//
// The run's time is measured from the start of Up to the last phase's verify,
// including the preflight check. A run over TimeTarget is marked in its
// record, and never stopped.
func Up(ctx context.Context, target Target, phases []Phase, opts ...UpOption) error {
	o := upOptions{clock: systemClock{}}
	for _, opt := range opts {
		opt(&o)
	}
	started := o.clock.Now()

	if err := checkTarget(target); err != nil {
		return err
	}
	if err := checkPhases(phases); err != nil {
		return err
	}
	last := len(phases) - 1
	if o.through != "" {
		last = indexOf(phases, o.through)
		if last < 0 {
			return fmt.Errorf("through %s: %w", o.through, ErrUnknownPhase)
		}
	}

	r, err := newRun(ctx, target, true)
	if err != nil {
		return err
	}
	if err := preflight(ctx, r, phases); err != nil {
		return err
	}

	// The record is written whatever became of the run's context, so a run
	// cancelled partway still records where it stopped.
	wctx := context.WithoutCancel(ctx)
	store := target.Store()
	rec := Record{ID: r.ID, Target: target.Name(), Started: started}
	rec.measure(o.clock.Now())
	// Until the first phase, phase 0, has applied and verified, the target
	// may have no state backend to keep the record in.
	if err := write(wctx, store, rec, true); err != nil {
		return err
	}

	for i, p := range phases[:last+1] {
		begun := o.clock.Now()
		step := StepApply
		err := p.Apply(ctx, r)
		if err == nil {
			step = StepVerify
			err = p.Verify(ctx, r)
		}
		now := o.clock.Now()
		rec.Phases = append(rec.Phases, PhaseRecord{
			Name:     p.Name(),
			Source:   r.Source,
			Duration: now.Sub(begun),
			Verified: err == nil,
		})
		rec.measure(now)
		if err != nil {
			rec.Stopped = &Stop{Phase: p.Name(), Step: step, Error: err.Error()}
			return errors.Join(
				fmt.Errorf("phase %s: %s: %w", p.Name(), step, err),
				write(wctx, store, rec, i == 0),
			)
		}
		if err := write(wctx, store, rec, false); err != nil {
			return err
		}
	}

	if last < len(phases)-1 {
		rec.Stopped = &Stop{Phase: phases[last].Name(), Step: StepThrough}
	} else {
		rec.Finished = true
	}
	return write(wctx, store, rec, last < 0)
}

// measure sets the run's time to now, and marks it when it's over the target.
func (rec *Record) measure(now time.Time) {
	rec.Total = now.Sub(rec.Started)
	rec.Over = rec.Total > TimeTarget
}
