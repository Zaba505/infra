---
title: "Technical Requirements"
description: >
    Technical requirements for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-26
---

**Parent capability:** [self-hosted-application-platform](_index.md)

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** Strict tenant isolation at data and compute layers.

### TR-02: Operators must roll out platform-contract changes without breaking existing tenants
**Source:** [UX: platform-contract-change-rollout](user-experiences/platform-contract-change-rollout.md)

### TR-03: Tenant-facing observability must be queryable per-tenant only
**Source:** [UX: tenant-facing-observability](user-experiences/tenant-facing-observability.md)

### TR-04: Operator-initiated tenant updates must complete with no end-user-visible downtime
**Source:** [UX: operator-initiated-tenant-update](user-experiences/operator-initiated-tenant-update.md)

### TR-05: Evicted tenants must be able to export their data
**Source:** [UX: move-off-the-platform-after-eviction](user-experiences/move-off-the-platform-after-eviction.md)

### TR-06: Migrations must be lossless and idempotent
**Source:** [UX: migrate-existing-data](user-experiences/migrate-existing-data.md)

### TR-07: Inter-service traffic must traverse Cloudflare → GCP
**Source:** prior shared decision

### TR-08: Platform must degrade gracefully when a GCP region becomes unreachable
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** Tenant workloads must continue serving (possibly with reduced functionality) when a GCP region is unreachable for up to 30 minutes.
