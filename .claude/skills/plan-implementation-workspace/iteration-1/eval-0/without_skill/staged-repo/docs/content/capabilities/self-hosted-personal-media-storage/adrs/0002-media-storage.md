---
title: "[0002] Media Storage Backend"
status: "accepted"
date: 2026-04-25
---

**Parent capability:** [self-hosted-personal-media-storage](../_index.md)
**Addresses requirements:** TR-01, TR-02

## Context and Problem Statement
Where does the media itself live, given TR-01 (isolation) and TR-02 (single-region-outage survival)?

## Considered Options
- GCS multi-region bucket per tenant
- Single multi-region bucket with tenant-prefixed object paths
- Self-hosted MinIO

## Decision Outcome
Chosen: single multi-region GCS bucket with tenant-prefixed object paths, fronted by a `photo-store` service that enforces tenant isolation on read/write.

## Realization
- A `photo-store` Go service under `services/photo-store/`.
- A `cloud/media-bucket` Terraform module provisioning the GCS bucket and IAM.
