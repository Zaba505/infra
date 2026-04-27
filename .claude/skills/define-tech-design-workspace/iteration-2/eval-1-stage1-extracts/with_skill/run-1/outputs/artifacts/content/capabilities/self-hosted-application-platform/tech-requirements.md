---
title: "Technical Requirements"
description: >
    Technical requirements extracted from the self-hosted-application-platform capability and its user experiences.
type: docs
reviewed_at: null
---

> **Living document.** Numbering is append-only. ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to a date *newer* than this file's last modification once you have read and edited it. The `define-tech-design` skill will refuse to proceed to ADRs (Stage 2) until that condition holds.

**Parent capability:** [self-hosted-application-platform](_index.md)

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [Capability §Business Rules](_index.md#business-rules--constraints)
**Requirement:** The platform must enforce strict tenant isolation at the data and compute layers. No tenant workload may observe or access another tenant's state, secrets, traffic, or telemetry under any normal or degraded operating condition.
**Why this is a requirement, not a decision:** Capability business rules explicitly state tenant isolation as an invariant.

### TR-02: Operators must be able to roll out a platform-contract change without breaking existing tenants
**Source:** [UX: platform-contract-change-rollout](user-experiences/platform-contract-change-rollout.md)
**Requirement:** When the platform publishes a new contract version, existing tenants must continue to operate against the prior contract version until they migrate. The platform must support multiple contract versions concurrently for a bounded migration window.
**Why this is a requirement, not a decision:** The UX flow requires graceful contract evolution; otherwise rollouts break tenants.

### TR-03: Tenant-facing observability data must be queryable per-tenant within their data scope only
**Source:** [UX: tenant-facing-observability §Journey](user-experiences/tenant-facing-observability.md#journey)
**Requirement:** Tenants must be able to query metrics, logs, and traces for their own workloads. Cross-tenant observability data must be inaccessible to a tenant.
**Why this is a requirement, not a decision:** Both the capability isolation invariant and the UX journey require this.

### TR-04: Operator-initiated tenant updates must complete without tenant-perceived downtime for online workloads
**Source:** [UX: operator-initiated-tenant-update §Success](user-experiences/operator-initiated-tenant-update.md#success)
**Requirement:** When the operator initiates an update to a tenant (config, version, or capability), tenants serving online traffic must observe no end-user-visible downtime during the update.
**Why this is a requirement, not a decision:** UX success criteria define zero downtime as the user-perceived outcome.

### TR-05: A tenant evicted from the platform must be able to take their data with them
**Source:** [UX: move-off-the-platform-after-eviction](user-experiences/move-off-the-platform-after-eviction.md)
**Requirement:** The platform must provide an export mechanism by which an evicted tenant can retrieve all of their data in a portable format within a defined export window.
**Why this is a requirement, not a decision:** UX explicitly requires the move-off journey to succeed even for evicted tenants.

### TR-06: New tenants must be able to migrate existing data into the platform without loss or corruption
**Source:** [UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#journey)
**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).
**Why this is a requirement, not a decision:** UX requires lossless, retry-safe migration as part of the journey.

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** [CLAUDE.md §Repository Overview](/CLAUDE.md) · prior shared decision
**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.
**Why this is a requirement, not a decision:** Inherited topology constraint from the repo's architecture; not subject to revisiting at the capability level.

### TR-08: Tenants must be provisioned from a packaged artifact plus a declared resource specification
**Source:** [Capability §Triggers & Inputs](_index.md#triggers--inputs) · [UX: host-a-capability §Journey steps 4–7](user-experiences/host-a-capability.md#4-hand-off-packaged-artifacts)
**Requirement:** The platform must accept a packaged artifact (in the platform-defined packaging form) plus a declaration of the tenant's compute, storage, network reachability, and availability needs, and stand up a running tenant from those inputs alone. Operator hand-rolling per-tenant configuration outside the declarations is not permitted.
**Why this is a requirement, not a decision:** Both capability and UX require declarative onboarding; the choice of *which* packaging form is a Stage 2 decision.

### TR-09: An existing tenant's declared needs must be modifiable as a delta without re-provisioning from scratch
**Source:** [UX: host-a-capability §Journey step 8 (change-later loop)](user-experiences/host-a-capability.md#8-change-later-loop-re-entry)
**Requirement:** The platform must support changing a live tenant's storage size, external endpoints, components added/removed, and version bumps as incremental changes against the existing tenant rather than re-provisioning it from a clean state.
**Why this is a requirement, not a decision:** The change-later loop is a first-class part of the host-a-capability journey.

### TR-10: Eviction must produce a clean teardown on a chosen date with a 30-day read-only/export-only retention window
**Source:** [UX: move-off-the-platform-after-eviction §Phases B–C](user-experiences/move-off-the-platform-after-eviction.md#phase-b--the-eviction-date) · [Capability §Business Rules — Eviction](_index.md#business-rules--constraints)
**Requirement:** On the eviction date the platform must tear down the tenant's compute and network reachability, freeze tenant data into a read-only/export-only state, and permanently remove all retained tenant state exactly 30 days after the eviction date. The 30-day clock is hard except for a platform-bug carve-out where the export tooling itself fails.
**Why this is a requirement, not a decision:** The UX commits to specific phase mechanics and a numeric retention window.

### TR-11: Export tooling must produce an archive plus integrity metadata (checksum/hash and total size in bytes)
**Source:** [UX: move-off-the-platform-after-eviction §Phase A step 3](user-experiences/move-off-the-platform-after-eviction.md#3-run-the-export-and-verify-it-themselves) · [Capability §Operator succession](_index.md#business-rules--constraints)
**Requirement:** Each invocation of the platform's export tool must yield a downloadable archive of the tenant's data plus a checksum/hash and a total size in bytes, so the capability owner has a baseline of integrity metadata they can verify against.
**Why this is a requirement, not a decision:** The UX names checksum/hash + size explicitly as what the platform produces.

### TR-12: Export tooling must exist for every kind of data the platform hosts and must be available on demand while a tenant is healthy
**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases](user-experiences/move-off-the-platform-after-eviction.md#edge-cases--failure-modes) · [Capability §Operator succession](_index.md#business-rules--constraints)
**Requirement:** The export tool is a core, generic platform feature — there must be no tenant whose data shape is unsupported by export. Export must be runnable on demand by the capability owner whenever the platform is healthy, and must continue to function during the 30-day post-eviction grace window.
**Why this is a requirement, not a decision:** Capability operator-succession rule and UX both treat export as a first-class, always-present feature.

### TR-13: A migration job is a one-shot, capability-owner-supplied process the platform runs and tears down on completion
**Source:** [UX: migrate-existing-data §Journey](user-experiences/migrate-existing-data.md#journey)
**Requirement:** The platform must provide a one-shot job runner that accepts a packaged migration artifact (same packaging form as any tenant component), executes it against a live destination tenant with declared resource and network reachability, exposes its progress through the standard observability surface, and is torn down on completion. The platform does not retain the migration job afterwards.
**Why this is a requirement, not a decision:** UX specifies one-shot lifecycle, packaging parity, and tear-down explicitly.

### TR-14: The platform must offer a secret-management facility tenants register secrets into and reference by name
**Source:** [UX: migrate-existing-data §Journey step 1](user-experiences/migrate-existing-data.md#1-register-old-host-credentials-with-the-platform-secret-management-offering) · [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** The platform must provide a secret-management offering. Tenants register secrets out-of-band; tenant artifacts reference secrets by name only. Secret values must never appear on engagement-surface issues or in artifacts.
**Why this is a requirement, not a decision:** UX names the offering explicitly; capability outputs treat it as expected.

### TR-15: The platform must offer compute, persistent storage, network reachability (internal + external with TLS), backup/DR, and observability as named offerings
**Source:** [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** The platform must provide each of compute, persistent storage, internal network reachability, external (end-user-reachable) network reachability with TLS termination, backup and disaster recovery for tenant data, and observability as discrete offerings tenants opt into. Each offering is part of the platform's contract surface.
**Why this is a requirement, not a decision:** Capability lists these as direct outputs; the implementations are Stage 2 decisions but their existence is forced.

### TR-16: The platform-provided identity service must be capable of honoring "lost credentials cannot be recovered" (Signal-style)
**Source:** [Capability §Business Rules — Identity service honors tenant credential-recovery rules](_index.md#business-rules--constraints)
**Requirement:** The identity offering the platform provides to tenants must be able to support the property that a tenant can lose credentials with no recovery path. Identity options that cannot honor this property are ineligible to be the platform-provided identity service.
**Why this is a requirement, not a decision:** Explicit capability business rule. The choice of *which* identity implementation that satisfies this is a Stage 2 decision.

### TR-17: The platform must be (re)buildable from a definitions repo with a single top-level entry point
**Source:** [Capability §Success Criteria — Reproducibility](_index.md#success-criteria--kpis) · [UX: stand-up-the-platform §Journey step 2](user-experiences/stand-up-the-platform.md#2-kick-off-the-top-level-rebuild)
**Requirement:** A complete platform rebuild — first build, disaster recovery, or drill — must be initiated from a single top-level entry point against a definitions repo. No manual snowflake configuration step may be required for the platform to reach ready-to-host-tenants.
**Why this is a requirement, not a decision:** Capability KPI plus UX both require it; *which* tool drives the entry point is a Stage 2 decision.

### TR-18: Platform rebuild must be a phased automated flow with explicit operator-validation checkpoints between phases
**Source:** [UX: stand-up-the-platform §Journey steps 3–6](user-experiences/stand-up-the-platform.md#3-phase-1--foundations)
**Requirement:** The rebuild flow must execute in named phases (Foundations → Core services → Cross-cutting → Canary tenant) with the automation pausing at the end of each phase for the operator to validate and signal `continue` before the next phase begins.
**Why this is a requirement, not a decision:** UX prescribes the phase boundaries and the pause/validate/continue interaction.

### TR-19: Each rebuild phase must support clean teardown of everything provisioned so far
**Source:** [UX: stand-up-the-platform §Constraints — Each phase must be reversible](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [§Edge Cases — Phase fails](user-experiences/stand-up-the-platform.md#edge-cases--failure-modes)
**Requirement:** At every checkpoint, "delete everything provisioned so far" must be a viable, reliable rollback. Partial state is not trusted; the recovery model on phase failure is to tear down and restart from the top.
**Why this is a requirement, not a decision:** UX edge-case rule and inherited constraint.

### TR-20: Standup readiness is signalled by deploying and exercising a canary tenant end-to-end
**Source:** [UX: stand-up-the-platform §Journey step 6](user-experiences/stand-up-the-platform.md#6-phase-4--readiness-verification-and-canary-tenant)
**Requirement:** The platform reaches "ready to host tenants" only after a known-good canary tenant has been deployed end-to-end, exercised against compute, storage, identity, backup, and observability, and torn down. No earlier signal is sufficient.
**Why this is a requirement, not a decision:** UX is explicit that the canary's success is the readiness signal; *what* the canary is, is a Stage 2 decision.

### TR-21: The platform must enforce tracked changes and definition immutability so it does not drift from its definitions
**Source:** [UX: stand-up-the-platform §Constraints](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [Capability §Success Criteria — Reproducibility](_index.md#success-criteria--kpis)
**Requirement:** All platform mutations must occur through the tracked-change flow against the definitions repo. Ad-hoc modification of running infrastructure outside that flow must not be a possible operating mode.
**Why this is a requirement, not a decision:** Reproducibility KPI is unmeetable without it; UX names it as a constraint.

### TR-22: Backup and disaster recovery must cover the platform itself, not only tenants
**Source:** [UX: stand-up-the-platform §Journey step 5](user-experiences/stand-up-the-platform.md#5-phase-3--cross-cutting-services) · [Capability §Outputs & Deliverables](_index.md#outputs--deliverables)
**Requirement:** Backup and DR must be wired in during the cross-cutting phase of standup — before any tenant arrives — and must cover platform-level state in addition to tenant data.
**Why this is a requirement, not a decision:** UX places backup at Phase 3 explicitly so it covers the platform itself.

### TR-23: The only engagement surface between operator and capability owner is GitHub issues against the infra repo, with distinct issue types per UX
**Source:** Every UX (e.g. [host-a-capability §step 1](user-experiences/host-a-capability.md#1-file-an-onboard-my-capability-issue-on-github), [migrate-existing-data §step 2](user-experiences/migrate-existing-data.md#2-file-a-migrate-my-data-issue-on-github), [operator-initiated-tenant-update §step 1](user-experiences/operator-initiated-tenant-update.md#1-file-a-platform-update-required-issue-per-affected-tenant), [platform-contract-change-rollout §step 1](user-experiences/platform-contract-change-rollout.md#1-file-a-platform-contract-change-umbrella-issue)) · [Capability §Business Rules — Operator-only operation](_index.md#business-rules--constraints)
**Requirement:** All operator/capability-owner interactions must occur as GitHub issues against the infra repo. The platform must distinguish issue types — at minimum: `onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, `eviction` — because the operator's review scope and lifecycle differ per type. There is no self-service portal and no other front door.
**Why this is a requirement, not a decision:** Every UX consumes this as the only surface; the type taxonomy is fixed by the UX set.

### TR-24: End users of tenant capabilities must have no direct access to the platform
**Source:** [Capability §Business Rules — No direct end-user access](_index.md#business-rules--constraints) · [UX: move-off-the-platform-after-eviction §Constraints](user-experiences/move-off-the-platform-after-eviction.md#constraints-inherited-from-the-capability)
**Requirement:** End users of tenant capabilities must reach the tenant only — never the platform itself. The platform must not present pages, APIs, or notifications to end users; the platform's only consumers are the operator and capability owners.
**Why this is a requirement, not a decision:** Capability business rule, reinforced by UX.

### TR-25: Successor takeover must converge on the same standup flow without any platform-side change
**Source:** [Capability §Business Rules — Operator succession](_index.md#business-rules--constraints) · [UX: stand-up-the-platform §Persona](user-experiences/stand-up-the-platform.md#persona)
**Requirement:** A designated successor who has broken the sealed credentials and asserted takeover must be able to operate the platform — including running the standup flow — through the same definitions and entry points the primary operator uses. The platform must not depend on the primary operator's session, machine, or transient state to be (re)bootable.
**Why this is a requirement, not a decision:** Capability operator-succession rule treats the seal-breaking as out-of-band but requires the post-takeover flow to be identical.

### TR-26: Foundations must include cloud-↔-home-lab connectivity as a first-class element provisioned in Phase 1
**Source:** [UX: stand-up-the-platform §Constraints](user-experiences/stand-up-the-platform.md#constraints-inherited-from-the-capability) · [Capability §Business Rules — May span public and private infrastructure](_index.md#business-rules--constraints)
**Requirement:** The Foundations phase of rebuild must provision the cross-environment connectivity (cloud ↔ home-lab) as part of the underlying network plumbing — not as an afterthought layered on later phases.
**Why this is a requirement, not a decision:** Capability allows spanning public/private; UX places connectivity in Phase 1.

### TR-27: Tenant-facing observability access (login + alert channel) must be provisioned automatically as part of onboarding
**Source:** [UX: tenant-facing-observability §Entry Point](user-experiences/tenant-facing-observability.md#entry-point) · [UX: host-a-capability §Journey step 5](user-experiences/host-a-capability.md#5-wait-while-the-operator-provisions)
**Requirement:** When a tenant is onboarded, the platform must provision the capability owner's observability login (scoped to that tenant) and configured alert channel as part of provisioning — not as a separate later request.
**Why this is a requirement, not a decision:** UX explicitly says access is in place by the time the tenant goes live.

### TR-28: Capability owners must be able to self-serve their own alert thresholds without going through an issue
**Source:** [UX: tenant-facing-observability §Journey step 3](user-experiences/tenant-facing-observability.md#3-pull-mode-capability-owner-tunes-thresholds-if-needed)
**Requirement:** The observability offering must expose a self-service surface by which a capability owner adjusts the threshold values that govern their own push alerts. This is the one self-service surface the platform exposes; everything else still goes through GitHub issues.
**Why this is a requirement, not a decision:** UX explicitly carves out threshold tuning as the lone self-service exception.

### TR-29: Push alerts must name both the signal and the capability so the recipient can act without first opening anything else
**Source:** [UX: tenant-facing-observability §Journey step 4](user-experiences/tenant-facing-observability.md#4-push-mode-an-alert-reaches-the-capability-owner)
**Requirement:** Each alert delivered on the capability owner's chosen channel must identify which signal crossed which threshold and which capability/tenant it pertains to.
**Why this is a requirement, not a decision:** UX defines the alert payload's minimum contents.

### TR-30: Contract-change rollouts must support concurrent operation of old and new contract forms during the rollout window (full removals excepted)
**Source:** [UX: platform-contract-change-rollout §Journey step 3](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues)
**Requirement:** During a platform-contract-change rollout, both the old form and the new form of the affected offering must run concurrently on the platform until the deadline, so tenants can migrate at their own pace within the window. Full offering removals (no replacement) are the only exception; for those, the change is all-or-nothing at the deadline.
**Why this is a requirement, not a decision:** UX specifies the concurrent-rollout shape and its full-removal exception.

### TR-31: A platform-contract-change deadline must be uniformly enforced across all tenants, with no per-tenant slip
**Source:** [UX: platform-contract-change-rollout §Journey step 2](user-experiences/platform-contract-change-rollout.md#2-capability-owners-acknowledge-in-thread) · [§step 4](user-experiences/platform-contract-change-rollout.md#4-deadline-arrives)
**Requirement:** The deadline declared on a `platform contract change` umbrella issue must apply uniformly to every affected tenant. The platform's mechanics may not provide a way to extend or relax the deadline for a single tenant — global push only.
**Why this is a requirement, not a decision:** UX makes uniformity an explicit property of the journey.

### TR-32: Operator-initiated tenant-update issues are filed per-affected-tenant, with the deadline inherited from the external event that forced the update
**Source:** [UX: operator-initiated-tenant-update §Journey step 1](user-experiences/operator-initiated-tenant-update.md#1-file-a-platform-update-required-issue-per-affected-tenant)
**Requirement:** The `platform update required` flow must produce one issue per affected tenant (not a single umbrella). The deadline carried by each issue is the external deadline (vendor sunset, CVE remediation window, EOL date), not an operator-chosen date.
**Why this is a requirement, not a decision:** UX is explicit that this flow is one-issue-per-tenant and that deadlines are inherited.

### TR-33: The platform must provide a status-update mechanism for contract-change rollouts that emits migration metrics on a chosen cadence in the umbrella thread
**Source:** [UX: platform-contract-change-rollout §Journey step 3](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues)
**Requirement:** During a contract-change rollout, the platform (or operator tooling) must emit periodic status updates carrying migration metrics: how many tenants remain on the old form, how many have migrated, which `modify` issues are still open, and time remaining until the deadline. Cadence is chosen at filing time, sized to the rollout's overall length.
**Why this is a requirement, not a decision:** UX names the cadence-driven status-update mechanism explicitly.

### TR-34: Per-tenant migration in a contract-change rollout uses the existing `modify my capability` inner loop, linked to the umbrella issue
**Source:** [UX: platform-contract-change-rollout §Journey step 3](user-experiences/platform-contract-change-rollout.md#3-tenants-migrate-via-separate-modify-my-capability-issues) · [UX: operator-initiated-tenant-update §Journey step 3](user-experiences/operator-initiated-tenant-update.md#3-run-the-modify-inner-loop)
**Requirement:** The actual artifact handoff / re-provision / test / close work for each tenant in both the contract-change and operator-initiated-update flows must reuse the `modify my capability` inner loop rather than introducing a parallel inner-loop mechanism.
**Why this is a requirement, not a decision:** Both UXes explicitly delegate the inner loop to `modify my capability`.

### TR-35: The platform contract is evergreen — capability owners do not re-accept the contract on each modification
**Source:** [Capability §Business Rules — Tenants must accept the platform's contract](_index.md#business-rules--constraints) · [UX: host-a-capability §Journey step 8](user-experiences/host-a-capability.md#8-change-later-loop-re-entry)
**Requirement:** The platform must not require a capability owner to re-accept the contract on every `modify my capability` action. Contract acceptance is implicit on first onboarding and remains in force across modifications; contract changes themselves are operator-driven through the contract-change rollout flow.
**Why this is a requirement, not a decision:** Capability rule and UX both treat acceptance as evergreen.

### TR-36: Reproducibility KPI — full platform rebuild must complete within 1 hour of wall-clock time
**Source:** [Capability §Success Criteria — Reproducibility](_index.md#success-criteria--kpis)
**Requirement:** The end-to-end rebuild from no platform to ready-to-host-tenants must fit within 1 hour of wall-clock time. Missing the budget does not block the platform from going into service but generates a tracked follow-up issue.
**Why this is a requirement, not a decision:** Numeric capability KPI.

### TR-37: Operator maintenance budget KPI — routine operation must consume no more than 2 hours per operator per week
**Source:** [Capability §Success Criteria — Operator maintenance budget](_index.md#success-criteria--kpis)
**Requirement:** The platform's design and operations must keep routine operator effort at or below 2 hours per week. Tenants whose accommodation routinely costs disproportionate operator time approach the eviction threshold.
**Why this is a requirement, not a decision:** Numeric capability KPI; constrains every Stage 2 decision.

### TR-38: Eviction threshold — sustained accommodation cost above 2× the maintenance budget OR breaking reproducibility is grounds for eviction
**Source:** [Capability §Business Rules — Eviction threshold](_index.md#business-rules--constraints)
**Requirement:** The platform must support an operator-driven eviction path triggered when accommodating a tenant would either push routine operation sustainably above 2× the maintenance-budget KPI or break the reproducibility KPI. Either condition alone is sufficient grounds.
**Why this is a requirement, not a decision:** Capability rule defines the trigger; the *mechanism* by which eviction happens is documented in the move-off UX.

## Open Questions

These were surfaced during extraction and are seeds for Stage 2 ADRs (or for amendments to UX docs if a requirement turns out to be missing). They are *not* requirements.

- **Tenant packaging form.** OCI image, Helm chart, Kustomize, Nix derivation, repo-conforming-to-convention? Drives TR-08, TR-09, TR-13, TR-36.
- **Compute substrate.** Kubernetes (which distribution?), Nomad, raw systemd, container runtime only? Must satisfy TR-15 (compute), TR-17, TR-36, TR-37.
- **Persistent storage architecture.** Block vs object, per-tenant volumes vs shared filesystem, location (home-lab disks, cloud, both). Must satisfy TR-15 (storage), TR-22.
- **Network architecture details.** Ingress termination, DNS strategy, internal service discovery, reuse vs extend of the existing WireGuard tunnel from CLAUDE.md. Must satisfy TR-07, TR-15 (network), TR-26.
- **Identity service implementation.** Which self-hosted identity option satisfies TR-16's recovery-rules constraint?
- **Secret-management implementation.** Vault, sealed-secrets, cloud-provider KMS, SOPS-in-repo? Must satisfy TR-14, TR-21.
- **Backup architecture.** What gets backed up, where, retention, restore-drill cadence. Must satisfy TR-22, TR-36.
- **Observability stack.** Metrics + logs + traces; tenant-scoping mechanism for TR-03; alert delivery channels for TR-29; self-serve threshold UX for TR-28.
- **Definitions-repo layout and top-level rebuild entry point.** Terraform modules + a wrapper, Nix flake, Make? Must satisfy TR-17–TR-21, TR-36.
- **Canary tenant identity.** Purpose-built no-op tenant vs small real tenant. Drives TR-20.
- **Drift detection mechanism** as a precondition to rebuild. Gates TR-21.
- **Successor credential escrow implementation.** Out-of-band by definition, but must be reconcilable with TR-25.
- **Tenant data export shape.** Whether export is on-demand only or continuously available; format and packaging of the archive (drives TR-11, TR-12).
- **Contract versioning scheme.** Semver vs date-based vs other; how the platform represents "old form" vs "new form" of an offering during TR-30 concurrent operation.
- **Status-update format and storage** for TR-33 — comments on the umbrella issue, in-place edits to the issue body, an external dashboard?
