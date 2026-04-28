---
title: "Technical Requirements"
description: Technical requirements for the self-hosted-application-platform capability.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-tech-design` / `plan-adrs` skill will refuse to proceed to ADRs (Stage 2) until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [self-hosted-application-platform](_index.md)

## How to read this

Each requirement is **forced** by the capability or a user experience — it constrains what the system must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review.

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [Capability §Business Rules](_index.md#business-rules--constraints)
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
> ⚠️ source no longer resolves — human review

**Source:** [UX: migrate-existing-data](user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists)
**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).
**Why this is a requirement, not a decision:** UX requires lossless, retry-safe migration as part of the journey. The original anchor (`#a-section-that-no-longer-exists`) no longer resolves — the section was renamed. Likely re-source candidates: [UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#journey) or [UX: migrate-existing-data §Edge Cases & Failure Modes](user-experiences/migrate-existing-data.md#edge-cases--failure-modes).

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** [CLAUDE.md §Repository Overview](../../../../CLAUDE.md) · prior shared decision
**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.
**Why this is a requirement, not a decision:** Inherited topology constraint from the repo's architecture; not subject to revisiting at the capability level.

### TR-08: The platform must be rebuildable from definitions in at most one hour, with no manual snowflake configuration
**Source:** [Capability §Success Criteria & KPIs](_index.md#success-criteria--kpis) · [UX: stand-up-the-platform §Goal](user-experiences/stand-up-the-platform.md#goal)
**Requirement:** Starting from no platform at all, the platform must be (re)provisioned end-to-end from its versioned definitions within a 1-hour wall-clock budget. No step in the rebuild may require manual configuration that is not captured in the definitions repo.
**Why this is a requirement, not a decision:** The capability's *Reproducibility* KPI sets the 1-hour bound and the rebuild-from-definitions invariant. The stand-up UX is the journey the KPI is measured against.

### TR-09: Every phase of platform provisioning must be cleanly reversible (full teardown must be a viable rollback)
**Source:** [UX: stand-up-the-platform §Constraints Inherited from the Capability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [UX: stand-up-the-platform §Edge Cases](user-experiences/stand-up-the-platform.md#edge-cases--failure-modes)
**Requirement:** The platform's definitions must support deleting all state provisioned within any single phase (and across all phases) reliably, so that "tear down everything and restart" is always a viable response to a failed rebuild phase. Partial state must not be retained or trusted across a restart.
**Why this is a requirement, not a decision:** The stand-up UX explicitly forces "delete everything and start over" as a guaranteed option at every checkpoint, in service of reproducibility honesty.

### TR-10: Platform readiness must be verifiable end-to-end by exercising a tenant deployment, not by self-checks alone
**Source:** [UX: stand-up-the-platform §Phase 4](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant)
**Requirement:** "Ready to host tenants" must be demonstrated by deploying, exercising, and tearing down a purpose-built canary tenant that uses every platform-provided service (compute, storage, identity, network reachability, backup, observability). Readiness must not be declared from infrastructure self-checks alone.
**Why this is a requirement, not a decision:** The UX makes canary success the readiness signal; the capability's *Default hosting target* outcome forces end-to-end host-a-tenant proof rather than infrastructure-level proof.

### TR-11: The platform's identity offering must be capable of honoring a "lost credentials cannot be recovered" property
**Source:** [Capability §Business Rules](_index.md#business-rules--constraints)
**Requirement:** Whatever identity service the platform offers to tenants must support a mode in which lost end-user credentials cannot be recovered (Signal-style), because at least one tenant capability requires this property. An identity option that cannot honor this property is not eligible to be the platform-provided identity service.
**Why this is a requirement, not a decision:** The capability business rules explicitly disqualify any identity option that cannot honor this property.

### TR-12: Tenants must be packaged and declared in a single platform-defined contract, accepted before hosting begins
**Source:** [Capability §Business Rules](_index.md#business-rules--constraints) · [UX: host-a-capability §Journey](user-experiences/host-a-capability.md) · [UX: migrate-existing-data §Constraints Inherited from the Capability](user-experiences/migrate-existing-data.md#constraints-inherited-from-the-capability)
**Requirement:** Every tenant component (including one-shot migration jobs) must be packaged in the form the platform accepts and must declare its resource needs (compute, storage, network reachability) and identity choice up front. The platform must reject components that do not conform to the packaging contract.
**Why this is a requirement, not a decision:** The capability rule "Tenants must accept the platform's contract" makes this a precondition of being hosted at all; multiple UXs lean on the same contract.

### TR-13: Platform changes that introduce, modify, or retire tenant-affecting state must be tracked and immutable, with drift detectable before any rebuild
**Source:** [UX: stand-up-the-platform §Constraints Inherited from the Capability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [UX: stand-up-the-platform §Entry Point](user-experiences/stand-up-the-platform.md#entry-point)
**Requirement:** All platform-altering operations must produce tracked, immutable changes against the definitions repo, such that a preflight drift check can verify the live platform still matches the definitions. Ad-hoc, untracked modification must not be possible.
**Why this is a requirement, not a decision:** The stand-up UX refuses to start a rebuild until a drift check passes; that check is only meaningful if every other UX upholds tracked-and-immutable changes as an invariant.

## Open Questions

Things the user volunteered as solutions during extraction (parked for Stage 2), or constraints the capability/UX docs don't yet make explicit.

- Whether tenant data export (TR-05) should be on-demand or continuously-available — captured during extraction from the move-off UX.
- Whether contract versioning (TR-02) requires semver or a different versioning scheme — surfaced but deferred to Stage 2.
- Whether the canary tenant for readiness (TR-10) should be a single fixed canary or rotated per-rebuild — UX defines a purpose-built canary but does not constrain its identity over time.
- Whether the 30-day post-eviction retention window (referenced in the move-off UX) should be encoded as a TR with a numeric bound, or left to the eviction policy. Surface for human review.
- Re-sourcing TR-06: confirm which section of `migrate-existing-data.md` should be the canonical link now that the original anchor was renamed.
