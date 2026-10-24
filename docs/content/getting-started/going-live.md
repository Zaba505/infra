---
title: Going live
type: docs
weight: 30
---

Going live is the same rebuild a drill runs, against the live target. The
commands are the same; what changes is the target, and two rules that keep a
drill from ever reaching it by mistake.

## Add the live target

Give the live target a file in the values checkout, `targets/live.yaml`:

```yaml
kind: live
```

Only one target is live.

## Commit everything first

A run against the live target refuses uncommitted changes, in this repository
or in the values checkout, so every change to the live platform is made from
history. Commit both, and check that neither has anything left over:

```sh
git status
git -C ../infra-values status
```

Otherwise `up` stops before it changes anything:

```text
a run against the live target refuses uncommitted changes: the target's values have uncommitted changes
```

## Check the target

```sh
dagger call preflight --target live
```

It changes nothing, and prints that the check passed. If the target's last run
didn't finish, or a phase that run applied no longer verifies, it fails and
says which. See [The preflight check]({{< ref "/fundamentals/rebuild.md#the-preflight-check" >}}).

## Run it

```sh
dagger call up --target live
```

As in [a drill]({{< ref "/getting-started/first-drill.md" >}}), `up` applies and
verifies one phase, stops at the checkpoint after it, and prints the run's
summary: what the phase made, how long it took and the run's time so far. Look
at what the phase made, then continue:

```sh
dagger call up --target live --resume
```

Each call runs the next phase, until the summary says the run finished and
every phase verified. That run is the target's last-known-good. A live rebuild
is meant to take at most 60 minutes, counting the time spent at checkpoints.
The summary marks a run that went over, and it's never stopped for time.

## After a failure

A failed apply or verify stops the run at once, prints what failed with the
summary, and leaves the target as it is. Look at what failed, then either:

- **resume,** once the target is fit to: commit the fix, in whichever checkout
  needs it, and run `up --resume`. It verifies every phase the run already
  verified, then applies the failed phase again and goes on.
- **tear down,** and start again with a fresh `up`.

## Tearing it down

`down` against the live target needs the target's name, to confirm it:

```sh
dagger call down --target live --confirm live
```

Without `--confirm live` it refuses:

```text
down was not confirmed with the target's name: live is the live target; confirm with "live" to tear it down
```

It refuses uncommitted changes too, as `up` does. It tears down every phase in
reverse order, phase 0 last, and it's safe to run again after it fails
partway.

Next: [Fundamentals]({{< ref "/fundamentals/_index.md" >}})
