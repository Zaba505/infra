---
title: "Tech Design"
description: >
    Composed final-state technical design for the {{capability_name}} capability. Synthesizes the outcomes of all accepted ADRs into a single human-friendly read-through.
type: docs
---

> **Composed document.** Synthesizes accepted ADRs in `adrs/` and the requirements in `tech-requirements.md`. For *why* a decision was made, follow the ADR link. This doc covers *what* the system looks like once the decisions are realized.

<!--
LINK CONVENTION (Hugo + Docsy):
- Cross-document links use {{</* relref "..." */>}}, NOT raw `path/file.md`.
- Linkable headings carry an explicit `{#id}` attribute.
- TR citations: `[TR-NN]({{</* relref "tech-requirements.md#tr-NN" */>}})`.
- ADR citations: `[ADR-NNNN]({{</* relref "adrs/NNNN-name.md" */>}})`.
- UX citations: `[UX: name]({{</* relref "user-experiences/name.md" */>}})`.
-->

**Parent capability:** [{{capability_name}}]({{< relref "_index.md" >}})
**Inputs:** [Technical Requirements]({{< relref "tech-requirements.md" >}}) · [ADRs]({{< relref "adrs/_index.md" >}}) · [User Experiences]({{< relref "user-experiences/_index.md" >}})

## Overview {#overview}

{{overview_paragraph}}

## Components {#components}

The pieces that make up this capability and how they connect.

### Component diagram {#component-diagram}

```mermaid
{{component_mermaid}}
```

### Inventory {#inventory}

For each component: what it is, where it lives in the repo, and which ADR(s) put it there.

{{component_inventory}}

<!--
Example entry:

#### tenant-registry service {#tenant-registry-service}
**Location:** `services/tenant-registry/`
**Established by:** [ADR-0002: Tenant state ownership]({{</* relref "adrs/0002-tenant-state-ownership.md" */>}})
**Responsibility:** Source of truth for tenant identity and lifecycle state. Owns the tenants table.
-->

## Key flows {#key-flows}

One sequence diagram per user experience, showing how the components above realize that journey.

{{flow_sections}}

<!--
Example entry per UX:

### Flow: {ux name} {#flow-ux-name}
Realizes [UX: {ux-name}]({{</* relref "user-experiences/{ux-name}.md" */>}}).

```mermaid
sequenceDiagram
    ...
```

Notes: any subtleties not visible in the diagram.
-->

## Data & state {#data-state}

What is stored, where, who owns it, and how its lifecycle is managed.

{{data_state}}

## How requirements are met {#how-requirements-are-met}

The audit trail. Every TR-NN must appear. If a TR has no ADR or no realization, this document is premature — return to Stage 2.

| TR | ADR(s) | Realized in |
|----|--------|-------------|
{{requirement_realization_table}}

<!--
Each row links the TR and the ADR(s) via relref:
| [TR-01]({{</* relref "tech-requirements.md#tr-01" */>}}) | [ADR-0002]({{</* relref "adrs/0002-compute-substrate.md" */>}}) | `homelab/cluster/` |
-->

## Deferred / Open {#deferred-open}

Decisions intentionally left out of this design (with the reason and rough timing for revisiting), and open questions still under discussion.

{{deferred}}
