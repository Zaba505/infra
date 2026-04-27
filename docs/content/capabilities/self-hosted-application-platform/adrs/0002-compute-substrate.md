---
title: "[0002] Compute Substrate"
description: >
    Tenant workloads run on a Kubernetes cluster on home-lab hardware, using a small distro (k3s/k0s class). Namespaces are the per-tenant boundary.
type: docs
weight: 2
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform](../_index.md)
**Addresses requirements:** [TR-01](../tech-requirements.md#tr-01-provide-compute-as-a-tenant-offering), [TR-12](../tech-requirements.md#tr-12-provide-a-one-shot-migration-process-offering-that-runs-tenant-supplied-migration-jobs), [TR-13](../tech-requirements.md#tr-13-admit-migration-jobs-only-when-their-peak-temporary-footprint-is-at-most-2-the-destination-tenants-steady-state-compute-and-storage), [TR-17](../tech-requirements.md#tr-17-definitions-driven-single-entry-point-rebuild-of-the-platform-end-to-end), [TR-19](../tech-requirements.md#tr-19-each-rebuild-phase-must-be-cleanly-torn-downable-the-partial-state-is-itself-untrusted), [TR-20](../tech-requirements.md#tr-20-maintain-a-purpose-built-canary-tenant-alongside-the-platforms-definitions-and-use-it-as-the-readiness-signal), [TR-22](../tech-requirements.md#tr-22-tracked-changes-and-immutability-for-all-platform-state-modifying-actions), [TR-24](../tech-requirements.md#tr-24-tenant-provisioning-must-run-only-through-the-platforms-existing-definitions), [TR-25](../tech-requirements.md#tr-25-during-platform-contract-change-rollouts-run-old-and-new-forms-of-the-offering-concurrently-until-the-deadline-except-where-the-change-is-a-full-removal), [TR-26](../tech-requirements.md#tr-26-tenants-declare-resource-needs-at-onboarding-and-on-every-modify-the-platform-admits-or-refuses-based-on-those-declarations), [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals), [TR-33](../tech-requirements.md#tr-33-routine-platform-operation-must-fit-within-2-hoursweek-of-operator-time)

## Context and Problem Statement

[ADR-0001](./0001-public-private-infrastructure-split.md) places tenant compute on the home-lab side. This ADR picks the substrate that actually runs tenant workloads there. It is the single biggest downstream constraint in this design: it dictates what tenant packaging form is plausible (ADR #3), how the migration runner is built (ADR #10), what observability and identity hooks are available (ADRs #5, #8), and how per-tenant isolation (TR-32) is enforced.

Many TRs touch this decision because the substrate either provides their primitives natively or pushes them onto bespoke implementation: TR-12 needs a one-shot job runner; TR-25 needs concurrent old-and-new offerings; TR-26 needs declared-resource admission; TR-32 needs an isolation boundary; TR-19 needs cheap, complete teardown; TR-22 / TR-24 need a definitions-driven provisioning surface.

## Decision Drivers

- **TR-32 isolation** — the substrate's native boundary determines whether per-tenant isolation is the default or has to be bolted on.
- **TR-12 migration runner** — the migration-process offering is a one-shot job runner integrated with secrets and observability. A substrate with a Job primitive shrinks the offering's implementation to mostly configuration.
- **TR-25 concurrent old/new** — contract-change rollouts need two versions of an offering serving simultaneously. A substrate with first-class deployment selectors makes this a configuration concern, not a coordination project.
- **TR-26 / TR-13 admission** — tenants declare resource needs and migration jobs admit at ≤2× steady state. A substrate with quota primitives turns this into manifest, not code.
- **TR-19 teardown** — "delete everything" must be reliable at every phase boundary. The substrate's teardown semantics directly answer this.
- **TR-17 rebuild ≤1 hour** — substrate bootstrap eats the rebuild budget. Heavier substrates lose this directly.
- **TR-33 ≤2 hr/week** — the substrate's routine maintenance cost (patching, upgrades, drift handling) is paid weekly. More features ≠ more maintenance, but bigger surface area generally does.
- **Capability tiebreaker — vendor independence > minimizing operator effort.** A substrate the operator runs themselves on owned hardware satisfies the tiebreaker over a managed-cluster product.

## Considered Options

### Option A — Bare-metal VMs, one (or a small fixed pool) per tenant

A hypervisor (e.g. Proxmox, libvirt+KVM) on home-lab nodes; each tenant gets one or more VMs.

- **TR-32:** strongest isolation of any option — separate kernels.
- **TR-12 / TR-25 / TR-26:** none of these are first-class. A migration-job offering, concurrent-deploy mechanic, and resource-quota admission would all be bespoke.
- **TR-19:** clean — VMs delete completely.
- **TR-17:** hard. Every VM needs an OS install path, and the rebuild has to bring up the hypervisor *and* the tenants' OS images.
- **TR-33:** heavy. Patching and upgrades happen per VM, multiplied by tenants.
- **Packaging form (downstream ADR #3):** would have to be a VM image, which is a much heavier handoff than capability owners are realistically producing.

### Option B — Single-node container runtime (Docker / Podman) on one or a few home-lab boxes, no orchestrator

- **TR-32:** kernel-shared isolation; weaker than VMs and weaker than namespaces+NetworkPolicy in K8s without significant hardening.
- **TR-12 / TR-25 / TR-26:** all hand-rolled. No native job, no native blue-green selector, no native quota admission.
- **TR-19:** very clean — `docker rm` is fast and complete.
- **TR-17:** very fast bootstrap.
- **TR-33:** low *for a small number of tenants*; rises sharply once cross-cutting concerns (observability wiring, secret injection, identity, network policy) accumulate, because each becomes a bespoke pattern instead of an offered primitive.

### Option C — Kubernetes on the home-lab using a small distro (k3s / k0s class)

A single home-lab cluster using a small Kubernetes distribution, with **namespaces as the per-tenant boundary**.

- **TR-32:** namespace + NetworkPolicy + RBAC is a coherent isolation answer ADRs #5 (identity) and #6 (network) can lean on.
- **TR-12:** the `Job` primitive *is* the migration runner; concurrent runs across tenants are free.
- **TR-25:** Deployment + Service + label selectors give concurrent old/new without coordination.
- **TR-26 / TR-13:** ResourceQuota and LimitRange enforce admission; the 2× rule is a quota expression, not custom code.
- **TR-19:** namespace deletion cascades; "delete everything" is one operation per tenant; cluster teardown likewise.
- **TR-17:** small-distro bootstrap is on the order of minutes; the rebuild budget is *tight* but not blown — and most of what would otherwise eat the budget (offering wiring) is configuration on top of a bootstrapped cluster, not separate provisioning.
- **TR-22 / TR-24:** declarative manifests are the substrate's native input; tracking changes via git and immutability via admission policy are well-trodden patterns.
- **TR-33:** the chief risk here. Mitigated by picking a small distro (one binary, embedded etcd or sqlite, no add-ons by default) rather than full kubeadm with a stack of operators. Maintenance cost is real but bounded.
- **Packaging form (downstream ADR #3):** OCI image + a small manifest is the natural form, which is what capability owners are realistically producing.

### Option D — Nomad on the home-lab

HashiCorp Nomad as the orchestrator.

- **TR-32:** weaker than K8s namespaces; isolation primitives exist but are less integrated with network and identity.
- **TR-12 / TR-25 / TR-26:** Nomad has jobs and resource constraints, but the ecosystem for admission control, identity integration, and observability is thinner than K8s.
- **TR-19:** clean.
- **TR-17:** lighter bootstrap than K8s.
- **TR-33:** comparable to a small K8s distro for the platform itself, but each cross-cutting concern (TR-04 identity, TR-06 observability, TR-11 secrets, TR-15 CSI-equivalent) requires more bespoke wiring than the K8s ecosystem ships.
- **Packaging form:** flexible (containers, VMs, raw exec), which is a feature for some platforms — but our packaging-form decision wants a *single* form (TR-10), so the flexibility is unused weight.

## Decision Outcome

Chosen option: **Option C — Kubernetes on the home-lab using a small distro (k3s / k0s class)**, with a single cluster per home-lab site, and **namespaces as the per-tenant boundary**.

This option is chosen because the count of TRs Kubernetes answers natively materially exceeds the count Options B and D leave to bespoke implementation:

- TR-12 (one-shot migration runner) is the `Job` primitive, with concurrent jobs across tenants free.
- TR-25 (concurrent old/new during rollout) is a Deployment+Service+selector pattern, not a coordination project.
- TR-26 / TR-13 (declared-resource admission + 2× migration cap) are ResourceQuota/LimitRange expressions.
- TR-32 (isolation) is namespaces + NetworkPolicy + RBAC, which subsequent ADRs (#5 identity, #6 network) can compose with.
- TR-19 (clean teardown) is namespace deletion (per-tenant) and cluster teardown (whole-platform).
- TR-22 / TR-24 (definitions-driven, tracked, immutable) is the substrate's native input model.

Each of those TRs would be hand-rolled under Option B (no orchestrator) and partially-rolled under Option D (Nomad), and hand-rolled answers tend to break TR-33 (2 hr/week) long before they break functionally. Option A (VMs) loses on TR-12, TR-25, TR-26 simultaneously and forces a packaging form heavier than capability owners are realistically producing.

The chief cost of Option C — operator-held control-plane surface area — is bounded by picking a small distro (one binary, embedded datastore, no add-on operators by default) rather than full kubeadm. This trades against TR-33; the trade is acceptable because the alternative — re-implementing the K8s primitives above as platform code — costs more weekly than running a small cluster.

The capability's tiebreaker (*vendor independence > minimizing operator effort*) is honored: the cluster runs on home-lab hardware the operator owns, and the substrate is open-source software with multiple distros so swapping one for another later is bounded work.

### Consequences

- **Good, because** every per-tenant offering — compute, storage, network, identity, secrets, observability — has a coherent home (a namespace) and a coherent isolation story (RBAC + NetworkPolicy scoped to that namespace).
- **Good, because** the migration-process offering (TR-12) is mostly configuration on top of `Job`, and concurrent migrations across tenants come for free.
- **Good, because** contract-change rollouts (TR-25) are realized as label-selector switches, not as parallel infrastructure.
- **Good, because** "delete everything" (TR-19) is `kubectl delete namespace <tenant>` per tenant and a cluster wipe at the platform level — fast and reliable.
- **Bad, because** the rebuild budget (TR-17) now includes K8s bootstrap. A small distro keeps this on the order of minutes, but it is not free. Phase 2 of [stand-up-the-platform](../user-experiences/stand-up-the-platform.md#4-phase-2--core-platform-services) must absorb this without blowing the 1-hour KPI.
- **Bad, because** the operator carries Kubernetes concepts (namespaces, RBAC, NetworkPolicy, CRDs, controllers) in their head as part of routine operation. TR-33 (2 hr/week) is pressured by this and is the failure mode to watch for.
- **Bad, because** namespace-based isolation is weaker than per-tenant-VM isolation. TR-32 holds *given* correctly configured RBAC + NetworkPolicy + (later) Pod Security; misconfiguration is a real risk and the cost of the chosen substrate.
- **Requires:** ADR #3 (packaging form) lands on OCI image + manifest, which is the substrate's native input. ADR #4 (storage) chooses a CSI-compatible storage stack. ADR #5 (identity) chooses something that integrates with K8s service-account-style identity for in-cluster needs. ADR #6 (network) chooses a CNI and an ingress mechanism. ADR #8 (observability) chooses a stack that integrates with K8s discovery. ADR #9 (secrets) chooses a secret store K8s pods can read by name. ADR #10 (migration runner) is a thin layer on `Job`. ADR #12 (definitions tooling) must drive K8s manifests into the cluster (and must not depend on the cluster being up, since the cluster is itself part of what is rebuilt). ADR #13 (rebuild orchestration) sequences cluster bootstrap before any per-tenant resources.

### Realization

How this decision shows up in the repo:

- **`homelab/`** (or the equivalent placed by ADR #12) houses the cluster definition: node configuration, the small-distro install, the cluster's bootstrap token / join configuration, and the in-cluster baseline (RBAC roles for the operator, default-deny NetworkPolicy, a ResourceQuota / LimitRange template).
- **Per-tenant manifests** live alongside the platform definitions — one namespace per tenant, with a ResourceQuota matching the tenant's declared resource needs (TR-26) and a LimitRange that, in combination with the migration-job pattern from ADR #10, encodes the 2× admission cap (TR-13).
- **The canary tenant (TR-20)** is a manifest in the platform definitions producing one namespace with a trivial workload that exercises every offering and tears down cleanly.
- **`services/`** Go services (where the platform itself ships custom services — e.g. potentially the export-tooling backend) run as Deployments on the same cluster, alongside tenant workloads in their own namespace.
- **Cluster bootstrap** is Phase 2 of [stand-up-the-platform](../user-experiences/stand-up-the-platform.md#4-phase-2--core-platform-services); the cluster is *not* assumed to exist when ADR #12's tooling runs against the home-lab side, so the tooling must drive the install of the small distro itself before applying any in-cluster resources.

## Open Questions

- **Which small distro** (k3s, k0s, microk8s, or another). Deferred to the deployment-time decision; this ADR commits to "small-distro class" as the constraint.
- **Single cluster vs. multi-cluster.** This ADR assumes one cluster per home-lab site. Multi-cluster (e.g. for blue/green of the platform itself, or to separate experimental tenants) is a deferred concern; revisit when there is a real reason.
- **Pod Security baseline.** The level of restriction (restricted vs. baseline vs. privileged for special cases) is a follow-on decision, likely settled while drafting ADR #6 (network) and ADR #9 (secrets). It is a TR-32 hardening choice.
