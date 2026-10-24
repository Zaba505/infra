package run

import (
	"context"
	"fmt"
)

// unfinished reads the target's unfinished run. It returns the run's record
// with the index in phases of the first phase the run hasn't verified, which
// is len(phases) when it has verified them all.
//
// It changes nothing and verifies nothing. It stops with ErrNoUnfinishedRun
// when there is no run to resume, and with ErrUnknownPhase when the run
// applied a phase that is no longer in the order.
//
// A run whose record says it is still going, with no Stopped, is resumed like
// any other: that is also what a run whose up was killed partway leaves, and
// there is no telling the two apart from the record. The operator decides
// whether a target is fit to resume.
func unfinished(ctx context.Context, target Target, phases []Phase) (Record, int, error) {
	name := target.Name()
	last, ok, err := LastRun(ctx, target.Store())
	if err != nil {
		return Record{}, 0, fmt.Errorf("resume: %w", err)
	}
	if !ok {
		return Record{}, 0, fmt.Errorf("resume: %w: %s holds no run. Start one with up", ErrNoUnfinishedRun, name)
	}
	if last.Finished {
		return Record{}, 0, fmt.Errorf("resume: %w: %s's last run, %s, finished. Start a new one with up",
			ErrNoUnfinishedRun, name, last.ID)
	}

	verified := make(map[string]bool, len(last.Phases))
	for _, applied := range last.Phases {
		if indexOf(phases, applied.Name) < 0 {
			return Record{}, 0, fmt.Errorf("resume: %w: run %s applied phase %s, which is not in the phase order",
				ErrUnknownPhase, last.ID, applied.Name)
		}
		verified[applied.Name] = applied.Verified
	}
	start := 0
	for start < len(phases) && verified[phases[start].Name()] {
		start++
	}
	return last, start, nil
}

// resume makes r the run rec records: it hands r the run's ID, and verifies
// every phase the run has verified so far, in the order it applied them. It
// stops with ErrDrift when one no longer verifies, changing nothing.
func resume(ctx context.Context, r *Run, phases []Phase, rec Record) error {
	r.ID = rec.ID
	for _, applied := range rec.Phases {
		if !applied.Verified {
			continue
		}
		p := phases[indexOf(phases, applied.Name)]
		if err := p.Verify(ctx, r); err != nil {
			return fmt.Errorf("resume: %w: phase %s no longer verifies: %w. Tear the target down with down",
				ErrDrift, applied.Name, err)
		}
	}
	return nil
}
