package main

import (
	"context"

	"dagger/rebuild/internal/dagger"
)

// The module's Go lives in two Go modules: the module itself, in rebuild/, and
// its tests module, in rebuild/tests/. Each is checked on its own, because go
// vet, golangci-lint and go test stop at a nested go.mod.
const (
	moduleDir = "/rebuild"
	testsDir  = "/rebuild/tests"
	lintFile  = "/.golangci.yml"
)

// GoPipeline runs devex's z5labs Go pipeline over the module's Go and the
// Go of its tests module: gofmt, go vet, golangci-lint with this
// repository's .golangci.yml, and go test with the race detector.
//
// Generated code is left to the generator. golangci-lint skips files marked
// as generated, and the Go SDK's own check, dagger-go-sdk:stale, is what
// holds them to what the generator writes.
//
// It is never cached, because the workspace the CLI fills in is a live view
// of the tree rather than an argument a cache key can describe: a cached pass
// would be a pass for a tree the check never looked at.
//
// +check
// +cache="never"
func (m *Rebuild) GoPipeline(
	ctx context.Context,
	// The workspace to check. The Dagger CLI fills this in from the workspace
	// the check runs in.
	callingWorkspace *dagger.Workspace,
) error {
	lint := callingWorkspace.File(lintFile)
	module := callingWorkspace.Directory(moduleDir, dagger.WorkspaceDirectoryOpts{
		// The tests module is a Go module of its own, checked below.
		Exclude: []string{"tests/"},
	})
	if err := goPipeline(ctx, module, lint); err != nil {
		return err
	}
	return goPipeline(ctx, callingWorkspace.Directory(testsDir), lint)
}

func goPipeline(ctx context.Context, source *dagger.Directory, lint *dagger.File) error {
	return dag.Z5Labs().
		Go(source).
		WithLint(dagger.Z5LabsGoChainWithLintOpts{Config: lint}).
		WithTest(true).
		Ci(ctx)
}
