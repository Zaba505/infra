# Tech-design kickoff for `self-hosted-application-platform` — Stage 1 complete (gated)

We're in **Stage 1: extracting technical requirements**. Stage 2 (ADRs) and Stage 3 (composed tech-design) are blocked behind the review gate.

## What I did

1. **Read everything first**, per the skill's preconditions:
   - The capability doc at `docs/content/capabilities/self-hosted-application-platform/_index.md`.
   - All 7 UX docs under `docs/content/capabilities/self-hosted-application-platform/user-experiences/`: `host-a-capability`, `stand-up-the-platform`, `migrate-existing-data`, `move-off-the-platform-after-eviction`, `operator-initiated-tenant-update`, `platform-contract-change-rollout`, `tenant-facing-observability`.
   - Repo house patterns from `CLAUDE.md` (Cloudflare → GCP topology, no humus framework, chi/bedrock service shape).

2. **Re-extracted requirements into the canonical Stage-1 file** at `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`. The previous version held only TR-01 through TR-07, which under-covered the UX set (no requirements for export tooling shape, secret management, observability self-serve thresholds, contract-change concurrency mechanics, rebuild-phase reversibility, evergreen-contract acceptance, the eviction 30-day window, etc.). Per the skill's append-only rule, **TR-01 through TR-07 were preserved verbatim** and **TR-08 through TR-38 were appended** with their own source links back to the capability and UX docs. No renumbering.

3. **Reset the review gate.** The previous file had `reviewed_at: 2026-04-26` but its mtime was also 2026-04-26 (today) — the skill requires `reviewed_at` to be *newer* than mtime, so the gate was already borderline. With the re-extraction now bumping mtime, I set `reviewed_at: null` so the gate is unambiguously open. You need to re-review before Stage 2 can begin.

## Notable findings during extraction

- **A second, parallel requirements doc exists** at `docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md` using its own identifier scheme (`REQ-NN` / `Q-NN`), alongside `tech-design/_index.md`. The skill specifies the canonical Stage-1 location is `docs/content/capabilities/{name}/tech-requirements.md` (the file I updated). The duplicate under `tech-design/` should probably be either removed or reframed as a working scratch — left untouched here because deletion is a human call. Worth resolving before Stage 2 so ADR provenance has a single source of truth.
- **Several open questions surfaced as Stage-2 ADR seeds** (packaging form, compute substrate, storage architecture, identity implementation, secret manager, backup architecture, observability stack, definitions-repo entry point, canary tenant identity, drift detection, successor credential escrow, export shape, contract versioning, status-update format). These are captured under "Open Questions" in the requirements doc, *not* as TRs.
- **Cross-capability touchpoints to flag for Stage 2:** TR-07 already cites the existing Cloudflare → GCP topology as inherited; the network-architecture open question may need to defer parts of itself to a shared `docs/content/r&d/adrs/` ADR rather than a capability-scoped one. I'll surface this at the start of Stage 2.

## Stage gate — what you need to do next

> I've extracted the technical requirements to `tech-requirements.md`. **Review and edit it before we proceed.** When you're done, set `reviewed_at:` in the frontmatter to a date *newer* than the file's last modification (today is fine if you save right after) — or tell me to do it on your behalf and I'll record your verbal confirmation. I won't propose ADRs until that's set.

## Files touched

- Updated (in-tree): `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`

## Files NOT touched (intentionally)

- `docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md` — pre-existing Stage-1 narrative; flagged above for human resolution.
- `docs/content/capabilities/self-hosted-application-platform/tech-design/technical-requirements.md` — parallel REQ-NN-style doc; flagged above.
- `docs/content/capabilities/self-hosted-application-platform/_index.md` — no `## Tech Design` cross-link added yet; per the skill that's optional and offered after Stage 3.
- No ADRs yet (`adrs/` directory not created) — Stage 2 is gated.
- No `tech-design.md` yet — Stage 3 is gated behind Stage 2.
