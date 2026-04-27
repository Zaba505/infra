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
The define-tech-design skill refuses Stage 3 until every ADR is accepted (or superseded with the superseder accepted).

LINK CONVENTION (Hugo + Docsy):
- Cross-document links use {{</* relref "..." */>}}, NOT raw `path/file.md`. Hugo
  build-validates relref and resolves it to the correct directory URL.
- Linkable headings carry an explicit `{#id}` attribute. Auto-generated slugs
  break with `:` `&` `§` `(` em-dashes etc.; explicit IDs are stable.
- TR citations use `#tr-NN` anchors set by tech-requirements.md.
-->

**Parent capability:** [{{capability_name}}]({{< relref "../_index.md" >}})
**Addresses requirements:** {{tr_citations}}
<!--
Each TR citation: `[TR-03]({{</* relref "../tech-requirements.md#tr-03" */>}})`
Multiple TRs: list comma-separated, each with its own relref.
-->

## Context and Problem Statement {#context}

{{context}}

State the problem in terms of the cited TRs. A reader should be able to follow: this TR (which links to a UX or capability rule) → this decision needs to be made.

## Decision Drivers {#decision-drivers}

<!--
Driver each line back to a TR or to a constraint inherited from a prior shared ADR / CLAUDE.md pattern.
TR references use `[TR-NN]({{</* relref "../tech-requirements.md#tr-NN" */>}})` format.
-->

{{decision_drivers}}

## Considered Options {#considered-options}

At least two. Each option must be evaluated against the cited TRs.

{{considered_options}}

## Decision Outcome {#decision-outcome}

Chosen option: "{{chosen_option}}", because {{rationale_in_terms_of_TRs}}.

### Consequences {#consequences}

* Good, because {{positive_consequence}}
* Bad, because {{negative_consequence}}
* Requires: {{follow_up_work}}
<!--
Cross-references to other ADRs in this capability:
`[ADR-0002]({{</* relref "0002-compute-substrate.md" */>}})`
References to UX docs:
`[UX: host-a-capability §4]({{</* relref "../user-experiences/host-a-capability.md#step-4-handoff" */>}})`
-->

### Realization {#realization}

How this decision shows up in the codebase: which `services/{name}/`, which `cloud/{module}/`, which `pkg/` packages. Stage 3's `tech-design.md` will compose this with the other ADRs into the final-state narrative.

{{realization}}

## Open Questions {#open-questions}

{{open_questions}}
