---
name: define-technical-requirements
description: >-
  Guide the user through extracting a living technical-requirements document
  from a capability's reviewed business-requirements doc, plus its capability
  definition and user experiences for context. The output is
  `tech-requirements.md` with each TR citing the BR-NN it derives from (and
  optionally the capability or UX section that frames it), ready for human
  review before any architectural decisions are made. Use this skill whenever
  the user wants to surface, list, extract, or document the technical
  constraints implied by a capability's business requirements — phrases like
  "extract tech requirements for {capability}", "what does {capability}
  require technically", "list the technical requirements", "we need TRs for
  this capability", or as the step that follows reviewing
  business-requirements.md. Do NOT use to make decisions, propose options, or
  pick technologies (those belong to `plan-adrs` / `define-adr`). Do NOT use
  to extract business requirements (use `extract-business-requirements`). Do
  NOT use to define the capability or UX themselves (use `define-capability`
  / `define-user-experience` first).
---

# Define Technical Requirements for a Capability

This skill produces `tech-requirements.md` for a capability — a flat, numbered list of technical requirements, each citing the BR-NN it derives from (and, where useful, the capability or UX section that frames it). The doc is **living**: it gets re-extracted as the BRs and UX docs evolve, and it gates the start of architectural decision-making (`plan-adrs` will refuse to run until a human has reviewed it).

This is **Step 6** of the capability development lifecycle (the first technical step). The step before it (`extract-business-requirements`) produces `business-requirements.md` — the reviewed, human-approved statement of what the business demands of the system. Steps after it (`plan-adrs`, `define-adr`, `plan-tech-design`) decide *how to build it*. This skill draws the line between the business layer and the technical layer: it translates each BR into the technical constraints **forced** by it, without choosing how.

## Why this matters

If the human starts choosing technologies before they have written down the technical constraints, the choices will be unjustified — and unjustified choices are where production bugs hide. Translating each BR into the technical constraint(s) it forces turns "use Postgres" from a hunch into something that either does or doesn't satisfy a written constraint that itself derives from a written business demand.

The list also makes the chain of evidence inspectable. Every TR cites BR-NN. Every later ADR cites TR-NN. Every line of the eventual tech-design traces BR → TR → ADR → component. If that chain breaks, the design has unjustified pieces. So this document is the foundation everything technical stands on; getting the discipline right here pays back in every later stage.

## Preconditions — read everything first

Before eliciting anything:

1. **Find the parent capability** at `docs/content/capabilities/{name}/_index.md`. If it is still a flat file (`{name}.md`), stop and ask the user to migrate it to page-bundle form first (see `define-user-experience`).
2. **Read `business-requirements.md`** at `docs/content/capabilities/{name}/business-requirements.md`. This is the authoritative input — every TR must cite the BR-NN it derives from. **Refuse to proceed** if this file is missing (route the user to `extract-business-requirements`) or if its `reviewed_at` frontmatter is `null` or older than the file's last modification time (the BRs haven't been human-reviewed yet — TRs derived from un-reviewed BRs are TRs with un-reviewed reasons).
3. **Read the capability doc end-to-end.** Internalize stakeholders, business rules, success criteria, out-of-scope list — used as context to validate that the BRs cover what they should.
4. **Read every UX doc** under `docs/content/capabilities/{name}/user-experiences/`. Not just the one the user mentioned — used as context for understanding why a BR exists, and for citing UX sections in TRs where the technical translation is shaped by a specific user-perceived behavior.
5. **Skim `docs/content/r&d/adrs/`** for prior shared decisions that constrain this capability (cloud provider, network topology, error response format, identifier standard, etc.). Cite them as constraints, don't re-decide them.
6. **Note the repo's house patterns** from `CLAUDE.md`: chi/bedrock Go service shape, `pkg/errorpb` for errors, no humus framework, Cloudflare → GCP topology. These are inherited constraints, not requirements you discover.

If `business-requirements.md` is missing or unreviewed, **stop and route the user to `extract-business-requirements`** (or to setting `reviewed_at` if BRs exist but haven't been signed off). If the capability or any UX docs are missing, route to `define-capability` / `define-user-experience` first. Tech requirements derived from missing inputs are tech requirements with missing reasons. Refuse to proceed; do not invent the missing inputs in your head.

## Goal

Produce or update `docs/content/capabilities/{name}/tech-requirements.md` from `assets/template.md`.

## What is and is not a technical requirement

A **TR is the technical translation of a BR** — the constraint(s) the technical solution must satisfy in order to deliver the business demand. Examples:
- BR: "An evicted tenant must be able to leave with all of their data." → TR: "The system must allow a tenant to retrieve a complete data export within a defined export window, without operator assistance."
- BR: "Existing tenants must keep working when the platform contract evolves." → TR: "The system must support N concurrent contract versions for a bounded migration window, with version pinning per tenant."
- BR: "No tenant can ever observe another tenant's state." → TR: "Per-tenant compute, storage, and network paths must be isolated such that a tenant has no readable surface onto another tenant's resources under normal or degraded operation."

A **TR can also derive from a prior shared ADR or repo-wide constraint**, not a BR. Examples:
- "All inter-service calls must traverse the Cloudflare → GCP path" (inherited from the repo's topology, not a BR).
- "Resource identifiers must follow the [ADR-0006 standard]({{< ref "/r&d/adrs/0006-resource-identifier-standard" >}})" (cite the ADR as the source).

A **decision is chosen** from multiple options that all satisfy a TR. Examples:
- "Use Postgres logical replication" (one of several ways to satisfy a no-downtime data-export TR).
- "Use mTLS for service-to-service auth" (one of several ways to satisfy a tenant-isolation TR).

**Decisions do not go in `tech-requirements.md`.** If the user volunteers a decision during extraction, capture it in **Open Questions** for later ADR work — never as a TR. This separation is what makes the eventual ADRs meaningful: an ADR's job is to pick one option among several that all satisfy the underlying TR. If you let a decision in here, you erase the alternatives that the ADR was supposed to weigh.

**A TR is not a restated BR.** If the TR text is byte-equivalent to the BR, you have not translated anything. The BR says *what the business demands*; the TR says *what the technical solution is forced to do* in order to deliver it (mechanisms, thresholds, surfaces — without choosing the technology).

## Append-only TR identity

Requirements are identified `TR-01`, `TR-02`, … and **numbers are append-only forever**. When re-extracting on a living doc:

- Preserve every existing TR-NN whose source link still resolves.
- Append newly-discovered requirements at the end with the next free number.
- If a TR's source no longer resolves (UX deleted, capability rule rewritten), **flag it** with `> ⚠️ source no longer resolves — human review` — do not delete it. The human resolves the flag.
- Never renumber. Gaps are honest history. Downstream ADRs cite TR-NN, so renumbering silently breaks ADR provenance.

## Source links

Every TR must cite at least one source. **The primary source is the BR-NN it derives from.** TRs may also cite the capability or UX section that frames the technical translation, or a prior shared ADR / repo pattern that constrains it. **Use Hugo's `ref` shortcode for every internal link — never raw paths.** Hugo will fail the build on a broken `ref`; raw paths break silently when content is reorganized.

- BR (primary): `[BR-03]({{< ref "business-requirements.md#br-03" >}})` — anchor must be the slugified BR heading from `business-requirements.md`. Add an explicit `{#br-03}` anchor on the BR heading there if one isn't present yet.
- Capability section (context): `[Capability §Business Rules]({{< ref "_index.md#business-rules" >}})`
- UX page or section (context): `[UX: upload-photo §Edge Cases]({{< ref "user-experiences/upload-photo.md#edge-cases" >}})`
- Prior shared ADR: `[ADR-0006]({{< ref "/r&d/adrs/0006-resource-identifier-standard" >}})`
- Repo pattern from CLAUDE.md: cite inline (CLAUDE.md is not Hugo content; do not link with `ref`).

**Section deep-links require an explicit anchor on the target heading.** Add `{#anchor-id}` to the heading you are linking to (e.g. `## Business Rules {#business-rules}` or `### BR-03: …  {#br-03}`) before linking — Hugo's default slugify-from-heading-text breaks every implicit anchor as soon as a heading is reworded. If you cite a section that has no explicit anchor yet, pause the extraction and either add the anchor in the source doc or capture the missing-annotation as an open question.

A TR with **no BR citation** must either cite a prior shared ADR / inherited constraint, or be flagged for review as an unsourced TR — it likely means a BR is missing and should be raised back to `extract-business-requirements`. Multi-sourced TRs are usually the most important ones — they are the ones forced from more than one direction.

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

- **Translate, don't design.** No options, no chosen approaches, no technology names. If you find yourself proposing a database, a protocol, or a library, stop — that is the next skill's job.
- **Walk the BR list.** The primary input is `business-requirements.md`. For each BR, ask: what does the technical solution have to do — in implementation-neutral terms — to deliver this? That answer (one or more sentences) is the TR. A BR may produce zero TRs (already covered by inherited constraints), one TR, or several TRs that decompose the demand.
- **Quote the source where possible.** If a BR says "uploads must survive intermittent connectivity without losing user-perceived progress", a derived TR might be "the system must persist upload progress to durable per-tenant storage and resume from the last persisted offset on reconnect" — and the link from TR-NN to BR-NN makes the derivation auditable.
- **Push back on premature solutions.** If the user says "we'll need Postgres for tenant state", redirect: "That's a decision for `plan-adrs` / `define-adr`. The TR underneath is what?"
- **Push back on missing BRs.** If the user volunteers a TR that doesn't trace to any BR (and isn't an inherited shared-ADR / repo-pattern constraint), the right answer is usually that a BR is missing. Surface it: "This TR has no BR to derive from — should we add a BR to `business-requirements.md` first via `extract-business-requirements`, or is this an inherited constraint?"
- **Don't invent TRs.** If nothing in the BR list, capability, UX, or inherited constraints implies a thing, it isn't a TR. Capture it as an open question.
- **Living-doc framing.** Tell the user the doc is meant to drift as BRs and UX docs evolve. Re-extract on demand. Git diff is the review surface — don't lament that the doc changed; that's the point.

## Producing the document

Use `assets/template.md`. Fill `{{capability_name}}`, `{{requirements}}` (the numbered TR sections), and `{{open_questions}}`. Each TR follows the shape commented in the template — heading, **Source:** line(s), **Requirement:** paragraph, **Why this is a requirement, not a decision:** line.

Save to `docs/content/capabilities/{name}/tech-requirements.md` (page-bundle form for the parent capability is required).

Keep the Hugo/Docsy frontmatter intact — it is required for the page to render in the published docs site. Set `reviewed_at: null` on a fresh extraction or when re-extraction has changed any TR; only carry forward an existing `reviewed_at` value if the file's TR contents are byte-equal to what you started with.

After producing the doc, tell the user:
- Where it was saved
- How many TRs were added/preserved/flagged
- Which open questions are parked for the ADR stage
- That `plan-adrs` is gated on `reviewed_at` and won't run until they set it

## Validate the rendered page

After saving (or re-extracting) `tech-requirements.md`, run `cd docs && hugo server` and open the page in a browser to confirm it renders correctly — every TR's `ref` source link must resolve, and any flagged stale-source markers must be visible. A clean build is necessary but not sufficient; pages can build clean while rendering broken (missing/incorrect front-matter `type`, unresolved `ref` shortcode, broken section anchor). The doc is not done until the rendered page has been visually validated. If you cannot run the preview in this session, say so explicitly to the user rather than reporting the doc as done.

## Output format

Plain markdown matching `assets/template.md`. No code fences around the whole document. Mermaid is unnecessary here; this stage is text-only.

## Examples of good redirection

**Example 1 — premature decision dressed as a TR:**
User: "TR-04 should be: use Postgres for tenant metadata storage with logical replication for cross-region failover."
Response: "That's a Stage-2 decision, not a TR. 'Use Postgres' is one of several options that could satisfy the underlying constraint. The TR is closer to: 'tenant metadata must be queryable by tenant ID with strong consistency, and survive a single-region failure within {N} minutes' — which BR forces this? I'll cite that BR and capture 'consider Postgres + logical replication' as an open question for the corresponding ADR."

**Example 2 — TR with no BR to derive from:**
User: "Add TR-08: API responses must include a request-id header."
Response: "I don't see a BR that forces this. Is it derived from a prior shared ADR (in which case I'll cite the ADR), or is there a missing BR — e.g. 'tenants must be able to correlate a failed request with the platform's logs' — that we should add via `extract-business-requirements` first? Either is fine; I just don't want to create an unsourced TR."

**Example 3 — vague slogan as TR:**
User: "TR-N: the system must be reliable."
Response: "'Reliable' is a slogan, not a TR until it's measurable and sourced to a BR. Which BR demands the reliability outcome — and what does 'reliable' mean in its terms (availability target, max recovery time, guaranteed durability)? With the BR cited and the constraint quantified, we can phrase the TR around it."

**Example 4 — re-extraction surfaces a stale link:**
Response: "TR-06's BR-link `business-requirements.md#br-04` no longer resolves — BR-04 was renumbered or removed. I've flagged it with ⚠️ rather than rewriting or deleting it. Want to re-source it now, or leave the flag for review? Either way, I'll null `reviewed_at` since the file changed; you'll need to re-set it before `plan-adrs` will run."

**Example 5 — BRs not yet reviewed:**
User: "Extract tech requirements for the self-hosted-application-platform capability."
Response: "I can't proceed yet — `business-requirements.md` exists but its `reviewed_at` is `null` (or older than the file's last modification). TRs derived from un-reviewed BRs are TRs with un-reviewed reasons. Please review and edit the BR doc, set `reviewed_at` to today's ISO date, and I'll pick up from there."
