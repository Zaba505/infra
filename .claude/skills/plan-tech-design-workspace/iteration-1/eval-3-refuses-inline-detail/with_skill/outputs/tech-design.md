---
title: Tech Design
description: Composed tech design for the self-hosted application platform capability.
reviewed_at: null
---

## Overview

The self-hosted application platform hosts tenant capabilities on home-lab hardware fronted by Cloudflare and bridged to GCP via Wireguard. Tenants are first-class records owned by a `tenant-registry` service; each tenant is pinned to a platform contract version tracked by a `contract-version-catalog`; eviction produces an export bundle via the `tenant-eviction-exporter`; and the platform itself is brought up by a `platform-bootstrap` Terraform module. All component APIs return `application/x-protobuf` and surface errors as `application/problem+protobuf` per the shared error standard.

This document is a **skeleton**. Per-component contracts (endpoints, schemas, module inputs/outputs) live in the component-design docs filed via `define-component-design` and linked from the inventory below.

## Component diagram

```mermaid
flowchart LR
  Operator[Operator]
  Tenant[Tenant]
  Bootstrap[platform-bootstrap\n(Terraform module)]
  Registry[tenant-registry\n(service)]
  Catalog[contract-version-catalog]
  Exporter[tenant-eviction-exporter]
  ErrorPB[pkg/errorpb]

  Operator --> Bootstrap
  Operator --> Registry
  Tenant --> Registry
  Registry --> Catalog
  Registry --> Exporter
  Registry -.uses.-> ErrorPB
  Catalog -.uses.-> ErrorPB
  Exporter -.uses.-> ErrorPB
```

## Component inventory

| Component | Location | Type | Established by | Responsibility | Design doc |
|---|---|---|---|---|---|
| tenant-registry | `services/tenant-registry/` | service | ADR-0001 | Owns tenant records and the CRUD surface over them. | (issue pending — `define-component-design`) |
| contract-version-catalog | `services/contract-version-catalog/` | service | ADR-0002 | Tracks platform contract versions and the version each tenant is pinned to. | (issue pending — `define-component-design`) |
| tenant-eviction-exporter | `services/tenant-eviction-exporter/` | service/job | ADR-0003 | Produces the eviction export bundle when a tenant is evicted. | (issue pending — `define-component-design`) |
| platform-bootstrap | `cloud/platform-bootstrap/` | Terraform module | ADR-0001 (indirect, stand-up UX) | Provisions platform-level dependencies (network, storage, identity). | (issue pending — `define-component-design`) |
| pkg/errorpb | `pkg/errorpb/` | package (existing) | shared ADR-0007 | Standard problem/error protobuf types used by every component API. | (existing — no new design doc) |

> **Note on tenant-registry endpoints.** The user requested that `GET /tenants`, `POST /tenants`, `GET /tenants/{id}`, `PATCH /tenants/{id}`, and `DELETE /tenants/{id}` be written here. They are deliberately **not** included in this document. Endpoint contracts are a per-component-design concern; they will be authored in the tenant-registry component design doc, where pagination, PATCH semantics, ID derivation, idempotency, and error taxonomy can be reviewed in context.

## Data and state

Tenant records are the platform's primary state and are owned by `tenant-registry`. Contract-version pinning is owned by `contract-version-catalog` (storage location TBD — see gap G3). Eviction exports are produced by `tenant-eviction-exporter`; retention and pickup channel are TBD (gap G4). No component owns shared mutable state across services. Detailed schemas live in each component's design doc.

## Requirement realization

| TR | ADR(s) | Realized in |
|---|---|---|
| TR-NN | ADR-NNNN | component(s) |

> *This table is populated row-by-row from `tech-requirements.md`. Every TR-NN must appear with a non-empty ADR column and a non-empty component column. If any row would be empty in any column, it is a gap and is listed below rather than smoothed over.*

## Key flows

Per-UX sequence diagrams are deferred to the component design docs (or to dedicated follow-up issues if they cross multiple components). They are intentionally not authored in this skeleton.

## Gaps blocking completeness

| ID | Gap | Resolution type |
|---|---|---|
| G1 | Tenant ID derivation from the GitHub onboarding issue is unspecified. | Per-component spec via `define-component-design` (tenant-registry) |
| G2 | tenant-registry API contract (endpoints, pagination, PATCH semantics, idempotency) is unspecified. | Per-component spec via `define-component-design` (tenant-registry) |
| G3 | Storage location of tenant->contract-version pinning is unspecified. | Per-component spec via `define-component-design` (tenant-registry or contract-version-catalog) |
| G4 | Eviction export bundle format, retention, and pickup channel are unspecified. | Per-component spec via `define-component-design` (tenant-eviction-exporter) |
| G5 | TR-coverage walk to be re-verified at issue-filing time; any TR without an ADR will be filed back to `plan-adrs`/`define-adr`. | Amending ADR via `define-adr` (only if a TR is uncovered) |

This tech design is **not complete** until every gap above is resolved.

## Deferred / out of scope

- Per-component endpoint and schema detail (lives in component design docs).
- Cross-component sequence diagrams per UX (filed as separate follow-up issues if needed).
