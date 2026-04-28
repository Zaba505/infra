# ADR Plan — self-hosted-application-platform

**Capability:** `docs/content/capabilities/self-hosted-application-platform/`
**Tech requirements source:** `tech-requirements.md` (reviewed_at: 2026-04-26, unchanged since)
**Approach:** Enumerate the architectural decisions the TRs *force* before solving any of them. Each decision below becomes one capability-scoped MADR ADR, authored individually via `define-adr`. TR-07 is excluded from new ADRs because it is flagged as a prior shared decision; it appears as a *constraint* on the ADRs below.

## Decision Enumeration

### ADR-A: Tenant compute runtime
- **What's decided:** What hosts tenant workloads — the unit of compute the platform offers (e.g. containers on an orchestrator, VMs, systemd units, a PaaS-like abstraction). This is the foundation other ADRs sit on top of.
- **Forced by:** TR-01 (compute-layer isolation), TR-04 (online workloads survive operator-initiated updates), TR-07 (must live in the Cloudflare → GCP / WireGuard topology).
- **Why distinct:** Picking a runtime is independent of *how* tenants are isolated within it (ADR-B) and of *what* the platform contract looks like (ADR-D). The runtime choice constrains, but does not determine, those.
- **Dependencies:** None — foundational. Should be filed/decided first.

### ADR-B: Tenant isolation model (data + compute)
- **What's decided:** The mechanism by which one tenant cannot read another's state, secrets, traffic, or telemetry — namespacing, separate clusters, separate VMs, per-tenant projects, etc., and the corresponding data-plane partitioning.
- **Forced by:** TR-01 (strict tenant isolation under normal *and degraded* conditions), TR-03 (cross-tenant observability data inaccessible).
- **Why distinct:** Even after the runtime is chosen (ADR-A), there are multiple isolation strategies with different blast-radius and reproducibility profiles. The "degraded operating condition" clause in TR-01 is the load-bearing constraint and deserves its own decision.
- **Dependencies:** ADR-A.

### ADR-C: Tenant persistent storage model
- **What's decided:** How per-tenant durable state is provisioned and partitioned (per-tenant volumes, per-tenant DB instances, per-tenant buckets, etc.), including the integrity primitives that ADR-G and ADR-H build on.
- **Forced by:** TR-01 (data-layer isolation), TR-05 (export must produce a portable format), TR-06 (import must be idempotent with verifiable integrity).
- **Why distinct:** Storage shape is a separable decision from isolation *mechanism* (ADR-B can be satisfied by several storage shapes) and from export/import *mechanics* (ADR-G/H choose format and protocol on top of whatever shape is picked).
- **Dependencies:** ADR-A, ADR-B.

### ADR-D: Platform contract definition and versioning scheme
- **What's decided:** What the "platform contract" actually is (the surface tenants depend on), how versions are identified, how multiple versions coexist, and the bounded migration window's shape.
- **Forced by:** TR-02 (operator must roll out contract changes without breaking existing tenants; multiple versions concurrent for a bounded window).
- **Why distinct:** The contract is a first-class artifact — ADR-A/B/C choose technologies, but the contract is the operator/tenant interface and has its own lifecycle.
- **Dependencies:** ADR-A (the runtime choice influences what's exposable in the contract), but otherwise independent.

### ADR-E: Tenant ingress and zero-downtime update mechanism
- **What's decided:** How tenant-facing traffic is routed and how an operator-initiated update (config / version / capability change) is rolled out so online workloads see no end-user-visible downtime — e.g. blue/green, rolling, draining, request-replay.
- **Forced by:** TR-04 (no end-user-visible downtime on operator-initiated updates), TR-07 (traffic must traverse Cloudflare → GCP).
- **Why distinct:** The update *mechanism* is separable from runtime choice; multiple runtimes support multiple rollout strategies, and the rollout strategy interacts with ingress/load-balancing in ways that warrant their own ADR.
- **Dependencies:** ADR-A, ADR-D (contract changes are one of the things being rolled out).

### ADR-F: Per-tenant observability scoping
- **What's decided:** How metrics, logs, and traces are tagged, stored, and authorized so that each tenant can query *only* their own data, and cross-tenant queries are prevented.
- **Forced by:** TR-03 (per-tenant queryability with strict scope), TR-01 (telemetry is one of the things that must be isolated).
- **Why distinct:** Observability tooling and tenancy model is independent of the runtime/storage choices and has its own vendor / self-host decision space.
- **Dependencies:** ADR-B (isolation primitives carry through to telemetry).

### ADR-G: Tenant data export mechanism
- **What's decided:** The export format, trigger, delivery channel, integrity guarantees, and the defined export window for an evicted tenant.
- **Forced by:** TR-05 (evicted tenant must be able to take their data, portable format, defined window).
- **Why distinct:** Even given a storage model (ADR-C), the export *protocol* (format, packaging, signing, retrieval) is its own decision and shows up in tenant-visible behavior.
- **Dependencies:** ADR-C.

### ADR-H: Tenant data import / migration mechanism
- **What's decided:** How a new tenant's pre-existing data is accepted and imported with idempotency and verifiable integrity (no silent loss, no duplication on retry).
- **Forced by:** TR-06 (idempotent import with integrity verification).
- **Why distinct:** Import is the inverse of export but has different correctness properties (idempotency under retry, integrity verification) and may use different transports than export.
- **Dependencies:** ADR-C. Often paired with ADR-G but kept separate so each can be decided on its own merits.

## Cross-cutting Constraint (not a new ADR)

- **TR-07** ("All inter-service communication must traverse the Cloudflare → GCP path") is marked as a prior shared decision. It is **not** a new ADR for this capability. It is recorded as a constraint that ADR-A, ADR-B, ADR-E, and ADR-F must respect.

## Decision Dependency Order (suggested authoring order)

1. ADR-A (runtime)
2. ADR-B (isolation) — depends on A
3. ADR-D (contract & versioning) — depends loosely on A
4. ADR-C (storage) — depends on A, B
5. ADR-E (ingress + zero-downtime updates) — depends on A, D
6. ADR-F (per-tenant observability) — depends on B
7. ADR-G (export) — depends on C
8. ADR-H (import) — depends on C

## TR → ADR Coverage Matrix

| TR    | Covered by                                  |
|-------|---------------------------------------------|
| TR-01 | ADR-A, ADR-B, ADR-C, ADR-F                  |
| TR-02 | ADR-D                                       |
| TR-03 | ADR-F                                       |
| TR-04 | ADR-E                                       |
| TR-05 | ADR-G (uses ADR-C)                          |
| TR-06 | ADR-H (uses ADR-C)                          |
| TR-07 | constraint on ADR-A, ADR-B, ADR-E, ADR-F    |

Every TR is covered; every ADR traces to at least one TR.

## Out of Scope for This Plan
- Drafting any ADR's options, decision, or consequences (that is `define-adr`'s job).
- Adding new technical requirements (that is `define-technical-requirements`'s job).
- Cross-capability / repo-wide decisions (those go in `docs/content/r&d/adrs/`, not capability-scoped).
