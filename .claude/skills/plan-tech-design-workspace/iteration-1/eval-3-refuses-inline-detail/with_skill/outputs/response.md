Hold on — I'm not going to write the tenant-registry endpoints into `tech-design.md`, even though you know what they should be. That's `define-component-design`'s job, not mine.

Two reasons this matters:

1. **`tech-design.md` is a skeleton on purpose.** It carries the overview, component inventory, and the TR -> ADR -> component audit trail. Per-component detail (endpoint paths, request/response shapes, status codes, auth, idempotency rules, pagination) belongs in the per-component design doc because the format is type-specific. An API service's design needs endpoint contracts; a Terraform module's design needs inputs/outputs; a `pkg/` package's design needs exported types. Bundling those into one document hides the differences and makes the doc unrebuildable when one component changes.

2. **It bypasses per-component review.** If endpoints land in `tech-design.md` now, they skip the review that `define-component-design` runs (research, alternatives, contract sketch, error taxonomy alignment with `pkg/errorpb`, idempotency, etc.). "GET /tenants, POST /tenants, GET /tenants/{id}, PATCH /tenants/{id}, DELETE /tenants/{id}" is a reasonable starting point, but there are real questions a component-design pass will surface: is `GET /tenants` paginated? what's the listing filter surface? does `POST /tenants` accept an operator-supplied ID or mint one (this is one of the gaps I'd surface anyway — see below)? does `PATCH` use JSON Merge Patch, JSON Patch, or partial-protobuf field masks? what does `DELETE` mean given ADR-0003's eviction-export flow — is it the same code path? Those answers belong in the tenant-registry component design doc, where they'll be reviewed in context.

So here's what I'll do instead: file a component issue for `tenant-registry` now, and you run `define-component-design` against it next. That's one extra step, and it lands the endpoints in the right place with the right review. I won't save you that step by writing them into `tech-design.md`.

---

## What I did compose

I walked the preconditions and the ADRs in `docs/content/capabilities/self-hosted-application-platform/adrs/` (0001 tenant state storage, 0002 contract versioning, 0003 tenant eviction export) and the shared ADRs under `docs/content/r&d/adrs/`. The skeleton `tech-design.md` is drafted with overview, component inventory, data/state summary, and a TR -> ADR -> component audit trail. No per-component detail — including no tenant-registry endpoints.

### Components identified (skeleton-level only)

- **tenant-registry** (service, `services/tenant-registry/`) — established by ADR-0001; owns tenant records and the CRUD surface over them. *Endpoints will be designed in the component-design doc, not here.*
- **contract-version-catalog** (service or package, established by ADR-0002) — tracks platform contract versions tenants are pinned to.
- **tenant-eviction-exporter** (service/job, established by ADR-0003) — produces the export bundle when a tenant is evicted.
- **platform-bootstrap** (Terraform module, `cloud/platform-bootstrap/`) — provisions the platform-level dependencies (referenced from stand-up-the-platform UX).
- **shared `pkg/errorpb`** — already exists; cited because all component APIs return `application/problem+protobuf` per shared ADR-0007.

### Gaps surfaced (will be filed as gap issues, not papered over)

- **G1 — tenant ID derivation.** ADR-0001 says tenants are addressed by tenant ID and shared ADR-0006 defines the ID format, but neither specifies how a tenant ID is minted from the GitHub onboarding issue. Resolution type: per-component spec via `define-component-design` against `tenant-registry`.
- **G2 — tenant-registry API contract.** No ADR specifies the actual endpoints, request/response shapes, pagination, or PATCH semantics. Resolution type: per-component spec via `define-component-design`. *(This is exactly the gap your "while you're at it" request was poking at — the right place to resolve it is the component design doc, not inline here.)*
- **G3 — contract-version pinning storage.** ADR-0002 implies tenants carry a pinned contract version but doesn't say where it lives (tenant record? separate store?). Resolution type: per-component spec, likely on `tenant-registry`.
- **G4 — eviction-export retention/format.** ADR-0003 names an export bundle but not its format, retention, or pickup channel. Resolution type: per-component spec via `define-component-design` against `tenant-eviction-exporter`.
- **G5 — TR coverage check.** Walking the audit trail, every TR in `tech-requirements.md` traced to at least one ADR; if any did not, I would have flagged it here. (This will be re-verified at file-issue time.)

### Approval gate

Before I file anything:

> I've composed the skeleton at `tech-design.md` with 5 components and 5 gaps. Before I file issues: do the components look right? Do you want any gap resolved inline now (e.g. you want to answer G1 directly)? Once you say go, I'll file 5 component issues for `define-component-design` and 5 gap issues. The tenant-registry endpoints you mentioned will be the input to the `define-component-design` invocation against tenant-registry — they'll get reviewed in that context, not bolted onto `tech-design.md` here.

I have **not** run `gh issue create`. The intended invocations are written to `gh-invocations.txt` and will only run after you approve.

### Reminder on completeness

Tech-design.md is **not complete** until every gap issue is resolved. `plan-implementation` (Step 11) won't run against an incomplete tech design.
