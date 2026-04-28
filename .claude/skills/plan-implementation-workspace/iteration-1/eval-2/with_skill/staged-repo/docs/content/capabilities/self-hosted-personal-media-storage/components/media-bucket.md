---
title: "Component: media-bucket module"
type: docs
---

**Parent tech design:** [tech-design.md](../tech-design.md)
**Type:** Terraform module (`cloud/media-bucket/`)
**Established by:** ADR-0002

## Responsibility
Provisions the multi-region GCS bucket holding all tenant media, with the IAM policy permitting `photo-store-sa` and `share-service-sa` the access shapes their components require.

## Inputs
- `project_id` (string)
- `location` (string, e.g. `nam-multi`)
- `photo_store_service_account` (string)
- `share_service_service_account` (string)

## Outputs
- `bucket_name`

## Resources
- `google_storage_bucket` with versioning + uniform bucket-level access.
- IAM bindings: photo-store SA → object create/get/delete; share-service SA → signed-URL minting.

## Operational concerns
- `create_before_destroy` lifecycle.
- Lifecycle rule deletes object versions older than 90 days unless tenant has retention override.
