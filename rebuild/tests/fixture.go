package main

import (
	"dagger/tests/internal/dagger"
)

// gitImage is the image the tests make their git repositories in. dagger.lock
// pins its digest.
const gitImage = "alpine/git:v2.54.0"

// fixtureDate is the date of every commit the tests make, so the same files
// make the same commit and the fixture comes from the cache.
const fixtureDate = "2026-01-01T00:00:00Z"

// repository makes a git repository with committed in one commit, and
// uncommitted written on top of it without being committed.
func repository(committed, uncommitted map[string]string) *dagger.GitRepository {
	files := dag.Directory()
	for path, contents := range committed {
		files = files.WithNewFile(path, contents)
	}
	ctr := dag.Container().
		From(gitImage).
		WithDirectory("/repo", files).
		WithWorkdir("/repo").
		WithEnvVariable("GIT_AUTHOR_NAME", "rebuild tests").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "rebuild-tests@example.com").
		WithEnvVariable("GIT_AUTHOR_DATE", fixtureDate).
		WithEnvVariable("GIT_COMMITTER_NAME", "rebuild tests").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "rebuild-tests@example.com").
		WithEnvVariable("GIT_COMMITTER_DATE", fixtureDate)
	ctr = fixtureExec(ctr, "git", "init", "--quiet", "--initial-branch=main")
	ctr = fixtureExec(ctr, "git", "add", "--all")
	ctr = fixtureExec(ctr, "git", "commit", "--quiet", "--message=fixture")

	repo := ctr.Directory("/repo")
	for path, contents := range uncommitted {
		repo = repo.WithNewFile(path, contents)
	}
	return repo.AsGit()
}

// fixtureExec runs args in ctr to make a test fixture. It touches no target,
// so it may come from the cache, as the rebuild module's cachedExec may. It is
// the tests module's one call to WithExec.
func fixtureExec(ctr *dagger.Container, args ...string) *dagger.Container {
	return ctr.WithExec(args) //nolint:forbidigo // the tests module's one exec helper, for fixtures that touch no target
}
