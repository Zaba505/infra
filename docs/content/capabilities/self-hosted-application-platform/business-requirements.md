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

### BR-01: Provide a single default hosting target for the operator's capabilities {#br-01}
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md#purpose" >}})

**Requirement:** The platform must be the default place where the operator's capabilities run, so that "where does this run?" is a solved question rather than re-litigated per capability. Any capability the operator defines must be eligible to run here unless it is explicitly exempted.

**Why this is a requirement, not a TR or decision:** This is the first stated business outcome of the capability and the rationale for the capability existing at all. It does not name a technology — it sets the demand that the platform exist and absorb every capability by default.

### BR-02: Platform must be reproducible from its definitions {#br-02}
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md#purpose" >}}) · [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** The platform itself must be rebuildable from its definitions, with no manual snowflake configuration, so that a total loss does not mean a permanent loss of the platform. Any state the platform depends on must be expressible as part of the definitions.

**Why this is a requirement, not a TR or decision:** Reproducibility is one of the capability's named outcomes and is the operative test of "self-hosted." It does not specify how (which IaC tool, which packaging form) — only the demand that the platform be rebuildable from authoritative inputs.

### BR-03: Operator must retain end-to-end control over platform components {#br-03}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}})

**Requirement:** Wherever the platform uses third-party components (vendor services, public-cloud offerings), the operator must retain control of configuration, data, and the ability to leave. Vendor lock-in that prevents departure is unacceptable.

**Why this is a requirement, not a TR or decision:** "Self-hosted" is defined in the capability as operator-controlled end-to-end, not as forbidding all vendors. The BR sets the demand (retain control + ability to leave), not the implementation choice.

### BR-04: Platform-level investments must accrue to all tenants {#br-04}
**Source:** [Capability §Purpose & Business Outcome]({{< ref "_index.md#purpose" >}})

**Requirement:** Improvements made at the platform level (resiliency, observability, backup, security) must benefit every tenant capability rather than be re-solved per tenant.

**Why this is a requirement, not a TR or decision:** This is one of the four stated outcomes the capability promises. It demands a property of the platform's offerings (shared, not per-tenant), without choosing how that sharing is achieved.

### BR-05: Only the operator may administer the platform {#br-05}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}})

**Requirement:** There must be no co-operators, no delegated administration, and no shared day-to-day administration of the platform. The operator is the sole administrator.

**Why this is a requirement, not a TR or decision:** The capability's "Operator-only operation" rule states this in absolute terms. It is a forced constraint — every UX is shaped around the operator being the only one with administrative reach.

### BR-06: End users of tenants must have no direct access to the platform {#br-06}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Move Off the Platform After Eviction §Constraints Inherited]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#constraints-inherited" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** End users of tenant capabilities reach the tenant, not the platform. The platform must have no notion of "end users" of itself, no UI for them, and no communication channel to them — including during eviction. Capability owners, not the platform, are responsible for notifying their own end users of any tenant lifecycle change (such as an impending shutdown).

**Why this is a requirement, not a TR or decision:** The capability explicitly excludes direct end-user access; the eviction UX reinforces it (the platform never tells end users "this tenant has been retired") and assigns the notification responsibility to the capability owner. The BR forbids a class of behavior, not a specific implementation.

### BR-07: A designated successor must be able to take over operation if the primary operator becomes unavailable {#br-07}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Stand Up the Platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md#constraints-inherited" >}})

**Requirement:** The platform must support a designated successor operator who holds sealed/escrowed credentials and a runbook sufficient to keep the platform running if the primary operator becomes unavailable. Successor credentials are not used for routine operation.

**Why this is a requirement, not a TR or decision:** Operator succession is one of the capability's business rules and the standup UX assumes a successor can run the rebuild flow. The BR demands that takeover be possible; how the seal works is a downstream concern.

### BR-08: Platform must produce on-demand exportable archives of tenant data while healthy {#br-08}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** While the platform is up, every tenant's users (via the capability owner) must be able to retrieve their content as a portable archive without operator involvement. Export availability is conditional only on the platform being healthy.

**Why this is a requirement, not a TR or decision:** This is half of the capability's "Operator succession" rule (the other half is the successor) and the central mechanism of the eviction UX. It states what the user must be able to obtain, not how the export is implemented. Pairs with [BR-09](#br-09), which forbids gaps in export-tooling coverage.

### BR-09: Export tooling must exist for every kind of data the platform hosts {#br-09}
**Source:** [UX: Move Off the Platform After Eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#edge-cases" >}})

**Requirement:** Export tooling must be a core platform feature, available for every data shape the platform hosts. There must be no tenant whose data shape lacks an export path at the time of eviction.

**Why this is a requirement, not a TR or decision:** The eviction UX explicitly states this cannot-happen-by-design property and treats any gap as a platform bug. The BR forbids a class of failure rather than naming a tool.

### BR-10: Exports must include platform-produced verification material {#br-10}
**Source:** [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** Each export the platform produces must be accompanied by a checksum/hash and total size in bytes, so the capability owner can verify integrity. Semantic correctness validation remains the capability owner's responsibility; the platform's verification is the ceiling of what it can offer on the user's behalf.

**Why this is a requirement, not a TR or decision:** The eviction UX makes this guarantee explicit — the platform produces the integrity envelope; the user judges semantic correctness. The BR demands the envelope, not a specific hash function.

### BR-11: Tenant data must remain retrievable for 30 days after eviction in a read-only state {#br-11}
**Source:** [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** From the eviction date forward, the platform must hold the tenant's data in an export-only, read-only state for 30 days, during which the export tool must continue to work. After 30 days, the platform must stop offering any tenant-accessible copy of that data.

**Why this is a requirement, not a TR or decision:** The eviction UX defines the 30-day window as a hard tenant-facing guarantee. It is a business commitment to the departing capability owner, not a technical translation.

### BR-12: Export-tooling defects must pause the post-eviction retention countdown {#br-12}
**Source:** [UX: Move Off the Platform After Eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#edge-cases" >}})

**Requirement:** If a failure rooted in the platform's export tooling or data hosting prevents a clean export, the operator must pause that tenant's retention-window countdown until a clean export can be produced. Failures rooted in the capability owner's own validation steps must not pause the countdown.

**Why this is a requirement, not a TR or decision:** This is the only carve-out in the eviction UX's otherwise-hard 30-day rule, and it allocates accountability — the platform absorbs slippage caused by its own defects, never the user's. That is a business commitment.

### BR-13: Tenants must declare resource needs, packaging form, identity choice, and availability expectations up front {#br-13}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Host a Capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md#constraints-inherited" >}})

**Requirement:** To be hosted, a tenant must arrive packaged in the form the platform accepts, with declared resource needs, an identity-service choice, and acceptance of the platform's current availability characteristics. The declarations are made in the tech design and reviewed before approval.

**Why this is a requirement, not a TR or decision:** The capability's "Tenants must accept the platform's contract" rule and the host-a-capability UX both make this the price of admission. The BR demands the declaration; what shape "packaging" takes is a downstream decision.

### BR-14: Tenant onboarding must require explicit operator authorization {#br-14}
**Source:** [Capability §Triggers & Inputs]({{< ref "_index.md#triggers" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** No capability may begin running on the platform without the operator explicitly authorizing it. There must be no self-service onboarding path.

**Why this is a requirement, not a TR or decision:** The capability lists this as a precondition; the UX's "approved" comment is its operationalization. It is a control demand, not a technical translation.

### BR-15: All capability-owner ↔ platform engagement must occur through a single, recorded, asynchronous issue-thread workflow {#br-15}
**Source:** [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}}) · [UX: Move Off the Platform After Eviction §Entry Point]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#entry-point" >}})

**Requirement:** Every operator/capability-owner exchange (onboarding, modification, migration, platform-driven update, contract change, eviction) must occur on a single recorded, asynchronous-by-default issue thread that both parties can read and append to. There must be no self-service portal and no other front door, and exchanges must not happen over ephemeral channels (chat, email-only, voice) where the trail is lost.

**Why this is a requirement, not a TR or decision:** The UXes demand the *properties* — single thread, recorded, asynchronous, no other front door — without naming a tracker. Choosing a specific tracker (e.g. GitHub Issues) is a downstream decision recorded in an ADR, not here.

### BR-16: Issue types must distinguish review scopes legibly {#br-16}
**Source:** [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** Distinct issue types must exist for the distinct conversations the platform has: onboarding a capability, modifying a hosted capability, migrating data, operator-initiated forced updates, platform contract changes, and eviction. The type itself signals the operator's review scope and the capability owner's expectations.

**Why this is a requirement, not a TR or decision:** Each UX names its issue type explicitly and explains why it is distinct. It is a coordination demand on the platform, not a tooling decision.

### BR-17: The platform contract must be evergreen for already-hosted tenants {#br-17}
**Source:** [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Platform-Contract-Change Rollout §Constraints Inherited]({{< ref "user-experiences/platform-contract-change-rollout.md#constraints-inherited" >}})

**Requirement:** A capability owner must not be required to re-accept the platform's contract on each modify request. Changes to the platform's contract are the platform's responsibility to communicate ahead of time and migrate tenants through; the contract is never sprung on a tenant during a modify request.

**Why this is a requirement, not a TR or decision:** The host-a-capability UX states the evergreen property; the contract-change rollout UX is its operationalization. It is a promise to the user, not a technical translation.

### BR-18: Capability owners must be able to update tenant needs after onboarding via a delta-only review {#br-18}
**Source:** [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** Once a tenant is live, its capability owner must be able to file a modify request that the operator reviews scoped to the delta only — not as a full re-evaluation of the tenant.

**Why this is a requirement, not a TR or decision:** The host-a-capability UX's change-later loop is built around this property. It is a user-experience demand on the modify path, not a tooling choice.

### BR-19: Platform must work with tenants on fall-behind cases rather than evict {#br-19}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** When a tenant's components have fallen behind what the platform supports, the default operator response must be to bring the tenant current rather than evict. Eviction in fall-behind cases must occur only as a downstream consequence of a missed operative delivery date, never as the first response.

**Why this is a requirement, not a TR or decision:** The capability's eviction rule carves out this behavior explicitly; the operator-initiated-tenant-update UX is the carve-out's operationalization. It is a user-facing commitment, not a technical translation.

### BR-20: Forced-update issues must carry the external reason and inherited deadline {#br-20}
**Source:** [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** When the operator opens a `platform update required` issue, it must name the external pressure forcing the change (vendor sunset, CVE, EOL) and the deadline inherited from that pressure. Each forcing event gets its own issue, even when the same tenant is hit by multiple events at once.

**Why this is a requirement, not a TR or decision:** The UX makes both properties explicit and motivates them — the capability owner needs to see *why* and *by when*. It is a transparency demand.

### BR-21: Extensions to inherited deadlines must be bounded by the external pressure's safe slack {#br-21}
**Source:** [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** When a capability owner cannot ship within an inherited deadline, any extension granted must be sized to the slack the external pressure actually allows — never invented by the operator independent of that pressure. If the pressure leaves no safe slack, extensions must be refused.

**Why this is a requirement, not a TR or decision:** The UX explicitly states this constraint. It is a control on operator discretion to prevent the platform from absorbing risk it cannot honestly carry.

### BR-22: A missed operative delivery date must result in a separate, linked eviction issue {#br-22}
**Source:** [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** When a capability owner misses the operative date for a forced update (the original inherited deadline or an agreed extension), the operator must open a separate eviction issue linked back to the update issue, and close the update issue as superseded. Eviction must not be re-policed inside the update flow.

**Why this is a requirement, not a TR or decision:** The UX prescribes this exact split, with the rationale that update-flow scope and eviction-flow scope must remain distinct. It is a coordination demand on the platform.

### BR-23: Operator-driven contract changes must be communicated ahead of time via a single umbrella issue {#br-23}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** When the operator chooses to retire an offering, change a packaging form, or alter availability characteristics, the change must be announced via a single umbrella issue tagging every affected capability owner, containing what is changing, what it is changing to, the deadline, the reason, and the migration guideline (where applicable).

**Why this is a requirement, not a TR or decision:** The contract-change rollout UX names this shape explicitly and motivates the umbrella over per-tenant issues. It is the operationalization of the evergreen promise.

### BR-24: Contract-change deadlines must give every affected tenant at least two status-update cycles {#br-24}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** When the operator picks a contract-change deadline, it must allow every affected tenant at least two full status-update cycles before cutoff — one to acknowledge and start, one to finish or surface blockers with time still to respond.

**Why this is a requirement, not a TR or decision:** The UX states this minimum explicitly. It is a fairness commitment that bounds operator discretion.

### BR-25: Contract-change deadlines must not be negotiable per-tenant {#br-25}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** Capability owners must not be able to negotiate per-tenant slips of a contract-change deadline. The deadline applies uniformly; only a global extension (covering every affected tenant) is available, and only when the migration guideline itself proves insufficient.

**Why this is a requirement, not a TR or decision:** The UX makes this rule absolute. It enforces that the deadline remains a deadline rather than degrading into a per-tenant negotiation.

### BR-26: Capability owners must explicitly acknowledge contract-change umbrella issues {#br-26}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** Each tagged capability owner on a contract-change umbrella issue must explicitly acknowledge the change in-thread. Silence in a multi-tenant thread is treated as non-engagement and feeds the same laggard branch as failing to migrate.

**Why this is a requirement, not a TR or decision:** The UX states the acknowledgment requirement and its consequence. It is an engagement contract, not a technical mechanism.

### BR-27: During contract-change rollout, old and new forms must run concurrently when a replacement exists {#br-27}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** When a contract change replaces an old offering with a new one, the platform must serve both forms concurrently throughout the rollout window. Full offering removals (no replacement) are exempt — the change is all-or-nothing at the deadline.

**Why this is a requirement, not a TR or decision:** The UX prescribes the concurrent rollout window and names the carve-out. It is a user-facing commitment that gives tenants room to migrate at their own pace.

### BR-28: Replacement offerings must be implemented and running before a contract-change umbrella issue is filed {#br-28}
**Source:** [UX: Platform-Contract-Change Rollout §Entry Point]({{< ref "user-experiences/platform-contract-change-rollout.md#entry-point" >}})

**Requirement:** Where a contract change replaces an old offering with a new one, the replacement must already be implemented and running on the platform alongside the old one before the umbrella issue is filed. Tenants must never be asked to migrate against an unbuilt replacement.

**Why this is a requirement, not a TR or decision:** The UX states this as a precondition of the journey. It is a quality-of-rollout commitment.

### BR-29: Operator must post regular status updates throughout a contract-change rollout {#br-29}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** During a contract-change rollout the operator must post status updates on a regular schedule (cadence sized to the timeline), with the current snapshot in the umbrella issue body and each scheduled update also as a comment. Each update must report how many tenants are still on the old form, how many have migrated, which `modify` issues are open, and how much time remains.

**Why this is a requirement, not a TR or decision:** The UX prescribes both the cadence shape and the metrics. It is a transparency demand on rollout coordination.

### BR-30: At the contract-change deadline, the old form must be removed and laggards must transition to eviction {#br-30}
**Source:** [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** On the contract-change deadline, the old form must be removed regardless of remaining tenants on it. For each tenant that has not migrated, the operator must open a separate eviction issue (linked to the umbrella) and the umbrella must close. No tenant may be silently broken on a removed offering.

**Why this is a requirement, not a TR or decision:** The UX prescribes this exact closeout behavior. It is the inverse of the evergreen promise — once communicated and given time, the deadline is real.

### BR-31: Eviction must be allowed when needs and capabilities fundamentally diverge {#br-31}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Move Off the Platform After Eviction §Persona]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#persona" >}})

**Requirement:** The platform must be able to decline continued hosting for a tenant whose requirements it cannot meet — specialized hardware, regulatory constraints, an availability target stronger than the platform offers. Eviction is initiated by the operator, not the capability owner.

**Why this is a requirement, not a TR or decision:** The capability defines this rule and the eviction UX operationalizes it. It is a control demand on what the platform is allowed to refuse.

### BR-32: Eviction-date negotiation must occur upstream of the eviction journey {#br-32}
**Source:** [UX: Move Off the Platform After Eviction §Entry Point]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#entry-point" >}})

**Requirement:** By the time the eviction issue is filed, the eviction date must already be agreed and not subject to renegotiation inside the eviction journey. The 30-day post-eviction retention is the only post-date slack and is fixed.

**Why this is a requirement, not a TR or decision:** The UX states this as a hard wall. It is a coordination commitment that protects both parties from re-litigation.

### BR-33: Eviction issues must contain the date, the reason, and a link to export tooling with documentation {#br-33}
**Source:** [UX: Move Off the Platform After Eviction §Entry Point]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#entry-point" >}})

**Requirement:** An eviction issue must carry exactly the eviction date, the reason for eviction, and a link to the export tool with documentation describing how to use it and the export shape. Nothing else is required of the issue.

**Why this is a requirement, not a TR or decision:** The UX names these contents and treats the issue as self-sufficient. It is a content commitment to the departing user.

### BR-34: Tenant compute and network must be torn down on the eviction date {#br-34}
**Source:** [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** On the eviction date, compute and network for the tenant must be torn down. Tenant data then enters the export-only, read-only state covered by [BR-11](#br-11) for the duration of the retention window — no further writes by anyone.

**Why this is a requirement, not a TR or decision:** The UX prescribes this teardown distinctly from the data-state guarantee. Compute/network teardown is what makes the dataset stable for export; the read-only data window itself is BR-11's commitment, referenced here rather than restated.

### BR-35: ~~Capability owner — not the platform — must notify their own end users of eviction~~ {#br-35}
**Status:** Removed on 2026-04-28 — absorbed into [BR-06](#br-06).

The platform-no-communication-with-end-users rule was already absolute in BR-06; the eviction-context clarifier and the capability-owner-notifies content are now folded into BR-06 directly. Number retained per the doc's append-only rule so existing TR citations (if any) stay valid.

### BR-36: Platform must offer a one-shot migration-process runner for capability-owner-supplied jobs {#br-36}
**Source:** [UX: Migrate Existing Data §Goal]({{< ref "user-experiences/migrate-existing-data.md#goal" >}}) · [UX: Migrate Existing Data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md#constraints-inherited" >}})

**Requirement:** The platform must provide a runner for one-time migration jobs that the capability owner writes and packages. The platform runs the process; it does not write, debug, or shepherd it.

**Why this is a requirement, not a TR or decision:** The migration UX is built around this offering and is explicit about the seam — platform runs, owner authors. It is a service commitment that bounds platform responsibility.

### BR-37: Migration jobs must be packaged in the same form as any other tenant component {#br-37}
**Source:** [UX: Migrate Existing Data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md#constraints-inherited" >}})

**Requirement:** The [BR-13](#br-13) packaging requirement applies to migration jobs without exception — the contract must not relax for migration. A process that cannot be packaged in the form the platform accepts cannot be run by the platform.

**Why this is a requirement, not a TR or decision:** The UX states this no-carve-out constraint explicitly. BR-13 establishes the packaging form; this BR forbids relaxing it for the migration case, which is the only place a relaxation might plausibly be argued for.

### BR-38: Migration jobs must declare their re-run contract and any temporary spikes up front {#br-38}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** A migration request must declare whether the process is safe to re-run against an already-populated destination (or requires a wiped destination), and any temporary migration-only spike beyond the tenant's steady-state footprint. Approval of spikes is bounded by what the platform can accommodate.

**Why this is a requirement, not a TR or decision:** The UX names both declarations as part of the operator's review scope. It is a content demand on what tenants must communicate, not a technical translation.

### BR-39: Migration peak footprint must not exceed twice the destination tenant's steady-state footprint {#br-39}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The peak temporary footprint of a migration (steady-state plus declared spike) must be no more than 2× the destination tenant's steady-state compute and storage. If either dimension exceeds that threshold, the request must be rejected as written; the capability owner is asked to split, reduce, or resize the tenant first.

**Why this is a requirement, not a TR or decision:** The UX states the 2× limit as a hard review rule. It bounds the burden one tenant's migration may place on the platform.

### BR-40: Concurrent migrations across tenants must be supported {#br-40}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The platform must support multiple migrations running at once across different tenants without changing each tenant's experience of their own journey. Tenants must not expect exclusive use of the migration runner.

**Why this is a requirement, not a TR or decision:** The UX specifies this property explicitly. It is a capacity commitment that prevents migrations from serializing.

### BR-41: Recovery from migration failure must follow the capability owner's plan, not a platform-prescribed model {#br-41}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** When a migration job fails or its output fails validation, the next step must be whatever plan the capability owner provides (wipe-and-retry, resume, accept partial, abandon). The platform must not auto-clean, auto-retry, or prescribe a recovery model.

**Why this is a requirement, not a TR or decision:** The UX places the recovery decision squarely with the data owner. It is an allocation-of-responsibility commitment.

### BR-42: A migration job artifact must be torn down on completion {#br-42}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** Once a migration job is closed (successful or abandoned), the platform must tear down the job. The platform must not retain it; re-running later means filing a fresh migration issue.

**Why this is a requirement, not a TR or decision:** The UX states the one-shot lifespan and tear-down explicitly. It is a lifecycle commitment that prevents migrations from accumulating into unmanaged state.

### BR-43: Platform must provide secret management for tenant-supplied credentials referenced by their components {#br-43}
**Source:** [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The platform must offer a secret-management surface that capability owners can populate independently of the operator, so their components and migration processes can reference credentials by name without leaking the secrets through engagement-thread comments.

**Why this is a requirement, not a TR or decision:** The migration UX assumes such an offering and operationalizes its use. It is a capability-level demand for handling credentials safely.

### BR-44: Each tenant must be provided compute, persistent storage, network reachability, identity, backup/DR, and observability {#br-44}
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md#outputs" >}})

**Requirement:** For each hosted tenant, the platform must provide compute (a place for the application to run), persistent storage durable to the platform's defined standard, network reachability both internal and external, identity and authentication for the tenant's end users, backup and disaster recovery for tenant data, and observability that lets the operator and capability owner tell whether the tenant is healthy.

**Why this is a requirement, not a TR or decision:** This is the capability's stated direct outputs. It is the inventory of what every tenant must receive, named in business terms. The identity entry has a tenant-choice carve-out captured separately in [BR-46](#br-46) (BYO identity); this BR commits to availability of the inventory, not to the platform being the sole source of identity.

### BR-45: Platform-provided identity service must support the "lost credentials cannot be recovered" property {#br-45}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Host a Capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md#constraints-inherited" >}})

**Requirement:** Any identity option the platform offers to tenants must be capable of honoring a Signal-style "lost credentials cannot be recovered" property. An identity option that cannot honor this property is not eligible to be the platform-provided identity service.

**Why this is a requirement, not a TR or decision:** The capability rule names this property explicitly because at least one tenant requires it. It is a forced constraint on the identity offering, not a vendor selection.

### BR-46: Tenants must be able to bring their own identity if they choose {#br-46}
**Source:** [Capability §Outputs & Deliverables]({{< ref "_index.md#outputs" >}}) · [Capability §Triggers & Inputs]({{< ref "_index.md#triggers" >}})

**Requirement:** Tenants must have the option to bring their own identity service rather than use the platform-provided one. Their decision is recorded in their tech design, not at onboarding time.

**Why this is a requirement, not a TR or decision:** The capability lists BYO identity as a tenant choice and the host-a-capability UX confirms it is recorded upstream of onboarding. It is a flexibility commitment.

### BR-47: Platform must be rebuildable to "ready to host tenants" within 1 hour {#br-47}
**Source:** [Capability §Success Criteria & KPIs]({{< ref "_index.md#success-criteria" >}}) · [UX: Stand Up the Platform §Goal]({{< ref "user-experiences/stand-up-the-platform.md#goal" >}})

**Requirement:** Starting from no platform at all (with definitions repo and root-level access in hand), the platform must be rebuildable to a ready-to-host-tenants state within 1 hour. The KPI is a target — exceeding it does not block the platform from going into service, but it must be tracked as a follow-up.

**Why this is a requirement, not a TR or decision:** This is the stated *Reproducibility* KPI of the capability and the standup UX is the journey it is measured against. It is a business commitment to recovery speed.

### BR-48: Rebuild readiness must be validated end-to-end by a purpose-built canary tenant {#br-48}
**Source:** [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** Standup must conclude with the deployment, exercise, and teardown of a purpose-built canary tenant maintained alongside the platform definitions. "Ready to host tenants" must be demonstrated by hosting a tenant — not declared from infrastructure self-checks alone.

**Why this is a requirement, not a TR or decision:** The standup UX defines this as the binding readiness signal. It is a confidence demand that infrastructure-only checks would not satisfy.

### BR-49: Each rebuild phase must support clean teardown of partial state {#br-49}
**Source:** [UX: Stand Up the Platform §Edge Cases]({{< ref "user-experiences/stand-up-the-platform.md#edge-cases" >}}) · [UX: Stand Up the Platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md#constraints-inherited" >}})

**Requirement:** Every phase of the rebuild must be reversible — "delete everything provisioned so far and start over" must be a viable, reliable option at every checkpoint. Partial state must not be trusted across a phase failure.

**Why this is a requirement, not a TR or decision:** The standup UX prescribes this property as part of "phase fails → tear down everything and restart." It is a reliability demand on the rebuild flow.

### BR-50: A reproducibility drill must run after every significant platform change and at least quarterly {#br-50}
**Source:** [UX: Stand Up the Platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md#entry-point" >}}) · [UX: Stand Up the Platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md#constraints-inherited" >}})

**Requirement:** The reproducibility KPI must be honestly evaluated by running a parallel rebuild drill on scratch infrastructure after every significant platform change (any change that would alter what is rebuilt, what must be validated, or what must be trusted) and at least quarterly while the live platform keeps serving.

**Why this is a requirement, not a TR or decision:** The standup UX names this cadence explicitly as the integrity check on the KPI. It is a discipline commitment, not a tool.

### BR-51: Platform must enforce tracked changes and immutability so drift can be detected before rebuild {#br-51}
**Source:** [UX: Stand Up the Platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md#entry-point" >}}) · [UX: Stand Up the Platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md#constraints-inherited" >}})

**Requirement:** Every platform UX that can introduce platform state must enforce tracked changes and immutability. The standup journey must perform a preflight drift check whenever prior platform state exists; the check must pass (no unexplained differences) before rebuild begins.

**Why this is a requirement, not a TR or decision:** The standup UX prescribes the preflight drift check and is explicit that drift must be prevented and detected outside the rebuild flow. It is an integrity commitment.

### BR-52: Rebuild must span both public and private infrastructure as part of foundations {#br-52}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** The platform may span public-cloud and home-lab infrastructure, and rebuild must establish the foundations — including connectivity between the two — as part of the standard standup flow. Cross-environment connectivity is foundational, not an afterthought.

**Why this is a requirement, not a TR or decision:** The capability allows the span; the standup UX makes Phase 1 explicitly cross-environment. It is a scope demand on what the rebuild must produce.

### BR-53: Tenant-facing observability must include a platform-standard health bundle {#br-53}
**Source:** [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** Each capability owner with a live tenant must receive, automatically, a tenant-scoped view of a platform-standard health bundle: availability, latency, error rate, resource saturation, and restart/deployment events. Capability owners must not have to instrument their own capability to see these signals.

**Why this is a requirement, not a TR or decision:** The observability UX defines the bundle and names automatic provisioning. It is a content commitment of the observability offering.

### BR-54: Capability owners must be able to self-serve their own alert thresholds {#br-54}
**Source:** [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** Within the observability offering, the capability owner must be able to tune the thresholds at which alerts fire to them, without operator involvement. The platform must not prescribe what counts as unhealthy enough to alert on.

**Why this is a requirement, not a TR or decision:** The observability UX names this as the one self-service surface and motivates it as a maintenance-budget pressure-relief. It is an authority demand — the user decides their own alerting.

### BR-55: Platform must push alerts to capability owners when their thresholds are crossed {#br-55}
**Source:** [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** When a tenant signal crosses a capability-owner-set threshold, the platform must send an alert to the capability owner that names which signal and which capability. The alert path is a best-effort nudge, not the source of truth.

**Why this is a requirement, not a TR or decision:** The observability UX names email as the current channel but the BR captures the demand (push alerts on threshold crossings, name the signal and capability). The channel is a downstream decision.

### BR-56: Tenant view must indicate degraded alert delivery when known {#br-56}
**Source:** [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}}) · [UX: Tenant-Facing Observability §Edge Cases]({{< ref "user-experiences/tenant-facing-observability.md#edge-cases" >}})

**Requirement:** When the observability offering knows its alert delivery to a tenant is degraded, the tenant view must surface that fact, so silence from the alert path is not mistaken for evidence of health. The pull view must remain authoritative for current health.

**Why this is a requirement, not a TR or decision:** The UX states this property explicitly and motivates it as a trust commitment. It is a transparency demand.

### BR-57: Tenant observability access must be scoped to the tenant; cross-tenant visibility is operator-only {#br-57}
**Source:** [UX: Tenant-Facing Observability §Entry Point]({{< ref "user-experiences/tenant-facing-observability.md#entry-point" >}}) · [UX: Tenant-Facing Observability §Constraints Inherited]({{< ref "user-experiences/tenant-facing-observability.md#constraints-inherited" >}})

**Requirement:** A capability owner authenticated to the observability offering must land directly in their own tenant's view and stay confined there for the rest of the session. There must be no mode-switch that broadens scope; only the operator sees across tenants.

**Why this is a requirement, not a TR or decision:** The UX names this isolation property and the capability's operator-only rule reinforces it. It is a confidentiality commitment.

### BR-58: Tenant observability access must be provisioned automatically as part of onboarding {#br-58}
**Source:** [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Tenant-Facing Observability §Entry Point]({{< ref "user-experiences/tenant-facing-observability.md#entry-point" >}})

**Requirement:** A capability owner whose onboarding has closed must already have a working login to the observability offering and a wired alert-delivery address — without filing a separate request. Observability is part of being hosted.

**Why this is a requirement, not a TR or decision:** Both UXes assume this is true at the moment a tenant is live. It is an integration commitment between onboarding and observability.

### BR-59: Routine operator maintenance must remain within 2 hours per week {#br-59}
**Source:** [Capability §Success Criteria & KPIs]({{< ref "_index.md#success-criteria" >}})

**Requirement:** The total routine operation of the platform — across all hosted tenants and platform-internal work — must take no more than 2 hours per week of the operator's time. If maintenance regularly exceeds this, the platform must be simplified, not grown.

**Why this is a requirement, not a TR or decision:** This is the *Operator maintenance budget* KPI. It is a hard upper bound on what the platform may demand of its operator and is referenced in nearly every UX as a pressure constraint.

### BR-60: A tenant whose accommodation would push routine maintenance sustainably above twice the maintenance budget must be evictable {#br-60}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}}) · [UX: Migrate Existing Data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md#constraints-inherited" >}}) · [UX: Platform-Contract-Change Rollout §Constraints Inherited]({{< ref "user-experiences/platform-contract-change-rollout.md#constraints-inherited" >}})

**Requirement:** When continuing to accommodate a tenant would push routine maintenance sustainably above 2× the maintenance budget, or break reproducibility (require manual snowflake configuration), the platform must be able to evict that tenant. Either condition alone must be sufficient grounds.

**Why this is a requirement, not a TR or decision:** The capability defines the eviction threshold and several UXes name it as the operative trigger. It is the control that prevents the maintenance budget from being eroded indefinitely.

### BR-61: Tenant adoption must be measured against implemented capabilities, with explicit-loss capture {#br-61}
**Source:** [Capability §Success Criteria & KPIs]({{< ref "_index.md#success-criteria" >}}) · [UX: Host a Capability §Edge Cases]({{< ref "user-experiences/host-a-capability.md#edge-cases" >}})

**Requirement:** Adoption is measured by counting only *implemented* capabilities (deployed and serving end users in production) — defined-only or designed-only capabilities are neutral. An implemented capability that runs elsewhere counts negatively, and a tenant lost because the operator went silent must be recorded explicitly on the issue rather than being silently dropped.

**Why this is a requirement, not a TR or decision:** The capability defines the KPI's mechanic; the host-a-capability UX defines the explicit-loss recording. It is a measurement-discipline commitment.

### BR-62: Operating cost must remain proportional to delivered convenience and resiliency {#br-62}
**Source:** [Capability §Success Criteria & KPIs]({{< ref "_index.md#success-criteria" >}})

**Requirement:** Total operating cost must remain within what the operator considers acceptable given the convenience and resiliency the platform delivers. There is no fixed dollar target; the test is whether the operator would still choose to run the platform knowing the bill.

**Why this is a requirement, not a TR or decision:** This is the *Cost stays proportional to value* KPI. It is a business commitment that bounds investment without prescribing a number.

### BR-63: Buy-vs-build trade-offs must be judged on convenience, resiliency, and cost only {#br-63}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}})

**Requirement:** When the platform decides between buying and building a component, the inputs must be convenience, resiliency, and cost. Operator skill development must not influence the trade-off; "I want to learn this" is not, on its own, a valid reason to choose build over buy.

**Why this is a requirement, not a TR or decision:** The capability rule explicitly forbids skill-development as an input. It is a control on decision-making, not a technical translation.

### BR-64: When a tenant needs something the platform does not yet provide, the default response must be to evolve the platform {#br-64}
**Source:** [Capability §Business Rules & Constraints]({{< ref "_index.md#business-rules" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Migrate Existing Data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md#constraints-inherited" >}})

**Requirement:** When a tenant capability requires something the platform does not yet offer, the default response must be to update the platform to provide it — bounded by the reproducibility and maintenance KPIs. The platform must not push the requirement back onto the tenant as a first response, but is not obligated to grow without bound.

**Why this is a requirement, not a TR or decision:** The capability defines this rule and the host-a-capability UX's "new offering needed" branch operationalizes it. It is the rule that keeps the platform tenant-aligned without making it infinitely extensible.

## Open Questions

- **Authoritative deeper-backup-tier policy after the 30-day post-eviction window.** Carried from [Move Off the Platform After Eviction §Open Questions]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}}). The retention duration, deletion behavior, and operator-access/privacy constraints of any backup-tier copies that survive past day 30 are not yet defined. BR-11 covers only the tenant-accessible-copy guarantee.
- **Tenant-facing pending-update visibility.** Carried from [Operator-Initiated Tenant Update §Entry Point]({{< ref "user-experiences/operator-initiated-tenant-update.md#entry-point" >}}). If the platform later adds a tenant-side surface for pending platform updates, BRs in the BR-20..BR-22 cluster will need a companion BR for that signal. Until then, the operator-filed issue remains the first official signal.
- **Volunteered-but-parked technical translations.** None volunteered during this extraction. Placeholder so re-extractions have a home for things like specific cadences, durability levels, or protocols that surface during conversation.
