package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"dagger/rebuild/internal/dagger"
)

// The module calls WithExec in exactly two places, the two helpers below, and
// .golangci.yml bans it everywhere else. Every exec the module runs is
// therefore one or the other, and which one is a decision written at the call.
//
// The split exists because +cache="never" governs a function's result and
// nothing beneath it. The execs a function builds stay content addressed, so
// a never-cached function that runs the same exec on the same inputs as an
// earlier call is handed that call's output without anything running. For an
// exec that reads a target, that is a drift check reporting the state of the
// last run; for one that changes a target, it is an apply that never happened.

// nonceEnv is the environment variable targetExec puts its nonce in.
const nonceEnv = "REBUILD_EXEC_NONCE"

// targetExec runs args in ctr as an exec that reads or changes a target. It
// puts a fresh nonce in the exec's environment, so no two calls share a cache
// key and every call really runs. Use it for every exec that touches a
// target, reads included: a cached read reports no drift when there is some.
//
// The function that calls it must be marked +cache="never" as well, or the
// function's own result is cached and the exec is never reached.
func targetExec(ctr *dagger.Container, args []string, opts ...dagger.ContainerWithExecOpts) (*dagger.Container, error) {
	n, err := nonce()
	if err != nil {
		return nil, err
	}
	return ctr.WithEnvVariable(nonceEnv, n).WithExec(args, opts...), nil //nolint:forbidigo // one of the module's two exec helpers
}

// cachedExec runs args in ctr as work that may come from the cache, such as
// installing a tool while assembling a container. It adds nothing to the
// exec, so a second call on the same inputs is served from the first.
//
// Never use it for an exec that reads or changes a target; that is
// targetExec's.
func cachedExec(ctr *dagger.Container, args []string, opts ...dagger.ContainerWithExecOpts) *dagger.Container {
	return ctr.WithExec(args, opts...) //nolint:forbidigo // one of the module's two exec helpers
}

// nonce returns 32 hex characters of cryptographically random data.
func nonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate an exec nonce: %w", err)
	}
	return hex.EncodeToString(b), nil
}
