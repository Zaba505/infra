---
title: "Self-Hosted Personal Media Storage"
description: A capability for storing personal media on self-hosted infrastructure.
type: docs
reviewed_at: 2026-04-26
---

A capability for capability owners to host personal media on the self-hosted application platform.

## Stakeholders

- **Operator** — runs the platform.
- **Capability owner** — brings a media collection to host.
- **End user** — consumes media (browses, views, shares).

## Triggers

- A capability owner has new media (photos, video) to bring online.
- An end user wants to access stored media or share it with someone.

## Outputs

- Durably stored media accessible via the platform.
- Per-tenant isolation of media collections.

## Business Rules

- Media is owned by a tenant; each tenant is a single capability owner's collection.
- End users access media only via shares granted by the owning tenant.

## Success Criteria

- Media survives single-region outages.
- End users can recover access after losing a device.
