package run

import "time"

// A Clock tells a run the time. Up uses the system clock unless WithClock
// gives it another; the tests give it runtest's fake.
type Clock interface {
	Now() time.Time
}

// systemClock is the clock a run uses by default.
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }
