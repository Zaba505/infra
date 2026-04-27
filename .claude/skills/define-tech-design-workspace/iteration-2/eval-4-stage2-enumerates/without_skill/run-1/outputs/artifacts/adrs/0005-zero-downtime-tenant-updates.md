---
title: "[0005] Zero-Downtime Tenant Update Strategy"
description: >
    Choose the mechanism by which operator-initiated updates to a tenant's config, version, or capability complete without end-user-visible downtime for online workloads.
type: docs
weight: 5
category: "user-journey"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

[TR-04] requires that operator-initiated tenant updates (config, version, or capability) complete without end-user-visible downtime for tenants serving online traffic. Given ADR-0001 (per-tenant VMs) and ADR-0002 (Nomad as the scheduler), this ADR picks the update strategy.

Note: not every tenant is online; batch / cron-style tenants accept downtime windows. This ADR is scoped to *online* tenants — those whose tech design declares an online-traffic SLA.

## Decision Drivers

* TR-04: zero end-user-visible downtime for online tenants
* ADR-0001: tenant boundary is the VM — updates are VM lifecycle events, not container restarts
* ADR-0002: Nomad is the orchestrator
* NFR-02: maintenance budget — the operator should not have to choreograph each update
* TR-07: traffic flows through Cloudflare; the ingress can drain a backend before the VM goes away

## Considered Options

* **In-place update (mutate the running VM, e.g. config push + service restart)** — fastest, leaves the same VM in place.
* **Rolling replacement: provision new VM, drain old via ingress, swap, retire old (blue/green per tenant)** — two VMs exist briefly, ingress shifts traffic.
* **Live migration of the VM (hypervisor-level state transfer)** — tenant is unaware of the move; underlying host changes.

## Decision Outcome

Chosen option: **rolling replacement with ingress drain (blue/green per tenant)**.

The update sequence is:

1. Nomad provisions a new VM with the updated config / version / capability spec.
2. The new VM passes its health checks.
3. The Cloudflare-side ingress (TR-07) is updated to add the new VM to the backend pool.
4. The old VM is marked draining; in-flight requests complete.
5. Once drained, the old VM is removed from the backend pool and torn down.

Live migration was considered but is the wrong tool: it preserves *the same image* across hosts, which doesn't help when the *image* is what changed. In-place updates can't honor TR-04 for any change that requires a process restart, which is most non-trivial changes. Rolling replacement is the only option that is honest about what an "update" usually is and that uses the ingress as the cutover point — exactly what Cloudflare gives us for free.

### Consequences

* Good, because TR-04 is satisfied for any update that produces a new VM image, including OS patches, runtime version bumps, and capability code changes
* Good, because rollback is "leave the old VM up, drop the new one" — a clean, fast revert
* Good, because the cutover happens at the ingress, not at the application — tenants do not need to implement graceful shutdown beyond the standard practice of finishing in-flight requests
* Good, because Nomad's update stanza (`max_parallel`, `canary`, `auto_revert`) implements this directly
* Neutral, because the tenant briefly runs at 2x VM count during the swap — capacity headroom must exist; this is part of the platform's standing reserve
* Bad, because the old VM cannot be reused for the new spec — every update is a full VM provision; this puts pressure on NFR-01 if updates batch
* Bad, because a misconfigured health check could swap traffic to a broken new VM; mitigated by Nomad's `auto_revert` and by canary checks before full cutover

### Confirmation

* Each online tenant's Nomad job spec declares `update { canary = 1, auto_revert = true, healthy_deadline = … }`
* The standup canary (REQ-18 long-form) exercises a no-op update during validation and asserts zero failed end-user requests during the swap (synthetic load required)
* The operator-initiated-update issue type (REQ-05 long-form) includes a runbook step that checks the Cloudflare backend pool before and after the swap

## Pros and Cons of the Options

### In-place update

* Good, because no second VM is provisioned — fastest, lowest resource cost
* Bad, because most updates require a process restart; restart = downtime unless the tenant implements zero-downtime restarts itself, which violates "no end-user-visible downtime *for the platform's part*"
* Bad, because rollback means re-mutating the VM, which can fail and leave the VM in a partial state

### Rolling replacement (blue/green)

* Good, because it works for any change including OS / runtime / image changes
* Good, because rollback is clean
* Good, because the ingress is the natural cutover point
* Neutral, because it briefly doubles tenant VM count
* Bad, because it depends on health-check correctness — mitigated by canary + auto-revert

### Live migration

* Good, because the tenant is unaware
* Bad, because it preserves the image across hosts — solves the wrong problem (the operator's update is *about* changing the image)
* Good, because it remains useful for *host* maintenance (separate ADR territory) but not for tenant updates
* Bad, because cross-environment live migration (home-lab ↔ GCP) is not practical

## More Information

* The standing capacity reserve required for blue/green is part of the home-lab + GCP capacity planning; it is not a per-tenant cost
* Tenants whose tech design declares "batch/offline" can opt into in-place updates with a downtime window — that path is allowed but not the default
* Updates that change the tenant's contract version (ADR-0003) follow this same swap mechanism, with the new VM running against the new contract adapter
