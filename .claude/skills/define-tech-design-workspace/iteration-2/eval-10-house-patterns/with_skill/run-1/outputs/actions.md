# Actions log — eval-10-house-patterns / with_skill / run-1

## Inputs

- Skill: `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design/SKILL.md`
- Repo patterns: `/home/carson/github.com/Zaba505/infra/CLAUDE.md`
- Reviewed fixture (Stage 1 output): `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design/evals/fixtures/tech-requirements-reviewed.md`
- Capability doc: `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/_index.md`
- User prompt: stage 2 — list options for tenant-control-plane API exposure and pick one.

## Steps

1. **Read SKILL.md** to confirm the three-stage flow and Stage 2 discipline (one ADR at a time, anchored in TRs, MADR template).
2. **Attempted `cp` of the reviewed fixture** into the worktree at `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` — Bash was sandbox-denied, so used the `Write` tool instead with the fixture's exact contents.
3. **Read** the reviewed fixture to confirm `reviewed_at: 2026-04-26` (current date). Stage 2 gate is open.
4. **Read** the capability `_index.md` for stakeholders, business rules, and KPIs (notably the *Operator maintenance budget* KPI, which made the house-pattern adherence load-bearing for the rationale).
5. **Listed shared ADRs** under `docs/content/r&d/adrs/` (0001-0007) to confirm no prior shared decision dictates API transport — there isn't one. ADR-0007 standardizes the *error response format*, which I cited indirectly through `pkg/errorpb`.
6. **Read the ADR template** at `.claude/skills/define-tech-design/assets/adr.template.md` to mirror frontmatter, status lifecycle, sections.
7. **Confirmed no `adrs/` directory exists** for this capability yet — first ADR, so I also create `adrs/_index.md` per the skill instructions.
8. **Identified the in-scope decision boundary**: only "API transport." Adjacent decisions (state storage, contract-version encoding scheme, long-running-op shape, export mechanism) are deliberately deferred — bundling would violate the "one decision at a time" discipline in SKILL.md.
9. **Enumerated five options** (A: HTTP+protobuf+chi house pattern, B: gRPC, C: JSON+OpenAPI, D: Connect, E: GraphQL) and evaluated each strictly against the cited TRs (TR-02, TR-04, TR-07) plus the CLAUDE.md house-pattern constraint.
10. **Picked Option A** (house pattern). Wrote the ADR with `status: proposed` per skill discipline — user must accept before Stage 3.
11. **Wrote** `adrs/_index.md` and `adrs/0001-tenant-control-plane-api-transport.md` to natural Hugo paths.
12. **Copied** all written artifacts under `outputs/artifacts/` for the eval harness.
13. **Wrote** `response.md` with the option list, choice, file paths, and next-step prompt; **wrote** this `actions.md`.

## Files written (natural paths)

- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` — copied via Write from reviewed fixture
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/adrs/_index.md` — new
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/adrs/0001-tenant-control-plane-api-transport.md` — new

## Files written (eval mirror)

- `outputs/artifacts/tech-requirements.md`
- `outputs/artifacts/_index.md`
- `outputs/artifacts/0001-tenant-control-plane-api-transport.md`
- `outputs/response.md`
- `outputs/actions.md`

## Notes / deviations

- Bash tool was sandbox-denied; used `Write` for the file copy, equivalent outcome.
- ADR `status` is `proposed` (not `accepted`) — the user is the only one who can accept per the skill, and the prompt asked for option-listing and a recommendation, not a fait accompli.
- Did NOT advance to Stage 3 (`tech-design.md`) — Stage 2 has open ADRs (TR-01, TR-03, TR-05, TR-06 are not addressed by any ADR yet), so the gate is not open even after this ADR is accepted.
