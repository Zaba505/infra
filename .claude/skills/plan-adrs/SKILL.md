---
name: plan-adrs
description: >-
  Guide the user through enumerating the architectural decisions a capability's
  technical requirements force, then file one GitHub issue per ADR so each can
  be authored individually via `define-adr`. The output is a confirmed list of
  decisions and a set of issues — not the ADRs themselves. Use this skill
  whenever the user wants to plan, enumerate, or scope the ADRs needed for a
  capability — phrases like "plan ADRs for {capability}", "what ADRs do we
  need", "list the architectural decisions", "break this into ADRs", or as
  the step that follows reviewing tech-requirements. Do NOT use to draft or
  decide individual ADRs (use `define-adr` for that). Do NOT use to add new
  technical requirements (use `define-technical-requirements`).
---

# Plan ADRs for a Capability

This skill turns a reviewed `tech-requirements.md` into a confirmed list of architectural decisions and files one GitHub issue per ADR. It does **not** draft or decide any ADR — it only plans the set so that each can later be tackled one-at-a-time by `define-adr`.

This is **Step 7** of the capability development lifecycle. It sits between `define-technical-requirements` (Step 6) and `define-adr` (Step 8). The "plan, then per-item loop" pattern is intentional: enumerating the full set first surfaces decision dependencies and prevents the agent from drafting an ADR that will need to be redrafted once a sibling decision is made.

## Why this matters

A capability typically forces several architectural decisions, and they are rarely independent. "Where does tenant state live" implies "how is it partitioned" implies "how do services access it." If you draft those one at a time without first listing them, you write the first ADR with assumptions that the second ADR will violate, and you end up rewriting the first to match. Listing them up front makes those dependencies visible and lets the human decide the order.

Filing one issue per ADR also means the work is parallelizable across humans and across sessions: a partial standalone ADR is more useful than half of a bundled one, and an issue per decision creates a checkable manifest of "what's left."

## Preconditions — refuse to run without them

Before enumerating anything:

1. **Find `tech-requirements.md`** at `docs/content/capabilities/{name}/tech-requirements.md`. If it doesn't exist, stop and route to `define-technical-requirements`.
2. **Read it end-to-end.** You need every TR-NN in working memory to enumerate decisions against them.
3. **Check the review gate.** The frontmatter must have `reviewed_at:` set to an ISO date *newer* than the file's last modification time (`git log -1 --format=%aI -- tech-requirements.md` or `stat`-equivalent). If `reviewed_at: null` or older than the last edit, **stop**. The human has not reviewed the current contents; planning ADRs against an unreviewed list of TRs produces decisions sourced from drafts.

   Tell the user explicitly:

   > "I won't enumerate ADRs yet — `tech-requirements.md` shows `reviewed_at: {value}` but the file was last modified {when}. Review the current contents and set `reviewed_at:` to today's ISO date (or tell me you've reviewed and I'll record your verbal confirmation), then re-invoke me."

4. **Skim `docs/content/r&d/adrs/`** for prior shared decisions. They constrain what the capability decides itself; they are not re-decided here.
5. **Note the repo's house patterns** from `CLAUDE.md` (chi/bedrock service shape, `pkg/errorpb`, no humus, Cloudflare → GCP topology). These bound the option space for capability-scoped ADRs.

## Goal

Produce two things, in this order:

1. **A list of proposed ADRs** to confirm with the human. Each item has a *short title* and a *one-liner* explaining the decision and which TR-NNs it addresses. No options yet, no decision rationale yet — that is `define-adr`'s job.
2. **One GitHub issue per ADR**, filed only after the human approves the list, using `gh issue create`. Each issue body includes the parent capability link, the TR-NNs the ADR will address, and a short statement of the decision to be made.

## What is and is not an ADR-worthy decision

An **ADR-worthy decision** is one where:
- Multiple plausible options exist, and
- The choice has consequences that ripple beyond a single component, and
- The choice ties back to one or more TR-NNs (so the rationale is auditable).

Examples (capability-scoped):
- "Where tenant state is stored" — addresses TR-01 (isolation) and TR-04 (no-downtime updates)
- "Tenant identifier scheme" — addresses TR-05 (per-tenant exports) and TR-03 (per-tenant observability)

**Not an ADR:** an implementation detail with one reasonable choice and no cross-component impact (e.g. "what JSON field name to use for tenant ID in this one endpoint"). Those belong in component design, not ADRs.

**Not a capability-scoped ADR:** decisions that ripple across capabilities (cloud provider, network topology, error response format, identifier standard). See "Flag-and-stop for shared decisions" below.

## Enumeration discipline

- **Read the TRs first, then propose the decision set.** Don't go TR-by-TR mechanically — multiple TRs typically motivate one decision, and one decision can address multiple TRs. Cluster naturally.
- **One decision per ADR.** Bundling — "Tenant data architecture" covering storage *and* partitioning *and* access pattern — produces ADRs that nobody can supersede without revisiting unrelated subdecisions. Resist it. If the user proposes a bundle, split it explicitly and explain why:

  > "Let's split that. I see at least three decisions: where tenant state lives, how it's partitioned, and how services access it. Bundling them dilutes the rationale for each, and we can't supersede one without revisiting all of them. Which do you want first?"

- **Mirror back before filing issues.** State the proposed list aloud; let the user add, remove, reorder, or merge. Only file issues once the user has approved the list.
- **Don't propose options yet.** If you find yourself starting to weigh "Postgres vs. Firestore," stop — that's `define-adr`. The list at this stage is decisions to be made, not the answers.

## Flag-and-stop for shared decisions

If a proposed decision is obviously cross-capability (touches Cloudflare topology, identity, networking, error response format, the resource identifier standard, etc.), do not include it in the capability-scoped list. Surface it:

> "This decision looks shared across capabilities — it touches {topic}. It belongs as a shared ADR in `docs/content/r&d/adrs/`, which is a different flow than this skill. Want to defer it and proceed with the capability-scoped list, or pause to handle it separately?"

The shared-ADR flow is out of scope here. Do not draft or file capability-scoped issues for shared decisions.

## Filing the issues

Once the user approves the list, file one GitHub issue per ADR via `gh issue create`. Each issue:

- **Title:** `story(adr): {short ADR title} — {capability-name}` (matches the repo's `story(scope): description` convention).
- **Body:** parent capability link, the TR-NNs to be addressed, a one-paragraph statement of the decision to be made, and a link back to the parent capability planning issue if one exists.
- **Body must reference `define-adr`** as the skill that will author the ADR and explain that one invocation of `define-adr` produces one ADR file.

After filing, print the issue numbers/URLs back to the user as a manifest.

### Issue body template

```markdown
### Decision

{One paragraph: what decision needs to be made. Phrased as a question or as a missing choice, not a chosen answer.}

### Technical requirements addressed

- [TR-{NN}](../tech-requirements.md#tr-{NN}-{slug}) — {short label}
- [TR-{NN}](../tech-requirements.md#tr-{NN}-{slug}) — {short label}

### Parent capability

[{capability-name}](../docs/content/capabilities/{name}/_index.md)

### Authoring

This ADR will be authored via the `define-adr` skill — one invocation per ADR. The skill will identify research tasks, propose options tied back to the TR-NNs above, and stop for the human to make the final selection.

### Related

#{parent-capability-issue-or-722}
```

## Conversation discipline

- **Announce the gate result before doing anything else.** "Tech-requirements is reviewed (date `{X}`, file unchanged since) — proceeding." Or "Not yet reviewed — stopping." This makes the gate visible.
- **Read the TRs out loud before proposing.** Don't start by guessing the decision set; ground each proposed ADR in specific TR-NNs.
- **Stop at "approved list."** Until the user says "yes, file these," no issues are created. Don't pre-file in case the user changes their mind.
- **No options, no rationale, no draft ADRs in this skill.** Those are `define-adr`'s job. The line is sharp on purpose.

## Output format

Conversational user-facing output (the proposed list, the gate state, etc.) plus, on approval, GitHub issues filed via `gh issue create`. No markdown artifact written to the repo — the artifacts are the issues themselves.

## Examples of good redirection

**Example 1 — gate not satisfied:**
User: "Plan the ADRs for self-hosted-application-platform."
Response (after reading the file): "Stopping — `tech-requirements.md` has `reviewed_at: null` and was last modified today. The skill before this one (`define-technical-requirements`) instructs you to set `reviewed_at` to today's ISO date once you've read and edited the file. Do that, then re-invoke me. (If you want to confirm verbally that you've reviewed it, say so and I'll set the field on your behalf.)"

**Example 2 — bundle proposed:**
User: "Let's do one ADR called 'Tenant data architecture' that covers storage, partitioning, access, and backup."
Response: "Let's split. I count at least four decisions in there: (1) where tenant state lives, (2) how it's partitioned, (3) how services access it, (4) backup approach. Each addresses different TRs and any of them might be superseded later without touching the others. I'll file four issues — pick which order to tackle them in, or I'll suggest one based on dependencies. Want to keep going as four, or merge any?"

**Example 3 — shared decision surfaced:**
User: "Add an ADR for how all platform services authenticate to each other — let's go with mTLS at the Cloudflare edge."
Response: "That's a shared decision, not capability-scoped — it touches Cloudflare topology and inter-service identity, both of which affect every capability hosted on the platform. It belongs in `docs/content/r&d/adrs/` as a shared ADR via a separate flow, not in this capability's `adrs/`. Also, 'go with mTLS' is the answer; this skill only enumerates decisions to be made, not the chosen options. Want to defer the shared ADR and proceed with the capability-scoped list, or pause this skill to handle the shared one separately?"
