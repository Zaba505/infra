# ADR Draft: tenant-control-plane API Exposure

**Capability:** self-hosted-application-platform
**Service:** tenant-control-plane
**Decision:** How does `tenant-control-plane` expose its API?

Below are candidate options. Each is tied back to the relevant TR-NNs from `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`. Pick one and I will write up the full MADR.

---

## Relevant Technical Requirements

- **TR-01** — strict tenant isolation (data + compute)
- **TR-02** — concurrent platform-contract versions during migration
- **TR-03** — per-tenant scoping of observability queries
- **TR-04** — no tenant-perceived downtime on operator-initiated updates
- **TR-07** — all inter-service traffic must traverse the Cloudflare → GCP path

The house pattern (per `CLAUDE.md` and `services/machine/`) is: Go service, `chi` router, protobuf-over-HTTP (`application/x-protobuf`), `application/problem+protobuf` errors via `pkg/errorpb`, fronted by a GCP HTTPS load balancer (`cloud/https-load-balancer`) with Cloudflare mTLS trust (`cloud/mtls/cloudflare-gcp`). Any option below assumes that fronting topology to satisfy **TR-07**.

---

## Option 1 — Protobuf-over-HTTP on chi (house pattern, single versioned path prefix)

A single Cloud Run service exposes endpoints like `/api/v1/...` using `chi` + protobuf request/response, identical to `services/machine/`. New contract versions are introduced as `/api/v2/...` mounted on the same mux; old paths remain registered until the migration window closes.

- **Pros**
  - Matches the existing house pattern exactly — no new framework, no new content-type, reuses `pkg/errorpb`.
  - Concurrent contract versions (**TR-02**) are just additional `Register*` functions on the mux.
  - Zero-downtime rollouts (**TR-04**) come for free via Cloud Run + `create_before_destroy`.
  - Tenant identity can be enforced in chi middleware before any handler runs (**TR-01**, **TR-03**).
- **Cons**
  - No formal IDL contract beyond hand-curated `.proto` request/response messages; no service-level RPC schema.
  - Streaming (e.g. log tails for **TR-03**) is awkward — must be modeled as chunked HTTP or SSE on top of protobuf.
- **TR coverage:** TR-01, TR-02, TR-03, TR-04, TR-07.

## Option 2 — gRPC service (with grpc-gateway for any browser/HTTP clients)

Define the API in `.proto` as a gRPC service. Cloud Run terminates HTTP/2; tenants hit gRPC directly, or grpc-gateway exposes a JSON/HTTP shim where needed.

- **Pros**
  - Strong IDL with generated clients; native streaming for log/metric/trace tails (**TR-03**).
  - Concurrent contract versions (**TR-02**) handled cleanly via `package v1`/`package v2` services on the same server.
- **Cons**
  - **Diverges from the house pattern** — no other service in this repo runs gRPC; introduces new tooling, new error model (gRPC `Status` vs. `pkg/errorpb`), and a second content-type story.
  - Cloudflare mTLS path (**TR-07**) needs validation for HTTP/2 + gRPC framing end-to-end.
  - Higher ongoing cost: every future service author has to choose between two patterns.
- **TR coverage:** TR-01, TR-02, TR-03, TR-04, TR-07 (with extra integration work for TR-07).

## Option 3 — ConnectRPC over the existing chi mux

Use Connect (`connectrpc.com/connect`) which serves protobuf RPCs as plain HTTP/1.1 + HTTP/2 endpoints and mounts onto a `chi` mux. Same `.proto` IDL benefits as gRPC, but compatible with the existing Cloudflare→GCP HTTP path and the chi-based service shape.

- **Pros**
  - Keeps the chi + protobuf substrate of the house pattern; mountable alongside any existing `Register*` handlers.
  - Generated clients and a real service IDL (better than Option 1 for **TR-02** version management).
  - Native streaming on HTTP/2 supports per-tenant observability streams (**TR-03**).
- **Cons**
  - New dependency and new error model — `pkg/errorpb` would either need a Connect adapter or be replaced for these endpoints, creating two error conventions.
  - Still a deviation from the established `services/machine/` shape, just a smaller one than full gRPC.
- **TR coverage:** TR-01, TR-02, TR-03, TR-04, TR-07.

## Option 4 — Split control vs. data planes: chi+protobuf for control, separate streaming endpoint for observability

Keep the house pattern (Option 1) for tenant control operations (provision, update, evict, export). Expose **TR-03** observability streams via a dedicated endpoint (e.g. SSE or chunked protobuf) on the same service, or punt observability to a sibling service entirely.

- **Pros**
  - Lowest deviation from the house pattern for the bulk of the API surface.
  - Lets observability (**TR-03**) evolve independently without forcing the whole API onto a streaming RPC framework.
  - Each plane can be versioned independently for **TR-02**.
- **Cons**
  - Two API styles to document and authenticate; tenant-isolation middleware (**TR-01**) must be duplicated or shared as a library.
  - Risk of the observability surface eventually warranting its own service anyway, making this an interim shape.
- **TR coverage:** TR-01, TR-02, TR-03, TR-04, TR-07.

---

## Open questions before I draft the MADR

1. Is per-tenant streaming observability (**TR-03**) served by `tenant-control-plane` itself, or a sibling service? This is the main thing pushing toward Options 2/3/4.
2. Are there non-Go consumers of this API today (operator CLI, future tenant SDKs)? If yes, IDL-generated clients (Options 2/3) get more valuable.
3. How long is the **TR-02** bounded migration window expected to be? Long windows tilt toward stronger versioned IDL.

Tell me which option (1–4) and I will draft the ADR in MADR 4.0.0 format under `docs/content/capabilities/self-hosted-application-platform/adrs/`.
