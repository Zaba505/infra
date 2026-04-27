---
title: "Technical Requirements"
description: Technical requirements for the self-hosted-application-platform capability.
type: docs
reviewed_at: null
---

**Parent capability:** [self-hosted-application-platform](_index.md)

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** The platform must enforce strict tenant isolation at the data and compute layers. No tenant workload may observe or access another tenant's state, secrets, traffic, or telemetry under any normal or degraded operating condition.

### TR-02: Operators must be able to roll out a platform-contract change without breaking existing tenants
**Source:** [UX: platform-contract-change-rollout](user-experiences/platform-contract-change-rollout.md)
**Requirement:** When the platform publishes a new contract version, existing tenants must continue to operate against the prior contract version until they migrate. The platform must support multiple contract versions concurrently for a bounded migration window.

### TR-03: Tenant-facing observability data must be queryable per-tenant within their data scope only
**Source:** [UX: tenant-facing-observability §Journey](user-experiences/tenant-facing-observability.md)
**Requirement:** Tenants must be able to query metrics, logs, and traces for their own workloads. Cross-tenant observability data must be inaccessible to a tenant.

### TR-04: Operator-initiated tenant updates must complete without tenant-perceived downtime for online workloads
**Source:** [UX: operator-initiated-tenant-update §Success](user-experiences/operator-initiated-tenant-update.md)
**Requirement:** When the operator initiates an update to a tenant (config, version, or capability), tenants serving online traffic must observe no end-user-visible downtime during the update.

### TR-05: A tenant evicted from the platform must be able to take their data with them
**Source:** [UX: move-off-the-platform-after-eviction](user-experiences/move-off-the-platform-after-eviction.md)
**Requirement:** The platform must provide an export mechanism by which an evicted tenant can retrieve all of their data in a portable format within a defined export window.

### TR-06: New tenants must be able to migrate existing data into the platform without loss or corruption
**Source:** [UX: migrate-existing-data](user-experiences/migrate-existing-data.md)
**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** prior shared decision
**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.
