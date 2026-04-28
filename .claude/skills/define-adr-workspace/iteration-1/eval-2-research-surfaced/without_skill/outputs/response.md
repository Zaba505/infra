# Tenant State Storage ADR — Draft

I read `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (TR-01 through TR-07), the parent capability `_index.md`, and the existing MADR ADR `0006-resource-identifier-standard.md` to match the repo's MADR 4.0.0 frontmatter and section ordering.

## Output

- `0001-tenant-state-storage.md` — drafted in `proposed` status.

## Approach

The ADR follows the repo's existing MADR layout (frontmatter with `category`, `status`, `date`, `deciders/consulted/informed`; sections: Context, Decision Drivers, Considered Options, Decision Outcome, Consequences, Confirmation, Pros and Cons, More Information). I added an explicit **Research Tasks** section (R-1 through R-6) ahead of the options so the reviewer can see what was investigated before options were proposed — covering isolation models, export/import primitives, zero-downtime upgrade behavior, reproducibility cost, topology fit, and existing repo patterns.

## Options surfaced

Each option is tied back to specific TR-NNs in its pros and cons:

- **A** — Managed GCP Firestore, namespaced per tenant
- **B** — Managed GCP Cloud SQL Postgres, database-per-tenant
- **C** — Self-hosted Postgres + GCS on home-lab via WireGuard
- **D** — Per-tenant volume, tenant-owned filesystem

## Decision

Left **pending human selection** — `Decision Outcome` is unfilled and the status is `proposed`. Per the operator-only rule in the capability definition, this is the operator's call. Once selected, `Decision Outcome`, `Consequences`, and `Confirmation` need to be filled in and `status` advanced to `accepted`.

## Notes

- No `docs/content/capabilities/self-hosted-application-platform/adrs/` directory exists in the sandbox yet; the ADR was numbered `0001` as the first capability-scoped ADR. Confirm the destination path with the operator before moving it into the docs tree (capability-scoped vs. shared `docs/content/r&d/adrs/`).
- TR-02 (multiple contract versions concurrently) and TR-03 (per-tenant observability scoping) were intentionally not used as primary drivers for *this* ADR — they belong to separate decisions (contract versioning, observability scoping) and should get their own ADRs. They are mentioned only where they incidentally constrain this choice.
