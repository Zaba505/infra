---
title: "[0001] Public/Private Infrastructure Split"
description: >
    Place tenant compute, persistent storage, backup primary, and observability collection in the home-lab; place external ingress and the off-site backup archive in a small cloud edge; connect the two with a self-hosted VPN.
type: docs
weight: 1
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform]({{< relref "../_index.md" >}})
**Addresses requirements:** [TR-27]({{< relref "../tech-requirements.md#tr-27" >}}), [TR-01]({{< relref "../tech-requirements.md#tr-01" >}}), [TR-02]({{< relref "../tech-requirements.md#tr-02" >}}), [TR-03]({{< relref "../tech-requirements.md#tr-03" >}}), [TR-05]({{< relref "../tech-requirements.md#tr-05" >}}), [TR-17]({{< relref "../tech-requirements.md#tr-17" >}}), [TR-22]({{< relref "../tech-requirements.md#tr-22" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}), [TR-33]({{< relref "../tech-requirements.md#tr-33" >}})

## Context and Problem Statement {#context}

[TR-27]({{< relref "../tech-requirements.md#tr-27" >}}) forces the platform to span public-cloud and private/home-lab infrastructure, with the connectivity between them part of the foundation. The capability rule that grounds this — *"self-hosted means the operator controls end-to-end, not that every component runs on hardware the operator owns"* — leaves the *split* unspecified: any of the platform's offerings (TR-01 compute, TR-02 storage, TR-03 reachability, TR-04 identity, TR-05 backup, TR-06 observability) may be placed on either side.

This ADR decides where each major offering category lives, and how the two sides are connected. Subsequent ADRs (compute substrate, storage shape, identity product, etc.) will pick concrete implementations within whichever side this ADR places them on; the placement is what they need to know.

The capability's own tiebreaker is explicit: *reproducibility beats vendor independence beats minimizing operator effort*, and *cost is secondary to convenience and resiliency*. The split must be evaluated against those, not against ease of any particular offering's implementation.

## Decision Drivers {#decision-drivers}

- **TR-27** — must span both sides; connectivity is part of the foundation, not glue added later.
- **TR-17** — full rebuild ≤1 hour. The rebuild has to provision *both* sides, so the more cross-boundary coordination at startup, the harder the KPI.
- **TR-22** — tracked changes and immutability everywhere. Each side needs a definitions-driven mechanism the operator can lock down; the cheaper that mechanism is per side, the better.
- **TR-33** — ≤2 hr/week routine operation. More boundaries to mind = more steady-state work.
- **TR-32** — tenant isolation must hold. Cross-boundary request paths multiply the surface where isolation can be misconfigured.
- **Capability tiebreaker — vendor independence > minimizing operator effort.** Pure cloud is the operator-effort-minimizing option but trades against vendor independence.
- **Capability tiebreaker — cost secondary to convenience/resiliency.** Cost is allowed to grow when it buys meaningful convenience or resiliency; it is not allowed to grow gratuitously.

## Considered Options {#considered-options}

### Option A — Cloud-heavy (everything in cloud, home-lab as secondary/staging)

All offerings — compute, storage, identity, observability, backup — live in cloud. Home-lab is at most a backup site, possibly nothing.

- **Pros:** simplest to operate; no ISP/power dependency for tenants; one substrate to learn; rebuild (TR-17) is bounded by cloud APIs alone.
- **Cons against TRs and capability tiebreakers:** highest steady-state cost with no convenience/resiliency justification; tilts hard against the vendor-independence tiebreaker (every offering rides on one vendor's product roadmap); makes TR-27 spanning ceremonial — the home-lab carries no real load and nothing falls back to it; tenant data accumulates on a vendor's storage substrate, raising the cost of leaving them later.

### Option B — Home-lab-heavy with a small cloud edge

Home-lab hosts the bulk substrate: tenant compute, tenant persistent storage, backup *primary*, and observability *collection*. Cloud hosts only what the home-lab cannot reasonably do on its own: external ingress (so tenants are reachable when end users are on the public internet) and the *off-site backup archive* (so home-lab hardware loss is recoverable). The two sides are connected by a self-hosted VPN; identity and the operator-facing observability surfaces are placed in subsequent ADRs and may land on either side.

- **Pros against TRs and capability tiebreakers:** operator owns the bulk substrate, satisfying the vendor-independence tiebreaker; cloud bill stays small (only ingress + archive storage); TR-27 spanning is *load-bearing* — a tenant request really does cross the boundary on every external call, and DR really does pull from the cloud-side archive; the cloud surface is small enough that TR-22 (tracked changes/immutability) is cheap to enforce there; TR-32 isolation reasoning is contained on the home-lab side because cross-boundary traffic is bounded to ingress + archive.
- **Cons against TRs and capability tiebreakers:** ISP outage takes tenants offline (acceptable per the capability's "no specific availability SLA"); home-lab hardware loss is a real DR event, mitigated by the cloud-side archive plus TR-05; TR-17 1-hour rebuild has to include home-lab provisioning, which is harder than pure cloud — a real cost paid against this option.

### Option C — Symmetric split-by-suitability (each offering placed individually)

No global rule; each offering placed where it is locally optimal. E.g. compute home-lab, storage home-lab, identity cloud, observability cloud, backup cloud, ingress cloud.

- **Pros:** each offering optimal in isolation.
- **Cons against TRs and capability tiebreakers:** every request path crosses the boundary, so TR-27 connectivity becomes load-bearing for *every* call rather than just ingress and archive; TR-32 isolation gets harder to reason about when tenant data and tenant identity live on opposite sides; TR-33 (2hr/week) likely worse — more boundaries to mind during routine operation; TR-17 1-hour rebuild has to coordinate startup ordering across many cross-boundary dependencies; effectively asks the operator to be expert in two substrates simultaneously without the home-lab carrying the weight that would justify the investment.

### Connectivity sub-option — self-hosted VPN vs. vendor private interconnect

Whichever split is chosen, the cloud↔home-lab link must be operator-controllable, reproducible from definitions (TR-17, TR-22), and cheap. Vendor private-interconnect products fail the cost test at personal scale and fail the vendor-independence tiebreaker. A self-hosted VPN (e.g. WireGuard) terminated on both sides satisfies all three. The specific VPN product is left to the side that hosts the terminator and is not material to this ADR's split.

## Decision Outcome {#decision-outcome}

Chosen option: **Option B — Home-lab-heavy with a small cloud edge** (compute / persistent storage / backup primary / observability collection in home-lab; external ingress and off-site backup archive in cloud), connected by a self-hosted VPN.

This option is chosen because:

- It honors the capability's tiebreaker that **vendor independence > minimizing operator effort**: the bulk substrate sits on hardware the operator owns, and the cloud surface is small enough that switching providers later is bounded work rather than an existential exit.
- It makes TR-27 (spanning) *load-bearing* rather than ceremonial: external tenant traffic and the off-site archive really do cross the boundary, so the connectivity is exercised continuously and stays trustworthy.
- TR-22 enforcement is cheap on the cloud side because the cloud side is small; on the home-lab side, the operator already has root and can lock down the change paths directly.
- TR-32 isolation reasoning is concentrated on the home-lab side, where the bulk of tenant state lives and the isolation mechanism (chosen in later ADRs) operates within a single substrate.
- The acknowledged downsides — ISP-outage tenant reachability, harder home-lab provisioning during TR-17 rebuild — are explicitly within the capability's tolerance: there is no availability SLA, and the rebuild KPI is a target with a tracked-issue follow-up when missed (per [stand-up-the-platform §7]({{< relref "../user-experiences/stand-up-the-platform.md" >}})).

The cloud↔home-lab link is a self-hosted VPN, terminated by the operator on both sides; the specific product is deferred to where it is deployed and is not a separate ADR.

### Consequences {#consequences}

- **Good, because** the home-lab carries the substrate the operator most wants to control (tenant compute and tenant data), and the cloud carries only what objectively needs the public internet (ingress) or off-site placement (archive). Each side has a clear reason to exist.
- **Good, because** the cloud surface is small and bounded, which makes the cloud-side definitions (TR-22) compact and easy to audit.
- **Good, because** vendor swap on the cloud edge is a constrained problem (replace ingress + archive on the new vendor) rather than an everything-moves migration.
- **Bad, because** ISP or home-lab power loss takes hosted tenants offline. End users see whatever connection failure the underlying infra produces. This is consistent with the capability's no-SLA stance but should be communicated to capability owners during onboarding.
- **Bad, because** TR-17 rebuild must include home-lab provisioning. The 1-hour budget is tighter against this option than against Option A, and rebuild drills (per [stand-up-the-platform §Entry Point]({{< relref "../user-experiences/stand-up-the-platform.md" >}})) must really be run on parallel home-lab hardware to be honest.
- **Bad, because** the home-lab is a real DR event surface — hardware loss, fire, theft. Mitigation is the cloud-side off-site archive (TR-05) plus the export tooling (TR-14); both must actually work, not just exist.
- **Requires:** ADR #2 (compute substrate) chooses a substrate that runs on home-lab hardware. ADR #4 (storage offering) chooses a storage stack the home-lab can host with the off-site archive in cloud. ADR #6 (network reachability) addresses the cloud-edge ingress mechanism and home-lab egress for archive uploads. ADR #7 (backup & DR) pairs a home-lab primary with a cloud archive tier. ADR #12 (definitions tooling) chooses a tool capable of provisioning *both* sides from the same repo. The cloud↔home-lab VPN is realized as part of Phase 1 (Foundations) of [stand-up-the-platform]({{< relref "../user-experiences/stand-up-the-platform.md" >}}).

### Realization {#realization}

How this decision shows up in the repo:

- **`cloud/`** modules house the cloud-edge offerings only: external ingress, off-site backup archive bucket, the cloud terminator of the self-hosted VPN, and any cloud-side observability or identity components placed by later ADRs.
- **A new directory (e.g. `homelab/`)** houses the home-lab side: compute substrate provisioning, storage provisioning, backup-primary configuration, observability collectors, and the home-lab terminator of the VPN. The exact directory name and tool are decided by ADR #12.
- **A single top-level entry point** (per TR-17) drives both sides; the rebuild flow's Phase 1 (Foundations) is where the cross-boundary VPN is brought up and validated before any offering depending on it is provisioned.
- **`services/`** Go services, where any are needed for platform offerings (e.g. a custom export-tooling service), default to running on the home-lab compute substrate. Cloud-side Go services should be the exception, justified by their offering being placed on the cloud side by this ADR.

## Open Questions {#open-questions}

- **Which cloud provider hosts the cloud edge?** Not decided here. The cloud edge is small enough that this is a follow-on choice with bounded consequences; a one-line decision can name it (or it can be left at "the cloud edge" in the design and named when first deployed).
- **Identity placement.** Deferred to ADR #5. If it lands on a self-hosted product, it can be placed on either side at that point. The cloud edge is reserved for ingress + archive in this ADR; identity is not assumed to live there.
- **Operator-facing observability surface placement.** Deferred to ADR #8. Observability *collection* is on the home-lab side per this ADR; the visualization/alerting front-end may live on either side and is decided with the stack itself.
