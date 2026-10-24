---
title: The rebuild and its phases
type: docs
weight: 20
---

The platform is rebuilt, not repaired. One module in this repository stands it
up from nothing, one phase at a time, against one named target, and tears it
down again in reverse.

## The command

The rebuild is three functions of the `rebuild` Dagger module, which
`dagger.toml` makes the workspace's entry point. Run them from this
repository's checkout:

```sh
dagger call up        --target <name> [--resume] [--through <phase>] [--no-pause]
dagger call down      --target <name> [--confirm <name>]
dagger call preflight --target <name>
```

- **`up`** runs the preflight check, then applies and verifies phases in
  order. It stops at a checkpoint after each phase. `--no-pause` runs on
  without stopping, and `--through` names the last phase it runs. With
  `--resume`, it continues the target's unfinished run instead of starting a
  new one.
- **`down`** tears down every phase in reverse order, phase 0 last. It's safe
  to run again after it fails partway, and it succeeds when there's nothing
  left.
- **`preflight`** runs only the preflight check, and changes nothing.

There's no default target: each function refuses to run without `--target`.
None of them is ever served from Dagger's cache, and neither is any command a
phase runs against a target, so every call really reads and changes the
target it names.

## Targets

A target is the one place a run is pointed at, such as the live platform or a
scratch drill. Drill and live differ only in the target. Each has a file in
the values checkout, `targets/<name>.yaml`, which says which it is:

```yaml
kind: drill
```

`kind` is `live` or `drill`. The file is read strictly: a key it shouldn't
have, a second YAML document or any other kind is refused, and so is a target
with no file. The rest of the file, the values each phase reads, comes with
phase 0.

Only one target is live. Every call reads every target's file, as the values
checkout's commit has it and as the checkout stands, and refuses to run when
more than one target is live. A target is live when either version of its
file says so, so an uncommitted edit can't turn the live target into a drill.

The values checkout is a git repository outside this one. A Dagger function
can't read the machine it's called from, so the module takes the checkout as
an argument, which the operator sets once as a user-level Dagger setting; see
[What you need]({{< ref "/getting-started/prerequisites.md#the-values-checkout" >}}).

Three rules keep a drill off the live platform:

- **There's no default target.**
- **A run against the live target refuses uncommitted changes,** in this
  repository or in the values checkout, so every change to the live platform
  was made from history. A drill may run from uncommitted changes, and its
  record says so.
- **`down` against the live target needs `--confirm <name>`.** A drill
  doesn't, but refuses a name that isn't its own.

## The phase order

Each phase is one part of the platform. The order is fixed in one place,
`rebuild/phases.go`:

| Phase | Part |
|---|---|
| 0 | the state backend |
| 1a | the cloud anchor and the edge |
| 1b | the operations tunnel |
| 1c | home-lab hosts |
| 2 | tenant compute, then the offerings tenants use |
| 3 | backup, and observability |
| 4 | the canary tenant |

A phase joins the order when its part lands, and a rebuild runs every phase
the order holds. No phase has landed yet, so the order is empty: `up` against
any target passes the preflight check as a rebuild from nothing, and finishes
with nothing to do.

## The phase contract

Every phase meets one interface, in `rebuild/run/`:

- **`Apply`** makes the target match the phase's definitions.
- **`Verify`** fails unless the target is what the phase last applied, and it
  works.
- **`Teardown`** removes everything `Apply` made. It succeeds when there's
  nothing left to remove, including when `Apply` never ran.

Each runs against one target, and each is safe to run again after it fails
partway. A phase reads every value it needs from the target, and nothing from
the machine the run started on.

## Checkpoints

A Dagger function runs to completion, so a checkpoint is the end of a call,
not a prompt inside one.

After each phase verifies, `up` writes the run's record and returns a summary,
which `dagger call` prints: what the phase made, how long it took, and the
run's time so far. Then it stops. The operator chooses:

- **continue,** with `up --resume`, which runs the next phase;
- **tear down,** with `down`.

With `--no-pause`, as in CI and in drills, `up` doesn't stop between phases.
Each phase's verify is then the only gate.

A failed apply or verify stops the run at once. `up` writes the record and
fails with what failed, followed by the summary, and leaves the target as it
is to look at.

## Resuming

A run that stopped, at a checkpoint, where `--through` said, or on a failure,
is **unfinished**. `up --resume`:

1. stops unless the target's last run is unfinished;
2. runs `Verify` for every phase the run has verified so far, and stops on any
   failure, because that's drift;
3. continues from the first phase that didn't verify, starting with its
   `Apply`.

Whether a target is fit to resume is the operator's call. A failure can leave
a phase in a state its `Apply` can't fix. Resuming then fails again, and
`down` followed by a fresh `up` is the way back.

A resumed run is the same run, with the same ID. It may resume from a
different commit than it started from, such as one that fixes the phase that
failed, and its record keeps the commits each phase was applied from.

## The preflight check {#the-preflight-check}

Before `up` applies anything, it checks the target. `preflight` runs the same
check on its own.

1. **If the target has no state backend,** this is a rebuild from nothing, and
   the check passes.
2. **If the target's last run is unfinished,** the check stops. `up` never
   starts a new run on top of one: resume it with `up --resume`, or tear the
   target down with `down`.
3. **Otherwise it runs `Verify` for every phase the last run applied.** Any
   failure is drift, and the check stops before anything changes.

Drift is a difference between the target and what the rebuild last applied to
it. A difference between what was last applied and today's definitions isn't
drift: it's the change this run applies.

## Records

Each run keeps a record in the target's state backend, which phase 0 makes. It
holds:

- the commits of this repository and of the values checkout that each phase
  was applied from, and whether either had uncommitted changes;
- each phase it applied, what the phase made and how long it took;
- where it stopped, at a checkpoint, where `--through` said or on a failure,
  or that it finished.

The record of the last run that finished is the target's **last-known-good**.
The record lives only in the target, and `down` removes it with phase 0. Until
phase 0 lands, no target has a state backend, so no run is recorded and every
`up` is a rebuild from nothing.

## Time

A rebuild is meant to take at most 60 minutes. A run's time is measured from
the start of its first `up` to the last phase's verify. It includes the
preflight check, the time between calls at every checkpoint, and the time
between a failure and its resume. The summary shows each phase's time and the
run's, and marks a run over 60 minutes, and the record keeps both.

A run is never stopped for time. A live rebuild that's running long is still
the fastest way back.

Next: [The platform]({{< ref "/platform/_index.md" >}})
