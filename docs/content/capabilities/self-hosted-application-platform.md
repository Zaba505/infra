---
title: "Self-Hosted Application Platform"
description: >
    A reproducible, operator-controlled platform that hosts the operator's other capabilities, so each one does not have to solve hosting on its own or fall back to a vendor.
type: docs
weight: 10
---

> **One-line definition:** Provide a reproducible, operator-controlled platform on which the operator's other capabilities run by default, so that no capability has to depend on a vendor-specific hosting solution to be delivered.

## Purpose & Business Outcome
What business outcome does this capability deliver? Why does it exist?

This capability exists so that the operator's other capabilities (e.g. self-hosted personal media storage) have a well-defined, reproducible place to run that the operator controls end-to-end, instead of each capability independently choosing a vendor (e.g. a hosted Plex provider, a hosted Minecraft provider, a hosted Nextcloud provider). The outcomes it delivers, in order of importance:

1. **Default hosting target for the operator's capabilities.** Any capability the operator defines should be able to run here, so that "where does this run?" is a solved question rather than re-litigated per capability.
2. **Reproducibility.** The platform itself can be rebuilt from its definitions; it is not a snowflake. A total loss does not mean a permanent loss of the platform.
3. **Independence from hosting vendors.** The operator is not locked into any single provider's product roadmap, pricing, or terms for the things their capabilities depend on.
4. **A coherent place to invest infrastructure effort.** Improvements (resiliency, observability, backup) made once at the platform level benefit every tenant capability, instead of each capability re-solving them.

When these outcomes conflict: tenant adoption beats reproducibility (a perfect platform with no tenants is a failure); reproducibility beats vendor independence (a platform that can't be rebuilt is worse than one that uses some vendor components); vendor independence beats minimizing operator effort.

## Stakeholders

- **Owner / Accountable party:** The operator (Carson). Sole accountable party for the platform existing, running, and continuing to run.
- **Primary actors (initiators):** Capability owners — currently the operator wearing a different hat — who bring a capability to the platform to be hosted, or change what an already-hosted capability needs.
- **Secondary actors / consumers:** The tenant capabilities themselves, while running, consume platform services (compute, storage, network, identity, backup, observability).
- **Affected parties (impacted but not directly involved):** End users of the tenant capabilities (e.g. family and friends using self-hosted personal media storage). They never interact with the platform directly, but a platform outage or data loss directly affects them.

## Triggers & Inputs
What initiates the capability, and what information must be available?

- **Triggers:**
  - A capability owner brings a new capability to be hosted.
  - A capability owner changes the requirements of an already-hosted capability (more storage, new external endpoint, etc.).
  - The operator stands up the platform from scratch (initial build or full rebuild after loss).
  - The operator performs routine maintenance on the platform.
  - A tenant capability's components fall behind what the platform supports and need to be updated.
- **Required inputs:**
  - From the capability owner: the capability packaged in the form the platform accepts, a declaration of its resource needs (compute, storage, network reachability), and its availability expectations.
  - For tenants whose end users need to authenticate: either use of the platform-provided identity service, or a declared decision to bring their own.
- **Preconditions:**
  - The operator has authorized the capability to run on the platform (no self-onboarding by tenants — the operator is the only person making this decision).
  - The capability accepts the platform's contract (see Business Rules).

## Outputs & Deliverables
What does the capability produce? What changes in the world after it runs?

- **Direct outputs:** For each tenant capability, the platform provides:
  - **Compute** — a place for the application to run.
  - **Persistent storage** — durable storage for the application's data.
  - **Network reachability** — both internal (between tenants) and external (reachable by the tenant's end users).
  - **Identity & authentication for end users** — available to any tenant that wants it; tenants may opt to bring their own.
  - **Backup and disaster recovery** — of tenant data, to a standard the platform defines.
  - **Observability** — the operator can tell whether each tenant is up and healthy without the tenant having to instrument that itself.
- **Downstream effects / state changes:**
  - The operator's capabilities have a default answer to "where does this run?" and stop being individually coupled to vendor choices.
  - Investments in resiliency, backup, and observability accrue across all tenants instead of being repeated per capability.
  - The operator accumulates operational knowledge of one platform rather than fragmented knowledge of many vendor products.

## Business Rules & Constraints

- **Default hosting target.** All capabilities defined in this repo are expected to run on the platform unless explicitly exempted. A capability owner may choose to host elsewhere, but the platform is the default and the burden of justification is on opting out.
- **Operator-only operation.** Only the operator operates the platform and has administrative access to it. There are no co-operators and no delegated administration.
- **Tenants must accept the platform's contract.** To be hosted, a tenant must be packaged in the form the platform accepts, declare its resource needs up front, and accept the platform's availability characteristics. A tenant that needs guarantees stronger than the platform offers must host elsewhere.
- **Eviction is allowed when needs and capabilities diverge.** The platform may decline to continue hosting a tenant whose requirements it cannot meet (e.g. specialized hardware, regulatory constraints, an availability target the platform does not offer). However, where the divergence is merely that the tenant's components have fallen behind what the platform supports, the platform works with the tenant to bring them current rather than evicting.
- **The platform may span public and private infrastructure.** "Self-hosted" means the operator controls the platform end-to-end, not that every component runs on hardware the operator owns. Public-cloud components are allowed where the operator retains control of configuration, data, and the ability to leave.
- **No direct end-user access to the platform.** End users of tenant capabilities reach the tenant, not the platform. The platform has no notion of "end users" of itself; its consumers are tenant capabilities (and behind them, the operator).
- **Cost is secondary to convenience and resiliency.** Because there is one operator, added cost is acceptable when it buys meaningful convenience or resiliency. Cost should still be minimized where it does not cost convenience or resiliency.
- **The capability evolves with its tenants.** When a tenant capability needs something the platform does not yet provide, the default response is to update this capability's definition (and the platform) rather than push the requirement back onto the tenant.

## Success Criteria & KPIs

- **Tenant adoption.** Every capability defined in this repo runs on this platform. A capability defined in this repo that runs elsewhere is, by default, a failure of this capability — either the platform did not meet the tenant's needs, or the tenant was never asked to use it.
- **Reproducibility.** The platform can be stood up from its definitions in **at most 1 hour**, starting from no platform at all. This is the operational form of "reproducible" — if it takes longer than that, the platform is a snowflake regardless of how much of its config is in version control.
- **Operator maintenance budget.** Routine operation of the platform takes **no more than 2 hours per week** of the operator's time. If maintenance regularly exceeds this, the platform is consuming more attention than it earns and must be simplified, not grown.
- **Cost stays proportional to value.** Total operating cost remains within what the operator considers acceptable given the convenience and resiliency it delivers. There is no fixed dollar target; the test is whether the operator would still choose to run it knowing the bill.

## Out of Scope

- **Hosting for anyone other than the operator's own capabilities.** The platform does not offer hosting to third parties, the public, or family/friends directly. Family and friends reach the platform only as end users of a tenant capability (e.g. via self-hosted personal media storage), never as platform users.
- **Dictating the implementation.** "Homelab," "Kubernetes," and any specific stack are possible implementations of this capability, not part of its definition. The capability is satisfied by anything that meets its rules and KPIs.
- **A specific availability or performance SLA.** The platform offers whatever availability its current implementation can deliver within the operator's maintenance budget. Tenants needing stronger guarantees host elsewhere (per Business Rules).
- **End-user-facing features of tenant capabilities.** Photo viewing, game server gameplay, document editing, etc. are tenant concerns, not platform concerns.
- **Multi-operator administration, role delegation, or self-service onboarding.** Explicitly excluded by the operator-only rule.

## Open Questions

- **Operator skill development as an explicit goal.** Building and running the platform is also a way for the operator to learn and practice infrastructure skills. Is that a stated outcome of the capability (which would change how trade-offs like "buy vs. build" are judged), or an unstated personal benefit that should not influence the capability's success criteria? To be decided.
- **Compatibility with tenant identity rules.** Self-hosted personal media storage requires that lost user credentials cannot be recovered (Signal-style). If a tenant uses the platform-provided identity service, the platform must be capable of honoring that property. Confirm this is true of any identity implementation chosen.
- **Operator succession.** The personal-media-storage capability raises a "hit by a bus" question for longevity. Since this platform underlies that capability, the same question applies here — possibly more sharply, since no one else operates it. To be addressed jointly with that capability's open question.
- **Definition of "tenant adoption."** The KPI says every capability in this repo runs on the platform. Does a capability that is *defined* but not yet *implemented* count toward adoption, or only implemented capabilities? Clarify before measuring.
- **Threshold for eviction vs. accommodation.** The rules say the platform works with tenants on out-of-date components but may evict on fundamental mismatches. Where exactly that line sits will need real cases to calibrate.
