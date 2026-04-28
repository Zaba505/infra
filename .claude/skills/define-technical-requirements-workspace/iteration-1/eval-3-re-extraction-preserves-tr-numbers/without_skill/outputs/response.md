# Re-extraction of tech-requirements.md for self-hosted-application-platform

I re-walked the capability doc and all six UX docs under `docs/content/capabilities/self-hosted-application-platform/` and re-extracted the technical requirements. The existing TR-01 through TR-07 numbering is preserved — re-numbering would invalidate any external references (ADR drafts, issues, prior reviews) that already cite these TR IDs, and the existing TRs all remain valid against the current capability + UX text.

## What changed

- **`reviewed_at`** bumped from `2026-04-20` to `2026-04-27` to reflect today's re-extraction.
- **TR-06 source link fixed.** The previous link pointed at `migrate-existing-data.md#a-section-that-no-longer-exists`, which no longer resolves. The current `migrate-existing-data.md` organizes its content under sections like *Goal*, *Entry Point*, *Journey*, *Success*, *Edge Cases & Failure Modes*, etc. The lossless / retry-safe / idempotent-import content TR-06 depends on lives in **§Journey** (specifically the re-run contract declaration in step 2, the operator's re-run-contract check in step 3, and the resolution branches in step 7). I updated the source link to `[UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#journey)` to match the convention used by TR-03 (`§Journey`) and TR-04 (`§Success`).
- **TR-06 requirement text** lightly extended to surface the *re-run contract* declaration the UX now makes explicit, since that is the mechanism by which "no duplication on retry" is enforced. The original requirement statement is preserved verbatim as the leading sentence; the re-run-contract sentence is additive.

## What did not change

- **TR-01 through TR-05 and TR-07** — source links all still resolve (`_index.md#business-rules`, `platform-contract-change-rollout.md`, `tenant-facing-observability.md`, `operator-initiated-tenant-update.md`, `move-off-the-platform-after-eviction.md`, `CLAUDE.md`). Wording is byte-identical to the prior file.
- **TR numbering** — TR-01..TR-07 preserved. No insertions, no deletions, no renumbering.
- **Open Questions** — both still open; no new ones surfaced during re-extraction.

## No new TRs surfaced

I checked each UX (including `host-a-capability.md` and `stand-up-the-platform.md`, which the existing file does not cite directly) for technical constraints not already covered by TR-01..TR-07. The constraints I found there (operator-only access, packaging contract, secret-management offering, 1-hour reproducibility, 2hr/week maintenance budget) are either already covered transitively by existing TRs or are KPI/business-rule statements rather than technical requirements per the skill's scope. None warranted a new TR-08.

## Files

- Updated: `/home/carson/github.com/Zaba505/infra/.claude/skills/define-technical-requirements-workspace/iteration-1/eval-3-re-extraction-preserves-tr-numbers/without_skill/outputs/tech-requirements.md`

Ready for human review.
