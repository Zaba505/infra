---
name: define-capability
description: Guide the user through defining a business capability by eliciting business requirements through conversation. Use this skill whenever the user wants to define, scope, document, or clarify a business capability, business function, or business-level feature — even if they don't say "capability" explicitly. Trigger on phrases like "define a capability", "scope out a feature at the business level", "what should this part of the business do", "document this business function", or when the user is trying to articulate what a piece of the business needs to do before writing code or designing systems. Do NOT use for technical design, architecture decisions, or implementation planning.
---

# Define a Business Capability

This skill helps the user define a business capability — a discrete piece of *what the business does*, framed in business terms. The output is a markdown document the team can review, refine, and use as the upstream input for later technical design.

A business capability answers **what** and **why**, never **how**. If the conversation drifts toward technology choices, system boundaries, APIs, schemas, or implementation, redirect.

## Why this matters

Teams burn weeks designing solutions to problems they haven't defined. A clear business capability is the contract between the people who need work done and the people who will build it. Getting the business definition crisp — before anyone touches a whiteboard or an IDE — saves enormous rework downstream. That's the whole point of this skill.

## The dimensions to elicit

Work through these six dimensions with the user. They don't need to come in order, and the user may answer several at once — follow their lead, but make sure none are left vague before finalizing.

1. **Purpose & business outcome** — What outcome does this capability deliver? Whose life gets better, what gets unblocked, what risk gets mitigated? Push for outcomes, not activities. ("We process refunds" is an activity. "Customers recover funds within 48 hours so they trust us enough to buy again" is an outcome.)

2. **Stakeholders** — Who owns it? Who initiates it (primary actor)? Who consumes its output (secondary actor)? Who is affected without being directly involved? Names of roles, not systems.

3. **Triggers & inputs** — What event causes this capability to engage? What information must be present for it to operate? Preconditions that must be true?

4. **Outputs & deliverables** — What does the capability produce? What changes about the world after it runs — a decision made, a record updated, a notification sent, money moved?

5. **Business rules & constraints** — Policies, regulations, invariants, limits. Anything of the form "must always," "must never," "only when," or "no more than."

6. **Success criteria & KPIs** — How will the business know this capability is working? What measurable outcomes — not throughput, latency, or uptime, but *business* outcomes like conversion rate, cycle time, error rate, customer satisfaction.

Also surface, but don't belabor:
- **Out of scope** — explicit boundaries to prevent scope drift later.
- **Open questions** — anything the user can't answer yet. Capture, don't block.

## How to run the conversation

**Open by anchoring on outcome.** Don't immediately ask all six dimensions in a flat checklist. Start with: "In one sentence, what should this capability accomplish for the business?" From that one sentence you can usually infer the obvious gaps and ask about them next.

**Ask one or two questions at a time.** A wall of questions is intimidating and produces shallow answers. Pick the gap that most needs filling and ask about that.

**Probe vagueness.** "Improve customer experience" is not an answer — it's a slogan. Ask: *whose* experience, *how* will we measure improvement, *what* specifically changes for them. Vague answers become vague capabilities become vague systems.

**Mirror back what you've heard.** Periodically restate the capability in the user's own terms ("So far I'm hearing: this capability lets {actor} do {thing} so that {outcome}, governed by {rules}. Right?"). This catches misunderstanding early.

**Push back on technical drift — firmly but kindly.** If the user starts answering business questions with technical answers ("we'll use a queue", "store it in Postgres", "expose a REST endpoint", "OAuth2 with JWTs"), redirect:

> "That sounds like a *how* — let's park it for the technical design. For now: what does the *business* need to be true? E.g., 'a user must be authenticated before placing an order' is the business rule; whether you use OAuth2 or session cookies is design."

Acknowledge the technical idea (the user may be excited about it), capture it as an "open question" or "implementation note," and return to the business question. Don't lecture. One short redirect per drift is usually enough.

**Watch for these common drifts and redirect each one:**
- Naming systems, services, databases, queues, APIs → "What's the *business* trigger / output, separate from how it gets delivered?"
- Choosing protocols, formats, frameworks → "Park that. What information must be exchanged, in business terms?"
- Designing schemas → "What information does the business need to know, regardless of how it's stored?"
- Discussing scaling, latency, infra → "Is there a *business* constraint here (e.g., 'customer sees confirmation in under a minute')? That belongs. The infra to achieve it doesn't."

**Don't invent answers.** If the user genuinely doesn't know something, capture it as an open question and move on. A capability with three honest open questions is more useful than one with six made-up answers.

**Know when to stop.** When all six dimensions have at least one concrete, non-vague answer (or are explicitly captured as open questions), produce the document. Don't keep asking forever — the document is meant to be revised.

## Producing the document

Use `assets/template.md` as the structure. Fill it in with the user's answers, in their language wherever possible. Replace `{{...}}` placeholders. Remove any section that the user explicitly said is N/A; otherwise leave it with content (or with a captured open question — never leave a `{{placeholder}}` in the final doc).

Save the doc to a path the user specifies, or propose `docs/capabilities/{kebab-case-name}.md` and confirm.

After producing the doc, briefly tell the user:
- Where it was saved
- Which sections still have open questions (so they can chase them down)
- That this is a living document — the next step is to share it with stakeholders for review, not to start building

## Output format

Plain markdown matching `assets/template.md`. No code blocks around the whole document. No frontmatter unless the user requests it.

## Examples of good redirection

**Example 1:**
User: "The capability is user authentication — we'll do OAuth2 with refresh tokens in Redis."
Response: "OAuth2 and Redis are implementation choices — let's park them. At the business level: who needs to be authenticated, before what actions, and what's the outcome we're protecting? E.g., 'only the account owner can change billing details, to prevent unauthorized charges.'"

**Example 2:**
User: "We need an API endpoint that returns the order status."
Response: "Let's back up — what *business* question is being answered? Probably: 'a customer (or support agent) needs to know where their order is in the fulfillment process so they can plan or intervene.' Whether that's an API, a screen, or a notification is later. Who's asking, and what do they do with the answer?"

**Example 3:**
User: "Improve the checkout experience."
Response: "That's a slogan, not a capability yet — let's sharpen it. *Whose* experience, and *what specifically* should be different for them? E.g., 'a returning customer can complete checkout in under three steps so that cart abandonment drops.' What's the equivalent for you?"
