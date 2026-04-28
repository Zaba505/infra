---
title: "Technical Design"
description: Technical design for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-27
---

**Parent capability:** [self-hosted-application-platform](../_index.md)

## Overview

This document composes the technical design for the self-hosted-application-platform capability from its accepted ADRs (0001 tenant state storage, 0002 contract versioning, 0003 tenant eviction export). It describes the component inventory implied by those decisions, traces each technical requirement through the ADR(s) that addresses it down to the component(s) that realize it, and surfaces requirements that no accepted ADR yet covers (gaps).

Per-component detail (interfaces, data shapes, failure modes) is intentionally out of scope here — each component listed below will get its own component-design document filed via `define-component-design`.

## Component Inventory

The following components are implied by the accepted ADRs and the existing platform topology (Cloudflare → GCP, WireGuard to home lab):

1. **Tenant State Store** — Per-tenant Firestore namespace provisioner and accessor. Realizes ADR-0001. Owns isolation boundary at the data layer (TR-01) and supports namespace-level update operations that avoid blocking online traffic (TR-04).
2. **Platform Contract Registry** — Versioned contract package distribution and resolution. Realizes ADR-0002. Hosts multiple semver-tagged contract versions concurrently and resolves which version a given tenant is bound to (TR-02).
3. **Tenant Export Service** — On-demand export to GCS bucket signed URLs. Realizes ADR-0003. Generates a portable archive of a tenant's data and a time-bounded signed URL within the export window (TR-05).
4. **Tenant Lifecycle Controller** — Coordinates tenant onboarding, contract-version migration, operator-initiated updates, and eviction handoff to the Tenant Export Service. Implied by ADR-0001/0002/0003 collectively; no single ADR owns it but every ADR's outcome flows through it. Supports TR-04 by sequencing updates without tenant-perceived downtime.
5. **Platform Edge / Network Path** — Existing Cloudflare front and GCP backends with WireGuard to home lab. Pre-existing; carries all inter-service and tenant traffic (TR-07). No new ADR; called out so the audit trail is complete.

## TR -> ADR -> Component Audit Trail

| TR | Requirement (short) | ADR(s) | Component(s) |
|---|---|---|---|
| TR-01 | Tenant data isolation | ADR-0001 | Tenant State Store |
| TR-02 | Concurrent contract versions during migration | ADR-0002 | Platform Contract Registry, Tenant Lifecycle Controller |
| TR-03 | Per-tenant observability scoped to tenant only | *(none — gap)* | *(unassigned)* |
| TR-04 | Operator updates without tenant-visible downtime | ADR-0001 (storage doesn't block updates), ADR-0002 (contract version pinning during update) | Tenant Lifecycle Controller, Tenant State Store |
| TR-05 | Evicted tenant can take their data | ADR-0003 | Tenant Export Service, Tenant Lifecycle Controller |
| TR-06 | New tenant data import without loss/corruption | *(none — gap)* | *(unassigned)* |
| TR-07 | All inter-service traffic on Cloudflare → GCP path | prior shared decision | Platform Edge / Network Path |

## Gaps

The following technical requirements are accepted but no capability-scoped ADR yet decides how they are met. Each will be filed as a follow-up issue so the gap is tracked and an ADR can be authored:

- **TR-03 — Tenant-facing observability scoping.** ADRs 0001–0003 do not cover the observability data plane. A new ADR is needed to decide how per-tenant metrics, logs, and traces are collected, stored, and scoped on query.
- **TR-06 — Tenant data import / migration in.** No ADR yet decides the import format, idempotency mechanism, or integrity-verification approach. A new ADR is needed before the Tenant Lifecycle Controller can implement onboarding for tenants that arrive with existing data.

## Follow-up Issues

One issue is filed per component (for `define-component-design`) and one per gap (for `define-adr`). See `gh-invocations.txt` for the exact `gh issue create` commands.
