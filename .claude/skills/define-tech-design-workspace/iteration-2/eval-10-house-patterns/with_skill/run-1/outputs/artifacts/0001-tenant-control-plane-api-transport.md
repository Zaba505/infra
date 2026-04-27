---
title: "[0001] Tenant Control Plane API Transport"
description: >
    Choose how the tenant-control-plane service exposes its API to operator tooling and to platform-internal callers.
type: docs
weight: 1
category: "api-design"
status: "proposed"
date: 2026-04-26
deciders: []
consulted: []
informed: []
---

**Parent capability:** [self-hosted-application-platform](../_index.md)
**Addresses requirements:** [TR-02](../tech-requirements.md#tr-02-operators-must-be-able-to-roll-out-a-platform-contract-change-without-breaking-existing-tenants), [TR-04](../tech-requirements.md#tr-04-operator-initiated-tenant-updates-must-complete-without-tenant-perceived-downtime-for-online-workloads), [TR-07](../tech-requirements.md#tr-07-all-inter-service-communication-must-traverse-the-cloudflare--gcp-path)

## Context and Problem Statement

The `tenant-control-plane` service is the platform-internal API surface through which the operator (and other platform services) drives tenant lifecycle: register a tenant, declare its resource needs, initiate an operator-initiated update (TR-04), publish or roll a platform-contract version (TR-02), trigger an export, evict a tenant. It is *not* a tenant-facing API — end users never touch it; tenants themselves do not call it.

The decision is which on-the-wire transport and encoding the service exposes. The transport must:

- Sit cleanly behind the existing Cloudflare → GCP path so that all calls into the service traverse that path (TR-07). No alternative network plane.
- Carry strongly-typed contract versions, since rolling a platform-contract version (TR-02) requires the contract to be a *thing* the API can name and validate, not free-form JSON drift.
- Support long-running, observable operations (operator-initiated updates per TR-04 are not always sub-second), without requiring a bespoke streaming infrastructure.
- Match how Go services in this repo are already built (chi + protobuf-over-HTTP + `pkg/errorpb`), so the operator's *Maintenance budget* KPI is not eroded by maintaining a second service shape.

## Decision Drivers

* TR-07 — must work over Cloudflare → GCP (HTTPS only path; no arbitrary L4).
* TR-02 — contract evolution requires versioned, machine-checked schemas.
* TR-04 — operator-initiated updates need request/response shapes that can carry operation handles for follow-up polling, without forcing bidirectional streams.
* CLAUDE.md house pattern — `services/machine/` already uses `chi` + `application/x-protobuf` + `pkg/errorpb`. Diverging here means a second pattern to learn, test, and operate.
* Operator maintenance budget KPI (capability §Success Criteria) — every novel framework increases weekly maintenance cost.

## Considered Options

### Option A — HTTP + protobuf over chi (the existing house pattern)

`POST /api/v1/...` endpoints registered on a `*chi.Mux`, request and response bodies are `application/x-protobuf`, errors are `application/problem+protobuf` via `pkg/errorpb`. Schemas live in `services/tenant-control-plane/endpoint/endpointpb/`.

* Good — identical to `services/machine/`; zero new operational concepts for the operator.
* Good — protobuf gives versioned, machine-checked schemas (TR-02): contract version is just the proto package name + reserved-field discipline.
* Good — clean over Cloudflare's HTTPS proxy (TR-07); no L4 surprises.
* Good — long-running operations (TR-04) are modeled as "submit returns operation handle, poll handle" — plain request/response, no streaming infra required.
* Good — `pkg/errorpb` (`NewValidationError`, `NewConflictError`, `NewInternalError`) already exists and dispatches via `errorHandler`.
* Bad — protobuf-over-HTTP is not a public standard (no off-the-shelf curl-friendly debugging); operator must use a small internal CLI or a generated client. Acceptable because the API is operator-internal.
* Bad — no built-in streaming if a future requirement needs server-push. Mitigation: operation-handle polling covers the foreseeable case (TR-04).

### Option B — gRPC (HTTP/2, native)

Full gRPC server, codegen'd stubs, native streaming.

* Good — first-class versioned schemas (TR-02), same as Option A.
* Good — bidirectional streaming if ever needed.
* Bad — gRPC over Cloudflare's standard proxy is awkward (HTTP/2 + trailers): Cloudflare gRPC support exists but is a non-default code path, fighting TR-07's "uses the existing path" intent.
* Bad — diverges from the house pattern. The repo currently has zero gRPC services; introducing one adds a second server runtime, observability shape, and error model to maintain — directly contrary to the operator-maintenance-budget KPI.
* Bad — `pkg/errorpb` would need a parallel gRPC-status mapping.

### Option C — JSON over HTTP, OpenAPI-described

`POST /api/v1/...` with `application/json`, schema described in OpenAPI; errors as `application/problem+json` (RFC 7807).

* Good — most curl-friendly; lowest operator friction for ad-hoc debugging.
* Good — over Cloudflare HTTPS without surprises (TR-07).
* Mixed on TR-02 — JSON+OpenAPI can be versioned, but enforcement is weaker than protobuf: drift is easy and only caught at runtime by validators we'd have to wire ourselves.
* Bad — diverges from the house pattern. `pkg/errorpb` is protobuf-shaped; a JSON service either re-implements problem responses or maintains two error machineries.
* Bad — every other Go service in the repo would still be protobuf, so the operator pays the cost of two encodings.

### Option D — Connect (connectrpc)

Connect protocol — single server speaks gRPC, gRPC-Web, and Connect's own HTTP+JSON/protobuf, from one schema.

* Good — protobuf schemas (TR-02), HTTPS-friendly Connect protocol (TR-07).
* Good — curl-friendly via the JSON variant when debugging.
* Bad — net-new framework dependency in a repo that has zero Connect services today; CLAUDE.md explicitly warns against importing frameworks the repo isn't already using ("humus framework: not used in this repo").
* Bad — strictly more surface than Option A for a service that does not need three protocols at once.

### Option E — GraphQL

Single endpoint, query language, schema-first.

* Bad — does not match any in-repo pattern; introduces a resolver runtime and an N+1 vigilance discipline for what is fundamentally a small, command-shaped API (register tenant, update tenant, roll contract).
* Bad — GraphQL's strength (flexible client-driven queries) is wasted: there is one operator client, not a population of clients with divergent shape needs.
* Neutral on TR-07 (HTTPS) but loses on TR-02 enforcement (GraphQL schemas don't carry the same wire-level versioning discipline as protobuf).

## Decision Outcome

Chosen option: **Option A — HTTP + protobuf over chi**, because it is the only option that satisfies TR-02 (versioned, machine-checked contract), TR-04 (operation-handle polling over plain request/response), and TR-07 (Cloudflare → GCP HTTPS path) *without* introducing a second service shape that would erode the operator-maintenance-budget KPI. The CLAUDE.md house pattern already encodes this choice for `services/machine/`; following it here is the cheapest way to satisfy the requirements.

### Consequences

* Good — operator runs one service shape across the platform; chi mux, protobuf endpoints, `pkg/errorpb` error dispatch are reused verbatim.
* Good — the `tenant-control-plane` proto package is the canonical artifact carrying the platform contract version (TR-02) — version bumps are visible in code review.
* Bad — operator-side ad-hoc debugging requires either an internal CLI or `protoc --decode_raw`; we accept this because the API is operator-internal, not tenant-facing.
* Bad — a future requirement for server-push notifications would force a follow-up ADR (websocket or SSE) rather than coming "for free." Acceptable; not implied by current TRs.
* Requires:
  * New `services/tenant-control-plane/` skeleton mirroring `services/machine/` (`main.go`, `app/app.go`, `endpoint/`, `endpoint/endpointpb/`, `service/`).
  * `endpoint/endpointpb/` proto definitions for the initial command set (register tenant, update tenant, publish contract version, request export, evict).
  * Deployment via `cloud/rest-api/` Cloud Run module (already used for HTTP services), fronted by `cloud/https-load-balancer/` with mTLS from Cloudflare per the existing topology.

### Realization

* Service: `services/tenant-control-plane/` — chi mux on `HTTP_PORT` (default 8080), self-signed TLS, `app.Main(ctx) int` entrypoint per CLAUDE.md.
* Endpoints: `services/tenant-control-plane/endpoint/*.go`, each with a `RegisterX(mux *chi.Mux, deps...)` function and an `endpointpb/` package for proto types. Routes live under `/api/v1/...`.
* Errors: `pkg/errorpb` — `NewValidationError` for malformed requests, `NewConflictError` for contract-version conflicts during a TR-02 rollout, `NewInternalError` otherwise; `errorHandler(ctx, w, err)` dispatches.
* Backend: `services/tenant-control-plane/service/` for Firestore (or other) clients holding tenant registry state.
* Network path: Cloudflare → GCP HTTPS Load Balancer → Cloud Run, satisfying TR-07.

## Open Questions

* Whether the platform-contract version (TR-02) is encoded as a proto package suffix (`.v1`, `.v2`) or as a header / first-field discriminator. Defer to a follow-up ADR once the second contract version is on the horizon — the chosen transport supports either.
* Long-running operation representation (operation handle shape, polling endpoint vs. callback) for TR-04 — captured here as "operation handle + poll" but the precise shape is a follow-up ADR.
