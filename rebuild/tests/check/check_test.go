package check

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequire(t *testing.T) {
	t.Run("returns nil when every assertion holds", func(t *testing.T) {
		err := Require(func(rt require.TestingT) {
			require.Equal(rt, "a", "a")
		})
		require.NoError(t, err)
	})

	t.Run("returns the failed assertion and stops at it", func(t *testing.T) {
		reached := false
		err := Require(func(rt require.TestingT) {
			require.Equal(rt, "a", "b", "the values differ")
			reached = true
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "the values differ")
		require.False(t, reached, "Require ran past a failed assertion")
	})

	t.Run("keeps the first failure", func(t *testing.T) {
		err := Require(func(rt require.TestingT) {
			rt.Errorf("first")
			rt.Errorf("second")
		})
		require.EqualError(t, err, "first")
	})

	t.Run("lets any other panic through", func(t *testing.T) {
		boom := errors.New("boom")
		require.PanicsWithValue(t, boom, func() {
			_ = Require(func(require.TestingT) {
				panic(boom)
			})
		})
	})
}
