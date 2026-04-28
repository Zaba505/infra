# Plan ADRs — self-hosted-application-platform

## Gate check

Tech-requirements is reviewed (`reviewed_at: 2026-04-26`, file last modified 2026-04-25 — review is newer than the last edit). Proceeding.

## Pushback: I won't file a single bundled ADR

You asked for one ADR called "Tenant data architecture" covering where tenant state lives, how it's partitioned, how it's accessed, and how it's backed up. I'm going to push back on that — the `plan-adrs` skill is explicit that bundling is the wrong shape, and this is a textbook case.

Bundling these four concerns into one ADR causes three concrete problems:

1. **Rationale dilution.** Each sub-decision ties back to a different set of TR-NNs. Storage location is mainly TR-01 (isolation) and TR-04 (no-downtime updates). Partitioning is TR-01 and TR-03 (per-tenant observability scope) and TR-05 (per-tenant export). Access pattern is TR-01, TR-02 (multi-version contract), and TR-07 (Cloudflare → GCP path). Backup is TR-05 (eviction export) and TR-06 (lossless import). A single ADR with a single "Decision" section can't carry four distinct rationales without becoming a small design doc — which is what `tech-design.md` is for, not an ADR.
2. **Unsupersedeable.** If we later want to change the partitioning scheme but keep the storage choice, we'd have to supersede the whole bundle and re-litigate storage and backup along with it. ADRs work best when each one is a small, replaceable unit.
3. **Hidden dependencies.** Listing them separately surfaces the order: storage location constrains partitioning, partitioning constrains access patterns, and storage + partitioning together constrain backup. Bundling hides that ordering and lets us write a "decision" that quietly assumes facts not yet decided.

## Proposed split — four ADRs

I count four decisions in your bundle. Each one is independently ADR-worthy: multiple plausible options, cross-component consequences, and clear TR-NN ties.

1. **Tenant state storage location** — Where does tenant state physically live (per-tenant database, shared database with tenant column, per-tenant GCP project, etc.)? *Addresses:* TR-01, TR-04.
2. **Tenant state partitioning scheme** — How is tenant state partitioned and identified within whatever store we pick (tenant ID as primary key prefix, schema-per-tenant, database-per-tenant, namespace-per-tenant)? *Addresses:* TR-01, TR-03, TR-05.
3. **Tenant state access pattern** — How do platform services read and write tenant state (direct DB access, a tenant-data service in front of the store, per-tenant credentials vs. platform credentials with tenant-scoped queries)? *Addresses:* TR-01, TR-02, TR-07.
4. **Tenant state backup and export approach** — How is tenant state backed up operationally and exported on eviction (continuous replication, periodic snapshots, on-demand export job; same mechanism for both or separate)? *Addresses:* TR-05, TR-06.

A natural ordering is 1 → 2 → 3 → 4: backup approach in particular depends on the storage and partitioning choices, and access pattern depends on storage. You can override the order — just tell me.

## Not yet filed

Per the skill, no issues are filed until you approve the list. **Please confirm**: file these four issues, merge any of them, drop any, or reorder them?

## Intended `gh issue create` invocations (NOT executed)

Held in `gh-invocations.txt` in this outputs folder. They will only run after you approve the split.
