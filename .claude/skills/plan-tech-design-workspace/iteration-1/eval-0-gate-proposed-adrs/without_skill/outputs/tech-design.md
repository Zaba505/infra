---
title: "Technical Design"
description: Composed technical design for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-27
---

**Parent capability:** [self-hosted-application-platform](_index.md)

> WARNING: ADR-0002 (Contract Versioning) and ADR-0003 (Tenant Eviction Export) are still in **proposed** status as of 2026-04-27. The decisions below reflect the proposed outcomes and are subject to change once those ADRs are accepted. This tech-design should be re-reviewed after acceptance.

## Overview

The self-hosted-application-platform provides a multi-tenant runtime in which operators host isolated tenant workloads on the home-lab + GCP topology fronted by Cloudflare. This document composes the accepted (and currently-proposed) architectural decisions into a single picture of the intended end state, mapping each technical requirement to the ADR that addresses it and the component(s) that will realize it.

The platform sits on the established Internet -> Cloudflare (mTLS + DDoS) -> GCP (Cloud Run) <-> Home Lab (WireGuard) path. Tenant state is stored in per-tenant Firestore namespaces; the platform contract is versioned via semver in the contract package path; evicted tenants retrieve their data via on-demand exports to a signed GCS URL.

## Component Inventory

| Component | Responsibility | Backed by ADR(s) | Satisfies TR(s) |
|---|---|---|---|
| Tenant State Store | Per-tenant Firestore namespace; enforces isolation at the data layer and supports concurrent reads/writes during no-downtime updates. | 0001 | TR-01, TR-04 |
| Contract Registry & Router | Hosts multiple concurrent contract versions (semver in package path) and routes tenant traffic to the version they are pinned to. | 0002 (proposed) | TR-02 |
| Tenant Update Orchestrator | Performs operator-initiated tenant updates without end-user-visible downtime, coordinating contract version, config, and state. | 0001, 0002 (proposed) | TR-04, TR-02 |
| Tenant Export Service | On-demand export of a tenant's Firestore namespace to a GCS bucket, returning a signed URL within the eviction window. | 0003 (proposed) | TR-05 |
| Tenant Import Service | Idempotent import path for new tenants migrating pre-existing data into their Firestore namespace with integrity verification. | (gap - no ADR) | TR-06 |
| Tenant-Scoped Observability | Per-tenant query surface for metrics/logs/traces, enforcing scope so a tenant cannot read another's telemetry. | (gap - no ADR) | TR-03 |
| Network Path Conformance | Ensures all platform <-> tenant and platform <-> platform traffic stays on Cloudflare -> GCP -> WireGuard. | (shared, prior) | TR-07 |

## TR -> ADR -> Component Audit Trail

| TR | Requirement (short) | ADR | Status | Component |
|---|---|---|---|---|
| TR-01 | Tenant isolation | 0001 | accepted | Tenant State Store |
| TR-02 | Concurrent contract versions | 0002 | **proposed** | Contract Registry & Router; Tenant Update Orchestrator |
| TR-03 | Per-tenant observability scope | - | **GAP** | Tenant-Scoped Observability |
| TR-04 | No-downtime operator updates | 0001 | accepted | Tenant State Store; Tenant Update Orchestrator |
| TR-05 | Evicted tenant data export | 0003 | **proposed** | Tenant Export Service |
| TR-06 | Tenant data migration in | - | **GAP** | Tenant Import Service |
| TR-07 | Cloudflare -> GCP path | (shared) | n/a | Network Path Conformance |

## Surfaced Gaps

The following technical requirements are not yet covered by an ADR in this capability and must be resolved before the tech-design can be finalized:

- **TR-03 (Tenant-Scoped Observability)** - no ADR exists choosing how per-tenant query scoping is enforced (Cloud Logging tenant labels vs. dedicated per-tenant log buckets vs. proxying through a tenant-aware query gateway).
- **TR-06 (Tenant Data Migration In)** - no ADR exists choosing the import mechanism (resumable upload + import job vs. direct Firestore writes via signed credentials vs. a staged-bucket + validator pipeline).

## Open Items Before Acceptance

1. ADR-0002 and ADR-0003 must move from **proposed** to **accepted**; this design should be re-validated against any changes to their decision outcomes.
2. ADRs covering TR-03 and TR-06 must be authored and accepted (see Surfaced Gaps).
3. Per-component design documents (one per row in the Component Inventory) should be filed as follow-up work.
