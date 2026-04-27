# Actions log — eval-9-re-extraction (with_skill)

## Stage detection
- Read `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design/SKILL.md`.
- Existing `tech-requirements.md` present with TR-01..TR-07 and `reviewed_at: 2026-04-20`.
- User explicitly requested re-extraction. Stage 1.

## Preconditions
- Confirmed parent capability is in page-bundle form: `docs/content/capabilities/self-hosted-application-platform/_index.md` exists.
- Read capability `_index.md` end-to-end (Purpose, Stakeholders, Triggers/Inputs, Outputs, Business Rules, Success Criteria/KPIs, Out of Scope).
- Listed and read UX docs under `user-experiences/`:
  - `_index.md`
  - `host-a-capability.md`
  - `migrate-existing-data.md` (full read — to verify the broken anchor)
  - `move-off-the-platform-after-eviction.md` (header sweep)
  - `operator-initiated-tenant-update.md` (header sweep)
  - `platform-contract-change-rollout.md` (header sweep)
  - `stand-up-the-platform.md` (header sweep)
  - `tenant-facing-observability.md` (header sweep)
- Verified `migrate-existing-data.md` has no heading matching `#a-section-that-no-longer-exists`. The page exists; the anchor does not. The TR text is still supported by the page as a whole (especially §Journey and §Success).
- Noted CLAUDE.md constraints (Cloudflare → GCP topology, no humus framework) — already encoded as TR-07.

## Re-extraction decisions
- **Preserved verbatim:** TR-01, TR-02, TR-03, TR-04, TR-05, TR-07. All source links still resolve; text still grounded.
- **TR-06:** preserved text and source link as-is; added inline flag block:
  `> ⚠️ source no longer resolves — human review. ...`
  Did NOT delete. Did NOT renumber. Did NOT silently rewrite the link. Per skill: human resolves the flag.
- **Appended TR-08:** Reproducibility ≤1hr. Sources: capability §Success Criteria & KPIs and stand-up-the-platform UX §Journey. Not present in prior file; clearly forced.
- **Appended TR-09:** Identity must honor "lost credentials cannot be recovered." Source: capability §Business Rules. Forced (rule explicitly excludes non-conforming identity options).
- **Appended TR-10:** Designated successor with sealed/escrowed credentials. Source: capability §Business Rules. TR-05 covers export half of succession; TR-10 covers takeover half (kept as separate TR — they're separately forced and could be addressed by separate ADRs).
- **Appended TR-11:** Single declared packaging form for tenant components with declared resource needs. Sources: host-a-capability UX §Journey and capability §Business Rules.
- **Did NOT add** TRs around availability/performance SLA, cost ceilings, identity bring-your-own choices, end-user reachability — capability either explicitly excludes them or treats them as out-of-scope. Not forced; not requirements.

## Gate handling
- Reset `reviewed_at: null` in frontmatter. Document content has materially changed (4 new TRs + flag on TR-06); the prior `2026-04-20` review is no longer valid.
- Stage 2 explicitly NOT entered. Response message tells the user the gate is closed and what they need to do to open it.

## Files written
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (natural path — re-extracted document with TR-01..TR-11, TR-06 flagged, `reviewed_at: null`).
- `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design-workspace/iteration-1/eval-9-re-extraction/with_skill/outputs/artifacts/capabilities/self-hosted-application-platform/tech-requirements.md` (eval copy).
- `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design-workspace/iteration-1/eval-9-re-extraction/with_skill/outputs/response.md` (final user-facing response).
- `/home/carson/github.com/Zaba505/infra/.claude/skills/define-tech-design-workspace/iteration-1/eval-9-re-extraction/with_skill/outputs/actions.md` (this log).

## What I deliberately did NOT do
- Did NOT enter Stage 2 / propose ADRs (gate is closed).
- Did NOT delete or renumber TR-06 (would break ADR provenance and skill's append-only rule).
- Did NOT silently fix the TR-06 link to `#journey` or similar (skill says: flag, let human resolve).
- Did NOT capture decisions as requirements (no Postgres, no Kubernetes, no specific identity vendor mentioned in the file).
- Did NOT add a `## Tech Design` cross-link section to `_index.md` (premature — `tech-design.md` doesn't exist yet).
</content>
