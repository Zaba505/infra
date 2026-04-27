---
title: "[0005] Identity Offering — Self-Hosted Authentik"
description: >
    The platform offers tenants a self-hosted Authentik instance running in the home-lab cluster, with realms as the per-tenant boundary and per-realm configuration providing the "lost credentials cannot be recovered" property required by the capability.
type: docs
weight: 5
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform]({{< relref "../_index.md" >}})
**Addresses requirements:** [TR-04]({{< relref "../tech-requirements.md#tr-04" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}), [TR-17]({{< relref "../tech-requirements.md#tr-17" >}}), [TR-22]({{< relref "../tech-requirements.md#tr-22" >}}), [TR-33]({{< relref "../tech-requirements.md#tr-33" >}})

## Context and Problem Statement {#context}

[TR-04]({{< relref "../tech-requirements.md#tr-04" >}}) requires the platform to offer tenants an identity-and-authentication service for *their* end users, and constrains that service to be capable of honoring "lost credentials cannot be recovered" (Signal-style). This is a hard gate: any product that cannot be configured to enforce no-recovery is ineligible.

[ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}}) reserves the cloud edge for ingress and the off-site backup archive only, so identity, if self-hosted, lives in the home-lab cluster. [ADR-0002]({{< relref "0002-compute-substrate.md" >}}) makes the home-lab a Kubernetes cluster with namespaces as the per-tenant boundary; whatever IdP we choose has to plug into that.

The capability's vendor-independence tiebreaker presses against vendor-managed identity products. Identity is precisely the kind of foundational offering that, once tenants are on it, is painful to leave — making the tiebreaker more salient here than for offerings whose data is more portable.

The decision is also potentially cross-capability — other capabilities the operator defines later may want to consume the same identity service. This ADR keeps the decision capability-scoped for now, with the explicit caveat that it is a candidate for lifting to a shared ADR if other capabilities standardize on it.

## Decision Drivers {#decision-drivers}

- **TR-04 hard gate** — no-recovery property must be configurable. This rules out products without that capability; it does not select among the rest.
- **Capability tiebreaker — vendor independence > minimizing operator effort.** Identity is foundational; vendor lock-in here is harder to unwind than at most other layers.
- **TR-33 (≤2 hr/week).** The IdP is a stateful product running continuously. A heavy IdP eats weekly budget directly.
- **TR-17 (≤1 hr rebuild).** The IdP and its persistent state must come up within Phase 2 of the rebuild ([stand-up-the-platform §4]({{< relref "../user-experiences/stand-up-the-platform.md" >}})) without blowing the budget.
- **TR-32 isolation.** Tenants must not be able to read each other's identity data. The IdP's per-tenant-separation primitive (realm, organization, tenant, …) is the relevant boundary.
- **TR-22 / TR-24 definitions-driven.** IdP configuration (realms, providers, flows) must be expressible declaratively and tracked, not clicked into a console.
- **Capability rule "evolves with its tenants" — modulated.** For shape-y offerings (storage, ADR-0004) deferral works. For identity, a cross-cutting protocol that tenants' tech designs commit to at design time, deferral creates path-dependence that costs more to unwind later than the upfront cost saves.

## Considered Options {#considered-options}

### Option A — Self-hosted Keycloak in the home-lab cluster

Mature, full-featured OIDC/SAML provider. Realms are the per-tenant boundary. Per-realm flow configuration provides no-recovery (disable email-based recovery, lock recovery flows).

- **TR-04:** met — no-recovery is realm-configurable.
- **TR-32:** met — realms isolate tenants by design.
- **TR-22 / TR-24:** met — Keycloak supports declarative realm configuration (operators, JSON imports, terraform providers).
- **TR-33:** the highest-cost option — JVM, large surface area, version upgrades are non-trivial events; stateful with a Postgres dependency.
- **TR-17:** rebuild bootstrap is heavier than Authentik; not a deal-breaker, but presses the budget.
- **Vendor independence:** strong — open-source, large community.

### Option B — Self-hosted Authentik in the home-lab cluster

Modern open-source IdP. Tenants are configurable as separate Authentik *tenants* / providers; flows can be configured per scope so the no-recovery property is realizable.

- **TR-04:** met — recovery flows are declaratively configurable to omit credential recovery entirely.
- **TR-32:** met — Authentik's tenant/provider model gives the per-tenant boundary; combined with the K8s namespace-per-tenant boundary the isolation reasoning is uniform with the rest of ADR-0002.
- **TR-22 / TR-24:** met — Authentik exposes declarative configuration, including a Kubernetes operator for managing core resources.
- **TR-33:** materially lighter than Keycloak — single Python/Go process pair, modern UI, simpler upgrades. Real but bounded weekly cost.
- **TR-17:** lighter bootstrap than Keycloak; still includes the stateful DB dependency (see Realization).
- **Vendor independence:** strong — open-source, active community.

### Option C — A lighter / pluggable IdP (Dex + an account store, or Ory Kratos + Hydra, etc.)

Smaller, more composable identity primitives.

- **TR-04:** depends on the combination. Dex alone is a federation broker, not a primary IdP — it does not own credentials, so "credentials cannot be recovered" doesn't apply at its layer. Ory Kratos + Hydra can be configured to omit recovery flows.
- **TR-32:** workable but requires more careful design than Authentik/Keycloak's built-in tenancy.
- **TR-22 / TR-24:** good — these projects are designed declarative-first.
- **TR-33:** smaller per-component surface, but more components to operate (Kratos + Hydra + Oathkeeper-or-equivalent) — net surface area is comparable to or higher than Authentik.
- **Vendor independence:** strong.
- **Cost:** more wiring to assemble a complete tenant-facing IdP than the all-in-one options.

### Option D — Cloud-hosted IdP (Auth0, Okta, AWS Cognito, similar)

A managed identity vendor at the cloud edge.

- **TR-04:** generally met — these products support no-recovery configuration.
- **TR-32:** met — vendors model tenants explicitly.
- **TR-22 / TR-24:** met via vendor IaC.
- **TR-33:** lowest weekly cost — somebody else's pager.
- **Vendor independence:** **fails the tiebreaker** — every tenant who adopts the platform-provided identity is now on this vendor's roadmap and pricing. Identity is precisely the offering whose data and integrations are *most* costly to migrate later. The capability ranks vendor independence above operator effort, so this option's headline benefit (no operator effort) is on the losing side of that tiebreaker.

### Option E — No platform-provided identity offering on day one; require all tenants to bring their own

The strict *evolves with its tenants* reading: defer the offering until a tenant pulls it.

- **TR-04:** technically deferred. The TR is forced; deferring it indefinitely puts the platform out of compliance with the capability's named output.
- **Path-dependence cost:** every early tenant designs against bring-your-own-identity, which makes later adoption of the platform-provided offering an explicit migration per tenant. Identity is unlike storage shapes here — storage shapes can be retrofitted per tenant; identity is woven into a tenant's design.
- **Cost / TR-33:** lowest day-one cost; offset by likely-higher eventual cost when the offering is finally drafted under pressure from a real tenant.

## Decision Outcome {#decision-outcome}

Chosen option: **Option B — Self-hosted Authentik in the home-lab cluster**, with Authentik tenants/providers as the per-tenant boundary on top of the K8s namespace boundary, and per-tenant flow configuration realizing the "lost credentials cannot be recovered" property.

This option is chosen because:

- It passes TR-04's hard gate (no-recovery is declaratively configurable).
- It satisfies the capability's vendor-independence tiebreaker, which Option D fails most clearly for the offering most costly to leave later.
- It is materially lighter than Keycloak (Option A) on TR-33, with a small but real edge on TR-17 rebuild bootstrap. Keycloak's feature breadth is not repaid at personal scale.
- It is a coherent all-in-one product, avoiding Option C's cost of assembling a tenant-facing IdP from primitives.
- It avoids Option E's path-dependence cost. Identity is a cross-cutting protocol that tenants' tech designs commit to at design time; deferring the offering pushes every early tenant onto bring-your-own and makes later platform-provided adoption a per-tenant migration.

The IdP runs in its own Kubernetes namespace in the home-lab cluster. Per-tenant configuration is captured in the platform definitions and applied declaratively (via Authentik's Kubernetes operator or equivalent), satisfying TR-22 / TR-24.

Authentik depends on Postgres. Per [ADR-0004]({{< relref "0004-persistent-storage-offering.md" >}}), the platform does not yet offer a tenant-facing relational storage offering — managed offerings appear when a *tenant* demands them. Authentik is not a tenant; it is a platform offering. The Postgres it requires is *internal to the identity offering* — it ships in the same namespace as Authentik, is not exposed to tenants, and is operated as part of the IdP's own footprint. The day a tenant needs Postgres-shaped storage, ADR-0004's trip-wire fires and a tenant-facing relational managed-offering ADR is drafted; that future offering may or may not consolidate with Authentik's internal Postgres, but that is a decision for that ADR, not this one.

### Consequences {#consequences}

- **Good, because** TR-04's no-recovery property is realized by configuration we own, not negotiated with a vendor.
- **Good, because** Authentik runs on the same substrate as every other home-lab offering (TR-32 isolation reasoning is uniform with the rest of the design — namespace + RBAC + NetworkPolicy).
- **Good, because** the platform retains the option to swap to another self-hosted IdP later. The cost of swapping is bounded by how many tenants are on the platform-provided identity at the time, which the contract-change UX already governs ([UX #5]({{< relref "../user-experiences/platform-contract-change-rollout.md" >}})).
- **Good, because** tenants who need stronger guarantees than Authentik provides can bring their own per [TR-04]({{< relref "../tech-requirements.md#tr-04" >}})'s explicit BYO carve-out — the platform-provided offering is not the only path.
- **Bad, because** the platform now operates a stateful identity product. Authentik plus its Postgres is a real surface on TR-33; upgrades require attention; backup of the IdP's state must be in the backup scope (ADR #7).
- **Bad, because** the rebuild (TR-17) must bring up Authentik *and* its Postgres in Phase 2 alongside compute and storage. The 1-hour budget is tighter for it.
- **Bad, because** Authentik's internal Postgres is a *second* DB engine on the platform when ADR-0004's tenant-facing relational ADR is eventually drafted. That future ADR may consolidate or may not; for now there is a small, deliberate duplication between "Authentik's DB" and "future tenant DB offering."
- **Bad, because** TR-04's BYO option means the platform must also support tenants whose tech designs name a different IdP. The schema (ADR-0003) must permit BYO without the platform doing anything for those tenants beyond network reachability — i.e. "use platform identity" is a tenant manifest field, not the assumption.
- **Requires:**
  - ADR #6 (network) provides ingress to Authentik so end users of tenant capabilities can reach it; NetworkPolicy isolates Authentik's Postgres so only Authentik pods can reach it.
  - ADR #7 (backup & DR) includes Authentik's Postgres in the backup scope.
  - ADR #9 (secrets) holds Authentik admin/bootstrap credentials and any per-tenant signing keys it needs.
  - ADR-0003's schema gains a field expressing which identity the tenant uses (`platform-authentik` vs. `byo` with reachability info).
  - The canary tenant (TR-20, ADR #15) authenticates against the platform-provided identity service end-to-end so identity is exercised on every rebuild.

### Realization {#realization}

How this decision shows up in the repo:

- **An `authentik` namespace** in the home-lab cluster, provisioned as part of Phase 2 of the rebuild ([stand-up-the-platform §4]({{< relref "../user-experiences/stand-up-the-platform.md" >}})), holds:
  - The Authentik server and worker Deployments (or whatever Authentik's recommended layout is at the chosen version).
  - A Postgres workload (StatefulSet or operator-managed) backing Authentik. Its PVCs use the block storage class from ADR-0004.
  - A NetworkPolicy permitting only Authentik pods to reach the Postgres.
  - The Authentik Kubernetes operator (or equivalent declarative configuration mechanism) so per-tenant realms/providers/flows are managed via tracked manifests, satisfying TR-22 / TR-24.
- **Per-tenant Authentik configuration** lives in the platform definitions as declarative manifests, alongside the per-tenant K8s namespace and quota. The translator from ADR-0003 emits or references these when a tenant manifest names `platform-authentik`.
- **Ingress to Authentik** flows through the cloud-edge ingress (ADR #6) and across the cloud↔home-lab VPN (ADR-0001).
- **Authentik's Postgres** is in the backup scope (ADR #7) on equal footing with tenant data — its loss is recoverable.
- **The canary tenant** (ADR #15) authenticates against the platform-provided identity service, so a broken Authentik fails the canary and prevents the rebuild from declaring readiness.

## Open Questions {#open-questions}

- **Authentik version pinning policy.** Deferred to deployment-time. The platform's contract for tenants (which Authentik features they can rely on) is captured in the schema (ADR-0003), not in a pinned version field.
- **Will any other operator capability adopt this identity service?** If yes, this ADR is a candidate for lifting to a shared ADR in `r&d/adrs/`. Out of scope here.
- **Should the platform later offer a managed Postgres offering, and if so, should Authentik consolidate onto it?** Future ADR-0004-trip-wire territory. For now, Authentik's Postgres is internal to the identity offering.
- **Recovery property tooling for tenants.** TR-04 says the offering must be *capable* of honoring no-recovery; per-tenant configuration is what realizes it. Whether the schema (ADR-0003) defaults to no-recovery on or off, and whether tenants opt in or out, is a small follow-on field decision documented with the schema, not in this ADR.
