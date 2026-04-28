---
title: "Technical Requirements"
description: >
    Technical requirements derived from the {{capability_name}} capability's business requirements (with capability and UX docs as context). Each TR cites the BR-NN it derives from. Decisions belong in ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from `business-requirements.md` (and the capability/UX docs) on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `plan-adrs` skill will refuse to enumerate decisions until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [{{capability_name}}]({{< ref "_index.md" >}})
**Business requirements:** [business-requirements.md]({{< ref "business-requirements.md" >}})

## How to read this

Each TR is **forced** — by a BR (the primary case), by a prior shared ADR, or by a repo-wide constraint. It says what the technical solution must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review. If something has no BR or inherited-constraint source, raise a missing BR back to `extract-business-requirements`.

## Requirements

{{requirements}}

<!--
Each requirement should follow this shape:

### TR-01: {short imperative phrase}
**Source:** [BR-03]({{< ref "business-requirements.md#br-03" >}}) · [UX: name §Section]({{< ref "user-experiences/name.md#section-anchor" >}})

<!--
The primary source for a TR is the BR-NN it derives from. Capability or UX section links may be added as additional context. Prior shared ADRs may also be cited.
Cross-links MUST use Hugo's `ref` shortcode — never raw paths like `(_index.md)` or `../foo.md`.
Section deep-links require an explicit anchor on the target heading (e.g. `## Business Rules {#business-rules}` or `### BR-03: …  {#br-03}`).
Hugo's build will fail on broken `ref`s; raw paths break silently when content is reorganized.
-->
**Requirement:** {one paragraph describing the technical constraint in implementation-neutral terms — what the solution is forced to do to deliver the BR}
**Why this is a TR, not a BR or decision:** {what makes it the technical translation rather than a restated BR or a chosen solution — usually quoting how the BR's demand forces this constraint}

For requirements whose source no longer resolves after re-extraction, prepend:
> ⚠️ source no longer resolves — human review

For removed requirements, keep the heading and add:
> 🗑️ removed on {date} — see git history. Number is reserved and will not be reused.
-->

## Open Questions

Things the user volunteered as solutions during extraction (parked for the ADR stage), or constraints the capability/UX docs don't yet make explicit.

{{open_questions}}
