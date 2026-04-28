# Response

Stopping before I draft this one — it's a shared decision, not a capability-scoped one.

How platform services authenticate to each other is cross-capability by definition: every capability that runs a service on this platform inherits the answer, and the proposal touches Cloudflare edge topology (mTLS termination at the edge) plus a global rotation cadence. Per the `define-adr` skill's flag-and-stop rule, decisions that touch Cloudflare topology, identity, or networking belong in `docs/content/r&d/adrs/` via the shared-ADR flow, not under `docs/content/capabilities/self-hosted-application-platform/adrs/`. If I wrote it as a capability-scoped ADR, every other capability would either silently inherit it or contradict it, and we'd end up deleting and re-homing it later.

A couple of additional concerns worth surfacing before any ADR is drafted, shared or otherwise:

1. **"mTLS terminated at the Cloudflare edge" doesn't actually authenticate platform services to each other.** Cloudflare-edge mTLS authenticates *external clients* to the platform's ingress. Service-to-service auth happens *inside* the trust boundary (home lab ↔ GCP over WireGuard, per CLAUDE.md). So the proposed mechanism likely doesn't solve the stated problem; we should clarify whether we mean (a) external client → platform mTLS (already partly covered by `cloud/mtls/cloudflare-gcp/`), or (b) service ↔ service auth inside the mesh, which is a different decision with a different option set (mTLS via SPIFFE/cert-manager, signed JWTs, GCP IAM/IAP, WireGuard-only with no app-layer auth, etc.).
2. **30-day rotation is a concrete parameter, not a decision.** It belongs in the chosen option's realization, after the mechanism is picked. Locking it in before the mechanism risks committing to a cadence that doesn't match the tooling (e.g. cert-manager defaults, Cloudflare Origin CA limits).
3. **I have not opened `tech-requirements.md`.** Even though the gate is satisfied, the right next step here is to reclassify this as a shared ADR and then, in the shared flow, identify which TRs across which capabilities it actually addresses. I'm deliberately not citing TR-NNs from the self-hosted-application-platform doc, because pinning the decision to one capability's TRs is what makes it look capability-scoped when it isn't.

**Recommended next steps:**

- Reclassify (or close) the capability-scoped issue, if `plan-adrs` filed one.
- Open a shared-ADR flow under `docs/content/r&d/adrs/` for "platform service-to-service authentication" with the scoping question above resolved first.
- If you want, I can draft the issue text for the shared ADR — but I won't write the ADR file at the capability path.

No ADR file was written.
