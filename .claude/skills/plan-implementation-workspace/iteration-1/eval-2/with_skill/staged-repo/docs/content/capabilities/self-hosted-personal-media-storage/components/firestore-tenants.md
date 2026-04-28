---
title: "Component: firestore-tenants module"
type: docs
---

**Parent tech design:** [tech-design.md](../tech-design.md)
**Type:** Terraform module (`cloud/firestore-tenants/`)
**Established by:** ADR-0001, ADR-0003

## Responsibility
Provisions the Firestore database used by `tenant-registry` and `share-service`, and sets up the `tenants`, `share-grants`, and `share-audit` collections with appropriate indexes and IAM.

## Inputs
- `project_id` (string)
- `location` (string, e.g. `nam5`)
- `tenant_registry_service_account` (string)
- `share_service_service_account` (string)

## Outputs
- `database_name`
- `collections` map (`tenants`, `share_grants`, `share_audit`)

## Resources
- `google_firestore_database` (Native mode).
- `google_firestore_index` per query path (e.g. tenants by state, share-grants by grantee).
- IAM bindings: tenant-registry SA → read/write on `tenants`; share-service SA → read/write on `share-grants`, append-only on `share-audit`.

## Operational concerns
- `create_before_destroy` lifecycle to permit zero-downtime updates.
- Provider: GCP v7.11.0.
