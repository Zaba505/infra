---
description: Create a new User Journey document
---

Create a new User Journey document to define user flows and technical requirements.

**Requirements:**
1. All User Journeys MUST be created in `docs/content/r&d/user-journeys/`
2. Use naming format: `NNNN-title-with-dashes.md` where NNNN is zero-padded sequential (e.g., `0001-user-registration.md`)
3. Ask the user for the journey title/topic if not provided
4. Determine the next sequential number by checking existing journeys in `docs/content/r&d/user-journeys/`
5. Fill in today's date in YYYY-MM-DD format
6. Set initial status to "draft"
7. Create the `docs/content/r&d/user-journeys/` directory if it doesn't exist
8. Create an `_index.md` file in `docs/content/r&d/user-journeys/` if it doesn't exist

**User Journey Template to use:**

```markdown
---
title: "[NNNN] [Journey Title]"
description: >
    [Brief summary of what this user journey accomplishes]
type: docs
weight: [NNNN]
status: "draft"
date: YYYY-MM-DD
owner: ""
stakeholders: []
---

## Overview

[Provide a brief description of the user journey, its purpose, and the business value it provides.]

## User Personas

**Primary Persona**: [Name/Role]
- **Description**: [Who they are and their context]
- **Goals**: [What they want to accomplish]
- **Pain Points**: [What challenges they face]

**[Additional Personas if applicable]**:
- [Describe any variations in the user experience based on different personas]
- [Note any persona-specific requirements or flow variations]

## Journey Flow Diagram

Use Mermaid syntax to create a flowchart representing the user journey:

```mermaid
graph TD
    A[User starts] --> B{Decision point}
    B -->|Option 1| C[Action 1]
    B -->|Option 2| D[Action 2]
    C --> E[Next step]
    D --> E
    E --> F[Journey complete]
```

**Flow Narrative**:

[Provide a step-by-step description of the flow diagram, explaining key decision points, user actions, and any persona-specific variations in the experience]

## Requirements

### Access Control

#### REQ-AC-001
- **Priority**: P0 | P1 | P2
- **Description**: [What access control is needed]
- **Rationale**: [Why this is important from a user perspective]

### Analytics

#### REQ-AN-001
- **Priority**: P0 | P1 | P2
- **Description**: [What analytics/tracking is needed]
- **Rationale**: [Why this is important for measuring user experience]

## Success Metrics

[Define how success will be measured for this journey. Include both quantitative metrics (e.g., completion rate, time to complete) and qualitative metrics (e.g., user satisfaction).]

**Quantitative Metrics:**
- **[Metric Name]**: [Description and target value]
  - [What it measures]

**Qualitative Metrics:**
- **[Metric Name]**: [Description and target]
  - [How it will be gathered]

## Related Documentation

**Existing ADRs:**
- [Link to related ADRs that support this journey]

**Required ADRs (Not Yet Created):**
- **[ADR Topic]** - [Brief description of what needs to be decided and why it's needed for this journey]

**Related User Journeys:**
- [Link to related user journeys]

**API Documentation:**
- [Link to API documentation if it exists, or note that it needs to be created]

## Notes

[Any additional context, constraints, or considerations that don't fit in the sections above]

**Future Enhancements:**
- [Potential improvements or features for later iterations]

**Technical Considerations:**
- [Important technical factors that should inform implementation]

**Business Considerations:**
- [Business implications or opportunities related to this journey]
```

**Status values:** `draft` | `in-review` | `approved` | `implemented` | `deprecated`

**Priority levels:**
- `P0` (Must Have): Critical requirements that must be included in the initial design
- `P1` (Should Have): Important requirements that should be included but can be phased in if necessary
- `P2` (Nice to Have): Optional enhancements for future iterations

**Requirement ID format:** `REQ-[CATEGORY]-NNN`
- `AC` = Access Control
- `AN` = Analytics

**Reference:** Mermaid diagram syntax - https://mermaid.js.org/syntax/flowchart.html