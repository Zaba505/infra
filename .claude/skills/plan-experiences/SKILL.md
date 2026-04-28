---
name: plan-experiences
description: >-
  Guide the user through enumerating the user experiences a defined business
  capability implies, then file one GitHub issue per experience so each can be
  authored individually via `define-user-experience`. The output is a confirmed
  list of experiences and a set of issues — not the UX docs themselves. Use
  this skill whenever the user wants to plan, enumerate, or scope the user
  experiences for a capability — phrases like "plan experiences for
  {capability}", "what user journeys does this capability need", "list the
  UXes", "break this capability into experiences", or as the step that follows
  reviewing a capability definition. Do NOT use to draft or define an
  individual user experience (use `define-user-experience` for that). Do NOT
  use to define a brand-new capability (use `define-capability` first).
---

# Plan User Experiences for a Capability

This skill turns a reviewed capability definition into a confirmed list of user experiences and files one GitHub issue per experience. It does **not** draft any UX — it only plans the set so each can later be tackled one-at-a-time by `define-user-experience`.

This is **Step 2** of the capability development lifecycle. It sits between `define-capability` (Step 1) and `define-user-experience` (Step 3). The "plan, then per-item loop" pattern is the same one used by `plan-adrs` and `plan-tech-design`: enumerating the full set first surfaces overlap and seams, prevents drafting a UX that needs to be redrafted once a sibling UX is defined, and gives the human a checkable manifest of what's left.

## Why this matters

A capability typically implies several distinct user journeys (e.g. *upload a photo*, *share an album*, *bulk import*, *delete content*). They are rarely independent — one journey's success state is often another's entry point, and one persona's UX often constrains a sibling persona's. Drafting them one at a time without first listing them produces UXes with mismatched seams: the share UX assumes content was uploaded a certain way, but the upload UX defined later doesn't match. Listing them up front makes those seams visible and lets the human decide priority.

Filing one issue per experience also parallelizes the work across humans and across sessions, and keeps `define-capability` focused on the capability itself rather than ballooning into a multi-journey enumeration.

## Preconditions — refuse to run without them

Before enumerating anything:

1. **Find the capability doc.** Capabilities live at `docs/content/capabilities/{name}.md` *or* `docs/content/capabilities/{name}/_index.md` (page-bundle form). If the user has not named the capability, ask which one. If no capability doc exists, stop and route to `define-capability`.
2. **Read it end-to-end.** You need the stakeholders, triggers, outputs, and business rules in working memory to enumerate experiences against them.
3. **Check the review gate.** The frontmatter must have `reviewed_at:` set to an ISO date *newer* than the file's last modification time (`git log -1 --format=%aI -- {capability-file}` or `stat`-equivalent). If `reviewed_at:` is missing, `null`, or older than the last edit, **stop**. Planning experiences against an unreviewed capability produces journeys sourced from a draft.

   Tell the user explicitly:

   > "I won't enumerate experiences yet — the capability doc shows `reviewed_at: {value}` but the file was last modified {when}. Review the current contents and set `reviewed_at:` to today's ISO date (or tell me you've reviewed and I'll record your verbal confirmation), then re-invoke me."

   If the capability doc has no `reviewed_at` field at all, treat it the same way: stop and ask the human to add it. The field is the gate; without it, there is no gate.

4. **Note the repo's house patterns** from `CLAUDE.md` (Hugo + Docsy docs site, capability docs as page bundles under `docs/content/capabilities/`, UX docs under `{capability}/user-experiences/`). These shape where the issue bodies will direct the next-step skill to save its output.

## Goal

Produce two things, in this order:

1. **A list of proposed user experiences** to confirm with the human. Each item has a *short verb-led title* (the journey name, e.g. "upload a photo") and a *one-liner* identifying the persona, the goal, and which capability sections motivate it (stakeholders, triggers, outputs, success criteria). No journey steps yet, no UI, no edge-case enumeration — that is `define-user-experience`'s job.
2. **One GitHub issue per experience**, filed only after the human approves the list, using `gh issue create`. Each issue body links back to the parent capability, names the persona and goal, and points at `define-user-experience` as the skill that will author the doc.

## What is and is not an experience-worthy journey

An **experience-worthy journey** is one where:
- One persona from the capability's stakeholder list is accomplishing one goal end-to-end, and
- The journey has a distinct entry point and a distinct success state, and
- It is grounded in the capability's triggers, outputs, or success criteria — so the rationale is auditable back to the capability doc.

Examples (capability-scoped):
- "Upload a photo" — primary actor *capability owner*, goal: get media into the platform; addresses Triggers (a capability owner brings new media) and Outputs (durable storage).
- "Recover from a lost device" — primary actor *end user*, goal: get back to their content after losing access; addresses Success Criteria (continuity of access).

**Not an experience:** an internal system flow with no human in the loop ("the system reconciles state every 5 minutes"). Those are implementation details.

**Not a single experience:** a bundle covering multiple personas or multiple goals at once ("user manages their account" — that's at least: sign up, change password, recover access, delete account). Split it.

**Not capability-scoped:** a journey that obviously belongs to a different capability (e.g. while planning experiences for *self-hosted-personal-media-storage*, the user proposes "operator stands up the hosting platform" — that's a journey of the *self-hosted-application-platform* capability). See "Flag-and-stop for misplaced journeys" below.

## Enumeration discipline

- **Read the capability first, then propose the experience set.** Don't go section-by-section mechanically. Multiple capability sections typically motivate one experience, and one experience can address multiple sections. Cluster naturally — but always be ready to point back to the specific section that motivates each proposed UX.
- **One persona + one goal per experience.** Bundling — "user manages content" covering upload *and* share *and* delete — produces UX docs with mismatched journeys, mismatched success states, and edge cases that contradict each other. Resist it. If the user proposes a bundle, split it explicitly:

  > "Let's split that. I see at least three distinct journeys: uploading a new photo, sharing an existing one, and deleting content. Each has a different entry point and a different success state. Bundling them would force the UX doc to cover incompatible flows. Which do you want first?"

- **Mirror back before filing issues.** State the proposed list aloud; let the user add, remove, reorder, rename, or merge. Only file issues once the user has approved the list.
- **Don't draft journey steps yet.** If you find yourself starting to describe what the user sees on screen 1, stop — that's `define-user-experience`. The list at this stage is journeys to be defined, not their step-by-step content.
- **Don't pick personas the capability doesn't list.** If a proposed journey has a primary actor not named in the capability's *Stakeholders* section, surface the gap rather than inventing the persona. Either the capability needs an update (route back to `define-capability`) or the journey isn't really in scope.

## Flag-and-stop for misplaced journeys

If a proposed journey obviously belongs to a different capability, do not include it. Surface it:

> "This journey looks like it belongs to a different capability — {other-capability}. I can either (a) skip it and proceed with the {current-capability}-scoped list, or (b) pause so you can plan it under {other-capability} instead. Which do you want?"

Cross-capability journeys are out of scope here. Do not draft or file issues for them under the current capability.

## Filing the issues

Once the user approves the list, file one GitHub issue per experience via `gh issue create`. Each issue:

- **Title:** `story(ux): {short verb-led journey title} — {capability-name}` (matches the repo's `story(scope): description` convention).
- **Body:** parent capability link, the persona, the goal in one sentence, and which capability sections (Stakeholders / Triggers / Outputs / Business Rules / Success Criteria) the journey is anchored in.
- **Body must reference `define-user-experience`** as the skill that will author the UX doc, and explain that one invocation of `define-user-experience` produces one UX file under `docs/content/capabilities/{capability-name}/user-experiences/`.

After filing, print the issue numbers/URLs back to the user as a manifest.

### Issue body template

```markdown
### Experience

{One sentence in the user's voice: what this persona is trying to accomplish under this capability. E.g., "As a capability owner, I want to bring a new capability to be hosted on the platform so that I don't have to solve hosting myself."}

### Persona

{Persona name, taken from the parent capability's Stakeholders section.}

### Anchored in the parent capability

- **Stakeholders:** {names referenced}
- **Triggers / Inputs:** {bullet(s) referenced, if applicable}
- **Outputs / Deliverables:** {bullet(s) referenced, if applicable}
- **Business Rules:** {rules referenced, if applicable}
- **Success Criteria:** {criteria referenced, if applicable}

### Parent capability

[{capability-name}]({capability-doc-relative-path})

> Use the actual capability doc path here: `../docs/content/capabilities/{name}.md` if the capability is a flat file, or `../docs/content/capabilities/{name}/_index.md` if it has been promoted to page-bundle form. Don't hardcode `_index.md` — the link will 404 for flat capabilities.

### Authoring

This UX will be authored via the `define-user-experience` skill — one invocation per UX. The skill will elicit persona, goal, entry point, journey, success, and edge cases, and save the doc under `docs/content/capabilities/{capability-name}/user-experiences/{ux-name}.md`.

### Related

#{parent-capability-issue-or-722}
```

## Conversation discipline

- **Announce the gate result before doing anything else.** "Capability is reviewed (date `{X}`, file unchanged since) — proceeding." Or "Not yet reviewed — stopping." This makes the gate visible.
- **Read the capability out loud before proposing.** Don't start by guessing the journey set; ground each proposed UX in specific capability sections.
- **Stop at "approved list."** Until the user says "yes, file these," no issues are created. Don't pre-file in case the user changes their mind.
- **No journey steps, no UI, no draft UX docs in this skill.** Those are `define-user-experience`'s job. The line is sharp on purpose.

## Output format

Conversational user-facing output (the proposed list, the gate state, etc.) plus, on approval, GitHub issues filed via `gh issue create`. No markdown artifact written to the repo — the artifacts are the issues themselves.

## Examples of good redirection

**Example 1 — gate not satisfied:**
User: "Plan the experiences for self-hosted-personal-media-storage."
Response (after reading the file): "Stopping — the capability doc has no `reviewed_at` frontmatter field, so I have no signal that the current contents have been reviewed. Add `reviewed_at: {today}` to the frontmatter once you've read the doc and confirmed it reflects the capability you want to plan against, then re-invoke me. (If you want to confirm verbally, say so and I'll set the field on your behalf.)"

**Example 2 — bundle proposed:**
User: "Let's do one experience called 'user manages their content' that covers uploading, sharing, browsing, and deleting."
Response: "Let's split. I count at least four distinct journeys in there: uploading new content, sharing existing content, browsing what's already stored, and deleting. Each has a different entry point ('I just took a photo' vs. 'I want to send something' vs. 'I want to look at memories' vs. 'I want this gone'), and a different success state. A single UX doc covering all four would force mismatched flows into one document. I'll file four issues — pick which order to tackle them in, or I'll suggest one based on the capability's outputs. Want to keep going as four, or merge any?"

**Example 3 — misplaced journey surfaced:**
User: "Add an experience for the operator standing up the hosting platform."
Response: "That journey looks like it belongs to the *self-hosted-application-platform* capability, not *self-hosted-personal-media-storage* — the operator standing up the platform is the platform's concern, and the media-storage capability inherits the platform as a precondition. Want to skip it and continue with the media-storage list, or pause so you can plan that journey under the platform capability instead?"

**Example 4 — persona not in capability:**
User: "Plan an experience for a billing administrator reviewing tenant usage."
Response: "The capability doc doesn't list a billing administrator in *Stakeholders* — the named actors are the operator, capability owners, tenant capabilities, and end users. Either (a) we update the capability to add this persona (back to `define-capability` for that), or (b) this journey isn't really in scope here. Which? I won't invent a stakeholder the capability hasn't sanctioned."
