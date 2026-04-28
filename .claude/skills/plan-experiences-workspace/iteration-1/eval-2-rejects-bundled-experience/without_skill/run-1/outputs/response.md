# Plan: User Experience for `self-hosted-personal-media-storage`

**Capability:** Self-Hosted Personal Media Storage
**Capability doc:** `.claude/skills/plan-experiences/evals/fixtures/self-hosted-personal-media-storage.md`
**Reviewed at:** 2026-04-28 (file unchanged since)

## Requested Experience

The user asked for a single user experience titled **"User manages their content"** that covers:

- Uploading
- Sharing
- Browsing
- Deleting

## Concern: This bundles four distinct experiences

A single "user manages their content" experience that spans uploading, sharing, browsing, and deleting bundles four end-to-end user journeys with different triggers, actors, inputs, success conditions, and risk profiles. The capability doc itself enumerates these as separate triggers (see *Triggers & Inputs*):

- **A user uploads a file** (single, manual upload) — and the related automated-backup and bulk-import variants.
- **A user views, downloads, or organizes their own content.**
- **A content owner shares content with another user or a shared group.**
- **A user deletes their own content.**

Each of these has materially different concerns:

| Journey | Primary actor | Secondary actor | Distinct concerns from capability doc |
|---|---|---|---|
| Upload | Owner | — | Single vs. automated backup vs. bulk import; capture metadata preservation |
| Browse / view / download | Owner | — | On-demand archive export (longevity / operator-succession path) |
| Share | Owner | Recipient(s) / shared group | Private-by-default rule; affected-party concerns; one-to-one vs. group |
| Delete | Owner | — | 30-day retention window, eventual purge, accident recovery |

Bundling them into one UX would:

1. **Hide the recipient actor that only sharing introduces.** Browsing/uploading/deleting are single-actor flows; sharing is the only one with a second authorized human in the loop, plus the affected-party stakeholder the capability doc explicitly calls out.
2. **Collapse distinct success criteria.** The "zero data loss" KPI bites hardest on delete (30-day retention) and upload (durability). It applies differently to share (access-control correctness, not durability) and to browse (export-on-demand availability).
3. **Obscure the export-as-longevity flow inside "browsing."** The capability's *Operator succession* rule treats user-initiated full-archive export as a first-class durability mechanism — that is not the same journey as casually viewing a photo, even though both start from "I want to look at my stuff."
4. **Make later artifacts harder to write.** Business requirements, technical requirements, and ADRs derived from one omnibus UX will be muddier than ones derived from four focused UXes.

## Recommendation

Plan **four** user experiences instead of one, each scoped to a single end-to-end journey:

1. **User uploads content** — manual single-file, automated device backup, and bulk import from a prior provider. One trigger family, one actor, durability-focused success.
2. **User browses and exports their content** — view, download a single item, and pull a full on-demand archive (the longevity path). One actor, availability-focused success.
3. **User shares content with another user or group** — the only multi-actor flow; covers private-by-default enforcement and the affected-party stakeholder.
4. **User deletes their content** — covers the 30-day retention window, accident recovery, and eventual purge.

If the user genuinely wants a single document, a reasonable narrower alternative is **"User uploads and deletes their own content"** (one actor, lifecycle-focused) and to handle browse/export and share as their own UXes — but the four-way split above maps most cleanly onto the capability doc's own structure.

## Awaiting confirmation

Before filing any GitHub issues, please confirm one of:

- **(A) Split into four UXes** as recommended above (preferred). I will draft titles and one-line scopes for each, then print the four `gh issue create` commands.
- **(B) Different split** you specify (e.g. three UXes, or a different grouping).
- **(C) Proceed with the single bundled UX** "user manages their content" as originally requested, accepting the trade-offs above. I will print one `gh issue create` command.

No issues have been filed. No `gh` commands have been executed. The placeholder commands corresponding to options (A) and (C) are written to `gh-commands.txt` alongside this file for reference; pick one set after you confirm.
