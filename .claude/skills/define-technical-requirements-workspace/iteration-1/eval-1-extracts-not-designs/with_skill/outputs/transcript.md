# Transcript — define-technical-requirements run for Self-Hosted Application Platform

This is the user-facing response the assistant would deliver after running the `define-technical-requirements` skill against the `self-hosted-application-platform` capability.

---

I've extracted the technical requirements from the **Self-Hosted Application Platform** capability and its seven user experiences and saved them to:

`docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`

## What I read

- `_index.md` — the capability doc (purpose, stakeholders, business rules, KPIs, out-of-scope).
- All seven UX docs under `user-experiences/`:
  - `host-a-capability.md`
  - `migrate-existing-data.md`
  - `move-off-the-platform-after-eviction.md`
  - `operator-initiated-tenant-update.md`
  - `platform-contract-change-rollout.md`
  - `stand-up-the-platform.md`
  - `tenant-facing-observability.md`
- The repo-level shared ADRs in `docs/content/r&d/adrs/` (0001–0007) for cross-capability constraints.
- `CLAUDE.md` for inherited house-pattern constraints.

## What I produced

- **41 TRs** (TR-01 through TR-41), all freshly assigned (this is a first extraction — no prior `tech-requirements.md` existed).
- **0 preserved** from a prior run — first extraction.
- **0 flagged** with stale source links — first extraction.
- **6 open questions** parked for the ADR stage. None of them are decisions in disguise; each is either an explicit "TBD" in the UX docs (e.g. deeper backup-tier retention policy) or a constraint that the capability/UX docs don't yet make explicit (e.g. backup RPO/RTO standard, concurrent-migration capacity).

## A note on what I deliberately did *not* include

I kept this document strictly to constraints — no chosen technologies, protocols, or libraries. A few examples of things that *felt* like requirements but are actually decisions and have been left out:

- "Use GitHub Issues" — TR-11 says "GitHub issues against this repository must be the sole engagement surface" because the UX docs name GitHub specifically and there is no other channel; this is forced by the UX, not a free choice. If you would prefer to phrase it as "a single ticketed engagement surface" and treat GitHub specifically as a Stage-2 ratification, flag it on review and I'll relax the wording.
- "Use email for alerting" — same situation. TR-24 names email because the UX explicitly says email; if you want it abstracted to "a push channel that can deliver per-signal alerts", say so on review.
- Specific reproducibility tooling, secret-management technology, observability stack — none mentioned. Those are explicitly Stage-2 ADRs.

## Coverage summary, by source

- **Capability business rules / KPIs / outputs:** TR-01 through TR-09, TR-13, TR-14, TR-18, TR-19, TR-31, TR-37, TR-38, TR-39 trace at least partly here.
- **Host-a-Capability:** TR-11 through TR-14, TR-20, TR-33, TR-41.
- **Migrate-Existing-Data:** TR-15, TR-16, TR-17, TR-40.
- **Move-Off-After-Eviction:** TR-08 through TR-10, TR-36, TR-38.
- **Operator-Initiated-Tenant-Update:** TR-12, TR-36.
- **Platform-Contract-Change-Rollout:** TR-32, TR-33, TR-34, TR-35, TR-36.
- **Stand-Up-the-Platform:** TR-02, TR-03, TR-25 through TR-31, TR-39.
- **Tenant-Facing-Observability:** TR-20, TR-21, TR-22, TR-23, TR-24.

## Next step — review gate

**Review and edit `tech-requirements.md` before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to today's ISO date (or tell me to do it on your behalf and I'll record your verbal confirmation). The next skill — `plan-adrs` — won't enumerate decisions until that's set.

Specific things I'd appreciate your eyes on during review:

1. **TR-11 / TR-24 wording.** Are GitHub and email named-tech requirements (as the UX docs literally say), or do you want them abstracted to "single ticketed engagement surface" and "push notification channel"? Either is defensible.
2. **TR-15** is numerically tight (2× steady-state). Confirm the cap is genuinely a forced constraint and not a starting heuristic.
3. The **six open questions** at the bottom — if any of them have answers in your head that I missed in the docs, please either edit the source doc and I'll re-extract, or tell me and I'll capture them inline.

This is a living document. Re-run this skill any time the capability or a UX changes; numbering is append-only forever, so existing ADRs will keep citing the right TRs across re-extractions.
