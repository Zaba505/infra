---
title: "Component: tenant-registry"
type: docs
---

**Parent tech design:** [tech-design.md](../tech-design.md)
**Type:** Go service (`services/tenant-registry/`)
**Established by:** ADR-0001

## Responsibility
Source of truth for tenant identity and lifecycle state. Owns the `tenants` Firestore collection.

## Data
**`tenants` collection** (Firestore document per tenant):
- `tenant_id` (string, primary key, derived from onboarding artifact per shared resource-identifier ADR)
- `display_name` (string)
- `state` (enum: `provisioning`, `active`, `evicted`)
- `created_at`, `updated_at` (timestamps)

## API surface
HTTP/protobuf, served on `HTTP_PORT`. Endpoints:
- `POST /api/v1/tenants` — create tenant. Validates display_name, returns canonical record.
- `GET /api/v1/tenants/{id}` — fetch tenant.
- `PATCH /api/v1/tenants/{id}` — update display_name or state.
- `DELETE /api/v1/tenants/{id}` — soft-delete (sets state=evicted).
- `GET /api/v1/tenants` — list, paged.

## Dependencies
- `cloud/firestore-tenants` module must exist before deployment.
- Reuses `pkg/errorpb` for error responses.

## Operational concerns
- Deploys as a Cloud Run service behind the platform's HTTPS load balancer.
- Health check on `HTTP_PORT`.
- Service account: `tenant-registry-sa`, with read/write to the `tenants` collection only.
