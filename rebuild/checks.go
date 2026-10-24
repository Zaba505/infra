package main

import (
	"context"
	"errors"
	"fmt"

	"dagger/rebuild/internal/dagger"
	"dagger/rebuild/internal/gotool"
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
	// Both run whatever the other finds, and each failure names its Go
	// module, so one red run reports everything there is to fix.
	var errs []error
	if err := goPipeline(ctx, module, lint); err != nil {
		errs = append(errs, fmt.Errorf("rebuild/: %w", err))
	}
	if err := goPipeline(ctx, callingWorkspace.Directory(testsDir), lint); err != nil {
		errs = append(errs, fmt.Errorf("rebuild/tests/: %w", err))
	}
	return errors.Join(errs...)
}

func goPipeline(ctx context.Context, source *dagger.Directory, lint *dagger.File) error {
	return dag.Z5Labs().
		Go(source).
		WithLint(dagger.Z5LabsGoChainWithLintOpts{Config: lint}).
		WithTest(true).
		Ci(ctx)
}

// coverageFloor is the statement coverage, in percent, below which a library
// package fails the module's checks.
const coverageFloor = 80.0

// runPackages are the packages RunCoverage holds to coverageFloor: the phase
// model and its test helpers, measured together.
const runPackages = "./run/..."

// RunCoverage fails when the statement coverage of rebuild/run/, the phase
// model, drops below 80%. It is a library package, and the repository holds
// every library package to that floor.
//
// It runs go test over rebuild/run/ and the packages under it with the
// coverage of all of them counted, so a test helper package is measured by
// the tests that use it.
//
// It touches no target, so its execs may cache. It is never cached itself, for
// GoPipeline's reason: the workspace is a live view of the tree.
//
// +check
// +cache="never"
func (m *Rebuild) RunCoverage(
	ctx context.Context,
	// The workspace to check. The Dagger CLI fills this in from the workspace
	// the check runs in.
	callingWorkspace *dagger.Workspace,
) error {
	module := callingWorkspace.Directory(moduleDir, dagger.WorkspaceDirectoryOpts{
		Exclude: []string{"tests/"},
	})
	goMod, err := module.File("go.mod").Contents(ctx)
	if err != nil {
		return err
	}
	version, err := gotool.GoVersion(goMod)
	if err != nil {
		return err
	}
	const profile = "/tmp/cover.out"
	ctr := dag.Container().
		From("golang:"+version).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("rebuild-go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("rebuild-go-build")).
		WithMountedDirectory("/src", module).
		WithWorkdir("/src")
	ctr = cachedExec(ctr, []string{
		"go", "test", "-coverpkg=" + runPackages, "-coverprofile=" + profile, runPackages,
	})
	out, err := cachedExec(ctr, []string{"go", "tool", "cover", "-func=" + profile}).Stdout(ctx)
	if err != nil {
		return err
	}
	total, err := gotool.CoverageTotal(out)
	if err != nil {
		return err
	}
	if total < coverageFloor {
		return fmt.Errorf("rebuild/run/ has %.1f%% statement coverage, below the %.0f%% floor", total, coverageFloor)
	}
	return nil
}
