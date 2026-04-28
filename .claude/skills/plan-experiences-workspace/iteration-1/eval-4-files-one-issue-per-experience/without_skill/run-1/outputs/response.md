# Plan: User Experiences for `self-hosted-personal-media-storage`

Source capability doc: `.claude/skills/plan-experiences/evals/fixtures/self-hosted-personal-media-storage.md` (reviewed 2026-04-28, unchanged since).

I enumerated experiences directly from the capability's Triggers, Stakeholders, and Business Rules sections. Each trigger or operator/successor responsibility that involves a distinct human journey gets its own UX. Pure system rules (e.g. off-site backup, "private by default" as a policy) are not experiences and are intentionally absent.

## Enumerated experiences (12)

End-user content journeys:
1. **Upload a single file** — manual, one-shot upload (Triggers §1).
2. **Automated device backup** — unattended sync from a device (Triggers §2).
3. **Bulk import from existing provider** — e.g. Google Photos takeout (Triggers §3).
4. **Browse, view, and organize own content** — viewing/downloading/organizing one's own media (Triggers §4).
5. **Share content with another user or group** — owner-initiated sharing, 1:1 or group (Triggers §5).
6. **View content shared with me** — recipient-side counterpart to share.
7. **Delete own content** — including the 30-day retention window (Triggers §6 + retention rule).
8. **Pull on-demand archive of own content** — the user-side longevity mechanism, including scheduled periodic pulls (Operator-succession rule).

Operator and lifecycle journeys:
9. **Operator provisions or removes a user** — closed user set + retention rule.
10. **Operator terminates access on illegal-content evidence** — "no illegal content" rule, given the operator cannot inspect content directly.
11. **Successor operator takes over** — sealed/escrowed credential handoff (Operator-succession rule).
12. **User leaves the system** — voluntary departure, final export, 30-day purge.

## Deliberate non-inclusions

- **Public sharing**, **collaborative editing**, **monetization**, **self-signup**, and **credential recovery** are explicitly Out of Scope in the capability and were not turned into UX issues.
- **Off-site backup** is an operational property, not a user journey.
- **"Affected-party recourse"** (capability open question §2) is not yet a defined journey — it is flagged inside the operator-termination UX issue rather than getting its own UX, since the capability hasn't decided whether such a process exists.
- **Discovering new shared content** is folded into "view content shared with me" rather than split out, to avoid pre-deciding notification UX.

## What I filed

12 GitHub issues, one per experience, each labeled `enhancement` and titled in the repo's `story(scope): description` style. Each body identifies the source trigger or rule, the scope, in/out of scope, inputs/outputs, and a "done when" pointing to `define-user-experience`.

Per the instructions I did **not** invoke `gh issue create`. The exact commands are saved alongside this file as `gh-commands.txt` and can be executed directly.

## Files

- Plan/response: `.claude/skills/plan-experiences-workspace/iteration-1/eval-4-files-one-issue-per-experience/without_skill/run-1/outputs/response.md`
- gh commands: `.claude/skills/plan-experiences-workspace/iteration-1/eval-4-files-one-issue-per-experience/without_skill/run-1/outputs/gh-commands.txt`
