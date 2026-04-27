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
**Source:** [UX: migrate-existing-data](user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists)
> ⚠️ source no longer resolves — human review. The anchor `#a-section-that-no-longer-exists` does not exist in `user-experiences/migrate-existing-data.md`. The UX page itself still exists and still describes a one-shot migration job (see §Journey, §Success). A human should pick the correct anchor (likely §Journey or §Success) or confirm the page-level link is sufficient. The TR text is preserved as-is per append-only policy; do not delete or renumber.
**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).
**Why this is a requirement, not a decision:** UX requires lossless, retry-safe migration as part of the journey.

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** [CLAUDE.md §Architecture overview](../../../../CLAUDE.md) · prior shared decision
**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.
**Why this is a requirement, not a decision:** Inherited topology constraint from the repo's architecture; not subject to revisiting at the capability level.

### TR-08: The platform must be reproducible from its definitions in at most one hour
**Source:** [Capability §Success Criteria & KPIs](_index.md#success-criteria--kpis) · [UX: stand-up-the-platform §Journey](user-experiences/stand-up-the-platform.md)
**Requirement:** Starting from no platform at all, the platform — including its foundations, core services, and cross-cutting services — must be able to be stood up from version-controlled definitions within one hour of wall-clock time. No step in the stand-up journey may require manual snowflake configuration that is not captured as a definition.
**Why this is a requirement, not a decision:** The capability's *Reproducibility* KPI sets the 1-hour bound as the operational test of "reproducible," and the stand-up UX is the journey that exercises it.

### TR-09: The platform's identity offering must be capable of honoring a "lost credentials cannot be recovered" property
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** Whatever identity service the platform offers to tenants must be capable of operating in a mode where a tenant's end-user credentials, once lost, cannot be recovered or reset by the operator. This property must be available to tenants that require it (e.g. self-hosted personal media storage), even if other tenants do not enable it.
**Why this is a requirement, not a decision:** The capability business rule explicitly excludes any identity option that cannot honor this property.

### TR-10: The platform must support a designated successor operator with sealed/escrowed credentials
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** A designated successor must be able to take over operating the platform if the primary operator becomes unavailable. The credentials and runbook necessary for takeover must be held by the successor in sealed/escrowed form, must not be used for routine operation, and must enable continuity of the platform itself (not merely export of tenant data).
**Why this is a requirement, not a decision:** Capability business rule on operator succession explicitly requires both sealed-credential successor takeover and the on-demand export mechanism (the latter is covered by TR-05).

### TR-11: The platform must accept tenant components packaged in a single declared form, with declared resource needs
**Source:** [UX: host-a-capability §Journey](user-experiences/host-a-capability.md) · [Capability §Business Rules](_index.md#business-rules)
**Requirement:** Tenant capabilities must be onboarded in a single, well-defined packaging form that the platform accepts, accompanied by an up-front declaration of compute, storage, and network reachability needs. The same packaging form applies to one-shot components such as migration jobs (see TR-06). The platform must not accept tenant components that bypass this contract.
**Why this is a requirement, not a decision:** The capability's "tenants must accept the platform's contract" rule and the host-a-capability UX both require a single declared packaging contract; multiple ad-hoc forms would break reproducibility (TR-08) and the operator's maintenance budget.

## Open Questions

- Whether tenant data export (TR-05) should be on-demand or continuously-available — captured during extraction from the move-off UX.
- Whether contract versioning (TR-02) requires semver or a different versioning scheme — surfaced but deferred to Stage 2.
- TR-06 source link uses an anchor (`#a-section-that-no-longer-exists`) that no longer resolves in the migrate-existing-data UX. The TR itself is still supported by the UX page as a whole; the human reviewer should select the correct anchor (e.g. `#journey` or `#success`) or accept a page-level link.
</content>
