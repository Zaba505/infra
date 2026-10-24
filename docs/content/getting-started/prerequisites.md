---
title: What you need
type: docs
weight: 10
---

The rebuild entry point is a Dagger module in this repository. Everything it
runs, it runs in containers it puts together, so the machine that calls it needs
only two things.

## The Dagger CLI

The CLI has to be the version the module pins. That's `v1.0.0-beta.15`, the
`engineVersion` in the repository's `dagger-module.toml`. Install it with
Dagger's install script, at that version:

```sh
curl -fsSL https://dl.dagger.io/dagger/install.sh \
  | DAGGER_VERSION=1.0.0-beta.15 BIN_DIR="$HOME/.local/bin" sh
```

Check that the CLI on your `PATH` is the right one:

```sh
dagger version
```

It should report `v1.0.0-beta.15`. Dagger 1.0 is in beta, and the module
moves to a newer engine only when a pull request moves it. Install a newer CLI
then, too.

## A container runtime

The CLI runs the Dagger engine as a container, and the engine runs every step
of a run in containers of its own. Any runtime the Dagger CLI can start the
engine with works, such as Docker or Podman. The CLI starts the engine on the
first call and reuses it after that.

## To work on the module

Go is needed only to change the module's code, not to run it. The module's
`go.mod` files name the Go version.

Next: [Fundamentals]({{< ref "/fundamentals/_index.md" >}})
