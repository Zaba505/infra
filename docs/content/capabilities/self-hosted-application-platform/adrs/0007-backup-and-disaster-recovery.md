---
title: "[0007] Backup & Disaster Recovery"
description: >
    Hybrid backup — CSI VolumeSnapshots for fast in-cluster restore, plus file-level encrypted off-site backup to the cloud archive bucket. Single retention tier of 30 days for both live-tenant rolling backup and post-eviction retention.
type: docs
weight: 7
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform]({{< relref "../_index.md" >}})
**Addresses requirements:** [TR-05]({{< relref "../tech-requirements.md#tr-05" >}}), [TR-15]({{< relref "../tech-requirements.md#tr-15" >}}), [TR-22]({{< relref "../tech-requirements.md#tr-22" >}}), [TR-26]({{< relref "../tech-requirements.md#tr-26" >}}), [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}), [TR-33]({{< relref "../tech-requirements.md#tr-33" >}})

## Context and Problem Statement {#context}

[TR-05]({{< relref "../tech-requirements.md#tr-05" >}}) requires the platform to back up tenant data and be able to restore it, to a uniform standard the platform defines. This ADR picks the *mechanism*, the *destination*, the *cadence*, and the *retention policy*.

Two specific surfaces the design has to cover at this ADR's scope:

- **Block PVCs** ([ADR-0004]({{< relref "0004-persistent-storage-offering.md" >}})) — the foundational tenant storage primitive. Future managed-storage offerings will each define their own backup semantics in their own ADRs but ride this ADR's destination and retention shape.
- **Authentik's internal Postgres** ([ADR-0005]({{< relref "0005-identity-offering.md" >}})) — the platform's own identity state, treated as platform-state (not tenant data) for retention purposes.

This ADR also resolves the eviction UX's deferred *Open Question* about deeper-backup-tier retention beyond the 30-day tenant-accessible window ([UX §Open Questions]({{< relref "../user-experiences/move-off-the-platform-after-eviction.md" >}})). A deliberate retention choice closes that loop.

[ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}}) places the off-site archive at the cloud edge; [ADR-0006]({{< relref "0006-network-reachability.md" >}}) confirms the cloud↔homelab VPN's primary continuous use is exactly this archive-egress path.

## Decision Drivers {#decision-drivers}

- **[TR-05]({{< relref "../tech-requirements.md#tr-05" >}})** — must back up and must be able to restore.
- **[TR-15]({{< relref "../tech-requirements.md#tr-15" >}})** — eviction-freeze is *separate from* backup but shares the storage substrate. Backup must not be confused with eviction-freeze; they have different jobs.
- **Eviction UX promise** — *no tenant-accessible copy after day 30* ([UX §7 Walk away]({{< relref "../user-experiences/move-off-the-platform-after-eviction.md" >}})). The retention choice must honor this without hidden caveats.
- **[TR-32]({{< relref "../tech-requirements.md#tr-32" >}})** — backups are tenant data; one tenant must not be able to read another's backup, and backup credentials must not be brokered through tenant pods.
- **[TR-33]({{< relref "../tech-requirements.md#tr-33" >}})** — backup pipelines that need weekly hand-holding lose. Fewer tiers and a single retention horizon shrink ongoing operator work.
- **Capability rule "cost is secondary to convenience and resiliency."** Off-site bytes cost money. The cost is acceptable when paid for resiliency; it is not acceptable when paid for tiers nobody uses.
- **Capability rule "operator succession — exports preserve user data."** End-user data preservation is delivered via the export tool ([TR-14]({{< relref "../tech-requirements.md#tr-14" >}})) and tenants pulling their own archives, *not* via deeper-tier platform backups. This shapes how aggressive deeper retention needs to be.

## Considered Options {#considered-options}

### Backup mechanism

**M-A — CSI VolumeSnapshots replicated off-site.** Periodic CSI snapshots on the cluster, with a backup tool replicating snapshot bytes to the cloud archive bucket.
- Pros: snapshot semantics are well-understood; rapid in-cluster restore.
- Cons: replicating raw snapshot bytes off-site is opaque (no per-file deduplication, no encryption-at-rest controlled by the platform unless the tool adds it); restore requires the snapshot's CSI driver context.

**M-B — File-level backup via a restic/borg/kopia-class tool from inside backup pods that mount snapshots read-only.** No raw snapshot replication; the backup pod walks the filesystem on top of the snapshot and ships deduplicated, encrypted file-level data to the archive.
- Pros: deduplication and encryption are native; restore is portable across CSI drivers; archive contents are inspectable as files.
- Cons: file-level restore is not as fast as a CSI snapshot rollback for "I just deleted a thing"; the backup pod must mount each tenant's snapshot, which means the platform briefly handles tenant data in a backup namespace.

**M-C — Hybrid: CSI snapshot for fast in-cluster restore (24-48hr horizon), file-level via restic-class tool to cloud archive for off-site DR.**
- Pros: each mechanism does what it is good at; in-cluster fast restore for "oops I deleted a file" lives in the snapshot; full-DR off-site copy is encrypted and deduplicated; restore is portable.
- Cons: two mechanisms to operate, but they share the snapshot as a coordination point.

### Retention tiers

**T-α — Single tier of 30 days.** Tenant-accessible (same window for live-tenant rolling backups and post-eviction retention). After 30 days, deleted. No deeper tier.
- Pros: simplest possible policy; the eviction UX promise *no tenant-accessible copy after day 30* is realized literally — there is *no copy* after day 30, accessible to anyone or not. TR-33 minimal. Storage cost minimal.
- Cons: an evicted tenant who realizes 45 days post-eviction they're missing data has no path. A platform-side bug discovered after 30 days is unrecoverable from backup.

**T-β — Two tiers: tenant-accessible 30 days, operator-only deeper 90 days.** After 30 days, the tenant cannot reach it; for 90 days total, the operator can restore on request.
- Pros: a courtesy-restore window exists.
- Cons: the eviction UX deferred this exact question and the choice the operator now makes shapes the platform's privacy posture; storing tenant data 90 days past eviction is real exposure the operator must justify; complexity creeps into TR-33; tenants must understand "we said no tenant-accessible copy, but operator-restoreable copies exist for another 60 days."

**T-γ — As T-β plus indefinite cold archive for the platform's *own* state** (definitions snapshots, identity DB dumps, configs).
- Pros: rebuild from definitions assumes the definitions repo survives — a cold archive is an extra hedge.
- Cons: definitions live in git already; the hedge is small. Mostly added cost.

### Cadence

Setting cadence at this layer to **nightly for both block PVCs and Authentik Postgres**, with per-tenant manifest declarations allowed to request stricter RPO subject to admission. Stricter cadence is a [TR-26]({{< relref "../tech-requirements.md#tr-26" >}}) declaration field; the cap on how strict is deferred until a tenant pushes it.

## Decision Outcome {#decision-outcome}

Chosen options: **Mechanism M-C (hybrid CSI snapshot + file-level off-site)** + **Retention T-α (single tier, 30 days)**.

This pair is chosen because:

- **M-C** is the boring, well-supported pattern. CSI snapshots cover "I just deleted a thing" recovery and are the same mechanism the eviction-freeze flow uses ([ADR-0004]({{< relref "0004-persistent-storage-offering.md" >}}), [TR-15]({{< relref "../tech-requirements.md#tr-15" >}})) — coordination between backup and freeze is straightforward because they speak the same primitive. File-level off-site (via a restic-class tool) covers home-lab-loss DR with encryption, deduplication, and CSI-driver-portable restore.
- **T-α** honors the eviction UX's *no tenant-accessible copy after day 30* promise *literally* — there is no copy at all after 30 days, accessible or otherwise. This makes the platform's privacy posture explicit: the platform does not retain tenant data past the published window. End-user data preservation, when needed past 30 days, is delivered through the export tool ([TR-14]({{< relref "../tech-requirements.md#tr-14" >}})) and tenants pulling their own archives, exactly as the *operator succession — exports preserve user data* rule prescribes. Backup is for catastrophic-loss recovery within the 30-day window, not for post-window restore.
- T-β was rejected because storing tenant data past the published tenant-accessible window is real privacy exposure the operator would have to justify per tenant; the courtesy-restore convenience does not pay for the exposure or the explanation cost.
- T-γ was rejected because the definitions repo lives in git and the cold-archive hedge is small relative to its added cost. If the definitions repo's git history is ever insufficient as a recovery source, that is a *separate* problem with its own answer (e.g. mirroring git remotes), not a backup-tier issue.

The cadence is **nightly for both block PVCs and Authentik Postgres**. Per-tenant manifest declarations may request stricter RPO; the admission cap is deferred until a tenant pushes it.

The destination is the cloud-edge object store committed by [ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}}), reached over the cloud↔homelab VPN per [ADR-0006]({{< relref "0006-network-reachability.md" >}}).

### Consequences {#consequences}

- **Good, because** the eviction UX's published guarantee is realized literally. There is no "but actually we keep it for another 60 days" caveat to surface to evicted tenants.
- **Good, because** retention is a single horizon (30 days, everywhere, for everything tenant-scoped). [TR-33]({{< relref "../tech-requirements.md#tr-33" >}}) is honored: there is no tier mismatch to track or audit.
- **Good, because** CSI snapshots cover same-cluster restore quickly and serve double-duty for the eviction-freeze flow — one mechanism, two uses.
- **Good, because** off-site backup is encrypted, deduplicated, and CSI-driver-portable; the platform can be restored onto a different storage product without the backup tool caring.
- **Good, because** [TR-32]({{< relref "../tech-requirements.md#tr-32" >}}) is enforced at the bucket-prefix level: per-tenant prefixes in the archive bucket, scoped credentials per backup operation, no tenant pod ever holds archive credentials.
- **Bad, because** an evicted tenant who realizes 45 days post-eviction they're missing data has no path. The platform's response is "the export tool was available for the entire window; pull it again next time." This is the deliberate trade T-α makes and is consistent with the operator-succession rule that puts data preservation in the tenant's hands via export.
- **Bad, because** a platform-side bug or operator mistake discovered after 30 days is unrecoverable from backup. Mitigation: the 30-day window is wide enough that any reasonable detection cycle catches the issue; deeper events are recovered from the export archives tenants pulled or are not recovered at all.
- **Bad, because** Authentik's Postgres is also under T-α — a 31-day-old Authentik corruption is unrecoverable from backup. Mitigation: identity state is reconstructable from configuration (declarative, in the definitions repo) plus end users re-authenticating; tenant secrets in identity flows are not stored long-term in Authentik's DB.
- **Bad, because** the backup pipeline is two coordinated mechanisms. The coordination point (the CSI snapshot the file-level pod mounts) is the load-bearing piece; if the CSI driver's snapshot semantics regress, both restore paths are affected.
- **Requires:**
  - **A backup tool selection** (Velero with file-level plugins, restic-class tool driven by an in-cluster operator, or kopia-class equivalent). Pin "restic-class file-level backup with encryption-at-rest, deduplication, per-tenant scoping, and CSI-driver-portable restore." Product deferred.
  - **A cloud archive bucket** (per [ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}})) with: lifecycle rule deleting per-tenant prefixes 30 days after the last write, object-lock or similar immutability for in-window backups so a compromised cluster cannot delete its own backups, server-side encryption.
  - **Backup credentials** held in the secrets offering (ADR-0009 (forthcoming)), scoped per backup operation; never reachable from tenant pods.
  - **The export tool** (ADR-0011 (forthcoming)) reads from the same CSI snapshots the backup pipeline takes, so the eviction-freeze flow does not duplicate snapshot work.
  - **A restore runbook** documenting the two restore paths (in-cluster from snapshot, off-site from archive). The restore *itself* is operator workflow under [stand-up-the-platform]({{< relref "../user-experiences/stand-up-the-platform.md" >}}) — this ADR provides the mechanism.
  - **A `modify my capability` review checkpoint** for tenants declaring stricter-than-default RPO, until the admission cap is set.

### Realization {#realization}

How this decision shows up in the repo:

- **A `backup` namespace** in the home-lab cluster runs the chosen backup tool. The tool's controllers manage a schedule of `VolumeSnapshot`s per tenant namespace and corresponding off-site backup runs.
- **Per-tenant backup configuration** is emitted by the [ADR-0003]({{< relref "0003-tenant-packaging-form.md" >}}) translator from each tenant manifest's storage declarations: a snapshot schedule (default nightly, override in manifest), a per-tenant prefix in the archive bucket, and a per-tenant retention rule of 30 days that mirrors the bucket lifecycle rule for defense-in-depth.
- **The cloud archive bucket** is provisioned by [ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}})'s cloud-side definitions, with lifecycle, object-lock, and encryption configured.
- **Authentik's Postgres** has a `Job` (or equivalent operator-managed schedule) producing nightly logical dumps written to the same off-site path with a `platform-state/authentik/` prefix and the same 30-day retention.
- **The canary tenant** (ADR-0015 (forthcoming)) exercises a backup-and-restore round-trip during Phase 4 of the rebuild ([stand-up-the-platform §6]({{< relref "../user-experiences/stand-up-the-platform.md" >}})) so the pipeline is verified by every rebuild — backup that hasn't been restored from is not yet a backup.
- **The cloud↔homelab VPN** ([ADR-0006]({{< relref "0006-network-reachability.md" >}})) is the egress path for archive uploads; this is the VPN's primary continuous use.

## Open Questions {#open-questions}

- **Specific backup tool product.** Pin "restic-class file-level, encrypted, deduplicated, CSI-driver-portable restore."
- **Cloud archive bucket provider.** Inherits from the still-open question of which cloud provider hosts the cloud edge ([ADR-0001]({{< relref "0001-public-private-infrastructure-split.md" >}}) Open Questions).
- **Per-tenant cadence cap.** When/if a tenant manifests stricter-than-nightly RPO, the operator may need to set a hard cap. Deferred until pushed.
- **Restore-as-its-own-UX.** Tenant data restoration is named out-of-scope for [stand-up-the-platform]({{< relref "../user-experiences/stand-up-the-platform.md" >}}) and not yet a defined UX. This ADR provides the mechanism so that UX can exist; defining the UX itself belongs in a separate `define-user-experience` flow.
- **Definitions repo durability.** [Stage 3 will note] git is the source of truth and is not part of this ADR's backup scope; if definitions-repo durability turns out to need its own answer (e.g. mirroring git remotes), that is a separate decision.
