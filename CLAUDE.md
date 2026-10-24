# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What infra is

infra defines a self-hosted application platform: the place the operator's own capabilities run
by default, so that none of them depends on a vendor's hosting product.

- **It's rebuilt, not repaired.** The platform can be stood up from nothing, from this
  repository, in at most 60 minutes.
- **It's self-hosted, not self-owned.** It spans a public cloud and the home lab. "Self-hosted"
  means the operator controls it end to end, not that every part runs on hardware the operator
  owns.
- **It's run by one operator,** in at most 2 hours a week. A designated successor holds sealed
  credentials and uses them only if the operator can't.
- **Its tenants are the operator's capabilities.** Their end users reach the tenants, never the
  platform.

**infra is being rebuilt from scratch.** The charter,
[Proposal: rebuilding infra from scratch](https://github.com/Zaba505/infra/discussions/970), and
its revisions are the source of truth for what the platform is, the principles and rules every
change is held to, and the order its parts come back in. Read it before proposing or
implementing anything. Each part of the platform returns through its own stories, at the path
its proposal names; until then, the repository holds no platform code.

## Workflow

Work moves from discussion to merged code in four steps:

1. **Proposal.** Each part of the platform is proposed in a discussion in General titled
   "Proposal: …", with the sections Summary, Background, Proposed design, Design decisions,
   Drill, Documentation, What it replaces, Open questions and Out of scope. It says where its
   files live in the repository, and Documentation says where its pages sit in the docs site.
   Review happens in comments; a revised proposal is posted as a comment, and the body gets a
   note pointing to it.
2. **Milestone.** An accepted proposal gets a milestone named `Wave N: <proposal>`, whose
   description links to the proposal and to the charter. A milestone groups one proposal's
   stories; the waves, not the milestones, set the delivery order.
3. **Stories.** Stories are written from the accepted proposal only, with the `story` label, in
   its milestone and on the infra project, and each links back to the discussion.
4. **Backlog loop.** The backlog plugin's loop picks up eligible stories, implements them and
   hands the merge to GitHub. Its configuration is `.claude/backlog.json`.

At each wave boundary the charter gets a revision saying which proposals are done, which stories
landed them, and which wave starts.

## Issue tracking

Stories live on the **[infra project board](https://github.com/users/Zaba505/projects/9)**
(user-scoped). Resolve it at run time rather than hard-coding the number:

```bash
OWNER=$(gh repo view --json owner --jq .owner.login)
PROJECT=$(gh project list --owner "$OWNER" --format json \
  --jq '.projects[] | select(.title=="infra") | .number')
```

Needs the `project` token scope — `gh auth refresh -s project`.

- **Dependencies are native GitHub `blocked by` relationships**, written with
  `gh issue edit {issue} --add-blocked-by {prereq}` and read back from the `blockedBy` GraphQL
  field. Prose in a body is not a dependency, and there is no `Depends on` section in issue
  bodies or in `.github/ISSUE_TEMPLATE/story.yaml` — do not introduce one.
- **A batch of stories is created first and wired in a second pass.** Create every story, then
  add the `--add-blocked-by` edges once they all have numbers.
- **No story cites a sibling by number.** During drafting those numbers do not exist yet, and
  `#8` links to whatever issue 8 happens to be. The typed edge is the only reference.

## Engineering rules

From the charter. These are enforced in CI rather than by review. Today `Build` runs the docs
build and the rebuild module's Dagger checks: `dagger check` (the Go SDK's `dagger-go-sdk:stale`
and `rebuild:go-pipeline`) and `dagger -m rebuild/tests check --module tests`. Each other check
is added to `Build` (and to `verify`) by the story that first lands code it applies to.

- **One required check.** A workflow named `Build` (`.github/workflows/build.yaml`) runs every
  check below, and the `default` ruleset requires it. The `verify` commands in
  `.claude/backlog.json` run the same checks, so a change that passes locally passes in CI.
- **Infrastructure definitions are checked and tested.** Every module passes its tool's format
  check, validation and a linter, and has tests that run without real credentials. The tool and
  the linter are chosen in Proposal: infrastructure modules.
- **Go code is held to humus's rules.** golangci-lint runs with bedrock's forbidigo list, on test
  code too. Tests assert with `testify/require`. CI fails when a library package drops below 80%
  statement coverage. A service the platform runs that is written in Go is a humus app.
- **The docs site builds** on every pull request that touches it.
- **Every vendor and third-party component is justified** in the proposal that adds it. The
  operator must control its configuration, be able to export its data, and be able to revoke its
  credentials.
- **The rebuild module runs every exec through one of two helpers** in `rebuild/exec.go`:
  `targetExec`, which puts a fresh nonce on an exec that reads or changes a target so it is
  never cached, or `cachedExec`, for work that may cache. `.golangci.yml` bans `WithExec`
  everywhere else. A function that touches a target is also marked `+cache="never"`.
- **Dagger is pinned.** `engineVersion` in every `dagger-module.toml` and `DAGGER_VERSION` in
  `Build` move together, `dagger.lock` is committed, and `dagger generate -y` regenerates the
  bindings. devex's `daggerverse/CLAUDE.md` records the pitfalls of Dagger 1.0 beta.

## Docs site

The Hugo site lives in `docs/` and is published at https://zaba505.dev/infra/.

```bash
cd docs && hugo server        # local preview at localhost:1313
cd docs && npm install --no-save --no-audit --no-fund postcss postcss-cli autoprefixer \
  && hugo --gc --minify       # the build the Build check runs
```

Requires Hugo 0.147.0+ extended, Dart Sass and Node.js (PostCSS).

The layout is planned up front in the charter's Docs site section (revision 1). Its sections, in
order: the home page (Overview), Getting started, Fundamentals, The platform, Operating, and
Hosting a tenant. Each is a directory with an `_index.md` landing page whose `weight` sets the
order; within a section, pages are weighted in steps of 10 in the sketch's order so later pages
can be slotted in.

- **The skeleton is fixed by the charter.** Each section's landing page says what its pages will
  cover and in what order. The order comes from the skeleton, never from the order pages land in.
- **Each proposal documents its own part.** A proposal's Documentation section names the pages
  it adds or changes and where they sit. It may change the sketch when its design calls for it,
  and the charter's next revision records the change.
- **Pages are written while implementing.** Each story writes or updates the pages for what it
  lands, so a story isn't done until its pages are.
- **The site reads in order.** Each page links to the next one in the reading order, and the
  last page of a section links to the first page of the next. A page's URL follows its place in
  the sidebar.
- **Pages document the platform as it is.** The home page is an overview of the platform, not a
  log of the rebuild. No page describes the tree before the reset, or how to move off it. While
  the charter is open, the home page has one line saying the platform is being rebuilt, with a
  link to the charter; the charter's last revision removes it.
- **Every page sets `type: docs`** in its front matter explicitly.
- **Links between pages use `ref` or `relref`**, never raw paths, so a moved page fails the
  build instead of breaking silently.

## Commit & Branch Conventions

- Branch: `story/issue-{number}/{description}` or `fix/issue-{number}/{description}`
- Commit: `feat(issue-123): description` or `fix(issue-123): description`
