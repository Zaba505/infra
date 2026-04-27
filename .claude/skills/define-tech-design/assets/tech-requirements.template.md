---
title: "Technical Requirements"
description: >
    Technical requirements extracted from the {{capability_name}} capability and its user experiences. Each requirement links back to its source. Decisions belong in ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-tech-design` skill will refuse to proceed to ADRs (Stage 2) until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [{{capability_name}}]({{< relref "_index.md" >}})

## How to read this {#how-to-read-this}

Each requirement is **forced** by the capability or a user experience — it constrains what the system must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review.

## Requirements {#requirements}

{{requirements}}

<!--
Each requirement should follow this shape. NOTE the heading-ID convention `{#tr-NN}` —
this is what ADRs link to, and it must stay stable across rewordings of the heading text.

### TR-01: {short imperative phrase} {#tr-01}
**Source:** [Capability §Section]({{</* relref "_index.md#section-id" */>}}) · [UX: name §Section]({{</* relref "user-experiences/name.md#section-id" */>}})
**Requirement:** {one paragraph describing the constraint in technical-but-implementation-neutral terms}
**Why this is a requirement, not a decision:** {what makes it forced rather than chosen — usually a citation to the source's intent}

If the source heading lacks an explicit `{#id}`, link to the page without an anchor rather
than guessing the auto-generated slug:
**Source:** [UX: name]({{</* relref "user-experiences/name.md" */>}})

For requirements whose source no longer resolves after re-extraction, prepend:
> ⚠️ source no longer resolves — human review

For removed requirements, keep the heading and add:
> 🗑️ removed on {date} — see git history. Number is reserved and will not be reused.
-->

## Open Questions {#open-questions}

Things the user volunteered as solutions during extraction (parked for Stage 2), or constraints the capability/UX docs don't yet make explicit.

{{open_questions}}
