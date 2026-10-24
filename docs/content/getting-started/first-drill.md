---
title: A first drill
type: docs
weight: 20
---

A drill is a rebuild against a scratch target, which leaves the live platform
untouched. Drill and live differ only in the target, so a drill is the same
commands a live rebuild runs. Run them from this repository's checkout, with
the values checkout set as
[What you need]({{< ref "/getting-started/prerequisites.md#the-values-checkout" >}})
says.

## Add the drill's target

Give the drill a file in the values checkout, `targets/drill.yaml`:

```yaml
kind: drill
```

The name of the file, without `.yaml`, is the target's name. A drill may run
from uncommitted changes, in either checkout, so the file doesn't need to be
committed first. The run's record says it was.

## Check the target

```sh
dagger call preflight --target drill
```

It changes nothing. Against a target that holds nothing yet, it prints:

```text
the preflight check passed against drill.
```

It fails if the target's last run didn't finish, or if a phase that run applied
no longer verifies. [The preflight check]({{< ref "/fundamentals/rebuild.md#the-preflight-check" >}})
says what each means.

## Run it

```sh
dagger call up --target drill
```

`up` runs the preflight check, then applies and verifies the first phase, and
stops at the checkpoint after it. `dagger call` prints the run's summary: each
phase the call ran, what it made and how long it took, the run's time so far,
and what you can do next. At the checkpoint after phase 0, it reads like this:

```text
run 7QH3W2Y5N4RJXK6VZ2B4TQ3MJE against drill
phase 0: verified, 2m31.4s
  made the state backend, bucket drill-state
the run's time so far: 2m40.2s
the run paused at the checkpoint after phase 0. Continue with up --resume, or tear the target down with down.
```

Look at what the phase made. Then continue with the next phase:

```sh
dagger call up --target drill --resume
```

Each call runs one phase and stops at the checkpoint after it, until the last
phase verifies and the summary says so:

```text
the run finished: every phase verified.
```

No phase has landed yet, so today the first `up` finishes at once, with
nothing to do:

```text
run JCUPWNYUR4NCJ3NCOSRCBM5DKF against drill
the run took 829.009578ms
the run finished: every phase verified.
```

Two flags change where a call stops:

- `--no-pause` runs on from one phase to the next without stopping, so each
  phase's verify is the only gate.
- `--through <phase>` stops after the phase it names, such as
  `--through 0`. The run is then unfinished, and `--resume` continues it.

## After a failure

A failed apply or verify stops the run at once. `up` fails with what failed,
followed by the summary, and leaves the target as it is for you to look at:

```text
the run failed at phase 1a's verify: <what failed>
Once the target is fit to resume, continue with up --resume, which applies phase 1a again. Or tear the target down with down.
```

From there, either:

- **resume,** once the target is fit to: fix what failed, in either checkout
  if it needs a change, and run `up --resume`. It verifies every phase the run
  already verified, then applies the failed phase again and goes on. If an
  earlier phase no longer verifies, it stops, changing nothing.
- **tear down,** with `down`, below, and start again with a fresh `up`.

`up` without `--resume` refuses to start a new run on top of an unfinished
one.

## Tear it down

```sh
dagger call down --target drill
```

`down` tears down every phase in reverse order, phase 0 last, and leaves the
target as a rebuild from nothing. It succeeds when there's nothing left, and
it's safe to run again after it fails partway.

Next: [Going live]({{< ref "/getting-started/going-live.md" >}})
