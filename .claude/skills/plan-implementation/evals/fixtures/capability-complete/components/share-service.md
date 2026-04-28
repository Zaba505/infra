---
title: "Component: share-service"
type: docs
---

**Parent tech design:** [tech-design.md](../tech-design.md)
**Type:** Go service (`services/share-service/`)
**Established by:** ADR-0003

## Responsibility
Mints, audits, and revokes share grants. Issues short-lived signed URLs against the media bucket on behalf of grant holders.

## Data
**`share-grants` collection**:
- `grant_id`, `tenant_id` (granter), `grantee_id`, `object_id`, `expires_at`, `state` (`active`/`revoked`).

**`share-audit` collection** (append-only):
- `event_id`, `grant_id`, `event_type` (`granted`/`revoked`/`accessed`), `actor`, `at`.

## API surface
- `POST /api/v1/shares` — mint a grant (granter is the tenant, grantee is an end-user identifier).
- `GET /api/v1/shares/{grant_id}` — fetch a grant.
- `DELETE /api/v1/shares/{grant_id}` — revoke (writes audit event).
- `GET /api/v1/shares/{grant_id}/url` — return a signed URL for the grant's object (writes `accessed` audit event).
- `GET /api/v1/shares?tenant_id=...` — list grants for a tenant.

## Dependencies
- `cloud/firestore-tenants` provides the database; this component owns the `share-grants` and `share-audit` collections within it.
- `cloud/media-bucket` IAM must permit signed-URL minting by `share-service-sa`.

## Operational concerns
- Cloud Run service. Audit collection has append-only IAM; nothing can delete or update.
