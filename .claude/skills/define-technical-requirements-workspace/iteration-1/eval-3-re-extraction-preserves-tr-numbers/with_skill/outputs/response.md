# Re-extraction complete: `tech-requirements.md` for `self-hosted-application-platform`

I re-read the capability doc, every UX under `user-experiences/`, and the existing `tech-requirements.md` before regenerating. Numbering is preserved; nothing was renumbered or deleted.

## What changed

**Preserved as-is (TR-01 through TR-05, TR-07):** still sourced correctly. Wording carried forward verbatim. (TR-01's source link was lightly normalized to point at the capability's true heading slug `#business-rules--constraints` rather than the previous `#business-rules`; the section title is "Business Rules & Constraints" and Hugo slugifies the `&`.)

**Flagged for human review — TR-06:** the source link `user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists` no longer resolves. The section was renamed. Per the skill's append-only rules I did **not** rewrite the source or delete the TR. Instead I prepended:

```
> ⚠️ source no longer resolves — human review
```

…and noted likely re-source candidates inline (`#journey` or `#edge-cases--failure-modes`). You decide which is canonical and update the link.

**Newly extracted (appended at the end, TR-08 through TR-13):** these are constraints that the capability + UX docs clearly force but the prior file did not yet capture. They are forced, not chosen — each one cites a specific business rule, KPI, or UX section:

- **TR-08** — 1-hour rebuild from definitions, no manual snowflake config (capability *Reproducibility* KPI + stand-up-the-platform UX).
- **TR-09** — every provisioning phase must be cleanly reversible / "delete everything" must always be a viable rollback (stand-up-the-platform UX edge cases + inherited constraints).
- **TR-10** — readiness must be proven by exercising a canary tenant end-to-end, not by self-checks (stand-up-the-platform §Phase 4).
- **TR-11** — the platform-provided identity service must be capable of honoring the "lost credentials cannot be recovered" property (capability business rule, explicitly disqualifying any identity option that cannot honor it).
- **TR-12** — tenants (and migration jobs) must be packaged in the platform's accepted form and declare their resource needs up front (capability rule "Tenants must accept the platform's contract" + host-a-capability + migrate-existing-data inherited constraints).
- **TR-13** — platform-altering changes must be tracked + immutable so a preflight drift check is meaningful (stand-up-the-platform UX inherited constraints).

## `reviewed_at` was reset

The previous file had `reviewed_at: 2026-04-20`. Because the file changed (TR-06 flagged, six TRs appended), I set `reviewed_at: null` per the skill rules. **You need to re-review the file and re-set this to today's ISO date** (or tell me to record your verbal confirmation). The downstream `plan-adrs` skill will refuse to enumerate decisions until that's set.

## Open Questions parked for Stage 2

- On-demand vs continuously-available export (TR-05).
- Versioning scheme for contract changes (TR-02).
- Single fixed canary vs rotated canary for readiness (TR-10).
- Whether the 30-day post-eviction retention window deserves its own TR.
- Re-sourcing TR-06 to the correct section.

## Counts

- Preserved: 6 TRs (TR-01..TR-05, TR-07)
- Flagged: 1 TR (TR-06)
- Appended: 6 TRs (TR-08..TR-13)
- Removed: 0 (numbers are append-only)

Saved to `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`. Review it, edit anything that reads like a chosen solution rather than a constraint, re-source TR-06, then set `reviewed_at:` to today.
