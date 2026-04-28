# Transcript

## User request

> Extract the technical requirements for the self-hosted-application-platform capability. The capability and its user experiences are already documented under docs/content/capabilities/self-hosted-application-platform/.

## Assistant response

Read the parent capability and all seven UX docs under
`docs/content/capabilities/self-hosted-application-platform/`:

- `_index.md` (parent capability)
- `user-experiences/host-a-capability.md`
- `user-experiences/migrate-existing-data.md`
- `user-experiences/move-off-the-platform-after-eviction.md`
- `user-experiences/operator-initiated-tenant-update.md`
- `user-experiences/platform-contract-change-rollout.md`
- `user-experiences/stand-up-the-platform.md`
- `user-experiences/tenant-facing-observability.md`

Extracted the technical requirements implied by those documents — without
making any architecture or technology choices — and wrote them to
`tech-requirements.md` in this outputs directory.

The document organizes ~50 requirements under nine concern areas:

1. Platform Outputs (compute, storage, networking, identity, backup,
   observability, secret management, migration runner, export tooling)
2. Platform Contract & Tenant Packaging (packaging form, declared resource
   needs, migration re-run contract, 2x peak-footprint cap, concurrent
   migrations, concurrent old/new contract serving during rollouts,
   replacement-offering precondition)
3. Tenant-Facing Observability (tenant-scoped view, standard health
   bundle, email alerts, self-service threshold tuning, degraded-delivery
   surfacing, pull-as-source-of-truth, auto-provisioning at onboarding)
4. Engagement Surface (GitHub issues only, distinct issue types,
   cross-linking, umbrella body holds snapshot, status-update metrics)
5. Operator Operation (single operator, sealed successor credentials,
   operator-only cross-tenant visibility, 2hr/week budget, no end-user
   access)
6. Reproducibility, Rebuild & Drift (definitions repo, ≤1hr rebuild,
   automated phased rebuild with validation checkpoints, reversibility,
   canary tenant as readiness gate, preflight drift check, tracked
   changes/immutability, drills after every change and at least quarterly,
   migration/export offerings reproducible while jobs/archives are not)
7. Eviction & Tenant Wind-Down (deprovisionable tenant resources,
   30-day post-eviction read-only retention, retention pause on
   platform-side export failure, migration teardown, operator-only
   eviction)
8. Public + Private Infrastructure (cross-environment platform,
   operator retains control over public-cloud components)
9. Availability, Cost, Scope (no SLA, operator's own capabilities only,
   cost vs. value)

Plus an `Open Questions / TBD` section (deeper backup-tier policy after
30-day window — already flagged as TBD in the source) and an
`Items deliberately not requirements` section preserving the source docs'
explicit out-of-scope items so they don't silently turn into requirements.

Each requirement includes a **Source** trace back to the originating
capability section, UX section, or KPI.

## Files produced

- `tech-requirements.md` — the extracted technical-requirements document
- `transcript.md` — this transcript
