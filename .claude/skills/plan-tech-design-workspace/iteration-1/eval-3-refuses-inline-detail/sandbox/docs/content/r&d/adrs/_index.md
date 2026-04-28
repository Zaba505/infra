---
title: "Architecture Decision Records"
description: >
    Documentation of architectural decisions made using MADR 4.0.0 standard
type: docs
weight: 20
---

## Architecture Decision Records (ADRs)

This section contains architectural decision records that document the key design choices made. Each ADR follows the MADR 4.0.0 format and includes:

- Context and problem statement
- Decision drivers and constraints
- Considered options with pros and cons
- Decision outcome and rationale
- Consequences (positive and negative)
- Confirmation methods

### ADR Categories

ADRs are classified into three categories:

- **Strategic** - High-level architectural decisions affecting the entire system (frameworks, authentication strategies, cross-cutting patterns). Use for foundational technology choices.
- **User Journey** - Decisions solving specific user journey problems. More tactical than strategic, but still architectural. Use when evaluating approaches to implement user-facing features.
- **API Design** - API endpoint implementation decisions (pagination, filtering, bulk operations). Use for significant API design trade-offs that warrant documentation.

### Status Values

Each ADR has a status that reflects its current state:
- `proposed` - Decision is under consideration
- `accepted` - Decision has been approved and should be implemented
- `rejected` - Decision was considered but not approved
- `deprecated` - Decision is no longer relevant or has been superseded
- `superseded by ADR-XXXX` - Decision has been replaced by a newer ADR

These records provide historical context for architectural decisions and help ensure consistency across the platform.