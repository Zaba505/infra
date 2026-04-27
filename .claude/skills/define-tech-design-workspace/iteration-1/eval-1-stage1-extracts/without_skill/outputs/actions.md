# Actions log — Stage 1 tech design (without skill)

## Reads
- `docs/content/capabilities/self-hosted-application-platform/_index.md` — capability doc.
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/_index.md`
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/host-a-capability.md` (full)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/stand-up-the-platform.md` (full)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/migrate-existing-data.md` (head 50)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/move-off-the-platform-after-eviction.md` (head 50)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/operator-initiated-tenant-update.md` (head 50)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/platform-contract-change-rollout.md` (head 50)
- `docs/content/capabilities/self-hosted-application-platform/user-experiences/tenant-facing-observability.md` (head 50)
- `docs/content/r&d/adrs/0001-use-madr-for-architecture-decision-records.md` — MADR template reference.
- Listed `docs/content/r&d/adrs/` to confirm shared ADRs exist there (cross-capability scope).
- Read `CLAUDE.md` (system reminder) — noted Wireguard cloud↔home-lab pattern as relevant prior context.

## Writes (natural paths)
- `docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md` — new section landing page describing the three-stage flow and current status.
- `docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md` — Stage 1 artifact: 23 REQ + 6 NFR + 8 C + 12 Q + traceability matrix.

## Writes (eval output mirror)
- `outputs/response.md` — final response to the user.
- `outputs/actions.md` — this log.
- `outputs/artifacts/_index.md` — copy of the new tech-design landing page.
- `outputs/artifacts/technical-requirements.md` — copy of the Stage 1 extract.

## Decisions made during the run
- Treated this as **Stage 1 only** of the three-stage flow described in the `define-tech-design` skill description. Did not invoke the skill (the eval path is `without_skill/`).
- Chose to put the tech-design under the capability folder (`.../self-hosted-application-platform/tech-design/`) rather than under `r&d/adrs/`, matching the skill's stated rule that capability-scoped designs live next to the capability and only cross-capability ADRs go under `r&d/adrs/`.
- Did **not** draft any ADRs this turn — they require user input on the open questions Q-01..Q-12.
- Did **not** write a composed Stage-3 tech-design document — premature before ADRs exist.

## Notes
- Several `Bash` commands were denied by permission policy; substituted with `Read`/`Write` and absolute-path operations as needed. No workarounds attempted that would bypass intent.
- Did not run `go_workspace` / `go_vulncheck` (gopls MCP guidance) — task is pure-docs, no Go code touched.
