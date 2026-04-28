---
title: "[0001] Tenant state storage"
description: >
    How the platform stores per-tenant state such that tenants are isolated, exportable on eviction, importable on onboarding, and migrate-able without downtime.
type: docs
weight: 1
category: "strategic"
status: "proposed"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [self-hosted-application-platform](../_index.md)
**Addresses requirements:** TR-01, TR-03, TR-04, TR-05, TR-06

## Context and Problem Statement

The self-hosted-application-platform hosts multiple tenants. Each tenant produces and owns state — application data, configuration, and per-tenant slices of observability data. The platform needs a storage strategy that simultaneously holds five constraints from `tech-requirements.md`:

- **TR-01** (isolation): no tenant may read another tenant's state under any normal or degraded condition.
- **TR-03** (per-tenant observability scope): observability data must be queryable per-tenant within the tenant's scope only — meaning the storage choice has to make per-tenant scoping enforceable at the data layer, not just at the query layer.
- **TR-04** (no-downtime operator updates): tenant updates initiated by the operator must complete without tenant-perceived downtime, so the storage layer must support online schema/contract migration or per-tenant rolling cutover.
- **TR-05** (export on eviction): an evicted tenant must be able to take their data with them in a portable format within a defined export window. The storage layer must make a clean per-tenant export cheap.
- **TR-06** (idempotent import): new tenants must be able to import existing data with verifiable integrity (no silent loss, no duplication on retry). The storage layer must support an idempotent write path or expose primitives that one can build idempotency on top of.

TR-07 (Cloudflare→GCP topology) constrains all options below: the storage backend must live in GCP and be reachable along the existing topology. ADR-0003 (cloud-provider-selection: GCP) and ADR-0006 (resource-identifier-standard) bind us further — anything we provision is GCP-native and per-tenant resources must be named according to the resource-identifier standard.

## Decision Drivers

- **TR-01 isolation strength** — how hard is it to violate isolation by accident or by a query bug? A storage choice that makes cross-tenant reads structurally impossible is preferable to one that relies on every query being filtered correctly.
- **TR-03 per-tenant observability scoping** — the storage layer for tenant-facing telemetry must let us issue queries that *cannot* return another tenant's data.
- **TR-04 online updates** — schema or contract changes must be deliverable per-tenant without taking tenants offline.
- **TR-05 export cleanliness** — exporting one tenant's data must be a bounded operation, not a scan-and-filter across a shared dataset.
- **TR-06 idempotent import** — the storage primitive must support deterministic-key writes (upsert by tenant-scoped key) so that retried imports do not duplicate.
- **CLAUDE.md house pattern** — services are chi+bedrock with protobuf endpoints; backend clients live in `services/{name}/service/`. Firestore is already an in-house pattern (`services/machine/`). Departures from this require explicit justification.
- **ADR-0003 (GCP)** and **ADR-0006 (resource identifier standard)** — option must be GCP-native and per-tenant resources must be named per the standard.
- **Operational cost / complexity** — provisioning N databases per tenant has a different operational profile than one shared database with N logical scopes, and the platform is intentionally over-engineered for learning but not for unbounded ops cost.

## Considered Options

### Option A — Database-per-tenant on Cloud SQL for PostgreSQL

A separate Cloud SQL Postgres database (or instance, depending on density) is provisioned per tenant. Per-tenant connection strings are issued and held in Secret Manager. Schema migrations run against each tenant's database independently; rollouts are per-tenant.

**Pros**
- **TR-01:** isolation is structural — a query against tenant A's database cannot return tenant B's rows because the connection itself is tenant-scoped. Strongest isolation primitive of the three options.
- **TR-04:** per-tenant rollout is the natural unit of change; an operator update to one tenant is a migration on that tenant's database only, leaving every other tenant untouched. Fits no-downtime semantics cleanly.
- **TR-05:** export is `pg_dump` of one database — bounded, well-understood, portable.
- **TR-06:** Postgres `INSERT … ON CONFLICT DO NOTHING` (or `UPSERT`) gives idempotent import on a tenant-scoped primary key directly.

**Cons**
- **Operational cost** scales linearly with tenant count: each tenant is at minimum a database (and potentially an instance if density is constrained). For a learning-scale platform this is real money and real ops surface.
- **TR-03:** Postgres is not the natural store for telemetry (logs/traces); we'd still need a separate observability storage decision. This option only solves application state, not the full TR-03 surface.
- **TR-04 nuance:** while per-tenant is great for tenant-update rollout, *platform-contract* changes (TR-02, not directly addressed here) become N migrations to coordinate.
- New `cloud/cloud-sql-tenant/` Terraform module needed; new bedrock-config plumbing for per-tenant DSNs.

### Option B — Shared Firestore with per-tenant document-path prefixing and tenant-scoped security rules

One Firestore database in the platform's GCP project. All tenant data lives under a top-level collection keyed by tenant-id (e.g. `/tenants/{tenant_id}/...`), with Firestore security rules enforcing that any access path must include the caller's tenant-id. Matches existing in-house pattern (`services/machine/` already uses Firestore).

**Pros**
- **House pattern fit:** Firestore is already wired in `services/machine/`; new services can reuse the same `service/` client shape per CLAUDE.md.
- **Operational cost:** one database, no per-tenant provisioning. Cheapest of the three at the platform-control-plane level.
- **TR-04:** schema is implicit (document store), so contract migrations don't require offline DDL. Per-tenant rollout is achievable by routing a tenant's writes through a migrated code path.
- **TR-06:** Firestore document writes are idempotent on a deterministic document ID — natural fit for `import-by-tenant-scoped-key`.

**Cons**
- **TR-01 weakness:** isolation is enforced by *rules and code paths*, not by the storage primitive. A buggy query that forgets the `tenant_id` prefix can return cross-tenant data. Mitigations exist (security rules, a tenant-scoped client wrapper), but this is structurally weaker than Option A.
- **TR-05:** export is a recursive read of `/tenants/{tenant_id}/**` — feasible but not as clean as `pg_dump`. Requires a tenant-export service and bounded by Firestore's read-throughput quotas.
- **TR-03:** like Option A, doesn't itself solve observability storage scoping — needs a separate decision.
- Risk that a future engineer reaches around the wrapper and queries the raw Firestore client, breaking TR-01 silently.

### Option C — Per-tenant GCP project (project-per-tenant)

Each tenant gets its own GCP project. Within the project, the tenant's chosen storage primitives (Firestore, Cloud SQL, Cloud Storage) are provisioned. Platform-control-plane services in a separate platform project manage tenant projects via service accounts with scoped IAM grants.

**Pros**
- **TR-01:** strongest possible isolation — IAM boundary at the GCP project level. Cross-tenant access requires an explicit cross-project grant that doesn't exist by default. Structural isolation across all data types (state, logs, traces, secrets).
- **TR-03:** Cloud Logging and Cloud Trace are project-scoped by default — per-tenant observability scoping is free. This is the only option that solves TR-03 at the storage layer without further work.
- **TR-05:** export is "hand the tenant their project" or "bulk-export the project's resources" — cleanest possible export semantics.
- **TR-04:** per-tenant project is the natural rollout unit; updates touch one project at a time.

**Cons**
- **Operational cost & complexity (highest):** each tenant onboarding provisions a GCP project, billing account binding, IAM, VPC peering back into the platform project, and a new entry in the WireGuard topology. ADR-0003 (GCP) is preserved but the per-tenant ops surface is large.
- **TR-06:** import idempotency is now a property of whichever storage primitive is chosen *inside* the tenant project — this option defers the TR-06 question rather than answering it, so this ADR (or a sibling) would still need to decide the in-project store.
- **Quotas:** GCP project creation is API-quota-limited; scaling tenant onboarding past the quota requires support intervention.
- **Steep learning-curve and Terraform module work:** new `cloud/tenant-project/` module that orchestrates project + IAM + VPC + WireGuard plumbing; non-trivial.

## Decision Outcome

Chosen option: **{{pending — human selection required}}**, because **{{pending — rationale will be filled in once the human picks}}**.

### Consequences

* Good, because **{{pending}}**
* Bad, because **{{pending}}**
* Requires: **{{pending — Terraform module work and service-layer plumbing depend on which option is chosen}}**

### Realization

How this decision shows up in the codebase depends on the option chosen:

- **Option A (Cloud SQL per tenant):** new `cloud/cloud-sql-tenant/` Terraform module; per-tenant DSN secrets in Secret Manager; new `services/tenant-control-plane/service/` Postgres client; bedrock config keys for per-tenant connection lookup. `services/{name}/endpoint/` handlers use chi router and protobuf per CLAUDE.md.
- **Option B (Shared Firestore):** new `services/tenant-control-plane/service/` Firestore client wrapping the tenant-scoped collection root; `pkg/tenantfs/` package that enforces the path prefix and refuses queries lacking a tenant-id; security rules deployed alongside (new `cloud/firestore-tenant-rules/`).
- **Option C (Project-per-tenant):** new `cloud/tenant-project/` Terraform module; new `services/tenant-onboarding/` service that drives project creation; WireGuard topology extension; per-project service-account credentials managed by the platform.

In all three, request/response stays protobuf over HTTP per CLAUDE.md, errors flow through `pkg/errorpb`, and traffic conforms to TR-07 / ADR-0003 (Cloudflare→GCP).

## Open Questions

- **Observability storage (TR-03)** — Options A and B do not themselves answer where logs/metrics/traces live for per-tenant scoping. If the chosen option is A or B, a sibling ADR is needed for tenant-facing observability storage. Option C answers this implicitly via per-project Cloud Logging/Trace.
- **Density-vs-isolation tuning** — for Option A, do we run one Cloud SQL instance per tenant or many tenant-databases per instance? The latter weakens isolation slightly (shared compute, noisy-neighbour risk) but reduces cost; needs a follow-up decision.
- **TR-02 (multi-version contract rollout)** is not addressed by this ADR; it remains for a contract-versioning ADR. The chosen storage option needs to be checked against that ADR's requirements once drafted.
- **Migration path** — if we ever change tenant-state storage strategy later, what does superseding this ADR look like operationally? Worth thinking through before the first tenant onboards.
