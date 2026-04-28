---
title: "Technical Design"
description: Tech design for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-27
---

**Parent capability:** [self-hosted-application-platform](../_index.md)

## Overview

This document composes the accepted ADRs for the self-hosted-application-platform capability into a skeleton technical design. It describes the final-state shape of the platform at the component level and provides a TR -> ADR -> component audit trail. Per-component detail (API surfaces, schemas, internal data flows, deployment topology) is intentionally NOT inlined here; each component listed below has its own component-design document filed via `define-component-design`.

The platform exists to give the operator's other capabilities a default, reproducible, vendor-independent place to run. The design state below reflects only what the accepted ADRs force; gaps are listed explicitly so they can be addressed before implementation.

## Accepted ADRs

- [0001 Tenant State Storage](adrs/0001-tenant-state-storage.md) — per-tenant Firestore namespace (TR-01, TR-04)
- [0002 Contract Versioning](adrs/0002-contract-versioning.md) — semver in contract package path (TR-02)
- [0003 Tenant Eviction Export](adrs/0003-tenant-eviction-export.md) — on-demand export to GCS bucket signed URL (TR-05)

## Component Inventory

Each component below is a separate unit of design. Detailed design (endpoints, schemas, internal logic) lives in the component's own design document, not here.

### Tenant Registry
- **Purpose:** System-of-record for tenants and the per-tenant Firestore namespace assignment that ADR-0001 requires.
- **Forced by:** ADR-0001 (per-tenant Firestore namespace), ADR-0002 (must record each tenant's bound contract version).
- **Status:** Component design pending — see component-design issue.

### Contract Version Resolver
- **Purpose:** Resolves which contract version a given tenant is operating against, supporting concurrent versions during a migration window.
- **Forced by:** ADR-0002 (semver in contract package path), TR-02.
- **Status:** Component design pending — see component-design issue.

### Tenant Export Service
- **Purpose:** Produces an on-demand export of a tenant's data to a GCS bucket and returns a signed URL within the export window.
- **Forced by:** ADR-0003, TR-05.
- **Status:** Component design pending — see component-design issue.

### Tenant State Store (Firestore)
- **Purpose:** The per-tenant Firestore namespaces themselves; not a service this capability builds, but a managed dependency that the design relies on.
- **Forced by:** ADR-0001.
- **Status:** Configuration captured in the Tenant Registry component design.

## TR -> ADR -> Component Audit Trail

| TR | ADR(s) | Component(s) |
|----|--------|--------------|
| TR-01 (tenant isolation) | 0001 | Tenant Registry, Tenant State Store |
| TR-02 (contract versioning) | 0002 | Contract Version Resolver, Tenant Registry |
| TR-03 (per-tenant observability) | — (gap) | — |
| TR-04 (no-downtime tenant updates) | 0001 (partial) | Tenant State Store; **gap: update orchestration** |
| TR-05 (evicted-tenant data export) | 0003 | Tenant Export Service |
| TR-06 (data import for new tenants) | — (gap) | — |
| TR-07 (Cloudflare -> GCP path) | — (prior shared decision, not capability-scoped) | All exposed components inherit |

## Surfaced Gaps

The following requirements are not yet covered by an accepted capability-scoped ADR. Each will be filed as a separate issue so it can be planned via `plan-adrs` and decided via `define-adr` before this design is implementable.

1. **TR-03 — per-tenant observability.** No ADR yet selects the observability stack or the per-tenant scoping mechanism.
2. **TR-04 — no-downtime update orchestration.** ADR-0001 covers the storage shape but not how rollouts are orchestrated to avoid tenant-perceived downtime.
3. **TR-06 — tenant data import.** No ADR yet selects the import mechanism or integrity-verification approach.

## Out of Scope for This Document

- Per-component API surfaces, request/response schemas, and internal logic — captured in each component's own design document.
- Implementation (code, Terraform) — follows component design.
- Cross-capability shared decisions (e.g. the Cloudflare -> GCP topology in TR-07) — those live in `docs/content/r&d/adrs/`.
