## Response to user

Tech-requirements is reviewed (`reviewed_at: 2026-04-26`, newer than the file's last modification). The decision "tenant state storage" for the self-hosted-application-platform capability primarily addresses **TR-01** (tenant isolation), **TR-04** (no-downtime operator-initiated updates), **TR-05** (data exportability on eviction), and **TR-06** (lossless data import / migration). It is also constrained by **TR-07** (Cloudflare → GCP topology) and the prior shared ADRs (especially [0003] cloud provider = GCP, [0006] resource identifier standard, [0007] standard API error response).

Per the `define-adr` skill's Step 1, I am **not** going to draft options yet. The honest comparison between candidate storage models depends on facts I don't currently have, and inventing them would defeat the point of an ADR. Before I propose options, the following need to be answered:

### Research tasks to resolve before drafting options

1. **Existing GCP state already provisioned for this capability.** Is there an existing Firestore database / Cloud SQL instance / GCS bucket already provisioned (or referenced under `cloud/`) that this ADR should reuse, or are we provisioning fresh? `services/machine/` uses Firestore — is that the same project we'd reuse, or does the platform get its own GCP project? (This directly shapes whether "per-tenant database" vs "shared database with tenant-keyed rows" is even on the table at reasonable cost.)

2. **Realistic per-tenant export / import cost.** TR-05 (eviction export) and TR-06 (lossless import) are the hardest to satisfy with a shared schema. What is the expected upper bound on tenant data size, and what's the acceptable export-window length? If exports must complete in minutes, that argues for tenant-scoped storage units (a database, schema, or bucket-prefix per tenant). If hours-to-days is fine, a shared store with a streaming export pipeline is viable.

3. **No-downtime update model implied by TR-04.** TR-04 forbids tenant-perceived downtime during operator-initiated updates. Does "update" here include schema migrations to the tenant state store, or only application/config rollouts? If schema migrations count, options that require a global schema lock (e.g. single shared Postgres with `ALTER TABLE`) are partial-fail on TR-04 and need to be flagged as such; options with per-tenant schemas (Firestore document-per-tenant, schema-per-tenant Postgres) avoid the global lock.

4. **Compatibility with [0006] resource identifier standard.** Does the resource-identifier standard mandate a particular tenant ID format / key prefix? This shapes how tenant scoping is enforced at the storage layer (row-level vs database-level) and whether a candidate store can express the constraint natively.

5. **TR-01 isolation: "no tenant can read another's state" — at what layer is this enforced?** Strict reading: a compromised tenant workload must not be able to reach another tenant's bytes. That argues against shared-database options unless we're willing to also adopt a separate ADR on credential scoping (e.g. one IAM principal per tenant) and accept the operational cost. A separate-database-per-tenant option pushes isolation down to the cloud-provider IAM layer, which is harder to misconfigure but has provisioning-throughput limits worth checking.

6. **Tenant cardinality and provisioning rate.** How many tenants are realistic (10? 1000? 10000?), and how fast does the platform need to provision a new tenant? Per-tenant-database options that look clean at 10 tenants become operationally painful at 1000 due to GCP quota limits (e.g. Firestore database creation rate, Cloud SQL instance limits). I want a rough order of magnitude before naming products.

### What I need from you

For each of the six items above, please tell me one of:
- **(a)** Here's the answer: ...
- **(b)** Research it now (I'll do the investigation and come back with findings before drafting options).
- **(c)** Defer — file a follow-up issue and proceed with the unknown captured as an explicit assumption in the ADR.

I won't propose options against guesses. Once these are resolved (or explicitly deferred as assumptions), I'll draft at least two options, anchor each option's pros/cons in TR-01/04/05/06 and the relevant prior shared ADRs, mirror them back, and wait for you to pick. The ADR will go to `docs/content/capabilities/self-hosted-application-platform/adrs/0001-tenant-state-storage.md` with `status: proposed` until you choose.
