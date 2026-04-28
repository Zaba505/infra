---
name: define-user-experience
description: Guide the user through defining one user experience (UX) — an end-to-end user journey — for an already-defined business capability, and document it under that capability. Use this skill whenever the user wants to define, scope, document, or clarify a user experience, user journey, user flow, or "how a user accomplishes X" for a capability that already exists. Trigger on phrases like "define the UX for {capability}", "what's the user journey for X", "how does a user upload/share/view…", "document the user experience", "add a UX to this capability". Do NOT use for technical design, wireframes, screen-by-screen UI specs, or implementation planning. Do NOT use to define a brand-new capability — for that, use `define-capability` first.
---

# Define a User Experience for a Capability

This skill helps the user define **one user experience** (a single end-to-end journey) for an **existing** business capability. The output is a markdown document saved alongside the parent capability so the two stay linked.

A user experience answers **what the user does and perceives, in what order, to accomplish a goal under this capability**. It does *not* answer **how the system implements that** — no screens, no schemas, no APIs, no protocols.

## Why this matters

A capability says *what the business does*. A UX says *how that lands for one specific user accomplishing one specific goal*. Skipping straight from capability to implementation produces systems that technically meet the spec but feel wrong to use. A clear UX doc catches that disconnect before code is written, and it gives designers and engineers a shared, persona-grounded picture to build against.

A capability typically has **multiple** user experiences (e.g. *upload a photo*, *share an album*, *bulk import from another provider*, *delete content*). Each gets its own document.

## Preconditions — read the parent capability first

Before eliciting anything, **find and read the parent capability doc**. Capabilities live at `docs/content/capabilities/{capability-name}.md` *or* `docs/content/capabilities/{capability-name}/_index.md` (page bundle form).

If the user has not named the parent capability, ask which one. If no capability doc exists yet, stop and tell them to run the `define-capability` skill first — UX docs are children of capabilities and shouldn't float on their own.

From the capability doc, internalize:
- **Stakeholders** — the personas you can pick from for this UX.
- **Business rules & constraints** — the invariants this UX must respect.
- **Triggers & inputs / outputs** — the touchpoints the UX must hook into.
- **Success criteria** — the business outcomes this UX should contribute to.

You will reference these explicitly in the UX doc's *Constraints Inherited from the Capability* section. Do not reinvent rules the capability already established.

## The dimensions to elicit

Work through these with the user. They don't need to come in order; follow the user's lead but don't ship the doc with any of these vague.

1. **Persona** — *Which* actor from the parent capability is having this experience? What's their context entering this journey? What do they care about? Pick one — if the user wants to describe two personas at once, that's usually two UX docs.

2. **Goal** — In one sentence, from the user's point of view, what are they trying to accomplish? ("I want my last vacation's photos backed up automatically without thinking about it" — not "the system performs camera roll sync.")

3. **Entry point** — How does the user arrive? Did they tap something, get a notification, plug in a device, hit a deadline? What state of mind are they in?

4. **Journey** — The end-to-end flow, in plain language, step by step. Each step says what the user *does* and what they *perceive*. No system internals. A Mermaid flowchart is encouraged for branches and decision points; the prose narrative is the source of truth.

5. **Success** — What does a successful completion look and feel like for the user? What do they walk away with — a confirmation, a sense of relief, a thing they can now share with someone?

6. **Edge cases & failure modes** — What can go wrong *from the user's perspective*? Network drops mid-upload, recipient isn't a user yet, file too large, accidental deletion. For each, what should the experience do — at the experience level, not the implementation? ("The user sees their upload paused and can resume from where it stopped" is UX. "We checkpoint chunks in Redis with a 24h TTL" is not.)

7. **Constraints inherited from the capability** — Which specific business rules, stakeholders, or success criteria from the parent capability shape this UX? Cite them by name. This is how the doc stays anchored.

Also surface, but don't belabor:
- **Out of scope** — adjacent journeys this doc does *not* cover, so the next UX has a clear seam.
- **Open questions** — capture, don't block.

## How to run the conversation

**Open by anchoring on persona + goal.** "Which user, accomplishing what?" is the whole frame. Once those two are crisp, every other dimension is much easier to elicit. Don't dive into the journey steps until persona and goal are nailed down — otherwise you'll redesign the flow when the user clarifies who it's for.

**Walk the journey one step at a time.** Don't ask for the whole flow up front. Start at the entry point and ask "what happens next, from the user's view?" Repeat until they reach the goal. Branches and decisions emerge naturally this way.

**Probe vagueness.** "The user shares a photo" is not a step — it's a label. Ask: how do they choose the recipient? What do they see when they're done? How do they know it worked? Vague journey steps become vague designs become vague experiences.

**Push back on technical drift — firmly but kindly.** If the user starts answering UX questions with technical answers ("we'll use signed URLs", "WebRTC for the upload", "OAuth scope check"), redirect:

> "That's an implementation detail — let's park it. At the experience level: what does the user *see and do* here? The mechanism behind it is a design choice for later."

Capture the technical idea as an "open question" or "implementation note" if they want, then return to the user's perspective.

**Watch for these common drifts and redirect each one:**
- Naming screens, components, modals, endpoints → "Describe what the user *does and perceives*, not the UI element. The UI gets designed from this doc, not specified in it."
- Choosing protocols, storage, sync mechanisms → "Park that. At the experience level, what does the user expect to be true after this step?"
- Designing schemas or data models → "What information does the user need to *see or provide*, in their own terms?"
- Discussing scaling, latency, infra → "Is there an *experience* expectation here (e.g., 'feels instant', 'can leave the app and come back')? That belongs. Implementation to achieve it doesn't."

**Anchor each step in the parent capability's rules.** When the user describes a step, mentally check the capability's business rules. If a step would violate one ("…and then the operator approves the share" — but the capability says all content is private to its owner unless *they* share it), surface the conflict immediately. The capability is the contract; the UX must respect it.

**Mirror back periodically.** Restate the journey in the user's own terms ("So: {persona} starts at {entry}, then {step}, then {step}, and walks away with {success} — right?"). Catches divergence early.

**Don't invent answers.** If the user doesn't know an edge case yet, capture it as an open question. A doc with three honest open questions is more useful than one with six made-up answers.

**Know when to stop.** When persona, goal, entry, journey, success, and at least the obvious edge cases are concrete (or explicitly captured as open questions), produce the document. The doc is meant to be revised.

## Producing the document

Use `assets/template.md`. Fill in `{{...}}` placeholders with the user's answers, in their language wherever possible. Never leave a `{{placeholder}}` in the final doc.

### Where to save it

UX docs live **under the parent capability** as a page bundle, at:

```
docs/content/capabilities/{capability-name}/user-experiences/{ux-name}.md
```

This means the parent capability must be in **page-bundle form** — i.e. `docs/content/capabilities/{capability-name}/_index.md`, not a flat `{capability-name}.md` file. If the parent capability is currently a flat file, perform this one-time migration before saving the UX:

1. Create directory `docs/content/capabilities/{capability-name}/`.
2. Move (`git mv`) the existing `{capability-name}.md` to `{capability-name}/_index.md`.
3. Confirm the migration with the user before doing it — moving the file changes its URL only marginally (Hugo treats `{name}.md` and `{name}/_index.md` as the same path), but it's still a structural change worth confirming.

Then ensure the `user-experiences/` subdirectory exists, and create a `user-experiences/_index.md` if it doesn't, with Docsy section frontmatter (`title: "User Experiences"`, `description`, `type: docs`, `weight`). Without this, Hugo will not render `/capabilities/{name}/user-experiences/` as a section page.

The UX doc filename should be kebab-case and verb-led — `upload-photo.md`, `share-album.md`, `bulk-import-from-google-photos.md`. Names describe the *journey*, not the persona.

### Cross-linking

In the UX doc, set the `Parent capability:` link using Hugo's `ref` shortcode — `[{{capability_name}}]({{< ref "../_index.md" >}})`. **Never use raw paths** (`(../_index.md)`, `(/capabilities/{name}/)`); Hugo's build will fail loudly on a broken `ref`, while raw paths break silently when content is reorganized.

If a link targets a specific section on a page, the target heading must carry an explicit anchor — e.g. `## Business Rules {#business-rules}` — and the link uses that anchor: `{{< ref "../_index.md#business-rules" >}}`. Do not rely on Hugo's slugified-from-heading-text default; rewording a heading breaks every implicit anchor pointing at it.

In the parent capability's `_index.md`, optionally add a `## User Experiences` section listing the children — but only if the user wants it; Docsy will surface them in the sidebar regardless.

After producing the doc, briefly tell the user:
- Where it was saved
- Which sections still have open questions
- That this is one of likely several UX docs for the capability — what other journeys did the capability imply that aren't covered yet?

## Validate the rendered page

After saving the UX doc, run `cd docs && hugo server` and open the page in a browser to confirm it renders correctly — including the parent-capability `ref` link and the Mermaid flow diagram. A clean build is necessary but not sufficient: pages can build with no errors yet render broken (missing/incorrect front-matter `type`, unresolved `ref` shortcode, malformed Mermaid). The doc is not done until the rendered page has been visually validated. If you cannot run the preview in this session, say so explicitly to the user rather than reporting the doc as done.

## Output format

Plain markdown matching `assets/template.md`. Keep the Hugo/Docsy frontmatter — it is required for the page to render. Mermaid blocks render natively in Docsy.

## Examples of good redirection

**Example 1:**
User: "For the upload UX, the client will chunk the file into 5MB pieces and POST each to a presigned URL."
Response: "That's the implementation. At the experience level: the user picks a file (or a phone takes a photo automatically) — what do they see while it's uploading, what do they see when it's done, and what happens if they close the app halfway through? Those are the UX questions. Chunking is how we *deliver* on the answer."

**Example 2:**
User: "The share screen has a search box, a recent contacts row, and a share button at the bottom."
Response: "Hold on — we're describing UI, not the experience. Back up: the user has decided to share something. Whom are they trying to reach (one person? a group?), how do they expect to find that recipient (by name? a list they've built?), and what do they want to be true after they tap 'share'? The screen design follows from those answers."

**Example 3:**
User: "The journey is: user uploads, system stores, user views."
Response: "That's the system's journey, not the user's. From the user's seat: they pulled out their phone because… something prompted them. They picked a photo because… they wanted what? They tapped upload and then… what did they see, and how did they know it worked? Each of those is a step in *their* experience."
