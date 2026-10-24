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
//
// The state backend is phase 0's: its Apply makes it and its Teardown removes
// it, records and all. So a Down that succeeds leaves no record behind, and
// the next Up's preflight check passes as a rebuild from nothing. A Down that
// fails partway leaves phase 0, and with it the unfinished record that keeps
// the next Up from building on what's left.
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
// the run starts, when a resume of it starts, after every phase, and when the
// run stops.
type Record struct {
	// ID is the run's ID. A resumed run keeps it.
	ID string `json:"id"`

	// Target is the name of the target the run was pointed at.
	Target string `json:"target"`

	// Started is when the run's first up started, before the preflight
	// check. Resuming the run doesn't move it.
	Started time.Time `json:"started"`

	// Phases holds each phase the run applied, in the order it applied them.
	// A phase applied again on a resume, after it failed, replaces the
	// phase's earlier entry, so the run holds one entry per phase, and it
	// holds the commits that phase was last applied from.
	Phases []PhaseRecord `json:"phases"`

	// Finished is true once the run has verified the last phase in the
	// order. A run that is not finished is unfinished, whether it failed,
	// paused at a checkpoint, stopped where Through said, or is still going.
	// Up with Resume continues an unfinished run.
	Finished bool `json:"finished"`

	// Stopped says where an unfinished run stopped. It is nil for a run that
	// finished, and for one that is still going.
	Stopped *Stop `json:"stopped,omitempty"`

	// Total is the run's time: from Started to the last phase's verify. It
	// includes the preflight check, the time between calls at every
	// checkpoint, and the time between a failure and its resume.
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
	// uncommitted changes. A resumed run may apply its phases from a
	// different commit than the one it started from, so each phase keeps
	// its own.
	Source Source `json:"source"`

	// Made is what the phase reported it made, through Run.Made, in the
	// order it reported it.
	Made []string `json:"made,omitempty"`

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

	// StepCheckpoint is the checkpoint after a phase: the run verified it
	// and paused there, for the operator to resume or tear down.
	StepCheckpoint Step = "checkpoint"

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

// put adds pr to the run's phases, in place of the phase's earlier entry when
// a resume applies it again.
func (rec *Record) put(pr PhaseRecord) {
	for i := range rec.Phases {
		if rec.Phases[i].Name == pr.Name {
			rec.Phases[i] = pr
			return
		}
	}
	rec.Phases = append(rec.Phases, pr)
}

// measure sets the run's time to now, and marks it when it's over the target.
func (rec *Record) measure(now time.Time) {
	rec.Total = now.Sub(rec.Started)
	rec.Over = rec.Total > TimeTarget
}

// write saves rec. Until phase 0 has made the target's state backend there is
// nowhere to keep it, and that is not a failure when beforeBackend says the
// run hasn't got that far: the target then holds nothing a later run would
// need explained. Once phase 0 has verified, a target with no backend is one
// whose record would be lost, so the run stops.
func write(ctx context.Context, store Store, rec Record, beforeBackend bool) error {
	err := store.Write(ctx, rec)
	if err == nil || beforeBackend && errors.Is(err, ErrNoBackend) {
		return nil
	}
	return fmt.Errorf("writing run %s's record: %w", rec.ID, err)
}
