---
name: define-adr
description: >-
  Author a single capability-scoped Architecture Decision Record (ADR) for a
  specific decision identified by `plan-adrs`. Identifies research tasks
  before proposing options, requires ≥2 options each tied back to TR-NNs,
  surfaces the trade-offs, and waits for the human to make the final
  selection — the skill never picks the option on its own. Use this skill
  whenever the user wants to write, draft, debate, or decide a single ADR
  for a capability — phrases like "draft the ADR for {decision}", "let's do
  ADR {NNNN}", "weigh options for {decision}", "write the ADR on tenant
  state storage", or as the per-decision step that follows `plan-adrs`. Do
  NOT use to plan or enumerate multiple ADRs (use `plan-adrs`). Do NOT use
  for shared cross-capability decisions (those belong in
  `docs/content/r&d/adrs/` via a separate flow).
---

# Define a Single ADR for a Capability

This skill authors **one ADR** for **one decision**. The ADR lives at `docs/content/capabilities/{name}/adrs/{NNNN}-{kebab}.md`, follows MADR 4.0.0, and starts in `status: proposed`. It moves to `status: accepted` only when the human explicitly chooses one of the proposed options.

This is **Step 8** of the capability development lifecycle. It runs once per ADR-issue created by `plan-adrs` (Step 7). Multiple invocations stay in lockstep with multiple issues; each invocation produces exactly one file.

## Why this matters

An ADR's job is to make **why** auditable. Every later component design and every implementation task will cite this ADR (or one of its descendants), and a future engineer will read it to understand what was tried and what was rejected. If the file doesn't show the alternatives that were considered and weighed against the underlying TRs, it's not an ADR — it's just a record of an opinion.

The discipline that makes this work is harsh: at least two options, every option's pros/cons phrased in terms of the TR-NNs the decision addresses, and the human picks. Skipping any of those produces ADRs that are confidently wrong months later when someone tries to supersede them.

## Preconditions — refuse to run without them

Before drafting anything:

1. **Identify the decision and its issue.** The user names a decision (e.g. "tenant state storage") or points to a `story(adr): ...` issue from `plan-adrs`. If neither, ask which decision and route the user to `plan-adrs` if the list of decisions hasn't been enumerated yet.
2. **Find `tech-requirements.md`** at `docs/content/capabilities/{name}/tech-requirements.md`. If it doesn't exist, stop and route to `define-technical-requirements`.
3. **Check the review gate.** `reviewed_at:` must be an ISO date *newer* than the file's last modification time. If `reviewed_at: null` or older, **stop**. The current TRs haven't been reviewed; an ADR sourced from drafts is meaningless.
4. **Read every TR-NN that the decision addresses.** Internalize the constraints they impose. Each option you eventually consider has to be evaluated against these TRs by name.
5. **Read `docs/content/r&d/adrs/`** for prior shared decisions that constrain the option set (cloud provider, network topology, error response format, identifier standard). Cite them as constraints — do not re-decide them in this ADR.
6. **Read `CLAUDE.md`** for repo house patterns: chi/bedrock service shape, `pkg/errorpb` for errors, no humus framework, Cloudflare → GCP topology, protobuf request/response. Options that contradict these patterns are not viable unless the ADR explicitly justifies departing.
7. **Read sibling ADRs** in `docs/content/capabilities/{name}/adrs/` (if any). Newly-drafted ADRs must not contradict accepted siblings — surface and resolve any tension before drafting.

## Goal

Produce one ADR file from `assets/template.md` with:
- Locally-numbered `{NNNN}` (next free number under `adrs/`, padded to 4 digits).
- TR-NN citations in `Addresses requirements:` and woven through the option analysis.
- **At least two considered options**, each evaluated explicitly against the cited TRs.
- `status: proposed` until the human picks; `status: accepted` only after explicit human selection.

## Step 1 — Identify research before proposing options

ADR options are usually not equally well-known to the agent or the user. Before drafting, ask: *what would I need to know to genuinely weigh these options?* Examples of facts that often require research:

- **Current-state facts**: what does the existing system already do, what is currently provisioned, what's already in `cloud/` or `services/` that constrains the choice?
- **External facts**: how does a candidate technology actually behave under the relevant constraints (e.g. logical-replication reconfiguration time, GCP regional failover semantics, mTLS handshake cost on the target hardware)?
- **Cost/quota facts**: pricing, quotas, free-tier limits relevant to the option set.
- **Compatibility facts**: does this option fit with an existing prior shared ADR, or does it require a sibling ADR change?

If the decision can be made meaningfully without these, skip ahead. If not, surface the research tasks **before** drafting options:

> "Before I draft options for this ADR, three things need to be checked: (1) {fact}, (2) {fact}, (3) {fact}. Do you want me to research these now, file follow-up issues for them, or do you already know the answers? I won't propose options against unknowns."

The user can answer the research inline, defer it (you file an issue and pause), or instruct you to make the unknowns explicit assumptions in the ADR. **Do not silently invent answers** — fabricated factual claims in an ADR's Considered Options section are exactly the bugs the discipline exists to prevent.

## Step 2 — Draft options (≥2, each anchored in TRs)

Once research is resolved (or explicitly deferred), draft at least two options. For each:

- Name the option in implementation-neutral language ("per-tenant Postgres database" not "Cloudbase v3.7"; or, when a specific product matters, name it but explain *why this product specifically* and not just the category).
- **Pros and cons phrased in terms of cited TRs.** Not "this is faster" — "Option A satisfies TR-04 (no-downtime updates) because rolling-upgrade primitives exist, but partially fails TR-01 (isolation) since shared compute is required."
- Note any TR an option *fails* or *partially satisfies*. Failing options stay in the doc — you don't delete them; the ADR's value comes from showing what was rejected and why.
- Note any prior shared ADR each option assumes or contradicts.

If you cannot phrase an option's pros/cons in TR-NN terms, that's a signal the decision is either premature (an underlying TR is missing — return to `define-technical-requirements`) or out of scope for this ADR (the trade-offs are about something other than what this ADR is meant to decide).

## Step 3 — Mirror back, do not select

State the options out loud and ask the user to choose. Phrase it as a question, not a recommendation:

> "Three options are on the table. Option A favors {TR set}; Option B favors {TR set} but loses ground on {TR}; Option C is the most conservative but requires {follow-up}. Which do you want to accept, or do you want to revise the option set?"

If the user explicitly asks for your opinion, you may give one — but it must be phrased as preference, not selection, and must explain the trade-off in terms of the TRs. **Do not write `status: accepted` until the user has picked an option.**

## Step 4 — Write the file

When the user has chosen:

1. Use `assets/template.md`.
2. Fill the frontmatter: `number`, `title`, `one_liner`, `weight`, `category` (strategic / user-journey / api-design), `date` (today's ISO date), `status: accepted`.
3. Cite TR-NNs in `Addresses requirements:` (e.g. `TR-03, TR-07`).
4. Fill `Decision Drivers` — each driver should trace to a TR or to a constraint inherited from a prior shared ADR / CLAUDE.md pattern.
5. Fill `Considered Options` with all options drafted in Step 2 (chosen and rejected). Keep the rejected options' analysis intact.
6. Fill `Decision Outcome` — name the chosen option and explain the rationale **in terms of the cited TRs**.
7. Fill `Realization` — which `services/{name}/`, `cloud/{module}/`, `pkg/` packages this decision shows up in. (`plan-tech-design` will eventually compose these into the tech-design narrative.)
8. Move any unresolved sub-questions to `Open Questions`.

Save to `docs/content/capabilities/{name}/adrs/{NNNN}-{kebab-name}.md`. If `adrs/_index.md` doesn't exist, create it with Docsy section frontmatter — Hugo won't render the ADR list otherwise.

If the user is undecided after debate, write the file with `status: proposed`, capture the debate in `Open Questions`, and tell the user explicitly that downstream skills (`plan-tech-design`) refuse to compose the tech design until status is `accepted`.

## Flag-and-stop for shared decisions

If the decision is obviously cross-capability (touches Cloudflare topology, identity, networking, error response format, the resource identifier standard, etc.), do not draft it as a capability-scoped ADR. Surface it:

> "This decision is shared across capabilities — it touches {topic}. It belongs in `docs/content/r&d/adrs/` via a separate flow. I'm not going to write it as a capability-scoped ADR here. Want to defer this ADR issue and proceed with the next one in the list, or pause to handle the shared ADR separately?"

If `plan-adrs` filed an issue for what is actually a shared decision (catching it late), help the user reclassify the issue rather than producing a capability-scoped ADR that will need to be deleted.

## House-pattern discipline

The repo's `CLAUDE.md` constrains the option set for capability-scoped ADRs:
- Go services follow `services/{name}/main.go → app/ → endpoint/ → service/` shape with chi router and bedrock config-from-env.
- Errors use `pkg/errorpb` and the `application/problem+protobuf` content type.
- Request/response are protobuf over HTTP.
- The humus framework is **not** used in this repo — never propose `rest.Run`, `rpc.Producer`, `rpc.Handler`, or embedded `config.yaml` as the implementation.
- Inter-region traffic conforms to the Cloudflare → GCP topology with WireGuard back to home lab.

Options that depart from these patterns are allowed but require explicit ADR-level justification — the ADR must say "we are diverging from {CLAUDE.md pattern} because {TR justification}" and surface the consequences. An option that silently invents non-house patterns (e.g. "use the humus framework") is wrong, not just unconventional.

## Conversation discipline

- **Announce the gate result and the decision being addressed up front.** "Tech-requirements is reviewed. Drafting ADR-{NNNN} for {decision}, addressing TR-{X}, TR-{Y}."
- **Do not skip Step 1.** If the only reason an option looks better is that you don't know the costs of the other, that's a research gap, not a decision.
- **Two options is the floor, not the goal.** If three or four are real, draft three or four. Don't force a straw option in to hit a count — but if you are tempted to write only one option, the decision is probably not yet ADR-worthy.
- **The human picks.** Even when the answer feels obvious, mirror back and ask. Decisions made by the human against an honest option set are robust; decisions made by the agent are brittle.

## Validate the rendered page

After saving the ADR file (and any new `adrs/_index.md`), run `cd docs && hugo server` and open the ADR page in a browser to confirm it renders correctly — the `ref` link to the parent capability must resolve, the ADR must appear in the `adrs/` section list, and any embedded Mermaid must render. A clean build is necessary but not sufficient; pages can build with no errors yet render broken (missing/incorrect front-matter `type` or `status`, unresolved `ref` shortcode, malformed Mermaid). The doc is not done until the rendered page has been visually validated. If you cannot run the preview in this session, say so explicitly to the user rather than reporting the doc as done.

## Output format

One MADR 4.0.0 markdown file from `assets/template.md`, with Hugo/Docsy frontmatter intact. No code fences around the whole document. Mermaid is allowed inside the body if a diagram clarifies an option.

## Examples of good behavior

**Example 1 — research surfaced before options:**
User: "Draft the ADR for tenant state storage."
Response: "Before I propose options, three things need to be answered: (a) is there an existing Firestore project we'd reuse or are we provisioning fresh? (b) what's the realistic logical-replication-aware failover time for the GCP regions in scope? (c) does the prior shared ADR on resource identifiers constrain table/key naming? Want me to research them, file issues for them, or do you already know? I'm not going to compare options against guesses."

**Example 2 — house-pattern violation caught:**
User: "Option B should be: build the tenant-control-plane service using the humus framework's `rest.Run` so we get retries for free."
Response: "humus isn't used in this repo per `CLAUDE.md` — the existing pattern is `chi.Mux` + `bedrock` config + protobuf endpoints. If you want retries, the option to draft is something like 'wrap chi handlers with a retry middleware' or 'add a retry concern at the client side', and the ADR can call out the missing primitive as a follow-up. I won't draft Option B as written."

**Example 3 — decision-by-agent refused:**
User: "Just pick whichever option is best. I trust you."
Response: "I'll give you my preference, but I can't set `status: accepted` on your behalf — the discipline of this skill is that the human picks. My preference is Option B because it satisfies TR-01 and TR-04 cleanly while accepting a known cost on TR-06 (which is partially mitigated by the import retry semantics in the migration UX). Want to accept Option B, want me to surface anything I might have missed, or want to revise the option set?"
