// Package gotool reads what the Go toolchain writes, for the module's checks.
// It imports nothing from Dagger, so its tests run with go test; the module's
// main package can't be tested that way, since its Dagger client needs an
// engine session to start.
package gotool

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GoVersion returns the version in a go.mod's go directive, which is the
// toolchain image devex's Go pipeline tests with too.
func GoVersion(goMod string) (string, error) {
	for line := range strings.Lines(goMod) {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "go" {
			return fields[1], nil
		}
	}
	return "", errors.New("go.mod has no go directive")
}

// CoverageTotal reads the total statement coverage, in percent, from the
// output of go tool cover -func, whose last line reads
// "total: (statements) 85.3%".
func CoverageTotal(funcOutput string) (float64, error) {
	for line := range strings.Lines(funcOutput) {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "total:" {
			continue
		}
		pct, ok := strings.CutSuffix(fields[len(fields)-1], "%")
		if !ok {
			break
		}
		total, err := strconv.ParseFloat(pct, 64)
		if err != nil {
			return 0, fmt.Errorf("reading the total coverage from %q: %w", strings.TrimSpace(line), err)
		}
		return total, nil
	}
	return 0, fmt.Errorf("go tool cover reported no total coverage:\n%s", funcOutput)
}
