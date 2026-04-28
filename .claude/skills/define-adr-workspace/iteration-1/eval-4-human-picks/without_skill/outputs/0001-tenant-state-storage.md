---
title: "[0001] Tenant State Storage"
description: >
    How the self-hosted-application-platform persists per-tenant state such that tenants are strictly isolated, exportable, and importable.
type: docs
weight: 1
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

The self-hosted-application-platform hosts multiple tenant capabilities (e.g. self-hosted personal media storage) on shared underlying infrastructure that spans the operator's home lab and GCP. Each tenant produces durable state — application data, configuration, secrets — that the platform must persist on the tenant's behalf.

The capability's rules and the extracted technical requirements force several non-negotiable properties on whatever storage shape we adopt:

- **TR-01:** No tenant may read another tenant's state under any normal or degraded condition.
- **TR-05:** An evicted tenant must be able to take their data with them in a portable format.
- **TR-06:** A new tenant must be able to import existing data idempotently with verifiable integrity.
- **TR-07:** All state access must conform to the existing Cloudflare → GCP / WireGuard → home-lab topology.
- **Capability §Business Rules:** the platform may span public and private infrastructure; only the operator administers it; reproducibility (KPI: ≤1h rebuild) must be preserved; cost is secondary to convenience and resiliency.

How should the platform store tenant state so that isolation, portability, and reproducibility are all simultaneously satisfied?

## Decision Drivers

- **Strict tenant isolation (TR-01).** A misconfiguration of one tenant must not be capable of exposing another tenant's data. This is the strongest driver — it rules out shared schemas with logical-only separation.
- **Portability in and out (TR-05, TR-06).** Export and import must be achievable without bespoke per-tenant tooling. The storage shape should expose a "give me everything for tenant X" boundary that maps cleanly to a portable archive.
- **Reproducibility from definitions (capability KPI).** The storage layout must be expressible as Terraform/code so the platform can be rebuilt in ≤1 hour. Snowflake per-tenant configuration is disqualifying.
- **Operator-only operation.** The blast radius of any storage operation must be small enough that one operator can reason about it. Per-tenant boundaries are easier to reason about than a shared store with policy filtering.
- **Heterogeneous tenant data shapes.** Different tenants need different storage primitives — blob/object storage for media, structured storage for catalogs and metadata. The decision must accommodate both.
- **Cost.** Secondary, but managed-database-per-tenant cost grows linearly with tenant count and must not dominate.

## Considered Options

* **Option A — Per-tenant isolated storage namespace.** Each tenant gets a dedicated set of storage resources (its own GCS bucket(s) for blobs, its own logical database — Firestore database or Cloud SQL database — for structured state) provisioned under a tenant-scoped IAM boundary. Cross-tenant access is impossible by IAM construction, not by query filter.
* **Option B — Shared multi-tenant store with tenant-scoped rows.** A single Firestore project (or single Cloud SQL instance) stores all tenants' data, with a `tenant_id` column / collection prefix and security rules that filter by the caller's identity.
* **Option C — Per-tenant managed database instance.** Each tenant gets a dedicated Cloud SQL instance (and bucket). Strongest possible isolation; closest to "give them the box."
* **Option D — Home-lab-only storage (ZFS datasets per tenant) with offsite backup.** Each tenant gets a ZFS dataset on the operator's hardware; structured state lives in a database running on the same hardware; GCS holds replicated backups only.

## Decision Outcome

Chosen option: **Option A — Per-tenant isolated storage namespace**, because it is the only option that simultaneously satisfies TR-01 isolation by IAM construction (rather than by application-layer filter), makes TR-05 export and TR-06 import trivially scoped to a tenant boundary, and remains expressible as Terraform modules so the platform stays reproducible within the 1-hour KPI. It also accepts heterogeneous data shapes — a tenant with only blobs gets only a bucket, a tenant with structured state gets a logical database — without forcing every tenant to pay for an instance per Option C.

In concrete terms, for each tenant the platform provisions:

1. One or more **per-tenant GCS buckets** (e.g. `{tenant}-blobs`, `{tenant}-exports`) with uniform bucket-level access and a tenant-scoped service account that is the only principal granted `roles/storage.objectAdmin` on those buckets.
2. A **per-tenant logical database** inside a shared managed database instance (Firestore database, or a database within a shared Cloud SQL instance) where the same tenant-scoped service account is the only principal granted access. The instance is shared for cost; the database (the IAM-scoped unit) is not.
3. A **tenant-scoped service account** that is the sole identity any tenant workload runs as. Workloads in tenant A cannot impersonate the service account of tenant B.

The shared Cloudflare → GCP / WireGuard topology (TR-07) is unaffected: tenant workloads continue to reach storage via the same network paths; only the IAM boundary changes per tenant.

### Consequences

* Good, because tenant isolation (TR-01) is enforced at the IAM/resource boundary — a query filter bug in application code cannot leak cross-tenant data.
* Good, because export (TR-05) is a `gsutil cp -r gs://{tenant}-blobs ...` plus a single-database dump, with no need to filter rows by `tenant_id`.
* Good, because import (TR-06) targets a fresh per-tenant namespace, so retries are naturally idempotent — re-running an import overwrites only that tenant's namespace.
* Good, because per-tenant resources are produced by a single Terraform module instantiated per tenant; reproducibility (≤1h rebuild) is preserved.
* Good, because heterogeneous tenants only pay for the primitives they use — a blob-only tenant gets no database.
* Neutral, because the shared database instance is a shared failure domain — an instance-level outage affects all tenants on it. Acceptable given the capability does not promise a per-tenant SLA.
* Bad, because per-tenant resource counts grow linearly with tenant count, increasing the surface area Terraform must manage. Mitigated by a single tenant module.
* Bad, because operator-side mistakes (e.g. running a maintenance script against the wrong service account) remain possible; isolation is between tenants, not between the operator and tenants.

### Confirmation

Compliance with this ADR will be confirmed through:

1. A `cloud/tenant-state-storage/` Terraform module that, given a tenant name, produces the bucket(s), logical database, and service account in one call. Onboarding a tenant must go through this module.
2. An IAM policy assertion (Terraform `google_storage_bucket_iam_policy` and the database equivalent) that lists the tenant service account as the sole non-operator principal — verified in CI by `terraform plan` review.
3. An integration test that, with tenant A's credentials, attempts to read from tenant B's bucket and database and asserts a 403/permission-denied response.
4. An export script that, given a tenant name, emits a portable archive — exercised on a fixture tenant during CI.

## Pros and Cons of the Options

### Option A — Per-tenant isolated storage namespace

Per-tenant buckets and per-tenant logical databases inside shared instances, gated by per-tenant service accounts.

* Good, because TR-01 isolation is enforced by IAM rather than by application-layer filtering — the strongest form available short of per-tenant infrastructure.
* Good, because TR-05 export and TR-06 import map to a single tenant boundary, so portability tooling is generic across tenants.
* Good, because reproducibility is preserved: one Terraform module instantiated per tenant.
* Good, because cost scales with usage, not with tenant count — shared instances amortize fixed costs.
* Neutral, because the shared database instance is a shared failure domain.
* Bad, because the tenant-resource count grows linearly with tenant count.

### Option B — Shared multi-tenant store with tenant-scoped rows

A single store, with `tenant_id` filtering and security rules.

* Good, because resource count is constant regardless of tenant count.
* Good, because cost is minimized — one bucket, one database.
* Bad, because TR-01 isolation depends on every query and every security rule being correct forever — a single bug leaks all tenants. This is the failure mode the capability's "no tenant may observe another's state under any normal or degraded condition" rule explicitly forbids.
* Bad, because TR-05 export requires writing a tenant-aware filter for every collection/table, which must be maintained as the schema evolves.
* Bad, because TR-06 idempotent import requires the import path to enforce `tenant_id` scoping on every write — yet another place for a bug to leak data.
* Bad, because operator mistakes have global blast radius — `DELETE FROM table` is a platform-wide event.

### Option C — Per-tenant managed database instance

Each tenant gets its own Cloud SQL instance (or equivalent) plus its own bucket.

* Good, because isolation is even stronger than Option A — a noisy neighbor cannot affect a tenant's database performance.
* Good, because TR-05 export is a full instance dump; TR-06 import is a full instance restore.
* Good, because tenants have no shared failure domain at the database level.
* Bad, because cost grows linearly with tenant count at instance granularity — a Cloud SQL instance has a non-trivial floor price even at zero load. This conflicts with the capability rule that cost should still be minimized where it does not cost convenience or resiliency.
* Bad, because reproducibility is harder: provisioning a Cloud SQL instance regularly takes longer than the 1-hour rebuild KPI allows when several tenants exist.
* Bad, because tenants with only blob needs are paying for a database they do not use.

### Option D — Home-lab-only storage (ZFS datasets per tenant) with offsite backup

Tenant state lives on operator-owned hardware; GCS holds replicated backups.

* Good, because the operator has total control of the medium.
* Good, because per-tenant ZFS datasets give strong isolation, snapshotting, and per-tenant export by design.
* Bad, because TR-07 is satisfied only awkwardly — tenant workloads running in GCP would have to reach storage in the home lab over WireGuard for every read/write, making latency and availability a function of the home-lab uplink.
* Bad, because reproducibility is harder — ZFS layout is not as cleanly described in code as GCP resources are.
* Bad, because a total loss of the home lab makes the platform unavailable until restore from backup completes, where Option A keeps state in GCP and the home lab is one of several reachable surfaces.
* Neutral on cost, depending on hardware vs. cloud-storage prices.

## More Information

### Tenant boundary, concretely

For a hypothetical tenant `media`:

```
gs://media-blobs              # tenant-owned blob storage
gs://media-exports            # tenant-owned export staging
firestore:(default)/media     # tenant-owned logical Firestore database
sa: media-tenant@<project>.iam.gserviceaccount.com  # sole tenant identity
```

Tenant `media`'s workloads run as `media-tenant@…`. That service account is the only non-operator principal with access to any of the four resources above. No other tenant's service account has any grant on them.

### Mapping back to TRs

| Requirement | How this ADR satisfies it |
|---|---|
| TR-01 (isolation)         | IAM-level boundary on per-tenant bucket and per-tenant logical database; cross-tenant access is structurally impossible |
| TR-05 (export on eviction)| Export = copy tenant bucket(s) + dump tenant database; both are scoped to the tenant by construction |
| TR-06 (idempotent import) | Import targets a fresh tenant namespace; retries overwrite only that namespace |
| TR-07 (network topology)  | Storage primitives are GCP-native and reached via the existing topology; unchanged |

### References

- Capability: [self-hosted-application-platform](../_index.md)
- Tech requirements: [tech-requirements.md](../tech-requirements.md)
- Related shared ADR: [0006 Resource Identifier Standard](../../../r&d/adrs/0006-resource-identifier-standard.md) for naming of per-tenant resources
