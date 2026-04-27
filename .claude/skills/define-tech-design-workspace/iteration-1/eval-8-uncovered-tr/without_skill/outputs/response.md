I composed `tech-design.md`, but I want to flag a gap before you read it: only 4 of the 8 technical requirements actually have an accepted ADR backing them.

## Coverage map

| Requirement | ADR |
|---|---|
| TR-01 Tenant isolation | ADR-0001 (per-tenant Firestore namespace) |
| TR-02 Contract change rollout | ADR-0002 (semver in contract package path) |
| TR-03 Per-tenant observability | **none** |
| TR-04 No-downtime updates | ADR-0001 (storage layer only — rollout orchestration undecided) |
| TR-05 Evicted-tenant export | ADR-0003 (on-demand GCS signed URL) |
| TR-06 Lossless, idempotent migrations | **none** |
| TR-07 Inter-service traffic via Cloudflare → GCP | **none capability-scoped** (only a "prior shared decision" reference) |
| TR-08 Graceful degradation on regional outage | **none** |

## What I did

- The "Designed components" section (§3) describes only what the accepted ADRs actually say. I did not extrapolate beyond their text — e.g. ADR-0003 doesn't talk about export format, so I noted that explicitly rather than invent one.
- TR-03, TR-06, TR-07, and TR-08 are listed in §4 as **Open items** with the specific questions each one needs to answer. I deliberately did **not** design solutions for them in-line, because doing so would either (a) fabricate decisions you haven't made or (b) front-run an ADR that should be authored separately.
- I also called out that TR-04 is only **half** covered: ADR-0001 enables per-tenant rollout but doesn't specify the orchestration (drain/dual-run/cutover/rollback). Flagged in §4.5.
- The doc opens with a `Status: Incomplete` callout so a reader knows up-front this isn't a finished design.

## Recommended next step

Author ADRs for TR-03, TR-06, TR-07, TR-08, and a TR-04 rollout ADR. Then re-compose this document — at that point the "Open items" section disappears.

## Files

- `docs/content/capabilities/self-hosted-application-platform/tech-design.md` (new)
- `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (from fixture)
- `docs/content/capabilities/self-hosted-application-platform/adrs/{0001,0002,0003,_index}.md` (from fixture)
