// Rebuild is the platform's rebuild entry point: the Dagger module that
// stands the platform up from this repository, and tears it down again.
//
// It lands its commands in later stories. Today it holds the rule every one
// of them depends on — an exec that reads or changes a target never comes from
// the cache — and the checks that hold the module's Go to the repository's
// rules.
package main

import (
	"context"

	"dagger/rebuild/internal/dagger"
)

// Rebuild holds no state.
type Rebuild struct{}

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
