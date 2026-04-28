# Response

I drafted ADR 0001 for tenant state storage, scoped to the `self-hosted-application-platform` capability. As requested, I picked the option I judged best.

## Decision

**Per-tenant isolated storage namespace** (Option A): each tenant gets its own GCS bucket(s) and its own logical database (inside a shared managed database instance), all gated by a tenant-scoped service account that is the sole non-operator principal with access. Instances are shared for cost; the IAM-scoped unit (the bucket / the logical database) is not.

## Why this option

The decision was driven by the four tech requirements that touch tenant state:

- **TR-01 (isolation)** is enforced by IAM at the resource boundary, not by application-layer `tenant_id` filtering. A query bug cannot leak cross-tenant data.
- **TR-05 (export on eviction)** becomes a generic operation: copy the tenant's bucket(s), dump the tenant's logical database. No per-tenant export tooling.
- **TR-06 (idempotent import)** targets a fresh per-tenant namespace, so retries are naturally idempotent.
- **TR-07 (Cloudflare → GCP / WireGuard topology)** is unchanged — only the IAM boundary differs per tenant.

It also respects the capability's reproducibility KPI (a single Terraform module per tenant) and its cost rule (instances are shared, so cost scales with usage rather than tenant count).

## Options considered and rejected

- **B — Shared multi-tenant store with `tenant_id` filtering.** Cheapest and simplest in resource count, but TR-01 isolation would depend on every query and every security rule being correct forever. The capability rule against any tenant observing another's state "under normal or degraded condition" forbids this failure mode.
- **C — Per-tenant managed database instance.** Strongest possible isolation, but instance-floor cost grows linearly with tenant count and provisioning time threatens the 1-hour rebuild KPI.
- **D — Home-lab-only storage with offsite backup.** Satisfies TR-07 awkwardly — GCP-resident tenant workloads would reach state over WireGuard for every read/write, making latency and availability a function of the home-lab uplink.

## Files

- `/home/carson/github.com/Zaba505/infra/.claude/skills/define-adr-workspace/iteration-1/eval-4-human-picks/without_skill/outputs/0001-tenant-state-storage.md` — the ADR (MADR 4.0.0 format, matching the `docs/content/r&d/adrs/` convention).

The ADR was numbered `0001` because it is the first capability-scoped ADR for `self-hosted-application-platform` (capability ADRs are scoped to their capability directory, separate from the shared `r&d/adrs/` series).

## Note on the requested file location

The user asked for the ADR to be saved to the outputs directory only; I did not write into the sandbox's `docs/content/capabilities/self-hosted-application-platform/adrs/` location. If you want the ADR placed there as well, the intended path would be:

`docs/content/capabilities/self-hosted-application-platform/adrs/0001-tenant-state-storage.md`
