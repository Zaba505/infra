package main

import (
	"context"
	"errors"
	"fmt"

	"dagger/rebuild/internal/dagger"
	"dagger/rebuild/run"
)

// errNoValues is returned by up, down and preflight when the module was given
// no values checkout.
var errNoValues = errors.New("no values checkout: set it once with dagger settings -g rebuild values <path>")

// target reads the named target from the values checkout: its file,
// targets/<name>.yaml, read strictly. A target with no file is refused, and so
// is one whose file says neither kind: live nor kind: drill.
//
// The file is read from the checkout as it stands, uncommitted changes
// included. A drill may run from them, and a live run refuses them.
func (m *Rebuild) target(ctx context.Context, name string) (run.Target, error) {
	if m.Values == nil {
		return nil, errNoValues
	}
	path, err := run.TargetFilePath(name)
	if err != nil {
		return nil, err
	}
	checkout := m.Values.Uncommitted().After()
	ok, err := checkout.Exists(ctx, path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeRegularType})
	if err != nil {
		return nil, fmt.Errorf("target %s: reading the values checkout: %w", name, err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: %s: the values checkout has no %s", run.ErrUnknownTarget, name, path)
	}
	data, err := checkout.File(path).Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("target %s: reading %s: %w", name, path, err)
	}
	f, err := run.ParseTargetFile([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("target %s: %s: %w", name, path, err)
	}
	return &target{name: name, kind: f.Kind, workspace: m.Workspace, values: m.Values}, nil
}

// target is the run.Target a call is pointed at.
type target struct {
	name      string
	kind      run.Kind
	workspace *dagger.Workspace
	values    *dagger.GitRepository
}

func (t *target) Name() string   { return t.name }
func (t *target) Kind() run.Kind { return t.kind }

// Source reads this repository's commit and uncommitted changes from the
// workspace's git, and the values' commit and uncommitted changes from the
// values checkout.
func (t *target) Source(ctx context.Context) (run.Source, error) {
	git := t.workspace.Git()
	repo, err := commit(ctx, git.Head(), git.Uncommitted())
	if err != nil {
		return run.Source{}, fmt.Errorf("this repository: %w", err)
	}
	values, err := commit(ctx, t.values.Head(), t.values.Uncommitted())
	if err != nil {
		return run.Source{}, fmt.Errorf("the values checkout: %w", err)
	}
	return run.Source{Repo: repo, Values: values}, nil
}

// commit reads one checkout's commit, and whether it has uncommitted changes,
// untracked files included.
func commit(ctx context.Context, head *dagger.GitRef, uncommitted *dagger.Changeset) (run.Commit, error) {
	sha, err := head.Commit(ctx)
	if err != nil {
		return run.Commit{}, fmt.Errorf("reading its commit: %w", err)
	}
	clean, err := uncommitted.IsEmpty(ctx)
	if err != nil {
		return run.Commit{}, fmt.Errorf("reading its uncommitted changes: %w", err)
	}
	return run.Commit{SHA: sha, Uncommitted: !clean}, nil
}

// Store is where the target's run records are kept, in its state backend.
// Phase 0 makes the state backend, and its story replaces this with the real
// store. Until then no target has one.
func (t *target) Store() run.Store { return noBackend{} }

// noBackend is the store of a target with no state backend. It holds no
// records and keeps none, which run treats as a rebuild from nothing.
type noBackend struct{}

func (noBackend) Records(context.Context) ([]run.Record, error) { return nil, run.ErrNoBackend }

func (noBackend) Write(context.Context, run.Record) error { return run.ErrNoBackend }
