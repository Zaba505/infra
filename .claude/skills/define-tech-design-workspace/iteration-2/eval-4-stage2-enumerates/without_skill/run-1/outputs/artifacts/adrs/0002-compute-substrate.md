---
title: "[0002] Compute Substrate"
description: >
    Choose the compute substrate that orchestrates tenant VMs across home-lab and cloud, given the per-tenant-VM isolation model from ADR-0001.
type: docs
weight: 2
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

ADR-0001 chose per-tenant VMs on a shared hypervisor as the isolation primitive. The platform now needs an orchestration layer that:

* Provisions and lifecycle-manages tenant VMs reproducibly (NFR-01: full rebuild ≤ 1 hr)
* Runs tenants in both home-lab and cloud (TR-07: Cloudflare → GCP path; capability *Business Rules*: may span public and private)
* Supports zero-downtime updates per TR-04 (see ADR-0005 for the mechanism)
* Stays inside the operator's 2 hr/week maintenance budget (NFR-02)

How is the platform's compute layer assembled — what runs the hypervisor, what schedules tenant VMs onto it, and what makes home-lab and cloud look like one fleet from the operator's perspective?

## Decision Drivers

* TR-04: tenant-perceived zero downtime for online updates — needs live VM lifecycle, not just create/destroy
* TR-07: Cloudflare → GCP path is non-negotiable for inter-service traffic
* NFR-01: ≤ 1 hr rebuild — substrate must be definitions-driven, not interactively configured
* NFR-02: ≤ 2 hr/week — substrate must be operable by one person in their spare time
* C-01: spans public and private — substrate must work on both home-lab metal and cloud
* Tiebreaker: reproducibility > vendor independence > minimizing effort

## Considered Options

* **Cloud-managed VMs only (GCE) on GCP, no home-lab compute** — ignore home-lab for compute; use it for storage / boot only.
* **Home-lab KVM + GCE, federated by Terraform only** — two pools, no unified scheduler; operator picks per tenant.
* **Home-lab KVM + GCE, federated by Nomad** — Nomad as the cross-environment scheduler with `qemu` and `gce` task drivers.
* **Home-lab KVM + GCE, federated by Kubernetes + KubeVirt** — Kubernetes as the substrate, VMs as KubeVirt resources.

## Decision Outcome

Chosen option: **home-lab KVM + GCE, federated by Nomad**.

Nomad treats the home-lab and the GCP region as two node pools in one fleet. Tenant VMs are scheduled by declaring a `qemu` job (home-lab) or a `gce` job (cloud); placement is constrained by the tenant's declared resource needs (REQ-01 from the long-form requirements). The operator interacts with one CLI and one job-spec format regardless of where a tenant lands. Definitions live in the infra repo; `nomad job run` is the rebuild verb (NFR-01).

Why not Kubernetes + KubeVirt: KubeVirt makes VMs second-class citizens in a system that is fundamentally container-oriented, and the operator's maintenance burden of running Kubernetes against NFR-02 is hard to justify for a single-operator platform. Why not "GCE-only": violates the spirit of C-01 (the home-lab exists; not using it for compute makes "self-hosted" a label rather than a property). Why not Terraform-only: Terraform is the *provisioner*, not the *scheduler*; it cannot do TR-04's zero-downtime updates without a higher-level orchestrator.

### Consequences

* Good, because home-lab and cloud are one schedulable fleet from the operator's perspective
* Good, because Nomad's job-spec format is small enough to fit in the rebuild budget and the operator's head
* Good, because vendor independence is preserved — Nomad runs anywhere; cloud nodes can be swapped for another provider without re-architecting
* Good, because TR-07 is preserved: Nomad does not impose a network topology, so the existing Cloudflare → GCP → WireGuard → home-lab path is reused unchanged
* Neutral, because Nomad is one more piece of software the operator must learn and maintain
* Bad, because the Nomad ecosystem for VM workloads is smaller than Kubernetes' container ecosystem; some operational tooling will need to be operator-built
* Bad, because if HashiCorp's licensing changes again, the operator may need to switch to OpenBao-style forks (mitigated: Nomad is small enough to fork)

### Confirmation

* The infra repo's definitions include Nomad job specs for each tenant under `cloud/` (or a parallel `homelab/` dir)
* `nomad node status` shows nodes from both home-lab and GCP attached to one cluster
* A scripted standup test (REQ-18 long-form: canary tenant) provisions the canary via `nomad job run` and validates it serves traffic through the Cloudflare path
* The rebuild script (NFR-01) is `terraform apply && nomad job run` — measured end-to-end; if it exceeds 1 hr, ADR-0009-style follow-up is filed

## Pros and Cons of the Options

### Cloud-managed VMs only (GCE) on GCP

* Good, because operator effort is lowest — GCE handles host lifecycle
* Good, because zero-downtime live migration is built-in (helps TR-04)
* Bad, because it forfeits the home-lab as a compute participant — pushes the operator toward "self-hosted in name only"
* Bad, because cost scales with every tenant in cloud-only billing; against NFR-05's spirit

### Home-lab KVM + GCE, federated by Terraform only

* Good, because it uses tools the operator already runs
* Bad, because Terraform is not a scheduler; rolling updates, draining, and zero-downtime reschedules (TR-04) become operator-scripted
* Bad, because there is no live view of which tenant is healthy on which host

### Home-lab KVM + GCE, federated by Nomad

* Good, because one scheduler covers both environments
* Good, because TR-04 has a native primitive (job updates with `max_parallel`, `canary`, `auto_revert`)
* Good, because the surface area fits one operator
* Neutral, because the operator must learn Nomad
* Bad, because VM-driver ecosystem is smaller than container-driver ecosystem
* Bad, because HashiCorp licensing risk (mitigated above)

### Home-lab KVM + GCE, federated by Kubernetes + KubeVirt

* Good, because Kubernetes ecosystem is the largest for orchestration tooling
* Good, because KubeVirt brings real VM lifecycle primitives
* Bad, because running Kubernetes for one operator burns the maintenance budget (NFR-02)
* Bad, because KubeVirt VMs are second-class to containers in the Kubernetes mental model
* Bad, because the substrate's complexity is disproportionate to the platform's scale

## More Information

* The hypervisor is KVM on home-lab nodes and the GCE-managed hypervisor in cloud — the operator does not run a hypervisor in cloud
* WireGuard between home-lab and GCP (existing topology, TR-07) is the network plane Nomad sees as one flat network
* ADR-0005 builds on this: zero-downtime updates use Nomad's update stanza
