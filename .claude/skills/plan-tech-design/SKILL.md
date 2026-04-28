---
name: plan-tech-design
description: >-
  Compose a skeleton tech-design document for a capability from its accepted
  ADRs and file one GitHub issue per component (for `define-component-design`)
  and one issue per surfaced gap. The output is `tech-design.md` with overview,
  component inventory, and a TR → ADR → component audit trail — NOT the
  per-component detail (those live in component design docs filed under
  `define-component-design`). Use this skill whenever the user wants to
  compose, plan, or scope the tech design for a capability — phrases like
  "compose the tech design for {capability}", "plan the components", "what
  components does this capability need", "write the tech-design.md", or as
  the step that follows accepting all ADRs. Do NOT use to author per-component
  design docs (use `define-component-design`). Do NOT use to draft ADRs (use
  `define-adr`).
---

# Plan the Tech Design for a Capability

This skill composes `tech-design.md` for a capability — a **skeleton** narrative that names the components, maps every TR-NN through an accepted ADR to a component, and surfaces gaps. It does **not** write the per-component detail (sequence flows, schemas, API contracts). Each component gets a GitHub issue so `define-component-design` can author its design doc one at a time.

This is **Step 9** of the capability development lifecycle. It runs once all the ADRs in `adrs/` are accepted, and produces both an artifact (`tech-design.md`) and a manifest of follow-up work (component issues + gap issues).

## Why this matters

A tech design is the seam between architecture and implementation. The architecture (ADRs) says *what was decided and why*; the implementation (component designs) says *exactly what we will build*. The tech-design.md is what makes that seam navigable: it lists the components a future engineer will need to read about, and it carries the audit trail showing every requirement traces through a decision into a component.

The skeleton-only scope of this skill is intentional. Per-component detail belongs in per-component docs because component formats differ — a database table needs columns and indexes, an API service needs endpoints and contracts, a Terraform module needs inputs and outputs. Bundling them into one giant tech-design.md hides the differences and makes the doc unrebuildable when one component changes. Splitting them out also makes the work parallelizable across humans and across sessions.

The hard discipline of this skill is **gap-surfacing**. While composing the skeleton, missing pieces become visible — an API contract that no ADR specifies, a database schema that no component owns, a flow that two ADRs would split inconsistently. These are not failures; they are the value of composition. The skill's job is to file them as issues (or resolve them inline) before declaring the tech design complete.

## Preconditions — refuse to run without them

Before composing anything:

1. **Find the capability** at `docs/content/capabilities/{name}/_index.md` (page-bundle form).
2. **Read the capability doc**, every UX doc, and `tech-requirements.md` end-to-end.
3. **Check the review gate on tech-requirements.** `reviewed_at:` must be an ISO date *newer* than the file's last modification time. If `reviewed_at: null` or older, **stop**.
4. **Read every ADR under `adrs/`.** Verify each has `status: accepted` (or `superseded` with the superseder accepted). If any ADR is `proposed`, **stop and list them**:

   > "Cannot compose tech-design.md — these ADRs are still proposed: {list}. Either accept them via `define-adr`, or revise them. I won't compose a design against unsettled decisions."

5. **Read `docs/content/r&d/adrs/`** for prior shared decisions referenced by capability ADRs (cloud provider, network topology, error format, identifier standard). The composed design must respect them.
6. **Read `CLAUDE.md`** for repo house patterns. Components proposed must fit `services/{name}/`, `cloud/{module}/`, `pkg/` shape unless an ADR justifies otherwise.

## Goal

Produce two things:

1. **`docs/content/capabilities/{name}/tech-design.md`** from `assets/template.md` with: overview paragraph, component diagram, component inventory (one entry per component, linking to the issue filed below), data & state summary, and the TR → ADR → component audit trail.
2. **Two sets of GitHub issues, after the human approves both:**
   - One issue per component, filed against `define-component-design` (Step 10), with title `story(component): {component name} — {capability-name}`.
   - One issue per surfaced gap, with title `story(gap): {gap name} — {capability-name}`. Gap issues block the tech design from being considered complete.

Tech-design.md is **not complete** until every gap issue is resolved. Tell the user this explicitly at the end.

## Step 1 — Identify components from ADRs

Walk every accepted ADR's `Realization` section. Each names one or more components — a service, a Terraform module, a `pkg/` package, an external system. Cluster them:

- **Services** — `services/{name}/` Go applications.
- **Modules** — `cloud/{module}/` Terraform reusable modules.
- **Packages** — `pkg/{name}/` shared Go libraries.
- **External systems** — Cloudflare, GCP-hosted services, GitHub, etc.

Components may appear in multiple ADRs (e.g. a `tenant-registry` service might be established by ADR-0001 and used by ADR-0003). That's fine — the inventory entry cites every ADR it derives from.

If a component appears in no ADR but feels obvious (e.g. "we'll need an `_index.md` for the section"), that's a sign it's *not* an architectural decision — it's a structural one. Mention it in the overview but don't file a component issue for it.

## Step 2 — Surface gaps before composing

While walking the ADRs and TRs, look for:

- **TRs without ADRs.** Every TR must trace through some ADR. If TR-08 has no ADR addressing it, that's a gap — either the ADR set is incomplete (file a follow-up to `plan-adrs`) or TR-08 was abandoned (return to `define-technical-requirements` to remove or amend).
- **ADRs without realizations.** An accepted ADR with no `Realization` section, or with a vacuous one, has no component to put it into. The ADR may need amending.
- **Implementation-detail gaps that ADRs assume but don't specify.** Examples: an ADR says "tenants are addressed by tenant ID" but no shared ADR defines the tenant ID format; an ADR says "the tenant API exposes CRUD over tenants" but no contract names the actual endpoints; an ADR says "tenant data is stored as protobuf-encoded records" but no schema defines the records. **These are *not* ADR-worthy** (deciding to call an API is architectural; deciding what API to call is implementation), but they must be specified before a component design can be written.
- **Two ADRs in tension.** ADR-0003's component depends on a behavior ADR-0005 has retired, etc. These need amending ADRs, not papering-over prose.

For each gap: state it specifically, name what would resolve it (a follow-up `define-technical-requirements` run, an amending ADR via `define-adr`, an immediate inline answer from the user, or a per-component spec deferred to `define-component-design`), and add it to the gap list. **Do not invent solutions in `tech-design.md`.**

## Step 3 — Compose the skeleton

Use `assets/template.md`. Fill:

- `{{overview_paragraph}}` — one paragraph summarizing what gets built. No surprises here; this is a digest of the ADR outcomes.
- `{{component_mermaid}}` — Mermaid component diagram. Boxes for components, edges for the principal data/control flows derived from ADRs. Keep it digestible — at most ~10 boxes; if the system has more, group sub-components.
- `{{component_inventory}}` — one entry per component. Each entry: name, repo location, Established-by ADR list, one-sentence responsibility, and a link to the component-design issue filed in Step 5. **No deeper detail.** That belongs in the per-component design doc.
- `{{data_state}}` — one short paragraph naming what data exists, who owns it, lifecycle. Detailed schemas live in per-component docs; this is just enough for a reader to know which component to click into.
- `{{requirement_realization_table}}` — one row per TR-NN. Columns: TR, ADR(s), Realized in (component(s)). **Every TR must appear.** If a row would be empty in any column, it's a gap — return to Step 2.
- `{{deferred}}` — gaps the user explicitly chose to defer plus any out-of-scope items.

Do **not** write `## Key flows` content beyond a placeholder pointing at the component issues. Per-UX sequence diagrams are a per-component-design concern (or deserve their own follow-up issues if they cross multiple components).

## Step 4 — Mirror back, get explicit approval

State out loud, before any issues are filed:

- The component list (with one-sentence summaries).
- The gap list (one line per gap).
- The audit-trail row count vs. TR count (must match).

Ask:

> "I've composed the skeleton at `tech-design.md` and identified {N} components and {M} gaps. Before I file issues: do the components look right? Do you want any gap resolved inline now? Once you say go, I'll file {N} component issues for `define-component-design` and {M} gap issues."

Wait for explicit approval before filing.

## Step 5 — File issues

Once approved, file via `gh issue create`:

- **One component issue per component**, title `story(component): {component name} — {capability-name}`. Body includes: capability link, the ADR(s) that established this component, the responsibility, and a pointer to `define-component-design` as the authoring skill.
- **One gap issue per gap**, title `story(gap): {gap name} — {capability-name}`. Body includes: the gap statement, what type of resolution is needed (`define-technical-requirements`, amending `define-adr`, or per-component spec), and a link back to the parent capability planning issue.

Print the issue numbers/URLs back as a manifest.

### Component issue body template

```markdown
### Component

**Location:** `{repo-path}` (e.g. `services/tenant-registry/`)
**Type:** {service / module / package / external-system}
**Responsibility:** {one sentence}

### Established by

- [ADR-{NNNN}](../adrs/{NNNN}-{slug}.md) — {short title}
- [ADR-{NNNN}](../adrs/{NNNN}-{slug}.md) — {short title}

### Authoring

This component's design doc will be authored via `define-component-design` — one invocation per component. The doc format is type-specific (e.g. table definition vs. API service). The composed `tech-design.md` will be updated to link to the component design once written.

### Parent capability

[{capability-name}](../docs/content/capabilities/{name}/_index.md)

### Related

#{parent-capability-issue-or-722}
```

### Gap issue body template

```markdown
### Gap

{One paragraph describing the gap and why it blocks tech-design completeness.}

### Type of resolution needed

- [ ] Amending TR via `define-technical-requirements`
- [ ] Amending ADR via `define-adr`
- [ ] Per-component spec via `define-component-design`
- [ ] Inline answer (resolved in conversation)

### Surfaced during

Composition of `tech-design.md` for [{capability-name}](../docs/content/capabilities/{name}/_index.md).

### Related

#{parent-capability-issue-or-722}
```

## Completion

After filing, tell the user explicitly:

> "I filed {N} component issues and {M} gap issues. The tech design is **not complete** until all gap issues are resolved — `plan-implementation` (Step 11) won't run against an incomplete tech design. Each component issue is the input to one `define-component-design` invocation. Run them in any order; the audit trail in `tech-design.md` will be updated as components are designed."

## Conversation discipline

- **Announce gates before doing anything else.** "Tech-requirements is reviewed; ADRs 0001..0007 all accepted. Composing." Or "Stopping — ADR-0003 is still proposed."
- **Don't write per-component detail.** If you find yourself writing endpoint paths, table columns, or sequence diagrams beyond the audit-trail level, you're doing `define-component-design`'s job. Stop and file the issue.
- **Don't paper over gaps.** Prose is not a substitute for a missing decision or a missing spec. Surface, don't smooth.
- **The audit trail is non-negotiable.** Every TR appears in the table; every ADR with `status: accepted` is realized in some component. If you cannot make these true, return to `plan-adrs` / `define-adr`.

## Validate the rendered page

After saving `tech-design.md` (and before filing component/gap issues), run `cd docs && hugo server` and open the page in a browser to confirm it renders correctly — every `ref` link to TRs, ADRs, and UXs must resolve, the component Mermaid diagram must render, and the audit-trail table must format cleanly. A clean build is necessary but not sufficient; pages can build with no errors yet render broken (missing/incorrect front-matter `type`, unresolved `ref` shortcode, malformed Mermaid). The doc is not done until the rendered page has been visually validated. If you cannot run the preview in this session, say so explicitly to the user rather than reporting the doc as done.

## Output format

`tech-design.md` is plain markdown with Hugo/Docsy frontmatter. Mermaid renders natively in Docsy. Component-inventory and audit-trail tables use standard markdown tables. No code fences around the whole document.

GitHub issues are filed via `gh issue create` after explicit user approval, one per component and one per gap.

## Examples of good behavior

**Example 1 — gap surfaced during composition:**
Response: "Walking ADR-0003: it says tenants are addressed by tenant ID. The shared ADR `r&d/adrs/0006-resource-identifier-standard` defines an ID format, but neither ADR specifies how the tenant ID is derived from the GitHub onboarding issue (auto-incremented? hash of the capability name? operator chooses?). That's an implementation gap — not ADR-worthy, but a per-component spec is needed. Filing a gap issue, type 'per-component spec via `define-component-design`'."

**Example 2 — premature per-component detail caught:**
User: "While you're at it, write out the tenant-registry service's API endpoints — I know what they should be."
Response: "Hold on — that's `define-component-design`'s job. I'll file a component issue for tenant-registry now and you can run `define-component-design` against it next; it'll capture endpoint design with the right API-service format. If I write endpoints into `tech-design.md`, it bypasses the per-component review and makes this doc unrebuildable when the API changes."

**Example 3 — uncovered TR detected:**
Response: "Building the audit trail surfaced a gap: TR-08 ('graceful degradation when a GCP region is unreachable') is in `tech-requirements.md` but no ADR addresses it. I'm not going to invent a story in `tech-design.md`. Two paths: (a) accept that TR-08 is unaddressed and file an ADR issue for `plan-adrs`/`define-adr` to handle later, or (b) drop TR-08 from `tech-requirements.md` if the constraint isn't real. Which?"
