---
title: "[0008] Observability Stack — LGTM + OpenTelemetry Collector"
description: >
    Self-hosted Grafana LGTM stack (Loki for logs, Grafana for visualization, Tempo for traces, Mimir for metrics) with an OpenTelemetry Collector as the single ingestion surface for tenants. Full OTel/OTLP compliance; native multi-tenancy via per-tenant tenant-IDs at the LGTM components.
type: docs
weight: 8
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform]({{< relref "../_index.md" >}})
**Addresses requirements:** [TR-06]({{< relref "../tech-requirements.md#tr-06" >}}), [TR-07]({{< relref "../tech-requirements.md#tr-07" >}}), [TR-08]({{< relref "../tech-requirements.md#tr-08" >}}), [TR-09]({{< relref "../tech-requirements.md#tr-09" >}}), [TR-22]({{< relref "../tech-requirements.md#tr-22" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}), [TR-33]({{< relref "../tech-requirements.md#tr-33" >}})

## Context and Problem Statement {#context}

Five tightly-coupled responsibilities ride on this ADR:

- **Standard health bundle** ([TR-06]({{< relref "../tech-requirements.md#tr-06" >}})) — availability, latency, error rate, resource saturation, restart/deployment events — uniformly across tenants.
- **Tenant-scoped access for capability owners; cross-tenant for the operator** ([TR-07]({{< relref "../tech-requirements.md#tr-07" >}})) — one offering, scope enforcement.
- **Self-serve threshold tuning** ([TR-08]({{< relref "../tech-requirements.md#tr-08" >}})) — the only capability-owner self-service surface in the entire platform.
- **Email push alerting with degraded-delivery indication** ([TR-09]({{< relref "../tech-requirements.md#tr-09" >}})) — alerts go by email; the tenant view must surface when delivery is degraded so silence is not mistaken for health.
- **Tenant isolation** ([TR-32]({{< relref "../tech-requirements.md#tr-32" >}})) — no tenant or capability owner can read another tenant's signals.

Because the platform's compute is Kubernetes ([ADR-0002]({{< relref "0002-compute-substrate.md" >}})) and the packaging form is platform-defined ([ADR-0003]({{< relref "0003-tenant-packaging-form.md" >}})), the ingestion contract from tenant pods is a platform decision: tenants emit signals through whatever protocol the platform standardizes on, and the platform fans them into the storage layer.

The operator's stated direction is to commit to **full OpenTelemetry/OTLP compliance** for the ingestion surface. That choice rules out Prometheus-scrape-only stacks at the ingestion layer (Prometheus has OTLP receive support but the model is scrape-first) and points at the Grafana **LGTM** stack (Loki / Grafana / Tempo / Mimir), all of which natively accept OTLP and share the same multi-tenancy model — per-tenant *tenant IDs* threaded through every component.

## Decision Drivers {#decision-drivers}

- **OTel/OTLP-first ingestion.** A single ingestion protocol that covers metrics, logs, and traces, kept open and vendor-neutral, simplifies the tenant contract ([TR-06]({{< relref "../tech-requirements.md#tr-06" >}})) and keeps the platform's signal-data egress path swappable later.
- **Native multi-tenancy.** The storage components must enforce tenant-scope at the API layer, not only at the visualization layer. ([TR-07]({{< relref "../tech-requirements.md#tr-07" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}).)
- **One offering, two roles.** The operator's cross-tenant view and the capability owner's tenant-scoped view live behind the same offering ([UX: tenant-facing-observability]({{< relref "../user-experiences/tenant-facing-observability.md" >}})). Stacks that need a separate operator-only product fail this.
- **Self-serve thresholds in-product.** [TR-08]({{< relref "../tech-requirements.md#tr-08" >}}) is realized inside the visualization tool. Stacks where alerting is a separate front-end fragment the surface and break the "one offering" promise.
- **Degraded-delivery surface.** The tenant view must show when email delivery is broken ([TR-09]({{< relref "../tech-requirements.md#tr-09" >}})) — the alerting component must expose delivery health as queryable signals.
- **Capability tiebreaker — vendor independence > minimizing operator effort.** Same logic as [ADR-0005]({{< relref "0005-identity-offering.md" >}}) rejecting vendor identity: observability data is sensitive (traffic patterns, error payloads, trace contents). Self-hosted is preferred.
- **[TR-33]({{< relref "../tech-requirements.md#tr-33" >}}).** Each component is paid weekly. Single-replica deployments with [ADR-0007]({{< relref "0007-backup-and-disaster-recovery.md" >}}) covering catastrophic loss is the right posture at this scale; HA via clustering of each component is gold-plating.

## Considered Options {#considered-options}

### Option A — Prometheus + Grafana + Alertmanager + Loki (self-hosted)

Industry-standard, deep K8s integration. Prometheus scrapes; multi-tenant via Grafana Organizations and per-tenant scrape configs; logs in Loki; alerting via Alertmanager.

- **OTel/OTLP:** partial. Prometheus has experimental OTLP receive; the natural model is scrape-first. Tenants emitting OTLP would land on a translator pattern.
- **Multi-tenancy:** at visualization layer (Grafana orgs) and at scrape-config layer; *not* at the metrics-store API layer (single Prometheus is one tenant from its own perspective).
- **Traces:** not included; would require adding Tempo or Jaeger separately.
- **Surface:** moderate.

### Option B — Grafana Cloud (vendor SaaS)

- **OTel/OTLP:** native.
- **Multi-tenancy:** native.
- **Vendor-independence tiebreaker:** **fails**, on the same logic that rejected vendor identity in [ADR-0005]({{< relref "0005-identity-offering.md" >}}). Observability data is sensitive; egress to a vendor moves the trust boundary outside the operator's control.
- **Cost:** real, scales with signal volume.

### Option C — VictoriaMetrics + VictoriaLogs + Grafana + VMAlert

A lighter Prometheus-compatible stack.

- **OTel/OTLP:** supported.
- **Multi-tenancy:** native (per-tenant accountID/projectID).
- **Traces:** not first-class; would require pairing with another component.
- **Community:** smaller than Grafana's; "boring widely-known answer" benefit is weaker.

### Option D — LGTM (Loki + Grafana + Tempo + Mimir) with OTel Collector

- **L = Loki** for logs.
- **G = Grafana** for visualization, dashboards, alerting UI.
- **T = Tempo** for traces.
- **M = Mimir** for metrics (Prometheus-compatible TSDB with native multi-tenancy).
- **OTel Collector** in front, accepting OTLP from tenant pods and fanning out to Loki / Tempo / Mimir.

- **OTel/OTLP:** native everywhere. The Collector is the single ingestion contract; LGTM components consume the Collector's output natively.
- **Multi-tenancy:** native at *every storage component* — Mimir, Loki, and Tempo all accept and enforce a `X-Scope-OrgID` (or equivalent) tenant-ID header. The OTel Collector tags incoming OTLP from a tenant namespace with that namespace's tenant-ID and forwards. Grafana data sources per Organization scope to the matching tenant-ID. Cross-tenant isolation is enforced at the API of every storage component, not only at the UI.
- **Traces included.** [TR-06]({{< relref "../tech-requirements.md#tr-06" >}})'s standard bundle is metrics-and-events focused; traces are not strictly required by it but are extremely useful for "is this my tenant or the platform" root-causing ([UX: tenant-facing-observability]({{< relref "../user-experiences/tenant-facing-observability.md" >}})). They come for free with this option.
- **Alerting:** Grafana Alerting (or Mimir's bundled Alertmanager). Email delivery and queue health expose Prometheus-format metrics into Mimir, which the tenant view then queries to surface a "degraded delivery" panel ([TR-09]({{< relref "../tech-requirements.md#tr-09" >}})).
- **Surface:** four LGTM components + the Collector. Real, but each is operationally well-understood and the ecosystem documentation is dense.

## Decision Outcome {#decision-outcome}

Chosen option: **Option D — LGTM stack (Loki + Grafana + Tempo + Mimir) with OpenTelemetry Collector as the single ingestion surface**. Self-hosted in the home-lab cluster, single-replica each, on the block storage primitive from [ADR-0004]({{< relref "0004-persistent-storage-offering.md" >}}) and within the backup scope of [ADR-0007]({{< relref "0007-backup-and-disaster-recovery.md" >}}).

This option is chosen because:

- It is the only option where **OTel/OTLP is native at every layer** — ingestion, storage, querying. The tenant contract for emitting signals is exactly OTLP, with no translator pattern bolted in front.
- **Native multi-tenancy at each storage component** ([TR-07]({{< relref "../tech-requirements.md#tr-07" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}})). Mimir, Loki, and Tempo each enforce tenant-scope at their HTTP API via tenant-ID headers; visualization-only scope (Option A) is weaker because a misconfigured Grafana org can leak across tenants whereas a missing tenant-ID header in LGTM yields no data at all.
- **Traces (Tempo) come along** without a separate stack, directly serving the tenant-vs-platform root-causing path ([UX: tenant-facing-observability]({{< relref "../user-experiences/tenant-facing-observability.md" >}})).
- **One Grafana, two roles**: capability owners log into a tenant-scoped Grafana Organization; the operator has a cross-tenant Organization. Same product, both [TR-07]({{< relref "../tech-requirements.md#tr-07" >}}) views.
- **Alerting + degraded-delivery indicator are first-class.** Grafana Alerting (or Mimir Alertmanager) exposes delivery metrics; the tenant view binds a panel to those metrics so [TR-09]({{< relref "../tech-requirements.md#tr-09" >}})'s "do not treat email silence as health" promise is realized in the same product.
- **Vendor independence holds** ([ADR-0005]({{< relref "0005-identity-offering.md" >}}) precedent). All components are open source and self-hosted; the OTel ingestion protocol is the open standard, so even if individual LGTM components were swapped later, the tenant contract does not change.

The stack runs **single-replica each** in the home-lab cluster, with PVCs on the block storage primitive ([ADR-0004]({{< relref "0004-persistent-storage-offering.md" >}})) and inclusion in the backup scope of [ADR-0007]({{< relref "0007-backup-and-disaster-recovery.md" >}}). HA via component clustering is rejected at this scale; loss of recent observability history is acceptable, and full restore from backup is the recovery path. Catastrophic observability loss does not destroy tenant data — it loses recent signals.

Authentication into Grafana federates to **Authentik via OIDC** ([ADR-0005]({{< relref "0005-identity-offering.md" >}})). Capability-owner role and tenant-Organization membership are derived from Authentik claims; the operator role is a separate Authentik group with cross-tenant Organization membership.

### Consequences {#consequences}

- **Good, because** OTLP is the single ingestion contract for tenants. The tenant manifest schema ([ADR-0003]({{< relref "0003-tenant-packaging-form.md" >}})) does not need separate fields per signal type — tenants emit OTLP, the platform routes.
- **Good, because** isolation is enforced at the storage API of every component, not only at the visualization layer. A misconfigured Grafana Organization is a recoverable incident; a missing tenant-ID header at the storage API yields *no signals at all*, which is fail-closed.
- **Good, because** the cross-tenant operator view and the tenant-scoped capability-owner view live in the same Grafana, which honors [UX: tenant-facing-observability]({{< relref "../user-experiences/tenant-facing-observability.md" >}})'s "no separate URL per tenant" rule literally.
- **Good, because** [TR-08]({{< relref "../tech-requirements.md#tr-08" >}}) self-serve thresholds and [TR-09]({{< relref "../tech-requirements.md#tr-09" >}}) degraded-delivery indication are realized inside one product surface, not split across multiple UIs.
- **Good, because** traces (Tempo) directly support the [UX: tenant-facing-observability]({{< relref "../user-experiences/tenant-facing-observability.md" >}}) "is this me or the platform?" question — the capability owner can follow a request from ingress to their pod to confirm where the failure lives.
- **Bad, because** the platform now operates four LGTM components + the OTel Collector. [TR-33]({{< relref "../tech-requirements.md#tr-33" >}}) is pressured; mitigation is single-replica posture, declarative configuration, and Grafana's well-trodden upgrade discipline. Watch closely.
- **Bad, because** the OTel Collector is now load-bearing — if it falls behind, every tenant's signals fall behind. It must be in the canary's exercised path so a broken Collector fails the rebuild rather than ships unnoticed.
- **Bad, because** OTLP-first ingestion subtly shifts work onto tenants — they must instrument their pods to emit OTLP. The platform mitigates this by having the Collector also scrape K8s baseline metrics (kubelet, kube-state-metrics, cAdvisor) so the standard health bundle's *availability / restart / saturation / resource* signals come from substrate observation regardless of what the tenant instruments. Tenant-emitted signals are bonus depth, not the floor.
- **Bad, because** Mimir, Loki, and Tempo each have their own retention configuration to tune. Defaulting all three to a uniform horizon (proposal: 14 days local, in line with [ADR-0007]({{< relref "0007-backup-and-disaster-recovery.md" >}})'s 30-day shape but shorter to keep PVC sizes bounded) reduces the chance of inconsistent windows.
- **Requires:**
  - **[ADR-0003]({{< relref "0003-tenant-packaging-form.md" >}}) schema** gains a tenant-ID field (the platform-controlled identifier that becomes the LGTM tenant-ID header) and an optional OTLP-emit flag for the tenant's pod. The schema does not let tenants set their own tenant-ID; the platform owns it.
  - **[ADR-0005]({{< relref "0005-identity-offering.md" >}}) Authentik realm** for capability-owner login to Grafana, with group claims for tenant-Organization membership and an operator group with cross-tenant access.
  - **[ADR-0006]({{< relref "0006-network-reachability.md" >}})** ensures Grafana's external hostname is exposed only to capability-owner / operator authentication paths, never to end users ([TR-28]({{< relref "../tech-requirements.md#tr-28" >}})).
  - **Future ADR (secrets)** stores Grafana SMTP credentials (for outbound email) and the OIDC client secret for Authentik integration.
  - **Future ADR (export tooling)** is *not* responsible for exporting observability data on eviction. Observability data is platform-side and not part of the tenant's exportable archive; this is consistent with the eviction UX, which exports tenant *data* (the bytes the tenant capability stored), not platform-collected signals. Worth confirming if pushed.
  - **Future ADR (canary tenant)** emits OTLP and triggers a deliberate alert during readiness verification, so both ingestion and alert delivery are exercised on every rebuild.

### Realization {#realization}

How this decision shows up in the repo:

- **An `observability` namespace** in the home-lab cluster runs:
  - The **OpenTelemetry Collector** (DaemonSet for substrate metrics + Deployment for tenant OTLP ingestion).
  - **Mimir** (single-replica, monolithic-mode), with PVCs on the block storage class.
  - **Loki** (single-replica, monolithic-mode), with PVCs.
  - **Tempo** (single-replica, monolithic-mode), with PVCs.
  - **Grafana** with OIDC integration to Authentik.
  - **Alertmanager** (or Grafana Alerting) with SMTP configured for outbound email.
- **Per-tenant configuration** is emitted by the [ADR-0003]({{< relref "0003-tenant-packaging-form.md" >}}) translator:
  - A tenant-ID derived from the tenant's namespace name.
  - The Collector is configured to tag signals from that namespace with the tenant-ID before forwarding to LGTM components.
  - A Grafana Organization with the tenant-ID's data sources scoped to that tenant-ID.
  - A capability-owner Authentik group that maps to the Grafana Organization.
- **The OTel Collector** also scrapes substrate sources (kubelet, kube-state-metrics, cAdvisor) and produces the standard-health-bundle signals for every tenant regardless of in-pod instrumentation.
- **Email delivery health metrics** from Alertmanager / Grafana Alerting feed back into Mimir, surfaced in the tenant Organization's default dashboard as the [TR-09]({{< relref "../tech-requirements.md#tr-09" >}}) degraded-delivery indicator.
- **Alert thresholds** ([TR-08]({{< relref "../tech-requirements.md#tr-08" >}})) are managed through Grafana Alerting within each tenant Organization. Capability owners self-serve.
- **The canary tenant** emits OTLP from its pods, exercises the Collector → LGTM path, and triggers a known alert that is verified to arrive by email; absence of the alert fails the rebuild.

## Open Questions {#open-questions}

- **Local retention horizon per LGTM component.** Proposal: 14 days for Mimir / Loki / Tempo locally; coordinated with [ADR-0007]({{< relref "0007-backup-and-disaster-recovery.md" >}})'s 30-day off-site backup shape. Confirm at deployment.
- **Single Collector vs. per-namespace.** Single cluster-wide Collector is simpler; per-namespace gives tighter isolation in case of Collector compromise. Default single; revisit if a tenant requires it.
- **Whether LGTM signals are part of the export archive on eviction.** Default no — observability is platform-collected and not tenant-owned data. Worth flagging in case a future tenant disagrees.
- **Grafana edition** (OSS vs. Enterprise). OSS is sufficient; Enterprise features (advanced RBAC, SSO advanced) are not required given Authentik integration via OIDC.
- **HA posture for any individual component** if the operator finds, in practice, that loss of recent observability history during a rebuild is operationally painful. Single-replica is the day-one posture; HA is a pull from operational pain, not a pre-decision.
