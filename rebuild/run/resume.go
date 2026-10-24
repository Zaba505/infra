package run

import (
	"context"
	"fmt"
)

// resume reads the target's unfinished run and makes r that run: it hands r
// the run's ID, verifies every phase the run has verified so far, and returns
// the run's record with the index in phases of the first phase the run hasn't
// verified. That is len(phases) when it has verified them all.
//
// It changes nothing. It stops with ErrNoUnfinishedRun when there is no run to
// resume, with ErrUnknownPhase when the run applied a phase that is no longer
// in the order, and with ErrDrift when a phase it verified no longer does.
func resume(ctx context.Context, r *Run, phases []Phase) (Record, int, error) {
	name := r.Target.Name()
	last, ok, err := LastRun(ctx, r.Target.Store())
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

	r.ID = last.ID
	verified := make(map[string]bool, len(last.Phases))
	for _, applied := range last.Phases {
		i := indexOf(phases, applied.Name)
		if i < 0 {
			return Record{}, 0, fmt.Errorf("resume: %w: run %s applied phase %s, which is not in the phase order",
				ErrUnknownPhase, last.ID, applied.Name)
		}
		if !applied.Verified {
			continue
		}
		if err := phases[i].Verify(ctx, r); err != nil {
			return Record{}, 0, fmt.Errorf("resume: %w: phase %s no longer verifies: %w. Tear the target down with down",
				ErrDrift, applied.Name, err)
		}
		verified[applied.Name] = true
	}

	start := 0
	for start < len(phases) && verified[phases[start].Name()] {
		start++
	}
	return last, start, nil
}
