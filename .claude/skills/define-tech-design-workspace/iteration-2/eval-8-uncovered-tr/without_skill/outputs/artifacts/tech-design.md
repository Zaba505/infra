---
title: "Technical Design"
description: >
    Composed technical design for the self-hosted-application-platform capability,
    derived from the accepted capability-scoped ADRs and the technical requirements.
type: docs
reviewed_at: 2026-04-26
---

**Parent capability:** [self-hosted-application-platform](_index.md)
**Inputs:** [tech-requirements.md](tech-requirements.md), [adrs/](adrs/)

> Status: **Incomplete.** Of TR-01 through TR-08, only TR-01, TR-02, TR-04, and
> TR-05 are covered by an accepted ADR. TR-03, TR-06, TR-07, and TR-08 have no
> capability-scoped decision yet and are intentionally left as open items below
> rather than designed-in-place. The composed design will need to be revised once
> those decisions land.

## 1. Overview

The self-hosted application platform hosts isolated tenant workloads on top of
the home-lab + GCP footprint. This document composes the accepted ADRs into a
single human-readable description of the intended end state.

The design is organized by the requirement(s) each piece addresses.

## 2. Requirement-to-decision coverage

| Requirement | Covered by | Status |
|-------------|------------|--------|
| TR-01 Tenant isolation                       | [ADR-0001](adrs/0001-tenant-state-storage.md)     | Designed |
| TR-02 Contract change rollout                | [ADR-0002](adrs/0002-contract-versioning.md)      | Designed |
| TR-03 Per-tenant observability               | _none_                                            | **Open** |
| TR-04 No-downtime operator-initiated updates | [ADR-0001](adrs/0001-tenant-state-storage.md)     | Designed (storage piece only) |
| TR-05 Evicted-tenant data export             | [ADR-0003](adrs/0003-tenant-eviction-export.md)   | Designed |
| TR-06 Lossless, idempotent migrations        | _none_                                            | **Open** |
| TR-07 Inter-service traffic via Cloudflare → GCP | _none capability-scoped_ (referenced as a prior shared decision) | **Open at this layer** |
| TR-08 Graceful degradation on regional outage | _none_                                           | **Open** |

## 3. Designed components

### 3.1 Tenant state storage (TR-01, TR-04)
Per [ADR-0001](adrs/0001-tenant-state-storage.md), each tenant gets its own
Firestore namespace.

- **Isolation (TR-01):** namespace boundary is the trust boundary. Service code
  resolves the namespace from the authenticated tenant identity on every
  request; cross-namespace reads are not a supported code path.
- **No-downtime updates (TR-04, storage layer):** because each tenant's data
  lives in its own namespace, schema evolution and operator-initiated updates
  can be staged tenant-by-tenant without touching the shared dataset. Note this
  ADR only covers the **storage** contribution to TR-04; the rollout
  orchestration that actually delivers no-downtime updates is not yet decided
  (see §4).

### 3.2 Platform contract versioning (TR-02)
Per [ADR-0002](adrs/0002-contract-versioning.md), the platform contract is
versioned via semver in the contract's package path (e.g. `.../contract/v1`,
`.../contract/v2`).

- Multiple contract versions can be served concurrently from the same control
  plane.
- Tenants opt into a new major version explicitly; minor/patch versions are
  backward-compatible by definition and can be rolled out without tenant
  action.
- Deprecation of an old major version is a separate operator workflow that
  references this versioning scheme.

### 3.3 Evicted-tenant data export (TR-05)
Per [ADR-0003](adrs/0003-tenant-eviction-export.md), an evicted tenant
retrieves their data via an on-demand export to a GCS bucket, surfaced as a
signed URL.

- Export is initiated by the eviction workflow (not by the evicted tenant
  directly).
- The signed URL has a finite TTL; the operator-facing eviction UX is
  responsible for delivering it to the tenant.
- The export format is the tenant's Firestore namespace contents; no
  re-shaping into a "portable" schema is implied by the ADR.

## 4. Open items (uncovered requirements)

The following requirements have **no accepted capability-scoped ADR** and
therefore have no design in this document. They must be decided before the
capability can be considered designed.

### 4.1 TR-03 — Per-tenant observability
No decision on how telemetry is partitioned, stored, or queried such that a
tenant can only see their own signals. Open questions include: signal pipeline
(push vs. pull), per-tenant tenancy in the observability backend, query-time
authorization, and retention.

### 4.2 TR-06 — Lossless, idempotent migrations
No decision on the migration mechanism. Per-tenant Firestore namespaces
(§3.1) constrain the option space (migrations are per-namespace), but the
runner, idempotency strategy (e.g. migration ledger per namespace), and
verification approach are undecided.

### 4.3 TR-07 — Cloudflare → GCP inter-service traffic
The tech-requirements document attributes this to a "prior shared decision."
No capability-scoped ADR re-affirms how this capability's services are wired
through that path (ingress hostnames, mTLS identity per tenant vs. per
service, egress from GCP back to home lab). This needs either a pointer to
the shared ADR or a capability-scoped ADR layered on top of it.

### 4.4 TR-08 — Graceful degradation on regional outage
No decision. The 30-minute single-region-unreachable budget is not addressed
by any current ADR. Open questions: multi-region replication of the
per-tenant Firestore namespaces (which today is single-region by default),
control-plane availability during a regional outage, what "reduced
functionality" is acceptable per tenant tier, and whether home-lab fallback
is in scope.

### 4.5 TR-04 — Update rollout (orchestration half)
ADR-0001 covers the storage shape that *enables* per-tenant updates, but the
actual no-downtime rollout mechanism (drain, dual-run, cutover, rollback) is
not decided. Listed here as a partial gap.

## 5. Cross-cutting notes

- **Trust boundary:** the per-tenant Firestore namespace is currently the
  only enforced isolation primitive in the design. Network-, compute-, and
  observability-layer isolation are not yet specified.
- **Operator surface area:** three workflows are implied by accepted ADRs —
  contract version rollout (TR-02), tenant update (TR-04, partial), and
  tenant eviction + export (TR-05). Their concrete operator UX is described
  in the linked user-experience documents but the orchestration services
  behind them are not yet designed.

## 6. Next steps

1. Author ADRs for TR-03, TR-06, TR-07, and TR-08.
2. Author a follow-up ADR (or extend ADR-0001) covering the rollout
   orchestration half of TR-04.
3. Re-compose this document once the above are accepted; remove the "Open
   items" section and fold the new decisions into §3.
