// Package run is the rebuild's phase model: the run loop that applies and
// verifies phases in order, the teardown that unwinds them in reverse, the
// preflight check, the targets a run is pointed at, its records and its time.
//
// It is plain Go and imports nothing from Dagger; .golangci.yml keeps it that
// way. The rebuild module builds each phase with the containers, secrets and
// directories it needs, so nothing here ever sees an engine, and its tests
// drive it with the fakes in package runtest.
package run

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
)

// A Phase is one step of the rebuild. Every method runs against one target,
// and every method is safe to run again after it fails partway.
type Phase interface {
	// Name is the phase's place in the rebuild order, such as "0" or "1a".
	Name() string

	// Apply makes the target match this phase's definitions.
	Apply(ctx context.Context, r *Run) error

	// Verify fails unless the target is what this phase last applied, and
	// it works.
	Verify(ctx context.Context, r *Run) error

	// Teardown removes everything Apply made. It succeeds when there's
	// nothing left to remove, including when Apply never ran.
	Teardown(ctx context.Context, r *Run) error
}

// Run is what every phase method is handed: the target it runs against, the
// run's ID and the commits being rebuilt. A phase reads every value it needs
// from the target, and nothing from the environment it runs in.
type Run struct {
	// Target is the target the run is pointed at.
	Target Target

	// ID identifies the run. Its record carries the same ID.
	ID string

	// Source is the commit of this repository being rebuilt, and the commit
	// of the target's values, each with whether it had uncommitted changes.
	Source Source
}

// Source is what a run is rebuilt from.
type Source struct {
	// Repo is this repository's commit.
	Repo Commit `json:"repo"`

	// Values is the commit of the target's values.
	Values Commit `json:"values"`
}

// Commit is one checkout's commit, and whether the checkout had uncommitted
// changes on top of it.
type Commit struct {
	SHA         string `json:"sha"`
	Uncommitted bool   `json:"uncommitted"`
}

// Uncommitted reports whether either checkout had uncommitted changes.
func (s Source) Uncommitted() bool {
	return s.Repo.Uncommitted || s.Values.Uncommitted
}

// Kind says whether a target is the live platform or a drill. Its zero value
// is neither, so a target that never said which is refused rather than taken
// for one or the other.
type Kind int

const (
	// Live is the live platform. Only one target is live.
	Live Kind = iota + 1

	// Drill is a target a rebuild is practised against.
	Drill
)

// String returns "live" or "drill".
func (k Kind) String() string {
	switch k {
	case Live:
		return "live"
	case Drill:
		return "drill"
	default:
		return fmt.Sprintf("Kind(%d)", int(k))
	}
}

// A Target is the one place a run is pointed at. It holds every value a phase
// may read. Drill and live differ only in the target.
//
// There is no default target: Up, Down and Preflight each refuse a nil one.
type Target interface {
	// Name is the target's name. Down against the live target is confirmed
	// with it.
	Name() string

	// Kind says whether the target is live or a drill.
	Kind() Kind

	// Source reports the commits a run against this target is rebuilt from:
	// this repository's, and that of the target's values.
	Source(ctx context.Context) (Source, error)

	// Store is where the target's run records are kept, in its state
	// backend.
	Store() Store
}

// Errors a run stops with. Each is wrapped with what stopped it, so test for
// them with errors.Is.
var (
	// ErrNoTarget is returned when no target is given. There is no default.
	ErrNoTarget = errors.New("no target: a run needs one, and there is no default")

	// ErrUnknownKind is returned for a target that is neither live nor a
	// drill.
	ErrUnknownKind = errors.New("the target is neither live nor a drill")

	// ErrUncommitted is returned when a run against the live target would
	// run from uncommitted changes.
	ErrUncommitted = errors.New("a run against the live target refuses uncommitted changes")

	// ErrNotConfirmed is returned when down against the live target is not
	// confirmed with the target's name.
	ErrNotConfirmed = errors.New("down against the live target must be confirmed with the target's name")

	// ErrUnfinishedRun is returned by the preflight check when the target's
	// last run did not finish.
	ErrUnfinishedRun = errors.New("the target's last run is unfinished")

	// ErrDrift is returned by the preflight check when a phase the last run
	// applied no longer verifies.
	ErrDrift = errors.New("the target has drifted from what was last applied to it")

	// ErrUnknownPhase is returned when a phase is named that is not in the
	// phase order.
	ErrUnknownPhase = errors.New("no such phase")

	// ErrInvalidPhases is returned for a phase order that holds a nil phase,
	// a phase with no name, or two phases with the same name.
	ErrInvalidPhases = errors.New("invalid phase order")
)

// checkTarget refuses a target a run can't be pointed at.
func checkTarget(target Target) error {
	if target == nil {
		return ErrNoTarget
	}
	if target.Name() == "" {
		return fmt.Errorf("%w: the target has no name", ErrNoTarget)
	}
	switch target.Kind() {
	case Live, Drill:
		return nil
	default:
		return fmt.Errorf("target %s: %w", target.Name(), ErrUnknownKind)
	}
}

// checkPhases refuses a phase order with a nil phase, an unnamed phase or a
// name used twice.
func checkPhases(phases []Phase) error {
	seen := make(map[string]bool, len(phases))
	for i, p := range phases {
		if p == nil {
			return fmt.Errorf("%w: phase %d is nil", ErrInvalidPhases, i)
		}
		name := p.Name()
		if name == "" {
			return fmt.Errorf("%w: phase %d has no name", ErrInvalidPhases, i)
		}
		if seen[name] {
			return fmt.Errorf("%w: phase %s appears twice", ErrInvalidPhases, name)
		}
		seen[name] = true
	}
	return nil
}

// indexOf returns the index of the phase named name, or -1.
func indexOf(phases []Phase, name string) int {
	for i, p := range phases {
		if p.Name() == name {
			return i
		}
	}
	return -1
}

// newRun reads the target's source and hands back a Run with a fresh ID. A
// run that will change the live target refuses uncommitted changes, so every
// change to the live platform was made from history.
func newRun(ctx context.Context, target Target, changes bool) (*Run, error) {
	src, err := target.Source(ctx)
	if err != nil {
		return nil, fmt.Errorf("target %s: reading the commits to run from: %w", target.Name(), err)
	}
	if changes && target.Kind() == Live {
		if src.Repo.Uncommitted {
			return nil, fmt.Errorf("%w: this repository has uncommitted changes", ErrUncommitted)
		}
		if src.Values.Uncommitted {
			return nil, fmt.Errorf("%w: the target's values have uncommitted changes", ErrUncommitted)
		}
	}
	return &Run{Target: target, ID: rand.Text(), Source: src}, nil
}
