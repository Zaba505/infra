---
name: define-tech-design
description: Guide the user through defining the technical design for an already-defined business capability and its user experiences. Produces three artifacts in a strict three-stage flow: (1) a living technical-requirements document extracted from the capability + UX docs, (2) one or more capability-scoped MADR ADRs, and (3) a composed human-friendly tech-design document describing the final state. Use this skill whenever the user wants to design, architect, or technically scope how a capability will be built — phrases like "design the system for {capability}", "what's the architecture for X", "tech design for this capability", "ADRs for this capability", "how do we build {capability}". Do NOT use for implementation (writing code or Terraform), for defining capabilities or UX (use `define-capability` / `define-user-experience` first), or for shared cross-capability decisions (those belong in `docs/content/r&d/adrs/` and are out of scope here).
---

# Define a Technical Design for a Capability

This skill produces the technical design that realizes a business capability and its user experiences. Output is **three artifacts**, generated in **three strictly gated stages**:

1. `tech-requirements.md` — extracted technical requirements, each linked back to its source. Living doc, tracked in git, reviewed by a human before any decisions are made.
2. `adrs/{NNNN}-{name}.md` — one MADR 4.0.0 ADR per decision, numbered locally per capability.
3. `tech-design.md` — a composed, human-friendly narrative of the final-state design that synthesizes ADR outcomes plus diagrams.

Cross-capability decisions are **out of scope**. If a decision is shared, flag it and stop — it belongs in `docs/content/r&d/adrs/` via a separate flow.

## Why this matters

Tech design errors are expensive. A bad UX gets caught in design review; a bad tech design gets caught in production. The three-stage flow exists to make the chain of evidence inspectable: every line of the final design traces back through an ADR, through a technical requirement, to a capability rule or UX step. If that chain breaks, the design has unjustified pieces — and unjustified pieces are where the bugs live.

The hard stage gates exist because skipping ahead is the failure mode. A skill that helpfully drafts ADRs before the human has checked the requirements will produce confidently-wrong designs.

## Preconditions — read everything first

Before eliciting anything:

1. **Find the parent capability** at `docs/content/capabilities/{name}/_index.md`. If it is still a flat file (`{name}.md`), stop and ask the user to migrate it to page-bundle form first (see `define-user-experience`).
2. **Read the capability doc end-to-end.** Internalize stakeholders, business rules, success criteria, out-of-scope list.
3. **Read every UX doc** under `docs/content/capabilities/{name}/user-experiences/`. Not just the one the user mentioned — the design has to serve all of them.
4. **Skim `docs/content/r&d/adrs/`** for prior shared decisions that constrain this capability (cloud provider, network topology, error response format, identifier standard, etc.). Cite them as constraints, don't re-decide them.
5. **Note the repo's house patterns** from `CLAUDE.md`: chi/bedrock Go service shape, `pkg/errorpb` for errors, no humus framework, Cloudflare → GCP topology. These are constraints inherited from the codebase.

If the capability or any UX docs are missing, stop and route the user to `define-capability` / `define-user-experience` first. Tech design with missing inputs produces tech design with missing reasons.

## Stage 1 — Extract technical requirements (HARD GATE before Stage 2)

### Goal

Produce or update `docs/content/capabilities/{name}/tech-requirements.md` — a flat, numbered list of technical requirements, each linked back to its source in the capability doc or a UX doc.

### What is and is not a requirement

A **requirement is forced** by the capability or a UX. Examples:
- "The system must allow a tenant to migrate off the platform without downtime" (UX demands it)
- "All inter-service calls must traverse the Cloudflare → GCP path" (prior shared ADR)
- "Tenant data must remain isolated such that no tenant can read another's state" (capability business rule)

A **decision is chosen** from multiple options that all satisfy a requirement. Examples:
- "Use Postgres logical replication" (one of several ways to satisfy "no-downtime migration")
- "Use mTLS for service-to-service auth" (one of several ways to satisfy "isolated tenant data")

**Decisions do not go in `tech-requirements.md`.** If the user volunteers a decision during Stage 1, capture it as an open question for Stage 2 — never as a requirement. This separation is what makes ADRs meaningful later: an ADR's job is to pick one option among several that all satisfy the requirement.

### Append-only TR identity

Requirements are identified `TR-01`, `TR-02`, … and **numbers are append-only forever**. When re-extracting on a living doc:

- Preserve every existing TR-NN whose source link still resolves.
- Append newly-discovered requirements at the end with the next free number.
- If a TR's source no longer resolves (UX deleted, capability rule rewritten), **flag it** with `> ⚠️ source no longer resolves — human review` — do not delete it. The human resolves the flag.
- Never renumber. Gaps are honest history. ADRs cite TR-NN, so renumbering silently breaks ADR provenance.

### Source links

Every TR must link back. Use Hugo-relative paths:
- Capability section: `[Capability §Business Rules](_index.md#business-rules)`
- UX page or section: `[UX: upload-photo §Edge Cases](user-experiences/upload-photo.md#edge-cases)`
- Prior shared ADR: `[ADR-0006](/r&d/adrs/0006-resource-identifier-standard/)`
- Repo pattern from CLAUDE.md: `[CLAUDE.md §Critical Go Patterns](/CLAUDE.md)` (or just cite inline)

If a requirement has multiple sources, list them all. Multi-sourced requirements are usually the most important ones.

### The hard gate

`tech-requirements.md` carries a frontmatter field:

```yaml
reviewed_at: null   # set to an ISO date once a human has reviewed
```

The skill **refuses to enter Stage 2** until `reviewed_at` is a date *newer* than the file's last modification time (i.e. the human reviewed the current contents, not an old version). When Stage 1 finishes, tell the user explicitly:

> "I've extracted the technical requirements to `tech-requirements.md`. **Review and edit it before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to today's date (or tell me to do it on your behalf and I'll record your verbal confirmation). I won't propose ADRs until that's set."

If the user invokes the skill again later, first re-read `tech-requirements.md`, check the gate, and either re-extract (if requested) or move to Stage 2.

### Conversation discipline in Stage 1

- **Extract, don't design.** No options, no chosen approaches in this stage.
- **Quote the source where possible.** If a UX step says "the user expects to feel that the upload is safe even with flaky wifi", a derived requirement might be "the system must tolerate intermittent connectivity during uploads without losing user-perceived progress" — and the link makes the derivation auditable.
- **Push back on premature solutions.** If the user says "we'll need Postgres for tenant state", redirect: "That's a decision for Stage 2. The requirement underneath is what?"
- **Don't invent requirements.** If nothing in the capability or UX implies a thing, it isn't a requirement. Capture it as an open question or push back to expand a UX doc.

## Stage 2 — ADRs, one decision at a time (HARD GATE before Stage 3)

### Goal

Produce one ADR per technical decision, in `docs/content/capabilities/{name}/adrs/{NNNN}-{kebab-name}.md`, in MADR 4.0.0 format matching the existing examples in `docs/content/r&d/adrs/`. Numbering is **local to this capability**, starting at `0001`.

Ensure `docs/content/capabilities/{name}/adrs/_index.md` exists with Docsy section frontmatter. Create it on first ADR.

### Enumerate decisions before solving any

Open Stage 2 by reading `tech-requirements.md` (now possibly edited by the human) and proposing the **set of decisions the requirements force**. Resist bundling — "where does tenant state live" and "how is tenant identity propagated" are two ADRs, not one. Confirm the decision list with the user before writing any single ADR.

### Per-ADR shape

Use `assets/adr.template.md`. Each ADR's `## Context and Problem Statement` must cite the TR-NN identifiers it addresses — this is the chain of evidence. Each ADR's `## Considered Options` must list **at least two** options. Each option's pros/cons must be expressed in terms of the cited TRs.

`status:` lifecycle:
- `proposed` — drafted but not yet accepted
- `accepted` — user has confirmed; this ADR is part of the design
- `superseded` — a later ADR replaces it (cite the superseder)

The skill **refuses to enter Stage 3** until every ADR in `adrs/` has `status: accepted` (or `superseded`, with the superseder also accepted).

### Flag-and-stop for shared decisions

If a decision is obviously cross-capability (touches Cloudflare topology, identity, networking, error response format, the resource identifier standard, etc.), do not draft it as a capability-scoped ADR. Surface it:

> "This decision looks shared across capabilities — it touches {topic}. It belongs as a shared ADR in `docs/content/r&d/adrs/`. Want to defer it and proceed with the rest, or pause Stage 2 to handle it separately?"

### Conversation discipline in Stage 2

- **One ADR at a time, accepted before the next is started.** Tangled in-progress ADRs produce tangled designs.
- **Anchor every option in TRs.** "Option B fails TR-04 because it requires a stateful coordinator that can't survive a tenant migration." If you can't cite TRs, the decision is either premature or the requirement is missing — go back to Stage 1.
- **Mirror back the ADR before writing.** State the decision and rationale in plain language; let the user confirm or correct, then write the file.
- **Don't ship `proposed` and walk away.** If the user is undecided, capture the open question in the ADR's *Open Questions* section and leave status `proposed` — but make clear the design isn't ready for Stage 3.

## Stage 3 — Compose tech-design.md (final state)

### Goal

Produce `docs/content/capabilities/{name}/tech-design.md` — a single human-friendly read-through of the resulting system. Not a list of ADRs. A *narrative* a new engineer can read top-to-bottom to understand what is being built and why.

The doc **composes ADR outcomes**; it does not restate ADR rationale. For *why*, readers click through to the ADR.

### Required sections

Use `assets/tech-design.template.md`:

- **Overview** — one paragraph summarizing the design.
- **Components** — services in `services/`, modules in `cloud/`, external systems. Include a Mermaid component diagram.
- **Key flows** — one Mermaid sequence diagram per UX, showing how the components realize that journey.
- **Data & state** — what's stored where, who owns it, lifecycle.
- **Requirement → ADR → realization map** — a table mapping each TR-NN to the ADR(s) that resolve it and the component(s)/flow(s) that embody it. This is the audit trail. Every TR must appear; if any TR has no ADR, Stage 3 is premature — return to Stage 2.
- **Deferred / open** — anything left unresolved.

### Synthesis, not new decisions

If composing the final-state doc reveals a gap — a TR with no ADR, an ADR with no realization in the components, two ADRs in tension — **do not paper over it in prose.** Loop back to Stage 2 for a new or amending ADR. The composed doc must reflect actual decisions, not narrative glue.

### Cross-link

In the capability's `_index.md`, optionally add a `## Tech Design` section with three links: tech-requirements, the `adrs/` index, and `tech-design.md`. Offer this; don't force it.

## How to run the conversation across stages

- **Always announce the stage.** "We're in Stage 1: extracting requirements." This keeps the user oriented and makes it obvious when you'd be skipping ahead.
- **Re-entry is normal.** The skill will be invoked many times across a capability's life. Each invocation: read all three artifacts, figure out which stage is current, check the gates, proceed.
- **Living-doc discipline.** Stage 1's output is meant to drift as UX docs evolve. Re-extract on demand. Git diff is the review surface — don't lament that the doc changed; that's the point.

## Output formats

All three artifacts are plain markdown with Hugo/Docsy frontmatter. Mermaid renders natively in Docsy. ADRs match MADR 4.0.0 as used in `docs/content/r&d/adrs/`.

## Examples of good redirection

**Example 1 — premature decision in Stage 1:**
User: "TR-05 should be: use Firestore for tenant metadata."
Response: "That's a Stage 2 decision, not a requirement. The requirement underneath is something like 'tenant metadata must be queryable by tenant ID with strong consistency' — what does the UX actually demand here? I'll capture 'consider Firestore' as an open question for the corresponding Stage 2 ADR."

**Example 2 — bundling decisions in Stage 2:**
User: "Let's do one ADR covering the whole tenant state design."
Response: "Let's split it. I see at least three decisions in there: where tenant state lives, how it's partitioned, and how it's accessed by services. Bundling them into one ADR means the rationale for each gets diluted, and we can't supersede one without revisiting the others. Which do you want first?"

**Example 3 — gap discovered in Stage 3:**
Response: "Composing the design surfaced a gap: TR-08 (graceful degradation when GCP region is unreachable) is in `tech-requirements.md` but no ADR addresses it. I'm not going to invent a story in `tech-design.md`. Returning to Stage 2 to draft an ADR for it — should I propose options?"
