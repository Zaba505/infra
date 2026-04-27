# Stage 2: pushing back on bundling

I've read `tech-requirements.md` (reviewed_at: 2026-04-26 — review gate satisfied), the capability `_index.md`, and the seven UX docs. We're cleared for Stage 2.

But I don't want to write a single "Tenant data architecture" ADR. Bundling those four sub-decisions into one ADR dilutes the rationale for each, makes it impossible to supersede one without revisiting the others, and breaks the chain of evidence that ties each decision back to specific TRs. Per the skill's discipline: one decision per ADR.

## What I see inside your proposed ADR

You named four sub-decisions, and they're forced by different TRs with different option spaces:

1. **Where tenant state lives** (storage substrate selection) — forced primarily by TR-01 (isolation), TR-03 (per-tenant queryability), TR-07 (must live within the Cloudflare→GCP topology). Options include Firestore, Cloud SQL/Postgres, per-tenant database instances, etc. Independent rationale.

2. **How tenant state is partitioned** (isolation model) — forced by TR-01 (no cross-tenant reads) and TR-03 (per-tenant scoping). Options include per-tenant database, per-tenant schema, shared schema with row-level tenant_id, namespace-per-tenant. This is an isolation-strength vs operational-cost tradeoff that is *independent* of the storage substrate (you can pick any partitioning model on top of most substrates).

3. **How tenant state is accessed** (access path / authorization) — forced by TR-01 (enforcement point), TR-03 (per-tenant query scoping), TR-07 (network path). Options include direct DB access from tenant workloads, a platform data-plane service that brokers all access, sidecar-mediated access. This decision *consumes* the partitioning model but doesn't determine it.

4. **How tenant state is backed up** (backup, restore, and export) — forced by TR-05 (evicted-tenant export), TR-06 (lossless import), and indirectly TR-04 (no-downtime updates often need point-in-time snapshots). Options include managed-service snapshots, logical export pipelines, continuous CDC to object storage, per-tenant on-demand export jobs. Backup also has an open question already noted in `tech-requirements.md` (on-demand vs continuously-available export) that needs to be resolved here, not buried in a four-headed ADR.

Each of these four can be superseded independently as we learn things. Each has a distinct option set. Each cites a distinct (though overlapping) TR subset.

## What I'd like to do instead

Draft four ADRs, in this order, accepting each before starting the next:

- `0001-tenant-state-storage-substrate.md` — TR-01, TR-03, TR-07
- `0002-tenant-state-partitioning-model.md` — TR-01, TR-03 (depends on 0001)
- `0003-tenant-state-access-path.md` — TR-01, TR-03, TR-07 (depends on 0001, 0002)
- `0004-tenant-state-backup-and-export.md` — TR-05, TR-06, TR-04 (depends on 0001; resolves the open question on export modality)

A couple of these may flag as cross-capability once we get into options — in particular TR-07 means the access-path ADR may touch shared topology, and the partitioning model may interact with the resource-identifier standard. If so I'll surface and stop per the skill's flag-and-stop rule, not silently absorb them.

## Decision needed from you

Confirm one of:

1. Proceed with the four-ADR split above, starting with `0001-tenant-state-storage-substrate`.
2. A different split (e.g. you want backup folded into the substrate ADR because they're tightly coupled in your mind — I'd push back, but that's a conversation we can have).
3. You still want one bundled ADR — in which case I'll need you to override the skill's discipline explicitly, and I'll record that override in the ADR's context so future readers know why the rationale is tangled.

I'll wait for your call before writing any file.
