# Response: Push back on bundled "Tenant data architecture" ADR

I recommend **against** filing a single ADR titled "Tenant data architecture" that bundles where tenant state lives, how it's partitioned, how it's accessed, and how it's backed up. That's four distinct architectural decisions, each with its own forces, options, and reversibility profile. Bundling them produces an ADR that:

- Is hard to review (reviewers must agree to all four choices at once or block everything).
- Is hard to supersede (revisiting backup strategy later forces re-opening storage location).
- Hides the trade-offs each decision actually faces.
- Doesn't fit the MADR 4.0.0 pattern used in `docs/content/r&d/adrs/`, which is one decision per record.

## Suggested split

Reading `tech-requirements.md` (TR-01 isolation, TR-03 per-tenant query scope, TR-05 export on eviction, TR-06 idempotent import), the bundled topic naturally decomposes into four ADRs, each forced by a different cluster of TRs:

| # | ADR title | Decision | Primary TRs |
|---|-----------|----------|-------------|
| 1 | Tenant state storage system | Which datastore(s) hold tenant state (e.g., Firestore vs. per-tenant Postgres vs. object storage), and on what substrate | TR-01, TR-07 |
| 2 | Tenant data partitioning model | Logical vs. physical isolation: shared-store-with-tenant-id vs. database-per-tenant vs. project-per-tenant | TR-01, TR-03 |
| 3 | Tenant data access control | How tenant identity is established at the data layer and how queries are scoped (IAM, row-level rules, signed scopes, etc.) | TR-01, TR-03 |
| 4 | Tenant data backup, export, and import | Backup cadence, retention, restore SLO, plus the eviction-export and onboarding-import mechanisms | TR-05, TR-06 |

ADR 2 depends on ADR 1; ADR 3 depends on ADR 2; ADR 4 depends on ADRs 1 and 2. Filing them as separate issues lets each be authored, reviewed, and (if needed) superseded independently via `define-adr`.

## Recommended next step

Confirm the split above (or propose your own grouping), and I'll plan four GitHub issues — one per ADR — rather than one bundled issue. Per the sandbox instructions I will write the intended `gh issue create` invocations to `gh-invocations.txt` instead of running them.

If you genuinely want a single ADR anyway (e.g., because you've already decided this is one cohesive choice in your head), say so explicitly and I'll plan the one issue — but I want to flag the trade-off first.

## Draft gh invocations (not executed)

See `gh-invocations.txt` in this directory for the four-issue plan, ready to run once you confirm.
