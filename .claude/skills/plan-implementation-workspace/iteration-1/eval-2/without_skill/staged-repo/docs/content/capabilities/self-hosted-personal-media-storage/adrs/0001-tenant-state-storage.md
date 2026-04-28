---
title: "[0001] Tenant State Storage"
status: "accepted"
date: 2026-04-25
---

**Parent capability:** [self-hosted-personal-media-storage](../_index.md)
**Addresses requirements:** TR-01, TR-04

## Context and Problem Statement
Where does tenant state (registry of tenants, their identifiers, their lifecycle status) live, given TR-01 (isolation) and TR-04 (no-downtime onboarding)?

## Considered Options
- Per-tenant Firestore namespace
- Shared Firestore collection partitioned by tenant ID
- Per-tenant SQLite shard

## Decision Outcome
Chosen: shared Firestore collection partitioned by tenant ID, owned by a `tenant-registry` service that exposes a CRUD HTTP API.

## Realization
- A `tenant-registry` Go service under `services/tenant-registry/`.
- A Firestore collection provisioned via a `cloud/firestore-tenants` Terraform module.
