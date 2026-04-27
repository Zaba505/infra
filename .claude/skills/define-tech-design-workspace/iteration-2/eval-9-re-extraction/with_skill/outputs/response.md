# Stage 1 — re-extraction complete (gate is closed)

We're in **Stage 1** of `define-tech-design`. I re-read the capability doc and all seven UX docs under `user-experiences/`, then re-extracted `tech-requirements.md` for `self-hosted-application-platform`.

## What I changed

**Preserved unchanged (append-only policy):** TR-01, TR-02, TR-03, TR-04, TR-05, TR-07. Their source links still resolve and their text is still backed by the cited capability rule or UX section. No renumbering — ADRs cite TR-NN, and silently re-shuffling would break that provenance.

**TR-06 — flagged, not deleted.** Its source link points to `user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists`. The UX page still exists and still describes a one-shot, lossless, retry-safe migration job (the TR text is still accurate), but the section anchor doesn't match any heading in the current page. Per the skill, I added the required flag inline on TR-06:

> ⚠️ source no longer resolves — human review.

…with a note that the TR itself is still supported by the page (likely §Journey or §Success), and that the human should pick the correct anchor or accept a page-level link. **TR-06 was not deleted and not renumbered.** I also restated the open issue under *Open Questions* so it's visible at the document level.

**Newly appended TRs (TR-08 → TR-11)** — surfaced by re-reading every UX, not just the one with the broken link:

- **TR-08** — Platform must be reproducible from definitions in ≤1hr. Forced by the capability's *Reproducibility* KPI and exercised by the `stand-up-the-platform` UX. Was missing from the prior file.
- **TR-09** — Identity offering must be capable of honoring "lost credentials cannot be recovered." Forced by the capability business rule that explicitly excludes identity options unable to honor this.
- **TR-10** — Designated successor operator with sealed/escrowed credentials. Forced by the operator-succession business rule. (TR-05 already covers the export half of succession; TR-10 covers the takeover half.)
- **TR-11** — Tenant components must be packaged in a single declared form with declared resource needs. Forced by the "tenants must accept the platform's contract" rule and exercised by `host-a-capability` (and `migrate-existing-data` for one-shot jobs).

I deliberately did **not** invent TRs that aren't grounded in the capability doc, a UX, or a prior shared decision. A few candidates (e.g. specific availability targets, cost ceilings) were considered and rejected because the capability explicitly declines to set them.

## Why I did not auto-fix the TR-06 link

Two reasons. First, the skill is explicit: when a source no longer resolves, **flag it for human review, do not silently rewrite or delete**. Second, "the section was renamed" implies a human decision about which section the requirement now derives from — `§Journey` and `§Success` both reasonably support the TR-06 text, and I shouldn't pick on your behalf.

## Hard gate is now closed

The frontmatter is `reviewed_at: null`. **I won't propose ADRs (Stage 2) until you've reviewed the re-extracted requirements and set `reviewed_at:` to today's ISO date** (or told me explicitly to record your verbal confirmation). Specifically, please:

1. Resolve the TR-06 source-link flag — pick the correct anchor or convert to a page-level link.
2. Sanity-check TR-08 → TR-11. They're new since your last review on 2026-04-20; if any of them are wrong-shape (decisions disguised as requirements, or not actually forced) push back and I'll redraft.
3. Re-confirm TR-01 → TR-05, TR-07. They're carried forward verbatim, but a re-extraction is a fine moment to re-read them.
4. Set `reviewed_at: 2026-04-26` (or today's date) and let me know you're ready for Stage 2.

Files written:

- `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (updated; gate closed)
</content>
