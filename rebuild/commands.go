package main

import (
	"context"
	"fmt"

	"dagger/rebuild/run"
)

// Up runs the rebuild against the target: the preflight check, then each
// phase's apply and verify, in order. After each phase verifies, it writes the
// run's record and stops at a checkpoint, and returns the run's summary: what
// the phase made, how long it took and the run's time so far. Continue the run
// with --resume, or tear the target down with down.
//
// A failed apply or verify stops the run at once, and fails the call with what
// failed and the summary. The target is left as it is to look at; --resume
// applies the failed phase again, and down is always the way back to nothing.
//
// Against the live target, it refuses uncommitted changes, in this repository
// or in the values checkout.
//
// +cache="never"
func (m *Rebuild) Up(
	ctx context.Context,
	// The target to run against: the name of a file in the values checkout's
	// targets/ directory, without .yaml. There is no default.
	target string,
	// Continue the target's unfinished run instead of starting a new one,
	// whether it paused at a checkpoint, stopped where --through said, or
	// failed. It verifies every phase the run has verified so far, then
	// continues from the first one it hasn't.
	//
	// +optional
	resume bool,
	// The last phase to run, such as 0 or 1a. Without it, up runs through the
	// last phase in the order.
	//
	// +optional
	through string,
	// Run on from one phase to the next without stopping at the checkpoint
	// between them, as CI and drills do. Each phase's verify is then the only
	// gate.
	//
	// +optional
	noPause bool,
) (string, error) {
	t, err := m.target(ctx, target)
	if err != nil {
		return "", err
	}
	var opts []run.UpOption
	if resume {
		opts = append(opts, run.Resume())
	}
	if through != "" {
		opts = append(opts, run.Through(through))
	}
	if noPause {
		opts = append(opts, run.NoPause())
	}
	summary, err := run.Up(ctx, t, m.phases(), opts...)
	if err != nil {
		if summary.ID == "" {
			return "", err
		}
		return "", fmt.Errorf("%w\n\n%s", err, summary)
	}
	return summary.String(), nil
}

// Down tears down every phase against the target, in reverse order, so phase
// 0 is torn down last. It stops at the first teardown that fails, leaving the
// phases before it in place. It's safe to run again after it fails partway,
// and it succeeds when there's nothing left to remove.
//
// Against the live target it refuses unless --confirm names the target, and
// refuses uncommitted changes, as up does.
//
// +cache="never"
func (m *Rebuild) Down(
	ctx context.Context,
	// The target to tear down: the name of a file in the values checkout's
	// targets/ directory, without .yaml. There is no default.
	target string,
	// The target's name, to confirm the teardown. The live target refuses
	// without it. A drill doesn't need it, but refuses a name that isn't its
	// own.
	//
	// +optional
	confirm string,
) (string, error) {
	t, err := m.target(ctx, target)
	if err != nil {
		return "", err
	}
	var opts []run.DownOption
	if confirm != "" {
		opts = append(opts, run.Confirm(confirm))
	}
	if err := run.Down(ctx, t, m.phases(), opts...); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s is torn down: every phase's teardown succeeded.\n", target), nil
}

// Preflight runs only the preflight check against the target, as up does
// before it applies anything, and changes nothing. It passes against a target
// with no state backend, as a rebuild from nothing. It fails when the target's
// last run is unfinished, and when a phase that run applied no longer
// verifies, which is drift.
//
// +cache="never"
func (m *Rebuild) Preflight(
	ctx context.Context,
	// The target to check: the name of a file in the values checkout's
	// targets/ directory, without .yaml. There is no default.
	target string,
) (string, error) {
	t, err := m.target(ctx, target)
	if err != nil {
		return "", err
	}
	if err := run.Preflight(ctx, t, m.phases()); err != nil {
		return "", err
	}
	return fmt.Sprintf("the preflight check passed against %s.\n", target), nil
}
