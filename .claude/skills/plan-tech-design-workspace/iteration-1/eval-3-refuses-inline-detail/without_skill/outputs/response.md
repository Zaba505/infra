# Response

I composed `tech-design.md` for the self-hosted-application-platform capability, but I did **not** inline the tenant-registry's API endpoints into it. Here is why, and what I did instead.

## Why I am not inlining the endpoint list

The tech-design document is a capability-level skeleton: overview, component inventory, and a TR -> ADR -> component audit trail. Per-component detail — endpoints, request/response shapes, schemas, internal logic — is the job of `define-component-design`, which produces a separate component-design document per component. Mixing the two has real costs:

1. **It bypasses the validation step.** The component-design flow checks the proposed API surface against the capability's TRs and the relevant ADRs (here: TR-01 isolation, TR-02 contract versioning, ADR-0001's per-tenant namespace shape). Dropping `GET/POST/GET/PATCH/DELETE /tenants[/{id}]` straight into tech-design.md skips that check. For example, TR-02 requires the platform to support multiple contract versions concurrently — that almost certainly means version-bearing endpoints (e.g. a `contract-version` header, or `/v{semver}/tenants/...`). I do not want to enshrine an unversioned shape in the tech-design before that conversation happens.
2. **It conflates two review audiences.** The tech-design is reviewed for "does the capability hang together?" The component-design is reviewed for "is this component correct in isolation?" Inlining endpoints muddies both reviews.
3. **It becomes the wrong thing to maintain.** When the registry's API evolves, the source of truth should be its component-design doc, not a frozen snippet inside tech-design.md.

So your suggested endpoints are not lost — I captured them as the operator's starting proposal in the body of the tenant-registry component-design issue, explicitly marked as input to the `define-component-design` discussion rather than a settled decision. That keeps your context attached to the right artifact without short-circuiting the review.

## What I did

- Wrote `tech-design.md` with overview, accepted-ADR list, component inventory (Tenant Registry, Contract Version Resolver, Tenant Export Service, Tenant State Store), and the TR -> ADR -> component audit trail.
- Surfaced three gaps where TRs are not yet covered by accepted capability-scoped ADRs: TR-03, TR-04 (orchestration portion), TR-06.
- Wrote intended `gh issue create` invocations to `gh-invocations.txt` (not executed): one per component, one per gap. The tenant-registry issue carries your endpoint suggestion as input.

## Recommended next step

Run `define-component-design` against the tenant-registry issue. That is the right place to confirm (or refine) the `/tenants` surface against TR-01, TR-02, and ADR-0001/0002.
