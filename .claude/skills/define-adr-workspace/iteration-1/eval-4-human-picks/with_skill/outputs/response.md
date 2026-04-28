# Response

Tech-requirements is reviewed (`reviewed_at: 2026-04-26`, file mtime 2026-04-25 — gate passes). Drafting **ADR-0001 for tenant state storage**, addressing **TR-01, TR-03, TR-04, TR-05, TR-06** (with TR-07, ADR-0003, and ADR-0006 as inherited constraints).

I'll give you my preference, but I can't set `status: accepted` on your behalf — the discipline of this skill is that the human picks. The ADR is saved as `status: proposed` with three options drafted and the Decision Outcome / Consequences left as `{{pending}}` placeholders for you to fill in once you choose.

## The three options on the table

- **Option A — Database-per-tenant on Cloud SQL Postgres.** Strongest TR-01 isolation that's still cheap-ish; clean TR-05 export (`pg_dump`); native TR-06 idempotent upsert. Doesn't solve TR-03 (observability) — needs a sibling ADR. Linear ops cost in tenant count.
- **Option B — Shared Firestore with tenant-id-prefixed paths.** Cheapest ops; matches the existing `services/machine/` Firestore house-pattern in CLAUDE.md. TR-01 is enforced by *rules and code paths*, not the storage primitive — a buggy query can leak. TR-05 export feasible but quota-bound. TR-03 still unresolved.
- **Option C — Project-per-tenant.** Strongest TR-01 (IAM boundary at the project), and the **only option that solves TR-03 for free** (Cloud Logging/Trace are project-scoped). TR-05 is "hand the tenant their project". But it has the highest ops complexity, GCP project-creation quotas come into play, and TR-06 is *deferred* into whatever in-project store gets picked.

## My preference (not a selection)

**Option A.** It cleanly satisfies TR-01, TR-04, TR-05, and TR-06 with mature, well-understood Postgres primitives, and it scales the per-tenant blast radius down to a single database per tenant — which fits the operator-update-without-downtime story (TR-04) better than Option B's "be careful with the tenant_id filter" model. I'd rather take the linear ops cost of N databases than the silent-isolation-bug risk of Option B.

The reason I'm *not* picking Option C — even though it's TR-01 and TR-03 strongest — is that it defers TR-06 (idempotent import) into a sub-decision and bakes a per-tenant GCP-project-provisioning ops surface into the platform from day one. For a learning-scale platform that's a lot of complexity to absorb up front, and the TR-03 win can be recovered later by a separate observability-storage ADR layered on Option A.

## What I need from you

1. **Pick A, B, or C** — or tell me to revise the option set (e.g. "drop C, add a hybrid: shared Firestore for control-plane state, per-tenant Postgres for application state").
2. If you pick A or B, acknowledge that a sibling ADR for tenant-facing observability storage (TR-03) is going to be filed.
3. Once you pick, I'll fill in `Decision Outcome`, `Consequences`, and the matching `Realization` block, flip status to `accepted`, and the ADR is done.

## Files written

- `outputs/adrs/0001-tenant-state-storage.md` — ADR draft (`status: proposed`).
- `outputs/adrs/_index.md` — Docsy section index for the new `adrs/` directory.
