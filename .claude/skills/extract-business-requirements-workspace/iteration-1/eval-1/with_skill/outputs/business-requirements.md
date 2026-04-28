---
title: "Business Requirements"
description: >
    Business requirements extracted from the Self-Hosted Application Platform capability and its user experiences. Each requirement links back to its source. Technical requirements and decisions belong in tech-requirements.md and ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from the capability and UX docs on demand. Numbering is append-only — once a BR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). Technical requirements cite BR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `define-technical-requirements` skill will refuse to extract TRs until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [Self-Hosted Application Platform]({{< ref "_index.md" >}})

## How to read this

Each requirement is **forced** by the capability or a user experience — it states, in business or user-outcome terms, what the system must guarantee. Decisions about the *technical translation* (cadences, durability levels, protocols) belong in `tech-requirements.md`. Decisions about *how* (which database, which library, which provider) belong in `adrs/`. If something in this list reads like a technical constraint or a chosen solution rather than a business demand, flag it for review.

## Requirements

### BR-01: The platform must be the default hosting target for the operator's capabilities
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md" >}}) · [Capability §Business Rules & Constraints]({{< ref "_index.md" >}})

**Requirement:** Any capability the operator defines must be able to run on this platform by default, so that "where does this run?" is a solved question per capability rather than re-litigated each time. A capability hosted elsewhere counts as a failure of the platform to meet that capability's needs (or a failure to ask).

**Why this is a requirement, not a TR or decision:** The capability's first stated outcome explicitly demands a default hosting answer; the *Tenant adoption* KPI converts that demand into a measurable success criterion that counts elsewhere-hosted implemented capabilities negatively. It is a business-level promise about the platform's role, not a technical translation.

### BR-02: The platform must be rebuildable from its definitions
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md" >}}) · [Capability §Success Criteria & KPIs]({{< ref "_index.md" >}}) · [UX: stand-up-the-platform §Goal]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** A total loss of the platform must not be a permanent loss. The operator must be able to rebuild the platform from nothing back to a state ready to host tenants, working only from the definitions repo and root-level access to the underlying infrastructure. Bespoke per-rebuild snowflake configuration is not acceptable as a substitute.

**Why this is a requirement, not a TR or decision:** The capability frames reproducibility as a primary outcome and the *Reproducibility* KPI ties the property to a measurable, auditable rebuild. The 1-hour budget is the TR translation of "fast enough"; the BR underneath is the existence of a clean, definitions-only rebuild path at all.

### BR-03: The platform must remain independent of any single hosting vendor
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md" >}}) · [Capability §Business Rules & Constraints]({{< ref "_index.md" >}})

**Requirement:** The operator must not be locked into any single vendor's product roadmap, pricing, or terms for the things their capabilities depend on. The platform may use vendor components, but the operator must retain control over configuration, data, and the ability to leave any one vendor.

**Why this is a requirement, not a TR or decision:** The capability lists vendor independence as an explicit outcome and clarifies the meaning of "self-hosted" as operator-controlled end-to-end rather than fully on operator-owned hardware. It is an autonomy demand on the system, not a chosen technology stack.

### BR-04: Only the operator may administer the platform
**Source:** [Capability §Stakeholders]({{< ref "_index.md" >}}) · [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: tenant-facing-observability §Constraints Inherited]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** Administrative access to the platform must be confined to the operator. There are no co-operators, no delegated administration, no self-service onboarding by tenants, and no day-to-day operator role for the designated successor. Capability owners and end users must not be granted any administrative surface on the platform.

**Why this is a requirement, not a TR or decision:** The capability states this as a hard rule and every UX inherits it by name. It is a governance boundary on the system, not a technical access-control choice.

### BR-05: A designated successor must be able to take over the platform if the primary operator becomes unavailable
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}})

**Requirement:** The platform must support a designated successor operator who holds sealed/escrowed credentials and a runbook sufficient to keep the platform running if the primary operator becomes unavailable. Successor takeover is a discrete event triggered by unavailability — it must not be a routine shared-administration mode.

**Why this is a requirement, not a TR or decision:** The capability's *Operator succession* rule explicitly demands this as a continuity property of the platform. The mechanism (password manager handoff, physical envelope, etc.) is a decision; the requirement underneath is that succession is supported at all.

### BR-06: Each tenant's users must be able to retrieve their own content while the platform is healthy, without operator involvement
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: move-off-the-platform-after-eviction §Constraints Inherited]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** While the platform is healthy, the platform must offer on-demand exportable archives of each tenant's data such that the tenant's users can pull their content without operator assistance. This is the user-data complement to operator succession: even if no successor takes over, previously-pulled exports survive.

**Why this is a requirement, not a TR or decision:** The capability's *Operator succession* rule names on-demand exportable archives as a required mechanism, and the eviction UX consumes the same mechanism. It is a user-data continuity demand, framed in user-outcome terms — what is *exported* and *how* are translations.

### BR-07: A tenant must be hostable only after the operator has authorized it
**Source:** [Capability §Triggers & Inputs]({{< ref "_index.md" >}}) · [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}})

**Requirement:** No capability may begin running on the platform without an explicit authorization decision by the operator. There must be no self-onboarding path. Acceptance of the platform's contract is implicit in the tech-design submission the operator reviews and approves.

**Why this is a requirement, not a TR or decision:** The capability lists this as a precondition and an explicit business rule; the host-a-capability UX is structured around it. It is a control demand on the system, not a workflow tooling choice.

### BR-08: Tenants must declare their resource needs and accept the platform's contract up front
**Source:** [Capability §Triggers & Inputs]({{< ref "_index.md" >}}) · [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}})

**Requirement:** To be hosted, a tenant must arrive in the form the platform accepts, with its resource needs (compute, storage, network reachability, availability expectations) declared, and must accept the platform's contract — including the platform's availability characteristics. Tenants needing stronger guarantees than the platform offers must host elsewhere.

**Why this is a requirement, not a TR or decision:** The capability codifies this as a business rule; multiple UXs inherit it. The form-of-packaging and what counts as a "declaration" are translations; the BR is that the contract exists and is accepted before hosting.

### BR-09: The platform contract must be evergreen — changes are communicated ahead of time and tenants are migrated, not surprised
**Source:** [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: platform-contract-change-rollout §Goal]({{< ref "user-experiences/platform-contract-change-rollout.md" >}}) · [UX: platform-contract-change-rollout §Constraints Inherited]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** When the platform changes a term of its contract — retiring an offering, changing a packaging form, altering availability characteristics — the operator must communicate the change to every affected tenant ahead of time, provide a migration path where applicable, and migrate tenants onto the new contract before retiring the old one. Tenants must never have a contract change sprung on them mid-lifecycle.

**Why this is a requirement, not a TR or decision:** The host-a-capability UX names the contract as evergreen; the platform-contract-change-rollout UX is the operationalization of that promise. The umbrella-issue mechanic, the cadence of status updates, and the deadline-picking algorithm are translations; the BR is the no-surprise migration property.

### BR-10: The platform may decline or evict tenants whose accommodation would break the platform's reproducibility or maintenance budget
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: operator-initiated-tenant-update §Constraints Inherited]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}}) · [UX: platform-contract-change-rollout §Constraints Inherited]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** The platform must be allowed to decline a new tenant or evict an existing tenant when continued accommodation would either push routine operation sustainably above twice the operator-maintenance-budget KPI or break the reproducibility KPI (e.g. by requiring snowflake configuration that cannot be captured as definitions). Either condition alone is sufficient grounds.

**Why this is a requirement, not a TR or decision:** The capability codifies the eviction threshold as a business rule that ties directly to its KPIs; multiple UXs reference it. It is a business-level boundary on what the platform must absorb, not a technical limit.

### BR-11: When a tenant has fallen behind what the platform supports, the default response is to bring the tenant current rather than evict
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: operator-initiated-tenant-update §Goal]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}}) · [UX: operator-initiated-tenant-update §Constraints Inherited]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}})

**Requirement:** When the divergence between a tenant and the platform is merely that the tenant's components have aged out of platform support, the operator must work with the capability owner to bring the tenant current. Eviction in this case is reserved for the missed-final-date branch, not invoked at the first sign of fall-behind.

**Why this is a requirement, not a TR or decision:** The capability explicitly carves fall-behind out of the eviction default, and the operator-initiated-tenant-update UX is the operationalization of that carve-out. It is a business commitment about the operator/tenant relationship.

### BR-12: When a tenant needs something the platform does not yet provide, the default response is to evolve the platform
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: migrate-existing-data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md" >}})

**Requirement:** When a tenant capability needs an offering the platform does not yet provide, the default response must be to consider expanding the platform's definitions to support that need rather than refusing the tenant — bounded by the reproducibility and maintenance-budget KPIs. The platform is not obligated to grow without bound, but it must not reflexively push needs back onto tenants.

**Why this is a requirement, not a TR or decision:** The capability states this as a business rule; the host-a-capability "new offering needed" branch and the existence of the migration-process offering both operationalize it. It is a posture commitment about how the platform evolves.

### BR-13: The platform must provide compute, persistent storage, and network reachability to each tenant
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md" >}})

**Requirement:** For each hosted tenant, the platform must provide a place for the tenant's application to run, durable storage for the tenant's data, and both internal (between tenants) and external (reachable by the tenant's end users) network reachability.

**Why this is a requirement, not a TR or decision:** The capability lists these as direct outputs the platform delivers to every tenant. They are business-level commitments about what hosting *means* on this platform; the form (containers, VMs, block vs. object storage) is a translation.

### BR-14: The platform must offer an identity and authentication service for tenants whose end users need to authenticate
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md" >}}) · [Capability §Triggers & Inputs]({{< ref "_index.md" >}})

**Requirement:** The platform must provide an identity and authentication service for tenant end users. Tenants may opt to bring their own identity instead, and that decision must be declared as part of the tenant's contract acceptance.

**Why this is a requirement, not a TR or decision:** The capability lists identity/authentication as a direct output and as an input that must be declared at onboarding. It is a stated platform offering, not a chosen technology.

### BR-15: The platform-provided identity service must be capable of honoring "lost credentials cannot be recovered"
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: host-a-capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md" >}})

**Requirement:** Any identity implementation the platform offers to tenants must be capable of honoring a Signal-style "lost credentials cannot be recovered" property, because at least one tenant capability requires it. An identity option that cannot honor this property is not eligible to be the platform-provided identity service.

**Why this is a requirement, not a TR or decision:** The capability codifies this as a business rule and points to the specific tenant that demands it. It is a property the identity offering must hold, framed in user-outcome terms — *which* identity software is a decision.

### BR-16: The platform must back up tenant data and provide disaster recovery to a standard the platform defines
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md" >}})

**Requirement:** For each hosted tenant, the platform must back up the tenant's data and provide a disaster-recovery path for it. The platform is responsible for defining and meeting that standard; the standard itself is part of the platform's contract.

**Why this is a requirement, not a TR or decision:** The capability lists backup and DR as a direct output and as an investment that accrues across all tenants. The specific RPO/RTO and storage choices are translations; the BR is that the platform commits to the function.

### BR-17: The platform must provide observability sufficient for the operator to tell whether each tenant is up and healthy
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md" >}}) · [UX: tenant-facing-observability §Constraints Inherited]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** The platform must produce observability for each tenant that lets the operator determine whether the tenant is up and healthy without the tenant having to instrument that itself.

**Why this is a requirement, not a TR or decision:** The capability lists observability as a direct output and frames it as the operator's view; the tenant-facing-observability UX inherits that framing. The signal bundle and tooling are translations; the BR is that the operator-side view exists.

### BR-18: Capability owners must have tenant-scoped observability of their own capability's health, including pull access and push alerts
**Source:** [UX: tenant-facing-observability §Goal]({{< ref "user-experiences/tenant-facing-observability.md" >}}) · [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** Each capability owner of a live tenant must be able to (a) pull a current view of their own tenant's health on demand, scoped strictly to their tenant, and (b) receive push alerts when a signal they care about crosses a threshold they have set. They must learn about unhealth without depending on their end users to report it.

**Why this is a requirement, not a TR or decision:** The tenant-facing-observability UX makes both arrival modes load-bearing for the capability owner's experience. The view software, alert delivery channel, and signal bundle composition are translations; the BR is the dual-mode capability-owner-facing observability surface.

### BR-19: A capability owner must not be able to see another tenant's observability data
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md" >}}) · [UX: tenant-facing-observability §Constraints Inherited]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** When a capability owner accesses the platform's observability surface, they must be confined to their own tenant's data for the entire session. Cross-tenant visibility must remain exclusive to the operator. There must be no mode-switch that lets a capability owner widen their scope.

**Why this is a requirement, not a TR or decision:** The operator-only-operation rule and the tenant-facing-observability UX together force this isolation as a tenant-trust property. It is a confidentiality demand framed in business terms — the multi-tenancy model — not an authentication scheme.

### BR-20: Capability owners must be able to self-serve their own alert thresholds
**Source:** [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** A capability owner must be able to adjust the thresholds that fire alerts for their own tenant without operator involvement. This is the one self-service surface the platform exposes to capability owners; everything else still goes through the operator.

**Why this is a requirement, not a TR or decision:** The tenant-facing-observability UX explicitly justifies this as the one carve-out from the otherwise issue-driven engagement model, and frames it as a business-level need (capability owners need to iterate on noise without operator round-trips). The UI for setting thresholds is a translation.

### BR-21: When email alerting is degraded for a tenant, the platform must surface that degradation in the tenant view
**Source:** [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md" >}}) · [UX: tenant-facing-observability §Edge Cases]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** If the platform knows alert delivery is failing for a tenant, the tenant-scoped observability view must indicate that alerting is degraded, so capability owners do not mistakenly treat email silence as evidence of health. The pull view is the source of truth; email alerts are a best-effort acceleration path.

**Why this is a requirement, not a TR or decision:** The tenant-facing-observability UX explicitly establishes the trust contract between pull view and push alerts. It is a user-trust demand on the system; how degradation is detected and rendered is a translation.

### BR-22: Capability owners must engage the platform exclusively through GitHub issues, with distinct issue types per workflow
**Source:** [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: migrate-existing-data §Journey]({{< ref "user-experiences/migrate-existing-data.md" >}}) · [UX: operator-initiated-tenant-update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}}) · [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}}) · [UX: move-off-the-platform-after-eviction §Entry Point]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** GitHub issues against the infra repo must be the only channel by which capability owners and the operator coordinate platform work — onboarding, modifications, migrations, contract changes, platform-update requests, and eviction. The issue type must be distinct per workflow (`onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, eviction issue) so that the operator's review scope and the journey shape are unambiguous to all parties.

**Why this is a requirement, not a TR or decision:** Every UX is structured around this engagement model and the distinctness of issue type is load-bearing in each. It is a business-level contract about how the operator and capability owners interact; the choice of GitHub Issues *as the platform* would be an ADR-level decision were the platform itself ever changed, but the requirement that distinct, persistent issue threads carry the work is what the UXs force.

### BR-23: Onboarding a tenant must use only the platform's existing definitions — no per-tenant snowflake configuration
**Source:** [UX: host-a-capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** Provisioning a new tenant during onboarding must run against the platform's definitions; the operator must not hand-roll bespoke configuration to make a particular tenant fit. If onboarding requires bespoke manual config that cannot be captured as definitions, that is a platform-level reproducibility failure to be fixed at the platform level rather than absorbed silently.

**Why this is a requirement, not a TR or decision:** The host-a-capability UX traces this directly to the *Reproducibility* KPI inherited from the capability, and the stand-up-the-platform UX reinforces that "delete everything and start over" must always work. It is a business-level commitment about what hosting *costs* the platform's reproducibility property.

### BR-24: The platform must offer a one-shot migration-process runner for capability-owner-supplied data migrations
**Source:** [UX: migrate-existing-data §Goal]({{< ref "user-experiences/migrate-existing-data.md" >}}) · [UX: migrate-existing-data §Journey]({{< ref "user-experiences/migrate-existing-data.md" >}}) · [UX: migrate-existing-data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md" >}})

**Requirement:** The platform must offer a generic facility for running a one-shot migration job that a capability owner has packaged and handed off — reading from a prior host, writing into the capability owner's already-provisioned tenant. The migration job's lifecycle must be: provision when approved, run, observable while running, torn down on completion. The platform runs what it is given; logic correctness belongs to the capability owner.

**Why this is a requirement, not a TR or decision:** The migrate-existing-data UX is the operationalization of "the capability evolves with its tenants" for the migration use case. The packaging form, scheduler, and runner technology are translations; the BR is that this offering exists and follows the one-shot lifecycle.

### BR-25: The platform must offer secret management to which capability owners can register credentials their workloads need
**Source:** [UX: migrate-existing-data §Journey]({{< ref "user-experiences/migrate-existing-data.md" >}})

**Requirement:** The platform must provide a secret-management offering through which capability owners can register credentials (e.g. credentials a migration process needs to read from an old host), referenced by name from packaged artifacts so the secrets themselves never appear in issue threads or artifacts.

**Why this is a requirement, not a TR or decision:** The migrate-existing-data UX makes this an explicit precondition step that must exist before issue filing. It is a stated platform offering with a confidentiality-shaped business demand; the choice of secret-management product is a decision.

### BR-26: A migration job's peak temporary footprint must be bounded relative to the destination tenant's steady-state footprint
**Source:** [UX: migrate-existing-data §Journey]({{< ref "user-experiences/migrate-existing-data.md" >}})

**Requirement:** The platform must enforce, during migration review, that the migration job's peak temporary compute and storage footprint stays within a bounded multiple of the destination tenant's steady-state footprint. Migrations whose declared spike exceeds that bound must be rejected as written and split, reduced, or preceded by a tenant resize.

**Why this is a requirement, not a TR or decision:** The UX states the bound (currently 2x) as the business rule the operator enforces during migration review. The exact multiple is the TR translation; the BR is the existence of a declared, enforced cap on migration footprint relative to the tenant.

### BR-27: The platform must provide on-demand data export for every tenant, regardless of data shape
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: move-off-the-platform-after-eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}}) · [UX: move-off-the-platform-after-eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** Export tooling must be a core platform feature available for every kind of data the platform hosts, present from the moment a tenant is live (not assembled at eviction). The capability owner must be able to invoke it at any time the platform is healthy and receive an archive they can download immediately.

**Why this is a requirement, not a TR or decision:** The capability ties this to *Operator succession*, and the eviction UX makes it a hard property: an export-tooling gap is a platform bug, not an excuse for a slow exit. The format and packaging are translations; the BR is universal export availability.

### BR-28: Each export must include checksum/hash and total size produced by the platform
**Source:** [UX: move-off-the-platform-after-eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** Whenever the platform produces an export archive for a capability owner, it must publish alongside it a checksum/hash and total size in bytes. This is the ceiling of integrity verification the platform commits to; the capability owner is responsible for semantic validation beyond that.

**Why this is a requirement, not a TR or decision:** The eviction UX makes this the contract that bounds platform vs. capability-owner responsibility for "is this export complete." The hash algorithm choice is a decision; the BR is the existence of platform-verified bytes-level integrity.

### BR-29: After eviction, the platform must keep tenant-accessible data available for a fixed retention window with no slip
**Source:** [UX: move-off-the-platform-after-eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}}) · [UX: move-off-the-platform-after-eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** From the eviction date, the platform must hold the evicted tenant's data in an export-only, read-only state for a fixed retention window (currently 30 days), during which the export tool continues to work. After that window, no tenant-accessible copy of the data must remain. The window must not be extended for slow extracts or for the capability owner asking for more time.

**Why this is a requirement, not a TR or decision:** The eviction UX establishes both the existence of the window and its hard-wall property. The exact 30-day number is the TR translation; the BR is the existence of a fixed, no-slip retention window with a defined end-state.

### BR-30: When a platform-side defect prevents a clean export, the eviction retention countdown must pause until a clean export is achievable
**Source:** [UX: move-off-the-platform-after-eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** If a failure to produce a complete, valid export is rooted in the platform's export tooling or its data hosting (not in the capability owner's own validation), the platform must pause that tenant's removal-of-tenant-accessible-data countdown until the platform-side defect is resolved and a clean export has been produced. Capability-owner-rooted failures must not pause the countdown.

**Why this is a requirement, not a TR or decision:** The eviction UX explicitly carves this out as the one exception to the otherwise-hard retention wall. It is a fairness commitment about how the platform's defects affect the tenant's exit rights.

### BR-31: After eviction, compute and network must be torn down on the eviction date, and tenant data must transition to read-only
**Source:** [UX: move-off-the-platform-after-eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** On the eviction date, the platform must deprovision the evicted tenant's compute, network, and other live-serving resources, and transition the tenant's data to an export-only, read-only state in which no further writes (by anyone, including the capability owner) can occur. The eviction issue must carry the cutover confirmation.

**Why this is a requirement, not a TR or decision:** The eviction UX makes the cutover and the read-only freeze load-bearing for the rest of the journey (the data the capability owner extracts in Phase C is the *final* dataset). It is a state-transition demand framed in user-experience terms.

### BR-32: The platform must operate within the operator's weekly maintenance budget
**Source:** [Capability §Success Criteria & KPIs]({{< ref "_index.md" >}})

**Requirement:** Routine operation of the platform — across all hosted tenants and all routine UXs — must stay within a defined weekly budget of the operator's time (currently no more than 2 hours per week). Sustained over-budget operation is a signal that the platform must be simplified rather than grown, and is the operational basis for the eviction threshold rule.

**Why this is a requirement, not a TR or decision:** The capability defines this as a KPI and explicitly ties the eviction threshold to it. The number is the TR/KPI translation; the BR is that there is a budget at all and that the platform must live inside it.

### BR-33: The platform's cost must remain proportional to the convenience and resiliency it delivers
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [Capability §Success Criteria & KPIs]({{< ref "_index.md" >}})

**Requirement:** Total operating cost must remain within what the operator considers acceptable given the convenience and resiliency the platform delivers. There is no fixed dollar target, but cost must be minimized where doing so does not cost convenience or resiliency, and added cost is acceptable only when it buys meaningful convenience or resiliency.

**Why this is a requirement, not a TR or decision:** The capability codifies the cost-secondary-to-convenience-and-resiliency tiebreaker and the proportionality KPI. It is a value-judgment property of the platform's spending, not a technical budget constraint.

### BR-34: End users of tenant capabilities must not have any direct interface with the platform
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [Capability §Stakeholders]({{< ref "_index.md" >}}) · [UX: move-off-the-platform-after-eviction §Constraints Inherited]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}}) · [UX: tenant-facing-observability §Constraints Inherited]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** The platform must have no notion of "end users" of itself. End users reach the tenant capability, not the platform. The platform must not communicate with end users (e.g. no platform-side eviction notice page, no platform-side status page for end users) and must not surface administrative or observability surfaces to them.

**Why this is a requirement, not a TR or decision:** The capability defines the platform's actors as not including end users, and multiple UXs reinforce that boundary. It is a scope boundary on the system, not a UX choice.

### BR-35: The platform may span public and private infrastructure
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** The platform must be permitted to comprise both operator-controlled public-cloud components and operator-owned private infrastructure. "Self-hosted" must mean operator-controlled end-to-end, including configuration, data, and the ability to leave a vendor — not exclusively operator-owned hardware.

**Why this is a requirement, not a TR or decision:** The capability codifies this rule and the standup UX explicitly crosses cloud and home-lab boundaries during foundation provisioning. It is a scope demand on the platform's allowed shape.

### BR-36: Connectivity between public and private parts of the platform must itself be part of the foundation
**Source:** [UX: stand-up-the-platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** Where the platform spans cloud and home-lab, the connectivity between them must be provisioned as part of the platform's foundation during standup, not bolted on afterward. A reproducible rebuild must produce a working cross-environment platform without manual cross-environment wiring.

**Why this is a requirement, not a TR or decision:** The standup UX makes this load-bearing for Phase 1 success. The transport (Wireguard, VPN, peering) is a decision; the BR is that connectivity is foundational.

### BR-37: Standup must verify readiness end-to-end with a purpose-built canary tenant before declaring the platform ready
**Source:** [UX: stand-up-the-platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Edge Cases]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** A platform standup must not be declared ready on the basis of infrastructure self-checks alone. A purpose-built canary tenant maintained alongside the platform definitions must be deployed, exercised end-to-end (run, reachability, storage read/write, identity authentication, backup pickup, observability pickup), and torn down. Until that canary is green, the platform must not be marked ready, regardless of time pressure.

**Why this is a requirement, not a TR or decision:** The standup UX makes the canary load-bearing for the readiness signal and explicitly forbids bending this rule under time pressure. The canary's contents and tooling are decisions; the BR is the existence of a green-canary readiness gate.

### BR-38: Standup must run in phases with operator validation checkpoints between them
**Source:** [UX: stand-up-the-platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** Platform standup must proceed in distinct automated phases (foundations, core services, cross-cutting services, canary), pausing between phases for the operator to validate that the previous phase really succeeded. The operator must signal `continue` before the next phase runs.

**Why this is a requirement, not a TR or decision:** The standup UX makes the phased validation pattern load-bearing for the "confidence beats speed" tiebreaker the operator brings to the journey. The number and naming of phases are translations; the BR is that automation is interleaved with operator validation, not all-or-nothing.

### BR-39: A failed phase during standup must be recoverable only via tear-down-and-restart
**Source:** [UX: stand-up-the-platform §Edge Cases]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** When a standup phase fails validation, the operator must tear down everything provisioned so far and restart the rebuild from the top after fixing the underlying definition. Partial state must not be trusted as a starting point for resumption. Each phase must therefore support a clean teardown of its partial output.

**Why this is a requirement, not a TR or decision:** The standup UX explicitly justifies this in terms of the parent capability's reproducibility-beats-effort tiebreaker. The teardown mechanism is a decision; the BR is the no-resume-from-partial-state property.

### BR-40: Drift of platform state from definitions must be detected and resolved before any rebuild begins
**Source:** [UX: stand-up-the-platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** Whenever prior platform state exists, a preflight drift check against the live platform or last known-good environment must pass before any rebuild begins. Drift must not be discovered partway through a rebuild. The continuous machinery that prevents drift between rebuilds — tracked changes and immutability — must be enforced by every UX that can introduce platform state.

**Why this is a requirement, not a TR or decision:** The standup UX refuses to start without this check and ties it to the integrity of every other UX. The change-tracking mechanism is a decision; the BR is the drift-free-before-rebuild guarantee and the cross-UX immutability discipline.

### BR-41: Platform-rebuild capability must be drilled regularly to prove the reproducibility KPI is real
**Source:** [UX: stand-up-the-platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** The operator must rebuild the platform in parallel on scratch infrastructure after every significant platform change and at least quarterly, using the same flow as a real rebuild. The reproducibility KPI must not be a hope inferred from version control; it must be a property re-proved by drill on a bounded cadence.

**Why this is a requirement, not a TR or decision:** The standup UX makes this part of how the *Reproducibility* KPI is honestly evaluated. The frequency words ("every significant change", "quarterly") are scoping language; the BR is the drill-as-proof discipline.

### BR-42: A successor operator who has taken over must be able to run the standup flow identically
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}}) · [UX: stand-up-the-platform §Persona]({{< ref "user-experiences/stand-up-the-platform.md" >}}) · [UX: stand-up-the-platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** Once a successor has taken over and gained the operator's context, the rebuild flow must work identically for them as for the primary operator. The standup UX must not depend on operator-specific state outside what the sealed credentials and the definitions repo provide.

**Why this is a requirement, not a TR or decision:** The capability's *Operator succession* rule and the standup UX persona section make this load-bearing for continuity. The seal mechanism is a decision; the BR is the convergence of successor and primary on the same standup flow.

### BR-43: Operator skill development must not influence buy-vs-build trade-offs at the capability level
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md" >}})

**Requirement:** The operator may personally learn from building and running the platform, but skill development must not be a stated outcome of this capability and must not influence buy-vs-build decisions. Those decisions must be judged on convenience, resiliency, and cost only.

**Why this is a requirement, not a TR or decision:** The capability codifies this as a hard rule precisely to keep the platform from being shaped by curiosity. It is a decision-discipline boundary, not a technical or implementation constraint.

### BR-44: Modifications to a live tenant must not require re-acceptance of the platform contract
**Source:** [UX: host-a-capability §Journey]({{< ref "user-experiences/host-a-capability.md" >}})

**Requirement:** When a capability owner files a `modify my capability` issue, the operator's review must cover only the delta requested; the capability owner must not be required to re-accept the platform's contract on each modification. Any change to the platform's contract itself must come through the platform-contract-change rollout, not through modify reviews.

**Why this is a requirement, not a TR or decision:** The host-a-capability UX explicitly establishes the evergreen-contract promise on the change-later loop, and the contract-change UX inherits the same boundary. It is a relationship commitment, not a workflow tooling choice.

### BR-45: A platform-contract-change rollout must run both old and new contract forms concurrently during the rollout window, except for full removals
**Source:** [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}}) · [UX: platform-contract-change-rollout §Edge Cases]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** When the platform changes a contract term and a replacement offering exists, the platform must serve both the old and the new form concurrently during the rollout window, so tenants can migrate at their own pace within the deadline. A full offering removal (no replacement) is the carve-out: it is all-or-nothing at the deadline.

**Why this is a requirement, not a TR or decision:** The contract-change UX makes concurrent old/new the operationalization of the no-surprises promise. The mechanics of running both are a decision; the BR is the side-by-side property and its narrowly-defined exception.

### BR-46: A platform-contract-change deadline must give every affected tenant at least two full status-update cycles before cutoff
**Source:** [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** Because contract-change rollouts have no externally-imposed deadline to inherit, the operator must pick a deadline that gives every affected tenant at least two full status-update cycles before cutoff: one to acknowledge and start, one to finish or surface blockers in time to respond.

**Why this is a requirement, not a TR or decision:** The contract-change UX makes this the floor on operator-chosen deadlines. The exact cadence is a decision; the BR is the minimum tenant-reaction-window built into any contract-change deadline.

### BR-47: A contract-change rollout must publish status updates on a regular operator-chosen cadence in the umbrella thread
**Source:** [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** Throughout a contract-change rollout, the operator must post status updates on a regular schedule (sized to the rollout's overall timeline) in the umbrella issue thread. The current snapshot must live in the umbrella issue body so a reader landing cold sees the latest state; each scheduled update must also be posted as a comment so the rollout history remains visible. Updates must include how many tenants remain on the old form, how many have migrated, which `modify` issues are open, and how much time remains.

**Why this is a requirement, not a TR or decision:** The contract-change UX establishes status updates as how every party sees rollout progress without chasing it. The cadence is a decision; the BR is the standing visibility commitment and its required content.

### BR-48: Acknowledgment of a contract-change umbrella issue must be required from every affected capability owner, with non-acknowledgment treated as non-engagement
**Source:** [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}}) · [UX: platform-contract-change-rollout §Edge Cases]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** Each capability owner tagged on a contract-change umbrella issue must acknowledge in-thread. Silence in a multi-tenant umbrella is ambiguous and must not be treated as tacit consent; if no acknowledgment arrives by the deadline, the missing acknowledgment must be treated as non-engagement and that tenant must enter the laggard branch (separate eviction issue) like any other unmigrated tenant.

**Why this is a requirement, not a TR or decision:** The contract-change UX makes explicit acknowledgment load-bearing for distinguishing engagement from drift. It is a relationship commitment about how silence is interpreted in multi-tenant rollouts.

### BR-49: A contract-change deadline must not be negotiable per-tenant; only global extensions are allowed
**Source:** [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}}) · [UX: platform-contract-change-rollout §Edge Cases]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** During a contract-change rollout, the deadline must apply uniformly to all affected tenants. Per-tenant extensions must not be granted. The operator may push the deadline globally if the migration guideline turns out to be insufficient mid-rollout, and the new global deadline must then be announced in the umbrella thread under the same hard-deadline rule.

**Why this is a requirement, not a TR or decision:** The contract-change UX names this as the property that keeps a deadline a deadline. It is a fairness-and-coordination commitment, not a tooling constraint.

### BR-50: A platform-update-required issue must inherit its deadline from the external forcing event and must use a distinct issue type
**Source:** [UX: operator-initiated-tenant-update §Entry Point]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}}) · [UX: operator-initiated-tenant-update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}})

**Requirement:** When a platform-level dependency event (vendor sunset, CVE, EOL) forces a tenant update, the operator must file a `platform update required` issue per affected tenant whose deadline is inherited from the external event and whose external reason is recorded on the issue. The distinct issue type must signal to the capability owner that this is a required update, not an externally-driven optional cleanup or routine modify.

**Why this is a requirement, not a TR or decision:** The operator-initiated-tenant-update UX makes the inherited deadline and the distinct type load-bearing for the journey's interpretation. The CVE/EOL feed source is operational detail; the BR is the inherited-deadline contract and the typed-signal property.

### BR-51: Slack negotiation on a platform-update-required deadline must be bounded by the safe slack the external pressure allows
**Source:** [UX: operator-initiated-tenant-update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}})

**Requirement:** When a capability owner cannot ship within the inherited deadline, the operator and capability owner must negotiate an extended delivery date only to the extent the external pressure leaves safe slack. If no safe slack exists (e.g. an actively-exploited CVE), the original deadline must stand. The agreed date must be recorded on the issue.

**Why this is a requirement, not a TR or decision:** The operator-initiated-tenant-update UX establishes safe-slack-bounded extension as the way "we work with you, we don't evict for fall-behind" is honored without breaking the underlying external constraint. The negotiation form is a decision; the BR is the bound on extension.

### BR-52: A missed operative delivery date on a platform-update-required issue must trigger eviction via a separate, linked eviction issue
**Source:** [UX: operator-initiated-tenant-update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}}) · [UX: platform-contract-change-rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** When a tenant misses the operative delivery date — either an inherited deadline or an agreed extension — the operator must open a separate eviction issue linking back to the originating issue, and close the original as superseded by eviction. The same shape applies to laggards in a contract-change rollout. Eviction must always live in its own issue, not as a state on the originating issue.

**Why this is a requirement, not a TR or decision:** Both the operator-initiated-tenant-update and platform-contract-change-rollout UXs prescribe the same separate-and-linked structure for eviction. It is a record-keeping and journey-boundary commitment about how state transitions between updates and eviction.

### BR-53: When a platform-update request reveals a need the platform does not yet offer, the update must be paused and resumed via the new-offering branch
**Source:** [UX: operator-initiated-tenant-update §Edge Cases]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}})

**Requirement:** If shipping a platform-update-required change requires a capability the platform does not yet offer, the platform-update issue must remain open while the new-offering branch (per host-a-capability) is exercised, then the update must resume at the modify inner-loop step. The update must not be force-closed because the platform was not yet ready.

**Why this is a requirement, not a TR or decision:** The operator-initiated-tenant-update UX explicitly handles this as a join with host-a-capability's new-offering branch, in service of "the capability evolves with its tenants." It is a workflow-continuity commitment, not a tooling choice.

### BR-54: An evicted capability owner must walk away with no obligations and no tenant-accessible copy of their data left on the platform after the retention window
**Source:** [UX: move-off-the-platform-after-eviction §Goal]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}}) · [UX: move-off-the-platform-after-eviction §Success]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** By the time the eviction journey ends, the capability owner must have a clean exit: validated data in their hands, a clear paper trail on the eviction issue, and nothing left to chase down on the platform. After the retention window, the platform must offer no tenant-accessible copy of their data, regardless of whether the capability owner ever closed the loop. The relationship must remain amicable enough to support a future `host-a-capability` re-onboarding if the divergence later resolves.

**Why this is a requirement, not a TR or decision:** The eviction UX makes the clean-exit property load-bearing for the journey's success and explicitly preserves the right to come back. It is a user-outcome demand; the deeper backup-tier policy is open (see Open Questions).

## Open Questions

Things the user volunteered as TRs or decisions during extraction (parked for the next stage), or constraints the capability/UX docs don't yet make explicit.

- **Explicit anchors on capability and UX section headings.** The capability doc and UX docs currently rely on Hugo's slugify-from-heading-text for section anchors. Per skill guidance, section deep-links require explicit `{#anchor-id}` annotations on the target headings. Source links in this document point at section names by description but anchor only at the page level until explicit anchors are added to the source headings. This should be resolved in the source docs before downstream skills cite specific BR sections.
- **Numeric thresholds belong with TRs/KPIs, not here.** Several BRs (BR-26 migration footprint cap, BR-29 retention window, BR-32 maintenance budget, BR-46 status-update cycles) deliberately leave the numbers (`2x`, `30 days`, `2 hr/week`, `at least two cycles`) as KPI/TR-stage translations. The numeric values currently live in the capability doc and UX docs; `define-technical-requirements` should pick them up from there.
- **Deeper backup-tier policy after the 30-day retention window.** The eviction UX itself parks this as an open question: retention duration, deletion behavior, and operator-access/privacy constraints for any deeper backup-tier copies after tenant-accessible data is removed are still TBD. This is a BR-shaped gap that the capability or eviction UX needs to fill before BR-29/BR-54 can be considered complete.
- **No formal tenant-facing pending-update view exists today.** The operator-initiated-tenant-update UX notes that capability owners receive their first official signal of a forced update only when the per-tenant issue is filed. If an earlier deprecation/pending-update signal is added later, it would extend BR-18 (tenant-facing observability) rather than this UX. Parked for tenant-facing-observability evolution.
- **Backup standard.** BR-16 says the platform must back up tenant data "to a standard the platform defines" — but the standard itself is not defined in the capability doc. This is a gap for a future capability-doc revision; BR-16's TR translation will need that standard to be measurable.
- **Specific signal bundle for tenant-facing observability.** BR-18 names the platform-standard health bundle (availability, latency, error rate, resource saturation, restart/deployment events) as listed in the tenant-facing-observability UX, but the bundle's exact contents and any platform-default thresholds are translation territory for `define-technical-requirements`.
- **Migration-process offering's concurrency model.** The migrate-existing-data UX promises concurrent migrations across tenants but does not bound capacity. Capacity sizing and queueing belong with TRs/ADRs.
- **Drift detection mechanism.** BR-40 demands that drift be detected and resolved before rebuild and that tracked-changes/immutability be enforced cross-UX, but the mechanism (what counts as the "last known-good environment", how drift is computed, where the policy is enforced) is open and is a TR/ADR concern.
