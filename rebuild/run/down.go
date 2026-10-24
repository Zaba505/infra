package run

import (
	"context"
	"fmt"
)

// A DownOption changes how Down runs.
type DownOption func(*downOptions)

type downOptions struct {
	confirm string
}

// Confirm confirms Down with the target's name. Down against the live target
// refuses without it. Against a drill it isn't needed, but a name that isn't
// the target's is still refused, since it means the operator meant another
// target.
func Confirm(name string) DownOption {
	return func(o *downOptions) { o.confirm = name }
}

// Down tears down every phase against target in reverse order, so the first
// phase, phase 0, is torn down last. It stops at the first teardown that fails
// and returns its error, leaving the phases before it in place: phase 0 holds
// the state backend, and removing it while later phases still stand would lose
// what they are.
//
// Down keeps no record of its own. Phase 0's teardown removes the state
// backend and the records in it, so a Down that finishes leaves the target as
// a rebuild from nothing; see Store.
//
// Every Teardown succeeds when there's nothing left to remove, so Down
// succeeds against a target with nothing on it, including one Apply never ran
// against, and it's safe to run again after it fails partway.
//
// Against the live target, Down refuses unless Confirm names the target, and
// refuses uncommitted changes, as Up does.
func Down(ctx context.Context, target Target, phases []Phase, opts ...DownOption) error {
	var o downOptions
	for _, opt := range opts {
		opt(&o)
	}
	if err := checkTarget(target); err != nil {
		return err
	}
	if err := checkPhases(phases); err != nil {
		return err
	}
	switch {
	case o.confirm != "" && o.confirm != target.Name():
		return fmt.Errorf("%w: confirmed with %q, but the target is %q", ErrNotConfirmed, o.confirm, target.Name())
	case o.confirm == "" && target.Kind() == Live:
		return fmt.Errorf("%w: %s is the live target; confirm with %q to tear it down", ErrNotConfirmed, target.Name(), target.Name())
	}

	r, err := newRun(ctx, target, true)
	if err != nil {
		return err
	}
	for i := len(phases) - 1; i >= 0; i-- {
		p := phases[i]
		if err := p.Teardown(ctx, r); err != nil {
			return fmt.Errorf("phase %s: teardown: %w", p.Name(), err)
		}
	}
	return nil
}
