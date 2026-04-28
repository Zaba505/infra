# Response to user

Tech-requirements is reviewed (`reviewed_at: 2026-04-26`). Drafting ADR-0001 for the `tenant-control-plane` service's API exposure, addressing **TR-01** (tenant isolation), **TR-02** (multi-version contract), **TR-03** (per-tenant observability scope), **TR-04** (no-downtime updates), and **TR-07** (Cloudflare → GCP topology).

## Inherited constraints (not up for re-decision in this ADR)

These come from prior shared ADRs and `CLAUDE.md` and bound the option set:

- **ADR-0003 (cloud provider) + ADR-0007 (error response):** services run on GCP and emit `application/problem+protobuf` errors via `pkg/errorpb`.
- **TR-07 / `CLAUDE.md` topology:** ingress is Internet → Cloudflare (mTLS) → GCP, with WireGuard back to home lab. The API has to terminate behind a Cloudflare-fronted GCP frontend — anything bypassing that is out.
- **House service shape (`CLAUDE.md`):** Go services use `services/{name}/main.go → app/ → endpoint/{endpointpb} → service/`, with `chi.Mux` for routing, `bedrock` for env config, `pkg/errorpb` for errors, and protobuf-over-HTTP for request/response. The humus framework (`rest.Run`, `rpc.Producer`, `rpc.Handler`) is **not** used — any option proposing it is wrong, not just unconventional.
- **ADR-0006 (resource identifier standard):** URL/resource shapes that need IDs must conform.

So every option below assumes: chi router, protobuf request/response, `pkg/errorpb` errors, deployed behind Cloudflare → GCP, registered like `services/machine/` is.

## Research that should land before you pick

I do not need research to enumerate the options below — they are all house-pattern compliant. But two questions will affect *which* you'd pick, and I'd rather flag them than silently assume:

1. **R-1 — Tenant identity in the request.** Is the tenant identifier carried by Cloudflare-validated mTLS client cert (mapped to a tenant), by a JWT, or by a path/header convention? TR-01 / TR-03 enforcement (per-tenant scoping) leans hard on whichever mechanism we pick. I have not assumed this — I will treat tenant-identity as an open question in the ADR regardless of which option you pick.
2. **R-2 — Operator vs. tenant audiences.** Does this service serve both operators (platform-contract changes per TR-02, tenant updates per TR-04) and tenants (observability queries per TR-03), or only operators? The answer changes whether one API or two are appropriate (see Option D).

If you'd rather punt these to follow-up issues, say so and I'll capture them as Open Questions in the ADR; if you have answers now, share them and I'll fold them into the chosen option.

## Considered options

### Option A — Single versioned protobuf-over-HTTP REST API on chi (`/api/v{n}/...`)

One `tenant-control-plane` service exposes a chi-routed protobuf-over-HTTP API at `/api/v1/...`, mirroring `services/machine/`. Path-based version prefix (`v1`, `v2`, ...) lets multiple contract versions coexist behind the same Cloudflare hostname during the migration window. Per-tenant scoping is enforced inside each handler against the tenant identity asserted by Cloudflare mTLS / a verified token.

- **Satisfies TR-02** (multi-contract): two version prefixes can be served from the same binary or two side-by-side deployments behind the same LB during the bounded migration window.
- **Satisfies TR-04** (no-downtime updates): house pattern (`create_before_destroy` Cloud Run revisions behind `cloud/https-load-balancer/`) already gives rolling deploys; chi mux is process-local, no extra dance.
- **Satisfies TR-07** (topology): plain HTTPS behind Cloudflare → GCP HTTPS LB → Cloud Run, exactly what `cloud/rest-api/` and `cloud/https-load-balancer/` already provision.
- **Satisfies TR-01 / TR-03** *if* tenant-scoping middleware is non-optional on every route — failure mode here is a forgotten check on a new endpoint, mitigated by a chi route group + a single auth middleware that all `/api/v{n}/...` routes mount under.
- **Neutral on TR-06** (data import): import endpoints are just more handlers; idempotency is per-endpoint discipline.
- **Cost:** least new ground. Reuses every existing pattern.
- **Failure mode:** path-prefix versioning makes it tempting to drift contracts in subtle ways (header-only changes, etc.) — TR-02 needs a written contract-change rubric to enforce that breaking changes always bump the prefix.

### Option B — Single gRPC API (with grpc-gateway HTTP/JSON transcoding)

Service exposes a gRPC server (TLS, behind Cloudflare → GCP) with a parallel HTTP/JSON gateway for non-gRPC callers. Versioning via protobuf package names (`tenant.v1`, `tenant.v2`).

- **Partially satisfies TR-07:** Cloudflare supports gRPC, and GCP HTTPS LB → Cloud Run supports gRPC — but it adds a streaming-aware path the rest of the repo's services don't use. We'd be diverging from the chi/protobuf-over-HTTP house pattern; the ADR would have to justify the divergence.
- **Satisfies TR-02:** package-name versioning is clean.
- **Satisfies TR-04:** standard rolling deploys.
- **Neutral on TR-01 / TR-03:** interceptors are the gRPC equivalent of chi middleware; same failure mode as Option A.
- **Cost:** introduces a new framework primitive (gRPC server + grpc-gateway) that no other service in the repo runs. New `cloud/` work for LB-level gRPC config. Justification for departing from `CLAUDE.md` would need to be explicit.
- **Why include it:** if there's a strong streaming-RPC need (e.g. tenant log tailing for TR-03) gRPC's bidi streaming is a real benefit that protobuf-over-HTTP doesn't give you cheaply.

### Option C — Single REST API with header-based versioning (no path prefix)

Same as Option A but contract version is in a header (e.g. `X-Platform-Contract-Version: 2026-04`) rather than the URL path.

- **Partially satisfies TR-02:** can serve multiple versions, but routing them inside one handler set is more error-prone than path-prefix routing — the chi mux can't dispatch on the header without a custom matcher.
- **Same as A on TR-04, TR-07, TR-01, TR-03.**
- **Cost:** same as A plus a custom version-dispatch shim. No real benefit unless we have a strong reason to keep URLs stable across versions, which the TRs don't argue for.
- **Why include it:** completeness — header-versioned APIs are a known idiom; rejecting it explicitly is more useful than pretending it isn't an option.

### Option D — Split into two services: `tenant-control-plane` (operator-facing) + `tenant-self-service` (tenant-facing)

Two services, each chi/protobuf-over-HTTP, deployed independently. Operator endpoints (TR-02 contract rollout, TR-04 tenant updates) live in `tenant-control-plane`; tenant-facing endpoints (TR-03 observability queries, TR-05 export, TR-06 import) live in a sibling service. Each gets its own Cloudflare hostname / mTLS trust profile.

- **Strongly satisfies TR-01 / TR-03:** operator-facing and tenant-facing data planes never share a process, so a bug in tenant-scoping middleware can't accidentally expose operator endpoints to tenants. Strongest isolation story of the four.
- **Satisfies TR-04:** same rolling-deploy pattern, applied to two services.
- **Satisfies TR-07:** two Cloud Run services behind two LB backends, both Cloudflare-fronted.
- **Cost on TR-02:** contract-version coordination across two services is harder — a contract change that touches both surfaces needs synchronized rollouts. Mitigation: keep the contract definitions in a shared `pkg/` and bump versions together.
- **Cost generally:** double the `services/`, `cloud/rest-api/` invocations, deploy automation. Worth it only if the TR-01 isolation argument is decisive.
- **Why include it:** R-2 above — if the audiences really are distinct, this is the option that uses isolation-by-process to enforce TR-01 instead of relying on middleware discipline.

## What I'm asking you to decide

Four options on the table:

- **Option A** is the most house-pattern-faithful — protobuf-over-HTTP on chi with path-prefix versioning. Lowest cost, best pattern fit, all TRs satisfied with middleware discipline as the soft spot for TR-01/TR-03.
- **Option B** trades pattern-fit for gRPC's streaming/typing story. Only worth it if streaming RPCs (e.g. log tail for TR-03) are a near-term need.
- **Option C** is included for completeness and is dominated by A — happy to drop it from the ADR if you agree, but rejecting it on the page is the honest move.
- **Option D** is the strongest TR-01/TR-03 story but costs two services. It's the right call if R-2's answer is "operator and tenant audiences are genuinely distinct".

Which do you want to accept, do you want me to revise the option set, and how would you like to handle R-1 (tenant identity) and R-2 (audience split) — answer now, or capture as Open Questions on whichever option you pick?

I will not set `status: accepted` until you choose.
