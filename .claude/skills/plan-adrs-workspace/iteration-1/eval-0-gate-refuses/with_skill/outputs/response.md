Stopping — I won't enumerate ADRs yet.

`docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` shows `reviewed_at: null` in its frontmatter. The `plan-adrs` skill's review gate requires `reviewed_at` to be set to an ISO date newer than the file's last modification before any ADR planning happens. Planning ADRs against an unreviewed TR list risks sourcing decisions from drafts that may still change.

To unblock this:

1. Read the current contents of `tech-requirements.md` end-to-end and edit anything that needs editing.
2. Set the frontmatter `reviewed_at:` to today's ISO date (`2026-04-27`).
3. Re-invoke `plan-adrs`.

Alternative: if you've already reviewed the file and just didn't update the field, tell me explicitly ("I've reviewed it, record my verbal confirmation") and I'll set `reviewed_at: 2026-04-27` on your behalf, then proceed.

No issues filed. No `gh-invocations.txt` written — the skill stops at the gate before any decisions are enumerated or any GitHub issues are drafted.
