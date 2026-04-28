---
title: "[0001] Tenant State Storage"
status: "accepted"
date: 2026-04-25
---

**Parent capability:** [self-hosted-application-platform](../_index.md)
**Addresses requirements:** TR-01, TR-04

## Context and Problem Statement
Where does tenant state live, given TR-01 (isolation) and TR-04 (no-downtime updates)?

## Considered Options
* Per-tenant Firestore namespace
* Shared Firestore with tenant-id partition key
* Per-tenant SQLite shard

## Decision Outcome
Chosen: per-tenant Firestore namespace.
