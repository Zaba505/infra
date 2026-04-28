---
title: "[0008] Platform Inter-Service Authentication via Cloudflare-Terminated mTLS"
description: >
    All self-hosted-application-platform services authenticate to each other using mTLS terminated at the Cloudflare edge, with client certificates rotated every 30 days.
type: docs
weight: 8
category: "strategic"
status: "proposed"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

## Context and Problem Statement

The self-hosted-application-platform capability is composed of multiple internal services (compute scheduler, storage broker, identity, observability, backup, etc.) that must call each other to deliver tenant-facing capabilities. Per [TR-07](../../capabilities/self-hosted-application-platform/tech-requirements.md#tr-07-all-inter-service-communication-must-traverse-the-cloudflare--gcp-path), all inter-service traffic traverses the Cloudflare → GCP path with WireGuard back to the home lab. We need a single, consistent way for those services to prove their identity to each other so that:

* Tenant isolation ([TR-01](../../capabilities/self-hosted-application-platform/tech-requirements.md#tr-01-tenants-must-be-isolated-such-that-no-tenant-can-read-anothers-state)) cannot be circumvented by a misrouted or spoofed call between platform components.
* The operator does not have to invent a per-service auth scheme each time a new platform service is added.

How should platform services authenticate to each other?

## Decision Drivers

* **TR-01** — tenant isolation depends on platform components trusting only legitimate peers.
* **TR-07** — all inter-service traffic already traverses Cloudflare; the auth solution must fit that topology.
* **Operator maintenance budget** (capability KPI: ≤ 2 hours/week) — credential management cannot become a recurring toil sink.
* **Reproducibility** (capability KPI: stand up in ≤ 1 hour) — the auth mechanism must be expressible as definitions, not manual cert handling.
* **Vendor independence** (capability outcome) — Cloudflare is already in the trust path; no *new* vendor lock-in is acceptable.
* Avoid bespoke auth code in each Go service; reuse what the existing `cloud/mtls/cloudflare-gcp/` Terraform module already establishes.

## Considered Options

* **mTLS terminated at the Cloudflare edge, certs rotated every 30 days** (chosen)
* mTLS terminated at each service (peer-to-peer mTLS inside the platform)
* Bearer tokens (JWT) issued by an internal auth service
* Network-level trust only (rely on WireGuard + private GCP networking, no per-call auth)

## Decision Outcome

Chosen option: **"mTLS terminated at the Cloudflare edge, certs rotated every 30 days"**, because Cloudflare is already the single ingress for all inter-service traffic per TR-07, the existing `cloud/mtls/cloudflare-gcp/` module already provisions the trust anchors, and 30-day rotation matches a cadence the operator can automate without exceeding the maintenance budget. Each platform service is issued a client certificate; Cloudflare validates the cert at the edge and forwards an authenticated identity header to the destination service over the existing GCP path.

### Consequences

* Good, because authentication is centralized at the edge — individual Go services do not implement TLS client auth themselves and stay focused on their domain logic.
* Good, because it composes with existing Terraform (`cloud/mtls/cloudflare-gcp/`) and follows the same pattern Cloudflare already uses for tenant ingress.
* Good, because 30-day rotation limits the blast radius of a leaked cert without requiring constant operator attention if rotation is automated.
* Good, because revocation is operationally simple: pull the cert from Cloudflare's allowed-CA bundle and the caller is locked out at the edge.
* Neutral, because the auth boundary is the Cloudflare edge — services on the GCP side trust the forwarded identity header, so the WireGuard + GCP private-network segment must remain non-tenant-reachable (already enforced by TR-07).
* Bad, because Cloudflare becomes a hard dependency for *internal* auth, not just ingress — a Cloudflare outage degrades intra-platform calls, not only external traffic.
* Bad, because rotation must be automated end-to-end (issue → distribute → upload to Cloudflare → reload caller); a half-built rotation pipeline will burn the maintenance budget.
* Bad, because debugging an auth failure now spans Cloudflare config + GCP load balancer + service logs rather than living in one place.

### Confirmation

Implementation compliance will be confirmed through:
1. Terraform in `cloud/mtls/cloudflare-gcp/` (or a sibling module) provisions per-service client certs and the Cloudflare-side trust bundle.
2. An automated rotation job runs on a ≤ 30-day cadence and is observable (last-rotation timestamp surfaced to the operator).
3. Integration test: a request to a platform service bearing no client cert, or a cert outside the trust bundle, is rejected at the Cloudflare edge before reaching GCP.
4. Code review: no platform service implements its own bearer-token or shared-secret auth for inter-service calls.

## Pros and Cons of the Options

### mTLS terminated at the Cloudflare edge, certs rotated every 30 days

Each platform service receives a client certificate from a CA Cloudflare trusts. Cloudflare validates the cert at the edge and forwards the authenticated identity to the destination service over the existing Cloudflare → GCP → WireGuard path. Certs are rotated every 30 days via automation.

* Good, because it reuses the trust path already mandated by TR-07.
* Good, because it offloads TLS client-auth complexity from each Go service.
* Good, because the existing `cloud/mtls/cloudflare-gcp/` module provides the trust-anchor pattern.
* Good, because 30 days is short enough to limit leaked-cert exposure but long enough that rotation failures don't cause same-week outages.
* Neutral, because services must trust an identity header forwarded from Cloudflare; this is acceptable given network segmentation but worth documenting.
* Bad, because Cloudflare becomes the single point of failure for inter-service auth.
* Bad, because rotation automation is non-trivial and must be built before the cadence is meaningful.

### mTLS terminated at each service (peer-to-peer mTLS inside the platform)

Each Go service performs TLS client-auth itself; certs distributed by an internal CA.

* Good, because there is no dependency on Cloudflare for *internal* auth — Cloudflare outage does not break intra-platform calls.
* Good, because the auth boundary is closer to the service, reducing reliance on a forwarded identity header.
* Bad, because every service must implement and maintain TLS client-auth, conflicting with the small-Go-services pattern documented in CLAUDE.md.
* Bad, because cert distribution to each service requires its own automation (Secret Manager pulls, reload-on-rotation, etc.) which is heavier than the centralized Cloudflare-side approach.
* Bad, because it duplicates a trust-anchor mechanism that already exists at the edge.

### Bearer tokens (JWT) issued by an internal auth service

A central auth service mints short-lived JWTs that callers attach to every inter-service request.

* Good, because tokens are short-lived by design and revocation is trivial (stop minting).
* Good, because identity is carried in-band and visible to the destination service without a header-trust assumption.
* Bad, because it introduces a new platform service (the token issuer) that itself needs to authenticate callers — the problem recurses.
* Bad, because every Go service must implement JWT verification (key rotation, clock skew, audience checks) — bespoke auth code in each service.
* Bad, because it does not reuse the Cloudflare/mTLS trust path that TR-07 already mandates.

### Network-level trust only (rely on WireGuard + private GCP networking)

No per-call auth; trust is established by network reachability alone.

* Good, because it is the simplest possible option — zero auth code, zero rotation.
* Good, because it relies entirely on infrastructure already in place.
* Bad, because any compromise of a single platform component (or a misconfigured firewall rule) yields full lateral movement across all platform services.
* Bad, because it gives no defense in depth for TR-01 — tenant isolation would rest on the network alone.
* Bad, because adding a new service to the trusted network is a manual, error-prone step.

## More Information

### Trust path

```
caller-service  →  Cloudflare edge (mTLS validation)  →  GCP HTTPS LB  →  destination-service
```

Cloudflare validates the client certificate against the per-service trust bundle and forwards an identity header (e.g. `Cf-Client-Cert-Subject`) to the destination service. Destination services treat this header as authoritative *only* for traffic arriving over the WireGuard/GCP private path; direct ingress paths must reject it.

### Rotation cadence

30 days. Rationale: long enough to avoid same-week breakage from a transient rotation failure (operator has time to react within the maintenance budget), short enough that a leaked cert has bounded value.

### Related

* TR-01 — Tenant isolation
* TR-07 — Cloudflare → GCP traffic path
* Existing module: `cloud/mtls/cloudflare-gcp/`
* ADR 0001 — Use MADR for ADRs (format)
