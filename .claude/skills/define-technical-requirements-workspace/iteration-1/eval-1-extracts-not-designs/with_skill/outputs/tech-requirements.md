---
title: "Technical Requirements"
description: >
    Technical requirements extracted from the Self-Hosted Application Platform capability and its user experiences. Each requirement links back to its source. Decisions belong in ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-tech-design` skill will refuse to proceed to ADRs (Stage 2) until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [Self-Hosted Application Platform](_index.md)

## How to read this

Each requirement is **forced** by the capability or a user experience — it constrains what the system must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review.

## Requirements

### TR-01: The platform must provide compute, persistent storage, network reachability (internal and external), identity/authentication for end users, backup/disaster recovery of tenant data, and observability as offerings consumable by tenant capabilities
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: stand-up-the-platform §Phase 2 — Core platform services](user-experiences/stand-up-the-platform.md#4-phase-2--core-platform-services) · [UX: stand-up-the-platform §Phase 3 — Cross-cutting services](user-experiences/stand-up-the-platform.md#5-phase-3--cross-cutting-services)
**Requirement:** The platform's surface to tenants is a fixed set of offering categories — compute, persistent storage, internal and external network reachability, end-user identity/authentication, backup and disaster recovery of tenant data, and observability. Every hosted tenant must be able to consume any subset of these without each tenant re-instrumenting them itself. The set is the platform's externally-visible contract for what it provides.
**Why this is a requirement, not a decision:** The capability's Direct Outputs section explicitly enumerates these. Specific implementations (which compute substrate, which storage technology, which identity protocol) are decisions; that the categories exist at all is forced.

### TR-02: The platform must be rebuildable end-to-end from version-controlled definitions, with no per-tenant snowflake state required to make it ready
**Source:** [Capability §Purpose & Business Outcome (Reproducibility)](_index.md#purpose--business-outcome) · [Capability §Success Criteria & KPIs (Reproducibility)](_index.md#success-criteria--kpis) · [UX: stand-up-the-platform §Journey](user-experiences/stand-up-the-platform.md#journey) · [UX: host-a-capability §Constraints Inherited from the Capability (KPI: 1-hour reproducibility)](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability)
**Requirement:** The entire platform — foundations, core services, cross-cutting services, and any per-tenant provisioning — must be expressible as definitions that, when applied from scratch to root-level access on the underlying infrastructure, produce a platform ready to host tenants. No step of standup or onboarding may require manual configuration that cannot be captured back into those definitions.
**Why this is a requirement, not a decision:** Reproducibility is named as the second-most-important business outcome and called out as a KPI. The standup UX makes it operationally testable.

### TR-03: A complete platform standup from definitions to "ready to host tenants" must complete in at most one hour of wall-clock time
**Source:** [Capability §Success Criteria & KPIs (Reproducibility)](_index.md#success-criteria--kpis) · [UX: stand-up-the-platform §Note the wall-clock and close out](user-experiences/stand-up-the-platform.md#7-note-the-wall-clock-and-close-out)
**Requirement:** The wall-clock duration from kicking off the top-level rebuild to the canary tenant going green must be no more than one hour. Exceeding this does not block the platform from going into service, but does generate a tracked follow-up.
**Why this is a requirement, not a decision:** The KPI states the operational form of "reproducible" as "at most 1 hour." The standup UX measures it.

### TR-04: Routine operation of the platform must be sustainable in no more than two hours per week of operator time
**Source:** [Capability §Success Criteria & KPIs (Operator maintenance budget)](_index.md#success-criteria--kpis) · [UX: host-a-capability §Constraints Inherited from the Capability](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability) · [UX: migrate-existing-data §Constraints Inherited from the Capability](user-experiences/migrate-existing-data.md#constraints-inherited-from-the-capability) · [UX: tenant-facing-observability §Constraints Inherited from the Capability](user-experiences/tenant-facing-observability.md#constraints-inherited-from-the-capability)
**Requirement:** All routine operator activities — onboarding, modify-loop reviews, migration runs, status updates during contract rollouts, alerting triage, drills — must collectively fit within a two-hour-per-week budget. Designs that force operator time disproportionate to this are non-compliant.
**Why this is a requirement, not a decision:** The KPI sets the numeric bound. Multiple UXs cite it explicitly as the constraint they must respect.

### TR-05: Tenants must be isolated such that one tenant's failure, load, or data cannot affect another tenant's availability or expose another tenant's data
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: tenant-facing-observability §Journey](user-experiences/tenant-facing-observability.md#journey) · [UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#4-operator-onboards-and-starts-the-migration-job)
**Requirement:** Compute, storage, network, identity, backup, and observability must be partitioned per tenant. A capability owner authenticating to the observability offering lands in their own tenant's view and stays confined there; concurrent migrations across different tenants do not leak resources or visibility between tenants; data of one tenant is never reachable by another.
**Why this is a requirement, not a decision:** The capability defines tenants as discrete consumers receiving discrete compute/storage/network/identity. Tenant-Facing Observability says capability owners cannot browse across tenants. The migration UX names concurrent multi-tenant runs as supported.

### TR-06: Only the operator may exercise administrative control over the platform; there is no co-operator surface and no delegated administration
**Source:** [Capability §Business Rules & Constraints (Operator-only operation)](_index.md#business-rules--constraints) · [UX: host-a-capability §Constraints Inherited from the Capability](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability) · [UX: tenant-facing-observability §Constraints Inherited from the Capability](user-experiences/tenant-facing-observability.md#constraints-inherited-from-the-capability)
**Requirement:** All administrative state changes — provisioning, deprovisioning, contract changes, eviction execution — must be reachable only by the operator's identity. No tenant-facing surface, no third-party surface, and no delegated-admin surface may grant administrative effects on the platform itself.
**Why this is a requirement, not a decision:** Stated as a hard business rule. Multiple UXs reaffirm "the operator is the only person who…".

### TR-07: A designated successor operator must be able to take over running the platform using sealed/escrowed credentials when the primary operator is unavailable, without those credentials being exercised during routine operation
**Source:** [Capability §Business Rules & Constraints (Operator succession)](_index.md#business-rules--constraints) · [UX: stand-up-the-platform §Persona](user-experiences/stand-up-the-platform.md#persona) · [UX: stand-up-the-platform §Constraints Inherited from the Capability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform must support a successor-operator handoff: there exist sealed/escrowed credentials sufficient to grant the successor full operator capability, and they are not used during normal day-to-day operation. After takeover, the successor uses the same operator surfaces; no separate "successor mode" is required.
**Why this is a requirement, not a decision:** Operator succession is named as a business rule with explicit sealed-credential mechanics; the standup UX makes the post-takeover flow identical to operator flow.

### TR-08: For every tenant, the platform must offer an on-demand export tool that produces a downloadable archive of the tenant's data along with a checksum/hash and total size, while the platform is healthy
**Source:** [Capability §Business Rules & Constraints (Operator succession)](_index.md#business-rules--constraints) · [UX: move-off-the-platform-after-eviction §Run the export and verify it themselves](user-experiences/move-off-the-platform-after-eviction.md#3-run-the-export-and-verify-it-themselves) · [UX: move-off-the-platform-after-eviction §Edge Cases](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes)
**Requirement:** The platform must provide a tenant-facing export tool, available at any time the platform is healthy, that produces a downloadable archive of that tenant's data accompanied by a platform-produced checksum/hash and total byte count. The tool must work for every kind of data the platform hosts — there is no tenant for which "no exporter exists" is an acceptable state. The generated archive itself is ephemeral; the tool must be re-runnable to regenerate it.
**Why this is a requirement, not a decision:** Operator-succession's "on-demand exportable archives" promise and the eviction UX's reliance on the same tool both force its existence. The eviction UX explicitly states export tooling is a core platform feature for every data shape.

### TR-09: After an eviction date, the evicted tenant's data must remain accessible to the capability owner via the export tool in a read-only state for 30 days, after which no tenant-accessible copy may remain
**Source:** [UX: move-off-the-platform-after-eviction §Phase B](user-experiences/move-off-the-platform-after-eviction.md#phase-b--the-eviction-date) · [UX: move-off-the-platform-after-eviction §Phase C](user-experiences/move-off-the-platform-after-eviction.md#phase-c--post-eviction-30-day-retention-window) · [UX: move-off-the-platform-after-eviction §Walk away](user-experiences/move-off-the-platform-after-eviction.md#7-walk-away)
**Requirement:** On the eviction date, compute and network for the evicted tenant must be torn down and the tenant's data must transition to a read-only state from which no further writes are accepted. For exactly 30 days afterward the export tool must continue to function against this frozen snapshot. After 30 days, the platform must offer no tenant-accessible copy of that data. (Whether deeper backup-tier copies persist beyond that point is an explicit open policy question — see Open Questions.)
**Why this is a requirement, not a decision:** The eviction UX defines this as a hard wall and a contract with the capability owner.

### TR-10: If the platform's export tooling or data hosting is shown to be at fault for a failed export, the 30-day retention countdown for that tenant must pause until a clean export can be produced
**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases (Export comes back wrong)](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes)
**Requirement:** The platform must support pausing an evicted tenant's tenant-accessible-data retention countdown when a platform-side defect (export tooling bug, hosting failure) prevents a clean export, and resuming or restarting it once a clean export has been produced. Capability-owner-side validation failures do not pause the clock.
**Why this is a requirement, not a decision:** Stated explicitly as the one exception to the 30-day hard wall in the eviction UX.

### TR-11: GitHub issues against this repository must be the sole engagement surface between capability owners and the platform; no other front door exists
**Source:** [UX: host-a-capability §File an "onboard my capability" issue on GitHub](user-experiences/host-a-capability.md#1-file-an-onboard-my-capability-issue-on-github) · [UX: host-a-capability §Constraints Inherited from the Capability](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability) · [UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#journey) · [UX: operator-initiated-tenant-update §Journey](user-experiences/operator-initiated-tenant-update.md#journey) · [UX: platform-contract-change-rollout §Journey](user-experiences/platform-contract-change-rollout.md#journey)
**Requirement:** All capability-owner-initiated and operator-initiated tenant-facing engagements (onboarding, modification, migration, operator-initiated update, contract change, eviction) must be coordinated through GitHub issues against the infra repository. No self-service portal, alternate channel, or non-issue mechanism is provided for these flows.
**Why this is a requirement, not a decision:** Every tenant-facing UX explicitly states GitHub issues are the channel. Host-a-Capability calls them "the only channel."

### TR-12: The platform must support distinct issue types for `onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, and eviction, and these types must be recognizable to participants
**Source:** [UX: host-a-capability §File an "onboard my capability" issue on GitHub](user-experiences/host-a-capability.md#1-file-an-onboard-my-capability-issue-on-github) · [UX: host-a-capability §Change-later loop](user-experiences/host-a-capability.md#8-change-later-loop-re-entry) · [UX: migrate-existing-data §File a "migrate my data" issue on GitHub](user-experiences/migrate-existing-data.md#2-file-a-migrate-my-data-issue-on-github) · [UX: operator-initiated-tenant-update §File a "platform update required" issue per affected tenant](user-experiences/operator-initiated-tenant-update.md#1-file-a-platform-update-required-issue-per-affected-tenant) · [UX: platform-contract-change-rollout §File a "platform contract change" umbrella issue](user-experiences/platform-contract-change-rollout.md#1-file-a-platform-contract-change-umbrella-issue) · [UX: move-off-the-platform-after-eviction §Entry Point](user-experiences/move-off-the-platform-after-eviction.md#entry-point)
**Requirement:** Distinct, named issue types must exist for each tenant-facing journey: `onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, and eviction. The distinct typing matters for the participants — review scopes differ by type, and downstream behaviors (eviction-issue cross-linking, umbrella issues) refer to type.
**Why this is a requirement, not a decision:** Each UX names its issue type and explains why distinct typing is needed.

### TR-13: The platform must accept a single, defined packaging form for tenant artifacts, identical across onboarding, modification, and migration-process handoff
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [Capability §Business Rules & Constraints (Tenants must accept the platform's contract)](_index.md#business-rules--constraints) · [UX: host-a-capability §Hand off packaged artifacts](user-experiences/host-a-capability.md#4-hand-off-packaged-artifacts) · [UX: migrate-existing-data §Entry Point](user-experiences/migrate-existing-data.md#entry-point) · [UX: migrate-existing-data §Constraints Inherited from the Capability](user-experiences/migrate-existing-data.md#constraints-inherited-from-the-capability)
**Requirement:** There is exactly one packaging form the platform accepts for tenant components, and the same form is used for migration-process artifacts. Tenants are responsible for packaging; the operator does not repackage on the tenant's behalf.
**Why this is a requirement, not a decision:** The capability rule says tenants must be packaged in the form the platform accepts. The migration UX explicitly mandates the *same* packaging form for migration processes.

### TR-14: When a tenant is onboarded, modified, or migrated, the tenant must declare its resource needs (compute, storage, network reachability, identity choice, availability expectations) up front, on the relevant issue
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [Capability §Business Rules & Constraints (Tenants must accept the platform's contract)](_index.md#business-rules--constraints) · [UX: host-a-capability §Entry Point](user-experiences/host-a-capability.md#entry-point) · [UX: migrate-existing-data §File a "migrate my data" issue on GitHub](user-experiences/migrate-existing-data.md#2-file-a-migrate-my-data-issue-on-github)
**Requirement:** The platform's review process must require declared resource needs as a first-class input. Hosting decisions and migration approvals depend on these declarations, including any temporary migration-only spikes.
**Why this is a requirement, not a decision:** A business rule. The UXs operationalize it as the issue contents and review scope.

### TR-15: A migration job's peak temporary footprint (steady-state plus declared spike) must be no more than 2× the destination tenant's steady-state compute and storage footprint, and approval depends on this check
**Source:** [UX: migrate-existing-data §Operator review on the issue](user-experiences/migrate-existing-data.md#3-operator-review-on-the-issue) · [UX: migrate-existing-data §Edge Cases](user-experiences/migrate-existing-data.md#edge-cases--failure-modes)
**Requirement:** The platform's migration-process review must enforce the 2× cap on peak temporary footprint relative to the destination tenant's steady-state. Migration requests above the cap are not approvable as written; the capability owner must split, reduce the spike, or resize the tenant first.
**Why this is a requirement, not a decision:** Numeric threshold stated explicitly in the migration UX.

### TR-16: The platform must offer a generic one-shot job runner ("migration-process offering") that runs tenant-supplied migration artifacts with the standard observability surface, and tears them down on completion
**Source:** [UX: migrate-existing-data §Operator onboards and starts the migration job](user-experiences/migrate-existing-data.md#4-operator-onboards-and-starts-the-migration-job) · [UX: migrate-existing-data §Operator tears down the migration job](user-experiences/migrate-existing-data.md#8-operator-tears-down-the-migration-job-and-closes-the-issue) · [UX: migrate-existing-data §Constraints Inherited from the Capability (KPI: 1-hour reproducibility)](user-experiences/migrate-existing-data.md#constraints-inherited-from-the-capability)
**Requirement:** The platform must provide a runner offering for one-shot tenant-supplied jobs (currently used for migration), with the standard tenant observability surface attached and explicit teardown after the job completes. The offering itself must be reproducible from definitions; specific jobs are per-tenant ephemera.
**Why this is a requirement, not a decision:** The migration UX names the offering and prescribes the lifecycle (runner exists, job is torn down, no history retained).

### TR-17: The platform must provide a tenant-facing secret-management offering through which capability owners register credentials that their packaged artifacts (e.g. migration processes) reference by name
**Source:** [UX: migrate-existing-data §Register old-host credentials with the platform secret management offering](user-experiences/migrate-existing-data.md#1-register-old-host-credentials-with-the-platform-secret-management-offering)
**Requirement:** A secret-management offering must exist such that capability owners can register secrets ahead of time, name them, and have their tenant artifacts read those secrets by name without exposing them on issues or in artifact contents.
**Why this is a requirement, not a decision:** The migration UX requires this surface as a precondition for filing a migration issue.

### TR-18: The platform's identity/authentication offering for end users must be capable of honoring a "lost credentials cannot be recovered" property
**Source:** [Capability §Business Rules & Constraints (Identity service honors tenant credential-recovery rules)](_index.md#business-rules--constraints) · [UX: host-a-capability §Constraints Inherited from the Capability](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability)
**Requirement:** Whatever identity option the platform offers tenants for end-user authentication must be technically capable of operating in a mode where lost credentials cannot be recovered (Signal-style). An identity option incapable of this property is ineligible to be the platform-provided identity service.
**Why this is a requirement, not a decision:** Stated as a hard business rule, motivated by an existing tenant's needs.

### TR-19: The platform must allow tenants to "bring their own identity" instead of using the platform-provided identity service, when declared in their tech design
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: host-a-capability §Entry Point](user-experiences/host-a-capability.md#entry-point)
**Requirement:** The platform must support tenants whose tech design declares a non-platform identity provider, without forcing them to use the platform-provided identity service. The choice is recorded at onboarding (in the tech design) and not re-litigated.
**Why this is a requirement, not a decision:** The capability lists "bring your own" as a tenant option in Triggers/Inputs and Outputs. The onboarding UX names it as a pre-arrival choice.

### TR-20: The tenant's observability access (login + email alerting + the platform-standard health bundle) must be provisioned automatically as part of onboarding, scoped to the tenant
**Source:** [UX: host-a-capability §Wait while the operator provisions](user-experiences/host-a-capability.md#5-wait-while-the-operator-provisions) · [UX: tenant-facing-observability §Access is already in place](user-experiences/tenant-facing-observability.md#1-access-is-already-in-place-set-up-during-onboarding) · [UX: tenant-facing-observability §Pull entry](user-experiences/tenant-facing-observability.md#entry-point)
**Requirement:** Onboarding must provision a working observability login scoped to the tenant, an email alerting channel wired to the capability owner's contact address, and the platform-standard health bundle (availability, latency, error rate, resource saturation, restart/deployment events) for every tenant. No further capability-owner action is needed for these to exist.
**Why this is a requirement, not a decision:** The observability UX explicitly states this provisioning happens during onboarding step 5 and is not opt-in.

### TR-21: The observability offering must be a single shared surface that confines each capability owner to their own tenant's view; only the operator may see across tenants
**Source:** [UX: tenant-facing-observability §Pull entry](user-experiences/tenant-facing-observability.md#entry-point) · [UX: tenant-facing-observability §Capability owner opens the observability view](user-experiences/tenant-facing-observability.md#2-pull-mode-capability-owner-opens-the-observability-view) · [UX: tenant-facing-observability §Constraints Inherited from the Capability](user-experiences/tenant-facing-observability.md#constraints-inherited-from-the-capability)
**Requirement:** The observability offering is one offering serving everyone; per-session scoping confines a capability owner to their own tenant's data with no UI-side mode switch into a wider view. Cross-tenant visibility is exclusive to the operator's identity.
**Why this is a requirement, not a decision:** Both the UX and the operator-only business rule force this shape.

### TR-22: Capability owners must be able to self-serve threshold tuning for their own tenant's email alerts within the observability offering
**Source:** [UX: tenant-facing-observability §Capability owner tunes thresholds](user-experiences/tenant-facing-observability.md#3-pull-mode-capability-owner-tunes-thresholds-if-needed)
**Requirement:** The observability offering must expose a tenant-scoped threshold-tuning surface that capability owners use without operator involvement. This is the single self-service surface the platform exposes to capability owners.
**Why this is a requirement, not a decision:** Explicitly named as the one exception to "everything goes through GitHub issues" — directly forced by the UX.

### TR-23: The observability offering must surface, in the tenant view, when email alert delivery is degraded for that tenant
**Source:** [UX: tenant-facing-observability §Capability owner tunes thresholds](user-experiences/tenant-facing-observability.md#3-pull-mode-capability-owner-tunes-thresholds-if-needed) · [UX: tenant-facing-observability §Edge Cases](user-experiences/tenant-facing-observability.md#edge-cases--failure-modes)
**Requirement:** When the offering knows email delivery is failing for a tenant, the tenant view must indicate that alerting is degraded so the capability owner does not interpret email silence as evidence of health. The pull view remains the source of truth.
**Why this is a requirement, not a decision:** The UX explicitly defines this trust contract.

### TR-24: Email-based alerting must be capable of delivering a per-signal, per-tenant alert that names the signal and the capability when a tenant's threshold is crossed
**Source:** [UX: tenant-facing-observability §An alert reaches the capability owner](user-experiences/tenant-facing-observability.md#4-push-mode-an-alert-reaches-the-capability-owner)
**Requirement:** When a tenant's threshold is crossed, the platform must emit an email to that tenant's capability owner identifying which signal fired and which capability is affected. The email is enough to begin investigation without opening the offering.
**Why this is a requirement, not a decision:** The UX prescribes this exact content shape; it constrains the alert payload, not the choice of delivery technology.

### TR-25: The platform standup must include a preflight drift check whenever prior platform state exists, and the rebuild must refuse to proceed if drift is detected
**Source:** [UX: stand-up-the-platform §Entry Point](user-experiences/stand-up-the-platform.md#entry-point) · [UX: stand-up-the-platform §Decide to rebuild and confirm preconditions](user-experiences/stand-up-the-platform.md#1-decide-to-rebuild-and-confirm-preconditions) · [UX: stand-up-the-platform §Constraints Inherited from the Capability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** Before any rebuild involving prior platform state, the platform must support comparing live (or last-known-good) state against the definitions and refusing to start the rebuild while unexplained differences remain. First-ever builds are vacuously clean.
**Why this is a requirement, not a decision:** The standup UX names this as a required preflight, not optional.

### TR-26: The platform's definitions and operations must enforce tracked changes and immutability, so that drift can be both prevented and detected
**Source:** [UX: stand-up-the-platform §Constraints Inherited from the Capability (Tracked changes and immutability)](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** Every UX that can introduce platform state must operate through tracked, immutable changes — ad-hoc modification outside the definitions is not permitted. The drift check (TR-25) is only meaningful if this is held.
**Why this is a requirement, not a decision:** The standup UX names this as a property the platform's definitions and operations must hold across all UXs.

### TR-27: The standup must execute as a single top-level entry point that runs in phases (foundations → core services → cross-cutting services → canary), pausing for explicit operator validation between phases
**Source:** [UX: stand-up-the-platform §Kick off the top-level rebuild](user-experiences/stand-up-the-platform.md#2-kick-off-the-top-level-rebuild) · [UX: stand-up-the-platform §Phase 1 — Foundations](user-experiences/stand-up-the-platform.md#3-phase-1--foundations) · [UX: stand-up-the-platform §Phase 4 — Readiness verification and canary tenant](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant)
**Requirement:** The standup automation must be driven by one top-level entry point, partitioned into the named phases, with a deterministic pause-and-validate checkpoint between each phase. The operator does not drive each step by hand, but does explicitly signal continuation.
**Why this is a requirement, not a decision:** The UX names the structure as the journey itself.

### TR-28: Each standup phase must be reversible — "tear down everything provisioned so far" must be a viable, reliable rollback at every checkpoint
**Source:** [UX: stand-up-the-platform §Edge Cases (Phase fails mid-rebuild)](user-experiences/stand-up-the-platform.md#edge-cases--failure-modes) · [UX: stand-up-the-platform §Constraints Inherited from the Capability (Each phase must be reversible)](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform's definitions must support a clean teardown of any partially-provisioned state at any checkpoint, so that on a phase failure the operator can return to an empty starting state and restart. Partial state is not trusted.
**Why this is a requirement, not a decision:** Stated as the rollback semantics in the UX.

### TR-29: A purpose-built canary tenant maintained alongside the platform definitions must be deployed end-to-end as the readiness gate; "ready to host tenants" cannot be declared from infrastructure self-checks alone
**Source:** [UX: stand-up-the-platform §Phase 4 — Readiness verification and canary tenant](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant) · [UX: stand-up-the-platform §Edge Cases (Canary tenant fails to come up)](user-experiences/stand-up-the-platform.md#edge-cases--failure-modes) · [UX: stand-up-the-platform §Constraints Inherited from the Capability (Default hosting target)](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** A canary tenant defined alongside the platform must be deployed as part of standup, exercise compute / reachability / storage read-back / authentication against the platform identity service / backup pickup / observability collection, then be torn down. Until the canary is green, the platform is not ready, regardless of phase-level signals.
**Why this is a requirement, not a decision:** The UX makes the canary a non-negotiable readiness gate.

### TR-30: Drills (full rebuilds on parallel scratch infrastructure) must be runnable on demand and scheduled at minimum quarterly and after every significant platform change, while the live platform continues serving
**Source:** [UX: stand-up-the-platform §Entry Point (Drift / reproducibility drill)](user-experiences/stand-up-the-platform.md#entry-point) · [UX: stand-up-the-platform §Constraints Inherited from the Capability (KPI: 1-hour reproducibility)](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform must support running the standup flow against scratch infrastructure (different from the live platform's underlying infra) without disturbing live operation, on demand and at least quarterly plus after each significant platform change. The drill is mechanically identical to a real rebuild.
**Why this is a requirement, not a decision:** The UX states the cadence and the parallel-infra requirement explicitly.

### TR-31: The platform may span both public-cloud and operator-owned (home-lab) infrastructure, with connectivity between them treated as part of the foundations
**Source:** [Capability §Business Rules & Constraints (The platform may span public and private infrastructure)](_index.md#business-rules--constraints) · [UX: stand-up-the-platform §Phase 1 — Foundations](user-experiences/stand-up-the-platform.md#3-phase-1--foundations) · [UX: stand-up-the-platform §Constraints Inherited from the Capability](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability)
**Requirement:** The platform's foundation must include both public-cloud and home-lab environments and the connectivity between them. Architectures that assume single-environment deployment are not compliant. Public-cloud components remain acceptable provided the operator retains control over configuration, data, and exit.
**Why this is a requirement, not a decision:** A business rule, reaffirmed by the standup UX.

### TR-32: During a platform-contract-change rollout, the old and the new form of an offering must run concurrently for the duration of the rollout window, except where the change is a full removal with no replacement
**Source:** [UX: platform-contract-change-rollout §Tenants migrate via separate `modify my capability` issues](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues) · [UX: platform-contract-change-rollout §Edge Cases (Full offering removal)](user-experiences/platform-contract-change-rollout.md#edge-cases--failure-modes)
**Requirement:** When a contract change retires an offering and replaces it with a new one, the platform must support running both forms simultaneously across the announced rollout window so tenants can migrate at their own pace before the deadline. Pure removals (no replacement) are exempt from concurrency.
**Why this is a requirement, not a decision:** The UX names concurrent old/new operation as a required property of the rollout.

### TR-33: The platform's contract must be evergreen: a tenant must not have to re-accept the contract on each modification, and contract changes must be communicated in advance and migrated, never sprung on tenants at runtime
**Source:** [UX: host-a-capability §Change-later loop](user-experiences/host-a-capability.md#8-change-later-loop-re-entry) · [UX: platform-contract-change-rollout §Constraints Inherited from the Capability (Evergreen contract)](user-experiences/platform-contract-change-rollout.md#constraints-inherited-from-the-capability)
**Requirement:** Modify-loop reviews touch only the delta. Contract changes are coordinated through an umbrella issue with a hard deadline, prior notice, and (where applicable) concurrent old/new running. The platform must not retire or alter a contract term in a way that is invisible to tenants until breakage.
**Why this is a requirement, not a decision:** Stated as a promise to capability owners in host-a-capability and operationalized by the contract-change UX.

### TR-34: A `platform contract change` rollout's status snapshot must live in the umbrella issue body and be updated on a chosen cadence, with each scheduled update also posted as a thread comment, carrying counts of migrated/un-migrated tenants, open `modify` issues, and time remaining
**Source:** [UX: platform-contract-change-rollout §Tenants migrate via separate `modify my capability` issues](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues)
**Requirement:** The rollout-status reporting must be present in two surfaces simultaneously — current state in the umbrella issue body for cold readers, and historical updates as thread comments for watchers — with the prescribed metrics on each update.
**Why this is a requirement, not a decision:** The UX prescribes both surfaces and the metric set.

### TR-35: A `platform contract change` rollout deadline must allow at least two full status-update cycles between filing and the deadline, when the deadline is not externally imposed
**Source:** [UX: platform-contract-change-rollout §File a "platform contract change" umbrella issue](user-experiences/platform-contract-change-rollout.md#1-file-a-platform-contract-change-umbrella-issue)
**Requirement:** When the operator chooses a deadline (proactive contract change), the deadline must be at least two full update cycles away from filing, so every tenant has a cycle to acknowledge/start and a cycle to finish or surface blockers.
**Why this is a requirement, not a decision:** The UX prescribes the cycle-count constraint explicitly.

### TR-36: When tenants miss a contract-change deadline or an operator-initiated update's operative delivery date, the platform's response must be to open a separate, linked eviction issue per laggard tenant — never silent breakage and never carrying the original issue indefinitely
**Source:** [UX: operator-initiated-tenant-update §Tip into eviction](user-experiences/operator-initiated-tenant-update.md#5-tip-into-eviction-after-the-last-workable-date-is-missed) · [UX: platform-contract-change-rollout §Deadline arrives](user-experiences/platform-contract-change-rollout.md#4-deadline-arrives)
**Requirement:** The platform's process must produce a separate eviction issue per laggard tenant, cross-linked to the originating update or umbrella issue, and the originating issue must close once eviction is in flight.
**Why this is a requirement, not a decision:** Both UXs prescribe the same eviction-handoff shape.

### TR-37: Eviction may only be triggered when accommodation would push routine maintenance sustainably above 2× the operator-maintenance-budget KPI or break the reproducibility KPI
**Source:** [Capability §Business Rules & Constraints (Eviction threshold)](_index.md#business-rules--constraints) · [UX: host-a-capability §Constraints Inherited from the Capability (Eviction threshold)](user-experiences/host-a-capability.md#constraints-inherited-from-the-capability) · [UX: operator-initiated-tenant-update §Tip into eviction](user-experiences/operator-initiated-tenant-update.md#5-tip-into-eviction-after-the-last-workable-date-is-missed)
**Requirement:** The eviction threshold is bounded by the two KPIs and must not be invoked merely because a tenant was late. The numeric thresholds inherit from the KPIs at any given time; this constraint is not restated in absolute hours so it cannot drift from them.
**Why this is a requirement, not a decision:** Stated explicitly as the eviction-threshold business rule.

### TR-38: The platform must not surface tenant-eviction or platform-maintenance state to end users; end users see only whatever connection behavior the underlying infrastructure produces
**Source:** [Capability §Business Rules & Constraints (No direct end-user access to the platform)](_index.md#business-rules--constraints) · [UX: move-off-the-platform-after-eviction §Edge Cases (End users keep hitting the tenant after the eviction date)](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes) · [UX: tenant-facing-observability §Out of Scope](user-experiences/tenant-facing-observability.md#out-of-scope)
**Requirement:** The platform has no notion of tenant end users and must not present platform-level status pages, retirement messages, or other communications to them. Communication with end users is exclusively the capability owner's responsibility.
**Why this is a requirement, not a decision:** Stated as a hard business rule and reaffirmed in multiple UXs.

### TR-39: The platform must provide backup of every tenant's data to a platform-defined standard, automatically, as part of being hosted
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables) · [UX: stand-up-the-platform §Phase 3 — Cross-cutting services](user-experiences/stand-up-the-platform.md#5-phase-3--cross-cutting-services) · [UX: stand-up-the-platform §Phase 4 — Readiness verification and canary tenant](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant)
**Requirement:** Backup of tenant data must be a cross-cutting service that engages automatically for every tenant, validated as part of standup readiness (the canary must be picked up by backup), and must apply to a platform-defined standard rather than per-tenant negotiation.
**Why this is a requirement, not a decision:** Direct outputs lists backup; the canary readiness check verifies it works.

### TR-40: The migration UX must support arbitrary, capability-owner-supplied migration logic without the platform inspecting or validating the logic itself; the operator reviews only the platform-side contract (resources, network, credentials, re-run shape)
**Source:** [UX: migrate-existing-data §Operator review on the issue](user-experiences/migrate-existing-data.md#3-operator-review-on-the-issue) · [UX: migrate-existing-data §Out of Scope](user-experiences/migrate-existing-data.md#out-of-scope)
**Requirement:** The platform's migration support is a runner; correctness, idempotency, and source-format handling are the capability owner's responsibility. The platform's review surface must therefore confine itself to the four named items: resources, network reachability, credential wiring, and re-run contract.
**Why this is a requirement, not a decision:** The migration UX defines the operator's review scope and explicitly excludes logic review.

### TR-41: Tenant onboarding must conclude with an explicit operator-driven test gate before the issue closes, and provisioning must not be considered complete until the capability owner confirms the deployed capability works
**Source:** [UX: host-a-capability §Test on request](user-experiences/host-a-capability.md#6-test-on-request) · [UX: host-a-capability §Operator closes the issue](user-experiences/host-a-capability.md#7-operator-closes-the-issue)
**Requirement:** The onboarding process must include a step in which the operator explicitly asks the capability owner to test, and the issue must remain open until the capability owner confirms or the parties iterate to a working state. Provisioning is not "done" until the test passes.
**Why this is a requirement, not a decision:** The UX prescribes this gate as part of the onboarding journey.

## Open Questions

Things the user volunteered as solutions during extraction (parked for Stage 2), or constraints the capability/UX docs don't yet make explicit.

- **Deeper backup-tier copy retention policy after the 30-day eviction window.** The eviction UX explicitly defers this: tenant-accessible copies are gone at day 30, but retention duration, deletion behavior, and operator-access/privacy constraints for any deeper backup-tier copies are TBD. Requirements derived from this will be added once a policy is set; until then, TR-09 only covers the tenant-accessible portion.
- **Tenant-facing pre-deprecation signal for `operator-initiated-tenant-update`.** The UX notes that if the platform ever adds an earlier deprecation/pending-update signal for capability owners, it would live in [Tenant-Facing Observability](user-experiences/tenant-facing-observability.md). Currently no such signal exists, so no TR is asserted; if/when it is added, new TRs will fall out of that UX update.
- **Status-update cadence sizing.** TR-34 requires status updates on a cadence chosen by the operator at filing time, and TR-35 requires at least two cycles before deadline; the UXs do not yet pin the lower bound on cadence frequency itself (e.g. minimum daily / maximum monthly). May warrant tightening once operational experience is gained.
- **Handling of capability-owner-driven concurrent migrations beyond two tenants at once.** Migration UX states concurrent migrations are supported but does not name a maximum or describe queuing/back-pressure behavior. May surface a TR (capacity guarantees, queue semantics) once the offering is implemented.
- **Backup standard.** TR-39 requires "a platform-defined standard" but the standard itself (RPO, RTO, retention, restore semantics) is not enumerated in the capability or UX docs. This is parked rather than asserted as a TR; defining the backup standard is itself a Stage-2 concern that will spawn additional TRs once defined.
- **Operator notification path during contract rollouts and operator-initiated updates.** Both UXs require the operator to chase silence, but the capability/UX docs do not specify the technical means by which the operator is notified when a tenant has acknowledged or shipped on an issue. Likely satisfied by GitHub's native notification surface, but not stated as a constraint anywhere yet.
