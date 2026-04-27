---
title: "ADRs"
description: >
    Capability-scoped Architecture Decision Records for the Self-Hosted Application Platform. Each ADR resolves one or more requirements in the capability's tech-requirements doc.
type: docs
weight: 10
---

These ADRs are scoped to the Self-Hosted Application Platform capability. Cross-cutting decisions (cloud provider, MADR format, network-boot architecture, etc.) live under [`docs/content/r&d/adrs/`](../../../../r&d/adrs/) and are inputs here, not outputs.

Each ADR cites the `TR-NN` requirements it satisfies from [`../tech-requirements.md`](../tech-requirements.md).

## Index

| ADR | Title | Drives |
|---|---|---|
| [0001](0001-tenant-isolation-model.md) | Tenant isolation model | TR-01, TR-03 |
| [0002](0002-compute-substrate.md) | Compute substrate | TR-01, TR-04, TR-07 |
| [0003](0003-contract-versioning-coexistence.md) | Platform-contract versioning and multi-version coexistence | TR-02 |
| [0004](0004-tenant-scoped-observability.md) | Tenant-scoped observability stack | TR-03 |
| [0005](0005-zero-downtime-tenant-updates.md) | Zero-downtime tenant update strategy | TR-04 |
| [0006](0006-tenant-data-export-and-import-format.md) | Tenant data export and migration import format | TR-05, TR-06 |
