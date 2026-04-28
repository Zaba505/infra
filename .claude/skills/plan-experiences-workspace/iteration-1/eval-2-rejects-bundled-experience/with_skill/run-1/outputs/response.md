**Gate check:** Capability `self-hosted-personal-media-storage` is reviewed (`reviewed_at: 2026-04-28`, file unchanged since) — proceeding.

**However, I'm not going to file a single experience called "user manages their content."** Let's split it.

That title bundles at least four distinct journeys, and the skill I'm running explicitly resists bundles because they produce UX docs with mismatched entry points, mismatched success states, and contradictory edge cases. Looking at the capability doc's *Triggers & Inputs* and *Outputs & Deliverables* sections, here's what's actually in there:

1. **Upload content** — Primary actor: any authorized user (operator, family, friend). Goal: get a new photo/video/file into the system. Anchored in Triggers ("a user uploads a file", "automated backup", "bulk import from existing provider") and Outputs ("the user's content is durably stored and retrievable by them"). Entry point: "I just took a photo / I have a takeout archive." Success: content is durably stored and visible to the uploader.

2. **Share content with another user or group** — Primary actor: a content owner. Goal: grant explicit access to a specific recipient or shared group. Anchored in Triggers ("a content owner shares content with another user or a shared group"), Business Rules ("Private by default", sharing one-to-one or via shared group), and Outputs ("Shared content becomes retrievable by the explicitly named recipients"). Entry point: "I want this person to see this." Success: the named recipient can retrieve it; nobody else can.

3. **Browse / view their own content** — Primary actor: any authorized user. Goal: find and view what they've already stored (and what's been shared with them). Anchored in Triggers ("a user views, downloads, or organizes their own content") and the Outcome of *Longevity* ("content remains accessible and intact over long time horizons"). Entry point: "I want to look at memories" or "I want to find that file from last summer." Success: the user locates and views the content.

4. **Delete content** — Primary actor: a content owner. Goal: get content out of the system. Anchored in Triggers ("a user deletes their own content"), Business Rules ("30-day retention after deletion"), and Outputs ("Content the user deleted is removed from their view immediately and purged from the system after the retention window"). Entry point: "I want this gone." Success: it's gone from the user's view immediately, and purged after 30 days.

Each of those has a different entry point and a different success state. Bundling them into one UX doc would force four mismatched flows into a single document — exactly what `define-user-experience` is shaped to avoid.

A couple of additional notes on the proposed split, before you approve:

- **"Browse" was not in your original four** — you said upload, share, *and* delete. I added browsing because the capability's Triggers section explicitly names "a user views, downloads, or organizes their own content" as a distinct trigger, and because viewing shared content is the success state of journey 2. If you'd rather drop browsing for now and revisit later, say so.
- **There are other capability-implied journeys I have *not* added** because you didn't ask for them, but flagging in case you want them in the same planning pass: *automated device backup* (distinct from manual upload — different entry point: the phone, not the user), *bulk import from another provider* (Google Photos takeout — distinct entry point and success state), *pull an on-demand archive of one's own content* (Business Rules: operator-succession / longevity guarantee), *operator provisions a new user* (Business Rules: closed user set), and *operator removes a user / handles departure* (Business Rules: 30-day retention on departure). These are real journeys grounded in the capability doc; whether they're in scope for *this* planning pass is your call.

**I have not filed any issues yet.** Per the skill, no `gh issue create` runs until you approve a confirmed list.

Please confirm one of:

- (a) File four issues for the split above (upload, share, browse, delete) — that's the literal expansion of your request.
- (b) File three issues — drop browsing, keep upload/share/delete as you originally said.
- (c) Expand the list to include some/all of the additional journeys I flagged (backup, bulk import, on-demand archive, provisioning, departure). Tell me which.
- (d) Something else — rename, reorder, merge any pair you think really is one journey.

Once you pick, I'll print the `gh issue create` commands (one per approved experience), each with a body anchored back to the specific capability sections above and pointing at `define-user-experience` as the next-step skill.
