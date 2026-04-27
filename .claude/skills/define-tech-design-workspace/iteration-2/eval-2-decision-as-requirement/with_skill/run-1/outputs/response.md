# Response

We're in **Stage 1: extracting requirements** for the self-hosted-application-platform capability — and what you've proposed is a **Stage 2 decision dressed up as a requirement**, so I'm not adding it as a TR. Two reasons:

1. **"Use Postgres with logical replication" is a chosen solution, not a forced constraint.** A requirement is something the capability or a UX *forces* on the design (e.g. "tenant metadata must be queryable per tenant with strong consistency", or "metadata must survive loss of a single GCP region with RPO ≤ N seconds"). Postgres + logical replication is one of several ways to satisfy such requirements — Spanner, Firestore multi-region, CockroachDB, and others would also be candidates. ADRs only earn their keep when they pick among real alternatives; if we encode the answer as a requirement, the Stage 2 ADR has nothing to decide.
2. **TR-04 already exists** ("Operator-initiated tenant updates must complete without tenant-perceived downtime"). Per the skill's append-only rule, **TR numbers are never reused or renumbered** — ADRs cite TR-NN, so renaming TR-04 would silently break provenance. New requirements get the next free number (TR-08), not an existing one.

### What I did instead

I captured your input as an **Open Question** at the bottom of `tech-requirements.md`:

> **Tenant metadata storage technology and cross-region failover strategy** — user proposed "Postgres with logical replication for cross-region failover" on 2026-04-26. This is a Stage 2 decision, not a requirement. The underlying requirements need to be extracted first (e.g. queryability, consistency, durability, RTO/RPO for cross-region failure)... "consider Postgres + logical replication" will be one of the considered options.

This way the proposal isn't lost — it'll show up as a considered option when we draft the corresponding ADR in Stage 2.

### What I need from you to actually progress

Before we can add a real new TR (which would become TR-08, not TR-04) or proceed to Stage 2, answer these:

1. **What user-visible or business-rule outcome is forcing tenant-metadata storage to support cross-region failover?** I scanned the existing UX docs and capability rules — none explicitly demand multi-region survival of tenant metadata. If there's an availability/disaster-recovery expectation buried in a UX (or in the capability success criteria), point me at it and I'll extract a TR with a proper source link. If there isn't one, the right move is to *first* add it to the capability or a UX doc via `define-capability` / `define-user-experience`, then re-extract here.
2. **What's the implicit non-functional requirement?** Probably something like: "Tenant metadata must remain available and consistent across the loss of a single GCP region, with RPO ≤ X and RTO ≤ Y." If you can articulate the X/Y, I'll draft it as TR-08 with a link to wherever the source rule lives.

Once you give me a sourced rationale, I'll add **TR-08** (not TR-04) and we can move to Stage 2 where Postgres-with-logical-replication competes against other options.

### Files

- Updated: `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (added open question; no TRs renumbered or modified)
