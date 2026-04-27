# Pushing back on the single-ADR shape

Before I draft anything, I want to flag a concern with the framing. "Tenant data architecture" as one ADR bundles four genuinely independent decisions, and I think we'll regret it within a quarter. Let me explain, then propose a split.

## Why one ADR is the wrong unit here

You listed four concerns:

1. **Where tenant state lives** (storage substrate — Firestore? Postgres-per-tenant? object storage? a mix keyed by data class?)
2. **How it's partitioned** (tenant-per-database? shared schema with a tenant_id discriminator? namespace-per-tenant?)
3. **How it's accessed** (direct client SDK? a data-access service in front? per-tenant credentials vs. platform-held credentials with row-level enforcement?)
4. **How it's backed up** (snapshot cadence, retention, restore RTO/RPO, cross-region copies, export format reuse for TR-05)

Each of these has its own forcing functions in the tech requirements:

- **Partitioning** is driven primarily by TR-01 (isolation invariant) and TR-03 (per-tenant query scope). The decision pivots on how strong an isolation boundary we need at the storage layer — physical, logical, or cryptographic.
- **Access** is driven by TR-01 and TR-07 (Cloudflare→GCP path). The decision pivots on whether tenants ever touch storage directly or always through a platform service, and on credential blast radius.
- **Storage substrate** is driven by the workload shape (online vs. analytical, structured vs. blob), TR-04 (zero-downtime updates constrain schema migrations), and TR-06 (import idempotency favors substrates with strong upsert semantics).
- **Backup** is driven by TR-05 (evicted-tenant export) and operational durability concerns that aren't even in the tech-requirements doc yet — RTO/RPO, geographic redundancy, restore drills.

These four pivot on **different requirements**, have **different reversibility profiles** (you can re-shape backup any time; you cannot easily re-partition a populated multi-tenant store), and will be **revisited on different cadences**. MADR's value is that each decision carries its own context, options-considered, and consequences. Bundling them means:

- The "Consequences" section becomes a soup — readers can't tell which trade-off came from which choice.
- When we later revisit, say, partitioning (because a tenant outgrows the shared model), we'd have to supersede an ADR that also locked in backup cadence, which is unrelated and still correct. That muddies the provenance trail the tech-requirements doc is trying to preserve.
- Review is harder. Four reviewers with four specialties can't sign off on their slice independently.
- We lose the ability to record "deferred" cleanly. Backup specifics, in particular, depend on RTO/RPO numbers we haven't elicited yet — I don't want to either block the storage decision on that or quietly hand-wave it inside a combined doc.

## What I'd propose instead

Four ADRs, capability-scoped, that can be authored and reviewed in parallel but have a clear dependency order:

1. **ADR-001: Tenant data partitioning model** — addresses TR-01, TR-03. Decides the isolation boundary (physical DB per tenant vs. logical namespace vs. shared+discriminator). This one is **least reversible** and should be decided first because the others key off it.
2. **ADR-002: Tenant state storage substrate(s)** — addresses TR-04, TR-06, plus workload fit. Decides what technology actually holds tenant state, possibly per data class (config vs. application data vs. blobs). Depends on ADR-001 because the partitioning model constrains substrate choice (e.g., "DB-per-tenant" rules out substrates with painful provisioning overhead).
3. **ADR-003: Tenant data access path** — addresses TR-01, TR-07. Decides whether tenants and platform services reach storage directly or through a data-access service, and how credentials/authz are scoped. Depends on ADR-001 and ADR-002.
4. **ADR-004: Tenant data backup and restore** — addresses TR-05 and the operational durability gap. Depends on ADR-002 (backup mechanisms are substrate-specific) and can also feed the export mechanism TR-05 requires, so it should explicitly cross-reference the move-off UX.

That gives us four crisp documents, each ~1–2 pages, each citing the specific TR-NNs it resolves, each independently supersedable.

## What I need from you to proceed

A few things before I start drafting:

- **Are you OK with the four-ADR split?** If you have a strong reason for the bundle (e.g., you want one review pass, or this is for an audience that wants narrative not decision records), tell me and we can compose a tech-design doc that *reads* as one story while the underlying ADRs stay split. That's actually the intended shape: ADRs are the atomic decisions; the tech-design doc is the human-friendly composition over them.
- **Open question on TR-05**: the tech-requirements doc flags "on-demand vs. continuously-available export" as open. Backup design depends on the answer. Do you want to resolve that now, or should ADR-004 record it as a deferred sub-decision and pick the more flexible option?
- **RTO/RPO targets**: not in tech-requirements yet. Without numbers, ADR-004 will be hand-wavy. I'd suggest we add a TR-08 (durability/recovery objectives) to the requirements doc before I draft ADR-004, even if the numbers are rough.
- **Workload classes**: do we already know whether tenants store (a) only platform-managed config, (b) arbitrary application state, (c) blobs/large objects, or some combination? This shapes ADR-002 heavily and I don't see it in the requirements doc.

If you want me to just push ahead with the four ADRs using my best read of the requirements and mark the gaps as deferred, I can do that — but I'd rather spend ten minutes nailing those four questions than write four ADRs that immediately need superseding.

Which way do you want to go?
