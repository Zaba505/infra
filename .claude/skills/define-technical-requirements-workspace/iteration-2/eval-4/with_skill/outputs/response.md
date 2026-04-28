# Re-extracted tech-requirements.md for self-hosted-application-platform

I re-extracted `tech-requirements.md` against the reviewed `business-requirements.md` (BR-01..BR-54, `reviewed_at: 2026-04-27`). The output is saved to `outputs/tech-requirements.md` (per the eval setup; nothing under `docs/` was modified).

## What changed vs. the prior file (TR-01..TR-07, reviewed_at 2026-04-20)

**Preserved (numbering held, append-only):**
- TR-01 through TR-07 are kept in place at their original numbers.
- Raw-path links in the existing TRs (`(_index.md)`, `(user-experiences/...)`) were rewritten to Hugo `ref` shortcodes, since the skill mandates `ref` and the prior file predated that rule. The TR text itself is unchanged for TR-04, TR-07; TR-01/02/03/05 were lightly extended to add the BR-NN primary citations the skill requires (each TR must cite a BR or an inherited constraint).

**Stale source flagged (TR-06):**
- TR-06's UX source `user-experiences/migrate-existing-data.md#a-section-that-no-longer-exists` does not resolve. Per the skill, I prepended `> ⚠️ source no longer resolves — human review` and left the broken link in place rather than deleting or silently rewriting it. I also added `BR-24` as the primary citation (which is the BR that forces this TR), so the TR is not unsourced while the UX anchor is being repaired. The Open Questions section flags this for human resolution.

**Appended (TR-08..TR-30):** 23 new TRs forced by BRs not previously translated:
- TR-08 rebuild-from-definitions (BR-02, BR-23)
- TR-09 successor takeover (BR-05, BR-42)
- TR-10 per-tenant compute/storage/network (BR-13)
- TR-11 identity service with unrecoverable-credentials mode (BR-14, BR-15)
- TR-12 backup/DR with published standard (BR-16)
- TR-13 operator-side per-tenant health observability (BR-17)
- TR-14 self-service alert thresholds (BR-20, BR-19)
- TR-15 alerting-degraded indicator on pull view (BR-21)
- TR-16 typed GitHub-issue engagement channel (BR-22)
- TR-17 onboarding from definitions only (BR-23, BR-02)
- TR-18 one-shot migration job lifecycle + footprint cap (BR-24, BR-26)
- TR-19 named-secret store with name-only references (BR-25)
- TR-20 export checksum/size emission (BR-28)
- TR-21 post-eviction state machine (BR-29, BR-30, BR-31)
- TR-22 clean-exit + re-onboarding eligibility (BR-54)
- TR-23 cross-environment connectivity foundational (BR-35, BR-36)
- TR-24 phased standup + no-resume-from-partial-state (BR-38, BR-39)
- TR-25 canary readiness gate (BR-37)
- TR-26 drift preflight + cross-UX immutability (BR-40)
- TR-27 scratch-infra rebuild drill (BR-41)
- TR-28 platform-update-required issue contract (BR-50, BR-51, BR-22)
- TR-29 contract-change umbrella status/acknowledgment/uniform-deadline (BR-46, BR-47, BR-48, BR-49)
- TR-30 separate linked eviction issue on missed deadline (BR-52, BR-22)

**`reviewed_at` reset to `null`:** Contents changed substantially (link-format conversion, stale-source flag, new BR citations on existing TRs, 23 appended TRs), so the prior `2026-04-20` review is no longer valid. The skill requires resetting on meaningful changes.

## BRs intentionally not translated into a new TR
- BR-01, BR-03, BR-10, BR-11, BR-12, BR-32, BR-33, BR-43, BR-44, BR-53 — these are governance/relationship/posture commitments (default-hosting, vendor-independence posture, eviction *thresholds*, fall-behind-default-is-update, evolve-the-platform default, maintenance-budget KPI, cost-proportionality value judgment, skill-development-not-a-driver, modify-doesn't-re-accept-contract, paused-update-resumes). They constrain operator behavior or framing rather than forcing technical surfaces — they're already covered by inherited TRs (TR-08, TR-16, TR-21) or belong in capability KPIs/decision discipline rather than `tech-requirements.md`.
- BR-04 contributes to TR-01's framing (operator-only administration is one face of the isolation property) but is not split into a separate TR; it's cited as a secondary source on TR-01.
- BR-06 is folded into TR-05 (on-demand exports while platform healthy is the same export surface that TR-05 forces).
- BR-09 is split across TR-02 (multi-version concurrency) and TR-04 (no tenant-perceived downtime).
- BR-45 is folded into TR-02 (concurrent old/new during rollout window).

## Open Questions parked for the ADR stage
TR-06 stale link, on-demand-vs-continuous export (TR-05), versioning scheme (TR-02), numeric thresholds embedded in BRs, undefined backup standard (BR-16/TR-12), tenant-facing signal bundle (TR-03/TR-13), migration concurrency model (TR-18), drift-detection mechanism (TR-26), deeper backup-tier policy after retention window (TR-21/TR-22), explicit-anchor annotations on capability/UX headings, and a possible future tenant-facing pending-update signal.

## Hugo render validation: skipped
Per eval instructions, I did not run `cd docs && hugo server` to validate the rendered page. **The skill normally requires this step** — every `ref` resolution and the stale-source marker visibility must be confirmed in-browser before the doc is considered done. Specifically unverified:
- All `ref` source links in TR-01..TR-30 (including the BR-NN deep links, which assume `{#br-NN}` anchors are present on the BR headings in `business-requirements.md`).
- The TR-06 stale-source admonition rendering as visible warning rather than swallowed markdown.
- Frontmatter `type: docs` rendering correctly under the page-bundle parent.

## Next step for the user
> Review and edit `tech-requirements.md` before we proceed. When you're done, set `reviewed_at:` in the frontmatter to today's ISO date (or tell me to do it on your behalf). The next skill — `plan-adrs` — won't enumerate decisions until that's set. Pay particular attention to TR-06 (stale UX source link), the BR-NN anchors in `business-requirements.md` (they need explicit `{#br-NN}` on each BR heading for the deep links to resolve), and the open-questions list (some BR gaps need to flow back upstream rather than being resolved at the TR layer).
