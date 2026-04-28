# Business Requirements Extracted — Self-Hosted Application Platform

I read the parent capability doc and all seven user-experience docs under `docs/content/capabilities/self-hosted-application-platform/` and produced a sourced business-requirements document at:

- `business-requirements.md` (this folder)

## What's in it

93 numbered business requirements (BR-01 through BR-93), each citing the specific section of either the capability doc or one of the UX docs it derives from. Requirements are organized into 14 sections:

- **A. Hosting Offerings** (BR-01–BR-11) — compute, storage, network, identity, backup, observability, secret management, migration runner, export tool.
- **B. Engagement Surface** (BR-12–BR-16) — GitHub-issues-only, distinct issue types, no SLA, evergreen contract, implicit acceptance via tech design.
- **C. Onboarding a Tenant** (BR-17–BR-24) — operator-only approval, narrow review scope, three outcomes, packaged-artifact handoff.
- **D. Migrating Existing Data** (BR-25–BR-32) — tenant-supplied process, 2x footprint cap, declared re-run contract, no auto-clean, one-shot lifecycle.
- **E. Moving Off After Eviction** (BR-33–BR-38) — operator-initiated, 30-day retention, hard wall except for platform-side export bugs, capability-owner-owned validation.
- **F. Operator-Initiated Tenant Update** (BR-39–BR-43) — one issue per affected tenant, work-with-don't-evict default, slack-bounded extensions.
- **G. Platform-Contract-Change Rollout** (BR-44–BR-52) — single umbrella issue, mandatory acknowledgment, two-cycles-of-slack rule, concurrent old/new during rollout, replacement-must-pre-exist.
- **H. Standing Up the Platform** (BR-53–BR-61) — definitions-only rebuild, drill cadence, preflight drift check, phased automation with checkpoints, reversible phases, canary readiness gate.
- **I. Tenant-Facing Observability** (BR-62–BR-69) — automatic provisioning, tenant-scoped, standard health bundle, self-serve thresholds (the only self-service surface), pull-view-authoritative.
- **J. Operator Model & Succession** (BR-70–BR-74) — operator-only operation, sealed successor credentials, on-demand archives, shared export tool with eviction.
- **K. Eviction Threshold** (BR-75–BR-76) — single cross-cutting rule (2x maintenance budget OR breaks reproducibility).
- **L. KPIs** (BR-77–BR-80) — tenant adoption, 1-hour reproducibility, 2-hr/wk maintenance, cost-proportional-to-value.
- **M. Scope Boundaries** (BR-81–BR-88) — explicit "shall not" clauses pulled from the various Out-of-Scope sections.
- **N. Cross-Cutting Properties** (BR-89–BR-93) — public+private infra allowed, cost tiebreaker, skill-development excluded from trade-offs, outcome priority order, tracked-changes-and-immutability requirement.

## Open Questions surfaced

The capability and UX docs all say "None at this time" for open questions *except* UX-3, which still has one open question about the deep backup-tier retention/deletion policy beyond the 30-day tenant-accessible window. Carried forward as **OQ-1** at the end of the document.

## Notes on traceability

Every BR cites a source — usually the specific section heading and quoted phrase. UX docs are referenced as UX-1 through UX-7 with a legend at the top of the document. Nothing in the document was added without a source; nothing was inferred beyond what the docs already commit to.

## What was NOT done

- No technical requirements (those belong in a follow-up `tech-requirements.md`).
- No architectural decisions or option weighing.
- Nothing under `docs/` was modified.

Output saved at `/home/carson/github.com/Zaba505/infra/.claude/skills/extract-business-requirements-workspace/iteration-1/eval-1/without_skill/outputs/`.
