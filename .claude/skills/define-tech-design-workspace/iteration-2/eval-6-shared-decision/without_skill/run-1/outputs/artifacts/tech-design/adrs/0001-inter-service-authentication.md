---
title: "[0001] Inter-Service Authentication for the Self-Hosted Application Platform"
description: >
    How platform services and tenant workloads authenticate to each other across the Cloudflare → GCP → home-lab path.
type: docs
weight: 1
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: []
consulted: []
informed: []
---

## Context and Problem Statement

The Self-Hosted Application Platform spans Cloudflare (edge), GCP (cloud-hosted services), and the home lab (over WireGuard). Per [TR-07](../technical-requirements.md#tr-07-all-inter-service-communication-must-traverse-the-cloudflare--gcp-path) all inter-service traffic traverses this path, and per [TR-01](../technical-requirements.md#tr-01-tenants-must-be-isolated-such-that-no-tenant-can-read-anothers-state) we must enforce strict tenant isolation including at the network layer.

Today there is no defined mechanism by which platform services authenticate one another, nor by which a tenant workload authenticates to a platform service (or vice versa). We need a single, capability-wide answer to: **how does service A prove its identity to service B before B trusts the request?**

The user proposed: *use mTLS terminated at the Cloudflare edge, with certificate rotation every 30 days.*

This ADR records that proposal, and — importantly — flags a structural concern with it before accepting it as the decision.

## Decision Drivers

* **TR-01 (tenant isolation):** authentication must be strong enough that a tenant cannot impersonate the platform or another tenant.
* **TR-07 (Cloudflare → GCP topology):** the chosen mechanism must work with the existing edge-fronted topology without requiring tenants or operators to bypass Cloudflare.
* **Operational burden:** rotation, issuance, and revocation must be automatable. Manual cert handling at home-lab scale is acceptable only if rare.
* **Reuse of existing primitives:** the repo already has [`cloud/mtls/cloudflare-gcp/`](../../../../../../cloud/mtls/cloudflare-gcp/) — Cloudflare-to-GCP mTLS trust anchors. Whatever we choose should compose with that, not duplicate it.
* **Defense in depth:** edge authentication alone does not authenticate east-west traffic between services *behind* the edge.

## Considered Options

* **Option A — mTLS terminated at the Cloudflare edge, 30-day rotation** (the proposal as stated).
* **Option B — mTLS at the edge *plus* mTLS (or workload identity) between services behind the edge** (the proposal, extended).
* **Option C — Token-based service identity (e.g. SPIFFE/SPIRE or OIDC workload identity) end-to-end**, with TLS used only for transport encryption.
* **Option D — Defer the decision** until Q-04 (identity service) and Q-02 (compute substrate) ADRs are written, since both materially constrain the answer.

## Decision Outcome

**Chosen option: Option D — defer, but with a recorded preference for Option B over Option A.**

The proposal as stated (Option A) does not, on its own, satisfy the problem statement. mTLS *terminated* at the Cloudflare edge authenticates Cloudflare to the origin (and optionally the origin to Cloudflare), but it does **not** authenticate one platform service to another, nor does it authenticate a tenant workload to a platform service once traffic is past the edge. Accepting Option A as written would leave east-west authentication undefined and would not satisfy TR-01 under a "tenant workload running inside the platform" threat model.

Two things should happen before this ADR moves from `proposed` to `accepted`:

1. **Resolve Q-04 (identity service)** in a sibling ADR. The east-west answer (workload identity, SPIFFE, mesh-issued certs, or OIDC tokens) is part of that decision; it is not appropriate to pin it here.
2. **Resolve Q-02 (compute substrate)** in a sibling ADR. The substrate (Kubernetes, Nomad, plain VMs, Cloud Run, etc.) determines which workload-identity primitives are even available and how cert issuance/rotation is automated.

In the interim, the recorded preference is:

* **Edge (north-south):** keep the existing Cloudflare ↔ GCP mTLS pattern from `cloud/mtls/cloudflare-gcp/`. Rotation cadence is a property of that module and should be set there, not asserted here.
* **Behind the edge (east-west):** plan for workload-identity-issued mTLS or signed tokens, to be specified by the Q-04 ADR.
* **30-day rotation:** noted as a target, not adopted. The right cadence depends on the issuer (e.g. SPIRE defaults are hours, ACME-issued certs commonly 90 days, internal CAs vary). Pinning 30 days before choosing the issuer is premature.

### Consequences

* Good, because we avoid encoding a decision that doesn't actually solve the stated problem (authenticating services *to each other*).
* Good, because we keep the door open for stronger east-west authentication once the substrate and identity decisions land.
* Good, because we explicitly link this ADR to the open questions (Q-02, Q-04) that block it.
* Neutral, because no new mechanism is adopted today — current behavior (edge mTLS only, nothing east-west) persists until the blockers resolve.
* Bad, because the platform remains without a documented east-west authentication story until Q-02 and Q-04 are decided. This should be tracked.

### Confirmation

This ADR will be revisited and either superseded or promoted to `accepted` once:

1. The Q-02 (compute substrate) ADR is accepted.
2. The Q-04 (identity service) ADR is accepted.
3. A follow-up review confirms the east-west mechanism satisfies TR-01 against the "malicious tenant workload" threat model.

## Pros and Cons of the Options

### Option A — mTLS terminated at the Cloudflare edge, 30-day rotation

* Good, because it reuses the existing `cloud/mtls/cloudflare-gcp/` module.
* Good, because Cloudflare handles edge cert lifecycle, lowering operational burden at the perimeter.
* Good, because 30-day rotation is a defensible cadence for edge-issued certs (short enough to limit blast radius, long enough to avoid churn).
* Bad, because "terminated at the edge" means the authentication context is lost beyond Cloudflare — service B cannot cryptographically verify that a request truly came from service A.
* Bad, because it does not address tenant-to-platform or tenant-to-tenant authentication once traffic is inside the trust boundary.
* Bad, because it pins a rotation cadence (30 days) before the issuer is chosen.
* Bad, because under TR-01 a compromised or hostile tenant workload inside the platform would face no cryptographic barrier to impersonating a platform service.

### Option B — Edge mTLS *plus* east-west workload-identity mTLS

* Good, because it preserves the edge benefits of Option A and adds a meaningful answer to the actual question.
* Good, because it composes with substrate-native identity (e.g. Kubernetes ServiceAccount-bound certs, SPIFFE SVIDs).
* Good, because rotation can be aggressive (hours) for east-west certs without operator burden, since issuance is automated by the identity plane.
* Neutral, because it requires choosing an issuer/identity plane — i.e., it depends on Q-04.
* Bad, because it is more components to operate than Option A.

### Option C — Token-based service identity (SPIFFE JWT-SVID, OIDC) end-to-end

* Good, because tokens carry verifiable identity through L7 proxies (including Cloudflare), avoiding the "termination loses context" problem.
* Good, because tokens compose naturally with per-request authorization.
* Neutral, because TLS is still needed for transport confidentiality — tokens don't replace it.
* Bad, because token validation requires a verifier (JWKS endpoint, etc.) reachable by every service, increasing the platform's internal surface.
* Bad, because token theft (e.g. via SSRF in a tenant workload) is a meaningful risk unless tokens are short-lived and audience-bound.

### Option D — Defer until Q-02 and Q-04 are decided

* Good, because it avoids locking in a mechanism that the substrate or identity decisions might invalidate.
* Good, because it makes the dependency between these ADRs explicit and traceable.
* Neutral, because the deferral is bounded — the blocking ADRs are already enumerated as open questions.
* Bad, because the platform has no documented east-west authentication story in the meantime.

## More Information

### Why "mTLS terminated at the edge" is not, by itself, an answer to the question asked

The user's question was: *how do all services authenticate to each other?* mTLS terminated at the Cloudflare edge answers a different question: *how does Cloudflare authenticate to my origin (and optionally vice versa)?* Once the TLS session terminates at Cloudflare, the origin sees a request from Cloudflare — not from the original caller. To authenticate caller-to-callee for two services that both sit behind the edge, you need either:

* a separate TLS session established between caller and callee (east-west mTLS), or
* a bearer token / signed assertion that the callee can verify independently of the transport (Option C).

This ADR exists in part to make that distinction explicit so a future reader does not assume "edge mTLS" suffices for TR-01.

### Open questions this ADR depends on

* **Q-02:** compute substrate — determines which workload-identity primitives are available.
* **Q-04:** identity service — determines the issuer, rotation cadence, and verification path for east-west credentials.

### References

* `cloud/mtls/cloudflare-gcp/` — existing edge mTLS module.
* [TR-01](../technical-requirements.md#tr-01-tenants-must-be-isolated-such-that-no-tenant-can-read-anothers-state), [TR-07](../technical-requirements.md#tr-07-all-inter-service-communication-must-traverse-the-cloudflare--gcp-path).
* SPIFFE/SPIRE: https://spiffe.io/
* RFC 8705 (OAuth 2.0 Mutual-TLS Client Authentication).
