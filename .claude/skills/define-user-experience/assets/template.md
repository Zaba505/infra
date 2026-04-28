---
title: "{{name}}"
description: >
    {{one_liner}}
type: docs
---

> **One-line definition:** {{one_liner}}

**Parent capability:** [{{capability_name}}]({{< ref "../_index.md" >}})

## Persona
Who is having this experience? Pick one of the actors named in the parent capability's Stakeholders section, and describe them in the context of *this* journey.

- **Role:** {{persona_role}}
- **Context they come from:** {{persona_context}}
- **What they care about here:** {{persona_motivations}}

## Goal
In one sentence, what is this user trying to accomplish? State it from their point of view.

{{goal}}

## Entry Point
How does the user arrive at this experience? What triggers them to start?

{{entry_point}}

## Journey
The end-to-end flow, in plain language. Each step describes what the user does and what they perceive — not how the system implements it.

{{journey_steps}}

### Flow Diagram
```mermaid
{{mermaid_flow}}
```

## Success
What does a successful completion of this experience look and feel like for the user? What do they walk away with?

{{success}}

## Edge Cases & Failure Modes
What can go wrong from the user's perspective? How should each case be handled at the experience level (not the implementation level)?

{{edge_cases}}

## Constraints Inherited from the Capability
Which business rules, stakeholders, or success criteria from the parent capability shape this experience? Cite them so future readers can trace the lineage.

{{inherited_constraints}}

## Out of Scope
What this experience explicitly does **not** cover. Other journeys for the same capability go in their own docs.

{{out_of_scope}}

## Open Questions
Decisions deferred or unknowns that need follow-up before this experience can be built.

{{open_questions}}
