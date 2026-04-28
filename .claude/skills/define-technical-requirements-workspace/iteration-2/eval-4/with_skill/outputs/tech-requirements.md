---
title: "Technical Requirements"
description: >
    Technical requirements derived from the self-hosted-application-platform capability's business requirements (with capability and UX docs as context). Each TR cites the BR-NN it derives from. Decisions belong in ADRs, not here.
type: docs
reviewed_at: null
---

> **Living document.** This is regenerated from `business-requirements.md` (and the capability/UX docs) on demand. Numbering is append-only — once a TR is assigned, it keeps that number forever, even if removed (mark removed ones explicitly). ADRs cite TR-NN, so renumbering would silently break provenance.

> **Review gate.** Set `reviewed_at:` in the frontmatter to today's ISO date once you have read and edited this document. The `plan-adrs` skill will refuse to enumerate decisions until `reviewed_at` is newer than the file's last modification.

**Parent capability:** [self-hosted-application-platform]({{< ref "_index.md" >}})
**Business requirements:** [business-requirements.md]({{< ref "business-requirements.md" >}})

## How to read this

Each TR is **forced** — by a BR (the primary case), by a prior shared ADR, or by a repo-wide constraint. It says what the technical solution must do, not how. Decisions about *how* (which database, which protocol, which library) belong in `adrs/`, not here. If something in this list reads like a chosen solution rather than a constraint, flag it for review. If something has no BR or inherited-constraint source, raise a missing BR back to `extract-business-requirements`.

## Requirements

### TR-01: Tenants must be isolated such that no tenant can read another's state
**Source:** [BR-19]({{< ref "business-requirements.md#br-19" >}}) · [BR-04]({{< ref "business-requirements.md#br-04" >}}) · [Capability §Business Rules]({{< ref "_index.md#business-rules" >}})

**Requirement:** The platform must enforce strict tenant isolation at the data, compute, network, and telemetry layers. No tenant workload, capability owner, or end user may observe or access another tenant's state, secrets, traffic, or telemetry under any normal or degraded operating condition. There must be no mode-switch by which a non-operator principal can widen scope beyond their own tenant.

**Why this is a TR, not a BR or decision:** The BR demands cross-tenant invisibility as a tenant-trust property; the TR translates that into the surfaces (data, compute, network, telemetry) on which isolation must hold and the no-scope-widening property the implementation must enforce.

### TR-02: Operators must be able to roll out a platform-contract change without breaking existing tenants
**Source:** [BR-09]({{< ref "business-requirements.md#br-09" >}}) · [BR-45]({{< ref "business-requirements.md#br-45" >}}) · [UX: platform-contract-change-rollout]({{< ref "user-experiences/platform-contract-change-rollout.md" >}})

**Requirement:** When the platform publishes a new contract version, existing tenants must continue to operate against the prior contract version until they migrate. The platform must support N concurrent contract versions for a bounded migration window with version pinning per tenant, with the carve-out that a full offering removal (no replacement) is all-or-nothing at the deadline.

**Why this is a TR, not a BR or decision:** The BR demands no-surprise migration; the TR translates that into concurrent multi-version support with per-tenant pinning over a bounded window — without choosing a versioning scheme.

### TR-03: Tenant-facing observability data must be queryable per-tenant within their data scope only
**Source:** [BR-18]({{< ref "business-requirements.md#br-18" >}}) · [BR-19]({{< ref "business-requirements.md#br-19" >}}) · [UX: tenant-facing-observability §Journey]({{< ref "user-experiences/tenant-facing-observability.md" >}})

**Requirement:** Capability owners must be able to query metrics, logs, and traces for their own tenant on demand, scoped strictly to their own tenant for the entire session. Cross-tenant observability data must be inaccessible to a capability owner. The pull surface must remain available even when push delivery is degraded.

**Why this is a TR, not a BR or decision:** The BRs demand both an on-demand pull surface for capability owners and strict per-tenant scope; the TR translates that into the queryability and scope-enforcement constraints the implementation must hold.

### TR-04: Operator-initiated tenant updates must complete without tenant-perceived downtime for online workloads
**Source:** [BR-09]({{< ref "business-requirements.md#br-09" >}}) · [UX: operator-initiated-tenant-update §Success]({{< ref "user-experiences/operator-initiated-tenant-update.md" >}})

**Requirement:** When the operator initiates an update to a tenant (config, version, or capability), tenants serving online traffic must observe no end-user-visible downtime during the update.

**Why this is a TR, not a BR or decision:** The UX success criterion frames zero downtime as the user-perceived outcome forced by the no-surprise migration BR; the TR translates it into the platform's behavior during update execution.

### TR-05: A tenant evicted from the platform must be able to take their data with them
**Source:** [BR-27]({{< ref "business-requirements.md#br-27" >}}) · [BR-29]({{< ref "business-requirements.md#br-29" >}}) · [UX: move-off-the-platform-after-eviction]({{< ref "user-experiences/move-off-the-platform-after-eviction.md" >}})

**Requirement:** The platform must provide an export mechanism by which an evicted tenant can retrieve all of their data in a portable format within a defined export window (the BR-29 retention window). Export tooling must be available for every kind of data the platform hosts and present from the moment a tenant is live, not assembled at eviction.

**Why this is a TR, not a BR or decision:** The BRs demand universal export availability and a fixed post-eviction retention window; the TR translates them into the property that an evicted tenant can complete extraction within the window using already-present tooling.

### TR-06: New tenants must be able to migrate existing data into the platform without loss or corruption
> ⚠️ source no longer resolves — human review

**Source:** [BR-24]({{< ref "business-requirements.md#br-24" >}}) · [UX: migrate-existing-data]({{< ref "user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists" >}})

**Requirement:** The platform must accept a tenant's pre-existing data and import it idempotently with verifiable integrity (no silent loss, no duplication on retry).

**Why this is a TR, not a BR or decision:** UX requires lossless, retry-safe migration as part of the journey. The original UX section anchor `#a-section-that-no-longer-exists` no longer resolves and must be re-sourced to the current `migrate-existing-data` UX section that frames this constraint; the BR-24 primary citation remains valid.

### TR-07: All inter-service communication must traverse the Cloudflare → GCP path
**Source:** CLAUDE.md §Architecture overview (repo-pattern, inherited) · prior shared decision

**Requirement:** Network traffic between platform services and tenant workloads, and between platform services themselves, must conform to the existing Cloudflare-fronted, GCP-hosted topology with WireGuard back to home lab.

**Why this is a TR, not a BR or decision:** Inherited topology constraint from the repo's architecture; not subject to revisiting at the capability level.

### TR-08: The platform must be rebuildable end-to-end from the definitions repo plus root infrastructure access alone
**Source:** [BR-02]({{< ref "business-requirements.md#br-02" >}}) · [BR-23]({{< ref "business-requirements.md#br-23" >}}) · [UX: stand-up-the-platform]({{< ref "user-experiences/stand-up-the-platform.md" >}})

**Requirement:** From a clean substrate, the platform must reach a tenant-ready state using only artifacts in the definitions repo plus root credentials to the underlying infrastructure. No bespoke per-rebuild or per-tenant manual configuration may be required at any step. Any state that cannot be expressed as a definition is a reproducibility defect and must be fixed at the platform level rather than absorbed.

**Why this is a TR, not a BR or decision:** The BR demands clean-rebuild-from-definitions as a continuity property; the TR translates that into the closure constraint on rebuild inputs and the no-snowflake property on every step.

### TR-09: Successor takeover must be possible from sealed credentials and the definitions repo with no operator-specific local state
**Source:** [BR-05]({{< ref "business-requirements.md#br-05" >}}) · [BR-42]({{< ref "business-requirements.md#br-42" >}})

**Requirement:** A designated successor in possession of the sealed/escrowed credentials and the definitions repo must be able to operate the platform — including running the standup flow identically to the primary operator — without any state, secret, or context that lives only on the primary operator's personal machine or accounts. Credential handoff must be a discrete event, not a steady-state shared-administration mode.

**Why this is a TR, not a BR or decision:** The BRs demand convergence of successor and primary on the same operating surface; the TR translates that into the no-personal-state constraint on the platform's operability.

### TR-10: Each hosted tenant must be provisioned with isolated compute, durable storage, and network reachability (internal and external)
**Source:** [BR-13]({{< ref "business-requirements.md#br-13" >}})

**Requirement:** For every hosted tenant, the platform must provision a place for the tenant's workload to run, durable storage for the tenant's data, and both internal (between tenants where declared) and external (reachable by end users) network reachability — all subject to the TR-01 isolation property.

**Why this is a TR, not a BR or decision:** The BR enumerates these as direct deliverables to every tenant; the TR translates them into the per-tenant provisioning surfaces the platform must produce — without choosing the form (containers vs VMs, block vs object).

### TR-11: The platform-provided identity service must support an unrecoverable-credentials mode
**Source:** [BR-14]({{< ref "business-requirements.md#br-14" >}}) · [BR-15]({{< ref "business-requirements.md#br-15" >}})

**Requirement:** The platform must offer an identity and authentication service for tenant end users. Any identity option presented as the platform-provided service must be capable of honoring a "lost credentials cannot be recovered" property for tenants that require it. Tenants may opt to bring their own identity, and that election must be a declared field of the tenant's contract acceptance.

**Why this is a TR, not a BR or decision:** The BRs demand both the offering and the unrecoverable-credentials property as eligibility criteria; the TR translates them into a capability constraint on whatever identity software is later chosen.

### TR-12: The platform must back up tenant data and provide a disaster-recovery path defined as part of the platform contract
**Source:** [BR-16]({{< ref "business-requirements.md#br-16" >}})

**Requirement:** For every hosted tenant, the platform must back up the tenant's data and provide a disaster-recovery path. The platform must publish the backup standard (RPO/RTO, scope, cadence) as a term of its contract so tenants accept it explicitly at onboarding.

**Why this is a TR, not a BR or decision:** The BR forces both the function and its publication as a contract term; the TR translates that into a published-standard constraint without picking specific RPO/RTO numbers (those are KPI/decision territory and currently a documented BR gap).

### TR-13: The platform must produce per-tenant operator-side health observability without requiring tenant instrumentation
**Source:** [BR-17]({{< ref "business-requirements.md#br-17" >}})

**Requirement:** The platform must produce a per-tenant up/healthy signal accessible to the operator for every hosted tenant, derived from platform-side observation rather than from instrumentation the tenant must add.

**Why this is a TR, not a BR or decision:** The BR demands an operator-side view that does not depend on tenant cooperation; the TR translates that into the platform-sourced signal constraint without naming a tooling stack.

### TR-14: Capability owners must be able to set their own alert thresholds within their tenant scope without operator involvement
**Source:** [BR-20]({{< ref "business-requirements.md#br-20" >}}) · [BR-19]({{< ref "business-requirements.md#br-19" >}})

**Requirement:** The platform must expose a self-service surface for capability owners to create, modify, and remove alert thresholds against signals scoped to their own tenant. The surface must not allow widening scope beyond the tenant (TR-01) and must not require operator action to take effect.

**Why this is a TR, not a BR or decision:** The BR carves out alert-threshold self-service as the one tenant-facing self-service surface; the TR translates that into a scoped, no-operator-loop write surface.

### TR-15: Degradation of push alert delivery for a tenant must be reflected in that tenant's pull observability view
**Source:** [BR-21]({{< ref "business-requirements.md#br-21" >}})

**Requirement:** When the platform detects that alert delivery is failing for a tenant, the tenant-scoped observability view must surface a clear "alerting degraded" indicator so capability owners do not interpret silence as health. The pull view must remain the source of truth.

**Why this is a TR, not a BR or decision:** The BR forces the trust contract between push and pull; the TR translates that into a detect-and-render constraint on the tenant view, leaving detection mechanism and rendering form open.

### TR-16: All capability-owner ↔ operator engagement must occur via typed GitHub issues against the infra repo
**Source:** [BR-22]({{< ref "business-requirements.md#br-22" >}})

**Requirement:** The platform must accept and produce its capability-owner-facing work only through GitHub issues filed against the infra repo, and each workflow (`onboard my capability`, `modify my capability`, `migrate my data`, `platform update required`, `platform contract change`, eviction) must use a distinct issue type so the operator's review scope and the journey shape are unambiguous to all parties.

**Why this is a TR, not a BR or decision:** The BR forces issues-as-the-engagement-channel and the per-workflow distinct-type property; the TR translates that into the constraint on the platform's intake/output surface for tenant work.

### TR-17: Tenant onboarding must execute against the platform's existing definitions with no per-tenant manual configuration
**Source:** [BR-23]({{< ref "business-requirements.md#br-23" >}}) · [BR-02]({{< ref "business-requirements.md#br-02" >}})

**Requirement:** The provisioning path that brings a new tenant live must consume only the platform's definitions; no operator-applied bespoke configuration may be required to fit a particular tenant. If onboarding requires bespoke manual config, that condition is a platform-level reproducibility defect, not a tenant accommodation.

**Why this is a TR, not a BR or decision:** The BR demands that tenants do not introduce snowflake config; the TR translates that into the constraint that the onboarding execution path is closed over the definitions repo.

### TR-18: The platform must run capability-owner-supplied one-shot migration jobs through a managed lifecycle (provision → run → observe → teardown)
**Source:** [BR-24]({{< ref "business-requirements.md#br-24" >}}) · [BR-26]({{< ref "business-requirements.md#br-26" >}})

**Requirement:** The platform must provide a generic facility for executing a capability-owner-packaged one-shot migration job. The job must be provisioned on approval, run against the destination tenant, be observable while running, and be torn down on completion. The platform must enforce, at review time, a bounded multiple cap on the migration's peak temporary compute and storage footprint relative to the destination tenant's steady-state footprint.

**Why this is a TR, not a BR or decision:** The BRs demand both the offering's lifecycle shape and the footprint cap; the TR translates them into the lifecycle and admission-control constraints — leaving the runner technology and the exact multiple as decisions/KPIs.

### TR-19: The platform must offer a named-secret store referenced from packaged artifacts so secret material never appears in artifacts or issue threads
**Source:** [BR-25]({{< ref "business-requirements.md#br-25" >}})

**Requirement:** The platform must expose a secret-management facility for capability owners to register credentials their workloads need. Packaged artifacts and issues must reference these secrets by name only; the secret values themselves must not appear in any artifact or issue thread the platform handles.

**Why this is a TR, not a BR or decision:** The BR forces a confidentiality-shaped offering with a name-only reference contract; the TR translates that into the indirection constraint between artifact/issue surfaces and secret material — leaving the secret-store technology as a decision.

### TR-20: Each export archive the platform produces must be accompanied by a platform-computed checksum/hash and total byte size
**Source:** [BR-28]({{< ref "business-requirements.md#br-28" >}})

**Requirement:** Whenever the platform produces an export archive for a tenant, the platform must publish alongside it a checksum/hash and total size in bytes computed by the platform. This is the bytes-level integrity ceiling the platform commits to; semantic validation remains the capability owner's responsibility.

**Why this is a TR, not a BR or decision:** The BR draws the responsibility line at platform-verified bytes-level integrity; the TR translates that into a per-export emission constraint without picking a hash algorithm.

### TR-21: After eviction, the platform must enforce the retention window as a hard wall, with the only carve-out being platform-rooted export-defect pauses
**Source:** [BR-29]({{< ref "business-requirements.md#br-29" >}}) · [BR-30]({{< ref "business-requirements.md#br-30" >}}) · [BR-31]({{< ref "business-requirements.md#br-31" >}})

**Requirement:** From the eviction date the platform must (a) deprovision the tenant's compute, network, and live-serving resources, (b) transition the tenant's data into an export-only, read-only state, and (c) keep that state available for a fixed retention window (the BR-29 number) during which export tooling continues to work. After the window, no tenant-accessible copy must remain. The retention countdown may be paused only when a failure to produce a complete, valid export is rooted in the platform's tooling or hosting; capability-owner-rooted failures must not pause the countdown.

**Why this is a TR, not a BR or decision:** The BRs demand cutover, the read-only freeze, the fixed window, and the narrow defect-pause carve-out as load-bearing properties; the TR translates them into the post-eviction state-machine constraint the platform must enforce.

### TR-22: A tenant's exit must end with a verifiable clean-exit state and preserve eligibility for future re-onboarding
**Source:** [BR-54]({{< ref "business-requirements.md#br-54" >}})

**Requirement:** When the eviction journey ends, the platform must record a closed eviction issue with the cutover confirmation and (post-retention-window) the absence of any tenant-accessible copy of the tenant's data. The eviction state must not preclude a future `host-a-capability` re-onboarding by the same capability owner.

**Why this is a TR, not a BR or decision:** The BR forces the clean-exit and re-eligibility properties; the TR translates them into the post-journey record state and the absence of any latent re-onboarding block.

### TR-23: Cross-environment connectivity between cloud and home-lab parts of the platform must be provisioned as part of the standup foundation
**Source:** [BR-35]({{< ref "business-requirements.md#br-35" >}}) · [BR-36]({{< ref "business-requirements.md#br-36" >}})

**Requirement:** Where the platform spans operator-controlled public cloud and operator-owned private infrastructure, the connectivity between them must be brought up as part of the platform's foundation phase, from definitions, with no manual cross-environment wiring step required for a clean rebuild.

**Why this is a TR, not a BR or decision:** The BRs demand both the allowed multi-environment shape and that connectivity be foundational; the TR translates them into the rebuild-closure constraint on cross-environment links — leaving the transport (Wireguard, peering, etc.) as a decision.

### TR-24: Standup must execute as discrete automated phases with operator-validated checkpoints between them, and no phase may resume from partial state on failure
**Source:** [BR-38]({{< ref "business-requirements.md#br-38" >}}) · [BR-39]({{< ref "business-requirements.md#br-39" >}})

**Requirement:** The standup flow must be partitioned into discrete automated phases (foundations, core services, cross-cutting services, canary), pausing between phases for an explicit operator `continue` signal. On phase failure, the operator must be able to tear down everything provisioned so far and restart from the top after a definition fix; partial state must never be accepted as a starting point.

**Why this is a TR, not a BR or decision:** The BRs demand the phase + checkpoint structure and the no-partial-state property as load-bearing for confidence-over-speed; the TR translates them into the standup orchestration constraints — leaving phase contents and teardown mechanics as decisions.

### TR-25: Standup must not declare the platform ready until a purpose-built canary tenant has been deployed, exercised end-to-end, and torn down
**Source:** [BR-37]({{< ref "business-requirements.md#br-37" >}})

**Requirement:** The platform must withhold its ready signal until a canary tenant — maintained alongside the platform definitions — has been deployed, exercised across run, reachability, storage read/write, identity authentication, backup pickup, and observability pickup, and then torn down cleanly. Infrastructure self-checks alone must not suffice. The canary gate must not be bypassable under time pressure.

**Why this is a TR, not a BR or decision:** The BR demands a green-canary readiness gate with explicit no-bypass discipline; the TR translates that into the readiness-signal constraint on the standup flow without picking the canary's contents or framework.

### TR-26: Drift between live platform state and definitions must be detected and resolved before any rebuild begins, and immutability discipline must hold across all UXs that can introduce platform state
**Source:** [BR-40]({{< ref "business-requirements.md#br-40" >}})

**Requirement:** Whenever prior platform state exists, a preflight drift check against the live platform or last known-good environment must pass before standup begins. The continuous machinery that prevents drift between rebuilds — tracked changes and immutability of platform state — must be enforced uniformly by every UX that can introduce platform state.

**Why this is a TR, not a BR or decision:** The BR demands both the preflight drift gate and cross-UX immutability; the TR translates that into the precondition constraint on rebuild and the cross-UX invariant — leaving the drift-computation mechanism as a decision/open question.

### TR-27: A scratch-infrastructure rebuild drill must be performed after every significant platform change and at least quarterly using the same flow as a real rebuild
**Source:** [BR-41]({{< ref "business-requirements.md#br-41" >}})

**Requirement:** The platform must support a scratch-infrastructure rebuild drill executed identically to a real rebuild (same definitions, same flow), runnable in parallel without disturbing live state, and conducted on the cadence the BR requires (after every significant change, no less than quarterly).

**Why this is a TR, not a BR or decision:** The BR demands drill-as-proof of the reproducibility KPI; the TR translates that into the constraint that the rebuild flow is parallel-executable on scratch infrastructure with full fidelity.

### TR-28: A platform-update-required issue must carry the externally-inherited deadline and external reason, and follow the typed-issue contract
**Source:** [BR-50]({{< ref "business-requirements.md#br-50" >}}) · [BR-51]({{< ref "business-requirements.md#br-51" >}}) · [BR-22]({{< ref "business-requirements.md#br-22" >}})

**Requirement:** When an external dependency event (vendor sunset, CVE, EOL) forces a tenant update, the platform's intake must produce a `platform update required` issue per affected tenant whose deadline is the inherited external deadline and whose external reason is recorded on the issue. Negotiated extensions must be recorded on the issue and bounded by the safe slack the external pressure allows; where no safe slack exists, no extension may be recorded.

**Why this is a TR, not a BR or decision:** The BRs demand the inherited-deadline contract, the typed signal, and the safe-slack bound on negotiation; the TR translates them into the issue-shape and admission-control constraints — leaving the CVE/EOL feed source as operational detail.

### TR-29: A contract-change rollout must publish status updates on a regular operator-chosen cadence in the umbrella thread, with the current snapshot in the issue body and each scheduled update as a comment
**Source:** [BR-46]({{< ref "business-requirements.md#br-46" >}}) · [BR-47]({{< ref "business-requirements.md#br-47" >}}) · [BR-48]({{< ref "business-requirements.md#br-48" >}}) · [BR-49]({{< ref "business-requirements.md#br-49" >}})

**Requirement:** The platform must support a contract-change umbrella issue model in which (a) the operator-chosen deadline gives every affected tenant at least two full status-update cycles before cutoff, (b) the umbrella body holds the current snapshot of remaining-on-old, migrated, open `modify` issues, and time remaining, with each scheduled update also posted as a comment, (c) each tagged capability owner's acknowledgment is required and the absence of acknowledgment by deadline is treated as non-engagement, and (d) the deadline applies uniformly to all affected tenants — only global extensions are allowed.

**Why this is a TR, not a BR or decision:** The BRs collectively force the umbrella status-visibility, acknowledgment, and deadline-uniformity properties; the TR translates them into the umbrella-issue contract the platform must enforce — leaving cadence values and tooling form as decisions.

### TR-30: An eviction triggered by a missed operative delivery date or by laggard status in a contract rollout must be filed as a separate, linked eviction issue
**Source:** [BR-52]({{< ref "business-requirements.md#br-52" >}}) · [BR-22]({{< ref "business-requirements.md#br-22" >}})

**Requirement:** When a tenant misses the operative delivery date on a `platform update required` issue or fails to migrate by a contract-change deadline, the platform's intake must produce a separate eviction issue linking back to the originating issue, and the originating issue must be closed as superseded by eviction. Eviction state must not be carried as a flag on the originating issue.

**Why this is a TR, not a BR or decision:** The BRs demand a separate-and-linked eviction record across both update and contract-change journeys; the TR translates that into the issue-shape and lifecycle constraint on eviction.

## Open Questions

Things the user volunteered as solutions during extraction (parked for the ADR stage), or constraints the capability/UX docs don't yet make explicit.

- **TR-06 source link is stale.** The original UX section anchor (`migrate-existing-data.md#a-section-that-no-longer-exists`) no longer resolves; the BR-24 primary citation remains valid. Human review needed to re-source the UX context to a current section anchor (and to add an explicit `{#anchor-id}` annotation on the target heading per skill guidance).
- **Whether tenant data export (TR-05) should be on-demand or continuously-available** — captured during extraction from the move-off UX. (Carried forward.)
- **Whether contract versioning (TR-02) requires semver or a different versioning scheme** — surfaced but deferred to the ADR stage. (Carried forward.)
- **Numeric thresholds embedded in BRs but not yet in TRs:** the BR-26 migration footprint multiple (currently `2x`), BR-29 retention window (currently `30 days`), BR-32 maintenance budget (currently `2 hr/week`), and BR-46 status-update cycles (currently `at least two cycles`). These belong with TR/KPI translation; the TRs above intentionally cite the BR-defined thresholds rather than re-stating numbers, since the values currently live in the capability/UX docs.
- **Backup standard (BR-16 / TR-12) is undefined.** TR-12 forces publication of an RPO/RTO/scope/cadence standard, but the standard itself is not yet authored. This is a BR-shaped gap to resolve before TR-12 can be measurable.
- **Specific signal bundle for tenant-facing observability (BR-18 / TR-03 / TR-13).** The platform-standard health bundle (availability, latency, error rate, resource saturation, restart/deployment events) is named in the UX but its exact contents and any platform-default thresholds are open and belong with TR/ADR work.
- **Migration-process concurrency model (BR-24 / TR-18).** Concurrent migrations are promised across tenants but capacity sizing and queueing are unbounded; belongs with ADR work.
- **Drift detection mechanism (BR-40 / TR-26).** What counts as "last known-good environment", how drift is computed, where the policy is enforced — all open and belong with ADR work.
- **Deeper backup-tier policy after the 30-day retention window (BR-29 / BR-54 / TR-21 / TR-22).** Retention duration, deletion behavior, and operator-access/privacy constraints for any backup-tier copies after tenant-accessible data is removed are still TBD upstream in the BRs; the TRs above only cover the tenant-accessible window.
- **Explicit anchors on capability and UX section headings.** Source links currently anchor at the page level for capability and UX references. Per skill guidance, section deep-links require explicit `{#anchor-id}` annotations on the target heading. Once those anchors are added in the source docs, TRs above should be updated to deep-link to the relevant sections (and the BR anchors `#br-NN` already used here should be confirmed present in `business-requirements.md`).
- **No formal tenant-facing pending-update view exists today.** If an earlier deprecation/pending-update signal is added later, it would extend BR-18 / TR-03 rather than the operator-initiated-tenant-update journey. Parked for tenant-facing-observability evolution.
