# Response

I want to push back before we write this one. "Tenant data architecture" as a single ADR sounds tidy, but it bundles four genuinely independent decisions, and MADR works against us when we do that.

## Why I'd split it

Each of the four sub-topics you listed has its own decision drivers, its own option set, and its own blast radius if we revisit it later:

1. **Where tenant state lives** — driven by TR-01 (isolation) and TR-07 (Cloudflare→GCP topology). Options: Firestore in the same GCP project, per-tenant GCP project, home-lab Postgres reachable over WireGuard, etc. This is fundamentally a *placement and trust-boundary* decision.

2. **How it's partitioned** — driven by TR-01 and TR-03 (per-tenant observability scoping). Options: shared schema with tenant-id row filter, schema-per-tenant, database-per-tenant, project-per-tenant. This is a *blast-radius vs. operational-cost* tradeoff and is largely orthogonal to where the data physically lives — you can pick "Firestore" and still independently choose "collection-per-tenant" vs. "tenant-id field."

3. **How it's accessed** — driven by TR-01, TR-04 (zero-downtime updates), and TR-07. Options: direct client access with row-level security, gateway service that injects tenant scope, sidecar proxy. This is an *enforcement-point* decision.

4. **How it's backed up** — driven by TR-05 (evicted-tenant export) and TR-06 (migration-in integrity). Options: managed snapshot + retention, per-tenant logical export pipeline, continuous CDC to object storage. This is a *durability and portability* decision and shares almost nothing with the access-control decision.

If we bundle them:

- The MADR sections collapse. "Decision Drivers" becomes a grab-bag of all seven TRs, "Considered Options" becomes a cartesian product, and "Consequences" becomes unreadable.
- Provenance breaks. When we later change the backup strategy (likely — TR-05's open question about on-demand vs. continuously-available is unresolved), we'd be amending an ADR that also covers partitioning and access. The git history stops telling a clean story.
- The tech-design composition step downstream gets messier — it can't cite "ADR-NNN §backup" cleanly because the ADR doesn't have crisp sub-decisions; it has one mega-decision.

## What I'd propose instead

Four ADRs, each capability-scoped under `docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/`:

- `tenant-state-placement.md` — answers "where does tenant state live?" (drivers: TR-01, TR-07)
- `tenant-data-partitioning.md` — answers "how is tenant state partitioned?" (drivers: TR-01, TR-03)
- `tenant-data-access-pattern.md` — answers "how do tenants and platform services read/write tenant state?" (drivers: TR-01, TR-04, TR-07)
- `tenant-data-backup-and-export.md` — answers "how is tenant state durable and how do tenants get it back out?" (drivers: TR-05, TR-06)

These can still be written as a sequence in one sitting — they just stay as separate decision records so each one's drivers, alternatives, and consequences stay crisp.

## What I need from you

Before I draft anything, can you confirm one of:

- **(A)** Yes, split into the four ADRs above — start with placement, since the others depend on it.
- **(B)** You really want one ADR; here's why bundling is acceptable for your use case.
- **(C)** A different split (e.g., merge access + partitioning, keep backup separate).

I'd lean (A) but I want your call before I commit anything to the tech-design folder. I haven't written any ADR file yet — only the reviewed `tech-requirements.md` is in place.
</content>
</invoke>