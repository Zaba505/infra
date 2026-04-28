---
title: "Technical Design"
description: >
    Composed technical design for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-27
---

**Parent capability:** [self-hosted-application-platform](../_index.md)

## Overview

This document composes the accepted ADRs for the self-hosted-application-platform
capability into a single human-friendly description of the final state. It is a
skeleton: each component listed below will be detailed in its own component
design document, filed as a separate issue.

The platform hosts tenant capabilities on operator-controlled infrastructure
spanning Cloudflare (edge) and GCP (compute, storage, identity), per the prior
shared decision that inter-service traffic traverses Cloudflare → GCP (TR-07).
Tenant state is isolated at the data layer via per-tenant Firestore namespaces
(ADR 0001), the platform contract is versioned via semver in the package path
(ADR 0002), and evicted tenants retrieve their data via on-demand GCS
signed-URL exports (ADR 0003).

## Component Inventory

The platform decomposes into the following components. Each will receive its
own component-design document.

1. **Tenant State Store** — per-tenant Firestore namespaces. Provides isolated
   persistent storage for each tenant. Owns TR-01 (isolation) and contributes
   to TR-04 (no-downtime updates).
   *Source: ADR 0001.*

2. **Platform Contract Library** — semver-versioned contract package(s) that
   tenants depend on. Allows multiple contract versions to coexist so operators
   can roll forward without breaking already-running tenants. Owns TR-02.
   *Source: ADR 0002.*

3. **Tenant Export Service** — on-demand exporter that materializes a tenant's
   data into a GCS bucket and returns a signed URL. Triggered during eviction
   and by the operator-succession archive flow. Owns TR-05.
   *Source: ADR 0003.*

4. **Edge Ingress (Cloudflare → GCP)** — the existing Cloudflare proxy and
   mTLS-secured GCP HTTPS load balancer that fronts every tenant. Owns TR-07.
   *Source: prior shared decision; no capability-scoped ADR — see Gaps below.*

## TR → ADR → Component Audit Trail

| TR    | Requirement (short)                                  | ADR       | Component                |
| ----- | ---------------------------------------------------- | --------- | ------------------------ |
| TR-01 | Strict tenant isolation                              | 0001      | Tenant State Store       |
| TR-02 | Contract changes without breaking tenants            | 0002      | Platform Contract Library |
| TR-03 | Per-tenant-only observability queries                | **(none)** | **(gap — see below)**    |
| TR-04 | Operator updates with no end-user-visible downtime   | 0001      | Tenant State Store       |
| TR-05 | Evicted tenants can export their data                | 0003      | Tenant Export Service    |
| TR-06 | Migrations are lossless and idempotent               | **(none)** | **(gap — see below)**    |
| TR-07 | Inter-service traffic traverses Cloudflare → GCP     | (prior shared decision) | Edge Ingress |
| TR-08 | Graceful degradation when a GCP region is unreachable | **(none)** | **(gap — see below)**    |

## Gaps

The following technical requirements have no accepted ADR addressing them.
This tech design cannot be considered complete until each is resolved (either
by a new ADR plus component, or by an explicit decision to defer).

- **TR-03 (per-tenant observability):** No ADR defines where tenant telemetry
  is stored, how queries are scoped to a single tenant, or how the operator
  views cross-tenant aggregates without leaking per-tenant data into a tenant
  view. A capability-scoped ADR is required.

- **TR-06 (lossless, idempotent migrations):** No ADR defines the migration
  framework — how a migration is described, how idempotency is enforced, how
  loss is detected, or how a partially-applied migration is rolled back. A
  capability-scoped ADR is required.

- **TR-08 (graceful regional degradation, ≤30 min unreachable):** No ADR
  defines the multi-region posture (active/active, active/passive, single
  region with read-only fallback, etc.), how Firestore namespace data is
  reachable from a surviving region, how the Edge Ingress fails over, or what
  "reduced functionality" means concretely. A capability-scoped ADR is
  required before this design is complete.

Each gap will be filed as a GitHub issue (see `gh-invocations.txt`) so it can
be picked up by `define-adr`.

## Open Questions

- Resolution of the three gaps above.
- Component-design issues (one per component) will be filed once the gaps are
  closed and the inventory is final.
