# Plan ADRs — self-hosted-application-platform

**Gate check:** Tech-requirements is reviewed (`reviewed_at: 2026-04-26`, file unchanged since). Proceeding.

**Note on TR-07:** TR-07 ("inter-service communication must traverse Cloudflare → GCP") is annotated in the source as a *prior shared decision*. It does not produce a new capability-scoped ADR; it constrains the option space for the others.

**Note on shared-decision risk:** ADR #3 (tenant identifier scheme) is borderline shared — if other capabilities end up needing the same identifier, it should be promoted to `docs/content/r&d/adrs/`. Surfacing now; keeping capability-scoped because no shared identifier ADR exists yet.

## Proposed ADRs (approved by user — "I trust your enumeration, just go")

1. **Tenant isolation boundary** — what mechanism enforces isolation across data, compute, and telemetry layers. Addresses TR-01, TR-03.
2. **Tenant state storage location** — where tenant state lives such that isolation, no-downtime updates, and exportability are achievable. Addresses TR-01, TR-04, TR-05.
3. **Tenant identifier scheme** — the identifier that scopes per-tenant queries, exports, and imports. Addresses TR-03, TR-05, TR-06. *(Possibly shared — flagged.)*
4. **Platform contract versioning strategy** — how multiple contract versions are exposed concurrently and how the migration window is bounded. Addresses TR-02.
5. **Online tenant update strategy** — the deployment/cutover mechanism that delivers operator-initiated updates without tenant-visible downtime. Addresses TR-04.
6. **Tenant data export mechanism** — format, trigger, delivery, and integrity for evicted-tenant exports. Addresses TR-05.
7. **Tenant data import mechanism** — idempotent ingestion with verifiable integrity for migrations onto the platform. Addresses TR-06.
8. **Per-tenant observability data plane** — partitioning of metrics/logs/traces and the tenant query surface. Addresses TR-03.

## Filed-issue manifest (simulated; see `gh-invocations.txt` for runnable forms)

| # | ADR title | Issue |
|---|-----------|-------|
| 1 | Tenant isolation boundary | #801 |
| 2 | Tenant state storage location | #802 |
| 3 | Tenant identifier scheme | #803 |
| 4 | Platform contract versioning strategy | #804 |
| 5 | Online tenant update strategy | #805 |
| 6 | Tenant data export mechanism | #806 |
| 7 | Tenant data import mechanism | #807 |
| 8 | Per-tenant observability data plane | #808 |

Each issue body links the parent capability, lists the TR-NNs the ADR will address, states the decision to be made (no options, no rationale), and references `define-adr` as the authoring skill — one invocation per issue.
