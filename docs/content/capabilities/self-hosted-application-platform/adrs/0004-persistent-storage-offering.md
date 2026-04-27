---
title: "[0004] Persistent Storage Offering — Block as the Foundational Primitive"
description: >
    The platform offers CSI-backed block storage as the foundational persistent storage primitive. Higher-level managed offerings (object, relational, graph, etc.) are separate platform offerings layered on top, each decided in its own ADR when a tenant first needs it.
type: docs
weight: 4
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform](../_index.md)
**Addresses requirements:** [TR-02](../tech-requirements.md#tr-02-provide-persistent-storage-as-a-tenant-offering), [TR-05](../tech-requirements.md#tr-05-provide-backup-and-disaster-recovery-of-tenant-data), [TR-15](../tech-requirements.md#tr-15-support-tenant-lifecycle-stage-live--eviction-frozen-computenetwork-deprovisioned-data-read-only--tenant-accessible-copy-removed-at-30-days), [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals), [TR-33](../tech-requirements.md#tr-33-routine-platform-operation-must-fit-within-2-hoursweek-of-operator-time)

## Context and Problem Statement

[ADR-0001](./0001-public-private-infrastructure-split.md) places bulk persistent storage on the home-lab side. [ADR-0002](./0002-compute-substrate.md) chose Kubernetes. This ADR decides what persistent-storage *offerings* the platform exposes to tenants and how the eviction-frozen state from [TR-15](../tech-requirements.md#tr-15-support-tenant-lifecycle-stage-live--eviction-frozen-computenetwork-deprovisioned-data-read-only--tenant-accessible-copy-removed-at-30-days) works against them.

There is real tension to resolve here. On one side, tenants will need different shapes of storage over time — object stores for media, relational DBs for transactional data, possibly graph or time-series later. On the other side, the parent capability's *evolves with its tenants* rule says the platform should grow when tenants need it, not before. Adding three managed offerings up front, ahead of any tenant who needs them, fails that rule by going ahead of the tenants and inflates [TR-33](../tech-requirements.md#tr-33-routine-platform-operation-must-fit-within-2-hoursweek-of-operator-time) routine maintenance for offerings nobody yet uses.

The reconciliation: pick a foundational primitive that *every* higher-level storage offering can be built on, and treat each higher-level offering as its own platform offering with its own ADR, brought in when the first tenant needs it.

## Decision Drivers

- **TR-02 — persistent storage as a tenant offering.** Some persistent storage must exist from day one, because no realistic tenant runs without persistence.
- **The capability rule "the capability evolves with its tenants."** The platform should grow when tenants demand it, not pre-emptively. Higher-level managed offerings should be ADRs in their own right.
- **TR-32 isolation.** Whatever the foundational primitive is, the substrate-level boundary (Kubernetes namespace + per-PVC scoping) must contain it. The fewer storage substrates in play, the easier this is.
- **TR-15 eviction-frozen state.** Block has a clean substrate-level answer (CSI volume snapshot, mount read-only); managed offerings will each define their own freeze semantics in their own ADRs.
- **TR-05 backup/DR.** Block backup is generic (snapshot + ship snapshot off-site). Managed-offering backup is offering-specific (logical dumps, bucket replication, etc.) and lives with each managed offering's ADR.
- **TR-33 (≤2 hr/week).** Each offering carries weekly cost. Adding a managed offering before a tenant is using it pays cost for nothing.
- **TR-17 (≤1 hr rebuild).** Each offering eats rebuild budget. Block alone keeps Phase 2 of the rebuild ([stand-up-the-platform §4](../user-experiences/stand-up-the-platform.md#4-phase-2--core-platform-services)) tight.

## Considered Options

### Option A — Block only as the foundational primitive; managed offerings are separate ADRs, born per tenant demand

The platform offers exactly one persistent-storage primitive on day one: **CSI-backed block storage**, exposed to tenants as PersistentVolumeClaims. Higher-level managed storage (object, relational, graph, time-series, etc.) is *not* part of this ADR; each becomes its own separate platform offering with its own ADR when a tenant first needs it. Until then, a tenant needing — say — an object store ships its own object-store image (against a block PVC) inside its own pod, and the moment a *second* tenant has the same need, the *capability evolves with its tenants* rule fires and a new managed-offering ADR is drafted.

- **Day-one offering surface:** small. One CSI driver to operate, one storage class, one snapshot/backup pattern.
- **TR-15 freeze:** clean — CSI snapshot at eviction, mount read-only for export, delete original PVC.
- **TR-05 backup:** generic — snapshot + ship to cloud archive (ADR #7).
- **TR-33:** lowest. No managed offerings paid for until a tenant uses one.
- **Tenant burden:** higher *for the first tenant of a given managed shape* — they ship their own object store / DB. Acceptable because by definition there is exactly one such tenant; the second triggers a managed-offering ADR.
- **Process correctness:** matches the *evolves with its tenants* rule precisely.

### Option B — Block + Object on day one, other managed shapes per ADR later

Same as A but ships an object-store offering (e.g. an S3-compatible service) up front because media-style tenants are already on the radar.

- Adds a second offering before any tenant has actually arrived to demand it.
- Pros: smoother day-one experience for the (currently-anticipated) media tenant.
- Cons: pays TR-33 cost for an offering with zero tenants; pre-decides product choice ahead of the actual tenant's tech-design submission, which is premature; fails the *evolves with its tenants* rule.

### Option C — Block + Object + Relational (e.g. via a Postgres operator) up front

- Three offerings on day one.
- Pros: most "ready" platform.
- Cons: three operational surfaces with no tenants on any of them; significant rebuild-budget hit; high TR-33 weekly cost; locks in DB version choices ahead of any tenant's need; the strongest violation of the *evolves with its tenants* rule.

### Option D — Block only, *and* refuse to add managed offerings later

The lean version: block forever, every tenant ships its own everything.

- Pros: smallest possible platform surface forever.
- Cons: actively contradicts the *evolves with its tenants* rule the moment a second tenant needs the same shape and the platform refuses to absorb the leverage; pushes maintenance back onto every tenant separately; produces N copies of "your own Postgres in a pod" — exactly the per-capability re-litigation the parent capability exists to avoid.

## Decision Outcome

Chosen option: **Option A — CSI-backed block storage is the foundational persistent-storage primitive. Higher-level managed offerings (object, relational, graph, time-series, etc.) are separate platform offerings, each with its own ADR, drafted when a tenant first needs them under the *evolves with its tenants* rule.**

This option is chosen because:

- It is the only option that fits the parent capability's *evolves with its tenants* rule cleanly. The platform grows when tenants pull it; it does not grow ahead of them.
- It keeps day-one TR-33 (2 hr/week) cost low. The rebuild budget (TR-17) absorbs one CSI driver and one storage class, not three operators with their own state.
- It gives every future managed-offering ADR (object, relational, graph, …) a known primitive to compose on: those offerings are not free-standing — their durable state lives on the same block primitive this ADR commits to.
- TR-15 eviction-freeze and TR-05 backup have a clean substrate-level pattern at this layer (snapshot + read-only mount + ship); managed offerings each refine those patterns in their own ADRs (e.g. "logical dump at freeze" for relational, "bucket policy flip" for object).
- The cost of the first tenant of any given managed shape — they ship their own object store or DB — is bounded and self-correcting: by definition there is exactly one such tenant, and the moment a second arrives the rule fires and an ADR is drafted.

This ADR commits only to the foundational primitive: **CSI-compatible block storage, snapshotable, exposed to tenants as PVCs**. The specific CSI driver / product (e.g. Longhorn, Rook-Ceph, local-path-provisioner with snapshot support, or another) is deferred to the deployment-time decision; this ADR pins the constraint as "CSI-compatible and snapshotable," not the product.

### Consequences

- **Good, because** the platform's day-one storage surface is exactly one offering, with one operational pattern. TR-33 is honored.
- **Good, because** TR-15 freeze and TR-05 backup at the block layer are simple and uniform — the export tool (TR-14) reads from a snapshot mounted read-only, regardless of what was on the volume.
- **Good, because** every future managed-offering ADR composes onto a stable foundation rather than re-deciding the foundation.
- **Good, because** the *evolves with its tenants* rule is honored mechanically: managed offerings exist iff tenants justify them.
- **Bad, because** the *first* tenant of any new managed shape pays the cost of shipping that shape inside their own pod. They knew this when they submitted their tech design (which named the platform as host), so it is a known cost rather than a surprise — but it is a real cost.
- **Bad, because** the platform now has a "second tenant of shape X" trip-wire that is partly social: the operator has to actually notice the pattern and draft the ADR, not let two tenants both ship their own Postgres pod. The host-a-capability "new offering needed" branch ([§3b](../user-experiences/host-a-capability.md#3-resolution--one-of-three-branches)) is the place this is meant to be caught, but absent a deliberate review the trip-wire can be missed.
- **Bad, because** managed-offering ADRs drafted later may discover that block alone is *not* the right substrate for them (e.g. a future graph DB might prefer a particular file-system layout or object-store-backed storage that itself rests on a different primitive). When that happens, that ADR is allowed to introduce a new substrate; this ADR does not promise that block is the *only* substrate forever, only that it is the day-one foundational primitive.
- **Requires:**
  - ADR #7 (backup & DR) defines how block snapshots reach the cloud archive — the offering-agnostic backup path.
  - ADR #11 (export tooling) consumes block snapshots for the eviction-frozen export.
  - ADR #12 (definitions tooling) provisions the CSI driver in the cluster as part of Phase 2 of the rebuild.
  - **Future managed-offering ADRs**, drafted as tenants demand them, will compose on this primitive. Each will own its own freeze, backup, and admission semantics in its own ADR.
  - The host-a-capability review (§2) gains an explicit check: "is this the second tenant shipping its own X-shaped storage in a pod?" If yes, draft a managed-offering ADR before approving.

### Realization

How this decision shows up in the repo:

- **The cluster's storage class** is a single `StorageClass` resource backed by the chosen CSI driver, expressed in the platform definitions and applied during Phase 2 of the rebuild ([stand-up-the-platform §4](../user-experiences/stand-up-the-platform.md#4-phase-2--core-platform-services)).
- **A `VolumeSnapshotClass`** is provisioned alongside, used by both backup (TR-05) and the eviction-freeze flow (TR-15).
- **Per-tenant PVCs** are emitted by the translator from ADR #3 against the tenant's declared storage need (TR-26). Quotas (TR-13 / TR-26) are expressed via `ResourceQuota` in the tenant's namespace.
- **The eviction-freeze flow**, when triggered by the eviction lifecycle (ADR #17), takes a `VolumeSnapshot` of every PVC in the tenant's namespace, deprovisions compute and network, and leaves the snapshot accessible to the export tool (TR-14) for 30 days.
- **Managed-offering ADRs**, when written, will land at `adrs/00NN-{name}-storage-offering.md` and will reference this ADR as the substrate they compose on.
- **A trip-wire note** lives in the host-a-capability review checklist (or equivalent operator-side runbook): if a tenant submission ships its own X-shaped storage and another tenant has done the same, draft a managed-offering ADR.

## Open Questions

- **Which CSI driver** (Longhorn, Rook-Ceph, local-path-provisioner with snapshots, OpenEBS, …). Deferred to the deployment-time decision; this ADR pins "CSI-compatible and snapshotable."
- **Single-node failure mode for block.** Whether home-lab block storage is replicated across nodes, or single-node-with-backup, depends on the CSI driver and is decided alongside the driver. Replicated is the convenience/resiliency-favoring choice; single-node-with-backup is the simpler one. Either is consistent with this ADR.
- **Future managed-offering ADRs** (object, relational, graph, time-series, …) are not yet drafted; they are *expected to exist* under the *evolves with its tenants* rule but their content is not in scope here.
- **What counts as "a second tenant of shape X."** The trip-wire is partly judgement — two tenants both running their own Postgres in a pod is clearly the trigger; two tenants both reading a config file from disk is clearly not. The operator's review during host-a-capability is where this judgement lives.
