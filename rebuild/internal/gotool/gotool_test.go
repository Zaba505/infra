package gotool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoverageTotal_ReadsTheTotalLine(t *testing.T) {
	out := "dagger/rebuild/run/run.go:91:\tString\t\t100.0%\n" +
		"dagger/rebuild/run/up.go:47:\tUp\t\t75.0%\n" +
		"total:\t\t\t\t\t(statements)\t79.9%\n"

	total, err := CoverageTotal(out)
	require.NoError(t, err)
	require.InDelta(t, 79.9, total, 0.001)
}

func TestCoverageTotal_FailsWithNoTotal(t *testing.T) {
	for name, out := range map[string]string{
		"empty":             "",
		"no total line":     "dagger/rebuild/run/up.go:47:\tUp\t\t75.0%\n",
		"no percent":        "total:\t(statements)\t79.9\n",
		"not a number":      "total:\t(statements)\tmany%\n",
		"a bare total line": "total:\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := CoverageTotal(out)
			require.Error(t, err)
		})
	}
}

func TestGoVersion_ReadsTheGoDirective(t *testing.T) {
	version, err := GoVersion("module dagger/rebuild\n\ngo 1.26.8\n\nrequire (\n\tgo.example v1\n)\n")
	require.NoError(t, err)
	require.Equal(t, "1.26.8", version)

	_, err = GoVersion("module dagger/rebuild\n")
	require.Error(t, err)
}
