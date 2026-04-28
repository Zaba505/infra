# Response

Composed `tech-design.md` for the self-hosted-application-platform capability from the three accepted ADRs (0001 tenant state storage, 0002 contract versioning, 0003 tenant eviction export).

**Component inventory** (5 components):
1. Tenant State Store — realizes ADR-0001 (TR-01, TR-04)
2. Platform Contract Registry — realizes ADR-0002 (TR-02)
3. Tenant Export Service — realizes ADR-0003 (TR-05)
4. Tenant Lifecycle Controller — orchestration across all three ADRs (TR-02, TR-04, TR-05)
5. Platform Edge / Network Path — pre-existing topology (TR-07)

**Gaps identified** (TRs with no covering ADR):
- **TR-03** — tenant-facing observability scoping
- **TR-06** — tenant data import / migration in

**Follow-up issues** written to `gh-invocations.txt` (NOT executed):
- 5 component-design issues (one per component, for `define-component-design`)
- 2 ADR-gap issues (one per uncovered TR, for `define-adr` / `plan-adrs`)

All TRs (TR-01 through TR-07) appear in the audit-trail table; TR-03 and TR-06 are explicitly marked as gaps so the trail is complete and the follow-ups are traceable.
