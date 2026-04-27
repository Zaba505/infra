# Stage 2: Capability-scoped ADRs for self-hosted-application-platform

The tech-requirements doc has `reviewed_at: 2026-04-26` and is otherwise unchanged, so we're clear to proceed to ADRs. I drafted six MADR-format ADRs, one per non-trivial decision implied by TR-01..TR-06. TR-07 (Cloudflare → GCP path) is treated as an inherited constraint, not a decision, so it does not get an ADR of its own — it shows up as a driver in the ADRs that touch network topology.

All six live under the capability's tech-design tree, separate from the cross-cutting ADRs in `docs/content/r&d/adrs/`. They are `status: proposed` — your call to flip them to `accepted`.

## ADR index

| ADR | Title | TRs satisfied | One-line decision |
|---|---|---|---|
| 0001 | Tenant isolation model | TR-01, TR-03 | Per-tenant VM on a shared hypervisor (isolation enforced by the hypervisor, not by the platform's own policy code) |
| 0002 | Compute substrate | TR-01, TR-04, TR-07 | Home-lab KVM + GCE federated by Nomad — one schedulable fleet across home-lab and cloud |
| 0003 | Contract versioning + coexistence | TR-02 | Date-stamped contract versions (`YYYY-MM`), exactly one current + one prior supported |
| 0004 | Tenant-scoped observability | TR-03 | Shared Grafana Mimir + Loki + Tempo with per-tenant `X-Scope-OrgID`, behind a per-tenant query gateway (defence-in-depth on top of ADR-0001) |
| 0005 | Zero-downtime tenant updates | TR-04 | Rolling replacement (blue/green per tenant) with ingress drain at the Cloudflare layer |
| 0006 | Tenant data export + import format | TR-05, TR-06 | Tar archive with `manifest.json` + per-file SHA-256 + content-addressed `objects/` layout — same format both directions |

## Why these six (and not more)

- TR-01 and TR-03 collapse into one *isolation primitive* decision (ADR-0001), then a *backend* decision for observability (ADR-0004). Splitting them lets the agent-layer isolation and the query-time scoping be independent defences.
- TR-04 needs the substrate decision (ADR-0002) before the update-strategy decision (ADR-0005) makes sense. Both are non-trivial.
- TR-02 is a versioning-scheme + coexistence-mechanism decision in one (ADR-0003) — they are tightly coupled.
- TR-05 and TR-06 are one decision (ADR-0006) — same archive format, two directions.
- TR-07 is a pre-existing constraint inherited from the repo's overall topology; no decision needed at this layer.

## Cross-cutting dependency chain

```
ADR-0001 (isolation = VM)
    ├── ADR-0002 (Nomad as VM scheduler)
    │       └── ADR-0005 (Nomad update stanza for blue/green)
    └── ADR-0004 (observability scoping leverages VM identity)

ADR-0003 (contract versioning) — independent
ADR-0006 (export/import format) — references ADR-0003 for manifest schema_version
```

## Files written

Natural paths (capability tree):

- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/_index.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0001-tenant-isolation-model.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0002-compute-substrate.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0003-contract-versioning-coexistence.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0004-tenant-scoped-observability.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0005-zero-downtime-tenant-updates.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0006-tenant-data-export-and-import-format.md`

## Suggested next steps

1. Read through the six ADRs in numerical order — that follows the dependency chain.
2. The biggest decision to push back on (if you're going to push back on any) is ADR-0002: choosing Nomad over Kubernetes is the operator-budget bet. If Kubernetes is already in your bones, the maintenance-budget calculation flips.
3. ADR-0001 has a knock-on for the existing services in this repo (`services/machine/` is container-shaped, not VM-shaped) — worth confirming whether those services run *inside* tenant VMs or are platform-control-plane components that live alongside the substrate. That may want its own ADR if it's not obvious.
4. Once you've accepted/revised the ADRs, Stage 3 is the composed tech-design document that weaves them into a single human-readable design.

## Open questions surfaced (worth tracking)

- Hypervisor choice on home-lab metal (KVM was assumed; if you want Firecracker or a different hypervisor for boot speed, that's a sub-decision of ADR-0001).
- Standing capacity reserve sizing for blue/green (ADR-0005) — concrete number depends on tenant count and per-tenant VM size.
- Operator key rotation cadence for the manifest signing in ADR-0006 and the observability ingest signing in ADR-0004 — these share a key; one rotation procedure covers both, but the cadence isn't pinned.
