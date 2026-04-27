# tenant-control-plane API exposure — options & recommendation

Stage 2 (capability-scoped ADR territory) for the self-hosted-application-platform. The tenant-control-plane is the service operators and tenants call to manage tenant lifecycle (create, update, evict, export, import, query observability scope). Below are the realistic options for *how that API is exposed*, evaluated against the reviewed tech requirements and this repo's house patterns (CLAUDE.md).

## Hard constraints from the reviewed requirements

- **TR-07** — all inter-service comms traverse Cloudflare → GCP (with WireGuard back to home lab). Whatever protocol we pick must terminate at Cloudflare and survive a Cloudflare-fronted hop. mTLS at the edge is already the house pattern (`cloud/mtls/cloudflare-gcp/`).
- **TR-01** — strict tenant isolation. The exposure layer must carry a verifiable tenant identity into every request; the API surface cannot offer ambient cross-tenant operations.
- **TR-02** — multiple contract versions concurrently for a bounded window. The wire protocol and routing layer must support versioned APIs side-by-side.
- **TR-04** — no tenant-perceived downtime during operator-initiated updates. Whatever fronts the API must support drain/rolling deploys (already true for our HTTPS LB pattern in `cloud/https-load-balancer/`).

## House patterns from CLAUDE.md

- Existing Go services (see `services/machine/`) expose **HTTP via chi, request/response in protobuf** (`application/x-protobuf`), errors as `application/problem+protobuf` via `pkg/errorpb`. No humus, no gRPC, no embedded YAML.
- TLS is terminated at the service with a self-signed cert behind the LB; mTLS to Cloudflare lives in front.
- Config from env via `bedrock/config`.

The "house API style" is therefore: **chi + protobuf-over-HTTP + problem+protobuf**, fronted by the existing Cloudflare → GCP HTTPS LB with mTLS.

## Options

### Option A — chi + protobuf-over-HTTP (the `services/machine` pattern)
Tenant-control-plane is a Cloud Run service exposing `/api/v1/...` routes on a chi mux; requests/responses are protobuf messages, errors are `application/problem+protobuf`. Fronted by `cloud/rest-api/` + `cloud/https-load-balancer/`, mTLS via `cloud/mtls/cloudflare-gcp/`. Versioning via URL prefix (`/api/v1`, `/api/v2`).

- **Pros:** matches the only existing Go service in the repo verbatim; reuses `pkg/errorpb`; reuses every Terraform module already in `cloud/`; protobuf gives us schema evolution for TR-02; trivially supports operator and tenant clients (anything that can speak HTTP+proto); plays nicely with the Cloudflare-fronted topology (TR-07).
- **Cons:** no built-in streaming; URL-prefix versioning is coarse (acceptable for a control plane where call volume is low).

### Option B — gRPC (with grpc-gateway for any browser callers)
Tenant-control-plane exposes a gRPC service; an optional REST gateway translates JSON for browser/CLI clients.

- **Pros:** strong typing, codegen, streaming, mature deadline/cancellation semantics.
- **Cons:** **breaks the house pattern** — no existing service uses gRPC. Cloudflare's gRPC support has rougher edges than plain HTTPS (gRPC-Web or H2 quirks at the edge). Adds a second protocol surface (gateway) which doubles the versioning story for TR-02. Not justified for a low-QPS control plane.

### Option C — GraphQL
Single endpoint, client-shaped queries.

- **Cons:** no precedent in the repo; tenant-isolation enforcement (TR-01) gets harder when the query shape is dynamic (every resolver becomes an authz hot spot); versioning model (TR-02) is "evolve the schema forever," which fights an explicit bounded migration window. Reject.

### Option D — Plain JSON REST (no protobuf)
chi + `encoding/json`, OpenAPI for schema.

- **Pros:** lowest friction for ad-hoc curl/CLI use.
- **Cons:** diverges from the house pattern and from `pkg/errorpb` (which is *protobuf* problem details). We'd lose the existing error-serialization machinery and have to maintain two error formats across the platform.

### Option E — Cloud Run direct invocation via GCP IAM (no public HTTP API)
Operators and tenants call Cloud Run with GCP-issued ID tokens; no Cloudflare fronting for control-plane calls.

- **Cons:** **violates TR-07** — bypasses the Cloudflare → GCP path. Also forces tenants to hold GCP identities, which is a non-starter for the isolation model. Reject.

## Recommendation: Option A — chi + protobuf-over-HTTP, fronted by the existing Cloudflare → GCP HTTPS LB

This is the only option that simultaneously:
1. Matches the established house pattern (`services/machine/` is the reference implementation, and CLAUDE.md is explicit that this is *the* current pattern).
2. Reuses `pkg/errorpb` for `*Problem` / `*ValidationProblem` / `*ConflictProblem` responses, which the rest of the platform's clients will already understand.
3. Satisfies TR-07 cleanly via `cloud/rest-api/` + `cloud/https-load-balancer/` + `cloud/mtls/cloudflare-gcp/` with no new infrastructure module.
4. Supports TR-02 via URL-prefix versioning (`/api/v1`, `/api/v2` mounted as parallel chi sub-routers in the same binary during the migration window) — protobuf field evolution covers minor changes within a version.
5. Carries tenant identity for TR-01 via the mTLS client cert (operator-issued, per-tenant) extracted in middleware and propagated as a request-scoped value into every handler. No ambient/global tenant queries.

### Concrete shape (sketch, for the eventual ADR)

```
services/tenant-control-plane/
├── main.go                 # os.Exit(app.Main(context.Background()))
├── app/app.go              # ConfigFromEnv, Main, chi mux, TLS listener (mirrors services/machine/app)
├── endpoint/
│   ├── v1/                 # RegisterTenantsV1(mux, deps...) → /api/v1/tenants/*
│   ├── v2/                 # RegisterTenantsV2(mux, deps...) → /api/v2/tenants/* during TR-02 window
│   └── endpointpb/         # generated request/response messages
├── service/                # Firestore (or chosen store) client for tenant records
└── middleware/
    └── tenantidentity.go   # extract tenant ID from mTLS cert, inject into ctx
```

- Errors via `pkg/errorpb` → `errorHandler(ctx, w, err)` → `application/problem+protobuf`.
- Config via `bedrock/config` env-only (`HTTP_PORT`, `GCP_PROJECT_ID`, etc.).
- Concurrency via `sourcegraph/conc/pool` and `wg.Go` (Go 1.26).
- One small interface per backend client, defined in `endpoint/`.

### What this defers
- The exact tenant-identity proof (mTLS client cert subject vs. signed JWT relayed by Cloudflare) is a separate ADR — both compose with this exposure choice.
- Streaming export for TR-05 (potentially large payloads): handle as chunked HTTP responses or signed object-store URLs rather than restructuring the API; does not change the choice here.

## Files
- Tech requirements (refreshed `reviewed_at` to today): `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`
- Reference implementation to mirror: `/home/carson/github.com/Zaba505/infra/services/machine/app/app.go`
- House conventions: `/home/carson/github.com/Zaba505/infra/CLAUDE.md`
