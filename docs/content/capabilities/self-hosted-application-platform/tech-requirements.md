---
title: "Technical Requirements"
description: >
    Technical requirements derived from the Self-Hosted Application Platform capability's business requirements, with the capability and UX docs as supporting context. Each TR cites the BR-NN it derives from. Decisions belong in ADRs, not here.
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

### TR-01: Platform state must be entirely expressible as version-controlled definitions {#tr-01}
**Source:** [BR-02]({{< ref "business-requirements.md#br-02" >}}) · [BR-51]({{< ref "business-requirements.md#br-51" >}}) · [UX: Stand Up the Platform §Constraints Inherited]({{< ref "user-experiences/stand-up-the-platform.md#constraints-inherited" >}})

**Requirement:** Every piece of platform runtime state — each offering, every per-tenant binding, every shared piece of configuration the platform depends on — must be expressible in a tracked-changes definitions repository. Anything modifiable outside that repository is drift, and any UX that introduces platform state must route through the same recorded-change surface.

**Why this is a TR, not a BR or decision:** BR-02 demands reproducibility; BR-51 demands tracked changes and immutability so drift can be detected. The technical translation is that the definitions repository is the only authoritative surface for platform-modifying writes. Which repository, which tracked-changes mechanism, and which immutability discipline are downstream decisions.

### TR-02: Platform must expose a single top-level rebuild entry point that runs end-to-end from definitions in ≤60 minutes {#tr-02}
**Source:** [BR-02]({{< ref "business-requirements.md#br-02" >}}) · [BR-47]({{< ref "business-requirements.md#br-47" >}}) · [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** A single operator-invocable entry point must drive the rebuild from a fresh pull of the definitions repository, sequence the foundations → core services → cross-cutting → canary phases automatically, and be capable of completing within 60 minutes of wall-clock time on the target infrastructure when run end-to-end. Manual checkpoints between phases are permitted; manual driving of each step is not.

**Why this is a TR, not a BR or decision:** BR-47 is the 1-hour rebuild target; BR-02 is the demand that rebuild is from definitions only. The TR is the operative property — one entry point, automated, time-bounded — without naming a specific automation tool, language, or orchestrator.

### TR-03: Rebuild Phase 1 must establish foundations across both public-cloud and home-lab environments and the connectivity between them {#tr-03}
**Source:** [BR-52]({{< ref "business-requirements.md#br-52" >}}) · [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** The first rebuild phase must provision the public-cloud-side and home-lab-side foundations and the cross-environment connectivity between them, before any later phase proceeds. Single-environment standup (public-only or home-lab-only) is not a supported rebuild outcome.

**Why this is a TR, not a BR or decision:** BR-52 asserts the platform may span both environments and that connectivity is part of foundations, not an afterthought. The TR forces foundations-phase scope to include both sides plus the link; choosing the specific cloud, home-lab hardware, or tunnel mechanism is downstream.

### TR-04: Each rebuild phase must support a deterministic, definitions-driven teardown of all state it produced {#tr-04}
**Source:** [BR-49]({{< ref "business-requirements.md#br-49" >}}) · [UX: Stand Up the Platform §Edge Cases]({{< ref "user-experiences/stand-up-the-platform.md#edge-cases" >}})

**Requirement:** Every rebuild phase must expose a deterministic, definitions-driven teardown that removes every resource the phase produced, callable at every checkpoint. "Delete everything provisioned so far and start over" must be a viable, reliable option at each phase boundary. Partial state must not be carried across a phase failure into the next phase.

**Why this is a TR, not a BR or decision:** BR-49 demands that partial state never be trusted across a phase failure. The TR is the operative property — every phase has a clean teardown — without prescribing a teardown mechanism.

### TR-05: Rebuild flow must perform a preflight drift check that fails closed when prior platform state exists and unexplained differences remain {#tr-05}
**Source:** [BR-51]({{< ref "business-requirements.md#br-51" >}}) · [UX: Stand Up the Platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md#entry-point" >}})

**Requirement:** Before any rebuild begins, the platform must compare current platform state against a last-known-good reference and refuse to proceed if unexplained differences remain. On a first-ever build the check is vacuously satisfied; in every other case the check must pass before later phases run.

**Why this is a TR, not a BR or decision:** BR-51 demands that drift be detected outside the rebuild flow rather than discovered partway through it. The TR makes the preflight check a property of the rebuild entry point. The mechanism by which "last-known-good reference" is captured and compared is a downstream decision.

### TR-06: Rebuild flow must be runnable on parallel/scratch infrastructure without affecting the live platform {#tr-06}
**Source:** [BR-50]({{< ref "business-requirements.md#br-50" >}}) · [UX: Stand Up the Platform §Entry Point]({{< ref "user-experiences/stand-up-the-platform.md#entry-point" >}})

**Requirement:** The same definitions and the same rebuild entry point must be invocable against scratch infrastructure to support post-significant-change drills and at-least-quarterly drills, without touching live platform state. Drill mode and live mode must differ only in the underlying target.

**Why this is a TR, not a BR or decision:** BR-50 demands honest evaluation of the reproducibility KPI via parallel rebuilds. The TR forces drill-vs-live parity at the entry-point level; how target selection is parameterized is a downstream decision.

### TR-07: Platform must include a purpose-built canary tenant maintained alongside the definitions, used as the rebuild's binding readiness signal {#tr-07}
**Source:** [BR-48]({{< ref "business-requirements.md#br-48" >}}) · [UX: Stand Up the Platform §Journey]({{< ref "user-experiences/stand-up-the-platform.md#journey" >}})

**Requirement:** A canary tenant must be maintained alongside the platform definitions and must be deployed, exercised end-to-end against every platform-provided service, and torn down by the rebuild's final phase. Readiness must not be declared on infrastructure self-checks alone; the canary's pass/fail is the binding signal.

**Why this is a TR, not a BR or decision:** BR-48 names the canary as the readiness mechanism. The TR makes the canary a first-class artifact of the platform definitions. What the canary's workload looks like, and which signals it must produce, are downstream decisions.

### TR-08: Platform must accept tenant components — including migration jobs — only in a single pre-declared packaging form, with no carve-outs {#tr-08}
**Source:** [BR-13]({{< ref "business-requirements.md#br-13" >}}) · [BR-37]({{< ref "business-requirements.md#br-37" >}}) · [UX: Migrate Existing Data §Constraints Inherited]({{< ref "user-experiences/migrate-existing-data.md#constraints-inherited" >}})

**Requirement:** Exactly one packaging form is admissible to the platform. Any tenant component, including migration job artifacts, must arrive in that form to be runnable; the migration path must not relax it. Components that cannot be packaged this way cannot run on the platform.

**Why this is a TR, not a BR or decision:** BR-13 commits the platform to a packaging form; BR-37 forbids relaxing it for migration. The TR forces single-form admission; the actual form (container image, OCI bundle, archive, etc.) is an ADR.

### TR-09: Onboarding must require machine-readable declarations of resource needs, packaged artifact, identity choice, and availability acceptance before any provisioning is possible {#tr-09}
**Source:** [BR-13]({{< ref "business-requirements.md#br-13" >}}) · [BR-46]({{< ref "business-requirements.md#br-46" >}}) · [UX: Host a Capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md#constraints-inherited" >}})

**Requirement:** A tenant onboarding submission must, before provisioning is possible, include machine-readable declarations of (a) the tenant's resource needs (compute, storage, network), (b) the packaged artifact in the platform's accepted form, (c) the identity choice — platform-provided or BYO — recorded in the tech design, and (d) acceptance of the platform's current availability characteristics. Approval binds the runtime to those declarations.

**Why this is a TR, not a BR or decision:** BR-13 names the four declarations as the price of admission; BR-46 makes the identity choice one of them. The TR makes the declarations a hard precondition of the provisioning gate; the schema and review surface are downstream decisions.

### TR-10: Tenant runtime must not be provisioned without an explicit operator-issued authorization signal tied to the onboarding artifact {#tr-10}
**Source:** [BR-14]({{< ref "business-requirements.md#br-14" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** Provisioning of a new tenant runtime must be gated on a per-tenant authorization signal issued by the operator's identity and bound to the specific onboarding submission. There must be no provisioning path that bypasses this gate, and there must be no self-service onboarding path.

**Why this is a TR, not a BR or decision:** BR-14 forbids self-onboarding and demands explicit authorization. The TR forces a control point on the provisioning surface; how the authorization signal is represented (issue comment, signed approval, etc.) is downstream.

### TR-11: Onboarding flow must support a "new offering needed" hold and resume without requiring the capability owner to refile {#tr-11}
**Source:** [BR-64]({{< ref "business-requirements.md#br-64" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** When an onboarding tenant requires an offering the platform does not yet provide, the onboarding record must be holdable in a pending state and resumable from that point once the offering is added — without the capability owner refiling, restarting, or re-accepting the contract. The hold is bounded by the reproducibility (TR-02) and maintenance-budget (TR-54) limits.

**Why this is a TR, not a BR or decision:** BR-64 makes platform evolution the default response to a tenant need; the host-a-capability UX names the hold as the operationalization. The TR forces the hold-and-resume property on the onboarding flow without prescribing how the pending state is represented.

### TR-12: BYO-identity declarations must produce a tenant runtime with no platform-side binding to the platform-provided identity offering {#tr-12}
**Source:** [BR-46]({{< ref "business-requirements.md#br-46" >}}) · [UX: Host a Capability §Constraints Inherited]({{< ref "user-experiences/host-a-capability.md#constraints-inherited" >}})

**Requirement:** A tenant whose declaration is "BYO identity" must be provisionable without the platform-provided identity offering being wired in for end-user authentication. The platform's responsibility for that tenant's identity is limited to network reachability to the chosen external identity service.

**Why this is a TR, not a BR or decision:** BR-46 commits the platform to BYO as a real option. The TR forces the provisioning flow to honor the choice without coupling tenants to the platform-provided service; which external services are reachable is downstream.

### TR-13: Platform-provided identity offering must support a no-recovery credential property {#tr-13}
**Source:** [BR-45]({{< ref "business-requirements.md#br-45" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})

**Requirement:** The platform-provided identity offering must, per tenant electing it, support a configuration where no actor (operator, platform, third-party vendor) can recover a lost end-user credential — no reset email, no recovery code, no admin override. An identity option that cannot honor this configuration is not eligible to be the platform-provided service.

**Why this is a TR, not a BR or decision:** BR-45 forces the property because at least one tenant requires it. The TR is the operative constraint on the offering's surface. Which identity service is selected to satisfy it is an ADR.

### TR-14: All platform administrative interfaces must reject any principal other than the operator {#tr-14}
**Source:** [BR-05]({{< ref "business-requirements.md#br-05" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})

**Requirement:** Every platform-administrative surface — provisioning, deprovisioning, contract change, eviction issuance, secret rotation, drift reconciliation, etc. — must authenticate the caller as the operator's identity (or, when invoked, the sealed successor's). No delegated-administrator role, co-operator role, or shared admin credential exists.

**Why this is a TR, not a BR or decision:** BR-05 makes operator-only operation absolute. The TR closes the surface at the authentication layer rather than restating the rule. Which authentication mechanism is chosen is downstream.

### TR-15: Platform offerings must expose no end-user-addressable surface {#tr-15}
**Source:** [BR-06]({{< ref "business-requirements.md#br-06" >}}) · [UX: Move Off the Platform After Eviction §Constraints Inherited]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#constraints-inherited" >}})

**Requirement:** Platform offerings — observability, secret management, export tool, identity, migration runner, etc. — must expose no UI, API endpoint, or notification channel addressable by tenant end users. Authenticated principals on platform offerings are limited to operator and capability-owner roles; communication to end users about tenant lifecycle (including eviction) is the capability owner's responsibility, not the platform's.

**Why this is a TR, not a BR or decision:** BR-06 forbids any end-user surface on the platform. The TR translates this into the platform's per-offering surface design. What the operator and capability-owner surfaces actually look like is downstream.

### TR-16: Platform must hold a sealed/escrowed successor credential set sufficient to assume full operator authority {#tr-16}
**Source:** [BR-07]({{< ref "business-requirements.md#br-07" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}}) · [UX: Stand Up the Platform §Persona]({{< ref "user-experiences/stand-up-the-platform.md#persona" >}})

**Requirement:** A sealed credential set must exist that, when invoked by the designated successor, grants full operator authority — including running the rebuild flow and exercising every administrative interface covered by TR-14. The seal must be unsealable by the successor without participation from the primary operator. Routine operations must not exercise these credentials.

**Why this is a TR, not a BR or decision:** BR-07 forces successor capability and the seal-vs-routine distinction. The TR forces the credential-set property. The specific seal mechanism (password manager handoff, physical envelope, escrow service) is an ADR.

### TR-17: Each tenant must receive provisioned compute, persistent storage, internal/external network reachability, identity (or BYO binding), backup/DR, and observability — implemented as shared platform offerings {#tr-17}
**Source:** [BR-04]({{< ref "business-requirements.md#br-04" >}}) · [BR-44]({{< ref "business-requirements.md#br-44" >}}) · [Capability §Outputs]({{< ref "_index.md#outputs" >}})

**Requirement:** For every approved tenant, the platform must provision and operate, for the tenant's lifetime, the full inventory: compute, persistent storage, internal and external network reachability, identity binding (platform-provided or per TR-12 BYO), backup with disaster recovery for tenant data, and observability. Each must be implemented as a shared platform offering consumed by every tenant — not duplicated per tenant.

**Why this is a TR, not a BR or decision:** BR-44 lists the inventory; BR-04 demands platform-level investments accrue to every tenant. The TR fixes the inventory and the shared-offering shape. Specific durability levels, network protocols, and backup retention windows are downstream.

### TR-18: Third-party components admissible to the platform must allow control of configuration, data export, and credential revocation/rotation without vendor cooperation {#tr-18}
**Source:** [BR-03]({{< ref "business-requirements.md#br-03" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})

**Requirement:** Any third-party component the platform integrates must allow the operator to (a) read and modify configuration through the platform's tracked-changes surface, (b) export platform-held data in a portable form, and (c) revoke or rotate platform-held credentials without vendor cooperation. Components that fail any of these are not admissible.

**Why this is a TR, not a BR or decision:** BR-03 forbids vendor lock-in that prevents departure. The TR turns "self-hosted" into a per-component admissibility test. Which vendors are chosen is downstream.

### TR-19: All operator/capability-owner engagement must occur on a single durable, append-only, ordered thread per lifecycle event, accessible asynchronously {#tr-19}
**Source:** [BR-15]({{< ref "business-requirements.md#br-15" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** Every operator/capability-owner exchange (onboarding, modify, migration, forced update, contract change, eviction) must occur on exactly one durable engagement thread per event, append-only and ordered, accessible asynchronously to both parties, with the full history preserved. Ephemeral channels (chat, voice, email-only) are not acceptable as the channel of record.

**Why this is a TR, not a BR or decision:** BR-15 demands single-thread, recorded, asynchronous engagement. The TR fixes the channel properties without choosing a tracker.

### TR-20: Engagement channel must distinguish onboarding, modify, migration, forced-update, contract-change, and eviction at the type level {#tr-20}
**Source:** [BR-16]({{< ref "business-requirements.md#br-16" >}}) · [BR-22]({{< ref "business-requirements.md#br-22" >}})

**Requirement:** The engagement channel must support categorization that legibly separates onboarding, modification, data migration, operator-initiated forced update, platform-contract change, and eviction. Issue types must not collapse review scopes; an eviction triggered by a missed forced-update or contract-change deadline must always be a separate, linked issue from the issue that motivated it.

**Why this is a TR, not a BR or decision:** BR-16 demands distinct issue types per review scope; BR-22 forbids re-policing eviction inside the update flow. The TR is the typing constraint on the engagement surface; the type names and tracker semantics are downstream.

### TR-21: Modify-request review must surface only the delta from the tenant's currently-accepted declarations {#tr-21}
**Source:** [BR-17]({{< ref "business-requirements.md#br-17" >}}) · [BR-18]({{< ref "business-requirements.md#br-18" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}})

**Requirement:** The modify flow must support reviewing the proposed delta from the tenant's currently-accepted declarations, without requiring the capability owner to re-accept the platform contract or the operator to re-evaluate the tenant's full prior state.

**Why this is a TR, not a BR or decision:** BR-17 makes the contract evergreen; BR-18 makes modify review delta-only. The TR forces the property on the modify-review surface; how the delta is computed and displayed is downstream.

### TR-22: Forced-update issues must record external pressure name and inherited deadline; one issue per forcing event per affected tenant {#tr-22}
**Source:** [BR-19]({{< ref "business-requirements.md#br-19" >}}) · [BR-20]({{< ref "business-requirements.md#br-20" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** The forced-update issue type must require fields for (a) the external pressure forcing the change (vendor sunset, CVE, EOL) and (b) the deadline inherited from that pressure. When the same tenant is hit by multiple unrelated forcing events at once, each event must produce its own issue, even when remediation overlaps. Forced-update issues must remain open across multiple artifact handoffs and not progress toward eviction until the operative delivery date is missed.

**Why this is a TR, not a BR or decision:** BR-20 demands the two fields and the per-event split; BR-19 forbids early eviction. The TR fixes the issue-type schema and the lifecycle property; what tracker enforces the schema is downstream.

### TR-23: Forced-update flow must record both inherited and any extended operative date, and the extension's external-slack justification, both queryable by the eviction trigger {#tr-23}
**Source:** [BR-21]({{< ref "business-requirements.md#br-21" >}}) · [UX: Operator-Initiated Tenant Update §Journey]({{< ref "user-experiences/operator-initiated-tenant-update.md#journey" >}})

**Requirement:** The forced-update issue must record the original inherited deadline, any negotiated extended operative date, and the external slack that justifies any extension. Both dates must be queryable as inputs to the eviction trigger. Extensions exceeding the named safe slack must be refused; if the external pressure leaves no safe slack, no extension is offered.

**Why this is a TR, not a BR or decision:** BR-21 bounds extensions by external slack. The TR makes the bounds machine-checkable. The shape of the slack record is downstream.

### TR-24: Contract-change rollout must be initiated as a single multi-recipient umbrella issue carrying change, replacement, deadline, reason, and migration guideline {#tr-24}
**Source:** [BR-23]({{< ref "business-requirements.md#br-23" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** The umbrella-issue type must support a single artifact tagging every affected capability owner and carrying (a) what is changing, (b) what it is changing to (or that it is being removed), (c) the deadline, (d) the reason, and (e) the migration guideline where applicable. Per-tenant fanout for the rollout coordination itself is forbidden in this flow.

**Why this is a TR, not a BR or decision:** BR-23 names the umbrella shape and its contents. The TR fixes the issue-type schema; how multi-recipient tagging is implemented is downstream.

### TR-25: Contract-change deadlines must be at least 2× the chosen status-update cadence after filing {#tr-25}
**Source:** [BR-24]({{< ref "business-requirements.md#br-24" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** The umbrella issue must record both the deadline and the operator-chosen status-update cadence. The interval between filing and the deadline must be no less than two full cadence cycles. Combinations of cadence and deadline that violate this must be rejected by the rollout flow.

**Why this is a TR, not a BR or decision:** BR-24 demands at least two status cycles before cutoff. The TR makes the relationship machine-checkable; the cadence values themselves are operator decisions per rollout.

### TR-26: Contract-change deadline must be a single global value; per-tenant overrides are not supported {#tr-26}
**Source:** [BR-25]({{< ref "business-requirements.md#br-25" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** The umbrella issue type must store exactly one deadline applicable uniformly to all tagged tenants. There is no schema for per-tenant deadline overrides; only a global extension covering every tagged tenant may modify the deadline value.

**Why this is a TR, not a BR or decision:** BR-25 forbids per-tenant slips. The TR closes off the schema-level path to one. How global extensions are reflected in-thread is downstream.

### TR-27: Umbrella issues must track per-tenant acknowledgment state; at deadline, the rollout flow must atomically remove the old form, close migrated modify issues, file linked eviction issues per laggard, and close the umbrella {#tr-27}
**Source:** [BR-26]({{< ref "business-requirements.md#br-26" >}}) · [BR-30]({{< ref "business-requirements.md#br-30" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** Each tagged tenant must have an acknowledgment state on the umbrella issue. At the deadline, the rollout flow must atomically (a) remove the old offering from the platform regardless of remaining occupants, (b) close the migrated tenants' modify issues in the normal way, (c) file a separate eviction issue per laggard tenant (including non-acknowledgers) linked to the umbrella, and (d) close the umbrella. No tenant may be silently broken on a removed offering.

**Why this is a TR, not a BR or decision:** BR-26 demands explicit acknowledgment; BR-30 prescribes the deadline closeout shape. The TR consolidates them into the rollout flow's atomic close behavior; how atomicity is achieved is downstream.

### TR-28: Replacement offering must already be a live, hosted offering on the platform before the umbrella issue may be filed {#tr-28}
**Source:** [BR-28]({{< ref "business-requirements.md#br-28" >}}) · [UX: Platform-Contract-Change Rollout §Entry Point]({{< ref "user-experiences/platform-contract-change-rollout.md#entry-point" >}})

**Requirement:** The contract-change flow must refuse to file an umbrella issue when (a) the change replaces an old offering with a new one and (b) the replacement is not yet a live, hosted offering on the platform. Full-removal contract changes (no replacement) are exempt from this gate.

**Why this is a TR, not a BR or decision:** BR-28 makes the precondition absolute. The TR makes it a filing gate; how "live" is verified is downstream.

### TR-29: Platform must support running an old offering and its replacement concurrently for the rollout window when a replacement exists {#tr-29}
**Source:** [BR-27]({{< ref "business-requirements.md#br-27" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** For replacement-style contract changes, the platform must support tenants running on the old offering and the new offering simultaneously throughout the rollout window. The old offering must be removable on the deadline regardless of remaining occupants. Full-removal changes (no replacement) are exempt.

**Why this is a TR, not a BR or decision:** BR-27 commits to concurrent rollout windows. The TR forces the dual-form runtime property; whether concurrency is achieved by side-by-side instances, traffic splitting, or other means is downstream.

### TR-30: Rollout view must produce, on the operator-chosen cadence, both a refreshed in-issue snapshot and a thread comment carrying tenants-on-old, tenants-migrated, open modifies, and time-remaining {#tr-30}
**Source:** [BR-29]({{< ref "business-requirements.md#br-29" >}}) · [UX: Platform-Contract-Change Rollout §Journey]({{< ref "user-experiences/platform-contract-change-rollout.md#journey" >}})

**Requirement:** On the operator-chosen cadence, the contract-change rollout flow must (a) refresh the umbrella issue body with the current snapshot — tenants on the old form, tenants migrated, open modify issues, time remaining — and (b) post a thread comment carrying the same metrics. Both surfaces must be present so a reader landing cold and a watcher tracking history see consistent rollout state.

**Why this is a TR, not a BR or decision:** BR-29 demands both the live snapshot in the issue body and a historical-comment trail. The TR fixes the dual-surface property and the metric set; how the snapshot is computed and rendered is downstream.

### TR-31: Eviction issuance must be operator-only; eviction issues must be locked to their date at filing; required content is exactly date, reason, and link to export-tool documentation {#tr-31}
**Source:** [BR-31]({{< ref "business-requirements.md#br-31" >}}) · [BR-32]({{< ref "business-requirements.md#br-32" >}}) · [BR-33]({{< ref "business-requirements.md#br-33" >}}) · [UX: Move Off the Platform After Eviction §Entry Point]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#entry-point" >}})

**Requirement:** Filing an eviction issue must be restricted to the operator role; capability owners must have no path to initiate eviction. The eviction date must be set at filing and must not be mutable by either party afterward. Required content is exactly (a) eviction date, (b) reason, (c) link to the export-tool documentation; no other field is required.

**Why this is a TR, not a BR or decision:** BR-31, BR-32, and BR-33 together fix the issue's authorship, immutability, and contents. The TR consolidates all three into a single constraint on the eviction-issue schema and authorization gate; the issue-type implementation is downstream.

### TR-32: On the eviction date, tenant compute and network reachability must be deprovisioned and tenant data must transition to a read-only, export-only state {#tr-32}
**Source:** [BR-34]({{< ref "business-requirements.md#br-34" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** On the eviction date, the platform must (a) deprovision the tenant's compute and network reachability so end users can no longer reach it, and (b) transition tenant data to a read-only state in which no actor — including the capability owner, the operator, and tenant components — can write to it, while the export tool continues to function.

**Why this is a TR, not a BR or decision:** BR-34 prescribes the day-zero state transition. The TR is the operative property; whether read-only is enforced via permissions, snapshots, immutable storage, or another mechanism is downstream.

### TR-33: Tenant data must remain readable via the export tool for 30 days post-eviction; on day 30, all tenant data must be permanently deleted across every storage tier the platform controls, with deletion verifiable to the operator {#tr-33}
**Source:** [BR-11]({{< ref "business-requirements.md#br-11" >}}) · [BR-65]({{< ref "business-requirements.md#br-65" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** For 30 days after the eviction date, tenant data must remain accessible only through the export tool. On day 30 the platform must permanently delete the tenant's data across every storage tier it controls — the tenant-accessible export-only copy and any deeper backup-tier copies — with the deletion verifiable to the operator. No residual platform-controlled copy of an evicted tenant's data may survive day 30 in any tier, and no operator-only access path may persist past that point.

**Why this is a TR, not a BR or decision:** BR-11 commits to the tenant-accessible side; BR-65 closes the symmetric question for backup-tier copies. The TR consolidates both into the cross-tier deletion property; how deletion is performed and verified per tier is downstream.

### TR-34: Per-tenant 30-day retention countdown must be operator-pausable, with the pause distinguishing platform-side defects from capability-owner-side issues {#tr-34}
**Source:** [BR-12]({{< ref "business-requirements.md#br-12" >}}) · [UX: Move Off the Platform After Eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#edge-cases" >}})

**Requirement:** The post-eviction 30-day retention countdown must be operator-pausable per tenant. The pause/resume action must record which class triggered it — platform-side defect (pauses the clock) or capability-owner-side issue (does not pause) — and the action must be auditable. Resumption must restart the remaining retention window, not the full 30 days.

**Why this is a TR, not a BR or decision:** BR-12 carves out exactly this pause behavior and allocates accountability. The TR makes the pause a controllable property of the retention-clock surface; the audit-record format is downstream.

### TR-35: Platform must expose a per-tenant export tool callable without operator participation throughout the tenant's hosted lifetime and the post-eviction retention window {#tr-35}
**Source:** [BR-08]({{< ref "business-requirements.md#br-08" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** The platform must expose, per tenant, an export-tool invocation that produces a portable archive of that tenant's data. The invocation must be available throughout the tenant's hosted lifetime and across the 30-day post-eviction retention window without operator participation, and must be re-invocable on demand any number of times. The platform need not retain previously-generated archives between invocations.

**Why this is a TR, not a BR or decision:** BR-08 forces the on-demand, no-operator-needed export property. The TR fixes the invocation surface and re-invokability without prescribing an archive format.

### TR-36: Each export must be accompanied by a platform-produced content checksum/hash and total byte count {#tr-36}
**Source:** [BR-10]({{< ref "business-requirements.md#br-10" >}}) · [UX: Move Off the Platform After Eviction §Journey]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#journey" >}})

**Requirement:** Every export artifact produced by the platform must be paired with (a) a content checksum or hash and (b) a total byte count, both produced by the platform at export time and delivered alongside the artifact. Semantic correctness validation remains the capability owner's responsibility; the platform's verification is bounded to the integrity envelope.

**Why this is a TR, not a BR or decision:** BR-10 forces the verification envelope. The TR fixes the two integrity outputs; which hash function is chosen is an ADR.

### TR-37: Tenant admission must verify that an export-tool path covers every data shape the tenant will introduce {#tr-37}
**Source:** [BR-09]({{< ref "business-requirements.md#br-09" >}}) · [UX: Move Off the Platform After Eviction §Edge Cases]({{< ref "user-experiences/move-off-the-platform-after-eviction.md#edge-cases" >}})

**Requirement:** Tenant admission must verify that an export-tool path covers every data shape the tenant will introduce. A gap in export-tooling coverage must be treated as a platform defect that blocks admission until closed; admission may not proceed on the assumption that the gap can be filled later.

**Why this is a TR, not a BR or decision:** BR-09 forbids gaps in export coverage at eviction time. The TR moves the verification earlier — into admission — so eviction never discovers a gap. How coverage is enumerated and verified is downstream.

### TR-38: Platform must offer a one-shot job runner distinct from long-running tenant components, with progress visible through standard observability {#tr-38}
**Source:** [BR-36]({{< ref "business-requirements.md#br-36" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The platform must offer a job-runner offering — distinct from the long-running tenant component runtime — that executes a packaged artifact end-to-end against a single tenant, exposes progress through the platform's standard observability surfaces, and is bounded in lifetime by a single migration request. The platform runs the job; it does not write, debug, or shepherd it.

**Why this is a TR, not a BR or decision:** BR-36 commits to a one-shot job-runner offering. The TR fixes the offering's separation from long-running runtime and the progress-visibility property; the runner's implementation is downstream.

### TR-39: Migration requests must declare re-run safety and any temporary-spike footprint up front {#tr-39}
**Source:** [BR-38]({{< ref "business-requirements.md#br-38" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** A migration request must, at filing, declare (a) whether the migration process is safe to re-run against an already-populated destination tenant or requires a wiped destination, and (b) any temporary footprint spike beyond the destination tenant's steady-state. Approval is bounded by available platform capacity for the declared spike.

**Why this is a TR, not a BR or decision:** BR-38 names both declarations as part of the operator's review scope. The TR fixes the migration-issue schema; the schema's representation is downstream.

### TR-40: Migration approval must reject any request whose declared peak (steady-state plus spike) exceeds 2× the destination tenant's steady-state in either compute or storage {#tr-40}
**Source:** [BR-39]({{< ref "business-requirements.md#br-39" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The migration review flow must reject — without negotiation — any request where steady-state plus declared spike exceeds 2× the destination tenant's steady-state compute or storage. Resolution requires the capability owner to split the migration, reduce the spike, or resize the tenant first via the modify flow.

**Why this is a TR, not a BR or decision:** BR-39 makes the 2× cap a hard review rule. The TR makes the rule machine-checkable in the review surface; how steady-state is measured is downstream.

### TR-41: Migration runner must support concurrent migrations across distinct tenants without serialization or per-tenant exclusivity {#tr-41}
**Source:** [BR-40]({{< ref "business-requirements.md#br-40" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The migration runner must support multiple migrations running concurrently across different tenants without serializing them or coupling their progress. Tenants must not depend on exclusive use of the runner for their own migration to proceed.

**Why this is a TR, not a BR or decision:** BR-40 commits to concurrent migrations. The TR forces the no-serialization property; how concurrency is implemented (shared infrastructure, per-tenant isolation, queueing) is downstream.

### TR-42: Migration runner must not auto-clean, auto-retry, or auto-progress on job failure; subsequent action must be operator-driven against an explicit capability-owner plan {#tr-42}
**Source:** [BR-41]({{< ref "business-requirements.md#br-41" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** On migration job failure or invalid output, the runner must hold the tenant data in whatever state the failed job left it. The next action — wipe-and-retry, resume, accept partial, abandon — must be operator-driven against a plan the capability owner provides on the issue. The platform must not auto-clean, auto-retry, or auto-prescribe a recovery model.

**Why this is a TR, not a BR or decision:** BR-41 places the recovery decision squarely with the data owner. The TR forbids the runner from acting on its own; the plan-record format is downstream.

### TR-43: Migration runner must deprovision job artifacts on issue closure; re-running requires fresh job creation {#tr-43}
**Source:** [BR-42]({{< ref "business-requirements.md#br-42" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** On migration issue closure (success or abandonment), the runner must remove all per-job artifacts. Subsequent re-runs must require a fresh migration issue and fresh approval; the platform must not retain a migration job past closure.

**Why this is a TR, not a BR or decision:** BR-42 fixes the one-shot lifespan. The TR makes the teardown an obligation of the closure flow; what counts as a "per-job artifact" is bounded by the runner's design.

### TR-44: Platform must offer a secret-management surface populated by capability owners and consumed by their components, with secret values not readable by any non-consuming party {#tr-44}
**Source:** [BR-43]({{< ref "business-requirements.md#br-43" >}}) · [UX: Migrate Existing Data §Journey]({{< ref "user-experiences/migrate-existing-data.md#journey" >}})

**Requirement:** The platform must offer a secret-management surface where capability owners deposit credentials referenced by name from their tenant components and migration processes. Secret values must not appear in engagement-thread comments or in any operator-facing surface, and must not be readable by any party other than the platform components that consume them on the tenant's behalf. Population must be doable by the capability owner without operator involvement.

**Why this is a TR, not a BR or decision:** BR-43 commits to this surface and motivates it as a leak-prevention measure for credentials. The TR makes the secrecy property and capability-owner population first-class. The implementation (key store, secret manager) is an ADR.

### TR-45: Tenant-facing observability must expose, automatically per tenant, the platform-standard health bundle: availability, latency, error rate, resource saturation, and restart/deployment events {#tr-45}
**Source:** [BR-44]({{< ref "business-requirements.md#br-44" >}}) · [BR-53]({{< ref "business-requirements.md#br-53" >}}) · [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** For each live tenant, observability must surface — without capability-owner instrumentation — at minimum: availability, latency, error rate, resource saturation, and restart/deployment events. The bundle must be present from the moment the tenant goes live and must remain present for the tenant's lifetime.

**Why this is a TR, not a BR or decision:** BR-53 defines the bundle's content and the no-tenant-instrumentation property. The TR fixes both. Specific signal definitions, sample rates, and visualization shape are downstream.

### TR-46: Capability owners must be able to mutate their own tenant's alert thresholds without operator participation; cross-tenant threshold mutation is operator-only {#tr-46}
**Source:** [BR-54]({{< ref "business-requirements.md#br-54" >}}) · [BR-57]({{< ref "business-requirements.md#br-57" >}}) · [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** The observability offering must allow each capability owner to mutate alert thresholds for the signals on their own tenant, without operator involvement. Mutation of cross-tenant or platform-wide thresholds must be limited to the operator role.

**Why this is a TR, not a BR or decision:** BR-54 makes thresholds the one self-service surface for capability owners; BR-57 keeps cross-tenant scope operator-only. The TR forces the role-scoped mutation property on the threshold surface.

### TR-47: On threshold crossings, observability must push an alert naming both the signal and the capability {#tr-47}
**Source:** [BR-55]({{< ref "business-requirements.md#br-55" >}}) · [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** When a tenant signal crosses a capability-owner-set threshold, the observability offering must send an alert to the capability owner's registered delivery address. The alert payload must name (a) which signal crossed and (b) which capability is affected. The alert path is best-effort; the pull view is authoritative.

**Why this is a TR, not a BR or decision:** BR-55 commits to threshold-driven push alerts and the content. The TR fixes the property and payload contents; the delivery channel is an ADR.

### TR-48: Tenant view must surface alert-delivery health when degradation is detectable, while remaining the authoritative read of current health {#tr-48}
**Source:** [BR-56]({{< ref "business-requirements.md#br-56" >}}) · [UX: Tenant-Facing Observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md#journey" >}})

**Requirement:** When the observability offering detects that alert delivery to a tenant is degraded, the tenant view must surface that fact so silence on the alert path is not interpreted as evidence of health. The pull view must remain authoritative for current health regardless of alert-path state.

**Why this is a TR, not a BR or decision:** BR-56 demands both the degradation indicator and the pull-authoritative property. The TR fixes them; how degradation is detected is downstream.

### TR-49: Authentication to the observability offering must place a capability owner directly in their tenant scope with no mode-switch broadening it {#tr-49}
**Source:** [BR-57]({{< ref "business-requirements.md#br-57" >}}) · [UX: Tenant-Facing Observability §Entry Point]({{< ref "user-experiences/tenant-facing-observability.md#entry-point" >}})

**Requirement:** A non-operator authenticated session on the observability offering must land in the authenticated capability owner's tenant scope and remain confined to it for the session's lifetime. There must be no UI or API path that broadens scope to another tenant or to a cross-tenant view from a non-operator session.

**Why this is a TR, not a BR or decision:** BR-57 names the isolation property and the operator-only carve-out. The TR forces the session-scope property without choosing an authorization mechanism.

### TR-50: Onboarding closure must produce a working observability login and a wired alert-delivery address as part of provisioning {#tr-50}
**Source:** [BR-58]({{< ref "business-requirements.md#br-58" >}}) · [UX: Host a Capability §Journey]({{< ref "user-experiences/host-a-capability.md#journey" >}}) · [UX: Tenant-Facing Observability §Entry Point]({{< ref "user-experiences/tenant-facing-observability.md#entry-point" >}})

**Requirement:** Closure of an onboarding issue must produce, as part of the same provisioning flow, a working observability login for the capability owner and a wired alert-delivery address. The capability owner must not need to file a separate request to obtain either.

**Why this is a TR, not a BR or decision:** BR-58 demands automatic provisioning of observability access at onboarding. The TR forces the bundle-with-onboarding property; the specific identity and delivery-channel mechanics are downstream.

### TR-51: Onboarding-issue close-out must support a "lost — operator silence" outcome distinct from approved and declined, recorded in-thread {#tr-51}
**Source:** [BR-61]({{< ref "business-requirements.md#br-61" >}}) · [UX: Host a Capability §Edge Cases]({{< ref "user-experiences/host-a-capability.md#edge-cases" >}})

**Requirement:** The onboarding flow must support recording, at issue closure, three distinct terminal outcomes: approved (live tenant), declined (host elsewhere), and lost-to-operator-silence. Each must be a first-class queryable outcome, recorded in-thread, so adoption metrics can distinguish silent-loss from any other failure mode.

**Why this is a TR, not a BR or decision:** BR-61 prescribes the explicit-loss capture. The TR makes the outcome a first-class artifact rather than free-form text. The label and query surface are downstream.

### TR-52: Eviction trigger must be invocable on the basis of either the maintenance-budget condition or the reproducibility-break condition; either alone is sufficient {#tr-52}
**Source:** [BR-60]({{< ref "business-requirements.md#br-60" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})

**Requirement:** The eviction-issuance flow must accept either grounds — projected routine maintenance sustainably exceeding 2× the maintenance-budget KPI, or any required snowflake configuration that cannot be expressed as definitions — as sufficient justification. The recorded grounds must be queryable for later review; both conditions need not be present together.

**Why this is a TR, not a BR or decision:** BR-60 makes either condition independently sufficient. The TR makes the grounds-record property explicit on the eviction surface; how the conditions are measured is downstream.

### TR-53: Platform must produce queryable per-component cost data sufficient for the operator to judge cost-vs-value {#tr-53}
**Source:** [BR-62]({{< ref "business-requirements.md#br-62" >}}) · [Capability §Success Criteria]({{< ref "_index.md#success-criteria" >}})

**Requirement:** The platform must produce queryable per-component cost data on a regular cadence, sufficient for the operator to judge whether continuing operation is worth its bill. There is no fixed numeric target; the operator is the judge. Per-tenant attribution where attributable is desirable but not required by this TR.

**Why this is a TR, not a BR or decision:** BR-62 commits the platform to a cost-judgment surface without naming a target. The TR fixes the queryable-cost property; refresh cadence, granularity, and per-tenant attribution are downstream.

### TR-54: Platform's expected weekly operator-facing work, summed across the currently-hosted tenant set, must be designed to fit within a 2-hour weekly budget {#tr-54}
**Source:** [BR-59]({{< ref "business-requirements.md#br-59" >}}) · [Capability §Success Criteria]({{< ref "_index.md#success-criteria" >}})

**Requirement:** The set of routine operator-facing surfaces (alert handling, status updates, modify reviews, periodic checks) must be designed such that the platform's expected weekly operator work, summed across the currently-hosted tenant set, fits within a 2-hour weekly budget. New surfaces whose costs are not predictable enough to bound this way must not be added without redesign or scope reduction.

**Why this is a TR, not a BR or decision:** BR-59 fixes the maintenance budget. The TR is the design-time obligation that follows: every operator-facing surface is bounded by its share of the 2-hour weekly envelope. How the budgeting is performed is downstream.

### TR-55: ~~Platform-managed resources must use the universal resource identifier standard~~ {#tr-55}
> 🗑️ removed on 2026-04-28 — sourced only to ADR-0006, which was deleted when the repo's existing ADRs were cleared in preparation for the new capability development workflow. Number is reserved and will not be reused.

### TR-56: ~~Platform APIs must use the standard API error response format~~ {#tr-56}
> 🗑️ removed on 2026-04-28 — sourced only to ADR-0007, which was deleted when the repo's existing ADRs were cleared in preparation for the new capability development workflow. Number is reserved and will not be reused.

## Open Questions

Things volunteered as solutions during extraction (parked for the ADR stage), or constraints the capability/UX docs don't yet make explicit.

- **Buy-vs-build decision discipline (BR-63).** [BR-63]({{< ref "business-requirements.md#br-63" >}}) constrains the *decision process*, not the runtime system — it forbids citing operator-skill development as a justification in buy-vs-build trade-offs. It is not surfaced as a TR because there is no runtime obligation it forces; it is parked here for the ADR stage so per-component selection ADRs cite convenience/resiliency/cost evidence.
- **Cost-data refresh cadence and granularity (TR-53).** BR-62 demands cost-vs-value judgment but doesn't quantify "regular," "queryable," or how granular per-component cost must be. Treat as ADR input alongside the observability-offering decisions.
- **Status-update cadence sizing rules (TR-30, TR-25).** BR-29 prescribes a regular cadence "sized to the timeline" and BR-24 imposes the ≥2-cycle deadline rule, but neither fixes a procedure for picking the cadence. Treat as ADR input or per-rollout operator guidance.
- **Last-known-good reference for preflight drift check (TR-05).** The standup UX names "the live platform or the last known-good environment" as the comparison surface but does not specify the form of "last known-good" (snapshot ID, signed manifest, etc.). Treat as ADR input for the drift-detection design.
- **Topology adoption (TR-17, TR-03).** The current repo pattern places an Internet-facing edge layer (with mutual-authentication and traffic-control duties) in front of a private home-lab environment connected to a public-cloud environment through a secure cross-environment tunnel. Whether the platform formally inherits this shape — or selects a different one for the cross-environment foundations — is an ADR decision; this TR doc deliberately does not assume the inherited pattern, and the specific vendors that currently realize each layer are out of scope here.
- **Maintained checklist (UX: Stand Up the Platform).** The standup UX references a "maintained checklist" used during phase validation. Its shape is unspecified; capture as ADR input alongside the rebuild-flow design.
- **Public-cloud account vs. home-lab boundary in TR-18.** Whether a public-cloud account itself counts as a "third-party component" for the BR-03 admissibility test (read/modify config, export data, revoke credentials without vendor cooperation) is ambiguous; the cloud is named in BR-52 as part of foundations and in the capability rules as allowed. Treat as ADR input when picking specific cloud-side components.
