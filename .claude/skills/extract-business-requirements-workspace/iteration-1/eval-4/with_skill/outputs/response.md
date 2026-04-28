I re-extracted `docs/content/capabilities/self-hosted-application-platform/business-requirements.md`.

**Summary of changes**
- **Preserved:** BR-01, BR-02, BR-03, BR-04 — text and source links unchanged. Numbering preserved (append-only).
- **Flagged:** BR-05's source link `user-experiences/move-off-the-platform-after-eviction.md#a-section-that-no-longer-exists` no longer resolves (the section was renamed). Per the skill, I did **not** rewrite or delete the link — I added the `> ⚠️ source no longer resolves — human review` marker above the BR and left the broken `ref` in place so the Hugo build will fail loudly until you decide how to re-source it. The likely current section is `## Edge Cases & Failure Modes` (line 103 of that UX doc), but it has no explicit `{#anchor}` yet, so re-sourcing also requires adding an anchor to the target heading.
- **Added:** no new BRs this pass — nothing in the capability or UX docs surfaced a forced requirement that wasn't already captured.
- **Open Questions:** the two existing entries (BR-05 export modality, BR-02 compatibility window) are carried forward, plus a new note tracking the BR-05 stale-source decision.

**`reviewed_at` was nulled** because the file's contents changed (the warning marker on BR-05 and the new open question). You'll need to review and reset it to today's ISO date before `define-technical-requirements` will run.

**What I need from you**
1. Review the file — especially BR-05's flag.
2. Decide whether to re-source BR-05 to `## Edge Cases & Failure Modes` (and add `{#edge-cases}` to that heading in the UX doc), pick a different section, or leave the flag for later.
3. Set `reviewed_at:` in the frontmatter to today's ISO date once you're satisfied (or tell me to record your verbal confirmation).

**Hugo render validation skipped** — I did not run `cd docs && hugo server` in this session, so I have not visually confirmed the page renders. The broken `ref` on BR-05 will, by design, fail the Hugo build until it's resolved; that's expected and is the mechanism keeping the stale source visible.
