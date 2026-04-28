---
name: define-technical-requirements
description: >-
  Guide the user through extracting a living technical-requirements document
  from an already-defined business capability and its user experiences. The
  output is `tech-requirements.md` with each requirement linked back to its
  source in the capability or a UX, ready for human review before any
  architectural decisions are made. Use this skill whenever the user wants to
  surface, list, extract, or document the technical constraints implied by a
  capability — phrases like "extract tech requirements for {capability}",
  "what does {capability} require technically", "list the technical
  requirements", "we need TRs for this capability", or as the first step
  before drafting ADRs. Do NOT use to make decisions, propose options, or
  pick technologies (those belong to `plan-adrs` / `define-adr`). Do NOT use
  to define the capability or UX themselves (use `define-capability` /
  `define-user-experience` first).
---

# Define Technical Requirements for a Capability

This skill produces `tech-requirements.md` for a capability — a flat, numbered list of technical requirements, each linked back to its source in the capability doc or a UX doc. The doc is **living**: it gets re-extracted as the capability and UX docs evolve, and it gates the start of architectural decision-making (`plan-adrs` will refuse to run until a human has reviewed it).

This is **Step 6** of the capability development lifecycle (the first technical step). Steps before it (`define-capability`, `define-user-experience`) define *what the business does* and *how users experience it*. Steps after it (`plan-adrs`, `define-adr`, `plan-tech-design`) decide *how to build it*. This skill draws the line between the two: it captures what the technical solution is **forced** to do, without choosing how.

## Why this matters

If the human starts choosing technologies before they have written down the technical constraints, the choices will be unjustified — and unjustified choices are where production bugs hide. Pulling requirements out of the capability and UX docs first turns "use Postgres" from a hunch into something that either does or doesn't satisfy a written constraint.

The list also makes the chain of evidence inspectable. Every later ADR cites TR-NN identifiers. Every line of the eventual tech-design traces TR → ADR → component. If that chain breaks, the design has unjustified pieces. So this document is the foundation everything stands on; getting the discipline right here pays back in every later stage.

## Preconditions — read everything first

Before eliciting anything:

1. **Find the parent capability** at `docs/content/capabilities/{name}/_index.md`. If it is still a flat file (`{name}.md`), stop and ask the user to migrate it to page-bundle form first (see `define-user-experience`).
2. **Read the capability doc end-to-end.** Internalize stakeholders, business rules, success criteria, out-of-scope list.
3. **Read every UX doc** under `docs/content/capabilities/{name}/user-experiences/`. Not just the one the user mentioned — the requirements have to serve all of them.
4. **Skim `docs/content/r&d/adrs/`** for prior shared decisions that constrain this capability (cloud provider, network topology, error response format, identifier standard, etc.). Cite them as constraints, don't re-decide them.
5. **Note the repo's house patterns** from `CLAUDE.md`: chi/bedrock Go service shape, `pkg/errorpb` for errors, no humus framework, Cloudflare → GCP topology. These are inherited constraints, not requirements you discover.

If the capability or any UX docs are missing, **stop and route the user to `define-capability` / `define-user-experience` first.** Tech requirements derived from missing inputs are tech requirements with missing reasons. Refuse to proceed; do not invent the missing inputs in your head.

## Goal

Produce or update `docs/content/capabilities/{name}/tech-requirements.md` from `assets/template.md`.

## What is and is not a requirement

A **requirement is forced** by the capability or a UX. Examples:
- "The system must allow a tenant to migrate off the platform without downtime" (UX demands it)
- "All inter-service calls must traverse the Cloudflare → GCP path" (prior shared ADR)
- "Tenant data must remain isolated such that no tenant can read another's state" (capability business rule)

A **decision is chosen** from multiple options that all satisfy a requirement. Examples:
- "Use Postgres logical replication" (one of several ways to satisfy "no-downtime migration")
- "Use mTLS for service-to-service auth" (one of several ways to satisfy "isolated tenant data")

**Decisions do not go in `tech-requirements.md`.** If the user volunteers a decision during extraction, capture it in **Open Questions** for later ADR work — never as a requirement. This separation is what makes the eventual ADRs meaningful: an ADR's job is to pick one option among several that all satisfy the underlying TR. If you let a decision in here, you erase the alternatives that the ADR was supposed to weigh.

## Append-only TR identity

Requirements are identified `TR-01`, `TR-02`, … and **numbers are append-only forever**. When re-extracting on a living doc:

- Preserve every existing TR-NN whose source link still resolves.
- Append newly-discovered requirements at the end with the next free number.
- If a TR's source no longer resolves (UX deleted, capability rule rewritten), **flag it** with `> ⚠️ source no longer resolves — human review` — do not delete it. The human resolves the flag.
- Never renumber. Gaps are honest history. Downstream ADRs cite TR-NN, so renumbering silently breaks ADR provenance.

## Source links

Every TR must link back. Use Hugo-relative paths:
- Capability section: `[Capability §Business Rules](_index.md#business-rules)`
- UX page or section: `[UX: upload-photo §Edge Cases](user-experiences/upload-photo.md#edge-cases)`
- Prior shared ADR: `[ADR-0006](/r&d/adrs/0006-resource-identifier-standard/)`
- Repo pattern from CLAUDE.md: `[CLAUDE.md §Critical Go Patterns](/CLAUDE.md)` (or just cite inline)

If a requirement has multiple sources, list them all. Multi-sourced requirements are usually the most important ones — they are the ones that show up from more than one direction.

## The exit gate

`tech-requirements.md` carries a frontmatter field:

```yaml
reviewed_at: null   # set to an ISO date once a human has reviewed
```

The downstream `plan-adrs` skill **refuses to enumerate decisions** until `reviewed_at` is a date *newer* than the file's last modification time (i.e. the human reviewed the current contents, not an old version).

When you finish extracting, tell the user explicitly:

> "I've extracted the technical requirements to `tech-requirements.md`. **Review and edit it before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to today's ISO date (or tell me to do it on your behalf and I'll record your verbal confirmation). The next skill — `plan-adrs` — won't enumerate decisions until that's set."

If the user invokes this skill again later (re-extraction), first re-read the existing `tech-requirements.md`, preserve numbering, surface stale source links, and reset/null out `reviewed_at` when meaningful changes are made (so the human re-reviews the new contents).

## Conversation discipline

- **Extract, don't design.** No options, no chosen approaches, no technology names. If you find yourself proposing a database, a protocol, or a library, stop — that is the next skill's job.
- **Quote the source where possible.** If a UX step says "the user expects to feel that the upload is safe even with flaky wifi", a derived requirement might be "the system must tolerate intermittent connectivity during uploads without losing user-perceived progress" — and the link makes the derivation auditable.
- **Push back on premature solutions.** If the user says "we'll need Postgres for tenant state", redirect: "That's a decision for `plan-adrs` / `define-adr`. The requirement underneath is what?"
- **Don't invent requirements.** If nothing in the capability or UX implies a thing, it isn't a requirement. Capture it as an open question or push back to expand a UX doc.
- **Living-doc framing.** Tell the user the doc is meant to drift as UX docs evolve. Re-extract on demand. Git diff is the review surface — don't lament that the doc changed; that's the point.

## Producing the document

Use `assets/template.md`. Fill `{{capability_name}}`, `{{requirements}}` (the numbered TR sections), and `{{open_questions}}`. Each TR follows the shape commented in the template — heading, **Source:** line(s), **Requirement:** paragraph, **Why this is a requirement, not a decision:** line.

Save to `docs/content/capabilities/{name}/tech-requirements.md` (page-bundle form for the parent capability is required).

Keep the Hugo/Docsy frontmatter intact — it is required for the page to render in the published docs site. Set `reviewed_at: null` on a fresh extraction or when re-extraction has changed any TR; only carry forward an existing `reviewed_at` value if the file's TR contents are byte-equal to what you started with.

After producing the doc, tell the user:
- Where it was saved
- How many TRs were added/preserved/flagged
- Which open questions are parked for the ADR stage
- That `plan-adrs` is gated on `reviewed_at` and won't run until they set it

## Output format

Plain markdown matching `assets/template.md`. No code fences around the whole document. Mermaid is unnecessary here; this stage is text-only.

## Examples of good redirection

**Example 1 — premature decision dressed as a requirement:**
User: "TR-04 should be: use Postgres for tenant metadata storage with logical replication for cross-region failover."
Response: "That's a Stage-2 decision, not a requirement. 'Use Postgres' is one of several options that could satisfy the underlying constraint. The requirement is closer to: 'tenant metadata must be queryable by tenant ID with strong consistency, and survive a single-region failure within {N} minutes' — what does the capability or UX actually demand here? I'll capture 'consider Postgres + logical replication' as an open question for the corresponding ADR."

**Example 2 — vague slogan as requirement:**
User: "TR-N: the system must be reliable."
Response: "'Reliable' is a slogan, not a requirement until it's measurable and sourced. What in the capability or a UX forces a specific reliability constraint — an availability target, a max recovery time, a guaranteed durability? Quote the source and we can phrase the TR around it."

**Example 3 — re-extraction surfaces a stale link:**
Response: "TR-06's source link is `user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists` — that section was renamed. I've flagged it with ⚠️ rather than rewriting or deleting it. Want to re-source it now, or leave the flag for review? Either way, I'll null `reviewed_at` since the file changed; you'll need to re-set it before `plan-adrs` will run."
