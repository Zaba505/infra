---
title: "[0001] Tenant Isolation Model"
description: >
    Choose the isolation primitive (and the layer at which it is enforced) that prevents one tenant from observing another tenant's state, secrets, traffic, or telemetry.
type: docs
weight: 1
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

[TR-01] requires strict tenant isolation at the data and compute layers: no tenant workload may observe or access another tenant's state, secrets, traffic, or telemetry under normal or degraded operation. [TR-03] adds that observability data must be queryable per-tenant within their data scope only. The platform hosts multiple tenant capabilities side-by-side on shared infrastructure (compute, storage, network, observability backends), so isolation must be enforced by the platform itself rather than left to tenant code.

How should the platform structure tenant isolation so that the invariant is enforced by the runtime, not by convention or per-tenant configuration?

## Decision Drivers

* TR-01: tenant isolation is a capability-level invariant, not a best-effort goal
* TR-03: per-tenant observability scoping must follow from the same isolation model so it is not bolted on
* The platform spans public cloud and home-lab infrastructure (capability *Business Rules*: may span public and private)
* Operator maintenance budget KPI (≤ 2 hrs/week) — the model must not require per-tenant manual hardening
* Tiebreaker order: reproducibility > vendor independence > minimizing operator effort
* Failure-mode honesty: "isolation under degraded conditions" forbids models that depend on healthy control planes to stay safe

## Considered Options

* **Per-tenant Kubernetes namespace + NetworkPolicies + RBAC on a shared cluster** — soft isolation primitives layered on a shared kernel and shared control plane.
* **Per-tenant Kubernetes cluster (vcluster or full cluster per tenant)** — control-plane isolation, shared node/host kernel.
* **Per-tenant VM (one VM, or a small VM group, per tenant) on a shared hypervisor** — hardware-virtualization boundary, separate kernels.
* **Per-tenant dedicated host** — physical isolation, no shared kernel or hypervisor.

## Decision Outcome

Chosen option: **per-tenant VM on a shared hypervisor**.

This puts the isolation boundary at the hypervisor — a boundary the platform does not itself implement and that does not weaken under control-plane degradation. Each tenant gets its own kernel, its own filesystem, its own network namespace by default, and its own observability agents whose scope is determined by *which VM they run in* rather than by a label that another component must honor.

Tenant-scoped observability (TR-03) falls out of the model: the agent in tenant A's VM cannot see tenant B's signals because it is not in tenant B's VM. The platform-side aggregation tier enforces the same boundary at query time (see ADR-0004) but the per-VM scoping is the primary defence; the query-time scoping is defence-in-depth.

### Consequences

* Good, because the isolation invariant survives control-plane bugs and misconfigurations; a broken NetworkPolicy or a stale RBAC rule cannot let tenant A read tenant B's disk
* Good, because observability scoping (TR-03) is structural, not policy-driven
* Good, because per-tenant kernel choice / kernel version freedom becomes possible without coordinating with other tenants
* Good, because eviction (capability rule) is a clean teardown: delete the VM and its disk, no cross-tenant entanglement to unwind
* Neutral, because the platform must run a hypervisor or use cloud-managed VMs — both fit C-01 ("may span public and private")
* Bad, because per-VM overhead is higher than per-namespace overhead; small tenants pay for a kernel they do not need
* Bad, because operator-facing primitives shift from `kubectl`-style to VM-lifecycle-style, which has implications for ADR-0002 (compute substrate)

### Confirmation

* Each tenant's runtime artifacts (disk image, network interface, observability agent identity) are namespaced by VM identity in the platform's definitions repo
* A tenant-isolation regression test in the standup canary (REQ-18 from the longer requirements doc) attempts cross-tenant reads from tenant A's VM and asserts they fail at the network and storage layers
* Observability backend queries are validated to scope by VM identity — see ADR-0004

## Pros and Cons of the Options

### Per-tenant Kubernetes namespace + NetworkPolicies + RBAC

Soft isolation on a shared cluster.

* Good, because it is the lowest-overhead option per tenant
* Good, because tooling (`kubectl`, Helm) is well-known
* Bad, because isolation depends on every NetworkPolicy and RBAC rule being correct, on every controller honoring them, and on the kernel not being a shared blast surface — TR-01 says "under any normal or degraded operating condition" and a shared kernel violates that
* Bad, because a CRD or admission-controller bug can break isolation across all tenants at once
* Bad, because per-tenant observability scoping has to be enforced by labels that other components must respect

### Per-tenant Kubernetes cluster (vcluster or full cluster per tenant)

Control-plane isolation per tenant; shared host kernel via shared nodes (vcluster) or per-tenant nodes.

* Good, because control-plane bugs are scoped to one tenant
* Good, because tenants can have independent CRDs / operators
* Neutral, because if nodes are shared the host kernel is still a shared blast surface; if nodes are per-tenant the model approaches per-VM anyway
* Bad, because the operator now runs N control planes — pushes against NFR-02 (maintenance budget)
* Bad, because observability scoping is still policy-driven within each cluster's view of the shared aggregation tier

### Per-tenant VM on a shared hypervisor

* Good, because isolation is enforced by the hypervisor — a boundary the platform inherits rather than implements
* Good, because tenant-scoped observability falls out of the topology
* Good, because eviction and rebuild are clean
* Good, because it works identically on home-lab hardware and on cloud-managed VMs (C-01)
* Neutral, because per-tenant overhead is higher than per-namespace
* Bad, because small tenants pay fixed VM overhead
* Bad, because the operator-facing model is VM-lifecycle, not container-orchestration

### Per-tenant dedicated host

* Good, because there is no shared substrate at all
* Bad, because cost and operator effort scale linearly with tenants — incompatible with NFR-02 and NFR-05
* Bad, because reproducibility (NFR-01: ≤ 1 hr full rebuild) is hard when each tenant requires hardware procurement

## More Information

* This ADR sets the boundary; ADR-0002 picks the compute substrate that runs *inside* each tenant VM and the orchestration layer that runs *across* them
* TR-07 (Cloudflare → GCP path) is unaffected — VMs sit behind the same ingress topology
* Out of scope: which hypervisor (KVM vs. Firecracker vs. cloud-managed VMs) — that is part of ADR-0002
