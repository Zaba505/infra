// Rebuild is the platform's rebuild entry point: the Dagger module that
// stands the platform up from this repository, and tears it down again.
//
// Its commands are up, down and preflight, each against one named target:
//
//	dagger call up        --target <name> [--resume] [--through <phase>] [--no-pause]
//	dagger call down      --target <name> [--confirm <name>]
//	dagger call preflight --target <name>
//
// Each is a thin layer over rebuild/run/, the phase model, which imports
// nothing from Dagger. The module reads the target from the values checkout,
// builds the phases in rebuild/phases.go, and hands both to run. It also holds
// the rule every phase depends on — an exec that reads or changes a target
// never comes from the cache — and the checks that hold the module's Go to the
// repository's rules.
package main

import (
	"context"

	"dagger/rebuild/internal/dagger"
)

// Rebuild is the module. It holds what every run reads from the machine that
// starts it, since a Dagger function can't read the host: this repository's
// workspace and the values checkout.
type Rebuild struct {
	// Workspace is this repository's workspace. A run reads the commit it
	// rebuilds from, and whether there are uncommitted changes on top, from
	// its git.
	//
	// +private
	Workspace *dagger.Workspace

	// Values is the values checkout: a git repository outside this one that
	// holds a file per target, targets/<name>.yaml. A run reads the target's
	// file from it, and the values' commit and uncommitted changes.
	//
	// +private
	Values *dagger.GitRepository
}

// New returns the module.
func New(
	// This repository's workspace. The Dagger CLI fills it in from the
	// workspace the call runs in.
	workspace *dagger.Workspace,
	// The values checkout, a git repository holding a file per target,
	// targets/<name>.yaml. Set it once, outside this repository, with
	// dagger settings -g rebuild values <path>. up, down and preflight refuse
	// to run without it; the module's checks don't need it.
	//
	// +optional
	values *dagger.GitRepository,
) *Rebuild {
	return &Rebuild{Workspace: workspace, Values: values}
}

// probeImage is the image CacheProbe runs its command in. dagger.lock pins
// its digest.
const probeImage = "alpine:3.22"

// CacheProbe returns the probe the tests module uses to show what each of the
// module's two exec helpers does with the cache. It touches no target.
//
// It is for the module's tests, not part of the rebuild interface. It is
// public only because a tests module can reach nothing else.
func (m *Rebuild) CacheProbe() *CacheProbe {
	return &CacheProbe{}
}

// CacheProbe runs one command, which prints a random UUID, through each of the
// module's exec helpers. Both of its functions are never cached themselves, so
// the only thing that can hand a second call the first call's UUID is the
// cache of the exec underneath.
type CacheProbe struct{}

// Fresh runs the probe's command through the helper for execs that read or
// change a target, and returns what it printed. Two calls print two different
// UUIDs.
//
// +cache="never"
func (p *CacheProbe) Fresh(ctx context.Context) (string, error) {
	return targetExec(p.container(), probeCommand).Stdout(ctx)
}

// Cached runs the probe's command through the helper for work that may cache,
// and returns what it printed. A second call prints the first call's UUID.
//
// +cache="never"
func (p *CacheProbe) Cached(ctx context.Context) (string, error) {
	return cachedExec(p.container(), probeCommand).Stdout(ctx)
}

// probeCommand prints a UUID the kernel makes fresh on every read, so an exec
// that really ran twice prints two different ones.
var probeCommand = []string{"cat", "/proc/sys/kernel/random/uuid"}

func (p *CacheProbe) container() *dagger.Container {
	return dag.Container().From(probeImage)
}
