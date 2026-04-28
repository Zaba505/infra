---
title: "Technical Requirements"
description: >
    Technical requirements derived from the Self-Hosted Application Platform capability's business requirements (with capability and UX docs as context). Each TR cites the BR-NN it derives from. Decisions belong in ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from `business-requirements.md` (and the capability/UX docs) on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `plan-adrs` skill will refuse to enumerate decisions until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [Self-Hosted Application Platform]({{< ref "_index.md" >}})
**Business requirements:** [business-requirements.md]({{< ref "business-requirements.md" >}})

## How to read this

Each TR is **forced** — by a BR (the primary case), by a prior shared ADR, or by a repo-wide constraint. It says what the technical solution must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review. If something has no BR or inherited-constraint source, raise a missing BR back to `extract-business-requirements`.

## Requirements

### TR-01: Tenant workloads must be runnable on platform-provided compute, storage, and network without per-tenant bespoke wiring
**Source:** [BR-01]({{< ref "business-requirements.md#br-01" >}}) · [BR-13]({{< ref "business-requirements.md#br-13" >}}) · [BR-23]({{< ref "business-requirements.md#br-23" >}})

**Requirement:** The platform must expose a single, uniform set of runtime surfaces — compute execution, persistent storage, internal-tenant-to-tenant network paths, and external (end-user-reachable) ingress — onto which any tenant whose declared needs fit the platform's contract can be provisioned without operator-authored, per-tenant snowflake configuration. Provisioning a new tenant must consume only parameters declared in that tenant's onboarding artifact and templates that already live in the platform's definitions.

**Why this is a TR, not a BR or decision:** BR-01 makes the platform the default hosting target, BR-13 enumerates compute/storage/network as outputs, and BR-23 forbids per-tenant snowflake configuration. The technical translation is a parameterized, definitions-driven provisioning surface; the BR side merely says "hosting is uniform and definitions-only."

### TR-02: The complete platform — including cross-environment connectivity — must be reproducible from the definitions repo plus root infrastructure credentials
**Source:** [BR-02]({{< ref "business-requirements.md#br-02" >}}) · [BR-35]({{< ref "business-requirements.md#br-35" >}}) · [BR-36]({{< ref "business-requirements.md#br-36" >}}) · [BR-42]({{< ref "business-requirements.md#br-42" >}})

**Requirement:** Every artifact required to bring the platform from "no platform exists" to "ready to host tenants" — including connectivity between public-cloud and private/home-lab environments — must be expressible as code or configuration in the definitions repo. The rebuild path must require only (a) the definitions repo and (b) root-level access credentials to the underlying infrastructure providers; no manual, undocumented, or operator-memory-resident steps may be on the critical path. Successor and primary operator must be able to run the same flow with the same inputs and reach the same outcome.

**Why this is a TR, not a BR or decision:** BR-02 demands rebuild from definitions; BR-35 and BR-36 force the cross-environment span to be in scope of "the platform"; BR-42 forces successor parity. The technical constraint is "all rebuild inputs are versioned artifacts plus root creds" — no specific IaC tool is chosen here.

### TR-03: Platform standup must execute as discrete, automated phases with explicit operator-validation checkpoints between them
**Source:** [BR-38]({{< ref "business-requirements.md#br-38" >}})

**Requirement:** The standup automation must be decomposed into named phases (foundations, core services, cross-cutting services, canary) such that each phase runs to completion automatically, then halts and awaits an explicit `continue` signal from the operator before the next phase begins. The system must not auto-advance between phases.

**Why this is a TR, not a BR or decision:** BR-38 demands phased execution with operator validation; the technical translation is "automation halts on phase boundary, waits for explicit human signal, resumes." The number, naming, and notification mechanism of phases are decisions.

### TR-04: A failed standup phase must support clean teardown of its partial output, and rebuild must always restart from phase zero
**Source:** [BR-39]({{< ref "business-requirements.md#br-39" >}})

**Requirement:** Every standup phase must produce output that the same definitions can fully tear down to the pre-phase state. The standup tooling must not offer a "resume from failed phase" mode against partial state. After teardown, a fresh rebuild must restart from the first phase.

**Why this is a TR, not a BR or decision:** BR-39 demands no-resume-from-partial-state and forces each phase to support clean teardown. The technical constraint is "every phase has an inverse, and there is no mid-flow resume entry point" — the teardown mechanism is the decision.

### TR-05: Standup readiness must be gated on a green end-to-end canary tenant exercise
**Source:** [BR-37]({{< ref "business-requirements.md#br-37" >}}) · [BR-23]({{< ref "business-requirements.md#br-23" >}})

**Requirement:** A standup run must not transition to "ready" status until a canary-tenant artifact (versioned alongside platform definitions) has been provisioned through the same tenant onboarding path real tenants use, exercised against compute (run), network (internal and external reachability), storage (read/write), identity (authentication), backup (pickup confirmation), and observability (signal pickup confirmation), and torn down. Failure or skip of any of these checks must keep the platform in a not-ready state regardless of elapsed time.

**Why this is a TR, not a BR or decision:** BR-37 forces the canary as the readiness gate and forbids self-checks-only. BR-23 forces the canary to use the same onboarding path so it actually validates that path. The canary's contents and tooling are decisions.

### TR-06: A drift check against live or last-known-good platform state must be a precondition of any rebuild
**Source:** [BR-40]({{< ref "business-requirements.md#br-40" >}})

**Requirement:** Where prior platform state exists, the rebuild flow must execute a drift check that compares the live platform (or a recorded last-known-good snapshot) against the definitions and refuse to begin rebuild if drift is detected. Detection must occur preflight, not partway through. The platform must additionally enforce, via every UX that can introduce platform state, mechanisms (tracked-changes review, immutability of provisioned resources, etc.) that prevent drift between rebuilds.

**Why this is a TR, not a BR or decision:** BR-40 demands drift-free-before-rebuild plus a cross-UX immutability discipline. The technical constraint is the existence of preflight comparison and platform-wide write-discipline; the diff implementation, snapshot store, and enforcement mechanism are decisions (parked in Open Questions).

### TR-07: A reproducibility drill — full parallel rebuild on scratch infrastructure — must be runnable on demand and on a recurring schedule
**Source:** [BR-41]({{< ref "business-requirements.md#br-41" >}}) · [BR-02]({{< ref "business-requirements.md#br-02" >}})

**Requirement:** The platform must support running its full standup flow against an isolated, scratch infrastructure target without disturbing the live platform, on operator demand and on a recurring (at minimum quarterly) cadence triggered after every significant platform-definition change. The drill must use the same definitions and the same flow as a real rebuild; a passing drill is the artifact by which the reproducibility KPI is asserted.

**Why this is a TR, not a BR or decision:** BR-41 demands the drill discipline. The technical constraint is "the standup flow can target a separate, isolated infrastructure target" — the scheduler, target-isolation mechanism, and pass/fail recording are decisions.

### TR-08: The platform must remain replaceable per-vendor — no component may make its provider's data, configuration, or runtime irrecoverable to the operator
**Source:** [BR-03]({{< ref "business-requirements.md#br-03" >}}) · [BR-35]({{< ref "business-requirements.md#br-35" >}})

**Requirement:** For every external vendor the platform depends on, the operator must retain the ability to extract that vendor's configuration and data and replace it with another implementation within a bounded migration effort. No vendor-proprietary data format, control surface, or service may be on the platform's critical path without an operator-controlled escape route. Vendor components are permitted; vendor lock-in (no escape route) is not.

**Why this is a TR, not a BR or decision:** BR-03 makes vendor independence an autonomy demand. The technical translation is "every vendor dependency has an extractable-and-replaceable boundary." Which vendors and which boundaries are decisions.

### TR-09: All platform administrative surfaces must be reachable only by an operator-held credential set, with no tenant-accessible administrative path
**Source:** [BR-04]({{< ref "business-requirements.md#br-04" >}}) · [BR-19]({{< ref "business-requirements.md#br-19" >}}) · [BR-34]({{< ref "business-requirements.md#br-34" >}})

**Requirement:** Every administrative interface the platform exposes (provisioning, configuration, infrastructure root, secret management admin, observability admin, eviction tooling) must require credentials held only by the operator. No capability owner or end user may possess any credential or token that grants administrative scope. There must be no in-product role-elevation path that grants tenant or end-user identities administrative scope.

**Why this is a TR, not a BR or decision:** BR-04 confines administration to the operator, BR-19 forbids cross-tenant widening, BR-34 forbids end-user platform surfaces. The technical translation is a hard authorization split between operator and everyone else, with no elevation path. The IAM scheme is a decision.

### TR-10: A successor operator must be able to acquire the full operator credential set via a sealed/escrowed handoff that is not active during normal operation
**Source:** [BR-05]({{< ref "business-requirements.md#br-05" >}}) · [BR-42]({{< ref "business-requirements.md#br-42" >}})

**Requirement:** The platform's credential model must support a sealed, escrowed copy of every credential needed to operate the platform, held in a way that the designated successor can unseal upon operator unavailability but cannot exercise during normal operation. Once unsealed, the successor's credentials must function identically to the primary operator's, and the standup/operate flows must not depend on operator-specific local state outside the definitions repo plus those credentials.

**Why this is a TR, not a BR or decision:** BR-05 forces successor takeover as a discrete event; BR-42 forces standup parity. The technical translation is a sealed, dormant, identical-on-unseal credential set with no operator-local out-of-band state in the standup path. The seal mechanism (password manager, physical envelope, threshold cryptography) is a decision.

### TR-11: Per-tenant data must be retrievable on demand by the capability owner as a complete, downloadable archive without operator action
**Source:** [BR-06]({{< ref "business-requirements.md#br-06" >}}) · [BR-27]({{< ref "business-requirements.md#br-27" >}})

**Requirement:** For every kind of data the platform hosts, the platform must provide an export tool, available from the moment a tenant goes live, that the capability owner can invoke at any time the platform is healthy without filing an issue or otherwise involving the operator. The tool must produce an archive that the capability owner can download immediately. Adding a new data kind to the platform must include adding export-tool support for it; there must be no platform-hosted data that lacks an export path.

**Why this is a TR, not a BR or decision:** BR-06 and BR-27 demand on-demand, universal, operator-free export. The technical constraint is "every data shape has a self-service export endpoint, and the catalogue of supported data shapes equals the catalogue of exportable shapes." Format, packaging, and transport are decisions.

### TR-12: Every export archive must be accompanied by a platform-computed checksum/hash and total byte size
**Source:** [BR-28]({{< ref "business-requirements.md#br-28" >}})

**Requirement:** When the platform produces an export archive, it must publish, alongside the archive, a content checksum/hash and the total size of the archive in bytes. Both values must be computed by the platform — not provided by the capability owner — and must be retrievable through the same surface that delivers the archive.

**Why this is a TR, not a BR or decision:** BR-28 requires platform-published bytes-level integrity verification. The hash algorithm and how the values are surfaced are decisions; the requirement is that both are produced and exposed.

### TR-13: The platform must enforce an authorized-by-operator gate before any tenant provisioning can occur
**Source:** [BR-07]({{< ref "business-requirements.md#br-07" >}})

**Requirement:** Tenant provisioning must require an authorization artifact tied to a specific operator action — there must be no API, UI, or workflow path by which a capability owner can cause a tenant to be provisioned without that artifact existing first. The authorization must be tied to a specific tenant submission so it cannot be replayed against a different submission.

**Why this is a TR, not a BR or decision:** BR-07 forbids self-onboarding. The technical translation is "no provisioning code path can run without an operator-authored authorization artifact bound to the specific request." Where the artifact lives and how it is signed are decisions.

### TR-14: Tenant submissions must arrive in a structured form that declares compute, storage, network reachability, availability expectations, and identity choice
**Source:** [BR-08]({{< ref "business-requirements.md#br-08" >}}) · [BR-14]({{< ref "business-requirements.md#br-14" >}})

**Requirement:** The onboarding submission format must require explicit, structured declarations of (a) compute requirements, (b) persistent storage needs, (c) internal and external network reachability needs, (d) availability expectations, and (e) whether the tenant uses platform-provided identity or brings its own. The platform must reject submissions that omit any of these. Acceptance of the platform's contract is implicit in the submission.

**Why this is a TR, not a BR or decision:** BR-08 forces the up-front declaration and contract acceptance; BR-14 forces the identity declaration. The technical constraint is a structured, validated submission shape with required fields. The exact schema/format is a decision.

### TR-15: Modifications to a live tenant must operate only on the declared delta and must not require re-validation of the full platform contract
**Source:** [BR-44]({{< ref "business-requirements.md#br-44" >}})

**Requirement:** The `modify my capability` workflow must accept a delta against the current tenant state, restrict review and change to that delta, and not re-prompt for or re-record platform-contract acceptance. Any change to the platform contract itself must enter through the contract-change rollout pipeline, not through the modify path.

**Why this is a TR, not a BR or decision:** BR-44 forbids re-acceptance on modification. The technical translation is "modify operates on diffs against tenant state, not on full re-submissions, and the contract-acceptance code path is unreachable from the modify flow." The diff format is a decision.

### TR-16: Platform-contract changes must be deliverable as concurrent old/new offerings during a rollout window, except for full removals
**Source:** [BR-09]({{< ref "business-requirements.md#br-09" >}}) · [BR-45]({{< ref "business-requirements.md#br-45" >}})

**Requirement:** The platform must support running both the old and the new form of a contract term concurrently for a bounded rollout window, allowing tenants to migrate at their own pace before the old form is retired. The exception is full removals (no replacement offering), which the platform must support as an all-or-nothing cut at a single deadline. The platform must support pinning a tenant to either old or new form during the window.

**Why this is a TR, not a BR or decision:** BR-09 forces no-surprise migration; BR-45 forces the side-by-side property and its single carve-out. The technical translation is "the runtime can serve N concurrent contract versions for a bounded window with per-tenant pinning." The mechanism (parallel deployment, feature flag, gateway routing) is a decision.

### TR-17: The platform must be able to enumerate, on demand, every tenant affected by a given contract term
**Source:** [BR-09]({{< ref "business-requirements.md#br-09" >}}) · [BR-46]({{< ref "business-requirements.md#br-46" >}}) · [BR-47]({{< ref "business-requirements.md#br-47" >}}) · [BR-48]({{< ref "business-requirements.md#br-48" >}})

**Requirement:** The platform must maintain a queryable mapping from contract terms (offerings, packaging forms, availability characteristics) to the tenants currently consuming them, such that the operator can produce an up-to-date "affected tenants" list at any time during a rollout, and a current count of (still-on-old, migrated, in-flight) for status-update publication.

**Why this is a TR, not a BR or decision:** BR-09 forces ahead-of-time communication to every affected tenant; BR-46/47/48 require deadline-floor calculation, status updates with concrete counts, and acknowledgment tracking. All of these need the affected-tenant set to be derivable, not guessed. The data model and query surface are decisions.

### TR-18: The platform must enforce that a tenant declines tenant-provisioning when continued accommodation would breach the maintenance-budget or reproducibility KPIs
**Source:** [BR-10]({{< ref "business-requirements.md#br-10" >}}) · [BR-32]({{< ref "business-requirements.md#br-32" >}})

**Requirement:** The platform's onboarding and modify workflows must surface, at review time, an assessment of whether accommodating the new or modified tenant would (a) push routine operation sustainably above twice the operator-maintenance-budget KPI (currently 2 hours/week) or (b) require configuration that cannot be captured in the definitions repo. Either condition must produce a hard-stop signal that the operator's review path makes visible. A tenant whose modification triggers either condition must not be silently accepted.

**Why this is a TR, not a BR or decision:** BR-10 grants and demands eviction/decline rights tied directly to BR-32's budget KPI. The technical translation is a review-time check producing a hard-stop signal; the assessment heuristic and budget-tracking method are decisions.

### TR-19: A migration job's declared peak temporary footprint must be checked against a bounded multiple of the destination tenant's steady-state footprint at review time
**Source:** [BR-26]({{< ref "business-requirements.md#br-26" >}})

**Requirement:** The migration-process review surface must compare the migration job's declared peak compute and storage footprint against the destination tenant's steady-state declarations and reject the submission when the declared peak exceeds the platform's bound (currently 2x). Rejection must surface specifically that the peak-vs-steady-state cap was breached, so the capability owner can split, reduce, or precede with a tenant resize.

**Why this is a TR, not a BR or decision:** BR-26 makes the bounded-footprint check load-bearing during migration review. The technical translation is a declared-capacity comparison with a configured ratio; the exact ratio and the comparison heuristic are decisions.

### TR-20: The platform must offer a one-shot migration job runner with a provision-run-observe-teardown lifecycle, owned by the platform
**Source:** [BR-24]({{< ref "business-requirements.md#br-24" >}})

**Requirement:** The platform must expose a generic facility that accepts a capability-owner-supplied migration artifact, provisions execution resources for it on operator approval, runs it once with read access to a configured prior host and write access to the capability owner's tenant, exposes its progress and logs to the platform's observability surface while running, and tears the execution resources down on completion (success or failure). The lifecycle must be platform-managed; the migration logic correctness is the capability owner's responsibility.

**Why this is a TR, not a BR or decision:** BR-24 demands the offering and the lifecycle. The technical translation is the platform-owned lifecycle state machine and the capability-owner-supplied artifact contract; the runner technology, scheduler, and packaging form are decisions.

### TR-21: The platform must offer a secret-management surface that allows capability owners to register named credentials referenced by tenant artifacts
**Source:** [BR-25]({{< ref "business-requirements.md#br-25" >}})

**Requirement:** The platform must provide a secret-store interface through which capability owners can register credentials (e.g. credentials a migration job needs to read from an old host) tied to their own tenant scope, retrievable at runtime by the tenant's workloads via name reference. The plaintext value of a registered credential must never be required to appear in issue threads or in any committed artifact; only the name reference does.

**Why this is a TR, not a BR or decision:** BR-25 makes secret management a precondition for migration filing and, by extension, for any workflow that needs runtime credentials. The technical translation is a name-keyed, scope-confined secret store with a write surface for capability owners; the product choice and the store's storage backend are decisions.

### TR-22: Tenant-scoped observability must be exposed to the responsible capability owner via both a pull view and a push alert path
**Source:** [BR-17]({{< ref "business-requirements.md#br-17" >}}) · [BR-18]({{< ref "business-requirements.md#br-18" >}})

**Requirement:** The platform must produce, for every live tenant, a health signal bundle (covering availability, latency, error rate, resource saturation, and restart/deployment events) usable by the operator (BR-17) and additionally exposed to the responsible capability owner via (a) an on-demand pull view that always reflects the current state of their own tenant and (b) a push alert path that fires when a tenant-scoped signal crosses a configured threshold. The capability owner must not need to instrument these signals themselves.

**Why this is a TR, not a BR or decision:** BR-17 establishes operator-side observability as an output; BR-18 forces the dual-mode capability-owner-facing surface. The technical translation is a platform-emitted standard signal bundle with both pull and push delivery surfaces; the specific signal definitions, view software, and alert delivery channel are decisions.

### TR-23: The capability-owner observability surface must restrict every session to exactly one tenant scope, with no widening path
**Source:** [BR-19]({{< ref "business-requirements.md#br-19" >}}) · [BR-04]({{< ref "business-requirements.md#br-04" >}})

**Requirement:** When the observability surface authenticates a capability owner, the resulting session must be hard-bound to the tenant(s) that owner is responsible for, with all queries, alert configurations, and dashboards filtered to that scope on the server side. There must be no UI control, query parameter, API method, or role-switch that allows a capability-owner-authenticated session to read data from another tenant. Cross-tenant visibility must remain reachable only through operator-scoped credentials.

**Why this is a TR, not a BR or decision:** BR-19 demands per-tenant isolation of capability-owner observability sessions; BR-04 forbids any administrative widening. The technical translation is server-side scope enforcement, not a client-side filter, and a closed widening path. Which authn/authz model achieves this is a decision.

### TR-24: Capability owners must be able to set, change, and delete their own tenant-scoped alert thresholds without operator action
**Source:** [BR-20]({{< ref "business-requirements.md#br-20" >}})

**Requirement:** The observability surface must expose a self-service interface — accessible to the capability owner under the same scope-confinement rules as TR-23 — that lets them create, modify, and remove the threshold rules driving their own tenant's push alerts. No operator approval step or out-of-band issue may sit on this path.

**Why this is a TR, not a BR or decision:** BR-20 carves out alert-threshold self-service as the one capability-owner-self-service surface. The technical translation is a self-service write surface for threshold rules, scope-confined per TR-23. The UI form and the rule-storage shape are decisions.

### TR-25: The tenant-scoped pull view must surface push-alert delivery degradation when the platform detects it
**Source:** [BR-21]({{< ref "business-requirements.md#br-21" >}})

**Requirement:** The platform must monitor its own push-alert delivery (for example, email send failures or bounces) per tenant and, when delivery for a tenant is failing or degraded, surface that condition prominently within that tenant's pull view. The pull view must remain the source of truth for tenant health; the alert-degradation indicator must signal to the capability owner that they cannot rely on alert silence to mean health.

**Why this is a TR, not a BR or decision:** BR-21 forces the platform to expose alert-delivery degradation in the pull view. The technical translation is "delivery health is a first-class signal on the pull view, computed from observed delivery outcomes." The detection mechanism and the rendering are decisions.

### TR-26: The platform-provided identity service must support a mode in which credential recovery is impossible
**Source:** [BR-15]({{< ref "business-requirements.md#br-15" >}})

**Requirement:** Any identity implementation that the platform offers as its tenant-end-user identity service must support a configuration in which a tenant's end users, having lost their credential, cannot recover account access — there must be no operator-, platform-, or vendor-mediated reset that restores access. An identity option lacking this configuration mode is ineligible to be the platform-provided identity service.

**Why this is a TR, not a BR or decision:** BR-15 makes "lost credentials cannot be recovered" a property the identity offering must support. The technical translation is a no-recovery configuration mode that the platform can guarantee per tenant. The choice of identity software is a decision.

### TR-27: The platform must back up tenant data and provide a tested disaster-recovery path against a published platform-side standard
**Source:** [BR-16]({{< ref "business-requirements.md#br-16" >}})

**Requirement:** For every kind of tenant data the platform hosts, the platform must take backups on a schedule, retain them per a defined policy, and support restoring them in a documented disaster-recovery procedure. The standard (frequency, retention, RPO/RTO) must be published as part of the platform's contract so tenants accept it explicitly via TR-14. The DR procedure must be exercisable end-to-end without requiring tenant participation.

**Why this is a TR, not a BR or decision:** BR-16 demands backup, DR, and a platform-defined standard. The technical translation is "scheduled backup of every supported data shape, defined retention, defined recovery procedure, and a published standard tied to the contract." The exact RPO/RTO numbers and storage choices are decisions, parked in Open Questions until the standard is set.

### TR-28: Eviction cutover must atomically deprovision the tenant's compute and network and freeze the tenant's data into a read-only export-only state on the eviction date
**Source:** [BR-31]({{< ref "business-requirements.md#br-31" >}})

**Requirement:** On the recorded eviction date for a tenant, the platform must (a) deprovision that tenant's compute and live-serving network resources, (b) transition every persistent-storage location belonging to that tenant into a state in which no further writes (by anyone, including the capability owner) succeed, and (c) record cutover confirmation in the eviction issue thread. Read access by the export tool must remain functional through the retention window.

**Why this is a TR, not a BR or decision:** BR-31 forces the cutover and the read-only freeze as load-bearing for the eviction journey. The technical translation is "compute/network teardown and write-disablement happen atomically per tenant on a scheduled date." The mechanism (IAM revoke, bucket policy, storage snapshot) is a decision.

### TR-29: After eviction cutover, tenant-accessible data must remain available exactly for a fixed retention window, then be removed from all tenant-accessible surfaces
**Source:** [BR-29]({{< ref "business-requirements.md#br-29" >}})

**Requirement:** From eviction cutover, every export-only surface for that tenant must remain reachable for a fixed retention window (currently 30 days) without any extension. At the end of the window, the platform must remove every tenant-accessible copy of that tenant's data from those surfaces. The system must not expose any operator workflow to extend the window for slow extracts or capability-owner requests.

**Why this is a TR, not a BR or decision:** BR-29 demands a fixed, no-slip retention window with a defined end-state. The technical translation is "the retention timer is computed from cutover date, is non-extendable through any tenant-facing path, and triggers an enforced removal." The exact 30-day number and the removal mechanism are decisions; deeper backup-tier handling is parked in Open Questions.

### TR-30: The eviction retention countdown must be pause-able by the operator only when a platform-rooted export defect is recorded
**Source:** [BR-30]({{< ref "business-requirements.md#br-30" >}})

**Requirement:** The retention-window timer (TR-29) must support being paused by the operator on the basis of a recorded platform-side defect that prevents a clean export. The pause must be tied to that defect record so it ends when the defect is resolved and a clean export has been produced. Capability-owner-rooted failures must not trigger or maintain a pause.

**Why this is a TR, not a BR or decision:** BR-30 carves out the one exception to the otherwise-hard retention wall and ties it specifically to platform-side defects. The technical translation is "the timer is pausable, the pause is gated by a defect-record link, and the resumption gate is a clean export." The defect-tracking mechanism and the operator UI are decisions.

### TR-31: After the retention window, the platform must hold no tenant-accessible copy of the evicted tenant's data on any surface reachable by capability owners or end users
**Source:** [BR-29]({{< ref "business-requirements.md#br-29" >}}) · [BR-54]({{< ref "business-requirements.md#br-54" >}})

**Requirement:** Once the retention window for an evicted tenant has elapsed (and any TR-30 pause has resolved), the platform must guarantee that no surface reachable by the former capability owner or any end user returns the tenant's data — exports must fail, observability of historical data must be unreachable to the former owner, and any administrative views remain operator-scoped. Whether deeper backup-tier copies persist for operator-only purposes is governed by a separate (currently open) policy.

**Why this is a TR, not a BR or decision:** BR-29 and BR-54 force the post-window absence of any tenant-accessible copy. The technical translation is a hard removal across every tenant-facing surface; the deeper backup-tier policy is explicitly out of scope of this TR and parked in Open Questions.

### TR-32: A platform-update-required workflow must accept an externally-inherited deadline per tenant and persist the external reason on the tenant's update record
**Source:** [BR-50]({{< ref "business-requirements.md#br-50" >}}) · [BR-51]({{< ref "business-requirements.md#br-51" >}})

**Requirement:** The operator-initiated-tenant-update workflow must accept, per affected tenant, (a) a deadline value sourced from the external forcing event (vendor sunset, CVE, EOL) and (b) a recorded external-reason string identifying that event. The system must distinguish this issue type from `modify` and `onboard` so downstream views can recognize required updates. The recorded deadline must be the basis for any later slack-bounded extension negotiation, and any agreed extension must be recorded as a separate operative-delivery date alongside the original.

**Why this is a TR, not a BR or decision:** BR-50 forces the inherited-deadline contract and the typed signal; BR-51 forces the bounded-extension recording. The technical translation is a typed, deadline-bearing, externally-sourced issue/record with a separate operative date. The CVE/EOL feed source and the negotiation flow are decisions.

### TR-33: A missed operative delivery date on a platform-update-required or contract-change rollout must trigger a separate, linked eviction record
**Source:** [BR-52]({{< ref "business-requirements.md#br-52" >}})

**Requirement:** When the operative delivery date on a platform-update-required record (or the global deadline on a contract-change rollout) passes without the tenant having shipped, the platform's workflow must (a) open a new eviction record linked back to the originating record and (b) close the originating record as superseded by eviction. Eviction must always live in its own record, never as a state field on the originating record.

**Why this is a TR, not a BR or decision:** BR-52 prescribes the separate-and-linked structure. The technical translation is a workflow transition that creates a new linked record and closes the source. The exact triggering mechanism (manual operator step, automated transition) is a decision.

### TR-34: A platform-update-required record must be pause-able pending completion of a new-offering branch and must resume at the modify inner loop
**Source:** [BR-53]({{< ref "business-requirements.md#br-53" >}})

**Requirement:** When shipping an in-flight platform-update-required record reveals a need for a platform offering that does not yet exist, the workflow must allow the record to remain open (not force-closed) while a parallel new-offering branch (per the host-a-capability flow) is exercised, then resume the modify inner-loop step on the original record once the new offering is available.

**Why this is a TR, not a BR or decision:** BR-53 forbids force-closing the update because the platform was not yet ready. The technical translation is "the update record has a pause-and-resume state tied to a separate new-offering record." The cross-record linking surface is a decision.

### TR-35: A contract-change rollout's umbrella record must support tracked acknowledgment per affected capability owner, with deadline-driven escalation on missing acknowledgments
**Source:** [BR-48]({{< ref "business-requirements.md#br-48" >}})

**Requirement:** A contract-change umbrella record must enumerate every affected capability owner (per TR-17) and track an acknowledgment state per owner. At the deadline, capability owners with no acknowledgment must be flagged as non-engagement and routed into the laggard branch (separate eviction record per TR-33). Silence must not be representable as tacit consent in the data model.

**Why this is a TR, not a BR or decision:** BR-48 demands acknowledgment as load-bearing and silence as non-engagement. The technical translation is a per-owner acknowledgment field with deadline-driven state transition. The acknowledgment UI is a decision.

### TR-36: A contract-change rollout's umbrella record must carry a current-snapshot summary in its body and a chronological status-update log
**Source:** [BR-47]({{< ref "business-requirements.md#br-47" >}})

**Requirement:** A contract-change umbrella record must support (a) a body field holding the current snapshot (so a cold reader sees latest state) and (b) an append-only log of dated status updates. Each scheduled status update must populate both: the body is overwritten with the new snapshot, and the same content is appended as a new log entry. Each entry must include counts of (still-on-old, migrated), open `modify` records related to the rollout, and time remaining to deadline.

**Why this is a TR, not a BR or decision:** BR-47 forces both the standing snapshot and the historical log with specific required content. The technical translation is "the record schema supports a mutable body plus an immutable comment/log stream, and the publish step writes both." The cadence and the snapshot template are decisions.

### TR-37: All capability-owner-to-operator engagement must flow through a typed-issue tracker, with one persistent issue type per workflow
**Source:** [BR-22]({{< ref "business-requirements.md#br-22" >}})

**Requirement:** The platform must rely on a single, persistent issue tracker (currently GitHub Issues against the infra repo) as the only sanctioned engagement channel between capability owners and the operator. The tracker must support distinct issue types per workflow — `onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, eviction — such that the type alone signals the operator's review scope and journey shape. The platform must not accept work requests through other channels (email, chat, DMs).

**Why this is a TR, not a BR or decision:** BR-22 makes the typed-issue-thread engagement model load-bearing for every UX. The technical translation is a single tracker dependency with typed templates and a no-other-channels operating discipline. The tracker product is a decision (currently GitHub).

### TR-38: The platform must have no surface (UI, API, identity, communication channel) reachable by tenant end users
**Source:** [BR-34]({{< ref "business-requirements.md#br-34" >}})

**Requirement:** The platform must expose no UI, API, status page, notification path, or identity surface to end users of tenant capabilities. End-user-reachable surfaces are exclusively the tenant's responsibility. Any platform-side state that might affect end users (eviction, outage, contract change) must be communicated only via the capability owner, not directly to end users.

**Why this is a TR, not a BR or decision:** BR-34 declares end users out of scope of the platform's actor set and forbids direct platform↔end-user communication. The technical translation is "no end-user identity in the platform's identity model, no end-user-reachable URL surface, no end-user-addressed notifications." The boundary enforcement (gateway, IAM scope) is a decision.

### TR-39: Cost minimization must be subordinate to convenience and resiliency in component selection, but every selected component's cost must be observable to the operator
**Source:** [BR-33]({{< ref "business-requirements.md#br-33" >}})

**Requirement:** The platform must make cumulative operating cost — broken down by major component or vendor — observable to the operator on an ongoing basis, so that proportionality between cost and delivered convenience/resiliency can be assessed. The platform must not require cost-minimizing component choices when those choices would degrade convenience or resiliency, but the cost of non-minimal choices must be visible enough to be re-evaluated.

**Why this is a TR, not a BR or decision:** BR-33 is a value-judgment on spending; the technical translation is "cost is an observable signal at a granularity sufficient for the operator's proportionality test." The cost-tracking source and reporting surface are decisions.

## Open Questions

Things the user volunteered as solutions during extraction (parked for the ADR stage), or constraints the capability/UX docs don't yet make explicit.

- **Drift-detection mechanism (TR-06).** What constitutes "live state" vs. "last known-good environment", what diff algorithm produces the drift signal, where the policy is enforced (CI hook vs. preflight script vs. controller), and the immutability mechanism that prevents drift between rebuilds. Carried forward from the BR doc; this is a TR/ADR concern that this skill stops short of choosing.
- **Backup-and-DR standard numbers (TR-27).** The platform's RPO, RTO, backup frequency, and retention policy. BR-16 demands a "standard the platform defines" but the capability doc does not yet define it. TR-27 cannot be fully satisfied until this is set; raise back to the capability doc rather than choosing here.
- **Deeper backup-tier policy after eviction retention window (TR-29, TR-31).** Whether platform-side backups outlive the tenant-accessible retention window, for how long, with what access controls, and what privacy posture. BR-29's "no tenant-accessible copy" leaves operator-accessible copies unspecified. Carried forward from the BR doc; awaits a capability/UX-level decision before TR-29/TR-31 can be considered complete.
- **Tenant-facing observability signal bundle and default thresholds (TR-22).** The exact signals the platform commits to emitting and any platform-default thresholds. BR-18 references the bundle but does not specify it; this is translation territory and the bundle composition is a decision for `plan-adrs` / `define-adr`.
- **Migration-process concurrency capacity (TR-20).** The migrate-existing-data UX promises concurrent migrations across tenants but does not bound capacity. Capacity sizing, queueing, and per-tenant concurrency limits are TR/ADR-level concerns that this extraction does not cover.
- **Pending-update tenant-facing view (TR-22 / future BR).** The operator-initiated-tenant-update UX notes that today, capability owners learn of forced updates only when their per-tenant issue is filed. Adding an earlier deprecation/pending-update signal would extend BR-18 / TR-22; parked for tenant-facing-observability evolution rather than added speculatively here.
- **Explicit anchors on capability and UX section headings.** Most of the citations in this doc anchor at page level only because the capability and UX source headings do not yet carry explicit `{#anchor-id}` annotations. BR-NN anchors in `business-requirements.md` exist (citations to `#br-NN` here assume the BR headings carry those anchors, per the BR doc's own structure). Section-level anchors on `_index.md` and the UX pages still need to be added before TRs can deep-link to them.
- **Platform-vs-tenant cost attribution model (TR-39).** What granularity of cost attribution is sufficient — per-component, per-vendor, per-tenant — and which signal source(s) feed it. BR-33 forces the proportionality test but not the attribution mechanism.
- **Vendor-escape boundary list (TR-08).** The catalogue of vendor dependencies and the specific escape route for each is enumerated as ADRs are filed; this TR establishes the constraint without naming the vendors.

