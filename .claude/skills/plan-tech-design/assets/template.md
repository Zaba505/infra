---
title: "Tech Design"
description: >
    Composed final-state technical design for the {{capability_name}} capability. Synthesizes the outcomes of all accepted ADRs into a single human-friendly read-through.
type: docs
---

> **Composed document.** Synthesizes accepted ADRs in `adrs/` and the requirements in `tech-requirements.md`. For *why* a decision was made, follow the ADR link. This doc covers *what* the system looks like once the decisions are realized.

**Parent capability:** [{{capability_name}}]({{< ref "_index.md" >}})
**Inputs:** [Technical Requirements]({{< ref "tech-requirements.md" >}}) · [ADRs]({{< ref "adrs/_index.md" >}}) · [User Experiences]({{< ref "user-experiences/_index.md" >}})

## Overview

{{overview_paragraph}}

## Components

The pieces that make up this capability and how they connect.

### Component diagram

```mermaid
{{component_mermaid}}
```

### Inventory

For each component: what it is, where it lives in the repo, and which ADR(s) put it there.

{{component_inventory}}

<!--
Example entry:

#### tenant-registry service
**Location:** `services/tenant-registry/`
**Established by:** [ADR-0002: Tenant state ownership]({{< ref "adrs/0002-tenant-state-ownership.md" >}})
**Responsibility:** Source of truth for tenant identity and lifecycle state. Owns the tenants table.
-->

## Key flows

One sequence diagram per user experience, showing how the components realize that journey. **Per-component-design concern** — sequence diagrams are added here as each component's design doc is written via `define-component-design`. Until then, this section lists the UXs and the components they touch.

{{flow_sections}}

<!--
Example entry per UX (skeleton form — full sequence diagram lives in the component design):

### Flow: {ux name}
Realizes [UX: {ux-name}](user-experiences/{ux-name}.md).
Components touched: {list of component names from inventory above}.
Sequence detail: see {component}'s design doc once filed.
-->

## Data & state

What is stored, where, who owns it, and how its lifecycle is managed. **Skeleton form** — schemas, indexes, and storage formats live in per-component design docs (e.g. table-definition components). This section names the data domains and points at the owning component.

{{data_state}}

## How requirements are met

The audit trail. Every TR-NN must appear. If a TR has no ADR or no realization, this document is premature — return to `plan-adrs` / `define-adr`.

| TR | ADR(s) | Realized in |
|----|--------|-------------|
{{requirement_realization_table}}

## Deferred / Open

Decisions intentionally left out of this design (with the reason and rough timing for revisiting), and open questions still under discussion.

{{deferred}}
