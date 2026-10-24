package run

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// TimeTarget is how long a rebuild is meant to take. A run over it is marked
// in its record, and never stopped: a live rebuild that's running long is
// still the fastest way back.
const TimeTarget = 60 * time.Minute

// ErrNoBackend is what a Store returns when the target has no state backend
// to keep records in, as before phase 0 has made one.
var ErrNoBackend = errors.New("the target has no state backend")

// A Store keeps a target's run records, in the target's state backend. Phase
// 0's proposal provides the real one; the tests use runtest's in-memory one.
type Store interface {
	// Records returns every record the target holds, oldest first. It
	// returns ErrNoBackend when the target has no state backend.
	Records(ctx context.Context) ([]Record, error)

	// Write saves rec, replacing the record with the same ID if there is
	// one, and otherwise adding it after every other. It returns
	// ErrNoBackend when the target has no state backend.
	Write(ctx context.Context, rec Record) error
}

// Record is what a run leaves in its target's state backend. Up writes it when
// the run starts, after every phase, and when the run stops.
type Record struct {
	// ID is the run's ID.
	ID string `json:"id"`

	// Target is the name of the target the run was pointed at.
	Target string `json:"target"`

	// Started is when the run's up started, before the preflight check.
	Started time.Time `json:"started"`

	// Phases holds each phase the run applied, in the order it applied them.
	Phases []PhaseRecord `json:"phases"`

	// Finished is true once the run has verified the last phase in the
	// order. A run that is not finished is unfinished, whether it failed,
	// stopped where Through said, or is still going.
	Finished bool `json:"finished"`

	// Stopped says where an unfinished run stopped. It is nil for a run that
	// finished, and for one that is still going.
	Stopped *Stop `json:"stopped,omitempty"`

	// Total is the run's time: from Started to the last phase's verify,
	// including the preflight check.
	Total time.Duration `json:"total"`

	// Over is true when Total is over TimeTarget.
	Over bool `json:"over"`
}

// PhaseRecord is one phase a run applied.
type PhaseRecord struct {
	// Name is the phase's name.
	Name string `json:"name"`

	// Source is what the phase was applied from: the commits of this
	// repository and of the target's values, and whether either had
	// uncommitted changes.
	Source Source `json:"source"`

	// Duration is how long the phase took, from the start of its apply to
	// the end of its verify.
	Duration time.Duration `json:"duration"`

	// Verified is true when the phase's verify passed.
	Verified bool `json:"verified"`
}

// Step is the part of a phase a run stopped at.
type Step string

const (
	// StepApply is a phase's Apply.
	StepApply Step = "apply"

	// StepVerify is a phase's Verify.
	StepVerify Step = "verify"

	// StepThrough is the phase Through named: the run verified it and
	// stopped there, as it was asked to.
	StepThrough Step = "through"
)

// Stop is where an unfinished run stopped.
type Stop struct {
	// Phase is the name of the phase it stopped at.
	Phase string `json:"phase"`

	// Step is the part of the phase it stopped at.
	Step Step `json:"step"`

	// Error is what failed, for a run that stopped on a failure.
	Error string `json:"error,omitempty"`
}

// LastRun returns the target's most recent record. It reports false when the
// target holds none, including when it has no state backend.
func LastRun(ctx context.Context, store Store) (Record, bool, error) {
	recs, err := records(ctx, store)
	if err != nil || len(recs) == 0 {
		return Record{}, false, err
	}
	return recs[len(recs)-1], true, nil
}

// LastKnownGood returns the record of the target's last run that finished,
// which is the target's last-known-good. It reports false when no run has
// finished, including when the target has no state backend.
func LastKnownGood(ctx context.Context, store Store) (Record, bool, error) {
	recs, err := records(ctx, store)
	if err != nil {
		return Record{}, false, err
	}
	for i := len(recs) - 1; i >= 0; i-- {
		if recs[i].Finished {
			return recs[i], true, nil
		}
	}
	return Record{}, false, nil
}

// records reads every record in store, treating a target with no state
// backend as one with no records.
func records(ctx context.Context, store Store) ([]Record, error) {
	recs, err := store.Records(ctx)
	if errors.Is(err, ErrNoBackend) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading the target's records: %w", err)
	}
	return recs, nil
}

// write saves rec. A target with no state backend has nowhere to keep it,
// which is the case until phase 0 has made one, so that is not a failure:
// the target then holds nothing a later run would need explained.
func write(ctx context.Context, store Store, rec Record) error {
	err := store.Write(ctx, rec)
	if err == nil || errors.Is(err, ErrNoBackend) {
		return nil
	}
	return fmt.Errorf("writing run %s's record: %w", rec.ID, err)
}
