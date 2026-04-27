# Tech Design — Self-Hosted Application Platform: Stage 1 started

I've kicked off the tech design by reading the capability doc and all seven UX docs, then producing the **Stage 1 artifact**: a living technical-requirements extract that will feed the ADRs (Stage 2) and the composed tech-design document (Stage 3).

## What I produced

Two new files under the capability:

- `docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md` — landing page describing the three-stage flow and current status.
- `docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md` — the Stage 1 extract.

The requirements doc contains:

- **23 functional requirements** (`REQ-01`…`REQ-23`) covering tenant lifecycle (provision, modify, evict, migrate, operator-initiated update, contract rollout), the seven platform-provided offerings (compute, storage, network, identity, secrets, backup, observability, export), reproducibility/standup mechanics, and the engagement surface.
- **6 non-functional requirements** (`NFR-01`…`NFR-06`) including the 1-hour reproducibility KPI, the 2 hr/wk operator budget, the 2× eviction threshold, "no specific availability SLA," and the 30-day post-eviction grace window.
- **8 constraints** (`C-01`…`C-08`) including hybrid public/private infra, cross-environment connectivity, the operator-only rule, the credential-recovery property required of identity, the evergreen-contract rule, and the platform tiebreaker order (tenant adoption > reproducibility > vendor independence > operator effort).
- **12 open technical questions** (`Q-01`…`Q-12`) — these are the seeds for Stage 2 ADRs.
- A **traceability matrix** mapping each capability/UX source to the requirements it drives, so the lineage stays auditable as the design evolves.

## What I deliberately did *not* do yet

- **No ADRs written.** The skill flow is sequential: requirements first, then one ADR per non-trivial decision. I did not start drafting them in this turn because Q-01..Q-12 are real choices that need your input (especially the compute substrate, packaging form, and identity service — these have downstream coupling to almost everything else).
- **No composed tech-design document.** That is Stage 3 and should not be written until the ADRs have settled.
- **No shared/cross-capability ADRs.** Things like cloud-provider choice and the MADR format itself already live in `docs/content/r&d/adrs/` and are inputs here, not outputs.

## Recommended next step

Tackle the open questions in dependency order. My suggested first pass:

1. **Q-09 — definitions-repo layout and top-level rebuild entry point.** This is the spine of REQ-15..REQ-19 and the 1-hour KPI; everything else assumes it.
2. **Q-02 — compute substrate.** Constrains packaging (Q-01), networking (Q-04), identity hosting (Q-05), secret management (Q-06), and how observability (Q-08) is deployed.
3. **Q-01, Q-04, Q-05** can then be decided in parallel against that substrate.
4. **Q-03 (storage), Q-06 (secrets), Q-07 (backup), Q-08 (observability)** follow.
5. **Q-10–Q-12** (canary tenant, drift detection, successor escrow) are smaller and can be slotted in any order.

When you're ready, point me at one of those questions and I'll drive the ADR conversation.

## Key files

- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md`
