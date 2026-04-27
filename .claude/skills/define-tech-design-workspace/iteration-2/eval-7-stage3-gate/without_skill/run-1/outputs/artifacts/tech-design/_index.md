---
title: "Technical Design"
description: Composed technical design for the self-hosted-application-platform capability.
type: docs
---

**Parent capability:** [self-hosted-application-platform](../_index.md)
**Inputs:** [tech-requirements](../tech-requirements.md), [ADRs](../adrs/)

## Overview

The self-hosted application platform hosts isolated tenant workloads on the
existing Cloudflare-fronted, GCP-hosted topology (TR-07). Tenants get their
own state, observe only their own telemetry, can be updated without downtime,
and can be onboarded or evicted with their data intact.

This document composes the decisions captured in the capability's ADRs into a
single picture of the intended system. Where an ADR is still proposed, the
section is marked accordingly and reflects the current direction rather than
a settled commitment.

## Component map

| Concern | Mechanism | Source |
| --- | --- | --- |
| Tenant state storage | Per-tenant Firestore namespace | [ADR-0001](../adrs/0001-tenant-state-storage.md) (accepted) |
| Platform contract versioning | Semver in contract package path *(proposed)* | [ADR-0002](../adrs/0002-contract-versioning.md) (proposed) |
| Tenant data export on eviction | On-demand export to GCS bucket signed URL *(proposed)* | [ADR-0003](../adrs/0003-tenant-eviction-export.md) (proposed) |
| Network path | Cloudflare → GCP, WireGuard back to home lab | TR-07, prior shared decision |
| Tenant data import | Idempotent import pipeline | TR-06 (no ADR yet) |
| Tenant-scoped observability | Per-tenant query scope | TR-03 (no ADR yet) |

## Tenant isolation and state (TR-01, TR-04)

Each tenant is allocated its own Firestore namespace, per ADR-0001. Tenant
workloads address state only through a tenant-scoped client whose namespace
is fixed at provisioning time, which makes cross-tenant reads structurally
impossible rather than policy-enforced. This also gives operator-initiated
updates (TR-04) a clean unit to update against: a tenant's state lives in
exactly one namespace, so version migrations and config rollouts can be
staged per tenant without coordination across tenants.

## Platform contract versioning (TR-02) — proposed

ADR-0002 proposes encoding the contract version as semver in the contract's
package path (e.g. `contract/v1`, `contract/v2`). Multiple contract versions
are served concurrently during the bounded migration window, so existing
tenants continue against their pinned version while new tenants — or
tenants that have migrated — use the newer one. Until ADR-0002 is accepted,
the alternatives (date-based versioning, single rolling version with feature
flags) remain on the table.

## Tenant-facing observability (TR-03)

Tenants query metrics, logs, and traces only within their own data scope.
No ADR has been written for the observability stack yet; the requirement is
recorded here so a future ADR can pick it up. The tenant-scoped namespace
boundary established by ADR-0001 is the natural seam for scoping
observability queries as well.

## Operator-initiated tenant updates (TR-04)

Because each tenant's state is isolated (ADR-0001), operator-initiated
updates run per tenant and do not need to coordinate global cutovers.
Online tenants are updated using the platform's standard zero-downtime
deployment path on GCP (Cloud Run / load-balanced backends with
`create_before_destroy`), so end users see no interruption.

## Tenant eviction and data export (TR-05) — proposed

ADR-0003 proposes that, when a tenant is evicted, the platform produces an
on-demand export of their Firestore namespace into a GCS bucket and returns
a signed URL valid for the export window. Until ADR-0003 is accepted, the
alternatives (continuous replication, scheduled snapshot + manual download)
remain on the table.

## Tenant data migration in (TR-06)

No ADR has been written yet for the import path. The requirement is for an
idempotent import with verifiable integrity (no silent loss, no duplication
on retry). A future ADR should choose between, e.g., a staged import into a
shadow namespace with checksum verification before cutover versus a direct
streaming import with per-record idempotency keys.

## Network topology (TR-07)

All inter-service communication and all tenant-facing traffic traverses the
existing Cloudflare → GCP path with WireGuard back to the home lab. This is
a prior shared decision and is not re-litigated here; new platform services
must conform to it.

## Open items

- ADR-0002 (contract versioning) is proposed and needs confirmation.
- ADR-0003 (tenant eviction export) is proposed and needs confirmation.
- TR-03 (tenant-facing observability) has no ADR yet.
- TR-06 (tenant data import) has no ADR yet.

## Requirement → ADR coverage

| Requirement | Covered by |
| --- | --- |
| TR-01 | ADR-0001 (accepted) |
| TR-02 | ADR-0002 (proposed) |
| TR-03 | — (no ADR yet) |
| TR-04 | ADR-0001 (accepted) |
| TR-05 | ADR-0003 (proposed) |
| TR-06 | — (no ADR yet) |
| TR-07 | prior shared decision |
