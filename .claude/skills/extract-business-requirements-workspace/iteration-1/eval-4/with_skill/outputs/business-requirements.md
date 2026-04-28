---
title: "Business Requirements"
description: Business requirements for the self-hosted-application-platform capability.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a BR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). Technical requirements cite BR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-technical-requirements` skill will refuse to extract TRs until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [self-hosted-application-platform]({{< ref "_index.md" >}})

## How to read this

Each requirement is **forced** by the capability or a user experience — it states, in business or user-outcome terms, what the system must guarantee. Decisions about the *technical translation* (cadences, durability levels, protocols) belong in `tech-requirements.md`. Decisions about *how* (which database, which library, which provider) belong in `adrs/`. If something in this list reads like a technical constraint or a chosen solution rather than a business demand, flag it for review.

## Requirements

### BR-01: No tenant can ever observe another tenant's state
**Source:** [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})
**Requirement:** A tenant must never be able to read, infer, or otherwise observe another tenant's data, secrets, traffic, or telemetry — under normal operation or any degraded mode the platform supports.
**Why this is a requirement, not a TR or decision:** The capability lists tenant isolation as an invariant; it is a business demand on the platform, independent of any technical mechanism that enforces it.

### BR-02: Existing tenants must keep working when the platform contract evolves
**Source:** [UX: platform-contract-change-rollout]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})
**Requirement:** When the platform publishes a new contract version, every tenant currently operating against the previous contract must continue to operate, without operator intervention, until they choose to migrate.
**Why this is a requirement, not a TR or decision:** The UX flow demands graceful contract evolution as the user-perceived outcome; the technical means (versioning scheme, compatibility windows) is a TR/ADR question.

### BR-03: A tenant must be able to query observability data for their own workloads
**Source:** [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})
**Requirement:** Tenants must be able to inspect metrics, logs, and traces for the workloads they own, and only those — without operator help and without exposure to other tenants' observability data.
**Why this is a requirement, not a TR or decision:** Both the capability isolation invariant and the UX journey demand this as a tenant-facing outcome.

### BR-04: Operator-initiated tenant updates must be invisible to the tenant's end users
**Source:** [UX: operator-initiated-tenant-update §Success]({{< ref "user-experiences/operator-initiated-tenant-update.md#success" >}})
**Requirement:** When an operator updates a tenant's configuration, version, or capability, online traffic served by that tenant must not experience end-user-visible downtime.
**Why this is a requirement, not a TR or decision:** UX success criteria define the absence of user-visible downtime as the outcome; the technical translation (rolling deploys, draining, etc.) is a TR/ADR question.

### BR-05: An evicted tenant must be able to leave with all of their data
> ⚠️ source no longer resolves — human review

**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#a-section-that-no-longer-exists" >}})
**Requirement:** A tenant who has been evicted from the platform must still be able to retrieve all of their data in a portable form within a defined export window — without operator assistance.
**Why this is a requirement, not a TR or decision:** The UX explicitly includes the eviction edge case; portability is a business commitment to the tenant, independent of how the export is implemented.

## Open Questions

Things the user volunteered as TRs or decisions during extraction (parked for the next stage), or constraints the capability/UX docs don't yet make explicit.

- Whether tenant data export (BR-05) should be on-demand or continuously-available — surfaced during extraction from the move-off UX; defer to TR stage.
- Whether contract versioning (BR-02) is bounded by a maximum compatibility window — surfaced but not yet decided in the UX.
- BR-05's source anchor (`#a-section-that-no-longer-exists`) no longer resolves in `user-experiences/move-off-the-platform-after-eviction.md`. The likely current section is `## Edge Cases & Failure Modes`, but the source needs a human re-source decision (and an explicit `{#anchor}` on the target heading) before the link is restored.
