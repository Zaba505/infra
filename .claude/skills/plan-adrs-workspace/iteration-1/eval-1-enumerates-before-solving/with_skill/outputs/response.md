# Plan ADRs — self-hosted-application-platform

**Gate check:** `tech-requirements.md` is reviewed (`reviewed_at: 2026-04-26`, file unchanged since `2026-04-25`). Proceeding with enumeration.

I read all seven TRs end-to-end. Before I propose anything, two callouts:

- **TR-07 ("All inter-service communication must traverse the Cloudflare → GCP path")** is marked *prior shared decision* in the source. I am **not** including it in the capability-scoped list — it belongs in `docs/content/r&d/adrs/` and is already constrained by an existing shared decision. It will appear as a constraint inside several of the ADRs below, not as one of them.
- I am **not** proposing options or answers. Each item below is a *decision to be made*, sized so that one invocation of `define-adr` can author it.

## Proposed ADR list (for your approval)

1. **Tenant state storage** — where a tenant's business state physically lives (datastore choice, account/project boundary). Addresses **TR-01** (isolation at data layer), **TR-04** (state must be reachable through a no-downtime update), **TR-05** (must be exportable on eviction), **TR-06** (must accept idempotent imports).

2. **Tenant state partitioning scheme** — how tenant data is keyed/partitioned/segregated within whatever store ADR 1 picks (per-tenant DB vs. shared schema with row-level keys vs. per-tenant namespace). Addresses **TR-01**, **TR-03** (observability data must be queryable *per tenant within their scope only* — partitioning model has to make that enforceable), **TR-05** (export boundary follows the partition).

3. **Tenant compute isolation boundary** — what unit of compute is a tenant boundary (pod, namespace, project, cluster) and how cross-tenant access is denied at runtime. Addresses **TR-01** (isolation at compute layer), **TR-04** (the unit being updated without downtime is this same boundary).

4. **Platform contract versioning & concurrent-version strategy** — how multiple contract versions coexist, how a tenant is pinned to a version, and how the migration window is bounded. Addresses **TR-02**.

5. **Tenant update rollout strategy** — how an operator-initiated update (config, version, capability) is applied to a running tenant without end-user-visible downtime (blue/green, rolling, traffic-shifting, drain semantics). Addresses **TR-04**, with **TR-02** as a constraint (an in-flight contract migration must not break this).

6. **Tenant-facing observability data scoping** — how metrics/logs/traces are tagged, stored, and queried such that a tenant sees only their own data. Addresses **TR-03**, with **TR-01** as a constraint (telemetry is also tenant state for isolation purposes).

7. **Tenant data export mechanism & format** — how an evicted tenant retrieves all of their data (trigger, format, integrity guarantee, export window). Addresses **TR-05**, with **TR-01**/**TR-02** as constraints (export must not violate isolation; format should be contract-version-aware).

8. **Tenant data import mechanism & integrity verification** — how pre-existing tenant data is imported idempotently with verifiable integrity. Addresses **TR-06**. Kept separate from ADR 7 deliberately — import and export have different failure modes (import has to handle duplicate retries and partial loads; export does not), and bundling them produces an ADR you cannot supersede on one side without revisiting the other.

## Dependency notes (for ordering, not for this skill to decide)

- ADR 1 (storage) constrains ADR 2 (partitioning), which constrains ADR 6 (observability scoping), ADR 7 (export), and ADR 8 (import).
- ADR 3 (compute boundary) is independent of 1/2 and can run in parallel.
- ADR 4 (contract versioning) is independent of storage/compute and can run in parallel; it becomes a *constraint* on ADR 5 and ADR 7.
- A reasonable order: **1 → 2 → {3, 4} in parallel → 5 → 6 → {7, 8} in parallel.** Final ordering is your call.

## Things I deliberately did **not** include

- **TR-07 / Cloudflare → GCP topology.** Prior shared decision; not re-decided here.
- **Inter-service auth / identity.** If you want this decided for the platform, it is almost certainly a shared ADR (it touches every capability hosted on the platform), so it belongs in `docs/content/r&d/adrs/` via a separate flow. Flag-and-stop: tell me if you want to defer it or pause to handle it separately.
- **Specific datastore vendor, specific compute platform, specific telemetry backend.** Those are *options* under the ADRs above and belong to `define-adr`.

## What I need from you

Reply with one of:

- **"Approved as-is"** — I will file 8 issues via `gh issue create`, one per ADR.
- **"Approved with changes: …"** — add, remove, merge, split, or reorder.
- **"Pause — handle the inter-service auth shared ADR first"** — I stop here.

Until you approve, **no issues are filed.** The intended `gh issue create` invocations have been pre-staged in `gh-invocations.txt` for review but have **not** been executed.
