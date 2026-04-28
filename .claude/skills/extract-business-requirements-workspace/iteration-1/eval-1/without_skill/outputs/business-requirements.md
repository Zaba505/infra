# Business Requirements — Self-Hosted Application Platform

This document extracts the business-level requirements implied by the [Self-Hosted Application Platform](../../../../../../docs/content/capabilities/self-hosted-application-platform/_index.md) capability and its user experiences. Each requirement is sourced — meaning every BR-NN below points back to a specific section of the capability doc or a UX doc, so the lineage stays traceable.

Requirements are stated at the business level — *what* the platform must do or guarantee, not *how*. Technical requirements (TRs) and architectural decisions (ADRs) follow from these in downstream documents.

Sources referenced:

- `_index.md` — the capability doc (purpose, stakeholders, triggers, outputs, business rules, KPIs, out-of-scope).
- `user-experiences/host-a-capability.md` (UX-1)
- `user-experiences/migrate-existing-data.md` (UX-2)
- `user-experiences/move-off-the-platform-after-eviction.md` (UX-3)
- `user-experiences/operator-initiated-tenant-update.md` (UX-4)
- `user-experiences/platform-contract-change-rollout.md` (UX-5)
- `user-experiences/stand-up-the-platform.md` (UX-6)
- `user-experiences/tenant-facing-observability.md` (UX-7)

---

## A. Hosting Offerings (what the platform must provide to tenants)

**BR-01. The platform shall provide compute on which a tenant capability can run.**
Source: `_index.md` § Outputs & Deliverables ("Compute — a place for the application to run").

**BR-02. The platform shall provide persistent, durable storage for tenant data.**
Source: `_index.md` § Outputs & Deliverables ("Persistent storage").

**BR-03. The platform shall provide network reachability — both internal (between tenants) and external (reachable by a tenant's end users).**
Source: `_index.md` § Outputs & Deliverables ("Network reachability").

**BR-04. The platform shall offer an identity & authentication service for tenant end users.**
Source: `_index.md` § Outputs & Deliverables; § Business Rules ("Identity service honors tenant credential-recovery rules").

**BR-05. Any platform-provided identity option shall be capable of honoring a "lost credentials cannot be recovered" property (Signal-style).**
Source: `_index.md` § Business Rules — *Identity service honors tenant credential-recovery rules*.

**BR-06. Tenants shall be permitted to opt out of platform-provided identity and bring their own.**
Source: `_index.md` § Outputs & Deliverables; § Triggers & Inputs ("either use of the platform-provided identity service, or a declared decision to bring their own").

**BR-07. The platform shall provide backup and disaster recovery for tenant data, to a standard the platform defines.**
Source: `_index.md` § Outputs & Deliverables ("Backup and disaster recovery").

**BR-08. The platform shall provide observability sufficient for the operator to know whether each tenant is up and healthy without the tenant having to instrument that itself.**
Source: `_index.md` § Outputs & Deliverables ("Observability").

**BR-09. The platform shall provide a secret-management offering that tenants can register secrets with by name.**
Source: UX-2 § Step 1 ("registers any credentials their migration process needs … with the platform's secret-management offering").

**BR-10. The platform shall provide a one-shot job runner ("migration-process offering") capable of running tenant-supplied migration processes with the platform's standard observability.**
Source: UX-2 § Constraints Inherited (*"a platform-provided one-shot-job runner with the platform's standard observability"*); § Step 4.

**BR-11. The platform shall provide an export tool, available at all times for every tenant, that produces a downloadable archive of the tenant's data plus checksum/hash and total size in bytes.**
Source: UX-3 § Phase A Step 3, § Edge Cases ("Export tooling does not exist for this tenant's data shape … cannot happen by design — export tooling is a core platform feature").

---

## B. Engagement Surface (how capability owners interact with the platform)

**BR-12. The only channel for a capability owner to engage the platform shall be GitHub issues against the infra repo. There shall be no self-service portal and no other front door.**
Source: UX-1 § Step 1; `_index.md` § Out of Scope ("Multi-operator administration, role delegation, or self-service onboarding").

**BR-13. The platform shall support distinct issue types for distinct engagement modes, at minimum:**
- `onboard my capability` (UX-1)
- `modify my capability` (UX-1 step 8)
- `migrate my data` (UX-2)
- `platform update required` (UX-4) — operator-initiated
- `platform contract change` (UX-5) — operator-initiated, umbrella
- *eviction issue* (UX-3, UX-4 step 5, UX-5 step 4) — operator-initiated

The distinction between issue types must be meaningful to capability owners because the operator's review scope and the lifecycle differ per type.
Source: UX-1 § Step 8; UX-2 § Step 2; UX-4 § Step 1; UX-5 § Step 1.

**BR-14. There shall be no formal response-time SLA on issues in either direction (operator → capability owner or capability owner → operator).**
Source: UX-1 § Step 1 ("There is no response-time guarantee — this is personal-scale, async by default"); UX-4 § Edge Cases.

**BR-15. The platform's contract with tenants shall be evergreen: capability owners shall not be required to re-accept the contract on each modification.**
Source: UX-1 § Step 8 ("The platform contract is evergreen — the capability owner does not re-accept it on each modification").

**BR-16. Contract acceptance shall be implicit in the tenant's tech-design submission (which declares resource needs, identity choice, packaging form, and availability expectations); there shall be no separate explicit contract-acceptance gate.**
Source: UX-1 § Step 3a; § Constraints Inherited.

---

## C. Onboarding a Tenant (UX-1)

**BR-17. The platform shall accept tenants only after operator approval; no self-onboarding by tenants is permitted.**
Source: `_index.md` § Triggers & Inputs ("The operator has authorized the capability to run on the platform").

**BR-18. The operator's onboarding review scope shall be narrow: (a) does each platform-hosted component align with an existing platform offering, and (b) are any components requesting a new platform offering.**
Source: UX-1 § Step 2.

**BR-19. The platform shall support three onboarding outcomes: approved-as-is, new-offering-needed (waited on with no timeline guarantee), or declined-host-elsewhere.**
Source: UX-1 § Step 3.

**BR-20. The default response to a tenant needing something the platform does not yet provide shall be to consider expanding the platform, not to refuse the tenant.**
Source: `_index.md` § Business Rules ("The capability evolves with its tenants"); UX-1 § Constraints Inherited.

**BR-21. A new-offering request shall be declinable when the resulting ongoing scope cannot be kept reproducible within the *Reproducibility* KPI or operated within the *Operator maintenance budget* KPI, even if the offering is technically buildable.**
Source: UX-1 § Step 3c; § Constraints Inherited.

**BR-22. Tenants shall hand off their components as packaged artifacts in the form the platform accepts; they shall not hand over raw source for the operator to package.**
Source: UX-1 § Step 4.

**BR-23. The platform shall require the capability owner to validate the deployed tenant via a test before the onboarding issue closes; the platform does not prescribe a test plan.**
Source: UX-1 § Step 6.

**BR-24. Modify-request reviews shall cover only the delta, not a full re-evaluation of the tenant.**
Source: UX-1 § Step 8.

---

## D. Migrating Existing Data (UX-2)

**BR-25. The platform shall support migrating a capability owner's pre-existing data into a newly-provisioned tenant by running a tenant-supplied, one-shot migration process; the platform provides the runner, not the migration logic.**
Source: UX-2 § Persona, § Goal, § Out of Scope ("Writing or debugging the capability owner's migration process").

**BR-26. Each migration job shall have its peak temporary footprint capped at no more than 2x the destination tenant's steady-state compute and storage footprint; the operator shall reject migration requests that would exceed this cap.**
Source: UX-2 § Step 3.

**BR-27. Migration processes shall be packaged in the same form the platform accepts for any tenant component; the contract shall not relax for migration.**
Source: UX-2 § Constraints Inherited.

**BR-28. Each migration request shall declare its re-run contract — whether the process is safe against an already-populated destination or requires a wiped/empty destination per run.**
Source: UX-2 § Step 2.

**BR-29. Migration jobs shall be observable to the capability owner through the platform's standard observability surface.**
Source: UX-2 § Step 5.

**BR-30. The platform shall support concurrent migrations across different tenants without offering exclusive capacity or a completion-time guarantee to any one tenant.**
Source: UX-2 § Step 4; § Constraints Inherited.

**BR-31. The platform shall not auto-clean partial migration state; recovery decisions (wipe-and-retry, resume, accept partial, abandon) shall be the capability owner's to make.**
Source: UX-2 § Step 7b, § Edge Cases.

**BR-32. A migration job shall be torn down on completion; re-running later shall require a fresh `migrate my data` issue.**
Source: UX-2 § Step 8.

---

## E. Moving Off After Eviction (UX-3)

**BR-33. Eviction shall be operator-initiated only and shall be communicated to the capability owner via an eviction issue containing the eviction date, the reason, and a link to the export tooling.**
Source: UX-3 § Entry Point.

**BR-34. The platform shall guarantee a 30-day post-eviction retention window during which the capability owner can continue to pull exports from a frozen, read-only snapshot.**
Source: UX-3 § Phase B Step 5, § Phase C Step 7.

**BR-35. After 30 days post-eviction, the platform shall stop offering any tenant-accessible copy of the evicted tenant's data, regardless of whether the capability owner closed the loop.**
Source: UX-3 § Phase C Step 7.

**BR-36. The 30-day clock shall be hard except where the failure is shown to be in the platform's export tooling or data hosting; in that case the operator shall pause the retention-window countdown until a clean export can be produced.**
Source: UX-3 § Edge Cases.

**BR-37. Validation of export completeness and correctness shall be the capability owner's responsibility, not the platform's; the platform's verification surface stops at checksum/hash and total size.**
Source: UX-3 § Phase A Step 3.

**BR-38. The platform shall not communicate with the tenant's end users at any point in the eviction journey; notification of end users shall be the capability owner's responsibility.**
Source: UX-3 § Phase A Step 2; § Constraints Inherited.

---

## F. Operator-Initiated Tenant Update (UX-4)

**BR-39. When a platform-level dependency event (vendor sunset, CVE, EOL) forces an update, the operator shall file one `platform update required` issue per affected tenant, naming the falling-behind component, the replacement, the requested update shape, and the deadline with its external reason.**
Source: UX-4 § Step 1.

**BR-40. The default response to a tenant fall-behind shall be working *with* the capability owner to bring the tenant current, not eviction.**
Source: `_index.md` § Business Rules ("Eviction is allowed when needs and capabilities diverge"); UX-4 § Persona, § Constraints Inherited.

**BR-41. The operator shall be permitted to negotiate an extended delivery date with the capability owner only to the extent the external pressure leaves safe slack; if no safe slack exists, the inherited deadline stands.**
Source: UX-4 § Step 4.

**BR-42. A missed operative delivery date (inherited or extended) shall be the operational signal that the eviction-threshold rule has been crossed; the operator shall then open a separate eviction issue linking back to the update issue.**
Source: UX-4 § Step 5.

**BR-43. Concurrent forcing events affecting the same tenant shall be tracked on separate `platform update required` issues, cross-linked where remediation overlaps.**
Source: UX-4 § Step 1; § Edge Cases.

---

## G. Platform-Contract-Change Rollout (UX-5)

**BR-44. When the operator proactively changes a contract term, the change shall be communicated ahead of time via a single umbrella `platform contract change` issue tagging every affected capability owner.**
Source: UX-5 § Step 1; `_index.md` § Business Rules (the evergreen-contract promise).

**BR-45. The umbrella issue shall contain at minimum: what is changing, what it is changing to (or that it is being removed), the migration guideline if applicable, the hard deadline, the reason for the change, and the operator's chosen status-update cadence.**
Source: UX-5 § Step 1.

**BR-46. The deadline on a contract-change rollout shall be chosen to give every affected tenant at least two full status-update cycles before cutoff (one to acknowledge and start, one to finish or surface blockers).**
Source: UX-5 § Step 1.

**BR-47. Capability owners shall be required to acknowledge contract-change umbrella issues in-thread; silence shall be treated as non-engagement and ultimately as failure to migrate.**
Source: UX-5 § Step 2; § Edge Cases.

**BR-48. The deadline shall not be negotiable per-tenant; it shall apply uniformly. Global extensions are permitted only when the migration guideline or replacement itself must change for the remaining tenants.**
Source: UX-5 § Step 2; § Edge Cases.

**BR-49. During a contract-change rollout window the platform shall serve both the old and the new form of the contract concurrently, except for full offering removals where no concurrency is possible.**
Source: UX-5 § Step 3.

**BR-50. The operator shall post status updates on a regular schedule in the umbrella thread, sized to the rollout's overall timeline; the current snapshot shall live in the issue body and each scheduled update shall also be posted as a thread comment.**
Source: UX-5 § Step 3.

**BR-51. On the deadline, the old form shall be removed from the platform and a separate eviction issue shall be opened per laggard tenant linking back to the umbrella.**
Source: UX-5 § Step 4.

**BR-52. Any replacement offering required by a contract change shall be implemented and running on the platform *before* the umbrella issue is filed; building it is a precondition, not a step in the rollout.**
Source: UX-5 § Entry Point; § Out of Scope.

---

## H. Standing Up / Rebuilding the Platform (UX-6)

**BR-53. The platform shall be rebuildable from its definitions repo with no manual snowflake configuration required.**
Source: `_index.md` § Purpose & Business Outcome (Reproducibility); UX-6 § Goal, § Success.

**BR-54. The same rebuild flow shall serve first-ever build, disaster recovery, and reproducibility drills.**
Source: UX-6 § Entry Point.

**BR-55. The operator shall run a reproducibility drill on parallel scratch infrastructure after every significant platform change and at least quarterly, identical in flow to the real rebuild.**
Source: UX-6 § Entry Point; § Constraints Inherited.

**BR-56. Before any rebuild that involves prior platform state, the operator shall run a preflight drift check against the live platform or the last known-good environment; the rebuild shall not proceed until unexplained drift is resolved.**
Source: UX-6 § Entry Point; § Step 1.

**BR-57. The rebuild shall be automated end-to-end with manual operator-validation checkpoints between phases (foundations → core services → cross-cutting services → readiness verification).**
Source: UX-6 § Journey overview, § Steps 3–6.

**BR-58. Each rebuild phase shall be reversible — a clean teardown of any partially-provisioned state shall always be a viable, reliable option at every checkpoint.**
Source: UX-6 § Edge Cases ("Phase fails mid-rebuild"); § Constraints Inherited.

**BR-59. Readiness shall be declared only after a purpose-built canary tenant maintained alongside the platform definitions has been deployed end-to-end, exercised, and torn down successfully.**
Source: UX-6 § Step 6; § Constraints Inherited.

**BR-60. If the canary fails, the platform shall not be marked ready for tenants regardless of how green prior phases looked.**
Source: UX-6 § Edge Cases.

**BR-61. A rebuild that exceeds the 1-hour KPI shall not block the platform from going into service, but shall result in a tracked GitHub issue capturing the cause for follow-up.**
Source: UX-6 § Step 7; § Edge Cases.

---

## I. Tenant-Facing Observability (UX-7)

**BR-62. Each capability owner shall have a tenant-scoped observability view with a working login, provisioned automatically as part of `onboard my capability`; access shall not be a separately-requested add-on.**
Source: UX-7 § Entry Point; § Step 1.

**BR-63. Capability owners shall not be able to see other tenants' data; only the operator shall have cross-tenant visibility.**
Source: UX-7 § Step 2; § Constraints Inherited.

**BR-64. The platform shall expose a standard health bundle for every tenant: availability, latency, error rate, resource saturation, and restart / deployment events.**
Source: UX-7 § Step 1.

**BR-65. The platform shall send email alerts to the capability owner when their tuned thresholds are crossed; alerts shall name the signal and the capability.**
Source: UX-7 § Step 4.

**BR-66. Capability owners shall be able to self-serve their alert thresholds inside the observability offering; this is the only self-service surface the platform exposes.**
Source: UX-7 § Step 3.

**BR-67. The pull view shall be authoritative for current health; email shall be a best-effort acceleration nudge. When email delivery is degraded for a tenant, the tenant view shall surface that fact.**
Source: UX-7 § Step 1, § Step 3, § Edge Cases.

**BR-68. Requests to expand the signal bundle or add non-email alert channels shall be handled via `modify my capability`, not via this UX.**
Source: UX-7 § Step 1, § Out of Scope.

**BR-69. End users of tenant capabilities shall not receive any observability access from the platform.**
Source: UX-7 § Constraints Inherited; `_index.md` § Business Rules ("No direct end-user access to the platform").

---

## J. Operator Model & Succession

**BR-70. Only the operator shall operate the platform and have administrative access; no co-operators and no delegated administration shall exist.**
Source: `_index.md` § Business Rules ("Operator-only operation").

**BR-71. The platform shall designate a successor operator who holds sealed/escrowed credentials and a runbook sufficient to keep the platform running if the primary operator becomes unavailable.**
Source: `_index.md` § Business Rules ("Operator succession").

**BR-72. Successor credentials shall not be used for routine operation; takeover shall be a discrete event triggered by operator unavailability.**
Source: `_index.md` § Business Rules ("Operator succession"); UX-6 § Persona.

**BR-73. The platform shall provide on-demand exportable archives so each tenant's users can retrieve their own content without operator involvement *while the platform is healthy*; users are expected to pull these proactively and may schedule periodic pulls.**
Source: `_index.md` § Business Rules ("Operator succession").

**BR-74. The export mechanism that powers operator-succession archives shall be the same export tool used by the eviction journey (UX-3).**
Source: UX-3 § Constraints Inherited.

---

## K. Eviction Threshold (cross-cutting)

**BR-75. A tenant shall be evicted when accommodating it would either push routine operation sustainably above 2x the *Operator maintenance budget* KPI, or break the *Reproducibility* KPI by requiring snowflake configuration that cannot be captured as definitions. Either condition alone is sufficient.**
Source: `_index.md` § Business Rules ("Eviction threshold").

**BR-76. The eviction threshold shall apply across all triggers (UX-1 § Edge Cases, UX-4 step 5, UX-5 step 4) — the threshold is the single rule, the triggering events differ.**
Source: `_index.md` § Business Rules; UX-4 § Constraints Inherited; UX-5 § Constraints Inherited.

---

## L. KPIs (Success Criteria)

**BR-77. KPI — Tenant adoption.** Every *implemented* capability defined in this repo shall run on this platform. An implemented capability hosted elsewhere shall count negatively against this KPI. A capability owner who explicitly gives up on onboarding because the operator stayed silent too long shall be counted as a lost tenant on the issue itself.
Source: `_index.md` § Success Criteria; UX-1 § Edge Cases, § Constraints Inherited.

**BR-78. KPI — Reproducibility.** The platform shall be standable from its definitions in **at most 1 hour**, starting from no platform at all.
Source: `_index.md` § Success Criteria.

**BR-79. KPI — Operator maintenance budget.** Routine operation shall take **no more than 2 hours per week** of the operator's time. If maintenance regularly exceeds this, the platform must be simplified, not grown.
Source: `_index.md` § Success Criteria.

**BR-80. KPI — Cost stays proportional to value.** Total operating cost shall remain within what the operator considers acceptable given delivered convenience and resiliency. There is no fixed dollar target; the test is whether the operator would still choose to run it knowing the bill.
Source: `_index.md` § Success Criteria.

---

## M. Scope Boundaries (what the platform shall NOT do)

**BR-81. The platform shall not host capabilities for anyone other than the operator's own capabilities — no third parties, no public, no family/friends as direct platform users.**
Source: `_index.md` § Out of Scope.

**BR-82. The platform shall not commit to any specific availability or performance SLA. Tenants needing stronger guarantees shall host elsewhere.**
Source: `_index.md` § Out of Scope; § Business Rules; UX-1 § Constraints Inherited.

**BR-83. The platform shall not prescribe its own implementation (homelab, Kubernetes, or any specific stack). The capability is satisfied by anything that meets its rules and KPIs.**
Source: `_index.md` § Out of Scope.

**BR-84. The platform shall not provide multi-operator administration, role delegation, or self-service onboarding.**
Source: `_index.md` § Out of Scope.

**BR-85. The platform shall not write, debug, or shepherd tenant migration processes; it provides only the runner.**
Source: UX-2 § Out of Scope.

**BR-86. The platform shall not offer rollback after eviction; the capability owner runs the old host concurrently during migration and "rollback" simply means not cutting over.**
Source: UX-2 § Out of Scope; UX-3 § Out of Scope.

**BR-87. The platform shall not restore tenant data as part of platform standup; tenant data restoration is a separate UX that begins where UX-6 ends.**
Source: UX-6 § Out of Scope.

**BR-88. The platform shall not provide end-user-facing observability or status pages; that is a tenant concern.**
Source: UX-7 § Out of Scope; `_index.md` § Out of Scope ("End-user-facing features of tenant capabilities").

---

## N. Cross-Cutting Properties

**BR-89. The platform may span public and private infrastructure. "Self-hosted" means the operator controls the platform end-to-end, not that every component runs on hardware the operator owns.**
Source: `_index.md` § Business Rules.

**BR-90. Cost shall be secondary to convenience and resiliency in trade-offs. Added cost is acceptable when it buys meaningful convenience or resiliency; cost shall still be minimized where it does not cost convenience or resiliency.**
Source: `_index.md` § Business Rules.

**BR-91. Operator skill development shall not influence buy-vs-build trade-offs. Trade-offs are judged on convenience, resiliency, and cost only.**
Source: `_index.md` § Business Rules.

**BR-92. Where outcomes conflict, the established tiebreaker order is: tenant adoption > reproducibility > vendor independence > minimizing operator effort.**
Source: `_index.md` § Purpose & Business Outcome.

**BR-93. Every UX that can introduce platform state shall enforce tracked changes and immutability rather than allowing ad-hoc modification (so the preflight drift check in UX-6 is meaningful).**
Source: UX-6 § Constraints Inherited.

---

## Open Questions (carried forward from source docs)

- **OQ-1.** What is the authoritative policy for deeper backup-tier copies after the 30-day post-eviction retention window ends — retention duration, deletion behavior, and operator-access/privacy constraints?
  Source: UX-3 § Open Questions.

All other source documents currently report "_None at this time._" under Open Questions.
