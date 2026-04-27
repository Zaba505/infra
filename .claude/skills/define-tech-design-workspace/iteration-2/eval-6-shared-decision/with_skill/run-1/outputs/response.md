# Stopping before I draft this ADR — it's a shared decision, not a capability-scoped one

We're in **Stage 2** of `define-tech-design` for `self-hosted-application-platform`. `tech-requirements.md` is reviewed (`reviewed_at: 2026-04-26`), so the Stage-1 gate is open. But the decision you've described shouldn't land as a capability-scoped ADR, and the skill is explicit about that. Here's why, and what I'd suggest doing instead.

## Why this is out of scope for a capability-scoped ADR

The `define-tech-design` skill has a flag-and-stop rule:

> "If a decision is obviously cross-capability (touches Cloudflare topology, identity, networking, error response format, the resource identifier standard, etc.), do not draft it as a capability-scoped ADR. Surface it."

"How all services authenticate to each other" hits **three** of those triggers at once:

1. **Cloudflare topology.** Terminating mTLS at the Cloudflare edge is a property of the edge layer that fronts every service in the monorepo, not just platform services.
2. **Identity.** Service identity (what cert each service presents, how its identity is established and rotated) is a foundation other capabilities will inherit.
3. **Networking.** This shapes the trust boundary for every inter-service call repo-wide.

It's also forced by precedent in the requirements doc itself: **TR-07** already records "all inter-service communication must traverse the Cloudflare → GCP path" as an *inherited shared decision*, sourced to `CLAUDE.md`. A decision about how authentication rides on that path is the same kind of artifact: shared, repo-wide, owned at `docs/content/r&d/adrs/`. If I drafted it under `capabilities/self-hosted-application-platform/adrs/`, I'd be putting the canonical answer in the wrong place — and the next capability that needs to call a service would either re-decide it (drift) or cite a capability-scoped ADR they don't own (worse).

`define-tech-design` explicitly carves this out: "Cross-capability decisions are **out of scope**. If a decision is shared, flag it and stop — it belongs in `docs/content/r&d/adrs/` via a separate flow."

So I'm stopping rather than writing the ADR you asked for. Please don't read that as obstructionism — it's the gate doing its job.

## A second concern worth surfacing before you draft anywhere

Before this gets drafted as a shared ADR, the framing is worth pressure-testing. **mTLS terminated at the Cloudflare edge does not authenticate services to each other.** It authenticates:

- the *caller* to Cloudflare (client cert at the edge), and
- *Cloudflare* to the origin (Authenticated Origin Pulls or similar).

Service A calling Service B through Cloudflare gets you "Cloudflare vouches that some authorized client made this call" — not "Service B knows Service A made this call." If your goal really is service-to-service identity (per-service certs, per-service authorization), termination needs to happen *at the services* (or at a sidecar), not at the edge. If your goal is "only authorized clients can reach our origins", edge mTLS is the right tool but it's a perimeter control, not a service-to-service auth scheme.

Worth pinning down which one you mean before any ADR — they lead to materially different designs (and to different rotation stories; 30 days is fine for either, but the *what* being rotated differs).

## What you can do next

Pick one:

1. **Defer and continue Stage 2 here.** If there are any *capability-local* decisions for self-hosted-application-platform you want to tackle (TR-01 isolation mechanism, TR-02 contract versioning, TR-04 zero-downtime update mechanism, TR-05 export mechanism, TR-06 idempotent import), I'll enumerate them and we'll work one ADR at a time.
2. **Pause Stage 2 and open the shared-ADR conversation.** Tell me to switch to drafting a shared ADR under `docs/content/r&d/adrs/` for inter-service auth. That's a separate flow — I'd want to clarify the perimeter-vs-service-identity question above first, then enumerate options (edge mTLS only, mTLS at services with SPIFFE/SPIRE, mTLS via service mesh, signed JWT bearer tokens over TLS, etc.) rather than ratifying one upfront.
3. **Both** — defer the shared decision, proceed with capability-local ADRs now, come back to the shared one later. The capability tech design can cite "[shared ADR pending]" against TR-01/TR-03 in the meantime; Stage 3 won't compose until that's resolved.

Which would you like?
