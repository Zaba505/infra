I've extracted the business requirements for the Self-Hosted Application Platform capability to `docs/content/capabilities/self-hosted-application-platform/business-requirements.md`.

## Summary

- **54 BRs added** (BR-01 through BR-54), all from a fresh extraction (no prior `business-requirements.md` to preserve).
- **0 BRs preserved** (none existed) and **0 BRs flagged** (no stale source links since this is the first extraction).
- Sources span the parent capability doc and all seven UX docs (`host-a-capability`, `migrate-existing-data`, `move-off-the-platform-after-eviction`, `operator-initiated-tenant-update`, `platform-contract-change-rollout`, `stand-up-the-platform`, `tenant-facing-observability`).
- `reviewed_at` is set to `null` on the new doc.

## Open questions parked for the next stage

- **Explicit anchor IDs are missing on the capability and UX section headings.** Source links currently resolve at the page level; deep links into specific sections will require `{#anchor-id}` annotations on the target headings before `define-technical-requirements` cites BRs by section. Worth fixing in the source docs.
- **Numeric thresholds** (2x migration footprint, 30-day retention, 2-hr/week budget, two status-update cycles) are deliberately left in BR text as descriptions; the numbers belong in TRs/KPIs and the capability doc already carries them.
- **Deeper backup-tier policy after the 30-day retention** is still TBD per the eviction UX itself — this is a BR-shaped gap to fill in the capability or UX before BR-29/BR-54 are complete.
- **Backup standard** (BR-16) — the capability says the platform defines a standard but does not state it.
- **Tenant-facing pending-update view, signal bundle composition, migration concurrency model, drift detection mechanism** — all parked for TR/ADR.

## Next steps

**Review and edit `business-requirements.md` before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to today's ISO date (or tell me to do it on your behalf and I'll record your verbal confirmation). The next skill — `define-technical-requirements` — won't extract TRs until that's set.

I skipped the rendered-page validation step (`hugo server` preview) per the harness instruction to do so. The `ref` shortcodes are page-level only (no section anchors), which means a Hugo build should resolve them cleanly even without explicit `{#}` anchors on the source headings — but please re-run the validation locally before promoting.
