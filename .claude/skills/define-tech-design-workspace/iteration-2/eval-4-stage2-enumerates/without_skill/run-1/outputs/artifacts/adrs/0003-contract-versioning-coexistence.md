---
title: "[0003] Platform-Contract Versioning and Multi-Version Coexistence"
description: >
    Choose the versioning scheme for the platform contract and the mechanism by which two contract versions run side-by-side during a tenant's migration window.
type: docs
weight: 3
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

[TR-02] requires that when the platform publishes a new contract version, existing tenants continue to operate against the prior version until they migrate, with multiple contract versions concurrently supported for a bounded migration window. The "platform contract" here is the surface that tenants declare against: packaging form, declared-needs schema, observability/secret/identity offering shapes, and ingress conventions.

Two questions are open: what versioning scheme does the contract use, and how does the platform run two versions side-by-side without forcing every tenant to upgrade in lockstep?

## Decision Drivers

* TR-02: bounded multi-version coexistence is a hard requirement
* C-07 (long-form): contract is evergreen — tenants do not re-accept on each modification; operator-driven rollouts
* NFR-02: maintenance budget — N versions in flight forever would consume the budget; the bound matters
* The capability has very few tenants (single operator, single owner) — the scheme must scale *down* as gracefully as it scales out
* Reproducibility (NFR-01): version transitions must be reproducible from the definitions repo

## Considered Options

* **Semver (`MAJOR.MINOR.PATCH`) on the whole contract, with N=2 majors supported** — classic semver; "breaking" means a new MAJOR; one-major-back coexistence window.
* **Date-stamped contract versions (`YYYY-MM`), one current + one prior supported** — calendar-based, no semantic claims; coexistence window is the gap between the prior version and the migration deadline.
* **Per-offering versioning (compute-v2, observability-v1, …), tenant declares each independently** — each offering versions itself; no unified "platform version."

## Decision Outcome

Chosen option: **date-stamped contract versions (`YYYY-MM`), one current + one prior supported**.

A date stamp avoids the "is this a MAJOR or a MINOR" debate that semver invites and that, in a single-operator system, has no second person to push back on. The bound — "one current + one prior" — gives every tenant a known migration window that ends when the operator publishes the *next* version (which retires the oldest). The platform always runs exactly two contract versions concurrently; the migration window is the time between two consecutive publications.

The two versions coexist by being *separate sets of platform-side adapters* (ingress route templates, observability label conventions, secret-resolution paths) keyed on the tenant's declared `contract: 2026-04` field. A tenant runs against one and only one version at a time; the platform-side adapter for that version translates the tenant's declarations into the underlying substrate primitives (ADR-0002). Per-offering versioning was rejected because it multiplies the coexistence matrix beyond what one operator can hold in their head.

### Consequences

* Good, because the bound is structural: there are always exactly two versions, never N
* Good, because date stamps communicate "this was the contract as of April 2026" without making semantic claims the operator cannot enforce
* Good, because the migration deadline is predictable: it is the date of the *next* publication
* Good, because per-offering surfaces (compute, observability, secrets) sit behind a single per-version adapter — no combinatorial coexistence
* Neutral, because every tenant must declare its contract version explicitly in its tech design
* Bad, because date stamps lose the "is this breaking?" signal that semver provides — mitigated by a CHANGELOG section in each contract version's docs
* Bad, because tenants on the prior version are on a clock — the operator must communicate the next-publication date clearly (handled by the `platform contract change` issue type, REQ-21 long-form)

### Confirmation

* Each contract version is a directory in the infra repo: `contracts/2026-04/`, `contracts/2026-09/`, etc., containing the schema, the adapter code, and the migration notes
* The Nomad-side per-tenant job spec includes the `contract` label; the platform's adapter selection is keyed on it
* When a new contract version is published, CI fails the build if more than two contract directories are marked `status: supported`
* A `platform contract change` rollout issue (REQ-06 long-form) is opened on publication; it is closed only when every tenant on the now-prior version has migrated or been evicted

## Pros and Cons of the Options

### Semver (`MAJOR.MINOR.PATCH`)

* Good, because it is industry-standard
* Good, because "MAJOR" carries an unambiguous "breaking" signal
* Bad, because semver requires a community of consumers to be useful; with one operator the MAJOR/MINOR call is unilateral and tends to drift
* Bad, because PATCH is meaningless for a contract document — the contract either is or is not in effect

### Date-stamped (`YYYY-MM`), one current + one prior

* Good, because the bound is structural and easy to reason about
* Good, because the publication cadence is the migration cadence
* Good, because a single operator can hold two versions in their head, never N
* Neutral, because tenants must declare their version explicitly
* Bad, because no built-in "breaking" signal — mitigated by per-version CHANGELOG

### Per-offering versioning

* Good, because tenants can adopt one offering's new version without touching others
* Bad, because the coexistence matrix is the cartesian product of per-offering versions — explodes fast
* Bad, because the operator must reason about every combination during incidents
* Bad, because rollouts are no longer "the platform contract changed" but "compute-v2 + observability-v1 + secrets-v3 became compatible" — incoherent

## More Information

* The migration window for a given tenant is bounded by the next publication; the operator commits to a minimum window length (initially 90 days) in the contract docs
* Contract changes that affect TR-07 (Cloudflare → GCP path) are out of scope for this versioning scheme — those are cross-cutting and handled in `docs/content/r&d/adrs/`
