---
title: "[{{number}}] {{title}}"
description: >
    {{one_liner}}
type: docs
weight: {{weight}}
category: "{{category}}"
status: "proposed"
date: {{date}}
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions for this capability (auth strategy, data ownership boundaries)
- user-journey: Solutions for specific user-experience problems within this capability
- api-design: API endpoint design decisions for this capability's services

Numbering is local to this capability — start at 0001 and increment.
Status lifecycle: proposed → accepted → (later) superseded
The plan-tech-design skill refuses to compose tech-design.md until every ADR is accepted (or superseded with the superseder accepted).
-->

**Parent capability:** [{{capability_name}}](../_index.md)
**Addresses requirements:** {{tr_citations}}  <!-- e.g. TR-03, TR-07 -->

## Context and Problem Statement

{{context}}

State the problem in terms of the cited TRs. A reader should be able to follow: this TR (which links to a UX or capability rule) → this decision needs to be made.

## Decision Drivers

<!--
Driver each line back to a TR or to a constraint inherited from a prior shared ADR / CLAUDE.md pattern.
-->

{{decision_drivers}}

## Considered Options

At least two. Each option must be evaluated against the cited TRs.

{{considered_options}}

## Decision Outcome

Chosen option: "{{chosen_option}}", because {{rationale_in_terms_of_TRs}}.

### Consequences

* Good, because {{positive_consequence}}
* Bad, because {{negative_consequence}}
* Requires: {{follow_up_work}} <!-- e.g. new module under cloud/, new service under services/ -->

### Realization

How this decision shows up in the codebase: which `services/{name}/`, which `cloud/{module}/`, which `pkg/` packages. Stage 3's `tech-design.md` will compose this with the other ADRs into the final-state narrative.

{{realization}}

## Open Questions

{{open_questions}}
