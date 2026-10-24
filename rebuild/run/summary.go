package run

import (
	"fmt"
	"strings"
	"time"
)

// Summary is what Up returns: the phases the call ran, each with what it made
// and how long it took, and where the run stands, with its time so far.
type Summary struct {
	// ID is the run's ID.
	ID string

	// Target is the name of the target the run is pointed at.
	Target string

	// Phases holds the phases this call applied, in the order it applied
	// them. Phases an earlier call applied are in the run's record.
	Phases []PhaseRecord

	// Finished, Stopped, Total and Over are the run's, as its record holds
	// them once the call returns.
	Finished bool
	Stopped  *Stop
	Total    time.Duration
	Over     bool

	// Recorded is false when the target kept no record of the run, because
	// it had no state backend to keep it in, as when phase 0 fails before
	// making one. Such a run can't be resumed: the next up starts again.
	Recorded bool
}

// summarize is the summary of a call that ran the phases in ran, leaving the
// run's record as rec, which the target kept when kept says so.
func summarize(rec Record, ran []PhaseRecord, kept bool) Summary {
	return Summary{
		ID:       rec.ID,
		Target:   rec.Target,
		Phases:   ran,
		Finished: rec.Finished,
		Stopped:  rec.Stopped,
		Total:    rec.Total,
		Over:     rec.Over,
		Recorded: kept,
	}
}

// String is the summary as the operator reads it at a checkpoint: each phase
// the call ran, what it made and how long it took, the run's time so far,
// where the run stands and what the operator can do next.
func (s Summary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "run %s against %s\n", s.ID, s.Target)
	for _, p := range s.Phases {
		state := "verified"
		if !p.Verified {
			state = "did not verify"
		}
		fmt.Fprintf(&b, "phase %s: %s, %s\n", p.Name, state, p.Duration)
		for _, made := range p.Made {
			fmt.Fprintf(&b, "  made %s\n", made)
		}
	}

	if s.Finished {
		fmt.Fprintf(&b, "the run took %s", s.Total)
	} else {
		fmt.Fprintf(&b, "the run's time so far: %s", s.Total)
	}
	if s.Over {
		fmt.Fprintf(&b, ", over the %s target", TimeTarget)
	}
	b.WriteString("\n")

	const next = "Continue with up --resume, or tear the target down with down."
	switch st := s.Stopped; {
	case s.Finished:
		b.WriteString("the run finished: every phase verified.\n")
	case st == nil:
		b.WriteString("the run is still going.\n")
	case (st.Step == StepApply || st.Step == StepVerify) && !s.Recorded:
		fmt.Fprintf(&b, "the run failed at phase %s's %s: %s\n", st.Phase, st.Step, st.Error)
		b.WriteString("The target has no state backend to keep the run's record in, so there is no run to resume. " +
			"Start again with up, or tear the target down with down.\n")
	case st.Step == StepApply || st.Step == StepVerify:
		fmt.Fprintf(&b, "the run failed at phase %s's %s: %s\n", st.Phase, st.Step, st.Error)
		b.WriteString("Once the target is fit to resume, continue with up --resume, which applies phase " +
			st.Phase + " again. Or tear the target down with down.\n")
	default:
		fmt.Fprintf(&b, "the run%s. %s\n", describeStop(st), next)
	}
	return b.String()
}
