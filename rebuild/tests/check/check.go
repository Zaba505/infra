// Package check turns testify/require assertions into the error a Dagger
// check returns.
//
// A Dagger function has no *testing.T, and a check reports its verdict as the
// error it returns. Require runs a function with a require.TestingT and hands
// back the first assertion that failed, so a test in the tests module asserts
// with require as a Go test would.
//
// It imports nothing from Dagger, so its own tests run with go test.
package check

import (
	"fmt"

	"github.com/stretchr/testify/require"
)

// Require runs fn with a require.TestingT and returns the first assertion
// that failed as an error, or nil when none did. A failed require assertion
// stops fn where it failed, as it would stop a test. Any other panic in fn is
// let through.
func Require(fn func(t require.TestingT)) (err error) {
	t := &recorder{}
	defer func() {
		if v := recover(); v != nil {
			if _, ok := v.(failNow); !ok {
				panic(v)
			}
			err = t.failure
		}
	}()
	fn(t)
	return t.failure
}

// failNow is what recorder panics with to stop fn at a failed assertion.
type failNow struct{}

// recorder keeps the first failed assertion's message.
type recorder struct {
	failure error
}

func (r *recorder) Errorf(format string, args ...any) {
	if r.failure == nil {
		r.failure = fmt.Errorf(format, args...)
	}
}

func (r *recorder) FailNow() {
	panic(failNow{})
}
