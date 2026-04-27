---
title: "Technical Requirements"
description: >
    Technical requirements extracted from the self-hosted-application-platform capability and its user experiences.
type: docs
reviewed_at: null
---

> **Living document.** Numbering is append-only. ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-tech-design` skill will refuse to proceed to ADRs (Stage 2) until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [self-hosted-application-platform](_index.md)

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** The platform must enforce strict tenant isolation at the data and compute layers. No tenant workload may observe or access another tenant's state, secrets, traffic, or telemetry under any normal or degraded operating condition.
**Why this is a requirement, not a decision:** Capability business rules explicitly state tenant isolation as an invariant.

### TR-02: Operators must be able to roll out a platform-contract change without breaking existing tenants
**Source:** [UX: platform-contract-change-rollout](user-experiences/platform-contract-change-rollout.md)
**Requirement:** When the platform publishes a new contract version, existing tenants must continue to operate against the prior contract version until they migrate. The platform must support multiple contract versions concurrently for a bounded migration window.
**Why this is a requirement, not a decision:** The UX flow requires graceful contract evolution; otherwise rollouts break tenants.

### TR-03: Tenant-facing observability data must be queryable per-tenant within their data scope only
**Source:** [UX: tenant-facing-observability §Journey](user-experiences/tenant-facing-observability.md)
**Requirement:** Tenants must be able to query metrics, logs, and traces for their own workloads. Cross-tenant observability data must be inaccessible to a tenant.
**Why this is a requirement, not a decision:** Both the capability isolation invariant and the UX journey require this.

### TR-04: Operator-initiated tenant updates must complete without tenant-perceived downtime for online workloads
**Source:** [UX: operator-initiated-tenant-update §Success](user-experiences/operator-initiated-tenant-update.md)
**Requirement:** When the operator initiates an update to a tenant (config, version, or capability), tenants serving online traffic must observe no end-user-visible downtime during the update.
**Why this is a requirement, not a decision:** UX success criteria define zero downtime as the user-perceived outcome.

### TR-05: A tenant evicted from the platform must be able to take their data with them
**Source:** [UX: move-off-the-platform-after-eviction](user-experiences/move-off-the-platform-after-eviction.md)
**Requirement:** The platform must provide an export mechanism by which an evicted tenant can retrieve all of their data in a portable format within a defined export window.
**Why this is a requirement, not a decision:** UX explicitly requires the move-off journey to succeed even for evicted tenants.

### TR-06: New tenants must be able to migrate existing data into the platform without loss or corruption
**Source:** [UX: migrate-existing-data](user-experiences/migrate-existing-data.md)
**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).
**Why this is a requirement, not a decision:** UX requires lossless, retry-safe migration as part of the journey.

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** [CLAUDE.md §Architecture overview](../../../../CLAUDE.md) · prior shared decision
**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.
**Why this is a requirement, not a decision:** Inherited topology constraint from the repo's architecture; not subject to revisiting at the capability level.

## Open Questions

- Whether tenant data export (TR-05) should be on-demand or continuously-available — captured during extraction from the move-off UX.
- Whether contract versioning (TR-02) requires semver or a different versioning scheme — surfaced but deferred to Stage 2.
</content>
</invoke>