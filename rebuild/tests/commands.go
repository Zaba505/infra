package main

import (
	"context"
	"errors"

	"dagger/tests/check"
	"dagger/tests/internal/dagger"

	"github.com/stretchr/testify/require"
)

// The tests below run the rebuild module's up, down and preflight against a
// values checkout the tests make for themselves, and a workspace made the same
// way, so what they show doesn't depend on the state of the checkout they run
// in.

// UpAgainstDrillFinishes shows that up against a drill, with no phase in the
// order, passes the preflight check as a rebuild from nothing and finishes. The
// values have uncommitted changes, which a drill may run from.
//
// +cache="never"
func (t *Tests) UpAgainstDrillFinishes(ctx context.Context) error {
	summary, err := rebuild(uncommittedValues()).Up(ctx, "drill")
	if err != nil {
		return err
	}
	return check.Require(func(rt require.TestingT) {
		require.Contains(rt, summary, " against drill\n")
		require.Contains(rt, summary, "the run finished: every phase verified.")
	})
}

// UpAgainstLiveWithCommittedValuesFinishes shows that a live run from history
// isn't refused, so the refusal below is the uncommitted values' doing.
//
// +cache="never"
func (t *Tests) UpAgainstLiveWithCommittedValuesFinishes(ctx context.Context) error {
	summary, err := rebuild(committedValues()).Up(ctx, "live")
	if err != nil {
		return err
	}
	return check.Require(func(rt require.TestingT) {
		require.Contains(rt, summary, " against live\n")
		require.Contains(rt, summary, "the run finished: every phase verified.")
	})
}

// UpAgainstLiveWithUncommittedValuesIsRefused shows that up against the live
// target refuses a values checkout with uncommitted changes.
//
// +cache="never"
func (t *Tests) UpAgainstLiveWithUncommittedValuesIsRefused(ctx context.Context) error {
	_, err := rebuild(uncommittedValues()).Up(ctx, "live")
	return refused(err, "a run against the live target refuses uncommitted changes: the target's values have uncommitted changes")
}

// DownAgainstLiveIsRefusedWithoutConfirm shows that down against the live
// target refuses unless --confirm names it, and goes ahead when it does.
//
// +cache="never"
func (t *Tests) DownAgainstLiveIsRefusedWithoutConfirm(ctx context.Context) error {
	m := rebuild(committedValues())
	_, err := m.Down(ctx, "live")
	if err := refused(err, `down was not confirmed with the target's name: live is the live target; confirm with "live" to tear it down`); err != nil {
		return err
	}
	_, err = m.Down(ctx, "live", dagger.RebuildDownOpts{Confirm: "drill"})
	if err := refused(err, `down was not confirmed with the target's name: confirmed with "drill", but the target is "live"`); err != nil {
		return err
	}
	_, err = m.Down(ctx, "live", dagger.RebuildDownOpts{Confirm: "live"})
	return err
}

// TargetWithNoFileIsRefused shows that up, down and preflight each refuse a
// target the values checkout has no file for.
//
// +cache="never"
func (t *Tests) TargetWithNoFileIsRefused(ctx context.Context) error {
	const want = "no such target: staging: the values checkout has no targets/staging.yaml"
	m := rebuild(committedValues())
	_, upErr := m.Up(ctx, "staging")
	_, downErr := m.Down(ctx, "staging", dagger.RebuildDownOpts{Confirm: "staging"})
	_, preflightErr := m.Preflight(ctx, "staging")
	return errors.Join(
		refused(upErr, want),
		refused(downErr, want),
		refused(preflightErr, want),
	)
}

// UpAgainstLiveFromAnUncommittedWorkspaceIsRefused shows that up against the
// live target reads this repository's uncommitted changes from the
// workspace's git, and refuses them.
//
// +cache="never"
func (t *Tests) UpAgainstLiveFromAnUncommittedWorkspaceIsRefused(ctx context.Context) error {
	workspace := repository(workspaceFiles, map[string]string{"NOTES.md": "not committed yet\n"})
	_, err := rebuildIn(workspace, committedValues()).Up(ctx, "live")
	return refused(err, "a run against the live target refuses uncommitted changes: this repository has uncommitted changes")
}

// AnUncommittedEditKeepsTheLiveTargetLive shows that the live target
// stays live when its file is edited to say kind: drill and not committed: up
// still refuses the uncommitted change, and down still needs --confirm.
//
// +cache="never"
func (t *Tests) AnUncommittedEditKeepsTheLiveTargetLive(ctx context.Context) error {
	m := rebuild(repository(targetFiles, map[string]string{"targets/live.yaml": "kind: drill\n"}))
	_, upErr := m.Up(ctx, "live")
	_, downErr := m.Down(ctx, "live")
	return errors.Join(
		refused(upErr, "a run against the live target refuses uncommitted changes: the target's values have uncommitted changes"),
		refused(downErr, `down was not confirmed with the target's name: live is the live target`),
	)
}

// MoreThanOneLiveTargetIsRefused shows that up, down and preflight each refuse
// to run when a second target says kind: live, even against a drill.
//
// +cache="never"
func (t *Tests) MoreThanOneLiveTargetIsRefused(ctx context.Context) error {
	const want = "more than one target is live: live, prod each say kind: live"
	m := rebuild(repository(targetFiles, map[string]string{"targets/prod.yaml": "kind: live\n"}))
	_, upErr := m.Up(ctx, "drill")
	_, downErr := m.Down(ctx, "drill")
	_, preflightErr := m.Preflight(ctx, "drill")
	return errors.Join(
		refused(upErr, want),
		refused(downErr, want),
		refused(preflightErr, want),
	)
}

// refused fails unless err is a refusal that says want.
func refused(err error, want string) error {
	return check.Require(func(rt require.TestingT) {
		require.Error(rt, err, "the call was not refused")
		require.ErrorContains(rt, err, want)
	})
}

// rebuild is the rebuild module, pointed at values and at a workspace of its
// own: a repository with one commit and nothing uncommitted. The cache probe
// needs no values, and takes nil.
func rebuild(values *dagger.GitRepository) *dagger.Rebuild {
	return rebuildIn(repository(workspaceFiles, nil), values)
}

// rebuildIn is the rebuild module, pointed at values and at workspace.
func rebuildIn(workspace, values *dagger.GitRepository) *dagger.Rebuild {
	return dag.Rebuild(workspace.AsWorkspace(), dagger.RebuildOpts{Values: values})
}

// workspaceFiles is the committed contents of the workspace the tests run the
// module in.
var workspaceFiles = map[string]string{"README.md": "the repository a run rebuilds from\n"}

// targetFiles is the values checkout's committed contents: one live target and
// one drill, and no file for staging.
var targetFiles = map[string]string{
	"targets/live.yaml":  "kind: live\n",
	"targets/drill.yaml": "kind: drill\n",
}

// committedValues is a values checkout with nothing uncommitted.
func committedValues() *dagger.GitRepository {
	return repository(targetFiles, nil)
}

// uncommittedValues is the same values checkout with an uncommitted change on
// top, as an operator's checkout has while they edit it.
func uncommittedValues() *dagger.GitRepository {
	return repository(targetFiles, map[string]string{"NOTES.md": "not committed yet\n"})
}
