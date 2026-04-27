# Stage 2 — Tenant Control Plane API Transport

We're in **Stage 2** for `self-hosted-application-platform`. `tech-requirements.md` is reviewed (`reviewed_at: 2026-04-26`), so we're clear to draft ADRs. This is the first capability-scoped ADR, so I also created `docs/content/capabilities/self-hosted-application-platform/adrs/_index.md`.

The decision in scope: **how does the `tenant-control-plane` service expose its API?** This is purely the platform-internal control surface (operator + platform services in, no tenant end users). I deliberately did *not* bundle in adjacent decisions like "where does tenant-registry state live", "how is the platform contract version encoded in the schema", or "how are long-running operations modeled" — those are separate ADRs and will get their own treatment.

## TRs in scope

- **TR-02** — operators must roll a platform-contract change without breaking existing tenants → forces *machine-checked, versioned schemas*.
- **TR-04** — operator-initiated updates without tenant-perceived downtime → forces *long-running-operation-friendly* request/response shapes.
- **TR-07** — all inter-service traffic traverses Cloudflare → GCP → forces *HTTPS-clean* transports; rules out anything that fights the proxy.

Plus the inherited CLAUDE.md house pattern (chi + protobuf-over-HTTP + `pkg/errorpb`), which the operator-maintenance-budget KPI tells us not to diverge from without cause.

## Options considered

| # | Option | Verdict |
|---|---|---|
| A | **HTTP + protobuf over chi** (the house pattern: `application/x-protobuf`, errors via `pkg/errorpb`) | **Chosen** |
| B | gRPC native (HTTP/2 + trailers) | Rejected — fights TR-07 (Cloudflare proxy + gRPC is non-default), and adds a second service shape against the maintenance-budget KPI |
| C | JSON over HTTP, OpenAPI-described, RFC 7807 errors | Rejected — weaker TR-02 enforcement than protobuf, and forces a second error machinery alongside `pkg/errorpb` |
| D | Connect (connectrpc) | Rejected — net-new framework dependency for a service that does not actually need three protocols at once |
| E | GraphQL | Rejected — wrong shape (the API is command-shaped, single-client), and doesn't match any in-repo pattern |

## Pick: Option A — HTTP + protobuf over chi

It's the only option that satisfies TR-02, TR-04, and TR-07 *without* introducing a second service shape. The repo already encodes this choice in `services/machine/`; following it here keeps the operator on one runtime, one error model, and one deploy pipeline (`cloud/rest-api/` → `cloud/https-load-balancer/`).

Service skeleton mirrors `services/machine/`:

```
services/tenant-control-plane/
├── main.go                       # os.Exit(app.Main(context.Background()))
├── app/app.go                    # ConfigFromEnv, Main(ctx) int — chi + TLS server
├── endpoint/
│   ├── *.go                      # RegisterX(mux *chi.Mux, deps...) per command
│   └── endpointpb/               # proto-generated request/response types
└── service/                      # Firestore (or other) tenant-registry client
```

Errors via `pkg/errorpb` (`NewValidationError`, `NewConflictError`, `NewInternalError`) dispatched through `errorHandler(ctx, w, err)`. Network: Cloudflare → GCP HTTPS LB → Cloud Run, per TR-07.

Status on the ADR is **`proposed`**. Per the skill, I'm not advancing to Stage 3 until you flip it to `accepted`. Two open questions are captured in the ADR rather than papered over:

1. How the platform-contract version (TR-02) is carried on the wire — proto package suffix vs. header. Punt to a follow-up ADR.
2. Long-running-operation shape (TR-04) — "operation handle + poll" sketched, exact shape deferred.

## Files written

- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (copied from reviewed fixture)
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/adrs/_index.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/adrs/0001-tenant-control-plane-api-transport.md`

## Next

Tell me to flip ADR-0001 to `accepted`, or push back on the option choice / rationale. Other Stage 2 decisions still pending for this capability include tenant-state storage (TR-01, TR-03), tenant-data export mechanism (TR-05), data-import idempotency (TR-06), zero-downtime update mechanism (TR-04), and contract-version encoding (TR-02). Want me to enumerate the full Stage 2 decision list next?
