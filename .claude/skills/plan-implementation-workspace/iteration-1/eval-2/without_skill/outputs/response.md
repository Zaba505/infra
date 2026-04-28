# Plan Implementation — self-hosted-personal-media-storage

## Gate state

- `tech-design.md` shows `reviewed_at: 2026-04-27`, current per the user's confirmation.
- No open `story(gap):` issues for this capability.
- All five components in the inventory have design docs on disk under `components/`:
  - `tenant-registry.md`
  - `photo-store.md`
  - `share-service.md`
  - `firestore-tenants.md`
  - `media-bucket.md`
- No open `story(component):` issues for this capability.

Slicing proceeded.

## Slice rationale

Two infra modules first (services depend on their IAM/collections). Then `tenant-registry` is built end-to-end because `photo-store` validates tenant IDs against it. `photo-store` follows once the bucket module and tenant-registry API exist. `share-service` lands last because it depends on both Firestore (collections + indexes) and the media bucket (signed-URL IAM). A capability docs task closes out the rollout.

Each slice is one PR-shaped seam from the component designs — schemas/clients separated from endpoints, deployment wiring separated from app code, modules separated from services — so each lands with its own review and test signal.

## Mirrored task list (ordered by hard dependency)

1. **Author `cloud/firestore-tenants` Terraform module** — firestore-tenants. Database + indexes + IAM. No prereqs.
2. **Author `cloud/media-bucket` Terraform module** — media-bucket. Multi-region GCS + IAM + lifecycle. No prereqs.
3. **tenant-registry — protobuf types and Firestore service client** — tenant-registry. No prereqs.
4. **tenant-registry — app skeleton (config, main, chi mux, TLS)** — tenant-registry. After (3).
5. **tenant-registry — implement CRUD HTTP endpoints** — tenant-registry. After (3) + (4).
6. **tenant-registry — deployment wiring (SA + Cloud Run + LB backend)** — tenant-registry. After (1) + (5).
7. **photo-store — protobuf types, GCS client, tenant-registry client** — photo-store / tenant-registry. After (5).
8. **photo-store — endpoints, app skeleton, deployment wiring** — photo-store. After (2) + (7).
9. **share-service — protobuf types and Firestore clients (grants + append-only audit)** — share-service. After (1).
10. **share-service — endpoints, signed-URL minting, app skeleton, deployment wiring** — share-service. After (2) + (9).
11. **Capability docs — operator runbook page and Hugo cross-link verification** — cross-cutting. Soft-depends on all earlier tasks.

User said "go" — treated as approval. Issues filed.

## Filed issues

Logged to `outputs/intended_issues.jsonl` (11 lines), one JSON object per intended issue, in dependency order. Per workspace policy, `gh` was not invoked.

## Notes / surfaced risk

- Task 8 (photo-store endpoints) implements per-tenant authz; the share-grant validation hook is stubbed there because share-service ships later. Step 12 for that task should plan the integration test that exercises the live hook once share-service is deployed.
- The tech design lists three Go services and two Terraform modules but does not name a single capability-level integration test. If desired, a 12th task could be added to wire an end-to-end smoke test once all services are up — surfacing here rather than inventing it.
