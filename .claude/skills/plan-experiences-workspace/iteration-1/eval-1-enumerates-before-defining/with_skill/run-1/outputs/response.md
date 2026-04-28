Capability is reviewed (`reviewed_at: 2026-04-28`, file unchanged since) — proceeding.

I read `self-hosted-personal-media-storage.md` end-to-end. Below is a proposed list of user experiences for the capability, anchored back to specific sections (Stakeholders, Triggers, Outputs, Business Rules, Success Criteria). No journey steps, no UI — just the journey set, so we can confirm before filing issues.

## Proposed user experiences

Stakeholders the capability sanctions:
- **Operator** (Carson, sole accountable party)
- **Authorized user** (operator, family, friend — any primary actor)
- **Successor operator** (named in Business Rules → Operator succession)
- *Affected parties* (people depicted in shared media) — listed but with no system access; no UX is filed for them on their own. The "object to a shared photo" pathway is an Open Question and is not yet experience-worthy until the capability decides what (if any) recourse exists.

Proposed journeys (verb-led, one persona + one goal each):

1. **Upload a photo or file** — *authorized user* getting a single piece of media into their durable store. Anchored in: Triggers ("a user uploads a file"), Outputs ("durably stored and retrievable"), Business Rules ("private by default", "no storage quotas").

2. **Auto-back up a device's media** — *authorized user* keeping their phone's camera roll continuously mirrored without manual action. Distinct from #1: different entry point (background sync, not a deliberate "I'm uploading this"), different success state (steady-state confidence rather than a single-file confirmation). Anchored in: Triggers ("automated backup, e.g. phone camera roll sync"), Outputs (durable storage), Success Criteria ("zero data loss").

3. **Bulk-import from an existing provider** — *authorized user* moving an existing Google Photos / iCloud / etc. archive in. Distinct from #1 and #2: one-shot, large, with provider-specific metadata to preserve. Anchored in: Triggers ("bulk import from an existing provider, e.g. Google Photos takeout"), Purpose & Outcome (cost avoidance, longevity).

4. **Browse and view my own content** — *authorized user* finding and viewing the content they've already stored. Anchored in: Triggers ("a user views, downloads, or organizes their own content"), Outputs ("durably stored and retrievable by them"). Distinct success state from upload — the user is consuming, not producing.

5. **Share content with another user or shared group** — *authorized user* (acting as content owner) granting access to a named recipient or shared album (e.g. "family album"). Anchored in: Triggers ("a content owner shares content with another user or a shared group"), Business Rules ("private by default", sharing is the explicit override), Outputs ("shared content becomes retrievable by the explicitly named recipients").

6. **View content shared with me** — *authorized user* (acting as recipient) discovering and viewing media someone else shared. Distinct from #4 (own content) and from #5 (the share action itself): different entry point ("someone gave me access"), different success state (consuming someone else's content within the privacy posture). Anchored in: Stakeholders ("Secondary actors / consumers"), Outputs ("shared content becomes retrievable by the explicitly named recipients").

7. **Delete my own content** — *authorized user* removing content they own, with the 30-day retention window in mind. Anchored in: Triggers ("a user deletes their own content"), Business Rules ("30-day retention after deletion / departure"), Outputs ("removed from view immediately and purged after the retention window").

8. **Export a complete archive of my own content** — *authorized user* pulling a self-service archive of everything they own, on demand, without operator involvement. Anchored in: Business Rules ("Operator succession" — every user can pull a complete on-demand archive), Purpose & Outcome (longevity). Distinct from #4 (browse/view) — this is the disaster-prep / vendor-independence journey.

9. **Operator provisions or removes a user** — *operator* adding a family member or friend (or removing one). Anchored in: Stakeholders (operator is sole accountable party), Business Rules ("Closed user set — only the operator can add or remove users; no public sign-up"), Preconditions ("the user has been explicitly invited and provisioned").

10. **Operator terminates a user for an illegal-content violation** — *operator* revoking a user's access on credible evidence. Distinct from #9 (routine removal): the trigger is adversarial, the operator cannot inspect content directly, and the journey has to reconcile that constraint. Anchored in: Business Rules ("No illegal content" — operator reserves the right to terminate on credible evidence, even though the operator cannot inspect user content directly).

11. **Successor operator takes over a system whose primary operator is unavailable** — *successor operator* breaking the seal on escrowed credentials and assuming operation. Anchored in: Business Rules ("Operator succession" — sealed/escrowed credentials, takeover is a discrete event, successor's access is to operate the system, not to read user content). The privacy-posture invariants (private by default, lost credentials = lost data) must hold across the takeover.

## Journeys I considered and deliberately did not include

- **"User recovers lost credentials"** — explicitly **Out of Scope** ("Account/credential recovery") and reinforced by the Business Rule "lost credentials = lost data". Not an experience under this capability by design.

- **"Affected party requests removal of a photo they appear in"** — listed under **Open Questions** ("Affected-party recourse"). The capability has not yet decided whether such a process exists, who initiates it, or how the affected party — who may not be a user — would even reach the system. Filing a UX issue now would force `define-user-experience` to invent capability rules. Recommend resolving the Open Question (back via `define-capability`) before adding this.

- **"Operator stands up the hosting platform"** — belongs to the *self-hosted-application-platform* capability (the capability doc lists "general-purpose self-hosted application hosting" as **Out of Scope** here). Flag-and-stop: not filed under this capability.

- **"Operator monitors capacity / plans storage"** — capacity planning is named as the operator's problem in Business Rules ("No storage quotas"), but this is operational/runbook concern, not a capability-scoped user journey grounded in Triggers / Outputs. If we want it as a UX, we'd need the capability to sanction it as such; otherwise it lives in operator runbooks.

- **"User changes who they've shared content with / revokes a share"** — arguably a sub-flow of #5 (Share content). I left it folded into #5 because the capability doesn't separate "share" from "unshare" in Triggers, and the success state of a share UX must already include the share being revocable (otherwise "private by default" is one-way). If you'd rather split it into its own UX, say so and I'll add a #12.

## Open questions I want your call on before filing

- **#6 vs. #5 split.** I split "share content" (owner-side) from "view shared content" (recipient-side) because they have different entry points and different personas-in-role. Reasonable to merge if you'd rather one UX cover both sides of the seam — but I'd push back: the recipient's success state ("I can find and view what was shared with me") and the owner's success state ("the recipient I named can now see this") are different system properties.

- **#10 (operator terminates).** This is a thin journey — most of the work is policy (credible evidence, due process) rather than UX. If you'd rather treat it as a runbook entry rather than a UX, drop it from the list. I included it because the capability's business rule explicitly carves out an operator action, and there's a human-in-the-loop decision to document.

- **Order.** I'd suggest tackling #1 (upload) first — it's the canonical entry to the system and several other UXes (#4 browse, #5 share, #7 delete, #8 export) take its success state as their precondition. #2 (auto-backup) and #3 (bulk import) can come after #1 since they share its success state. #9 (provision) is a precondition for *all* of the above and could also reasonably go first. Tell me which ordering you want and I'll file in that order.

## Once you approve

On approval I'll file 11 GitHub issues — one per experience above — using `gh issue create`. Each issue body links back to this capability, names the persona and goal in the user's voice, anchors the journey in specific capability sections, and points at `define-user-experience` as the skill that authors the UX doc under `docs/content/capabilities/self-hosted-personal-media-storage/user-experiences/{ux-name}.md`.

The exact `gh issue create` commands I would run are listed in `gh-commands.txt` alongside this response (not yet executed — awaiting your approval).

Reply with:
- "approve all" to file all 11,
- "approve {numbers}" to file a subset,
- or any add / remove / rename / merge / split / reorder you want first.
