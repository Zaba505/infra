# Planned User Experiences — Self-Hosted Personal Media Storage

Source capability: `.claude/skills/plan-experiences/evals/fixtures/self-hosted-personal-media-storage.md` (reviewed 2026-04-28).

Below is the proposed set of end-to-end user experiences (UXes) implied by the capability. Each is a discrete journey a primary actor (operator, family member, or friend) initiates. Operator-only administrative journeys are included where the capability explicitly assigns them to the operator (provisioning, succession). Each will become its own GitHub issue so the doc can be authored individually via `define-user-experience`.

## Enumeration

The following UXes are derived from the capability's Triggers & Inputs, Business Rules, and Outputs sections:

1. **Manual upload** — A user uploads one or more files by hand from a device (Triggers: "single, manual upload").
2. **Automated device backup** — A user's device continuously syncs new media (e.g. phone camera roll) into the system without per-file user action (Triggers: "automated backup").
3. **Bulk import from a prior provider** — A user ingests an existing archive (e.g. a Google Photos takeout) so they can leave their previous cloud provider (Triggers: "bulk import"; supports the cost-avoidance and longevity outcomes).
4. **Browse and view own content** — A user navigates, views, and downloads their own media (Triggers: "views, downloads, or organizes their own content"; supports the durable-retrievability output).
5. **Organize own content** — A user arranges their media into albums/folders/tags as a precondition for sharing and for finding things later (Triggers: "organizes their own content"). Splitting this from view is justified because organization mutates state while viewing does not, and shared groups/albums are explicitly a sharing primitive in Business Rules.
6. **Share content with another user or group** — A content owner grants explicit access to a named recipient or shared group (Triggers: "shares content"; Business Rules: "private by default", one-to-one or shared group).
7. **Access content shared with me** — A recipient discovers, views, and downloads content others have shared with them (Outputs: "shared content becomes retrievable by the explicitly named recipients").
8. **Delete own content** — A user removes their own content, with the 30-day retention window providing accident recovery (Triggers: "deletes their own content"; Business Rules: 30-day retention).
9. **Pull a personal archive export** — A user pulls a complete on-demand export of their own content without operator involvement, as the longevity safety net (Business Rules: "operator succession" mechanism (a)).
10. **Operator provisions or removes a user** — The operator invites a new user or removes an existing one (Business Rules: "closed user set"; Preconditions: "explicitly invited and provisioned"). When a user is removed, the 30-day retention applies to their data.
11. **Successor operator takeover** — The designated successor uses sealed/escrowed credentials to assume operation of the system when the primary operator is unavailable (Business Rules: "operator succession" mechanism (b)).

### Considered and intentionally excluded

- **Account/credential recovery** — Explicitly Out of Scope ("lost credentials = lost data").
- **Public sharing** — Explicitly Out of Scope.
- **Collaborative editing / commenting** — Explicitly Out of Scope.
- **Affected-party (depicted subject) takedown request** — Listed as an Open Question in the capability; not yet a defined journey. Hold until the capability resolves the open question.
- **Operator content moderation / termination for illegal content** — The capability gives the operator the right to terminate access on credible evidence, but not a defined journey for inspecting content (operator cannot read user content). This collapses into UX 10 (provision/remove) rather than being a distinct UX.
- **Per-user quota management** — Excluded by Business Rules ("no storage quotas").

### Open questions for the human before filing issues

1. Should "Browse/view own content" (UX 4) and "Organize own content" (UX 5) be one UX or two? Splitting matches the read/write asymmetry; merging keeps the surface smaller. Recommendation: keep split.
2. Should "Manual upload" (UX 1) and "Automated device backup" (UX 2) be one UX with two modes, or two? They share the same end-state (durable storage of one user's content) but differ sharply in initiator (human vs. background agent) and failure modes (interactive error vs. silent retry). Recommendation: keep split.
3. "Successor operator takeover" (UX 11) is a once-in-a-lifetime event, not a recurring user journey. Confirm it should be modeled as a UX rather than as an operations runbook artifact. Recommendation: model as a UX because the capability explicitly defines its trigger, actor, and outputs.
4. Should the "operator provisions a user" and "operator removes a user" be one UX or two (UX 10)? They share an actor and a system but the removal path triggers the 30-day retention sweep. Recommendation: keep as one UX with two flows; split later if the flow doc grows unwieldy.

These four questions should be confirmed with the human before filing the issues below.

## GitHub issues to file (one per confirmed UX)

The commands that would be run are saved alongside this file as `gh-commands.txt`. They are not executed.
