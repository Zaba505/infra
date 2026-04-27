---
title: "Technical Requirements"
description: >
    Technical requirements extracted from the Self-Hosted Application Platform capability and its user experiences. Each requirement links back to its source. Decisions belong in ADRs, not here.
type: docs
reviewed_at: 2026-04-27
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-tech-design` skill will refuse to proceed to ADRs (Stage 2) until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [Self-Hosted Application Platform](_index.md)

## How to read this

Each requirement is **forced** by the capability or a user experience — it constrains what the system must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review.

## Requirements

### TR-01: Provide compute as a tenant offering
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** The platform must give each hosted tenant a place for its application code to run. Compute is one of the platform's named direct outputs.
**Why this is a requirement, not a decision:** The capability lists compute as a direct output. Whether that compute is VMs, containers, functions, or something else is a Stage 2 decision.

### TR-02: Provide persistent storage as a tenant offering
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** Each tenant must have durable storage for its data, lasting across restarts and re-provisioning of compute.
**Why this is a requirement, not a decision:** Listed as a direct output. Storage *kind* (block, object, document, relational) is Stage 2.

### TR-03: Provide network reachability — internal between tenants and external for end users
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** Each tenant must be reachable both internally (from other tenants on the platform) and externally (by the tenant's own end users), at the network layer.
**Why this is a requirement, not a decision:** The capability explicitly distinguishes internal and external reachability as a single direct output. The mechanism (DNS, ingress, mesh, etc.) is Stage 2.

### TR-04: Provide an identity-and-authentication offering for end users that can honor "credentials cannot be recovered"
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [Capability §Business Rules — Identity service honors tenant credential-recovery rules](_index.md#business-rules--constraints)
**Requirement:** The platform must offer an identity-and-authentication service to tenants whose end users need to authenticate. Whatever implementation is offered must be capable of being configured so that lost credentials cannot be recovered (Signal-style). Tenants may opt out by bringing their own.
**Why this is a requirement, not a decision:** A named direct output, with a hard rule constraining the eligible implementations. Choice of identity product is Stage 2.

### TR-05: Provide backup and disaster recovery of tenant data
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** The platform must back up tenant data and be able to restore it. The standard the platform meets is platform-defined and uniform across tenants.
**Why this is a requirement, not a decision:** Listed as a direct output. Backup mechanism, retention, and RPO/RTO targets are Stage 2.

### TR-06: Provide an observability offering covering availability, latency, error rate, resource saturation, and restart/deployment events
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: tenant-facing-observability §1 Access is already in place](user-experiences/tenant-facing-observability.md#1-access-is-already-in-place-set-up-during-onboarding)
**Requirement:** The platform must observe each tenant such that the operator and the capability owner can both tell whether it is up and healthy without the tenant instrumenting itself. The platform-standard health bundle is fixed: availability, latency, error rate, resource saturation, and restart / deployment events.
**Why this is a requirement, not a decision:** Capability lists observability as an output; the UX fixes the bundle's contents. Choice of observability stack and storage is Stage 2.

### TR-07: Tenant-scoped observability access for capability owners; cross-tenant view only for the operator
**Source:** [UX: tenant-facing-observability §Pull entry](user-experiences/tenant-facing-observability.md#entry-point) · [UX: tenant-facing-observability §Constraints — Operator-only operation](user-experiences/tenant-facing-observability.md#constraints-inherited-from-the-capability) · [Capability §Business Rules — Operator-only operation](_index.md#business-rules--constraints)
**Requirement:** A capability owner authenticating to the observability offering must land in their own tenant's view and stay confined to it for the session. Cross-tenant browsing must be exclusive to the operator. There is no separate URL per tenant — one offering serves everyone with scope enforcement.
**Why this is a requirement, not a decision:** Forced by the operator-only rule and the tenant-facing-observability UX. Auth mechanism and scope-enforcement implementation are Stage 2.

### TR-08: Self-serve threshold tuning for capability-owner email alerts
**Source:** [UX: tenant-facing-observability §3 Capability owner tunes thresholds](user-experiences/tenant-facing-observability.md#3-pull-mode-capability-owner-tunes-thresholds-if-needed)
**Requirement:** Within the observability offering, a capability owner must be able to set, change, and remove the alert thresholds that decide when the platform emails them. This is the only capability-owner self-service surface; everything else still goes through GitHub issues.
**Why this is a requirement, not a decision:** The UX explicitly carves this out as the one self-service exception to the operator-only rule. Threshold storage and UI are Stage 2.

### TR-09: Email-channel push alerting with degraded-delivery indication
**Source:** [UX: tenant-facing-observability §1 Access is already in place](user-experiences/tenant-facing-observability.md#1-access-is-already-in-place-set-up-during-onboarding) · [UX: tenant-facing-observability §3 Capability owner tunes thresholds](user-experiences/tenant-facing-observability.md#3-pull-mode-capability-owner-tunes-thresholds-if-needed) · [UX: tenant-facing-observability §Edge Cases — Alert delivery is broken](user-experiences/tenant-facing-observability.md#edge-cases--failure-modes)
**Requirement:** The platform must push email alerts to capability owners when their thresholds are crossed. When the offering knows email delivery is degraded for a tenant, the tenant view must surface that so the capability owner does not treat email silence as evidence of health.
**Why this is a requirement, not a decision:** The UX names email as the channel the platform delivers and requires the offering to expose delivery health. Email-provider choice and degradation-detection mechanism are Stage 2.

### TR-10: Provide a packaging form the platform accepts for all tenant components
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [Capability §Business Rules — Tenants must accept the platform's contract](_index.md#business-rules--constraints) · [UX: host-a-capability §4 Hand off packaged artifacts](user-experiences/host-a-capability.md#4-hand-off-packaged-artifacts)
**Requirement:** The platform must define exactly one packaging form for tenant components. Capability owners hand off artifacts in this form; the platform consumes them as-is. The same form must also be acceptable for migration-process artifacts (TR-12).
**Why this is a requirement, not a decision:** Multiple UXs assume a single accepted form so that contract acceptance is unambiguous. Which form (container image, OCI bundle, archive layout, etc.) is Stage 2.

### TR-11: Provide a secret-management offering tenants can register secrets with and reference by name
**Source:** [UX: migrate-existing-data §1 Register old-host credentials](user-experiences/migrate-existing-data.md#1-register-old-host-credentials-with-the-platform-secret-management-offering) · [UX: migrate-existing-data §3 Operator review — Credentials](user-experiences/migrate-existing-data.md#3-operator-review-on-the-issue)
**Requirement:** The platform must offer secret management. Tenants register secret values out-of-band; tenant artifacts and migration artifacts reference them by name. Secret values must never appear in GitHub issues or other coordination surfaces.
**Why this is a requirement, not a decision:** The migration UX presupposes this offering and the name-based reference pattern. Secret store implementation is Stage 2.

### TR-12: Provide a one-shot migration-process offering that runs tenant-supplied migration jobs
**Source:** [UX: migrate-existing-data §Goal](user-experiences/migrate-existing-data.md#goal) · [UX: migrate-existing-data §4 Operator onboards and starts the migration job](user-experiences/migrate-existing-data.md#4-operator-onboards-and-starts-the-migration-job) · [UX: migrate-existing-data §Constraints — The capability evolves with its tenants](user-experiences/migrate-existing-data.md#constraints-inherited-from-the-capability)
**Requirement:** The platform must offer a runner for one-shot migration jobs that: (a) accepts the same packaging form as TR-10, (b) supports concurrent jobs across different tenants, (c) integrates with the secret-management offering (TR-11) so jobs read named secrets, (d) integrates with the observability offering (TR-06) so capability owners can watch progress, and (e) supports clean teardown of the job after success or abandonment without leaving residue.
**Why this is a requirement, not a decision:** The UX prescribes all of these properties. Implementation (Kubernetes Jobs, Cloud Run jobs, batch system) is Stage 2.

### TR-13: Admit migration jobs only when their peak temporary footprint is at most 2× the destination tenant's steady-state compute and storage
**Source:** [UX: migrate-existing-data §3 Operator review — Resources](user-experiences/migrate-existing-data.md#3-operator-review-on-the-issue) · [UX: migrate-existing-data §Edge Cases — Migration job needs more resources](user-experiences/migrate-existing-data.md#edge-cases--failure-modes)
**Requirement:** The platform must be able to express, per tenant, the steady-state compute and storage footprint, and to refuse a migration job whose declared peak footprint exceeds 2× that.
**Why this is a requirement, not a decision:** The UX names this exact threshold. Whether the check is a runtime quota, an admission webhook, or operator-side review tooling is Stage 2.

### TR-14: Provide export tooling, present for every tenant and every kind of data the platform hosts, producing archive + checksum/hash + total size
**Source:** [UX: move-off-the-platform-after-eviction §3 Run the export](user-experiences/move-off-the-platform-after-eviction.md#3-run-the-export-and-verify-it-themselves) · [UX: move-off-the-platform-after-eviction §Edge Cases — Export tooling does not exist](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes) · [UX: move-off-the-platform-after-eviction §Constraints — Operator succession — on-demand exportable archives](user-experiences/move-off-the-platform-after-eviction.md#constraints-inherited-from-the-capability) · [Capability §Business Rules — Operator succession](_index.md#business-rules--constraints)
**Requirement:** The platform must include export tooling that works for every tenant whose data the platform hosts, with no per-tenant special cases. Each invocation must produce a downloadable archive plus a checksum/hash and total size in bytes. Generated archives are ephemeral — produced for download then and there, not retained by the platform.
**Why this is a requirement, not a decision:** Multiple capability rules and the eviction UX force this. Archive format and tooling implementation are Stage 2.

### TR-15: Support tenant lifecycle stage `live` → `eviction-frozen` (compute/network deprovisioned, data read-only) → tenant-accessible copy removed at 30 days
**Source:** [UX: move-off-the-platform-after-eviction §Phase B](user-experiences/move-off-the-platform-after-eviction.md#phase-b--the-eviction-date) · [UX: move-off-the-platform-after-eviction §Phase C](user-experiences/move-off-the-platform-after-eviction.md#phase-c--post-eviction-30-day-retention-window)
**Requirement:** The platform must be able to move a tenant into a state where its compute and network are deprovisioned but its data persists in a read-only form against which export tooling still runs, and to remove the tenant-accessible copy of that data 30 days after entering that state.
**Why this is a requirement, not a decision:** The UX names this lifecycle and the 30-day window precisely. Storage mechanism for the read-only state is Stage 2.

### TR-16: Pause the 30-day retention countdown when the export-tooling failure is platform-side
**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases — Export comes back wrong / Export tooling does not exist](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes)
**Requirement:** The platform must support pausing and resuming the per-tenant 30-day retention-window countdown when the failure to produce a clean export is rooted in the platform's tooling or hosting (not in the capability owner's validation).
**Why this is a requirement, not a decision:** This is the one explicit exception to the "30-day hard wall" stated by the UX. How the pause is enacted is Stage 2.

### TR-17: Definitions-driven, single-entry-point rebuild of the platform end-to-end
**Source:** [UX: stand-up-the-platform §2 Kick off the top-level rebuild](user-experiences/stand-up-the-platform.md#2-kick-off-the-top-level-rebuild) · [UX: stand-up-the-platform §Constraints — KPI: 1-hour reproducibility](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [Capability §Success Criteria — Reproducibility](_index.md#success-criteria--kpis)
**Requirement:** The platform must be rebuildable from a single source of definitions via a single top-level entry point, with no manual snowflake configuration along the way. The full rebuild end-to-end must complete in ≤1 hour from "no platform" to "ready to host tenants."
**Why this is a requirement, not a decision:** KPI plus the UX's single-entry-point rebuild model. Tool choice (Terraform, Pulumi, Kubernetes operators, shell, mixture) is Stage 2.

### TR-18: Phased rebuild with operator-validatable checkpoints in fixed order: Foundations → Core (compute, storage, identity) → Cross-cutting (backup, observability) → Canary
**Source:** [UX: stand-up-the-platform §Journey](user-experiences/stand-up-the-platform.md#journey)
**Requirement:** The rebuild automation must pause between phases for operator validation, in this order: Phase 1 Foundations (cloud + home-lab base, networking between them); Phase 2 Core services (compute, storage, identity); Phase 3 Cross-cutting (backup, observability); Phase 4 Canary tenant exercise.
**Why this is a requirement, not a decision:** The UX prescribes the phases and the order. Pause/resume mechanism is Stage 2.

### TR-19: Each rebuild phase must be cleanly torn-downable (the partial state is itself untrusted)
**Source:** [UX: stand-up-the-platform §Edge Cases — Phase fails](user-experiences/stand-up-the-platform.md#edge-cases--failure-modes) · [UX: stand-up-the-platform §Constraints — Each phase must be reversible](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** "Delete everything provisioned so far" must be a viable, reliable operation at every phase boundary. There must be no state that, once partially created, cannot be cleanly destroyed without manual surgery.
**Why this is a requirement, not a decision:** Forced by the rebuild model. Implementation (lifecycle blocks, namespaces, labels, separate projects per attempt) is Stage 2.

### TR-20: Maintain a purpose-built canary tenant alongside the platform's definitions and use it as the readiness signal
**Source:** [UX: stand-up-the-platform §6 Phase 4](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant) · [UX: stand-up-the-platform §Constraints — Default hosting target](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform's definitions must include a canary tenant that exercises every platform offering end-to-end (run, reach, store/retrieve, authenticate via the identity offering, be picked up by backup and observability), and that can be deployed and torn down within the rebuild flow. The canary's success is the readiness signal — readiness cannot be declared from infrastructure self-checks alone.
**Why this is a requirement, not a decision:** The UX requires this exact construct. Canary content (which trivial app) is Stage 2.

### TR-21: Preflight drift check before any rebuild that has prior platform state
**Source:** [UX: stand-up-the-platform §Entry Point](user-experiences/stand-up-the-platform.md#entry-point) · [UX: stand-up-the-platform §1 Decide to rebuild](user-experiences/stand-up-the-platform.md#1-decide-to-rebuild-and-confirm-preconditions) · [UX: stand-up-the-platform §Constraints — Tracked changes and immutability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** Before any rebuild against prior platform state can begin, a preflight check must compare the live or last-known-good environment against the definitions and pass. The rebuild must refuse to start when unexplained drift is present.
**Why this is a requirement, not a decision:** The UX makes this a hard gate. Tooling and reference-state mechanism are Stage 2.

### TR-22: Tracked changes and immutability for all platform state-modifying actions
**Source:** [UX: stand-up-the-platform §Constraints — Tracked changes and immutability across all platform UXs](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [Capability §Success Criteria — Reproducibility](_index.md#success-criteria--kpis)
**Requirement:** Every UX that can introduce or change platform state must do so through a tracked, immutable mechanism (no ad-hoc console edits or untracked SSH changes). Drift, in the steady state, must be impossible-by-construction rather than detected-after-the-fact.
**Why this is a requirement, not a decision:** Required for reproducibility honesty. Mechanism (GitOps loop, console lockdown, audit policy) is Stage 2.

### TR-23: Single GitHub-issues engagement surface with distinct issue types per workflow
**Source:** [UX: host-a-capability §1 File an "onboard my capability" issue](user-experiences/host-a-capability.md#1-file-an-onboard-my-capability-issue-on-github) · [UX: host-a-capability §8 Change-later loop](user-experiences/host-a-capability.md#8-change-later-loop-re-entry) · [UX: migrate-existing-data §2 File a "migrate my data" issue](user-experiences/migrate-existing-data.md#2-file-a-migrate-my-data-issue-on-github) · [UX: operator-initiated-tenant-update §1 File a "platform update required" issue](user-experiences/operator-initiated-tenant-update.md#1-file-a-platform-update-required-issue-per-affected-tenant) · [UX: platform-contract-change-rollout §1 File a "platform contract change" umbrella issue](user-experiences/platform-contract-change-rollout.md#1-file-a-platform-contract-change-umbrella-issue) · [UX: move-off-the-platform-after-eviction §Entry Point](user-experiences/move-off-the-platform-after-eviction.md#entry-point)
**Requirement:** All capability-owner ↔ operator coordination must occur through GitHub issues against the infra repo. The repo must define distinct issue types covering, at minimum: onboard-my-capability, modify-my-capability, migrate-my-data, platform-update-required, platform-contract-change, eviction. The distinct types are themselves a load-bearing signal because review scopes and lifecycles differ across them.
**Why this is a requirement, not a decision:** Six UXs presuppose this surface and these distinct types. Templates, automation, and labels are Stage 2.

### TR-24: Tenant provisioning must run only through the platform's existing definitions
**Source:** [UX: host-a-capability §5 Wait while the operator provisions](user-experiences/host-a-capability.md#5-wait-while-the-operator-provisions) · [UX: host-a-capability §Constraints — KPI: 1-hour reproducibility](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability)
**Requirement:** Onboarding, modification, and contract-change-driven re-provisioning of a tenant must be expressible as edits to the platform's definitions. There must be no path that requires the operator to hand-roll bespoke per-tenant configuration outside the definitions repo.
**Why this is a requirement, not a decision:** Per the UX: bespoke tenant config breaks the reproducibility KPI. Mechanism (modules, declarative manifests) is Stage 2.

### TR-25: During platform-contract-change rollouts, run old and new forms of the offering concurrently until the deadline (except where the change is a full removal)
**Source:** [UX: platform-contract-change-rollout §3 Tenants migrate via separate `modify my capability` issues](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues) · [UX: platform-contract-change-rollout §Edge Cases — Full offering removal](user-experiences/platform-contract-change-rollout.md#edge-cases--failure-modes)
**Requirement:** When an offering is being replaced, the platform must be able to host the old form and the new form simultaneously until the rollout deadline. Permanent dual-form support is not a goal — only during the rollout window.
**Why this is a requirement, not a decision:** Forced by the contract-change UX. Coexistence mechanism (parallel deployments, namespaces, version selectors) is Stage 2.

### TR-26: Tenants declare resource needs at onboarding and on every modify; the platform admits or refuses based on those declarations
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [Capability §Business Rules — Tenants must accept the platform's contract](_index.md#business-rules--constraints) · [UX: host-a-capability §2 Operator review](user-experiences/host-a-capability.md#2-operator-review-on-the-issue)
**Requirement:** The packaged-artifact handoff must be accompanied by a declaration of the tenant's compute, storage, and network reachability needs (and migration spike, where applicable per TR-13). The platform's admission process must consume these declarations.
**Why this is a requirement, not a decision:** Forced by the contract. Declaration format (manifest fields, an issue-template schema) is Stage 2.

### TR-27: Span public-cloud and private/home-lab infrastructure, with the connectivity between them part of the foundation
**Source:** [Capability §Business Rules — The platform may span public and private infrastructure](_index.md#business-rules--constraints) · [UX: stand-up-the-platform §3 Phase 1 — Foundations](user-experiences/stand-up-the-platform.md#3-phase-1--foundations) · [UX: stand-up-the-platform §Constraints — The platform may span public and private infrastructure](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform's foundation must explicitly include connectivity between public-cloud and private/home-lab infrastructure. Tenants and platform offerings may be placed on either side, and reachability must work across the boundary.
**Why this is a requirement, not a decision:** The capability and rebuild UX both require this. Connectivity mechanism (VPN, private interconnect) and which offerings live where are Stage 2.

### TR-28: No direct end-user access to the platform itself
**Source:** [Capability §Business Rules — No direct end-user access to the platform](_index.md#business-rules--constraints) · [UX: move-off-the-platform-after-eviction §Constraints — No direct end-user access to the platform](user-experiences/move-off-the-platform-after-eviction.md#constraints-inherited-from-the-capability) · [UX: tenant-facing-observability §Constraints — No direct end-user access to the platform](user-experiences/tenant-facing-observability.md#constraints-inherited-from-the-capability)
**Requirement:** Platform-level surfaces (admin consoles, eviction notifications, observability) must not be exposed to end users of tenant capabilities. End users reach the tenant; the platform has no notion of "end users of itself."
**Why this is a requirement, not a decision:** Hard rule. Implementation (auth boundaries, network segmentation) is Stage 2.

### TR-29: Sealed/escrowed successor-credential mechanism that supports zero routine use
**Source:** [Capability §Business Rules — Operator succession](_index.md#business-rules--constraints) · [UX: stand-up-the-platform §Constraints — Operator succession](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The credentials that grant the designated successor administrative access to the platform must exist in a sealed/escrowed form that is not used during routine operation. The successor's takeover is a discrete event, not a sharing of day-to-day administration.
**Why this is a requirement, not a decision:** Capability rule. Sealing mechanism (password-manager handoff, physical envelope, KMS-wrapped key) is Stage 2.

### TR-30: Operator must have a cross-tenant view that capability owners do not
**Source:** [UX: operator-initiated-tenant-update §Entry Point](user-experiences/operator-initiated-tenant-update.md#entry-point) · [UX: operator-initiated-tenant-update §Constraints — Operator-only operation](user-experiences/operator-initiated-tenant-update.md#constraints-inherited-from-the-capability) · [Capability §Business Rules — Operator-only operation](_index.md#business-rules--constraints)
**Requirement:** The operator must be able to enumerate all tenants and learn which of them are using a given platform offering or component version. Capability owners must not have this view.
**Why this is a requirement, not a decision:** The fall-behind UX cannot be initiated without it. Implementation (a registry of tenants, label queries against the offering) is Stage 2.

### TR-31: Migration jobs declare a re-run contract that the platform records and respects
**Source:** [UX: migrate-existing-data §2 File a "migrate my data" issue](user-experiences/migrate-existing-data.md#2-file-a-migrate-my-data-issue-on-github) · [UX: migrate-existing-data §3 Operator review — Re-run contract](user-experiences/migrate-existing-data.md#3-operator-review-on-the-issue)
**Requirement:** The migration-process offering must accept and record a per-job declaration of whether the job is safe to run against an already-populated destination tenant or whether the destination must be empty before each run. The platform's review/admission step must consult this declaration.
**Why this is a requirement, not a decision:** Forced by the UX. How the declaration is captured and where it's stored are Stage 2.

### TR-32: Per-tenant authentication and isolation strong enough that no tenant (or its capability owner via the observability offering) can read another tenant's data or signals
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: tenant-facing-observability §2 Capability owner opens the observability view](user-experiences/tenant-facing-observability.md#2-pull-mode-capability-owner-opens-the-observability-view) · [Capability §Business Rules — Operator-only operation](_index.md#business-rules--constraints)
**Requirement:** The platform must enforce isolation between tenants such that one tenant's compute, storage, network, identity, secrets, observability data, and exports are not accessible to another tenant or to another tenant's capability owner.
**Why this is a requirement, not a decision:** Implicit but unavoidable: every offering's per-tenant scoping rule fails if isolation can be bypassed. Mechanism (namespaces, IAM, per-tenant projects, network policy) is Stage 2.

### TR-33: Routine platform operation must fit within ≤2 hours/week of operator time
**Source:** [Capability §Success Criteria — Operator maintenance budget](_index.md#success-criteria--kpis)
**Requirement:** Every offering and process this design produces must, in steady state, leave routine operator work at or under 2 hours per week across all hosted tenants.
**Why this is a requirement, not a decision:** KPI. Each Stage 2 ADR's options must be evaluated against this implicit budget.

## Open Questions

These were either solutions volunteered by the source docs (parked for Stage 2 decisions), or items the source docs themselves call out as undecided. Each one points at an ADR Stage 2 will need to draft.

- **OQ-1: Compute substrate.** What runs tenant workloads (VMs, container orchestrator, serverless)? Drives TR-01, and by transitive choice the packaging form (TR-10) and the migration runner (TR-12).
- **OQ-2: Packaging form.** What single artifact form does the platform accept (OCI image, Helm chart, OCI bundle, archive layout)? Drives TR-10, TR-12.
- **OQ-3: Persistent storage shape(s).** What kinds of storage does the platform offer (block, object, document, relational, multiple)? Drives TR-02, TR-05, TR-15.
- **OQ-4: Identity service implementation.** What identity product satisfies the "credentials cannot be recovered" eligibility rule? Drives TR-04.
- **OQ-5: Observability stack.** What metrics/logs/traces/alerting stack realizes the standard health bundle and the email-with-degraded-indicator path? Drives TR-06, TR-07, TR-08, TR-09.
- **OQ-6: Backup mechanism and retention.** What backs tenant data up, where does it live, and how long is it kept? Drives TR-05. Note: the eviction UX explicitly defers the "deeper backup-tier policy beyond the 30-day tenant-accessible window" as TBD ([UX §Open Questions](user-experiences/move-off-the-platform-after-eviction.md#open-questions)).
- **OQ-7: Public/private split.** Which offerings live in the cloud, which in the home-lab, and what connects them? Drives TR-27. The capability allows either; the cost/convenience/resiliency tiebreaker decides per offering.
- **OQ-8: Definitions tooling.** What language/tool drives the rebuild, the per-tenant provisioning, and the immutability enforcement (Terraform, Pulumi, a Kubernetes operator, a mix)? Drives TR-17, TR-18, TR-19, TR-22, TR-24.
- **OQ-9: Drift detection.** What tool produces the preflight drift check and what does it compare against? Drives TR-21.
- **OQ-10: Eviction-frozen storage state.** How is tenant data held in a read-only, exportable-but-not-writable form for 30 days? Drives TR-15, TR-16.
- **OQ-11: Export tooling shape.** Per data shape, how does export run, where does it produce its archive, and how does the capability owner download it before it goes ephemeral? Drives TR-14.
- **OQ-12: Tenant registry / cross-tenant view.** What records the set of hosted tenants and their declared resource needs (TR-26), and powers the operator's cross-tenant view (TR-30) and the contract-change rollout's "who is still on the old form" tracking ([UX: platform-contract-change-rollout §3](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues))?
- **OQ-13: Issue-type mechanics on GitHub.** Are the six issue types (TR-23) GitHub Issue Forms, labels, a project-board partition, or something else? Templates and automation hang off this.
- **OQ-14: Successor-credential seal.** What sealing mechanism (password-manager export, KMS-wrapped key in a sealed envelope, etc.) holds the successor credentials per TR-29?
- **OQ-15: Pending-update tenant signal (deferred).** The fall-behind UX explicitly defers a tenant-facing pending-update view to a possible future expansion of the observability offering ([UX: operator-initiated-tenant-update §Out of Scope](user-experiences/operator-initiated-tenant-update.md#out-of-scope)). Not in scope for this design unless the user pulls it in.
