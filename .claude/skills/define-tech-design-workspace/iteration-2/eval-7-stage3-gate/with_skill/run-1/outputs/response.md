# Stage 3 gate — refusing to compose `tech-design.md`

We're in **Stage 3** (compose `tech-design.md`), but the Stage 2 → Stage 3 gate is closed. I will not write `tech-design.md` yet, because doing so would synthesize a final-state narrative on top of decisions that haven't actually been made. That's the failure mode the gate exists to prevent.

## Why the gate is closed

Per the skill:

> The skill **refuses to enter Stage 3** until every ADR in `adrs/` has `status: accepted` (or `superseded`, with the superseder also accepted).

Current ADR status in `docs/content/capabilities/self-hosted-application-platform/adrs/`:

| ADR | Title | Status | Addresses |
|-----|-------|--------|-----------|
| 0001 | Tenant State Storage | **accepted** | TR-01, TR-04 |
| 0002 | Contract Versioning | **proposed** | TR-02 |
| 0003 | Tenant Eviction Export | **proposed** | TR-05 |

ADRs 0002 and 0003 are still `proposed`. Composing the design now would either (a) silently treat their proposed outcomes as accepted — papering over an unmade decision in prose, exactly what Stage 3 forbids — or (b) leave gaping holes in the narrative for TR-02 and TR-05. Neither is acceptable.

## A second, independent problem: missing ADRs

Even if 0002 and 0003 were accepted, Stage 3 would still be premature. The skill requires:

> Every TR must appear; if any TR has no ADR, Stage 3 is premature — return to Stage 2.

Mapping the current TRs to existing ADRs:

| TR | Covered by ADR? |
|----|-----------------|
| TR-01 (tenant isolation) | 0001 |
| TR-02 (contract versioning) | 0002 (proposed) |
| TR-03 (per-tenant observability scoping) | **none** |
| TR-04 (no-downtime tenant updates) | 0001 (partial — storage only; nothing covers the rollout/cutover mechanism) |
| TR-05 (evicted-tenant data export) | 0003 (proposed) |
| TR-06 (lossless tenant data import) | **none** |
| TR-07 (Cloudflare → GCP topology) | **none** (likely a flag-and-stop shared decision, see below) |

Three TRs have no ADR at all (TR-03, TR-06, TR-07), and TR-04 is only partially addressed.

## What needs to happen before Stage 3

1. **Resolve the proposed ADRs.** For each of 0002 and 0003, mirror back the proposed outcome, confirm or correct, then flip `status:` to `accepted`. If you're undecided, that's fine — but Stage 3 stays closed.
   - 0002: confirm semver in contract package path (vs. date-based or rolling+flags).
   - 0003: confirm on-demand export to a GCS signed URL (vs. continuous replication or scheduled snapshots). Note the open question recorded in `tech-requirements.md` was already pointing here.
2. **Draft the missing ADRs (back to Stage 2).** Proposed list, one decision each:
   - **0004 — Tenant observability scoping** (TR-03). Options likely include per-tenant Grafana Cloud stack, shared backend with tenant-id label enforcement at query time, or per-tenant OTLP collector with isolated storage.
   - **0005 — Tenant update rollout mechanism** (TR-04, alongside 0001). Options likely include blue/green per-tenant Cloud Run revisions with traffic split, rolling restart with readiness gates, or operator-driven canary.
   - **0006 — Tenant data import pipeline** (TR-06). Options likely include a one-shot import job that writes idempotently keyed by source-row hash, a streaming CDC ingest, or a customer-staged GCS handoff plus reconciliation.
3. **Flag-and-stop on TR-07.** TR-07 (Cloudflare → GCP topology) is explicitly an inherited shared constraint — it shouldn't be a capability-scoped ADR. It should be cited in `tech-design.md` as a constraint linking to the relevant shared ADR under `docs/content/r&d/adrs/`. If no such shared ADR exists yet, that's a separate flow, not this skill's job.

## What I did not do, and why

- **I did not draft `tech-design.md`.** Doing so would violate the Stage 3 gate and produce a doc whose chain of evidence is broken on TR-02, TR-03, TR-05, TR-06, and partially TR-04.
- **I did not silently flip 0002/0003 to `accepted`.** That decision is yours, not mine — the ADR template's `status` field is the explicit handshake.
- **I did not invent ADRs 0004–0006.** Stage 2 is "one ADR at a time, accepted before the next is started," and each one needs you to confirm the option choice. I've sketched the decision shape above so we can pick one to start with.

## Next step — your call

Tell me which of these you want first:

- (a) Walk through 0002 (contract versioning) — confirm or change the proposed option, then flip to accepted.
- (b) Walk through 0003 (eviction export) — same.
- (c) Start drafting 0004 (tenant observability scoping) for TR-03.
- (d) Start drafting 0005 (tenant update rollout) for TR-04's uncovered half.
- (e) Start drafting 0006 (tenant data import) for TR-06.

Once 0002–0006 are all `accepted` and TR-07 is flagged as a shared concern, I'll compose `tech-design.md` in one pass.
