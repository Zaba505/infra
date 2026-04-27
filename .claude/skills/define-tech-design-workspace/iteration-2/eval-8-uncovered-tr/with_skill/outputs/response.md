# Cannot proceed to Stage 3 — TR-08 is uncovered

I read the three artifacts in `docs/content/capabilities/self-hosted-application-platform/`:

- `tech-requirements.md` — TR-01 through TR-08, `reviewed_at: 2026-04-26` (gate satisfied).
- `adrs/0001-tenant-state-storage.md` — `accepted`, addresses TR-01, TR-04.
- `adrs/0002-contract-versioning.md` — `accepted`, addresses TR-02.
- `adrs/0003-tenant-eviction-export.md` — `accepted`, addresses TR-05.

Stage 2's accept-gate is satisfied (every ADR is `accepted`). However, Stage 3 has its own precondition embedded in the required "Requirement → ADR → realization map":

> Every TR must appear; if any TR has no ADR, Stage 3 is premature — return to Stage 2.

Mapping TRs to ADRs surfaces a gap:

| TR | Addressed by |
|----|---|
| TR-01 | ADR-0001 |
| TR-02 | ADR-0002 |
| TR-03 | — (no ADR) |
| TR-04 | ADR-0001 |
| TR-05 | ADR-0003 |
| TR-06 | — (no ADR) |
| TR-07 | — (prior shared decision; constraint, not a capability ADR) |
| TR-08 | **— (no ADR)** |

You asked me to compose `tech-design.md`. **I'm not going to**, because doing so would force me to either (a) silently omit TR-08 from the audit-trail map, or (b) paper over it with narrative prose. Both failure modes are exactly what the skill's Stage 3 discipline forbids — see the skill's Example 3, which describes this scenario verbatim using TR-08 as the example:

> "Composing the design surfaced a gap: TR-08 (graceful degradation when GCP region is unreachable) is in `tech-requirements.md` but no ADR addresses it. I'm not going to invent a story in `tech-design.md`. Returning to Stage 2 to draft an ADR for it — should I propose options?"

**Returning to Stage 2.** TR-08 requires a decision: how does the platform degrade gracefully when a GCP region is unreachable for up to 30 minutes? Candidate options I can flesh out into an ADR on your confirmation:

1. **Multi-region active/active Cloud Run with GCLB failover.** Tenant workloads pinned to N regions; HTTPS LB sheds traffic from an unreachable region automatically. Highest availability, highest cost, requires per-tenant Firestore namespace replication strategy (interacts with ADR-0001).
2. **Active/passive with warm standby in a second region + manual or automated promotion.** Lower steady-state cost, RTO measured in minutes, requires a promotion runbook and Firestore export/import or scheduled replication.
3. **Read-only degraded mode in surviving region.** When primary region is unreachable, the platform serves cached/read-only responses from a secondary region with a tenant-visible banner. Cheapest, but "reduced functionality" is exposed to end users — needs sign-off from the capability owner.

There are also two adjacent gaps worth flagging while we're back in Stage 2 (you can defer them, but I'd be remiss not to surface them):

- **TR-03 (per-tenant observability)** has no ADR.
- **TR-06 (lossless, idempotent migrations)** has no ADR.

**TR-07** (Cloudflare → GCP traffic) traces to a prior shared decision under `docs/content/r&d/adrs/`; it's a constraint inherited into the design, not a new capability-scoped ADR. I'll cite it in the eventual `tech-design.md` rather than re-decide it.

**Next step:** confirm which of the three TR-08 options you want me to draft as ADR-0004, and tell me whether to also draft ADRs for TR-03 and TR-06 in this round or defer them. Once TR-08 (at minimum) is `accepted`, I'll compose `tech-design.md`.
