// Tests is the rebuild module's tests module. Each test is a function that
// returns an error, so it can be run on its own with dagger call, and All runs
// every one of them as the module's check.
package main

import (
	"context"

	"dagger/tests/check"

	"github.com/dagger/dagger/util/parallel"
	"github.com/stretchr/testify/require"
)

// Tests holds no state.
type Tests struct{}

// All runs every test in this module.
//
// +check
// +cache="never"
func (t *Tests) All(ctx context.Context) error {
	return parallel.New().
		WithJob("an exec that touches a target runs on every call", t.TargetExecIsNeverCached).
		WithJob("an exec that may cache comes from the cache", t.CachedExecIsCached).
		WithJob("up against a drill with no phases finishes", t.UpAgainstDrillFinishes).
		WithJob("up against the live target with committed values finishes", t.UpAgainstLiveWithCommittedValuesFinishes).
		WithJob("up against the live target with uncommitted values is refused", t.UpAgainstLiveWithUncommittedValuesIsRefused).
		WithJob("down against the live target is refused without --confirm", t.DownAgainstLiveIsRefusedWithoutConfirm).
		WithJob("a target with no file in the values checkout is refused", t.TargetWithNoFileIsRefused).
		Run(ctx)
}

// TargetExecIsNeverCached shows that an exec through the helper for execs
// that touch a target runs again on a second call: the probe prints a fresh
// UUID each time it really runs, so two calls print two different ones.
//
// +cache="never"
func (t *Tests) TargetExecIsNeverCached(ctx context.Context) error {
	first, err := rebuild(nil).CacheProbe().Fresh(ctx)
	if err != nil {
		return err
	}
	second, err := rebuild(nil).CacheProbe().Fresh(ctx)
	if err != nil {
		return err
	}
	return check.Require(func(rt require.TestingT) {
		require.NotEmpty(rt, first)
		require.NotEqual(rt, first, second, "the second exec came from the cache")
	})
}

// CachedExecIsCached shows that an exec through the helper for work that may
// cache comes back from the cache: the function around it is never cached, so
// a second call that prints the first call's UUID was served by the exec's
// cache rather than run again.
//
// +cache="never"
func (t *Tests) CachedExecIsCached(ctx context.Context) error {
	first, err := rebuild(nil).CacheProbe().Cached(ctx)
	if err != nil {
		return err
	}
	second, err := rebuild(nil).CacheProbe().Cached(ctx)
	if err != nil {
		return err
	}
	return check.Require(func(rt require.TestingT) {
		require.NotEmpty(rt, first)
		require.Equal(rt, first, second, "the second exec ran again")
	})
}
