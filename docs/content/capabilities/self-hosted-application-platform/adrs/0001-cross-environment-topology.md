---
title: "[0001] Cross-Environment Topology"
description: >
    The platform inherits the existing three-tier shape — Internet-facing edge → public-cloud anchor → secure tunnel → private home-lab — with tenant workloads hostable in either environment, all end-user traffic entering through the edge, and the tunnel reserved for platform operations/maintenance.
type: docs
weight: 1
category: "strategic"
status: "accepted"
date: 2026-07-17
deciders: ["Zaba505"]
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions for this capability (auth strategy, data ownership boundaries)
- user-journey: Solutions for specific user-experience problems within this capability
- api-design: API endpoint design decisions for this capability's services

Numbering is local to this capability — start at 0001 and increment.
Status lifecycle: proposed → accepted → (later) superseded
The plan-tech-design skill refuses to compose tech-design.md until every ADR is accepted (or superseded with the superseder accepted).
-->

**Parent capability:** [Self-Hosted Application Platform]({{< ref "../_index.md" >}})
**Addresses requirements:** TR-03, TR-17

## Context and Problem Statement

[TR-03]({{< ref "../tech-requirements.md#tr-03" >}}) forces the rebuild's first phase to establish foundations across *both* a public-cloud environment and a home-lab environment plus the connectivity between them; single-environment standup is explicitly not a supported outcome. [TR-17]({{< ref "../tech-requirements.md#tr-17" >}}) requires each tenant to receive compute, persistent storage, **internal and external network reachability**, identity, backup/DR, and observability — implemented as shared platform offerings. Neither TR names a topology.

The platform therefore needs a decided **foundational environment shape**: what sits on the public-cloud side, what sits on the home-lab side, and how the two connect. This decision is capability-scoped by confirmed framing — the platform *owns* its cross-environment topology, and every capability hosted on the platform inherits whatever shape this ADR exposes rather than deciding its own. The topology exposed to end users may differ from the platform-internal topology decided here.

The open question is whether the platform inherits the shape the repository already realizes — an Internet-facing edge (mutual-auth + traffic-control duties) in front of a public-cloud anchor, with a secure tunnel back to a private home lab — or selects a different shape for its cross-environment foundations.

## Decision Drivers

* **TR-03** — both environments *and* the link between them are phase-1 foundations, not afterthoughts; a single-environment shape is disqualified outright.
* **TR-17** — the shape must provide a clean **external** reachability tier for tenant applications and **internal** cross-environment reachability for platform operation, and must not concentrate backup/DR so narrowly that a single environment's loss is unrecoverable.
* **TR-01 / TR-02** ([def]({{< ref "../tech-requirements.md#tr-01" >}}), [rebuild]({{< ref "../tech-requirements.md#tr-02" >}})) — the whole topology must be expressible as version-controlled definitions and rebuildable end-to-end within 60 minutes; reusing shapes already realized as reproducible definitions is favored over shapes that must be built from scratch.
* **TR-04** ([teardown]({{< ref "../tech-requirements.md#tr-04" >}})) — each environment and the connectivity between them must be independently teardown-able at a phase checkpoint.
* **TR-18** ([admissibility]({{< ref "../tech-requirements.md#tr-18" >}})) — the edge and tunnel components must allow configuration control, data export, and credential revocation/rotation without vendor cooperation; more vendor surface is more admissibility risk.
* **CLAUDE.md house pattern** — `Internet → Cloudflare (mTLS + DDoS) → Home Lab ↔ GCP (WireGuard)` is the repository's documented inter-environment topology, already realized in `cloud/` (`mtls/cloudflare-gcp`, `vpc-network` `allow-wireguard`, `network-load-balancer` UDP gateway). Departing from it requires explicit justification.
* **Capability tiebreaker** — *reproducibility beats vendor independence beats minimizing operator effort.*
* **Operator-effort asymmetry between environments** *(not TR-derived — recorded as an honest driver)*. Managed public-cloud compute and datastore primitives deliver much of the TR-17 inventory without operator-built machinery, where the home-lab equivalents must be built and operated. Their current low-to-zero cost at this platform's scale sharpens the asymmetry. No TR ranks hosting cost, so this driver may motivate the option set but cannot by itself justify the outcome — it is subordinate to the TR-anchored drivers above and to the tiebreaker.

## Considered Options

### Option A — Inherit the three-tier shape (edge → public-cloud anchor → secure tunnel → private home lab)

Keep the shape the repository already realizes: an Internet-facing edge tier carrying mutual-auth and traffic-control (DDoS) duties, a public-cloud environment as the anchor, and a secure tunnel back to the private home lab. Tenant workloads may be placed in either environment (see the placement sub-decision in the outcome below).

* Satisfies **TR-03**: two environments plus their connectivity, all already foundation-phase concerns.
* Strongest on **TR-01/TR-02**: the definitions already exist and are already reproducible, so it is the shortest path to a ≤60-minute rebuild.
* Provides the dedicated external-reachability + scrubbing tier that **TR-17** external reachability wants, while the tunnel carries the internal cross-environment reachability.
* Ranks highest on the **reproducibility** tiebreaker (reuse of proven definitions).
* Cost: carries the most vendor surface (a dedicated edge vendor), the weakest position on **TR-18** and the vendor-independence tiebreaker — the edge vendor must pass the TR-18 admissibility test (config control, export, credential rotation without vendor cooperation).

### Option B — Two-environment, edge-less (public-cloud ingress fronts directly)

Keep both environments and the tunnel, but drop the dedicated Internet-facing edge tier; the public-cloud ingress/load balancer terminates external traffic and mutual-auth directly.

* Still satisfies **TR-03**.
* Better on **TR-18** and **TR-02** (one fewer vendor, fewer moving parts to rebuild).
* Weakens **TR-17** external reachability: loses the edge scrubbing/DDoS tier and pushes all external reachability onto the cloud ingress. Diverges from the CLAUDE.md house pattern and would need that divergence justified.

### Option C — Home-lab-primary, cloud-as-thin-edge (inverted anchor)

Invert the anchor: all stateful tenant offerings live home-side; the public cloud shrinks to a minimal always-on relay for external reachability and tunnel termination.

* Satisfies **TR-03**; ranks highest on the **vendor-independence** tiebreaker (the cloud becomes a replaceable relay).
* Concentrates blast radius and **TR-17** backup/DR on the home lab; external reachability degrades during a home-lab outage; greater distance from the current definitions hurts **TR-02** reproducibility speed. Loses on the reproducibility tiebreaker.

### Option D — Single-environment (all-cloud or all-home-lab)

* **Rejected by TR-03**, which explicitly makes single-environment standup an unsupported rebuild outcome. Recorded here so the reason the simplest shape is off the table is auditable.

## Decision Outcome

Chosen option: **Option A — inherit the three-tier shape**, because it satisfies TR-03 directly, is the fastest and most faithful path to the TR-01/TR-02 reproducibility target (its definitions already exist and are proven in `cloud/`), and wins the capability's stated reproducibility-first tiebreaker. The TR-18 vendor-surface cost is accepted as a bounded, downstream admissibility check on the edge and tunnel vendors rather than a reason to rebuild the shape from scratch.

**The two cross-environment paths are strictly separated planes, and this separation is part of the decision:**

* **Application data plane (end-user traffic).** Deployed tenant applications are reached by end users **only** through the Internet-facing edge, regardless of which environment hosts them: `end user → edge (mTLS + DDoS) → tenant application`. This is the **external** reachability of TR-17. Tenant application traffic never traverses the operations tunnel.
* **Operations / maintenance plane.** The public-cloud ↔ home-lab tunnel (today WireGuard) exists **solely** for platform operation and maintenance — the operator's control of the home-lab environment from the public-cloud side. It is the **internal** reachability of TR-17 and carries no tenant application traffic.

**Workload placement sub-decision: both environments are valid tenant-hosting targets.**

The initial reading of this topology treated the home lab as the sole host for tenant workloads, with the public-cloud anchor limited to edge-facing reachability and the operations plane. That reading is **widened here**: the public-cloud anchor is also a first-class tenant-hosting environment, because the cloud side already offers managed compute and datastore primitives that satisfy the TR-17 inventory (compute, persistent storage, backup/DR, observability) with materially less operator-built machinery than the home-lab equivalents — which serves the tiebreaker's third term (minimizing operator effort) without spending anything on the first two.

The consistency rule is what keeps this from fragmenting the topology: **placement changes where a workload runs, never how it is reached.** A cloud-hosted tenant application is *not* exposed directly via the cloud provider's public ingress; it sits behind the same Internet-facing edge as a home-lab-hosted one, so mutual-auth and traffic-control duties stay in exactly one tier and the TR-17 external-reachability story is identical in both environments.

Which environment a given tenant lands in is a **placement policy**, deliberately not decided here — see Open Questions.

```mermaid
flowchart LR
    user([End user])
    subgraph edge[Internet-facing edge]
        cf[mTLS + DDoS / traffic control]
    end
    subgraph cloud[Public-cloud anchor]
        capp[Tenant applications]
        ops[Operator / platform ops]
    end
    subgraph home[Private home lab]
        happ[Tenant applications]
    end

    user -- application data plane --> cf
    cf -- application data plane --> capp
    cf -- application data plane --> happ
    ops -- operations plane only<br/>secure tunnel --- happ
```

Both hosting targets sit behind the same edge; no tenant application is reachable by bypassing it.

### Consequences

* Good, because the foundations phase reuses already-reproducible `cloud/` definitions, keeping the TR-02 ≤60-minute rebuild target reachable and honoring the reproducibility tiebreaker.
* Good, because separating the application data plane (edge) from the operations plane (tunnel) gives TR-17 a clean external-reachability tier and a distinct internal-reachability tier, and prevents tenant traffic from ever depending on the operations tunnel.
* Good, because routing both environments' tenant traffic through the one edge keeps mutual-auth and traffic-control duties in a single tier — a tenant's reachability story does not change when its placement changes, and placement stays a migratable property rather than a baked-in commitment.
* Bad, because the dedicated edge vendor is the largest vendor-surface commitment, making it the weakest point against TR-18 and the vendor-independence tiebreaker; it must clear the TR-18 admissibility test.
* Bad, because two valid hosting environments means the TR-17 inventory (compute, persistent storage, backup/DR, observability) must be realized **in both** — the requirement is that offerings are shared across tenants, not that they exist once, but a per-environment implementation of each offering is real duplicated surface, and every offering added later pays this cost twice. Divergence between the two implementations is a live risk that the tech design must contain.
* Bad, because hosting tenant workloads on managed cloud primitives deepens exposure to a single cloud provider, pressing on TR-18 (data export, credential rotation without vendor cooperation) more than a home-lab-only placement would. This is accepted because the capability tiebreaker ranks reproducibility above vendor independence, but it is a real concession and each managed primitive admitted as a tenant-hosting offering must pass TR-18 on its own.
* Requires: the home-lab-side foundations must be expressed as version-controlled definitions (TR-01) — there is no `cloud/` analog for the home lab today; a downstream component design must define that surface and its TR-04 teardown.
* Requires: an edge→cloud-anchor origin path for cloud-hosted tenant applications that enforces the same mutual-auth posture as the edge→home-lab path, so "reachable only through the edge" is an enforced property rather than a convention.

### Realization

* `cloud/mtls/cloudflare-gcp/` — the Internet-facing edge trust (mTLS origin certs); the application data plane's entry point.
* `cloud/https-load-balancer/`, `cloud/ip/`, `cloud/dns/` — external reachability plumbing behind the edge.
* `cloud/vpc-network/` (`allow-wireguard` firewall tag) and `cloud/network-load-balancer/` (UDP gateway) — the operations-plane tunnel endpoints on the public-cloud side.
* `cloud/rest-api/`, `cloud/https-load-balancer/`, `cloud/internal-application-load-balancer/` (Cloud Run backends and their network endpoint groups) and `cloud/firestore/` — the public-cloud-side tenant-hosting offerings: managed compute and persistent storage for cloud-placed tenants, fronted by the edge rather than exposed directly.
* Home-lab-side platform offerings (tenant compute/persistent storage per TR-17) — **not yet** a `cloud/` module; a downstream component design owns how the home-lab foundations are expressed as definitions and torn down, and must reconcile its offering surface with the public-cloud-side equivalents above.
* `tech-design.md` (composed later by `plan-tech-design`) will fold this shape into the final-state narrative alongside the other accepted ADRs.

## Open Questions

* **Edge/tunnel vendor admissibility (TR-18).** The specific vendors realizing each layer are out of scope for this topology decision, but each must pass the TR-18 admissibility test (config control, data export, credential revocation/rotation without vendor cooperation). This is a downstream check, not a reopening of the shape.
* **Home-lab definitions surface (TR-01/TR-04).** How the home-lab-side foundations are expressed as version-controlled definitions and given a deterministic per-phase teardown has no `cloud/` precedent yet; it belongs to a downstream component design.
* **Ops-plane sequencing within phase 1 (TR-03/TR-04).** Whether the operations tunnel is stood up and torn down as its own checkpoint relative to the edge and the cloud anchor, or as a single foundations unit, is left to the tech design.
* **Tenant placement policy (TR-17).** This ADR establishes that both environments are valid hosting targets; it does not decide how a given tenant's environment is chosen, who chooses it (operator, capability owner, or a declared workload property), or whether placement is re-negotiable after onboarding. This is a distinct decision and likely warrants its own ADR — it bears on the onboarding and modify engagement flows, not on the topology shape.
* **Cross-environment offering parity (TR-17).** Whether every offering in the TR-17 inventory must exist in both environments before a tenant may be placed there, or whether environments may expose deliberately unequal offering sets (with placement constrained by what a tenant needs), is unresolved. The answer determines how much duplicated surface the platform actually carries.
* **Migration between environments (TR-17).** Whether a tenant can be moved between environments after onboarding, and what that costs in data migration and downtime, is not decided here. Keeping reachability identical across environments is what makes this *possible*; it does not make it *supported*.
