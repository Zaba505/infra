---
title: "[0001] Tenant State Storage"
description: >
    Select the storage substrate the self-hosted-application-platform offers tenants for their persistent state, such that tenant isolation, exportability, importability, and zero-downtime updates are all satisfiable.
type: docs
weight: 1
category: "strategic"
status: "proposed"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions (frameworks, auth strategies, cross-cutting patterns)
- user-journey: Solutions for specific user journey problems (feature implementation approaches)
- api-design: API endpoint design decisions (pagination, filtering, bulk operations)
-->

## Context and Problem Statement

The self-hosted-application-platform must provide each tenant capability with **persistent storage** for its state (per the capability's *Outputs & Deliverables*). That storage is the substrate against which several of the platform's hardest technical requirements get enforced or violated: tenant isolation (TR-01), per-tenant observability scoping (TR-03), zero-downtime operator-initiated updates (TR-04), portable export on eviction (TR-05), and verifiable idempotent import (TR-06). Whatever we pick will also have to coexist with the existing Cloudflare → GCP topology with WireGuard back to the home lab (TR-07).

The platform itself runs partly in GCP and partly in the operator's home lab (per the repo's stated architecture and the capability's "may span public and private infrastructure" rule). The choice of storage substrate determines where tenant state physically lives, how it is isolated, how it is moved in and out, and how it is upgraded under live traffic.

What storage substrate should the platform offer tenants for their persistent state?

## Decision Drivers

* **TR-01** — Strict tenant isolation: no tenant may read another's state under any normal or degraded condition.
* **TR-04** — Operator-initiated updates must be zero-downtime for online workloads, which means the storage layer cannot require a tenant outage to upgrade or migrate.
* **TR-05** — Evicted tenants must be able to take their data with them in a portable format within a defined window.
* **TR-06** — New tenants must be able to import pre-existing data idempotently with verifiable integrity.
* **TR-07** — Must fit the existing Cloudflare → GCP + WireGuard topology; new bespoke ingress paths are out of scope for this decision.
* **Reproducibility KPI** — The platform must be re-standable from definitions in ≤1 hour. The storage substrate must be expressible as definitions, not as a hand-built snowflake.
* **Operator maintenance budget** — ≤2 hours/week routine maintenance across the whole platform. The storage substrate's share of that budget must be small.
* **Cost is secondary to convenience and resiliency**, but not free — wildly more expensive options need to buy meaningful resiliency or convenience.
* **Operator-only operation** — there are no co-operators; whatever we pick must be operable by one person with sealed/escrowed successor credentials.

## Research Tasks

The following items were identified as needing investigation before options could be weighed; their findings are folded into "Pros and Cons of the Options" below. They are surfaced here so the human reviewer can see what was (and was not) checked.

* **R-1: Per-tenant isolation models in managed GCP storage** — how Firestore (already a `go.mod` dependency, used by `services/machine`), Cloud SQL, and GCS enforce tenant separation (database-per-tenant vs. row-level vs. bucket-per-tenant), and how that interacts with TR-01 under a degraded/misconfigured state.
* **R-2: Export/import primitives** — what each option natively offers for full-fidelity export (TR-05) and idempotent import with integrity verification (TR-06). For portable formats, whether export can be triggered without operator involvement (per the capability's "on-demand exportable archives" rule).
* **R-3: Zero-downtime upgrade behavior** — for each option, what an operator-initiated tenant update (config / version / capability change, TR-04) looks like at the storage layer: schema migrations, version skew, online vs. offline.
* **R-4: Reproducibility cost** — how much of each option is expressible in Terraform (consistent with the existing `cloud/*` module pattern) vs. requiring imperative bootstrap. Bears directly on the 1-hour standup KPI.
* **R-5: Topology fit** — whether the option requires network paths outside the established Cloudflare → GCP + WireGuard topology (TR-07). Self-hosted databases on home-lab hardware are reachable via WireGuard; managed GCP services are reachable from home-lab via the same path; anything else is a new ingress.
* **R-6: Existing repo patterns** — what storage the repo already uses. `services/machine` uses Firestore via `service.FirestoreClient`, and `pkg/errorpb` is wired for `application/x-protobuf` over HTTP. A choice that reuses these patterns has lower marginal cost than one that introduces a new substrate.

## Considered Options

* **Option A — Managed GCP Firestore, namespaced per tenant** (one Firestore database or top-level collection per tenant, in the operator's existing GCP project)
* **Option B — Managed GCP Cloud SQL (PostgreSQL), database-per-tenant** (one logical database per tenant on a shared Cloud SQL instance, or one instance per tenant for stricter isolation)
* **Option C — Self-hosted Postgres + GCS, on the home-lab cluster, reached via WireGuard** (operator-run Postgres for relational state, GCS bucket-per-tenant for blobs; both fronted by the platform's own data-access service)
* **Option D — Per-tenant volume on the home-lab cluster, exposed as a tenant-owned filesystem** (the platform gives each tenant a dedicated persistent volume and lets the tenant choose its own storage engine on top)

Each option must be evaluated against TR-01, TR-04, TR-05, TR-06, TR-07 explicitly; that mapping is included in each option's pros/cons below.

## Decision Outcome

**Pending human selection.** This ADR is in `proposed` status; the operator is the sole decider on this capability and has not yet picked an option. The trade-offs are captured below so the operator can make the call. Update this section with the chosen option, the rationale, and any conditions before transitioning to `accepted`.

### Consequences

To be filled in once an option is selected. Common to all options: tenant onboarding will gain a "provision tenant state" step; tenant offboarding will gain an "export then revoke" step; the platform contract (per TR-02) will reference whichever substrate is chosen, so changing it later is a contract-version migration, not a silent swap.

### Confirmation

To be filled in once an option is selected. Confirmation will minimally require:

* A reproducibility test: the chosen substrate is brought up from definitions as part of the ≤1-hour platform standup KPI.
* An isolation test (TR-01): a tenant credential, by construction, cannot read another tenant's state — verified at the substrate level, not just at the application layer.
* An export/import round-trip test (TR-05 + TR-06): export a tenant, re-import the export into a fresh tenant slot, verify integrity, and verify idempotency on retry.
* A zero-downtime update drill (TR-04): an operator-initiated update on a live tenant produces no end-user-visible downtime.

## Pros and Cons of the Options

### Option A — Managed GCP Firestore, namespaced per tenant

Each tenant gets its own Firestore database (or, at minimum, a top-level collection scoped by tenant ID) within the operator's existing GCP project. Access is gated by per-tenant service accounts whose IAM bindings only see that tenant's database/collection.

* Good (TR-01), because GCP IAM at the database level is enforced by GCP, not by application code — a misconfigured tenant service still cannot read another tenant's database.
* Good (TR-04), because Firestore is a managed service with no tenant-visible upgrade window; operator-initiated tenant updates do not require a storage-layer outage.
* Good (TR-07), because Firestore is reached over the existing GCP path; no new ingress is introduced.
* Good (Reproducibility KPI), because Firestore databases and IAM bindings are Terraform-expressible, fitting the existing `cloud/*` module pattern.
* Good (Operator maintenance budget), because there is nothing to patch, back up, or capacity-plan at the storage layer.
* Good (repo fit), because Firestore is already a dependency and there is a working `services/machine` example to copy.
* Neutral (TR-05), because Firestore export to GCS is a built-in operation but produces a Firestore-specific format; "portable" requires an additional translation step the platform must own.
* Neutral (TR-06), because import is supported but not idempotent by default; the platform would need to wrap import with an idempotency key and integrity check.
* Bad, because tenants whose data model is fundamentally relational (joins, transactions across many entities) get a worse fit than a SQL store.
* Bad, because the operator becomes more dependent on a single GCP service for the platform's most critical asset (tenant state); this weakens "independence from hosting vendors" at the platform level, even if the operator retains the ability to leave.

### Option B — Managed GCP Cloud SQL (PostgreSQL), database-per-tenant

A shared Cloud SQL instance (or one instance per tenant for stricter blast-radius control) where each tenant gets its own logical database with its own role. The platform's data-access layer brokers connections.

* Good (TR-01), because database-per-tenant with per-tenant roles is the canonical Postgres isolation pattern; row-level isolation is not relied upon.
* Good (TR-05), because `pg_dump` produces a portable, well-understood format; export tooling is mature and tenant-runnable in principle.
* Good (TR-06), because `pg_restore` plus a transactional staging schema gives idempotent import with integrity verification (checksums + transactional boundary).
* Good (TR-07), because Cloud SQL is reached over the existing GCP path; private IP + WireGuard is well-trodden.
* Good, because relational guarantees (transactions, foreign keys) are available to tenants that need them.
* Neutral (TR-04), because Cloud SQL maintenance windows can be configured to be short and during low-traffic hours, but they still exist; truly zero-downtime requires read replicas or a connection-multiplexing layer the platform would own.
* Neutral (Reproducibility KPI), because instance creation is Terraform-expressible but instance bring-up time can approach the 1-hour budget on cold start; a shared instance amortizes this.
* Neutral (Operator maintenance budget), because Cloud SQL is managed but not zero-touch — version EOLs, parameter tuning, and storage growth still land on the operator.
* Bad, because cost grows roughly linearly with tenants if instance-per-tenant is chosen for stronger isolation; shared-instance is cheaper but couples tenants at the resource level.
* Bad, because key-value or document-shaped tenant data gets a relational impedance mismatch that Firestore would not have.

### Option C — Self-hosted Postgres + GCS, on the home-lab cluster, reached via WireGuard

The platform runs Postgres on home-lab hardware for relational state and uses bucket-per-tenant GCS for blobs. The platform's data-access service brokers all access; tenants never address either substrate directly.

* Good (vendor independence), because the relational substrate runs on operator-controlled hardware; "independence from hosting vendors" is maximally preserved for the most critical asset.
* Good (TR-05), because `pg_dump` + GCS-native export gives portable, well-understood export per data class.
* Good (TR-06), because the same `pg_restore` + transactional staging story applies; GCS supports idempotent object writes via object versioning + content hash.
* Good (cost), because the marginal cost per tenant on existing home-lab hardware is near zero until capacity runs out.
* Neutral (TR-07), because the home-lab path is via the existing WireGuard tunnel — no new ingress — but the platform now has a hard dependency on that tunnel for tenant data access, which raises the failure-domain stakes of the tunnel itself.
* Neutral (TR-01), because isolation is enforceable via database-per-tenant + bucket-per-tenant + per-tenant credentials, but the *operator* is now responsible for making sure that's true (vs. GCP enforcing it). One misconfigured role grants cross-tenant reads.
* Bad (TR-04), because operator-initiated updates against a self-run Postgres are inherently more outage-prone than against a managed service; achieving zero tenant-visible downtime requires the operator to run replicas and orchestrate failover, which eats into the maintenance budget.
* Bad (Operator maintenance budget), because patching, backup verification, capacity planning, and disaster recovery for Postgres all become operator work — this is the option most likely to push past 2h/week.
* Bad (Reproducibility KPI), because rebuilding a self-hosted stateful service from definitions in ≤1 hour requires a working backup-restore pipeline that itself must be reproducible; this is achievable but doubles the surface area to keep reproducible.

### Option D — Per-tenant volume on the home-lab cluster, exposed as a tenant-owned filesystem

The platform hands each tenant a dedicated persistent volume on home-lab hardware and lets the tenant pick its own storage engine on top (sqlite, embedded KV, files, whatever).

* Good (vendor independence), same as Option C.
* Good (cost), same as Option C.
* Good (tenant flexibility), because each tenant chooses the engine that fits its data model with no impedance mismatch.
* Neutral (TR-01), because volume-level isolation is straightforward, but with tenants running their own engines on their own volumes the platform has the *least* visibility into whether a tenant has accidentally exposed itself.
* Bad (TR-05), because "portable export" now means "give the tenant their volume contents in some format" — there is no platform-owned, uniform export format; each tenant's export is shaped like its engine, and the platform has to trust each tenant to have an export story.
* Bad (TR-06), because idempotent import with verifiable integrity is the tenant's responsibility per engine; the platform cannot give a uniform guarantee.
* Bad (TR-04), because zero-downtime updates against arbitrary tenant-chosen engines cannot be promised by the platform; it would have to push the requirement back onto each tenant, contradicting the capability's "platform evolves with its tenants" rule.
* Bad (Operator maintenance budget), because supporting N storage engines is N times the operational surface vs. picking one.

## More Information

* TR-01, TR-04, TR-05, TR-06, TR-07 — `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`
* Capability definition — `docs/content/capabilities/self-hosted-application-platform/_index.md` (especially *Outputs & Deliverables*, *Business Rules*, and *Success Criteria & KPIs*)
* Existing Firestore usage — `services/machine/service/` (per repo CLAUDE.md)
* Existing Terraform module pattern — `cloud/*` (GCP provider v7.11.0, `create_before_destroy`)
* Topology — Internet → Cloudflare → Home Lab ↔ GCP (WireGuard); per repo CLAUDE.md and TR-07
