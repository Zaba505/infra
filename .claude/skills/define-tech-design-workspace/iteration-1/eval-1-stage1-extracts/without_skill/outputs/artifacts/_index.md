---
title: "Tech Design"
description: >
    Technical design workspace for the Self-Hosted Application Platform capability. Currently at Stage 1 (technical requirements) of the three-stage tech-design flow.
type: docs
weight: 20
---

This section holds the technical design for the [Self-Hosted Application Platform](../_index.md) capability. It is built in three stages:

1. **Technical requirements** — a living extract of what the design must satisfy, traceable to the capability and UX docs. See [`technical-requirements.md`](./technical-requirements.md).
2. **Capability-scoped ADRs** — one MADR per non-trivial decision (compute substrate, packaging form, identity service, etc.). To be drafted next; tracked as `Q-01`–`Q-12` in the requirements doc.
3. **Composed tech-design document** — a human-friendly description of the resulting design, weaving the ADRs together. Written last, when the ADRs have settled.

## Current status

- **Stage 1:** Draft. The technical-requirements doc is populated from the capability + 7 UX docs. Open questions Q-01 through Q-12 are the seeds for Stage 2.
- **Stage 2:** Not started. Recommended next step: pick one of Q-01 through Q-12 to ADR first. Suggested ordering — Q-09 (definitions-repo layout / top-level rebuild) before Q-02 (compute substrate) before everything else, because the substrate decision constrains packaging, networking, identity hosting, secret management, and observability deployment.
- **Stage 3:** Not started.

## Scope reminder

Per the `define-tech-design` flow: shared cross-capability decisions (e.g. "use MADR 4.0.0 for ADRs", "GCP as cloud provider", "network-boot architecture") live in [`docs/content/r&d/adrs/`](../../../r&d/adrs/) and are inputs to this design, not outputs of it. ADRs created here are scoped to the platform capability itself.
