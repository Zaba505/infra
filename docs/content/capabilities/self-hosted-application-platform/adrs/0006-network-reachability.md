---
title: "[0006] Network Reachability — Cloudflare Edge to Homelab Cluster"
description: >
    External tenant traffic flows end-user → Cloudflare → homelab public IP → router → in-cluster LoadBalancer → ingress controller → tenant namespace. In-cluster networking uses a NetworkPolicy-supporting CNI with default-deny per tenant; inter-tenant traffic is opt-in via tenant-manifest declarations.
type: docs
weight: 6
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform](../_index.md)
**Addresses requirements:** [TR-03](../tech-requirements.md#tr-03-provide-network-reachability--internal-between-tenants-and-external-for-end-users), [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals), [TR-22](../tech-requirements.md#tr-22-tracked-changes-and-immutability-for-all-platform-state-modifying-actions), [TR-26](../tech-requirements.md#tr-26-tenants-declare-resource-needs-at-onboarding-and-on-every-modify-the-platform-admits-or-refuses-based-on-those-declarations), [TR-27](../tech-requirements.md#tr-27-span-public-cloud-and-privatehome-lab-infrastructure-with-the-connectivity-between-them-part-of-the-foundation), [TR-28](../tech-requirements.md#tr-28-no-direct-end-user-access-to-the-platform-itself), [TR-33](../tech-requirements.md#tr-33-routine-platform-operation-must-fit-within-2-hoursweek-of-operator-time)

## Context and Problem Statement

[TR-03](../tech-requirements.md#tr-03-provide-network-reachability--internal-between-tenants-and-external-for-end-users) requires the platform to deliver two flavors of reachability — externally for tenants' end users and internally between tenants — and [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals) requires per-tenant isolation strong enough that no tenant can read another's traffic or signals. [ADR-0001](./0001-public-private-infrastructure-split.md) set up a small cloud edge for ingress and an off-site backup archive, with a self-hosted VPN between cloud and homelab. [ADR-0002](./0002-compute-substrate.md) chose Kubernetes; [ADR-0003](./0003-tenant-packaging-form.md) keeps NetworkPolicy / RBAC entirely on the platform side, derived from the tenant-manifest schema rather than shipped by tenants.

This ADR has to decide the *concrete request path* end-to-end and the *cluster-internal isolation primitive*, then reconcile both with [ADR-0001](./0001-public-private-infrastructure-split.md)'s foundation.

The decisive constraint not visible from the TRs alone is the operator's intended public-internet edge: **Cloudflare**. End-user traffic terminates at Cloudflare, which proxies to the homelab's public IP, where the homelab router forwards into the cluster. This collapses the cloud-edge component of [ADR-0001](./0001-public-private-infrastructure-split.md) for the *request path*: there is no cloud-edge VM in the path. The cloud↔homelab VPN's role narrows accordingly to (a) homelab → cloud egress for backup-archive uploads and (b) any future cloud-side platform reachback to the homelab. The VPN is *not* on tenant request paths.

## Decision Drivers

- **TR-03** — both flavors of reachability must work.
- **TR-32** — per-tenant isolation is the default; cross-tenant traffic is the exception.
- **TR-22 / TR-26 / [ADR-0003](./0003-tenant-packaging-form.md)** — every routing and policy decision must be expressible declaratively from tenant manifests; the platform owns the policy primitives, tenants don't ship them.
- **[ADR-0001](./0001-public-private-infrastructure-split.md) — small cloud edge.** The request path should not require a cloud component beyond what is genuinely public-internet-facing. Cloudflare absorbs that role; the cloud edge from [ADR-0001](./0001-public-private-infrastructure-split.md) shrinks to backup archive only.
- **TR-28** — platform-level surfaces must not be exposed to end users. Ingress is for tenant traffic; admin/observability surfaces are not internet-routable.
- **TR-33 (≤2 hr/week)** — every component on the request path is paid weekly. Fewer is better; mesh-class additions need real justification at this scale.
- **Capability tiebreaker — vendor independence > minimizing operator effort.** Cloudflare is a vendor. The dependency must be bounded enough that leaving it is a config change, not a migration.

## Considered Options

### Option A — Cloudflare proxy to homelab public IP, in-cluster LoadBalancer + ingress controller, NetworkPolicy-default-deny CNI

- **Public edge:** Cloudflare (proxied DNS records for tenant hostnames). TLS terminates at Cloudflare. DDoS / WAF on the free tier.
- **Origin auth:** Authenticated Origin Pulls (mTLS where Cloudflare presents a cert the homelab origin verifies) + Cloudflare-IP allowlist on the router as defense-in-depth.
- **Homelab router:** forwards `:443` from the public IP to the cluster's external LoadBalancer IP. Router configuration is part of the platform definitions ([TR-22](../tech-requirements.md#tr-22-tracked-changes-and-immutability-for-all-platform-state-modifying-actions)).
- **In-cluster LB:** an in-cluster L4 LoadBalancer (e.g. MetalLB-class) assigns a stable IP from a homelab IP pool to the ingress controller's `Service` of type `LoadBalancer`.
- **Ingress controller:** a single cluster-wide Gateway API-class (or standard Ingress) controller, routing by tenant hostname declared in the tenant manifest.
- **CNI:** a NetworkPolicy-supporting CNI; default-deny per tenant namespace; opt-in inter-tenant rules emitted from declared dependencies in the tenant manifest.
- **Cloud↔homelab VPN:** backup-archive egress + any future cloud-side platform reachback only. Not on the tenant request path.
- **TR-32 isolation:** strong — namespaces + default-deny NetworkPolicy + RBAC + (optional later) Pod Security baseline. Mesh-free at this scale.
- **TR-33:** modest — one CNI, one ingress controller, one in-cluster LB; Cloudflare config is mostly DNS + AOP, low ongoing cost.

### Option B — Same as A but add a service mesh (Istio/Linkerd) for in-cluster traffic

- **TR-32:** marginally stronger — mTLS between pods on top of NetworkPolicy.
- **TR-33:** materially higher — mesh control plane, sidecars, upgrade discipline.
- **Justification at personal scale:** weak. A mesh pays for itself when there is a lot of inter-pod traffic with policy/observability/auth needs that NetworkPolicy and the observability ADR can't cover. With single-digit tenants and opt-in inter-tenant traffic, the mesh is cost without proportional return.

### Option C — Cloudflare Tunnel instead of proxied-public-IP

`cloudflared` runs inside the homelab (or as a Deployment in the cluster) and establishes an outbound connection to Cloudflare; no inbound public IP, no router port-forward.

- **Pros:** removes the homelab public-IP and router-port-forward dependency; origin auth handled by the tunnel itself.
- **Cons:** Cloudflare-specific protocol on the homelab side (deeper vendor coupling than DNS+AOP); the homelab now has to run `cloudflared` as a piece of foundation infrastructure; less transparent to debug than a TCP port-forward; failure modes (tunnel restart, version drift) become foundation concerns.
- **Vendor-independence:** weaker than Option A — leaving Cloudflare under Option A is "swap DNS provider, drop AOP, optionally add a local TLS terminator"; under Option C it is also "stop running cloudflared and replace it with whatever inbound-edge replacement, which means re-establishing public reachability from scratch."

### Option D — No Cloudflare; direct DNS A record at homelab public IP, no public proxy

- **Pros:** simplest dependency surface; no third-party.
- **Cons:** the homelab IP is directly exposed; no DDoS or WAF; the operator owns rate-limiting, abuse handling, and TLS issuance entirely; for tenants the operator does *not* yet operate, this exposes whatever bugs are in their software directly to the public internet.
- **Capability tiebreaker:** wins on vendor independence, loses on convenience and resiliency. The capability's tiebreaker explicitly says cost is secondary to *convenience and resiliency*; Cloudflare's free tier buys both.

### Option E — Cloud LB (vendor) at the cloud edge instead of Cloudflare

- **Pros:** more "in scope" of [ADR-0001](./0001-public-private-infrastructure-split.md)'s cloud-edge framing.
- **Cons:** doesn't add what Cloudflare adds (anycast, DDoS, WAF); puts the public edge on a vendor product whose abandonment cost is higher than swapping DNS providers; pays cloud-LB cost for what Cloudflare gives free at this scale.

## Decision Outcome

Chosen option: **Option A — Cloudflare proxy to homelab public IP, in-cluster LoadBalancer + single Gateway-class ingress controller, NetworkPolicy-default-deny CNI; cloud↔homelab VPN is backup-archive-egress and future-reachback only and is not on the request path.**

This option is chosen because:

- It is the only option that simultaneously satisfies [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals) (default-deny NetworkPolicy is sufficient at this scale) and [TR-33](../tech-requirements.md#tr-33-routine-platform-operation-must-fit-within-2-hoursweek-of-operator-time) (no service mesh weekly cost).
- The Cloudflare dependency is bounded: Cloudflare holds no tenant data; leaving Cloudflare is a DNS-record change plus origin-auth removal plus optionally adding a local TLS terminator. That's hours of work, not a tenant migration.
- The capability's *cost is secondary to convenience and resiliency* tiebreaker is honored: Cloudflare's free tier provides DDoS, WAF, and anycast for the cost of a DNS configuration. Option D wins vendor-independence narrowly but loses convenience and resiliency, which is the wrong trade per the capability.
- The narrowed VPN role (backup archive egress + future reachback only) is *more* consistent with [ADR-0001](./0001-public-private-infrastructure-split.md)'s "small cloud edge" stance than my own original [ADR-0001](./0001-public-private-infrastructure-split.md) framing implied. The cloud edge is now genuinely just storage; the VPN exists for that storage path.

Cloudflare is configured with **proxied DNS records** for tenant hostnames pointing at the homelab public IP, with **Authenticated Origin Pulls** as the origin-auth mechanism. Cloudflare Tunnel (Option C) is rejected because its deeper vendor coupling raises the cost of leaving Cloudflare without a corresponding benefit at this scale.

### Consequences

- **Good, because** the request path has exactly four clear hops (Cloudflare → router → in-cluster LB → ingress controller), each with a clear job and a clear failure mode.
- **Good, because** TLS termination, certificate management, DDoS, and WAF live at Cloudflare and are paid for in DNS configuration, not in operator weekly time.
- **Good, because** [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals) is enforced by Kubernetes-native primitives ([ADR-0002](./0002-compute-substrate.md))'s namespace + NetworkPolicy + RBAC. The translator from [ADR-0003](./0003-tenant-packaging-form.md) emits a default-deny NetworkPolicy per tenant namespace plus explicit allow rules for declared dependencies; no mesh sidecar to operate.
- **Good, because** the cloud↔homelab VPN is now justified by exactly one continuous use (archive egress) plus latent reachback, which is a load-bearing role rather than a vague foundation concern.
- **Good, because** image-registry traffic stays inside the cluster (or homelab LAN), so pod startup does not depend on the cross-boundary VPN.
- **Bad, because** the homelab public IP is now part of the foundation. Loss of that IP (ISP change, dynamic-IP churn) is a real event needing a rebuild step (DNS update at Cloudflare). [TR-22](../tech-requirements.md#tr-22-tracked-changes-and-immutability-for-all-platform-state-modifying-actions) requires this to be tracked rather than hand-rolled.
- **Bad, because** the homelab router is now load-bearing infrastructure. Its `:443` port-forward and its Cloudflare-IP allowlist must be in the platform definitions, not on a sticky note.
- **Bad, because** the platform now depends on Cloudflare. Bounded — but real.
- **Bad, because** [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals) at this layer relies on NetworkPolicy correctness. Misconfiguration is a real risk; the canary tenant ([TR-20](../tech-requirements.md#tr-20-maintain-a-purpose-built-canary-tenant-alongside-the-platforms-definitions-and-use-it-as-the-readiness-signal)) should exercise a "should be denied" path so isolation regressions are caught at rebuild time.
- **Requires:**
  - **[ADR-0001](./0001-public-private-infrastructure-split.md) narrowed in tech-design.md.** When Stage 3 composes the design, the cloud↔homelab VPN should be described as backup-archive egress + future reachback only. The original [ADR-0001](./0001-public-private-infrastructure-split.md) text is consistent with this; tech-design.md should not overstate the VPN's role.
  - **[ADR-0003](./0003-tenant-packaging-form.md) schema** gains a tenant hostname field, declared external ports, and an optional inter-tenant-dependencies block (which tenants this tenant calls and on which ports). The translator emits Cloudflare DNS records, ingress routes, and inter-tenant NetworkPolicy from these.
  - **[ADR-0007](./0007-...) (backup & DR)** uses the VPN for archive uploads and is the VPN's primary continuous justification.
  - **[ADR-0008](./0008-...) (observability) and [ADR-0009](./0009-...) (secrets)** must reach into clusters without making admin surfaces internet-routable per [TR-28](../tech-requirements.md#tr-28-no-direct-end-user-access-to-the-platform-itself); operator access to these surfaces flows over the VPN or via local-only access, not through Cloudflare.
  - **[ADR-0012](./0012-...) (definitions tooling)** must drive Cloudflare DNS records, the homelab router config, the cluster CNI, the ingress controller, and the in-cluster LB — i.e. the tooling reaches all five hops.
  - **[ADR-0015](./0015-...) (canary tenant)** exercises the full external request path *and* an inter-tenant denied path.

### Realization

How this decision shows up in the repo:

- **Cloudflare configuration** lives in the platform definitions: proxied A/AAAA records per tenant hostname, Authenticated Origin Pulls enabled, Cloudflare-IP allowlist for the origin documented for the router config.
- **Homelab router configuration** lives in the platform definitions (as configuration the operator applies during Phase 1 of the rebuild — see [stand-up-the-platform §3](../user-experiences/stand-up-the-platform.md#3-phase-1--foundations)). At minimum: `:443` forward to the cluster external LoadBalancer IP, Cloudflare-IP allowlist, the homelab side of the cloud↔homelab VPN.
- **`homelab/`** (or the equivalent placed by [ADR-0012](./0012-...)) provisions the cluster CNI with NetworkPolicy enabled, the in-cluster L4 LB controller (MetalLB-class) configured with the homelab IP pool, and a single cluster-wide Gateway-class ingress controller installed in its own namespace.
- **Per-tenant namespace bootstrap** (emitted by [ADR-0003](./0003-tenant-packaging-form.md)'s translator) includes:
  - A default-deny NetworkPolicy.
  - Allow rules: ingress controller → tenant pods on declared external ports; tenant pods → cluster DNS; tenant pods → declared inter-tenant targets only.
  - An ingress route binding the tenant hostname to the tenant Service.
- **Image registry** runs in-cluster on the homelab side in its own namespace, reachable from all tenant namespaces (image pulls are an explicit allow in the per-tenant NetworkPolicy template). Specific registry product is deferred.
- **Canary tenant** ([TR-20](../tech-requirements.md#tr-20-maintain-a-purpose-built-canary-tenant-alongside-the-platforms-definitions-and-use-it-as-the-readiness-signal)) exercises: an external request through Cloudflare → router → cluster → namespace; an outbound-allowed call (DNS + cluster identity); and a known-denied inter-tenant call (asserting NetworkPolicy default-deny is in effect).

## Open Questions

- **Specific CNI** (Cilium, Calico, …). Pin "K8s CNI with NetworkPolicy default-deny." Cilium is the leading candidate because it also gives substrate-level observability hooks ([ADR-0008](./0008-...)); decide alongside that ADR or at deployment.
- **Specific in-cluster LB and ingress controller products.** Pin "L4 LB assigning stable IPs from a homelab pool" and "single cluster-wide Gateway-class controller." Products deferred.
- **Specific image registry product.** Pin "operator-controlled, in-cluster, reachable from tenant pods over allow-listed NetworkPolicy."
- **VPN terminator placement.** Either the homelab router runs the VPN terminator, or a dedicated host on the homelab does. Both are consistent with this ADR; deferred.
- **Pod Security baseline level** (restricted vs. baseline). Hardening choice that interacts with this ADR's isolation story; finalize alongside [ADR-0009](./0009-...) (secrets).
- **Dynamic-IP at the homelab.** If the ISP gives a dynamic public IP, automated DNS updates at Cloudflare are required. Operationally workable; the specific mechanism is deferred.
