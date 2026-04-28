# Technical Requirements: Self-Hosted Application Platform

This document extracts the technical requirements implied by the
[Self-Hosted Application Platform](../../../../../../../docs/content/capabilities/self-hosted-application-platform/_index.md)
capability and its user experiences. Each requirement is traced back to the
business rule, KPI, output, or UX step it derives from. No technology choices,
architectures, or designs are made here — those belong to subsequent ADRs.

Convention:
- **Source** links cite the originating doc and section. `cap` = parent
  capability `_index.md`; UX docs are referenced by filename.
- Requirements are grouped by concern, not by source doc. A requirement may
  be reinforced by multiple sources; all are listed.

---

## 1. Platform Outputs (what the platform must provide to tenants)

### TR-1.1 — Provide tenant compute
The platform must provide a compute environment in which a tenant capability's
application code can run.
- **Source:** cap *Outputs & Deliverables — Compute*; cap *Triggers & Inputs* (resource needs declaration); `host-a-capability.md` step 5.

### TR-1.2 — Provide tenant persistent storage
The platform must provide durable persistent storage for each tenant's data.
- **Source:** cap *Outputs & Deliverables — Persistent storage*; `host-a-capability.md` step 5; `migrate-existing-data.md` (destination tenant storage interfaces).

### TR-1.3 — Provide internal network reachability between tenants
The platform must provide network reachability between tenants running on the
platform.
- **Source:** cap *Outputs & Deliverables — Network reachability (internal)*.

### TR-1.4 — Provide external network reachability for tenants
The platform must provide network reachability that allows a tenant's end
users (outside the platform) to reach the tenant.
- **Source:** cap *Outputs & Deliverables — Network reachability (external)*.

### TR-1.5 — Provide an identity & authentication offering for tenant end users
The platform must provide an identity and authentication service usable by
any tenant whose end users need to authenticate.
- **Source:** cap *Outputs & Deliverables — Identity & authentication*; cap *Triggers & Inputs* (identity choice).

### TR-1.6 — Allow tenants to bring their own identity
A tenant must be able to opt out of the platform-provided identity service
and bring its own.
- **Source:** cap *Outputs & Deliverables*; cap *Business Rules — Identity service honors tenant credential-recovery rules*; `host-a-capability.md` *Constraints — Identity service honors tenant credential-recovery rules*.

### TR-1.7 — Platform-provided identity must support unrecoverable credentials
The platform-provided identity service must be capable of honoring a
"lost credentials cannot be recovered" (Signal-style) property for tenants
that require it.
- **Source:** cap *Business Rules — Identity service honors tenant credential-recovery rules*.

### TR-1.8 — Provide tenant data backup & disaster recovery
The platform must back up tenant data and provide a disaster-recovery
capability for it, to a standard the platform itself defines.
- **Source:** cap *Outputs & Deliverables — Backup and disaster recovery*.

### TR-1.9 — Provide observability of every tenant
The platform must produce observability signals for every tenant such that
the operator can tell whether each tenant is up and healthy without the
tenant having to instrument that itself.
- **Source:** cap *Outputs & Deliverables — Observability*; `tenant-facing-observability.md` *Constraints — Direct outputs include observability*.

### TR-1.10 — Provide a secret-management offering
The platform must provide a secret-management offering that capability owners
can register secrets with, and that tenant components / migration jobs can
read by name.
- **Source:** `migrate-existing-data.md` step 1; `migrate-existing-data.md` step 3 (secrets registered, named on issue).

### TR-1.11 — Provide a one-shot migration-process runner offering
The platform must provide an offering that runs a tenant-supplied one-shot
migration process: a packaged job that reads from an external source and
writes into an existing tenant via that tenant's normal interfaces, with
the platform's standard observability, with the job torn down on completion.
- **Source:** `migrate-existing-data.md` *Goal*, steps 4–8, *Constraints — The capability evolves with its tenants*.

### TR-1.12 — Provide a tenant-data export offering
The platform must provide a tenant-data export tool that, when invoked,
produces a downloadable archive of a tenant's data along with a
checksum/hash and total size in bytes. Export tooling must exist for every
kind of data the platform hosts (it is a core platform feature, not a
per-tenant artifact).
- **Source:** `move-off-the-platform-after-eviction.md` steps 3, 6; *Edge Cases — Export tooling does not exist*; *Constraints — Operator succession*.

### TR-1.13 — Export must be invokable on demand and re-runnable
The export tool must be invokable on demand by capability owners (within the
permitted lifecycle of their tenant) and must support being run multiple
times, each producing a fresh archive.
- **Source:** `move-off-the-platform-after-eviction.md` steps 3, 4, 6; cap *Business Rules — Operator succession* (on-demand exportable archives, periodic pulls).

### TR-1.14 — Generated export archives are ephemeral
The platform is not required to retain previously-generated export archives
for later pickup. Archives must be downloaded by the capability owner at
the time of generation; missed downloads require re-running the tool.
- **Source:** `move-off-the-platform-after-eviction.md` steps 3, 6.

---

## 2. Platform Contract & Tenant Packaging

### TR-2.1 — Define a tenant packaging form
The platform must define a single, documented packaging form ("the form
the platform accepts") in which tenant components — including migration
processes — are submitted to be deployed.
- **Source:** cap *Triggers & Inputs* (capability packaged in the form the platform accepts); cap *Business Rules — Tenants must accept the platform's contract*; `host-a-capability.md` step 4; `migrate-existing-data.md` *Constraints — Tenants must accept the platform's contract*.

### TR-2.2 — Tenants declare resource needs up front
The platform's contract must require tenants to declare their resource needs
(compute, storage, network reachability — internal and external) and
availability expectations at submission time.
- **Source:** cap *Triggers & Inputs*; cap *Business Rules — Tenants must accept the platform's contract*; `host-a-capability.md` *Entry Point*; `migrate-existing-data.md` step 2.

### TR-2.3 — Migration jobs declare a re-run contract
A submitted migration process must declare whether it is safe to run against
an already-populated destination tenant or whether the destination must be
wiped/empty before each run.
- **Source:** `migrate-existing-data.md` step 2, *Edge Cases — top-up migration*.

### TR-2.4 — Migration peak-footprint cap
The platform must reject (during operator review) a migration whose peak
temporary footprint (steady-state tenant footprint plus declared migration
spike, in compute and in storage) exceeds 2x the destination tenant's
steady-state compute and storage footprint, or which exceeds the platform's
currently available migration-process capacity.
- **Source:** `migrate-existing-data.md` step 3; *Edge Cases — Migration job needs more resources than declared*.

### TR-2.5 — Support concurrent migrations
The migration-process offering must support concurrent migration jobs
across different tenants without making any one tenant's job exclusive.
- **Source:** `migrate-existing-data.md` step 4; *Edge Cases — Another tenant is migrating at the same time*.

### TR-2.6 — Run old and new contract forms concurrently during a contract-change rollout
When a platform-contract change replaces an old offering with a new one, the
platform must be able to serve both the old and new forms simultaneously
during the rollout window. (Full-removal contract changes — where there is
no replacement — are exempt.)
- **Source:** `platform-contract-change-rollout.md` step 3, *Constraints — KPI: 1-hour reproducibility* (concurrent old/new during rollout is fine; permanent dual-form support is not the goal).

### TR-2.7 — Replacement offering must exist before its contract change rolls out
Where a contract change replaces an old offering with a new one, the
replacement offering must already be implemented and running on the
platform before the umbrella issue is filed.
- **Source:** `platform-contract-change-rollout.md` *Entry Point*, *Out of Scope — Building the replacement offering*.

---

## 3. Tenant-Facing Observability

### TR-3.1 — Tenant-scoped observability view
The observability offering must expose a per-tenant view that a capability
owner can authenticate to, and that confines them to their own tenant's
signals for the entire session. There must be no mode-switch that broadens
their scope to other tenants or to operator-wide views.
- **Source:** `tenant-facing-observability.md` *Entry Point — Pull entry*, step 2; *Constraints — Operator-only operation*.

### TR-3.2 — Platform-standard health bundle per tenant
For every live tenant, the observability offering must provide at minimum:
availability, latency, error rate, resource saturation, and restart /
deployment events.
- **Source:** `tenant-facing-observability.md` step 1; *Constraints — Tenants must accept the platform's contract*.

### TR-3.3 — Email alerting to capability owners
The platform must deliver alerts to capability owners by email, scoped to
their tenant. Alert content must include which signal fired and which
capability it pertains to.
- **Source:** `tenant-facing-observability.md` steps 1, 4.

### TR-3.4 — Self-service threshold tuning
A capability owner must be able to self-serve threshold values for their
own email alerts inside the observability offering, without filing an
issue and without operator involvement. This is the only self-service
surface in the platform.
- **Source:** `tenant-facing-observability.md` step 3; *Constraints — Operator-only operation* (carve-out).

### TR-3.5 — Surface degraded alert delivery in the tenant view
When the observability offering knows that email delivery is degraded for
a tenant, the tenant view must indicate that alerting is degraded so the
capability owner does not treat email silence as evidence of health.
- **Source:** `tenant-facing-observability.md` step 3; *Edge Cases — Alert delivery is broken*.

### TR-3.6 — Pull view is authoritative for current health
The tenant view (pull) must function as the source of truth for current
tenant health; email alerts are a best-effort acceleration channel.
- **Source:** `tenant-facing-observability.md` step 1; *Edge Cases — Alert delivery is broken*.

### TR-3.7 — Auto-provision observability access at onboarding
A capability owner's observability access (login + email alert wiring +
health bundle) must be provisioned automatically as part of onboarding,
not as a separate later request.
- **Source:** `tenant-facing-observability.md` *Entry Point*, step 1; `host-a-capability.md` step 5.

### TR-3.8 — Operator receives the same signals
Whatever signals are surfaced to capability owners must also be surfaced
to the operator, so the operator can see tenant-level health without the
capability owner having to escalate.
- **Source:** `tenant-facing-observability.md` step 5b, *Constraints — KPI: 2-hr/week*.

---

## 4. Engagement Surface (operator <-> capability owner)

### TR-4.1 — GitHub issues are the only engagement surface
All capability-owner-initiated engagement with the platform must occur
via GitHub issues filed against the infra repo. There must be no
self-service portal or alternate front door, with the single carve-out of
TR-3.4 (threshold tuning).
- **Source:** `host-a-capability.md` step 1; cap *Business Rules — Operator-only operation*; `tenant-facing-observability.md` *Constraints — Operator-only operation*.

### TR-4.2 — Distinct issue types
The infra repo must define distinct issue types, each carrying a different
operator-review scope and lifecycle:
- `onboard my capability`
- `modify my capability`
- `migrate my data`
- `platform update required`
- `platform contract change` (umbrella)
- `eviction` (operator-filed)
- **Source:** `host-a-capability.md` steps 1, 8; `migrate-existing-data.md` step 2; `operator-initiated-tenant-update.md` step 1; `platform-contract-change-rollout.md` step 1; `move-off-the-platform-after-eviction.md` *Entry Point*.

### TR-4.3 — Cross-link related issues
The platform's process must support and use cross-linking between issues
where journeys hand off (umbrella ↔ per-tenant `modify`; `platform update
required` ↔ `eviction`; `platform contract change` ↔ `eviction`;
`migrate my data` ↔ closed `onboard my capability`).
- **Source:** `operator-initiated-tenant-update.md` step 5; `platform-contract-change-rollout.md` steps 3, 4; `migrate-existing-data.md` step 2.

### TR-4.4 — Umbrella-issue body holds current rollout snapshot
For a `platform contract change` umbrella issue, the issue body must hold
the current rollout snapshot (latest state) and each scheduled status
update must additionally be posted as a thread comment for history.
- **Source:** `platform-contract-change-rollout.md` step 3.

### TR-4.5 — Rollout status updates carry standard metrics
Each scheduled status update on a contract-change umbrella must include:
how many tenants remain on the old form, how many have migrated, which
`modify` issues are open, and time remaining until the deadline.
- **Source:** `platform-contract-change-rollout.md` step 3.

---

## 5. Operator Operation

### TR-5.1 — Single operator, no delegated administration
The platform must be operable by exactly one operator at a time. No
co-operator role, no delegated administration, no shared day-to-day
administrative access.
- **Source:** cap *Business Rules — Operator-only operation*; `stand-up-the-platform.md` *Persona*; `tenant-facing-observability.md` *Constraints*.

### TR-5.2 — Sealed/escrowed successor credentials
The platform must support a designated successor operator who holds sealed
or escrowed emergency credentials and a runbook. These credentials are not
used for routine operation; takeover is a discrete event.
- **Source:** cap *Business Rules — Operator-only operation*, *Operator succession*; `stand-up-the-platform.md` *Persona*, *Edge Cases — Successor at the keyboard*.

### TR-5.3 — Operator-only cross-tenant visibility
Cross-tenant visibility (e.g. an operator-wide observability view, the
ability to see which tenants are using a particular platform offering)
must be available only to the operator role.
- **Source:** cap *Business Rules — Operator-only operation*; `operator-initiated-tenant-update.md` *Constraints — Operator-only operation*; `tenant-facing-observability.md` *Out of Scope — Operator-side observability*.

### TR-5.4 — Routine maintenance fits within 2 hr/week
Routine operation of the platform, summed across all tenants and all
offerings, must take no more than 2 hours per week of operator time.
This is the *Operator maintenance budget* KPI and is the bound against
which eviction thresholds, rollout cadence, and offering admission are
judged.
- **Source:** cap *Success Criteria & KPIs — Operator maintenance budget*; cap *Business Rules — Eviction threshold*; `host-a-capability.md` *Constraints — KPI: 2-hr/week*; `tenant-facing-observability.md` *Constraints — KPI: 2-hr/week*; `platform-contract-change-rollout.md` *Constraints — KPI: 2-hr/week*; `migrate-existing-data.md` *Constraints — KPI: 2-hr/week*.

### TR-5.5 — End users have no access to the platform
End users of tenant capabilities must have no direct access to platform
surfaces (no platform-side accounts, no platform-side notifications, no
platform-rendered "tenant retired" pages). Their only interaction with
the platform is via the tenant capability itself.
- **Source:** cap *Business Rules — No direct end-user access to the platform*; `move-off-the-platform-after-eviction.md` *Constraints — No direct end-user access*; `tenant-facing-observability.md` *Out of Scope — End-user-facing observability*.

---

## 6. Reproducibility, Rebuild & Drift

### TR-6.1 — Platform fully expressed as definitions in a repo
The entire platform must be expressible as definitions stored in a
versioned definitions repo, sufficient to rebuild it from nothing.
There must be no manual snowflake configuration outside those definitions.
- **Source:** cap *Success Criteria & KPIs — Reproducibility*; cap *Purpose & Outcomes — Reproducibility*; `stand-up-the-platform.md` *Entry Point*; `host-a-capability.md` *Constraints — KPI: 1-hour reproducibility*.

### TR-6.2 — Rebuild from definitions in <= 1 hour
Rebuilding the platform from no platform at all to "ready to host tenants"
must complete within 1 hour. This is the *Reproducibility* KPI; missing it
does not block the platform from going into service but must trigger a
follow-up issue.
- **Source:** cap *Success Criteria & KPIs — Reproducibility*; `stand-up-the-platform.md` *Goal*, step 7, *Edge Cases — 1-hour KPI is missed*.

### TR-6.3 — Rebuild is automated with operator validation checkpoints
The rebuild must run as a single top-level entry-point that drives
provisioning automatically, pausing between phases for the operator to
validate before continuing.
- **Source:** `stand-up-the-platform.md` *Journey* preamble, step 2.

### TR-6.4 — Phased rebuild ordering
The rebuild must execute in this order, each phase pausing for operator
validation:
1. Foundations (cloud + home-lab base, networking including cloud↔home-lab
   connectivity).
2. Core platform services (compute, persistent storage, identity).
3. Cross-cutting services (backup, observability) — covering the platform
   itself before any tenant arrives.
4. Readiness verification via canary tenant.
- **Source:** `stand-up-the-platform.md` steps 3–6.

### TR-6.5 — Each rebuild phase must be reversible
The platform's definitions must support a clean teardown of any partially
provisioned state at every phase boundary, so that "delete everything and
start over" is always a viable rollback.
- **Source:** `stand-up-the-platform.md` *Edge Cases — Phase fails mid-rebuild*; *Constraints — Each phase must be reversible*.

### TR-6.6 — Maintain a purpose-built canary tenant
A purpose-built canary tenant must be maintained alongside the platform's
definitions. It must exercise: running, reachability, data store/read,
authentication via the platform-provided identity service, backup pickup,
and observability pickup; and then be torn down after each rebuild.
- **Source:** `stand-up-the-platform.md` step 6, *Constraints — Default hosting target*.

### TR-6.7 — Canary success is the readiness gate
The platform may not be marked ready to host tenants until the canary
tenant comes up green end-to-end and is cleanly torn down — regardless of
how green prior phase signals were.
- **Source:** `stand-up-the-platform.md` step 6, *Edge Cases — Canary tenant fails*.

### TR-6.8 — Required preflight drift check before any rebuild with prior state
Before any rebuild that has prior platform state to compare against
(i.e. any rebuild other than first-ever), a preflight drift check must run
against the live platform or last known-good environment. The rebuild
must not start unless the check passes.
- **Source:** `stand-up-the-platform.md` *Entry Point*, step 1, *Edge Cases — Preflight drift check fails*.

### TR-6.9 — Tracked changes & immutability across all platform UXs
Every UX that can introduce platform state must enforce tracked changes
(via the definitions repo) and immutability of platform configuration —
no ad-hoc out-of-band modification — so the preflight drift check
(TR-6.8) is meaningful.
- **Source:** `stand-up-the-platform.md` *Constraints — Tracked changes and immutability*.

### TR-6.10 — Drills run after every significant platform change and at least quarterly
The same rebuild flow must be exercised on parallel scratch infrastructure
after every change that would alter what is rebuilt, what must be
validated, or what must be trusted before declaring readiness — and at
minimum quarterly.
- **Source:** `stand-up-the-platform.md` *Entry Point — Drift / reproducibility drill*; *Constraints — KPI: 1-hour reproducibility*.

### TR-6.11 — Migration-process and export offerings are themselves reproducible
The migration-process offering and the export offering, like every other
offering, must themselves be reproducible from definitions. Specific
in-flight migration jobs and generated export archives are *not* part of
the reproducible state — they are ephemeral artifacts.
- **Source:** `migrate-existing-data.md` *Constraints — KPI: 1-hour reproducibility*; `move-off-the-platform-after-eviction.md` *Constraints — KPI: 1-hour reproducibility*.

---

## 7. Eviction & Tenant Wind-Down

### TR-7.1 — Operator-deprovisionable tenant compute & network
The platform must support deprovisioning a tenant's compute and network
resources on a chosen date while leaving the tenant's data in place in a
read-only, no-further-writes state.
- **Source:** `move-off-the-platform-after-eviction.md` step 5.

### TR-7.2 — 30-day post-eviction read-only retention
After tenant compute/network is torn down at eviction, the platform must
retain a tenant-accessible (export-only) copy of the tenant's data for
exactly 30 days, after which the tenant-accessible copy is removed.
This 30-day clock is hard.
- **Source:** `move-off-the-platform-after-eviction.md` *Journey* preamble, step 7, *Edge Cases — Capability owner asks for more time*.

### TR-7.3 — Pause retention countdown only on platform-side export failure
The 30-day tenant-accessible retention countdown may be paused (and only
paused) when the cause of an export failure is shown to be in the
platform's export tooling or its data hosting. Failures rooted in the
capability owner's own validation do not pause the countdown.
- **Source:** `move-off-the-platform-after-eviction.md` *Edge Cases — Export comes back wrong*, *Edge Cases — Export tooling does not exist*.

### TR-7.4 — Migration jobs torn down on completion
A migration job must be torn down on terminal success or after the
capability owner abandons the migration. Re-running later requires a
fresh `migrate my data` issue.
- **Source:** `migrate-existing-data.md` step 8, *Edge Cases — re-run months later*.

### TR-7.5 — Eviction is operator-initiated only
Eviction must be initiated only by the operator filing an eviction issue
that names the eviction date and the reason. There is no
capability-owner-initiated eviction journey.
- **Source:** `host-a-capability.md` *Edge Cases — Capability is evicted later*; `move-off-the-platform-after-eviction.md` *Entry Point*; cap *Business Rules — Eviction*.

---

## 8. Public + Private Infrastructure

### TR-8.1 — Platform may span public cloud and private home-lab
The platform's implementation must be allowed to span both public cloud
and private home-lab infrastructure simultaneously. Cloud-↔-home-lab
connectivity is part of the platform's foundation, not an afterthought.
- **Source:** cap *Business Rules — The platform may span public and private infrastructure*; `stand-up-the-platform.md` step 3, *Constraints — public and private infrastructure*.

### TR-8.2 — Operator retains control of public-cloud configuration & data
Where the platform uses public-cloud components, the operator must retain
control of configuration, data, and the ability to leave (no vendor
lock-in that prevents exit).
- **Source:** cap *Business Rules — The platform may span public and private infrastructure*; cap *Purpose — Independence from hosting vendors*.

---

## 9. Availability, Cost, Scope

### TR-9.1 — No availability or performance SLA
The platform must not commit to a specific availability or performance
SLA. Tenants needing stronger guarantees host elsewhere.
- **Source:** cap *Out of Scope — A specific availability or performance SLA*; cap *Business Rules — Tenants must accept the platform's contract*; `migrate-existing-data.md` *Constraints — No specific availability or performance SLA*; `platform-contract-change-rollout.md` *Constraints*.

### TR-9.2 — Platform hosts only the operator's own capabilities
The platform must not offer hosting to third parties, the public, or
family/friends as platform users. End users of tenant capabilities reach
the platform only through a tenant capability.
- **Source:** cap *Out of Scope — Hosting for anyone other than the operator's own capabilities*; cap *Business Rules — No direct end-user access*.

### TR-9.3 — Cost is bounded by perceived value, not by a fixed dollar target
The platform's total operating cost must remain within what the operator
considers acceptable given the convenience and resiliency it delivers.
There is no fixed dollar target, but cost must not be optimized at the
expense of convenience or resiliency.
- **Source:** cap *Business Rules — Cost is secondary to convenience and resiliency*; cap *Success Criteria — Cost stays proportional to value*.

---

## Open Questions / TBD

These are items the source documents leave open. They are flagged here
so they are not silently turned into requirements:

- **OQ-1.** Authoritative policy for deeper backup-tier copies after the
  30-day tenant-accessible retention window ends — retention duration,
  deletion behavior, and operator-access/privacy constraints.
  - **Source:** `move-off-the-platform-after-eviction.md` *Open Questions*; *Success* (TBD note); cap *Business Rules — Operator succession* (referenced in passing).

---

## Items deliberately *not* requirements

For traceability, the following appear in the source docs but are
explicitly *out of scope* for the platform and therefore generate no
technical requirements:

- A specific availability/performance SLA on the platform itself.
- A completion-time SLA for migration jobs or for export.
- End-user-facing observability or status pages.
- Operator-side incident-management workflow inside an issue.
- Helping a capability owner write/debug their migration process.
- Validation of *semantic* correctness of exported or migrated data
  (the platform's commitment is bytes + checksum + size only).
- Helping a capability owner choose where to run after eviction.
- Application/runtime/configuration migration tooling (only data export
  is a platform concern).
- Co-operator administration, role delegation, self-service onboarding.
- Recovery from loss of root-level foundations (cloud account itself,
  all home-lab access).
- Tenant-facing pending-update / deprecation signal ahead of an
  operator-filed `platform update required` issue. (If ever added,
  belongs in tenant-facing observability — not in this set.)

- **Source:** cap *Out of Scope*; per-UX *Out of Scope* sections.
