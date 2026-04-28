# ADR Plan: Self-Hosted Application Platform

**Source:** `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`

The capability has 7 technical requirements (TR-01..TR-07). Below is the proposed list of capability-scoped ADRs needed to satisfy them. Each ADR is a single architectural decision; together they form the technical baseline of the platform. Cross-capability/shared decisions (e.g. the Cloudflare→GCP topology referenced by TR-07) are explicitly out of scope here and not re-decided.

---

## Proposed ADRs

### ADR-1: Tenant workload isolation boundary
- **Drives:** TR-01, TR-03 (compute/data plane half), TR-04 (update isolation)
- **Decision question:** What is the unit of isolation for a tenant's compute and in-memory state, and what mechanism enforces that no tenant can read another's runtime state?
- **Why an ADR:** TR-01 demands strict isolation but does not name a mechanism (separate VMs, separate namespaces, separate clusters, separate projects, etc.). The choice has cascading consequences for cost, reproducibility, and the operator maintenance budget KPI.
- **Out of scope here:** storage isolation (separate ADR), network isolation (separate ADR).

### ADR-2: Tenant data storage and isolation model
- **Drives:** TR-01 (data half), TR-05 (export feasibility), TR-06 (import feasibility)
- **Decision question:** How is tenant persistent data stored such that it is (a) inaccessible across tenants, (b) exportable as a portable archive, and (c) importable idempotently with integrity verification?
- **Why an ADR:** All three TRs touch the same storage substrate. Picking it once — e.g. per-tenant object-storage prefix vs. per-tenant database vs. per-tenant volume — sets the shape of both export and import.

### ADR-3: Tenant network isolation and ingress model
- **Drives:** TR-01 (network half), TR-07 (path conformance)
- **Decision question:** How are tenants reachable externally and isolated from each other on the network, while conforming to the existing Cloudflare → GCP → WireGuard → home lab path?
- **Why an ADR:** TR-07 fixes the topology but not how multiple tenants share it. Need a decision on per-tenant hostname/routing, mTLS terminations, and east-west deny-by-default.

### ADR-4: Platform contract versioning and concurrent-version support
- **Drives:** TR-02
- **Decision question:** How is the platform contract versioned, and what mechanism lets the platform run two contract versions concurrently for the migration window?
- **Why an ADR:** TR-02 forces concurrent-version support but is silent on whether that means side-by-side control planes, version-tagged tenant manifests with a shim, etc.

### ADR-5: Zero-downtime tenant update strategy
- **Drives:** TR-04
- **Decision question:** What rollout strategy (rolling, blue/green, surge) does the platform use when an operator-initiated update touches a tenant serving online traffic, such that end users observe no downtime?
- **Why an ADR:** TR-04 forbids tenant-perceived downtime but does not pick a mechanism. The choice interacts with ADR-1 (isolation unit) and ADR-3 (ingress).

### ADR-6: Per-tenant observability scoping
- **Drives:** TR-03
- **Decision question:** Where do tenant metrics/logs/traces land, and what enforces that a tenant query can only see its own scope?
- **Why an ADR:** TR-03 mandates per-tenant query scope. Need to decide whether to use a single backend with tenant-id label enforcement, per-tenant backends, or a fronting query gateway.

### ADR-7: Tenant data export format and integrity contract
- **Drives:** TR-05
- **Decision question:** What is the on-disk format, transport, and integrity guarantee (checksums, manifest) of a tenant export, and what is the bounded export window?
- **Why an ADR:** TR-05 names the requirement but not the format. ADR-2 picks where data lives; ADR-7 picks how it leaves.

### ADR-8: Tenant data import idempotency and integrity verification
- **Drives:** TR-06
- **Decision question:** What import protocol does the platform expose, and how does it guarantee idempotency on retry and detect silent loss/corruption?
- **Why an ADR:** TR-06 demands idempotent import with verifiable integrity. Decision spans manifest schema, dedup keys, and verification step.

---

## Requirements → ADR coverage matrix

| Requirement | ADR(s) |
|---|---|
| TR-01 (tenant isolation)            | ADR-1, ADR-2, ADR-3 |
| TR-02 (contract version rollout)    | ADR-4 |
| TR-03 (per-tenant observability)    | ADR-6 |
| TR-04 (zero-downtime updates)       | ADR-5 (with ADR-1, ADR-3) |
| TR-05 (export on eviction)          | ADR-7 (with ADR-2) |
| TR-06 (import existing data)        | ADR-8 (with ADR-2) |
| TR-07 (Cloudflare→GCP path)         | ADR-3 (conformance only; topology itself is a prior shared decision and is not re-decided) |

Every TR is covered. No ADR is proposed without a TR forcing it.

---

## Explicitly NOT proposed as ADRs (yet)

- **Choice of cloud provider / WireGuard topology** — already decided as a shared cross-capability decision (referenced by TR-07). Lives in `docs/content/r&d/adrs/`, not here.
- **Identity service implementation** — the capability mentions an identity service in business rules, but no TR in `tech-requirements.md` currently forces it. If/when a TR is added, an ADR follows.
- **Backup/DR mechanism** — listed as a platform output in `_index.md`, but no TR currently forces a specific decision. Add a TR first.
- **Operator-succession credential escrow** — business-rule level; not a technical decision the platform's architecture forces today.

If the user wants ADRs for any of the above, the right move is to extend `tech-requirements.md` first (via `define-technical-requirements`), then re-plan.
