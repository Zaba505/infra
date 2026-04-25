---
title: "Self-Hosted Personal Media Storage"
description: >
    A self-operated place where the operator, family, and friends can store and share personal media under their own control.
type: docs
weight: 10
---

> **One-line definition:** Provide a self-operated place where the operator, family, and friends can store and share their personal media (photos, videos, files) under their own control instead of paying a third-party cloud provider.

## Purpose & Business Outcome
What business outcome does this capability deliver? Why does it exist?

This capability exists so that a small, trusted circle of people (the operator, their family, and their friends) can keep their personal media without surrendering it to a commercial cloud provider. The outcomes it delivers, in order of importance:

1. **Privacy** — Users' content is not visible to, mined by, or monetized by any third party. The operator's circle decides what happens to their data.
2. **Longevity** — Content remains accessible and intact over long time horizons, not subject to a vendor's pricing changes, product sunsets, or account terminations.
3. **Control** — The operator decides how the system runs, who is on it, and how content is governed; users decide who sees their own content.
4. **Cost avoidance** — Avoiding ongoing per-GB subscription fees to a commercial cloud storage provider is a real but secondary benefit.

When these outcomes conflict, the order above is the tiebreaker. Privacy beats convenience; longevity beats operator convenience; control beats cost.

## Stakeholders

- **Owner / Accountable party:** The operator (Carson). Sole accountable party for the system existing, running, and continuing to run.
- **Primary actors (initiators):** Any authorized user — the operator, family members, or friends — uploading, viewing, sharing, or deleting their own content.
- **Secondary actors / consumers:** Other authorized users who have been explicitly granted access to a piece of content (individually or via a shared group/album).
- **Affected parties (impacted but not directly involved):** Subjects depicted in shared media (e.g. a family member who appears in a photo someone else uploaded) — they are affected by the privacy posture even if they are not the user who uploaded the content.

## Triggers & Inputs
What initiates the capability, and what information must be available?

- **Triggers:**
  - A user uploads a file (single, manual upload).
  - A user's device performs an automated backup (e.g. phone camera roll sync).
  - A user performs a bulk import from an existing provider (e.g. a Google Photos takeout archive).
  - A user views, downloads, or organizes their own content.
  - A content owner shares content with another user or a shared group.
  - A user deletes their own content.
- **Required inputs:**
  - An authenticated identity for the acting user.
  - The content itself (file bytes + whatever metadata the source provides, e.g. capture timestamps).
  - For sharing actions: the identity of the recipient user or shared group.
- **Preconditions:**
  - The user has been explicitly invited and provisioned by the operator (the user set is closed; no self-signup).
  - The user holds their own credentials. Lost credentials cannot be recovered (see Business Rules).

## Outputs & Deliverables
What does the capability produce? What changes in the world after it runs?

- **Direct outputs:**
  - The user's content is durably stored and retrievable by them.
  - Shared content becomes retrievable by the explicitly named recipients.
  - Content the user deleted is removed from their view immediately and purged from the system after the retention window (see Business Rules).
- **Downstream effects / state changes:**
  - The user (and the people they share with) can rely on this system as their primary store and stop paying a commercial provider for the same content.
  - The operator's circle accumulates a long-lived, private archive of personal media that is not dependent on any external vendor.

## Business Rules & Constraints

- **Closed user set.** Only the operator can add or remove users. There is no public sign-up. Users may be the operator, family members, or friends — chosen by the operator.
- **Private by default.** All content is private to its owner unless the owner explicitly shares it. Sharing may be one-to-one or via a shared group (e.g. a "family album"). No content is visible to other users — including the operator — without an explicit share.
- **Lost credentials = lost data.** If a user loses access to their credentials, their data is unrecoverable. The operator cannot reset access in a way that exposes the user's content. This is a deliberate Signal-style trade-off in service of the privacy outcome.
- **No storage quotas.** Users are not subject to per-user storage limits. Capacity planning is the operator's problem, not the users'.
- **No illegal content.** Users may not store content that is illegal in the operator's jurisdiction. The operator reserves the right to terminate a user's access on credible evidence of a violation, even though the operator cannot inspect user content directly.
- **Off-site backup is allowed.** Content may be replicated to an off-site location for durability, provided the off-site copy preserves the same privacy properties (no third party can read user content).
- **30-day retention after deletion / departure.** When a user deletes content, or when a user leaves the system entirely, their data is retained for 30 days and then purged. The window exists for accident recovery; after it elapses, the data is gone.
- **Operator succession.** Longevity is preserved through two complementary mechanisms: (a) every user can pull a complete on-demand archive of their own content at any time, without operator involvement, so that loss of the operator does not mean loss of access to already-uploaded content; and (b) a designated successor operator holds the credentials and runbook needed to keep the system running if the primary operator becomes unavailable. Exports preserve user data even if no successor takes over; the successor preserves continuity of the system itself. The successor's access is to operate the system, not to read user content — the privacy properties (private by default, lost credentials = lost data) continue to hold.

## Success Criteria & KPIs

- **Number of active users.** People in the operator's circle who actually use the system as a place to put their content (not just provisioned accounts that sit idle). Growth — or at minimum non-attrition — over time indicates the capability is meeting a real need versus their previous cloud provider.
- **Zero data loss.** No user ever loses content that they did not themselves delete. This is a hard binary: any incident of unintended data loss is a failure of the capability, regardless of cause.

## Out of Scope

- **General-purpose self-hosted application hosting** (e.g. running a Minecraft server, hosting other apps for the operator's circle). This is a separate capability and will be defined in its own document.
- **Public sharing.** Sharing content with anyone outside the operator's invited user set is not supported.
- **Collaborative editing of files.** This capability stores and shares media; it does not provide document editing, comments, or other collaboration features.
- **Monetization or commercial use.** The capability serves a private circle; it is not a product offered to the public.
- **Account/credential recovery.** Explicitly excluded by the privacy posture above.

## Open Questions

- **Definition of "active user"** for the KPI — what threshold of activity (uploads in the last N days? data stored? logins?) counts a user as active rather than dormant?
- **Affected-party recourse.** If someone depicted in a shared photo objects to its presence, is there a defined process for them to request its removal, given they may not be a user of the system themselves?
- **Operator's jurisdiction for "no illegal content."** Whose laws apply when users may live in different jurisdictions than the operator?
