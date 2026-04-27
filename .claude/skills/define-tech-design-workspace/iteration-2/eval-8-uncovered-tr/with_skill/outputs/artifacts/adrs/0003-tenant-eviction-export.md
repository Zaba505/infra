---
title: "[0003] Tenant Eviction Export"
status: "accepted"
date: 2026-04-25
---

**Parent capability:** [self-hosted-application-platform](../_index.md)
**Addresses requirements:** TR-05

## Context and Problem Statement
How do we satisfy TR-05 (evicted tenants can take their data)?

## Considered Options
* On-demand export to GCS bucket signed URL
* Continuous replication to tenant-controlled destination
* Scheduled snapshot with manual download

## Decision Outcome
Proposed: on-demand export to GCS bucket signed URL. Awaiting confirmation.
