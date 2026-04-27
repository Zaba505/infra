---
title: "Technical Requirements"
description: >
    Living extract of technical requirements for the Self-Hosted Application Platform capability, derived from the capability doc and its user experiences. This document is the input to ADRs and the composed tech design.
type: docs
weight: 1
status: draft
---

> **Status:** Draft — Stage 1 of the tech-design flow. This document extracts what the platform's design *must* satisfy, traceable line-by-line back to the capability doc and the seven user-experience docs. It does not make design decisions; those are captured as ADRs (Stage 2) and composed into the tech-design document (Stage 3).

**Parent capability:** [Self-Hosted Application Platform](../_index.md)

## How to read this document

Each requirement is tagged `[REQ-NN]`, given a one-line statement, and traced to its source(s) in the capability or UX docs. Requirements are grouped by concern, not by source document, because most requirements are reinforced by multiple sources.

A requirement here is something the platform's *technical implementation* must do or make true. Business rules are in the capability doc; UX flows are in the UX docs; this is the bridge.

## Functional requirements (what the platform must do)

### Tenant lifecycle

- **[REQ-01] Provision a tenant from a packaged artifact + declared resource needs.** The platform must accept a packaged artifact (form TBD — see ADR-01) plus declared compute / storage / network reachability / availability expectations and stand up a running tenant from them. Source: capability *Triggers & Inputs*, *Outputs & Deliverables*; `host-a-capability` steps 4–7.
- **[REQ-02] Modify an existing tenant's declared needs without re-provisioning from scratch.** Storage size, external endpoints, components added/removed, version bumps must be applicable as a delta. Source: `host-a-capability` step 8 (change-later loop).
- **[REQ-03] Evict a tenant cleanly on an operator-chosen date.** At the eviction date, compute and network reachability for that tenant go away; data is retained read-only for a 30-day grace window; after the grace window all tenant state is permanently removed. Source: `move-off-the-platform-after-eviction` Phases A/B/C.
- **[REQ-04] Run a one-shot, capability-owner-supplied migration job against a live tenant.** The job is packaged the same way as a capability component, has access to platform-managed secrets it names, and is torn down on completion. Source: `migrate-existing-data`.
- **[REQ-05] Coordinate operator-initiated tenant updates per-tenant** (one issue per affected tenant, deadline inherited from the external event). Source: `operator-initiated-tenant-update`.
- **[REQ-06] Coordinate proactive contract-change rollouts via a single umbrella tracking artifact** with per-tenant `modify` work items underneath. Source: `platform-contract-change-rollout`.

### Platform-provided offerings (the contract surface)

The platform must offer each of the following to tenants. Each is an "offering" that tenants opt into by declaring the need in their tech design:

- **[REQ-07] Compute offering.** A place for tenant application processes to run.
- **[REQ-08] Persistent storage offering.** Durable storage for tenant data, covered by [REQ-12] backup.
- **[REQ-09] Network reachability offering.** Both internal (tenant-to-tenant) and external (end-user-reachable). External reachability must terminate TLS and route to the correct tenant.
- **[REQ-10] Identity & authentication offering for tenant end users.** Optional per-tenant. Must be capable of honoring "lost credentials cannot be recovered" (Signal-style). Source: capability *Business Rules — Identity service honors tenant credential-recovery rules*.
- **[REQ-11] Secret management offering.** Tenants register secrets out-of-band; tenant artifacts reference them by name. Source: `migrate-existing-data` step 1 names this explicitly.
- **[REQ-12] Backup & disaster recovery offering for tenant data.** To a standard the platform defines. Restore on rebuild is part of a separate UX but the backups themselves must exist. Source: capability *Outputs & Deliverables*.
- **[REQ-13] Observability offering, tenant-scoped.** Capability owner can pull a tenant view (authenticated, scoped by tenant) and configure push alerts on a channel chosen at onboarding. Signal *categories* are platform-defined; specific signals per tenant are nailed down in that tenant's tech design. Source: `tenant-facing-observability`.
- **[REQ-14] Tenant data export tooling.** On-demand export producing an archive plus a checksum/hash and total size in bytes. Available while the tenant is healthy and during the post-eviction grace window. Source: capability *Operator succession*; `move-off-the-platform-after-eviction` step 3.

### Reproducibility and operations

- **[REQ-15] The entire platform must be (re)buildable from a definitions repo with a single top-level entry point.** No manual snowflake configuration. Source: capability *Reproducibility* KPI; `stand-up-the-platform` step 2.
- **[REQ-16] Rebuild is phased with operator-validation checkpoints between phases.** Phases: Foundations → Core services (compute/storage/identity) → Cross-cutting (backup/observability) → Canary tenant. Each phase pauses for operator `continue`. Source: `stand-up-the-platform` steps 3–6.
- **[REQ-17] Each phase must support a clean teardown.** "Delete everything provisioned so far" must be a viable rollback at every checkpoint. Source: `stand-up-the-platform` *Constraints — Each phase must be reversible*.
- **[REQ-18] A canary tenant must be deployable end-to-end as the readiness signal.** Exercises compute, storage, identity, backup, observability, then is torn down. Source: `stand-up-the-platform` step 6.
- **[REQ-19] All platform mutations must be tracked and definitions must be immutable.** Drift erodes the reproducibility KPI; the platform must not allow ad-hoc modification outside the tracked-change flow. Source: `stand-up-the-platform` *Constraints — Tracked changes and immutability*.
- **[REQ-20] Successor takeover must converge on the same standup flow.** Sealed credentials are escrowed externally and grant access to the operator's context; the rebuild path is identical thereafter. Source: capability *Operator succession*; `stand-up-the-platform` *Persona*.

### Engagement surface (operator ↔ capability owner)

- **[REQ-21] The only engagement surface is GitHub issues against the infra repo.** Distinct issue types: `onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, `eviction`. Source: every UX.
- **[REQ-22] No self-service portal; no co-operator delegation.** Source: capability *Operator-only operation*.
- **[REQ-23] No direct end-user access to the platform.** End users reach tenants, not the platform. Source: capability *Business Rules*.

## Non-functional requirements (how well, how much)

- **[NFR-01] Reproducibility KPI: full rebuild ≤ 1 hour wall-clock.** Missing it does not block the platform from going into service but generates a tracked follow-up issue. Source: capability *Success Criteria*.
- **[NFR-02] Operator maintenance budget KPI: ≤ 2 hours/week routine operation.** Tenants whose modify cadence routinely consumes disproportionate operator time cross into eviction territory. Source: capability *Success Criteria* + *Eviction threshold*.
- **[NFR-03] Eviction trigger threshold: 2× the maintenance budget OR breaks reproducibility.** Either condition alone is sufficient. Source: capability *Eviction threshold*.
- **[NFR-04] No specific availability SLA.** Platform offers whatever its current implementation delivers within the maintenance budget. Tenants needing more host elsewhere. Source: capability *Out of Scope*.
- **[NFR-05] Cost is secondary to convenience and resiliency** but should still be minimized where it does not cost either. Source: capability *Business Rules*.
- **[NFR-06] Post-eviction data grace window: 30 days, read-only/export-only.** Source: `move-off-the-platform-after-eviction` Phase C.

## Constraints (negative space — what the platform must not do, or must accept)

- **[C-01] May span public and private infrastructure.** "Self-hosted" means operator-controlled end-to-end, not "everything on operator-owned hardware." Public-cloud components allowed where operator retains control of config, data, and exit. Source: capability *Business Rules*.
- **[C-02] Foundations must include cross-environment connectivity (cloud ↔ home-lab).** Not an afterthought; provisioned in Phase 1 of standup. Source: `stand-up-the-platform` *Constraints*.
- **[C-03] No third-party / public / family-and-friends hosting.** Only the operator's own capabilities. Source: capability *Out of Scope*.
- **[C-04] Operator-skill-development is not a valid build-vs-buy tiebreaker.** Decisions judged on convenience, resiliency, and cost only. Source: capability *Business Rules*.
- **[C-05] Tiebreaker order:** tenant adoption > reproducibility > vendor independence > minimizing operator effort. Source: capability *Purpose*.
- **[C-06] Identity service must be capable of "lost credentials cannot be recovered."** Any candidate that cannot honor this property is ineligible to be the platform-provided identity offering. Source: capability *Business Rules*.
- **[C-07] Contract is evergreen.** Capability owners do not re-accept the contract on each modification; contract changes are operator-driven and migrate existing tenants (`platform-contract-change-rollout`). Source: capability *Business Rules*; `host-a-capability` step 8.
- **[C-08] Backup must cover the platform itself, not only tenants** (Phase 3 of standup wires it before any tenant exists). Source: `stand-up-the-platform` step 5.

## Open technical questions (to drive ADRs in Stage 2)

These are the questions whose answers become Stage 2 ADRs. Each names which REQs/Cs it must satisfy.

- **[Q-01] What is the tenant packaging form?** OCI image? Helm chart? Kustomize? Nix derivation? A repo conforming to a convention? Drives REQ-01, REQ-02, REQ-04, NFR-01.
- **[Q-02] What is the compute substrate?** Kubernetes (which distribution?), Nomad, raw systemd units, container-runtime-only? Must satisfy REQ-07, REQ-15, NFR-01, NFR-02, C-01, C-02.
- **[Q-03] What is the persistent-storage architecture?** Block? Object? Per-tenant volumes vs. shared filesystem? Where does it live (home-lab disks, cloud, both)? Must satisfy REQ-08, REQ-12, NFR-04, C-01.
- **[Q-04] What is the network architecture?** Ingress termination, DNS strategy, internal service discovery, cloud↔home-lab tunnel (Wireguard is already present in the repo's overall architecture per CLAUDE.md — confirm reuse vs. extend). Must satisfy REQ-09, C-01, C-02.
- **[Q-05] What identity service?** Self-hosted (Authentik, Keycloak, Ory, Zitadel, custom?) — all must be evaluated against C-06. Must satisfy REQ-10, C-06.
- **[Q-06] What secret-management offering?** Vault, sealed-secrets, cloud-provider KMS, SOPS-in-repo? Must satisfy REQ-11, REQ-19.
- **[Q-07] What backup architecture?** What gets backed up, where, retention, restore drill cadence. Must satisfy REQ-12, NFR-01 (restore must fit), C-08.
- **[Q-08] What observability stack?** Metrics + logs + traces; tenant-scoping mechanism; alert delivery channels. Must satisfy REQ-13, NFR-02 (operator must be able to see across tenants without it consuming the budget).
- **[Q-09] What definitions-repo layout and top-level rebuild entry point?** Terraform modules + a wrapper? Nix flake? Make? Must satisfy REQ-15–REQ-19, NFR-01.
- **[Q-10] What is the canary tenant?** Purpose-built no-op vs. small real tenant. Source: `stand-up-the-platform` Open Questions.
- **[Q-11] How is drift detected as a precondition to rebuild?** Source: `stand-up-the-platform` Open Questions; gates REQ-19.
- **[Q-12] How is successor credential escrow implemented?** Out-of-band by definition, but the platform must not depend on the operator's session to be re-bootable. Must satisfy REQ-20.

## Traceability matrix (capability/UX → requirements)

| Source | Drives |
|---|---|
| Capability *Outputs & Deliverables* | REQ-07–REQ-13 |
| Capability *Business Rules* | C-01, C-03–C-07, REQ-22, REQ-23, NFR-04, NFR-05 |
| Capability *Success Criteria* | NFR-01, NFR-02, NFR-03 |
| Capability *Operator succession* | REQ-14, REQ-20, Q-12 |
| `stand-up-the-platform` | REQ-15–REQ-20, C-02, C-08, Q-09–Q-11 |
| `host-a-capability` | REQ-01, REQ-02, REQ-21, C-07 |
| `migrate-existing-data` | REQ-04, REQ-11 |
| `move-off-the-platform-after-eviction` | REQ-03, REQ-14, NFR-06 |
| `operator-initiated-tenant-update` | REQ-05 |
| `platform-contract-change-rollout` | REQ-06, C-07 |
| `tenant-facing-observability` | REQ-13 |
