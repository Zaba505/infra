// Package runtest holds fakes for testing code that uses package run: a
// target that records every call made to it, phases that fail on demand and
// report what they made, an in-memory record store and a clock that moves only
// when it's told to.
//
// A fake phase records each call on the target it runs against, as
// "<name> apply", "<name> verify" or "<name> teardown", so a test reads the
// order a run made its calls in from Target.Calls.
package runtest

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"dagger/rebuild/run"
)

// The two kinds of target, for NewTarget.
const (
	Live  = run.Live
	Drill = run.Drill
)

// A TargetOption changes the target NewTarget makes.
type TargetOption func(*Target)

// Named gives the target a name. Without it, the target is named after its
// kind, "live" or "drill".
func Named(name string) TargetOption {
	return func(t *Target) { t.name = name }
}

// WithSource sets the commits a run against the target is rebuilt from.
// Without it, both are clean commits.
func WithSource(src run.Source) TargetOption {
	return func(t *Target) { t.source = src }
}

// FailSource makes reading the target's source fail with err.
func FailSource(err error) TargetOption {
	return func(t *Target) { t.sourceErr = err }
}

// NoBackend makes a target with no state backend, as before phase 0 has made
// one: its store holds no records and keeps none until a phase made with
// ProvidesBackend applies.
func NoBackend() TargetOption {
	return func(t *Target) { t.store.backend = false }
}

// FailStore makes the target's store fail every read and write with err.
func FailStore(err error) TargetOption {
	return func(t *Target) { t.store.err = err }
}

// FailStoreWrites makes the target's store fail every write after the first n
// with err, and keep reading as before.
func FailStoreWrites(n int, err error) TargetOption {
	return func(t *Target) {
		t.store.writesLeft = n
		t.store.writeErr = err
	}
}

// Target is a fake run.Target. It records every call a runtest phase makes to
// it, and which phases have made something and not torn it down.
type Target struct {
	name      string
	kind      run.Kind
	source    run.Source
	sourceErr error
	store     *Store

	mu    sync.Mutex
	calls []string
	made  []string
}

var _ run.Target = (*Target)(nil)

// NewTarget makes a fake target of the given kind, with a state backend and no
// records in it.
func NewTarget(tb testing.TB, kind run.Kind, opts ...TargetOption) *Target {
	tb.Helper()
	t := &Target{
		name: kind.String(),
		kind: kind,
		source: run.Source{
			Repo:   run.Commit{SHA: "0000000000000000000000000000000000000001"},
			Values: run.Commit{SHA: "0000000000000000000000000000000000000002"},
		},
		store: &Store{backend: true, writesLeft: -1},
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Name returns the target's name.
func (t *Target) Name() string { return t.name }

// Kind returns the target's kind.
func (t *Target) Kind() run.Kind { return t.kind }

// Source returns the commits a run against the target is rebuilt from.
func (t *Target) Source(context.Context) (run.Source, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.source, t.sourceErr
}

// SetSource changes the commits the next run against the target is rebuilt
// from, as a commit made between two calls of a run does.
func (t *Target) SetSource(src run.Source) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.source = src
}

// Store returns the target's in-memory record store.
func (t *Target) Store() run.Store { return t.store }

// Calls returns every call runtest phases have made to the target, in order.
func (t *Target) Calls() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return slices.Clone(t.calls)
}

// Left returns the phases whose Apply has run, even partway, and whose
// Teardown hasn't run since, in the order they were first applied. After a
// Down that succeeded it's empty.
func (t *Target) Left() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return slices.Clone(t.made)
}

// call records a call to phase's method, which returned err. An Apply leaves
// the phase on the target even when it failed, since it may have failed
// partway; only a Teardown that passed removes it.
func (t *Target) call(phase, method string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.calls = append(t.calls, phase+" "+method)
	switch {
	case method == "apply" && !slices.Contains(t.made, phase):
		t.made = append(t.made, phase)
	case method == "teardown" && err == nil:
		t.made = slices.DeleteFunc(t.made, func(p string) bool { return p == phase })
	}
}

// Store is an in-memory run.Store.
type Store struct {
	mu         sync.Mutex
	backend    bool
	err        error
	writeErr   error
	writesLeft int
	records    []run.Record
}

var _ run.Store = (*Store)(nil)

// Records returns every record the store holds, oldest first.
func (s *Store) Records(ctx context.Context) ([]run.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.check(ctx); err != nil {
		return nil, err
	}
	recs := make([]run.Record, len(s.records))
	for i, rec := range s.records {
		recs[i] = clone(rec)
	}
	return recs, nil
}

// Write saves rec, replacing the record with the same ID.
func (s *Store) Write(ctx context.Context, rec run.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.check(ctx); err != nil {
		return err
	}
	if s.writesLeft == 0 {
		return s.writeErr
	}
	if s.writesLeft > 0 {
		s.writesLeft--
	}
	rec = clone(rec)
	for i := range s.records {
		if s.records[i].ID == rec.ID {
			s.records[i] = rec
			return nil
		}
	}
	s.records = append(s.records, rec)
	return nil
}

// setBackend makes or removes the store's backend. Removing it removes every
// record in it.
func (s *Store) setBackend(on bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backend = on
	if !on {
		s.records = nil
	}
}

// check fails a read or a write as a real store would: on a context that's
// done, on the error FailStore set, or with no backend.
func (s *Store) check(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.err != nil {
		return s.err
	}
	if !s.backend {
		return run.ErrNoBackend
	}
	return nil
}

// clone copies rec so the store shares nothing with the run that wrote it.
func clone(rec run.Record) run.Record {
	rec.Phases = slices.Clone(rec.Phases)
	for i := range rec.Phases {
		rec.Phases[i].Made = slices.Clone(rec.Phases[i].Made)
	}
	if rec.Stopped != nil {
		stop := *rec.Stopped
		rec.Stopped = &stop
	}
	return rec
}

// A PhaseOption changes the phase Phase makes.
type PhaseOption func(*FakePhase)

// FailApply makes every Apply fail with err.
func FailApply(err error) PhaseOption {
	return func(p *FakePhase) { p.apply = failure{err: err, times: -1} }
}

// FailApplyOnce makes the first Apply fail with err, and every later one
// pass.
func FailApplyOnce(err error) PhaseOption {
	return func(p *FakePhase) { p.apply = failure{err: err, times: 1} }
}

// FailVerify makes every Verify fail with err.
func FailVerify(err error) PhaseOption {
	return func(p *FakePhase) { p.verify = failure{err: err, times: -1} }
}

// FailVerifyOnce makes the first Verify fail with err, and every later one
// pass.
func FailVerifyOnce(err error) PhaseOption {
	return func(p *FakePhase) { p.verify = failure{err: err, times: 1} }
}

// FailTeardownOnce makes the first Teardown fail with err, and every later one
// pass.
func FailTeardownOnce(err error) PhaseOption {
	return func(p *FakePhase) { p.teardown = failure{err: err, times: 1} }
}

// ProvidesBackend makes the phase the one that holds the target's state
// backend, as phase 0 does: an Apply that passes makes the backend, and a
// Teardown that passes removes it and every record in it.
func ProvidesBackend() PhaseOption {
	return func(p *FakePhase) { p.backend = true }
}

// Makes makes every Apply report each of what through Run.Made, in order,
// before it returns, whether it passes or fails.
func Makes(what ...string) PhaseOption {
	return func(p *FakePhase) { p.makes = what }
}

// Takes makes every Apply and every Verify move clock on by d, as though each
// took that long.
func Takes(clock *Clock, d time.Duration) PhaseOption {
	return func(p *FakePhase) {
		p.clock = clock
		p.takes = d
	}
}

// FakePhase is a fake run.Phase. Make one with Phase.
type FakePhase struct {
	name     string
	apply    failure
	verify   failure
	teardown failure
	clock    *Clock
	takes    time.Duration
	backend  bool
	makes    []string
}

var _ run.Phase = (*FakePhase)(nil)

// Phase makes a fake phase named name. With no options, every method passes.
func Phase(name string, opts ...PhaseOption) *FakePhase {
	p := &FakePhase{name: name}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// ErrNotARuntestTarget is what a runtest phase returns when it's run against a
// target NewTarget didn't make, since it has nowhere to record its calls.
var ErrNotARuntestTarget = errors.New("runtest: a runtest phase only runs against a target made by runtest.NewTarget")

// Name returns the phase's name.
func (p *FakePhase) Name() string { return p.name }

// Apply records "<name> apply" on the target.
func (p *FakePhase) Apply(_ context.Context, r *run.Run) error {
	return p.do(r, "apply", &p.apply, true)
}

// Verify records "<name> verify" on the target.
func (p *FakePhase) Verify(_ context.Context, r *run.Run) error {
	return p.do(r, "verify", &p.verify, true)
}

// Teardown records "<name> teardown" on the target.
func (p *FakePhase) Teardown(_ context.Context, r *run.Run) error {
	return p.do(r, "teardown", &p.teardown, false)
}

func (p *FakePhase) do(r *run.Run, method string, f *failure, timed bool) error {
	t, ok := r.Target.(*Target)
	if !ok {
		return ErrNotARuntestTarget
	}
	if timed && p.clock != nil {
		p.clock.Advance(p.takes)
	}
	err := f.next()
	t.call(p.name, method, err)
	if method == "apply" {
		for _, what := range p.makes {
			r.Made(what)
		}
	}
	if p.backend && err == nil {
		switch method {
		case "apply":
			t.store.setBackend(true)
		case "teardown":
			t.store.setBackend(false)
		}
	}
	return err
}

// failure is when a fake method fails, and with what.
type failure struct {
	err error
	// times is how many more calls fail: -1 for every call.
	times int
}

func (f *failure) next() error {
	if f.err == nil || f.times == 0 {
		return nil
	}
	if f.times > 0 {
		f.times--
	}
	return f.err
}

// Clock is a fake run.Clock. It moves only when Advance moves it.
type Clock struct {
	mu  sync.Mutex
	now time.Time
}

var _ run.Clock = (*Clock)(nil)

// NewClock makes a clock that reads start until it's advanced.
func NewClock(start time.Time) *Clock {
	return &Clock{now: start}
}

// Now returns the clock's time.
func (c *Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Advance moves the clock on by d.
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
