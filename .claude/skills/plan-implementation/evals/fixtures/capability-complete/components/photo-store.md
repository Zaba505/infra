---
title: "Component: photo-store"
type: docs
---

**Parent tech design:** [tech-design.md](../tech-design.md)
**Type:** Go service (`services/photo-store/`)
**Established by:** ADR-0002

## Responsibility
Mediates media-object I/O against the GCS bucket with per-tenant isolation. Tenants cannot read/write outside their own prefix.

## API surface
- `POST /api/v1/media` — upload (multipart). Accepts `tenant_id` (validated against tenant-registry), returns object handle.
- `GET /api/v1/media/{object_id}` — fetch (caller must be the owning tenant or bear a valid share grant).
- `DELETE /api/v1/media/{object_id}` — delete (owning tenant only).
- `GET /api/v1/media?tenant_id=...` — list owned media, paged.

## Dependencies
- `cloud/media-bucket` module must exist before deployment.
- `tenant-registry` must be reachable for tenant-id validation.

## Operational concerns
- Cloud Run service; service account `photo-store-sa` with read/write only to its tenant prefix policy on the bucket.
- Object naming: `{tenant_id}/{ulid}.{ext}`.
