# Actions Log — eval-4-stage2-enumerates / without_skill / run-1

## Context

The user asked for ADRs for the `self-hosted-application-platform` capability. The tech-requirements file referenced (`docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`) is the 7-TR fixture variant, with `reviewed_at: 2026-04-26` matching today's date — review gate satisfied.

## Setup

- Setup `cp` command was denied by the sandbox. The target file already existed at the expected path with content matching the fixture (verified via Read), so no copy was needed.
- All Bash commands were denied throughout. Worked via Read / Write only.

## Investigation

- Read fixture: `.claude/skills/define-tech-design/evals/fixtures/tech-requirements-reviewed-naive.md` — 7 TRs (TR-01..TR-07).
- Read existing capability tech-requirements: `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` — confirmed `reviewed_at: 2026-04-26`.
- Read capability `_index.md` for outcomes, business rules, KPIs, tiebreaker order.
- Read existing `tech-design/_index.md` and `tech-design/technical-requirements.md` (the longer 23-REQ derivative) for additional context (REQ-NN and Q-NN labels) — used as supplementary context, not as the primary numbering source.
- Inspected existing repo ADRs: `0001-use-madr-for-architecture-decision-records.md` and `0007-standard-api-error-response.md` to mirror MADR frontmatter conventions (title `[NNNN] …`, `category`, `status`, `date`, `deciders`/`consulted`/`informed`).

## Decisions about scope

- Authored 6 capability-scoped ADRs covering the non-trivial decisions implied by TR-01..TR-06. TR-07 (Cloudflare → GCP path) is a constraint, not a decision, so no ADR for it on its own.
- Placed ADRs under `docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/` per the tech-design `_index.md` guidance ("Capability-scoped ADRs" as Stage 2).
- Used numbering 0001–0006 within the capability; cross-cutting ADRs in `docs/content/r&d/adrs/` keep their own numbering.
- All ADRs `status: proposed` — they are draft Stage-2 outputs; the operator (deciders: `[operator]`) accepts them by changing status.

## Files written (natural paths)

1. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/_index.md`
2. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0001-tenant-isolation-model.md`
3. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0002-compute-substrate.md`
4. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0003-contract-versioning-coexistence.md`
5. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0004-tenant-scoped-observability.md`
6. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0005-zero-downtime-tenant-updates.md`
7. `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0006-tenant-data-export-and-import-format.md`

## Files written (eval artifact copies)

Identical copies under `.../without_skill/run-1/outputs/artifacts/adrs/`:
- `_index.md`
- `0001-tenant-isolation-model.md`
- `0002-compute-substrate.md`
- `0003-contract-versioning-coexistence.md`
- `0004-tenant-scoped-observability.md`
- `0005-zero-downtime-tenant-updates.md`
- `0006-tenant-data-export-and-import-format.md`

## Notes

- I did not modify `tech-requirements.md` (the user said it is reviewed and unchanged; numbering is append-only).
- I did not yet produce the Stage-3 composed tech-design document — that comes after the ADRs settle.
- ADRs reference each other intentionally (0001 → 0002 → 0005; 0001 → 0004; 0003 → 0006) to make the dependency chain explicit.
