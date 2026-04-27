# Actions log — Stage 1 re-extraction for `self-hosted-application-platform`

## Reads

- `/.claude/skills/define-tech-design/SKILL.md` — loaded the skill instructions.
- `/docs/content/capabilities/self-hosted-application-platform/_index.md` — capability doc, end-to-end.
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/_index.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/host-a-capability.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/stand-up-the-platform.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/migrate-existing-data.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/move-off-the-platform-after-eviction.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/operator-initiated-tenant-update.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/platform-contract-change-rollout.md`
- `/docs/content/capabilities/self-hosted-application-platform/user-experiences/tenant-facing-observability.md`
- `/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` — pre-existing TR-01..TR-07.
- `/docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md` — observed parallel structure.
- `/docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md` — observed parallel REQ-NN-style doc.

## Writes

- `/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` — preserved TR-01..TR-07; appended TR-08..TR-38; set `reviewed_at: null` to re-open the Stage-2 gate; refreshed Open Questions list.
- `/.claude/skills/define-tech-design-workspace/iteration-2/eval-1-stage1-extracts/with_skill/run-1/outputs/artifacts/content/capabilities/self-hosted-application-platform/tech-requirements.md` — workspace copy of the above (mirrors the in-tree path under `content/...`).
- `/.claude/skills/define-tech-design-workspace/iteration-2/eval-1-stage1-extracts/with_skill/run-1/outputs/response.md` — final response to caller.
- `/.claude/skills/define-tech-design-workspace/iteration-2/eval-1-stage1-extracts/with_skill/run-1/outputs/actions.md` — this log.

## Skill-conformance notes

- **Append-only TR identity preserved.** TR-01..TR-07 left unmodified; TR-08..TR-38 appended.
- **Source links** present on every TR, using Hugo-relative paths into `_index.md` and `user-experiences/*.md`, plus one cite into `/CLAUDE.md` for the inherited Cloudflare→GCP topology.
- **Decisions kept out of TRs.** Where the source UX or capability volunteered something that looked like a chosen approach (e.g. "GitHub issues" as the engagement surface), it was either (a) inherited as a UX-mandated requirement because every UX consumes it as the only surface (TR-23), or (b) demoted to the Open Questions list as a Stage-2 ADR seed.
- **Stage-2 gate re-opened.** `reviewed_at: null`. The skill will refuse Stage 2 until the human sets `reviewed_at` to a date newer than the file's mtime.
- **No ADRs drafted, no tech-design.md written.** Stages 2 and 3 are explicitly gated and were not entered.
- **Cross-capability flag.** Network architecture (TR-07 / TR-26 territory) noted as a likely shared-ADR candidate to be flagged at the start of Stage 2 rather than scoped to this capability.
