---
name: extract-business-requirements
description: >-
  Guide the user through extracting a living business-requirements document
  from an already-defined business capability and its user experiences. The
  output is `business-requirements.md` with each requirement linked back to
  its source in the capability or a UX, ready for human review before any
  technical requirements are written. Use this skill whenever the user wants
  to extract, list, or document the business needs implied by a capability —
  phrases like "extract business requirements for {capability}", "what does
  {capability} need to do", "list the business requirements", "we need BRs
  for this capability", or as the step that precedes drafting technical
  requirements. Do NOT use to write technical requirements (use
  `define-technical-requirements`). Do NOT use to define the capability or
  UX themselves (use `define-capability` / `define-user-experience` first).
  Do NOT use to make architectural decisions (those belong to `plan-adrs` /
  `define-adr`).
---

# Extract Business Requirements for a Capability

This skill produces `business-requirements.md` for a capability — a flat, numbered list of business requirements, each linked back to its source in the capability doc or a UX doc. The doc is **living**: it gets re-extracted as the capability and UX docs evolve, and it gates the start of technical-requirements work (`define-technical-requirements` will refuse to run until a human has reviewed it).

This is **Step 5** of the capability development lifecycle (the last business-side step). Steps before it (`define-capability`, `define-user-experience`) define *what the business does* and *how users experience it*. Steps after it (`define-technical-requirements`, `plan-adrs`, `define-adr`, `plan-tech-design`) translate business needs into technical constraints and decisions. This skill draws the line between the two: it captures **what the business demands of the system**, in business and user-outcome terms, without yet framing the technical solution.

## Why this matters

A capability says *what the business does* and a UX says *how that lands for one specific user accomplishing one specific goal*. Neither is testable on its own. Business requirements are: they are the auditable, reviewable contract between business intent and technical work.

Without this layer, technical requirements end up doing two jobs at once — restating business needs in technical clothing, and deriving technical constraints from those needs. That collapse is where chains of evidence go to die: a TR like "checkpoint upload progress every 10 seconds to durable storage" is impossible to evaluate without the underlying business need ("uploads must survive intermittent connectivity without losing user-perceived progress"). Pull the business needs out first and every later TR has something concrete to point to.

The list also makes the chain inspectable. Every TR cites BR-NN identifiers; every ADR cites the TRs that cite the BRs; every component in the eventual tech-design traces BR → TR → ADR → component. If that chain breaks, the design has unjustified pieces. So this document is the foundation everything technical stands on.

## Preconditions — read everything first

Before eliciting anything:

1. **Find the parent capability** at `docs/content/capabilities/{name}/_index.md`. If it is still a flat file (`{name}.md`), stop and ask the user to migrate it to page-bundle form first (see `define-user-experience`).
2. **Read the capability doc end-to-end.** Internalize stakeholders, business rules, success criteria, out-of-scope list.
3. **Read every UX doc** under `docs/content/capabilities/{name}/user-experiences/`. Not just the one the user mentioned — the requirements have to serve all of them.

If the capability or any UX docs are missing, **stop and route the user to `define-capability` / `define-user-experience` first.** Business requirements derived from missing inputs are business requirements with missing reasons. Refuse to proceed; do not invent the missing inputs in your head.

## Goal

Produce or update `docs/content/capabilities/{name}/business-requirements.md` from `assets/template.md`.

## What is and is not a business requirement

A **business requirement is forced** by the capability or a UX, framed as a business or user-outcome demand on the system. Examples:
- "An evicted tenant must be able to leave the platform with all of their data, in a portable form, without operator assistance" (UX demands it)
- "No tenant can ever observe another tenant's state" (capability business rule)
- "An operator-initiated tenant update must complete without end-user-visible downtime for online workloads" (UX success criterion)

Three things are **not** business requirements — and conflating any of them here breaks the downstream chain:

- A **technical requirement** is the *technical constraint derived from a BR.* The BR says "uploads must survive intermittent connectivity without losing user-perceived progress." The TR derived from that says "the system must persist upload progress to durable storage at intervals of at most N seconds." If the user volunteers a TR-shaped statement, capture it in **Open Questions** for the next stage — never as a BR.
- A **decision** is one of multiple options that all satisfy a BR. "Use S3 multipart upload" is a decision; "uploads must be resumable from the last completed segment" is the BR underneath. Decisions belong in ADRs.
- A **restatement of the capability** is not a BR. A BR *derives* with provenance — it says, in a measurable way, what the system must guarantee in service of the capability's intent. If the BR text is byte-equivalent to a sentence in the capability doc, you have not extracted anything.

This separation is what makes the eventual TRs and ADRs meaningful: a TR's job is to translate one or more BRs into technical constraints; an ADR's job is to pick one option among several that all satisfy the underlying TR. If you let TRs or decisions in here, you erase the layer the technical work was supposed to stand on.

## Append-only BR identity

Requirements are identified `BR-01`, `BR-02`, … and **numbers are append-only forever**. When re-extracting on a living doc:

- Preserve every existing BR-NN whose source link still resolves.
- Append newly-discovered requirements at the end with the next free number.
- If a BR's source no longer resolves (UX deleted, capability rule rewritten), **flag it** with `> ⚠️ source no longer resolves — human review` — do not delete it. The human resolves the flag.
- Never renumber. Gaps are honest history. Downstream TRs cite BR-NN, so renumbering silently breaks TR provenance (and through it, ADR provenance).

## Source links

Every BR must link back. **Use Hugo's `ref` shortcode for every internal link — never raw paths.** Hugo will fail the build on a broken `ref`; raw paths break silently when content is reorganized.

- Capability section: `[Capability §Business Rules]({{< ref "_index.md#business-rules" >}})`
- UX page or section: `[UX: upload-photo §Edge Cases]({{< ref "user-experiences/upload-photo.md#edge-cases" >}})`

**Section deep-links require an explicit anchor on the target heading.** Add `{#anchor-id}` to the heading you are linking to (e.g. `## Business Rules {#business-rules}`) before linking — Hugo's default slugify-from-heading-text breaks every implicit anchor as soon as a heading is reworded. If you cite a section that has no explicit anchor yet, pause the extraction and either add the anchor in the source doc or capture the missing-annotation as an open question.

If a requirement has multiple sources, list them all. Multi-sourced requirements are usually the most important ones — they are the ones that show up from more than one direction.

## The exit gate

`business-requirements.md` carries a frontmatter field:

```yaml
reviewed_at: null   # set to an ISO date once a human has reviewed
```

The downstream `define-technical-requirements` skill **refuses to extract TRs** until `reviewed_at` is a date *newer* than the file's last modification time (i.e. the human reviewed the current contents, not an old version).

When you finish extracting, tell the user explicitly:

> "I've extracted the business requirements to `business-requirements.md`. **Review and edit it before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to today's ISO date (or tell me to do it on your behalf and I'll record your verbal confirmation). The next skill — `define-technical-requirements` — won't extract TRs until that's set."

If the user invokes this skill again later (re-extraction), first re-read the existing `business-requirements.md`, preserve numbering, surface stale source links, and reset/null out `reviewed_at` when meaningful changes are made (so the human re-reviews the new contents).

## Conversation discipline

- **Extract, don't translate.** No technical phrasing, no implementation language, no chosen approaches. If you find yourself writing "the system must persist X to durable storage" or "the API must return within Y ms", stop — that is the next skill's job. The BR underneath is what the user or business *experiences* or *demands*, in their terms.
- **Quote the source where possible.** If a UX step says "the user expects to feel that the upload is safe even with flaky wifi", a derived BR might be "the system must preserve user-perceived upload progress across intermittent connectivity" — and the link makes the derivation auditable.
- **Push back on premature TRs.** If the user phrases a BR with a number ("must respond within 200ms", "must checkpoint every 10s"), redirect: "That's a TR — what does the underlying business or user need actually demand? The number belongs in `define-technical-requirements`, derived from the BR."
- **Push back on premature decisions.** If the user names a technology, protocol, or product ("use webhooks", "use Postgres"), redirect: "That's a decision for `plan-adrs` / `define-adr`. The BR underneath is what?"
- **Don't invent requirements.** If nothing in the capability or UX implies a thing, it isn't a requirement. Capture it as an open question or push back to expand a UX doc.
- **Living-doc framing.** Tell the user the doc is meant to drift as UX docs evolve. Re-extract on demand. Git diff is the review surface — don't lament that the doc changed; that's the point.

## Producing the document

Use `assets/template.md`. Fill `{{capability_name}}`, `{{requirements}}` (the numbered BR sections), and `{{open_questions}}`. Each BR follows the shape commented in the template — heading, **Source:** line(s), **Requirement:** paragraph, **Why this is a requirement, not a TR or decision:** line.

Save to `docs/content/capabilities/{name}/business-requirements.md` (page-bundle form for the parent capability is required).

Keep the Hugo/Docsy frontmatter intact — it is required for the page to render in the published docs site. Set `reviewed_at: null` on a fresh extraction or when re-extraction has changed any BR; only carry forward an existing `reviewed_at` value if the file's BR contents are byte-equal to what you started with.

After producing the doc, tell the user:
- Where it was saved
- How many BRs were added/preserved/flagged
- Which open questions are parked for the TR or ADR stage
- That `define-technical-requirements` is gated on `reviewed_at` and won't run until they set it

## Validate the rendered page

After saving (or re-extracting) `business-requirements.md`, run `cd docs && hugo server` and open the page in a browser to confirm it renders correctly — every BR's `ref` source link must resolve, and any flagged stale-source markers must be visible. A clean build is necessary but not sufficient; pages can build clean while rendering broken (missing/incorrect front-matter `type`, unresolved `ref` shortcode, broken section anchor). The doc is not done until the rendered page has been visually validated. If you cannot run the preview in this session, say so explicitly to the user rather than reporting the doc as done.

## Output format

Plain markdown matching `assets/template.md`. No code fences around the whole document. Mermaid is unnecessary here; this stage is text-only.

## Examples of good redirection

**Example 1 — TR dressed as a BR:**
User: "BR-04 should be: the system must checkpoint upload progress every 10 seconds to durable storage."
Response: "That's a TR, not a BR. 'Checkpoint every 10 seconds to durable storage' is the technical translation. The BR underneath is closer to: 'uploads must survive intermittent connectivity without losing user-perceived progress' — what in the capability or UX actually demands this? I'll capture 'consider 10s checkpoint cadence to durable storage' as an open question for the corresponding TR work."

**Example 2 — decision dressed as a BR:**
User: "BR-N: use S3 multipart upload for resumable uploads."
Response: "That's a Stage-2 decision, not a BR. 'Use S3 multipart' is one of several options that could satisfy the BR underneath. The BR is closer to: 'uploads must be resumable after a connection drop without re-sending completed segments.' I'll park 'consider S3 multipart upload' as an open question for the ADR stage."

**Example 3 — vague slogan as BR:**
User: "BR-N: the system must be reliable."
Response: "'Reliable' is a slogan, not a BR until it's measurable and sourced. What in the capability or a UX forces a specific reliability outcome — a recovery expectation, a guaranteed durability, a maximum loss the user is willing to tolerate? Quote the source and we can phrase the BR around it."

**Example 4 — re-extraction surfaces a stale link:**
Response: "BR-06's source link is `user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists` — that section was renamed. I've flagged it with ⚠️ rather than rewriting or deleting it. Want to re-source it now, or leave the flag for review? Either way, I'll null `reviewed_at` since the file changed; you'll need to re-set it before `define-technical-requirements` will run."
