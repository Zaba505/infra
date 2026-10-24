package run

import (
	"context"
	"fmt"
)

// Preflight runs the preflight check against target, as up does before it
// applies anything. It changes nothing.
//
//  1. If the target has no state backend, this is a rebuild from nothing, and
//     the check passes.
//  2. If the target's last run is unfinished, the target holds a partial
//     rebuild, and the check stops with ErrUnfinishedRun. The operator tears
//     the target down first.
//  3. Otherwise it runs Verify for every phase the last run applied. Any
//     failure is drift, and the check stops with ErrDrift.
func Preflight(ctx context.Context, target Target, phases []Phase) error {
	if err := checkTarget(target); err != nil {
		return err
	}
	if err := checkPhases(phases); err != nil {
		return err
	}
	r, err := newRun(ctx, target, false)
	if err != nil {
		return err
	}
	return preflight(ctx, r, phases)
}

func preflight(ctx context.Context, r *Run, phases []Phase) error {
	last, ok, err := LastRun(ctx, r.Target.Store())
	if err != nil {
		return fmt.Errorf("preflight: %w", err)
	}
	if !ok {
		return nil
	}
	if !last.Finished {
		return fmt.Errorf("preflight: %w: run %s%s. Tear the target down with down first",
			ErrUnfinishedRun, last.ID, describeStop(last.Stopped))
	}
	for _, applied := range last.Phases {
		i := indexOf(phases, applied.Name)
		if i < 0 {
			return fmt.Errorf("preflight: %w: the last run applied phase %s, which is not in the phase order, so it can't be verified",
				ErrUnknownPhase, applied.Name)
		}
		if err := phases[i].Verify(ctx, r); err != nil {
			return fmt.Errorf("preflight: %w: phase %s: %w", ErrDrift, applied.Name, err)
		}
	}
	return nil
}

// describeStop says where an unfinished run stopped, for the preflight
// check's error.
func describeStop(s *Stop) string {
	switch {
	case s == nil:
		return " never recorded where it stopped"
	case s.Step == StepThrough:
		return fmt.Sprintf(" stopped after phase %s, as it was asked to", s.Phase)
	default:
		return fmt.Sprintf(" stopped at phase %s's %s", s.Phase, s.Step)
	}
}
