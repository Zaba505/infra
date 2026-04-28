---
title: "[0003] Share Grants"
status: "accepted"
date: 2026-04-25
---

**Parent capability:** [self-hosted-personal-media-storage](../_index.md)
**Addresses requirements:** TR-03, TR-06, TR-07

## Context and Problem Statement
How are end-user share grants represented and audited (TR-03, TR-06, TR-07)?

## Considered Options
- Signed URLs only, no central record
- Central share-grants table with signed URLs derived per-request
- ACL list embedded in each media object

## Decision Outcome
Chosen: central share-grants table owned by a `share-service`, with signed URLs derived per-request and an append-only audit log of grant/revocation events.

## Realization
- A `share-service` Go service under `services/share-service/`.
- Reuses the `cloud/firestore-tenants` collection grouping for the share-grants and audit-log collections.
