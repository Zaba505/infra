package run_test

import (
	"testing"

	"dagger/rebuild/run"
	"dagger/rebuild/run/runtest"

	"github.com/stretchr/testify/require"
)

func TestDown_TearsDownEveryPhaseInReverseWithPhaseZeroLast(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(nil)
	require.NoError(t, up(t.Context(), target, ps, run.NoPause()))

	require.NoError(t, run.Down(t.Context(), target, ps))
	require.Equal(t, []string{
		"0 apply", "0 verify", "1a apply", "1a verify", "1b apply", "1b verify",
		"1b teardown", "1a teardown", "0 teardown",
	}, target.Calls())
	require.Empty(t, target.Left())
}

func TestDown_SucceedsWhenApplyNeverRan(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)

	require.NoError(t, run.Down(t.Context(), target, phases(nil)))
	require.Equal(t, []string{"1b teardown", "1a teardown", "0 teardown"}, target.Calls())
	require.Empty(t, target.Left())
}

func TestDown_SucceedsWhenRunTwice(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(nil)
	require.NoError(t, up(t.Context(), target, ps, run.NoPause()))

	require.NoError(t, run.Down(t.Context(), target, ps))
	require.NoError(t, run.Down(t.Context(), target, ps))
	require.Empty(t, target.Left())
}

func TestDown_SucceedsAfterAnUpThatFailedAtEachPhaseInTurn(t *testing.T) {
	for _, failing := range order {
		for step, fail := range map[run.Step]func(error) runtest.PhaseOption{
			run.StepApply:  runtest.FailApply,
			run.StepVerify: runtest.FailVerify,
		} {
			t.Run(failing+" "+string(step), func(t *testing.T) {
				target := runtest.NewTarget(t, runtest.Drill)
				ps := phases(map[string][]runtest.PhaseOption{failing: {fail(errNoEdge)}})
				require.ErrorIs(t, up(t.Context(), target, ps, run.NoPause()), errNoEdge)
				require.NotEmpty(t, target.Left())

				require.NoError(t, run.Down(t.Context(), target, ps))
				require.Empty(t, target.Left())
			})
		}
	}
}

func TestDown_RunAgainAfterFailingPartwayFinishesTheTeardown(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill)
	ps := phases(map[string][]runtest.PhaseOption{"1a": {runtest.FailTeardownOnce(errBroken)}})
	require.NoError(t, up(t.Context(), target, ps, run.NoPause()))
	applied := len(target.Calls())

	err := run.Down(t.Context(), target, ps)
	require.ErrorIs(t, err, errBroken)
	require.ErrorContains(t, err, "phase 1a: teardown")
	require.Equal(t, []string{"1b teardown", "1a teardown"}, target.Calls()[applied:],
		"phase 0 is not torn down while a later phase still stands")
	require.Equal(t, []string{"0", "1a"}, target.Left())

	require.NoError(t, run.Down(t.Context(), target, ps))
	require.Empty(t, target.Left())
}

func TestDown_AgainstTheLiveTargetRefusesWithoutItsName(t *testing.T) {
	for name, opts := range map[string][]run.DownOption{
		"unconfirmed":           nil,
		"confirmed with blank":  {run.Confirm("")},
		"confirmed as another":  {run.Confirm("drill")},
		"confirmed differently": {run.Confirm("Live")},
	} {
		t.Run(name, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Live)

			err := run.Down(t.Context(), target, phases(nil), opts...)
			require.ErrorIs(t, err, run.ErrNotConfirmed)
			require.Empty(t, target.Calls())
			if len(opts) == 0 || name == "confirmed with blank" {
				require.ErrorContains(t, err, `live is the live target; confirm with "live"`)
			}
		})
	}
}

func TestDown_AgainstTheLiveTargetConfirmedWithItsNameTearsDown(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Live, runtest.Named("prod"))

	require.NoError(t, run.Down(t.Context(), target, phases(nil), run.Confirm("prod")))
	require.Len(t, target.Calls(), 3)
}

func TestDown_AgainstADrillNeedsNoConfirmationButRefusesTheWrongName(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.Named("scratch"))

	require.NoError(t, run.Down(t.Context(), target, phases(nil)))
	require.NoError(t, run.Down(t.Context(), target, phases(nil), run.Confirm("scratch")))
	err := run.Down(t.Context(), target, phases(nil), run.Confirm("prod"))
	require.ErrorIs(t, err, run.ErrNotConfirmed)
	require.ErrorContains(t, err, `confirmed with "prod", but the target is "scratch"`)
	require.NotContains(t, err.Error(), "live")
	require.Len(t, target.Calls(), 6)
}

func TestDown_AgainstTheLiveTargetRefusesUncommittedChanges(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Live, runtest.WithSource(run.Source{
		Repo: run.Commit{SHA: "a", Uncommitted: true},
	}))

	err := run.Down(t.Context(), target, phases(nil), run.Confirm("live"))
	require.ErrorIs(t, err, run.ErrUncommitted)
	require.Empty(t, target.Calls())
}

func TestDown_StopsWhenTheTargetsSourceCannotBeRead(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.FailSource(errBroken))

	require.ErrorIs(t, run.Down(t.Context(), target, phases(nil)), errBroken)
	require.Empty(t, target.Calls())
}

func TestDown_AfterAFailedUpLetsTheNextUpStartFromNothing(t *testing.T) {
	for _, failing := range order {
		t.Run(failing, func(t *testing.T) {
			target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())
			opts := map[string][]runtest.PhaseOption{"0": {runtest.ProvidesBackend()}}
			opts[failing] = append(opts[failing], runtest.FailVerifyOnce(errNoEdge))
			failed := phases(opts)
			require.ErrorIs(t, up(t.Context(), target, failed, run.NoPause()), errNoEdge)
			require.ErrorIs(t, up(t.Context(), target, failed, run.NoPause()), run.ErrUnfinishedRun)

			require.NoError(t, run.Down(t.Context(), target, failed))
			_, ok, err := run.LastRun(t.Context(), target.Store())
			require.NoError(t, err)
			require.False(t, ok, "phase 0's teardown took the records with the backend")

			require.NoError(t, up(t.Context(), target, failed, run.NoPause()))
			require.True(t, lastRun(t, target).Finished)
		})
	}
}

func TestDown_ThatFailedPartwayKeepsTheNextUpOffWhatIsLeft(t *testing.T) {
	target := runtest.NewTarget(t, runtest.Drill, runtest.NoBackend())
	ps := phases(map[string][]runtest.PhaseOption{
		"0":  {runtest.ProvidesBackend()},
		"1a": {runtest.FailTeardownOnce(errBroken)},
		"1b": {runtest.FailVerify(errNoEdge)},
	})
	require.ErrorIs(t, up(t.Context(), target, ps, run.NoPause()), errNoEdge)

	require.ErrorIs(t, run.Down(t.Context(), target, ps), errBroken)
	require.ErrorIs(t, up(t.Context(), target, ps, run.NoPause()), run.ErrUnfinishedRun)

	require.NoError(t, run.Down(t.Context(), target, ps))
	require.Empty(t, target.Left())
}
