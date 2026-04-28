---
name: plan-implementation
description: >-
  Break a completed capability tech design into implementation tasks and file
  one GitHub issue per task so each can be developed individually via the
  per-task implementation flow (Step 12). The output is a confirmed task
  breakdown plus a set of issues — not the implementation itself, and not any
  per-task plan. Use this skill whenever the user wants to plan, scope, or
  enumerate the implementation work for a capability whose tech design and
  component designs are complete — phrases like "plan implementation for
  {capability}", "break this into tasks", "what work is left to build
  {capability}", "file the implementation issues", or as the step that follows
  finishing every component design. Do NOT use to write code, draft per-task
  plans, or compose tech-design / component-design docs (those belong to Step
  12 / `plan-tech-design` / `define-component-design`). Do NOT use against an
  incomplete tech design — refuse and route the user back.
---

# Plan Implementation Work for a Capability

This skill turns a **completed** capability tech design into a confirmed list of implementation tasks and files one GitHub issue per task. It does **not** write code, and it does **not** write the per-task development plan — that is Step 12's job, run once per task.

This is **Step 11** of the capability development lifecycle. It runs once `tech-design.md` is composed, every gap surfaced by `plan-tech-design` is resolved, and every component listed in the inventory has its component design doc written by `define-component-design`. The output is a manifest of dev work, sized so each task fits a focused pull request.

## Why this matters

A tech design plus a stack of component designs answers *what we will build*. Implementation tasks answer *how we will land the work, in what order, in what slices*. The slicing is its own discipline: too coarse and a single task hides multi-week work behind one issue (no review checkpoints, no parallelism, hard to estimate); too fine and the issue tracker drowns in trivia and the reader can't see the shape of the project. Both failure modes hide risk.

Doing this enumeration up front — before any code is written — surfaces ordering constraints (the registry table must exist before the API can be wired), missing prerequisites (no Terraform module exists for this kind of service yet), and ambiguity that the component designs left underspecified. It also produces a checkable manifest of issues so multiple humans (or sessions) can work in parallel.

The skill stops at the manifest. The per-task plan — what files to touch, what tests to write, how to verify — is `Step 12`'s job, authored one task at a time so each plan can be tailored to the kind of work (Go service vs. Terraform module vs. docs page).

## Preconditions — refuse to run without them

Before enumerating anything:

1. **Find the capability** at `docs/content/capabilities/{name}/_index.md` (page-bundle form). If a capability is still a flat file with no page bundle, stop — Step 11 doesn't run before Step 9 has produced a `tech-design.md` under a page bundle.
2. **Read `tech-design.md` end-to-end.** Every component name, every audit-trail row, the data & state summary, and the deferred list must be in working memory. Tasks are sliced against this doc.
3. **Check the review gate on tech-design.** The frontmatter must have `reviewed_at:` set to an ISO date *newer* than the file's last modification time. If `reviewed_at` is missing, `null`, or older than the last edit, **stop**.

   > "I won't plan implementation tasks yet — `tech-design.md` shows `reviewed_at: {value}` but the file was last modified {when}. Review the current contents and set `reviewed_at:` to today's ISO date (or tell me you've reviewed and I'll record your verbal confirmation), then re-invoke me."

4. **Check that every gap is resolved.** Search open issues:

   ```bash
   gh issue list --state open --search 'story(gap): in:title "{capability-name}" in:title' --limit 100
   ```

   If any open `story(gap):` issue mentions this capability, **stop and list them**:

   > "Cannot plan implementation — these gap issues are still open: {list}. The tech design isn't complete until each gap is resolved (via `define-technical-requirements`, an amending `define-adr`, or a per-component spec). Resolve them, then re-invoke me."

5. **Check that every component has a design doc.** Walk the component inventory in `tech-design.md`. For each entry, verify the design doc it links to exists on disk (typically `docs/content/capabilities/{name}/components/{component}.md`, but follow whatever path `tech-design.md` actually links to). If any component design is missing, **stop and list them**:

   > "Cannot plan implementation — these components in the inventory have no design doc yet: {list}. Each must be authored via `define-component-design` before tasks can be sliced against them. Run `define-component-design` for each, then re-invoke me."

   Also check open `story(component):` issues for this capability — an open component issue is a strong signal the design doc hasn't been written yet, even if a stub file exists.

6. **Read every component design end-to-end.** Tasks come from these docs (endpoints to wire, tables to provision, modules to author, jobs to schedule). Don't slice tasks from `tech-design.md` alone — its component entries are skeleton-only and won't surface the full implementation surface.
7. **Read `CLAUDE.md`** for repo house patterns (Go service shape under `services/{name}/`, Terraform modules under `cloud/{module}/`, shared Go libraries under `pkg/{name}/`, Hugo docs under `docs/content/`, branch and commit conventions). Tasks must respect these so each one can land as a normal PR.

## Goal

Produce two things, in this order:

1. **A list of proposed implementation tasks** to confirm with the human. Each task names the component(s) it touches, the work it does in one sentence, the acceptance criteria, and the design source (component design doc + TR(s)/ADR(s)).
2. **One GitHub issue per task**, filed only after the human approves the list, via `gh issue create`. Each issue body links the parent capability, names the component(s) and design source, and points at Step 12 (per-task development) as the next move.

## What is and is not a single implementation task

A **task-worthy slice** is one where:
- It can plausibly land as a single focused PR — small enough to review in one sitting, large enough that a reviewer can judge whether the change does something coherent.
- It has an explicit definition of done — a test passes, an endpoint responds, a module applies cleanly, a deployment succeeds.
- It traces to a specific component design (and through it to an ADR and one or more TRs).
- It does not depend on work that hasn't been sliced into a task yet — if a prerequisite exists, file it as its own task and order them.

Examples (from a hypothetical `tenant-registry` component with a CRUD HTTP API backed by Firestore):
- "Add `TenantRecord` Firestore schema and `service.tenants` client" — small, testable in isolation, traces to the component's data section.
- "Wire `POST /tenants` endpoint with validation and error handling" — depends on the schema task; one focused PR.
- "Add Terraform module `cloud/firestore-tenants` with collection + indexes" — separate concern, separate review, separate PR.

**Not a task:** "Build the tenant-registry service." That's the *component*. Slice it.

**Not a task:** "Refactor the entire `pkg/errorpb`." That's not implementation of this capability's design; if the component design doesn't call for it, it doesn't belong in this manifest.

**Not a task — too small:** "Rename one field in the proto." Either it's part of a coherent slice (fold it in) or it's a chore (skip the issue entirely; the dev will handle it during the parent task).

## Slicing discipline

- **Slice along seams the component designs already drew.** Each component design typically names its data, its surface (endpoints, inputs/outputs), its dependencies, and its operational concerns (deployment, monitoring). One task per seam is usually right.
- **Order tasks by hard dependencies, not by guessed priority.** If task B can't be built until task A is merged (e.g. an endpoint can't ship without its schema), say so explicitly so the reviewer can see the chain.
- **Don't fabricate work the design doesn't call for.** If a component design has no tests section, don't invent "add 100% coverage" as a task; the per-task plan in Step 12 will decide testing strategy. Tasks here are the *what*, not the *how*.
- **Don't slice every component identically.** A Terraform module is one or two tasks (author module + wire it from a root module). A Go service with five endpoints is more. Match the slice count to the actual surface.
- **One task per PR-shaped unit of work.** If you're describing a task and find yourself writing "and also" three times, split it.
- **Mirror back before filing issues.** State the proposed list aloud; let the user add, remove, reorder, rename, merge, or split. Only file once approved.
- **Don't write per-task plans here.** If you find yourself listing files to edit, tests to write, or CLI commands to run, stop — that's Step 12. The list here is *tasks to be planned and executed*, not their step-by-step content.

## Surface gaps before slicing

While reading the component designs, you may find seams the designs left underspecified — ambiguity that would force the implementor to invent a decision. If so, **don't paper it over with a task description**. Surface it:

> "Reading `components/tenant-registry.md`: the component design names a `POST /tenants` endpoint but doesn't specify whether tenant ID is generated server-side or supplied by the caller. That's an ambiguity I can't slice around. Either (a) update the component design to specify, or (b) tell me which way it should go and I'll record your call as a per-task input. Slicing waits."

Composition gaps that show up only when you try to enumerate tasks (cross-component coordination, deployment ordering, secret provisioning never named anywhere) are the same: surface them, route them to whoever should answer (an amending component design, a new gap issue, an inline answer), don't invent.

## Step 1 — Walk the components, list candidate tasks

For each component in the inventory:
- Read its design doc.
- List the discrete units of work the doc implies: schemas, endpoints, modules, scripts, deployment wiring, docs updates.
- Note which TR(s) and ADR(s) each unit traces back to (via the audit trail in `tech-design.md`).

Cluster across components where it makes sense — a task that touches two components but is one PR-sized change should be one task, not two. Don't artificially fan out.

## Step 2 — Order and present

Order tasks by hard dependency. Within a dependency level, order by what unblocks the most downstream work (data-layer tasks before API tasks before integration tasks is a common shape, but follow what the actual designs imply).

Mirror back the list with: short verb-led title, component(s) touched, one-sentence summary, acceptance criteria, design source, and (where it matters) the prerequisite task(s).

Ask:

> "I've sliced {N} implementation tasks across {M} components. Before I file issues: do the slices look right? Anything to merge, split, reorder, or drop? Any ambiguity in the component designs you want to resolve first? Once you say go, I'll file {N} issues."

Wait for explicit approval before filing.

## Step 3 — File issues

Once approved, file one GitHub issue per task via `gh issue create`. Each issue:

- **Title:** `story(impl): {short verb-led task title} — {capability-name}` (matches the repo's `story(scope): description` convention; `impl` is the implementation-task scope, parallel to `component`, `ux`, `gap`, etc.).
- **Body:** uses the template below.

After filing, print the issue numbers/URLs back as a manifest, in dependency order.

### Issue body template

```markdown
### Implementation task

{One-sentence summary of what this PR will do.}

### Component(s) touched

- [{component-name}](../docs/content/capabilities/{capability-name}/components/{component-name}.md)
- {... additional components if the task crosses seams}

### Design source

- Component design: see component link(s) above.
- ADR(s): [ADR-{NNNN}](../docs/content/capabilities/{capability-name}/adrs/{NNNN}-{slug}.md) — {short title}
- TR(s) realized: TR-{NN}{, TR-{NN}}

### Acceptance criteria

- [ ] {Concrete, observable signal the work is done — e.g., "POST /tenants returns 201 with the canonical record on valid input."}
- [ ] {Test signal — e.g., "Unit tests cover validation paths, integration test exercises the round-trip."}
- [ ] {Operational signal where relevant — e.g., "Module applies cleanly via `terraform fmt -recursive -check` and a dry-run plan."}

### Prerequisite tasks

- #{issue-number} — {prerequisite task title}, OR "None."

### Authoring

This task's per-PR development plan will be authored via the Step 12 flow — one invocation per task. The plan will identify files to touch, tests to write, and verification steps, tailored to the task type (Go service / Terraform module / docs / etc.).

### Parent capability

[{capability-name}](../docs/content/capabilities/{capability-name}/_index.md)

### Related

#{parent-capability-issue-or-722}
```

## Conversation discipline

- **Announce the gate result before doing anything else.** "Tech design reviewed `{date}`, file unchanged since. No open gap issues. All {N} components have design docs. Slicing." Or "Stopping — `components/tenant-registry.md` is missing." This makes the gate visible.
- **Read the component designs out loud (in summary) before slicing.** Don't slice by guessing what a component does — anchor each task in the actual doc.
- **Stop at "approved list."** Until the user says "yes, file these," no issues are created.
- **No code, no per-task plans, no file edits beyond the manifest.** Those are Step 12's job. The line is sharp on purpose.

## Output format

Conversational user-facing output (the proposed task list, the gate state, surfaced gaps) plus, on approval, GitHub issues filed via `gh issue create`. No markdown artifact written to the repo — the artifacts are the issues themselves. (Optionally, the manifest of issue URLs may be linked from a parent capability tracking issue, but this skill does not create or modify the parent.)

## Examples of good redirection

**Example 1 — gate not satisfied (open gap issue):**
User: "Plan implementation for self-hosted-personal-media-storage."
Response (after `gh issue list`): "Stopping — there's an open `story(gap): tenant-id derivation — self-hosted-personal-media-storage` (#812). The tech design isn't complete until that gap is resolved (per-component spec or amending ADR). Close that, then re-invoke me."

**Example 2 — missing component design:**
Response: "Stopping — the inventory in `tech-design.md` lists four components: tenant-registry, photo-store, share-service, firestore-tenants. I found design docs for the first three under `components/`, but `components/firestore-tenants.md` doesn't exist. There's also an open `story(component): firestore-tenants` issue (#824). Run `define-component-design` against #824, then re-invoke me."

**Example 3 — task too coarse pushed back:**
User: "Just file one issue: 'build the tenant-registry service.'"
Response: "That's the component, not a task — it'll hide multi-week work behind one PR with no review checkpoints. The component design names a Firestore-backed schema, five HTTP endpoints, a service-account wiring, and a deployment manifest. I'd slice that into roughly five tasks, ordered by dependency. Want me to mirror the proposed slice list?"

**Example 4 — ambiguity surfaced during slicing:**
Response: "While slicing tenant-registry: the component design names `POST /tenants` but doesn't say whether the tenant ID is server-generated or caller-supplied. I can't slice the validation work without that. Two paths: (a) update `components/tenant-registry.md` to specify, or (b) tell me your call now and I'll record it on the resulting `story(impl):` issue as a per-task input. Which?"

**Example 5 — premature implementation request:**
User: "Great, file the issues — and while you're at it, go ahead and start on the schema task; you've got everything you need."
Response: "Holding — Step 11 stops at the manifest. The schema task gets its own per-PR plan in Step 12 (authored via the per-task flow), tailored to a Go-service-with-Firestore-client task. Filing the issues now; once they're up, re-invoke against the schema issue and we'll plan that one task end-to-end."
