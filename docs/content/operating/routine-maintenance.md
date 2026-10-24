---
title: Routine maintenance
type: docs
weight: 20
---

## Keeping Dagger current

The rebuild entry point is a Dagger module, and three things pin the Dagger it
runs on:

- **`engineVersion`**, in every `dagger-module.toml`: the repository root's,
  for the module, and `rebuild/tests/dagger-module.toml`, for its tests. It
  pins the engine.
- **`DAGGER_VERSION`**, in `.github/workflows/build.yaml`: the CLI that `Build`
  installs.
- **`dagger.lock`**: every image and git ref the workspace uses, by digest or
  commit. That includes the Go SDK that writes the module's generated code, and
  devex's modules, which the module's checks run.

Nothing moves any of them except a pull request.

### Moving the engine

Move `engineVersion` and the CLI pin in `Build` together, in one pull
request. `Build` fails when they differ. Then install the same CLI version
locally, as in [What you need]({{< ref "/getting-started/prerequisites.md" >}}),
and regenerate the module's code:

```sh
dagger generate -y
```

Commit whatever it changes. The Go SDK's check, `dagger-go-sdk:stale`, fails
`Build` when the committed code isn't what `dagger generate` writes.

### Moving the pins in `dagger.lock`

`dagger lock update` resolves the entries again, moving tags and branches to
what they point at now:

```sh
dagger lock update            # every entry
dagger lock update --list     # the selectors it knows
dagger lock update devex      # just devex's modules
```

It regenerates the module's code as well. Run the checks, then commit
`dagger.lock` with what changed:

```sh
dagger check
dagger -m rebuild/tests check --module tests
```

The engine adds an entry for anything it resolves that `dagger.lock` doesn't
record. `Build` fails when a run leaves `dagger.lock` changed, so a new image
or module is committed with its pin.

### Why this matters more while Dagger is in beta

Dagger 1.0 is in beta, and a beta can change between releases: a module's
API, what the Go SDK generates, how a check runs. Pinned, none of that reaches
a rebuild until a pull request brings it in, and `Build` runs the module's
checks on it first. Moving often, a release at a time, keeps each change small
enough to see what broke. Left for months, the move becomes several breaking
releases at once, and a rebuild is the worst time to find that the pinned
version is no longer published.

Next: [Hosting a tenant]({{< ref "/hosting-a-tenant/_index.md" >}})
