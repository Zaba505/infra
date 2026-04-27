---
title: "[0004] Tenant-Scoped Observability Stack"
description: >
    Choose the observability stack and the mechanism by which tenants can query their own metrics, logs, and traces while being prevented from seeing any other tenant's data.
type: docs
weight: 4
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

[TR-03] requires that tenant-facing observability data be queryable per-tenant within their data scope only. Tenants must be able to query metrics, logs, and traces for their own workloads; cross-tenant data must be inaccessible.

ADR-0001 already provides per-VM structural isolation at the *agent* layer (an agent in tenant A's VM cannot read tenant B's workloads). This ADR addresses the *backend* layer: where the signals are aggregated, how they are stored, and how query-time scoping enforces TR-03 as defence-in-depth.

## Decision Drivers

* TR-03: tenant-scoped queries; cross-tenant data inaccessible to tenants
* TR-01: defence-in-depth — even if the agent layer is compromised, the backend must not let tenant A read tenant B
* NFR-02: maintenance budget — running three separate single-tenant Prometheus / Loki / Tempo stacks per tenant burns the budget
* NFR-01: rebuild ≤ 1 hr — backends must be definitions-driven
* TR-07: signals traverse Cloudflare → GCP path (no out-of-band telemetry exfil)

## Considered Options

* **Per-tenant single-tenant stacks (one Prometheus + Loki + Tempo per tenant)** — strongest isolation, highest overhead.
* **Shared Grafana Mimir + Loki + Tempo with per-tenant `X-Scope-OrgID` (multitenancy mode)** — single shared backend; tenant scope enforced at the query gateway.
* **Shared backends fronted by a per-tenant query proxy that injects scope and rejects unscoped queries** — middle ground; backends are not multitenant-aware but the proxy enforces it.

## Decision Outcome

Chosen option: **shared Grafana Mimir + Loki + Tempo with per-tenant `X-Scope-OrgID`**, with a per-tenant query gateway that the operator controls and that the tenant's queries flow through.

The Grafana stack's multitenancy mode is purpose-built for this: each signal is written with a tenant ID derived from the source VM's identity (set by the agent in the tenant's VM, signed at the platform's ingest tier so a compromised agent cannot spoof a different tenant ID). Reads go through a per-tenant gateway that injects the `X-Scope-OrgID` header and rejects any query that attempts to override it. The tenant cannot bypass the gateway because their network egress to the observability backend is constrained at the ingress (TR-07's Cloudflare layer).

This composes cleanly with ADR-0001: agent-layer isolation is structural; backend-layer isolation is enforced by the gateway; an attacker would need to break both to see another tenant's data.

### Consequences

* Good, because TR-03 is enforced at two independent layers (agent placement + query-gateway scope)
* Good, because the operator runs one set of backends instead of N — fits NFR-02
* Good, because the same backends serve the operator's cross-tenant view (the operator's gateway has no scope restriction)
* Good, because the Grafana stack is well-documented for this exact use case
* Neutral, because the platform must operate the ingest signing and the gateway — both are small services but they are operator code
* Bad, because a bug in the gateway is a multi-tenant blast radius — mitigated by the agent-layer isolation from ADR-0001 and by gateway integration tests in the canary
* Bad, because Mimir/Loki/Tempo are heavier than per-tenant single-instance Prometheus — but the per-tenant alternative scales linearly and is worse at N>2 tenants

### Confirmation

* The canary tenant's standup test (REQ-18 long-form) attempts to query for the operator's own observability data using the canary's gateway; the query must be rejected with HTTP 403
* The platform-side ingest signing key is rotated on a documented cadence; rotation is captured in the standup runbook
* All observability ingestion traffic flows through Cloudflare → GCP per TR-07; verified by pcap on the platform's egress at standup
* The per-tenant gateway's allow-list of `X-Scope-OrgID` values is generated from the platform's tenant registry; CI fails the build if a tenant exists without a corresponding gateway entry

## Pros and Cons of the Options

### Per-tenant single-tenant stacks

* Good, because isolation is total — no shared backend at all
* Good, because tenant-specific retention, sampling, and alerting are independent
* Bad, because operator runs N stacks; NFR-02 dies at N>2
* Bad, because cross-tenant operator view requires a federated query layer that recreates the multi-tenant problem
* Bad, because rebuild time (NFR-01) scales with tenant count

### Shared Grafana Mimir + Loki + Tempo with `X-Scope-OrgID`

* Good, because purpose-built for this
* Good, because operator runs one set of backends
* Good, because the cross-tenant view comes for free (the operator's gateway has no scope restriction)
* Neutral, because the gateway and signing are operator-owned code
* Bad, because a gateway bug has multi-tenant blast radius — mitigated by ADR-0001's structural agent isolation

### Shared backends fronted by a per-tenant scope-injecting proxy

* Good, because the backends can be any time-series / log / trace store
* Good, because the proxy is small and reviewable
* Bad, because writes also need scoping and the backends do not natively understand tenant IDs — the proxy must enforce both directions, which is more code than the Grafana-stack option provides out of the box
* Bad, because rolling backends without losing scope semantics requires per-backend understanding

## More Information

* The agent in each tenant VM is the platform-controlled image; the tenant does not configure or replace it
* Alert delivery channels are tenant-declared at onboarding (REQ-13 long-form); the alert pipeline reads from the same tenant-scoped backend
* Retention defaults are platform-set; per-tenant overrides require a contract-version bump (ADR-0003)
