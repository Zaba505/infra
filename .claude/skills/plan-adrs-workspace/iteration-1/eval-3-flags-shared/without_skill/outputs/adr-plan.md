# ADR Plan: Inter-Service Authentication for self-hosted-application-platform

## Request
Plan an ADR for how all platform services authenticate to each other, with the proposed approach: **mTLS terminated at the Cloudflare edge, certificates rotated every 30 days.**

## Concerns flagged before filing

Two things should be resolved with the user before filing an ADR issue:

### 1. This decision looks cross-capability, not capability-scoped

The user framed this as how *all platform services* authenticate to *each other*. That is a foundational, repo-wide concern — not something specific to the `self-hosted-application-platform` capability. Evidence:

- `TR-07` in `tech-requirements.md` is sourced as a **"prior shared decision"** (not from this capability or any of its UXs), meaning networking topology is already treated as a shared decision living outside this capability.
- The same auth mechanism would apply to any other capability's services in this monorepo.
- Per repo conventions (`docs/content/r&d/adrs/`), shared cross-capability ADRs live there, while capability-scoped ADRs live under the capability folder.

**Recommendation:** file this ADR under `docs/content/r&d/adrs/` as a shared ADR, not under `docs/content/capabilities/self-hosted-application-platform/adrs/`. Confirm with the user before proceeding.

### 2. The proposed mechanism may not actually solve service-to-service auth

"mTLS terminated at the Cloudflare edge" authenticates **external clients to Cloudflare**, not **platform services to each other**. Once TLS is terminated at the edge, traffic forwarded from Cloudflare to origin (and any service-to-service traffic *inside* the home lab / GCP perimeter) is no longer covered by that mTLS handshake. The current topology is `Internet → Cloudflare (mTLS proxy) → Home Lab ↔ GCP (Wireguard)` — edge mTLS sits at the first hop only.

So the proposed solution as stated does not satisfy the stated goal. Likely options the user may actually mean:

- **(A)** Edge mTLS only (Cloudflare ↔ origin), and rely on WireGuard + network isolation for service-to-service traffic. This is *not* service-to-service auth, but may be acceptable.
- **(B)** End-to-end mTLS: every service presents and verifies a client cert to every other service (e.g., SPIFFE/SPIRE, cert-manager-issued workload certs, or a service mesh). Cloudflare edge mTLS is then a separate, additional layer.
- **(C)** Edge mTLS + a different intra-cluster auth (e.g., signed JWTs, OIDC workload identity, GCP service-account auth).

**Recommendation:** clarify with the user which of (A)/(B)/(C) they want before drafting the ADR. The 30-day rotation cadence applies meaningfully only if there's an issuing authority for service certs (i.e., (B)).

## Decisions to enumerate (pending clarification)

If the user confirms intent and scoping, the ADR(s) to file would be:

1. **ADR (shared): Service-to-service authentication mechanism** — picks among edge-only mTLS, full mesh mTLS (SPIFFE/cert-manager/Linkerd/Istio), or token-based auth. Drives TR-01 (tenant isolation at the network layer) and TR-07 (inter-service traffic on the Cloudflare→GCP path).
2. **ADR (shared): Certificate issuance and rotation policy** — issuing CA, rotation cadence (proposed 30 days), revocation, and bootstrap. Could be folded into #1 if the chosen mechanism implies the CA design.

Both belong under `docs/content/r&d/adrs/`, not under the capability.

## Proposed next step

Before filing any GitHub issue, get the user's confirmation on:
- (a) the ADR is shared and goes under `docs/content/r&d/adrs/`, and
- (b) which auth model — edge-only, end-to-end mTLS, or hybrid — they actually want decided.

Once confirmed, file one issue per ADR using the planned `gh issue create` invocation in `gh-invocations.txt`.
