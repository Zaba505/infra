---
title: "Business Requirements"
description: Business requirements for the self-hosted-application-platform capability.
type: docs
reviewed_at: 2026-04-20
---

**Parent capability:** [self-hosted-application-platform](_index.md)

## Requirements

### BR-01: No tenant can ever observe another tenant's state
**Source:** [Capability §Business Rules](_index.md#business-rules)
**Requirement:** A tenant must never be able to read, infer, or otherwise observe another tenant's data, secrets, traffic, or telemetry — under normal operation or any degraded mode the platform supports.
**Why this is a requirement, not a TR or decision:** The capability lists tenant isolation as an invariant; it is a business demand on the platform, independent of any technical mechanism that enforces it.

### BR-02: Existing tenants must keep working when the platform contract evolves
**Source:** [UX: platform-contract-change-rollout](user-experiences/platform-contract-change-rollout.md)
**Requirement:** When the platform publishes a new contract version, every tenant currently operating against the previous contract must continue to operate, without operator intervention, until they choose to migrate.
**Why this is a requirement, not a TR or decision:** The UX flow demands graceful contract evolution as the user-perceived outcome; the technical means (versioning scheme, compatibility windows) is a TR/ADR question.

### BR-03: A tenant must be able to query observability data for their own workloads
**Source:** [UX: tenant-facing-observability §Journey](user-experiences/tenant-facing-observability.md)
**Requirement:** Tenants must be able to inspect metrics, logs, and traces for the workloads they own, and only those — without operator help and without exposure to other tenants' observability data.
**Why this is a requirement, not a TR or decision:** Both the capability isolation invariant and the UX journey demand this as a tenant-facing outcome.

### BR-04: Operator-initiated tenant updates must be invisible to the tenant's end users
**Source:** [UX: operator-initiated-tenant-update §Success](user-experiences/operator-initiated-tenant-update.md)
**Requirement:** When an operator updates a tenant's configuration, version, or capability, online traffic served by that tenant must not experience end-user-visible downtime.
**Why this is a requirement, not a TR or decision:** UX success criteria define the absence of user-visible downtime as the outcome; the technical translation (rolling deploys, draining, etc.) is a TR/ADR question.

### BR-05: An evicted tenant must be able to leave with all of their data
**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases](user-experiences/move-off-the-platform-after-eviction.md#a-section-that-no-longer-exists)
**Requirement:** A tenant who has been evicted from the platform must still be able to retrieve all of their data in a portable form within a defined export window — without operator assistance.
**Why this is a requirement, not a TR or decision:** The UX explicitly includes the eviction edge case; portability is a business commitment to the tenant, independent of how the export is implemented.

## Open Questions

- Whether tenant data export (BR-05) should be on-demand or continuously-available — surfaced during extraction from the move-off UX; defer to TR stage.
- Whether contract versioning (BR-02) is bounded by a maximum compatibility window — surfaced but not yet decided in the UX.
