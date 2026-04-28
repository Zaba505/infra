# Plan Implementation — self-hosted-personal-media-storage

## Gate check

- `tech-design.md` `reviewed_at: 2026-04-27`, file unchanged since — review gate satisfied.
- `gh issue list` for `story(gap):` on this capability — none open.
- `gh issue list` for `story(component):` on this capability — none open.
- All 5 components in inventory have design docs on disk: `tenant-registry.md`, `photo-store.md`, `share-service.md`, `firestore-tenants.md`, `media-bucket.md`.

Slicing.

## Component summaries (anchors for slicing)

- **tenant-registry** (Go service) — Firestore-backed `tenants` collection; 5 HTTP/protobuf endpoints (POST, GET-by-id, PATCH, DELETE, LIST). Cloud Run, SA `tenant-registry-sa`. Depends on `cloud/firestore-tenants`.
- **photo-store** (Go service) — Mediates GCS I/O with per-tenant prefix isolation; 4 endpoints (upload, fetch, delete, list). Cloud Run, SA `photo-store-sa`. Depends on `cloud/media-bucket` and on tenant-registry reachability for tenant-id validation.
- **share-service** (Go service) — Owns `share-grants` and append-only `share-audit` collections; 5 endpoints (mint, fetch, revoke, signed-URL, list). Cloud Run, SA `share-service-sa`. Depends on `cloud/firestore-tenants` and `cloud/media-bucket` IAM.
- **firestore-tenants** (Terraform module) — `google_firestore_database` (Native), indexes per query path, IAM bindings for tenant-registry/share-service SAs (append-only on `share-audit`).
- **media-bucket** (Terraform module) — Multi-region `google_storage_bucket` with versioning + UBLA, IAM for photo-store and share-service SAs, 90-day version-deletion lifecycle.

## Proposed implementation tasks (dependency-ordered)

1. **Author `cloud/firestore-tenants` Terraform module**
   - Components: firestore-tenants
   - Summary: Provision Firestore Native database, per-collection indexes, and SA IAM bindings (tenant-registry RW on `tenants`; share-service RW on `share-grants`, append-only on `share-audit`).
   - Acceptance: `terraform fmt -recursive -check` passes; module exposes `database_name` + `collections` outputs; `create_before_destroy` set; provider pinned to GCP v7.11.0.
   - Source: `components/firestore-tenants.md`; ADR-0001, ADR-0003; TR-01, TR-04, TR-05, TR-06, TR-07.
   - Prerequisites: none.

2. **Author `cloud/media-bucket` Terraform module**
   - Components: media-bucket
   - Summary: Provision multi-region GCS bucket with versioning + UBLA, IAM bindings for `photo-store-sa` (object RW) and `share-service-sa` (signed-URL minting), and 90-day version-deletion lifecycle.
   - Acceptance: `terraform fmt -recursive -check` passes; `bucket_name` output; lifecycle rule asserted via plan.
   - Source: `components/media-bucket.md`; ADR-0002; TR-01, TR-02.
   - Prerequisites: none.

3. **Add `tenant-registry` Firestore schema + `service.tenants` client**
   - Components: tenant-registry
   - Summary: Define `TenantRecord` (tenant_id, display_name, state enum, created_at, updated_at) and a Firestore client in `services/tenant-registry/service/` covering create/get/update/list/soft-delete.
   - Acceptance: Unit tests cover client CRUD against a Firestore emulator or fake; record marshals/unmarshals via protobuf.
   - Source: `components/tenant-registry.md` (Data section); ADR-0001; TR-04, TR-05.
   - Prerequisites: none (module #1 only required for deployment, not unit dev).

4. **Wire `tenant-registry` HTTP endpoints**
   - Components: tenant-registry
   - Summary: Implement `services/tenant-registry/endpoint/` handlers for `POST/GET/PATCH/DELETE /api/v1/tenants[/{id}]` and `GET /api/v1/tenants` (paged), using chi mux + `pkg/errorpb`.
   - Acceptance: All 5 endpoints respond with documented status codes; validation errors → `application/problem+protobuf`; integration test exercises the round-trip.
   - Source: `components/tenant-registry.md` (API surface); ADR-0001; TR-04, TR-05.
   - Prerequisites: task 3.

5. **Deploy `tenant-registry` to Cloud Run (SA + module wiring)**
   - Components: tenant-registry, firestore-tenants
   - Summary: Provision `tenant-registry-sa` service account, root-module composition that consumes `cloud/firestore-tenants` and deploys the Cloud Run service behind the platform HTTPS LB with health check on `HTTP_PORT`.
   - Acceptance: Cloud Run service definition lints clean; SA has read/write on `tenants` only; `terraform fmt -recursive -check` passes.
   - Source: `components/tenant-registry.md` (Operational concerns); ADR-0001; TR-01.
   - Prerequisites: tasks 1, 4.

6. **Add `photo-store` GCS client + object naming**
   - Components: photo-store
   - Summary: Implement `services/photo-store/service/` GCS client with `{tenant_id}/{ulid}.{ext}` object naming and per-tenant prefix isolation enforcement.
   - Acceptance: Unit tests cover put/get/delete/list keyed by tenant prefix; cross-tenant access attempts fail closed.
   - Source: `components/photo-store.md`; ADR-0002; TR-01, TR-02.
   - Prerequisites: none.

7. **Wire `photo-store` HTTP endpoints with tenant-registry validation**
   - Components: photo-store, tenant-registry
   - Summary: Implement upload (multipart), fetch, delete, list endpoints; validate `tenant_id` against the tenant-registry service before any object I/O; honor share-grant bearer for fetch.
   - Acceptance: All 4 endpoints respond per design; tenant-id validation rejects unknown/evicted tenants; integration test covers round-trip with a stub tenant-registry.
   - Source: `components/photo-store.md` (API surface); ADR-0002; TR-01, TR-02.
   - Prerequisites: task 6, task 4.

8. **Deploy `photo-store` to Cloud Run (SA + module wiring)**
   - Components: photo-store, media-bucket
   - Summary: Provision `photo-store-sa`, root-module composition consuming `cloud/media-bucket`, Cloud Run deployment with health check on `HTTP_PORT`, env wiring for tenant-registry endpoint.
   - Acceptance: SA scoped to bucket per-prefix policy only; `terraform fmt -recursive -check` passes.
   - Source: `components/photo-store.md` (Operational concerns); ADR-0002; TR-01, TR-02.
   - Prerequisites: tasks 2, 7.

9. **Add `share-service` Firestore schema + clients (share-grants + share-audit)**
   - Components: share-service
   - Summary: Define grant + audit-event record types and `services/share-service/service/` clients; audit client is append-only.
   - Acceptance: Unit tests cover grant CRUD and audit append; audit client cannot update or delete; protobuf round-trip green.
   - Source: `components/share-service.md` (Data); ADR-0003; TR-03, TR-06, TR-07.
   - Prerequisites: none.

10. **Wire `share-service` HTTP endpoints (incl. signed-URL minting)**
    - Components: share-service, media-bucket
    - Summary: Implement mint, fetch, revoke, signed-URL, list endpoints. Revoke and signed-URL paths write `revoked` / `accessed` events to the audit collection. Signed URLs are short-lived against the media bucket.
    - Acceptance: All 5 endpoints respond per design; audit events written on revoke and access; signed URLs validate against bucket IAM; integration test covers grant → access → revoke flow.
    - Source: `components/share-service.md` (API surface); ADR-0003; TR-03, TR-06, TR-07.
    - Prerequisites: task 9.

11. **Deploy `share-service` to Cloud Run (SA + module wiring)**
    - Components: share-service, firestore-tenants, media-bucket
    - Summary: Provision `share-service-sa`, root-module composition consuming both Terraform modules, Cloud Run deployment with health check, env wiring for Firestore + bucket.
    - Acceptance: SA scoped to RW on `share-grants`, append-only on `share-audit`, signed-URL minting on bucket; `terraform fmt -recursive -check` passes.
    - Source: `components/share-service.md` (Operational concerns); ADR-0003; TR-03.
    - Prerequisites: tasks 1, 2, 10.

## Mirror & ask

I've sliced 11 implementation tasks across 5 components. Before I file issues: do the slices look right? Anything to merge, split, reorder, or drop? Once you say go, I'll file 11 issues.

---

## User: "go" — approval received. Filing issues.

11 issues logged to `outputs/intended_issues.jsonl` (one JSON object per line). After real `gh issue create`, the manifest of issue numbers/URLs would be printed back in dependency order. Each issue's body links the parent capability and points at Step 12 for per-task development planning.
