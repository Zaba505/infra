# Plan: User Experiences for `self-hosted-personal-media-storage`

I read `docs/content/capabilities/self-hosted-personal-media-storage.md` and walked through its triggers, stakeholders, business rules, and out-of-scope list. Below is the proposed list of end-to-end user experiences (UXes) that the capability implies. Each is one journey a real human (operator or invited user) drives end-to-end. I have intentionally kept these at the "user journey" level — not screens, not technical components — so each can later be authored individually with `define-user-experience`.

## Proposed User Experiences

Each entry: **slug** — primary actor — one-line summary, anchored to capability sections.

1. **`provision-a-user`** — Operator — The operator invites a new family member or friend, issues them credentials, and the new user completes first-time setup. (Triggers: provisioning precondition; Business Rules: closed user set, lost credentials = lost data.)

2. **`upload-content-manually`** — Authenticated user — A user picks one or more files from a device and uploads them as their own private content. (Triggers: "single, manual upload"; Outputs: durably stored and retrievable.)

3. **`automated-device-backup`** — Authenticated user — A user configures their phone/device to continuously sync new captures (e.g. camera roll) into the system in the background. (Triggers: "automated backup".)

4. **`bulk-import-from-existing-provider`** — Authenticated user — A user imports a large takeout/export archive (e.g. Google Photos takeout) preserving capture metadata. (Triggers: "bulk import".)

5. **`browse-and-view-own-content`** — Authenticated user — A user navigates and views/streams the content they own, including organizing it (albums, ordering). (Triggers: "views, downloads, or organizes their own content".)

6. **`share-content-with-another-user-or-group`** — Content owner — A user grants explicit access to a piece of content (or album) to a named recipient or shared group; the recipient sees it appear in their view. (Triggers: sharing; Business Rules: private by default, no public sharing.)

7. **`receive-and-view-shared-content`** — Recipient user — A user finds, browses, and views content that was shared with them, distinguishable from their own. (Stakeholders: secondary actors / consumers.)

8. **`delete-own-content-and-30-day-recovery`** — Authenticated user — A user deletes content; it disappears from their view immediately and remains recoverable for 30 days, after which it is purged. (Business Rules: 30-day retention.)

9. **`pull-on-demand-archive-export`** — Authenticated user — A user pulls a complete export of all their own content without operator involvement, while the system is healthy. (Business Rules: operator succession (a) — user-driven export.)

10. **`offboard-a-user`** — Operator + departing user — The operator removes a user (voluntary departure, illegal-content termination, or similar); the departing user's data is retained 30 days then purged. (Business Rules: closed user set, no illegal content, 30-day retention.)

11. **`operator-succession-takeover`** — Successor operator — A designated successor uses sealed credentials and the runbook to take over running the system after the primary operator becomes unavailable, without gaining the ability to read user content. (Business Rules: operator succession (b).)

## Experiences I considered and dropped

- **Self-signup / account recovery flow** — explicitly out of scope (closed user set; lost credentials = lost data).
- **Public sharing flow** — explicitly out of scope.
- **Collaborative editing** — explicitly out of scope.
- **Affected-party content takedown request** — capability lists this as an Open Question, not a confirmed UX. Worth flagging for the operator to confirm before adding as a UX.
- **Off-site backup as a user-facing journey** — off-site replication is an operational/durability mechanism, not a journey a user drives. Belongs in tech design, not UX.
- **Operator capacity planning** — explicitly the operator's problem (no quotas), but it is operational work, not an end-user journey.
- **Browse other users' libraries / operator content inspection** — explicitly forbidden by "private by default" and "operator cannot inspect content".

## Open questions surfaced while planning

These are *new* questions raised by the planning step (separate from the three Open Questions already in the capability doc):

- Should `automated-device-backup` and `upload-content-manually` collapse into one UX? They differ enough in initiation (foreground vs. background, device-level config vs. one-shot) that I kept them separate, but the operator may want them merged.
- Is `receive-and-view-shared-content` distinct enough from `browse-and-view-own-content` to warrant its own UX? I kept it separate because the discovery surface and provenance (who shared this with me) is meaningfully different.
- Should `pull-on-demand-archive-export` include scheduled/recurring exports, or is a recurring schedule a separate UX? The capability doc mentions users "may schedule periodic pulls" — likely the same UX, but worth confirming.
- Is `operator-succession-takeover` a "user experience" or operational runbook? It has a clear actor and journey, so I included it; the operator may prefer to document it elsewhere.

## Next step

Once the operator confirms (or edits) the list above, file one GitHub issue per confirmed experience so each can be authored with `define-user-experience`. I have **not** run `gh issue create` — per the task instructions, the commands I would have run are listed below and saved to `gh-commands.txt`.

### `gh issue create` commands (not executed)

```
gh issue create \
  --title "story: define UX 'provision-a-user' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'provision-a-user' under capability self-hosted-personal-media-storage.

Actor: Operator (with new invited user completing first-time setup).
Summary: The operator invites a new family member or friend, issues them credentials, and the new user completes first-time setup.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: closed user set; lost credentials = lost data; provisioning precondition.

Use the 'define-user-experience' skill. Output: docs/content/capabilities/self-hosted-personal-media-storage/user-experiences/provision-a-user.md"

gh issue create \
  --title "story: define UX 'upload-content-manually' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'upload-content-manually' under capability self-hosted-personal-media-storage.

Actor: Authenticated user.
Summary: A user picks one or more files from a device and uploads them as their own private content.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: triggers - single manual upload; outputs - durably stored and retrievable.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'automated-device-backup' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'automated-device-backup'.

Actor: Authenticated user.
Summary: A user configures their device to continuously sync new captures (e.g. camera roll) into the system in the background.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: triggers - automated backup.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'bulk-import-from-existing-provider' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'bulk-import-from-existing-provider'.

Actor: Authenticated user.
Summary: A user imports a large takeout/export archive (e.g. Google Photos takeout) preserving capture metadata.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: triggers - bulk import.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'browse-and-view-own-content' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'browse-and-view-own-content'.

Actor: Authenticated user.
Summary: A user navigates and views/streams the content they own, including organizing it (albums, ordering).
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: triggers - views, downloads, or organizes their own content.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'share-content-with-another-user-or-group' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'share-content-with-another-user-or-group'.

Actor: Content owner.
Summary: A user grants explicit access to a piece of content (or album) to a named recipient or shared group.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: triggers - sharing; business rules - private by default; out of scope - public sharing.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'receive-and-view-shared-content' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'receive-and-view-shared-content'.

Actor: Recipient user.
Summary: A user finds, browses, and views content that was shared with them, distinguishable from their own.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: stakeholders - secondary actors / consumers.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'delete-own-content-and-30-day-recovery' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'delete-own-content-and-30-day-recovery'.

Actor: Authenticated user.
Summary: A user deletes content; it disappears immediately and is recoverable for 30 days, then purged.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: business rules - 30-day retention after deletion.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'pull-on-demand-archive-export' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'pull-on-demand-archive-export'.

Actor: Authenticated user.
Summary: A user pulls a complete export of all their own content without operator involvement.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: business rules - operator succession (a) user-driven export.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'offboard-a-user' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'offboard-a-user'.

Actor: Operator + departing user.
Summary: The operator removes a user; the departing user's data is retained 30 days then purged.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: closed user set; no illegal content; 30-day retention.

Use the 'define-user-experience' skill."

gh issue create \
  --title "story: define UX 'operator-succession-takeover' for self-hosted-personal-media-storage" \
  --body "Author the user experience document for 'operator-succession-takeover'.

Actor: Successor operator.
Summary: A designated successor uses sealed credentials and the runbook to take over running the system, without gaining the ability to read user content.
Capability: docs/content/capabilities/self-hosted-personal-media-storage.md
Anchors: business rules - operator succession (b).

Use the 'define-user-experience' skill."
```

Confirm the list (or tell me what to drop / add / merge) and I will then file the issues for real.
