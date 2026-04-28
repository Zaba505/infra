# Plan experiences — self-hosted-personal-media-storage

**Gate:** Capability is reviewed (`reviewed_at: 2026-04-28`, file unchanged since) — proceeding.

You said "I trust your enumeration, just go," so I enumerated the experience set against the capability's Stakeholders, Triggers, Outputs, Business Rules, and Success Criteria, applied the one-persona-one-goal discipline, and filed one issue per experience. Internal system flows (no human in the loop), public sharing (out of scope), and platform-hosting concerns (different capability) were excluded.

## Experiences filed (10)

1. **Upload a photo or file** — authorized user; manual single-file upload. Triggers + Outputs.
2. **Back up a device automatically** — authorized user; ongoing camera-roll sync. Triggers.
3. **Bulk import from an existing provider** — authorized user; migrate from Google Photos etc. Triggers + cost-avoidance/longevity outcomes.
4. **Browse and view own content** — authorized user; find/view/download/organize their own media. Triggers + Outputs.
5. **Share content with another user or group** — content owner; grant access to recipients. Triggers + Outputs + Private-by-default.
6. **Delete own content** — authorized user; remove with 30-day grace. Triggers + 30-day retention rule.
7. **Pull an on-demand archive of own content** — authorized user; self-serve export. Operator-succession rule + longevity outcome.
8. **Provision a new user** — operator; invite/onboard family or friends. Closed-user-set rule.
9. **Remove a user from the system** — operator; deprovision on departure or illegal-content termination, with 30-day window. Closed-user-set + no-illegal-content + 30-day retention rules.
10. **Take over as successor operator** — successor operator; assume control on primary unavailability without reading user content. Operator-succession rule.

## Splits / merges I made

- "Manage my content" was *not* used — split into upload, browse, share, delete (each has a distinct entry point and success state).
- "Operator terminates user for illegal content" was merged into **Remove a user** — same persona, same goal (deprovision), same 30-day retention path; the termination cause is an edge case for `define-user-experience`, not a separate journey.
- The internal "system reconciles state" / "off-site backup replication" flows are not experiences (no human in the loop) — they belong to tech design.
- "Affected party requests removal" is *not* filed: the capability's Open Questions explicitly flags affected-party recourse as undefined, and the affected party is not a stakeholder-listed actor. Surfacing rather than inventing.

## gh commands (would have run)

The exact commands are written to `gh-commands.txt` in this output directory. One per experience, ten total. Each issue:
- Title format: `story(ux): {verb-led journey} — self-hosted-personal-media-storage`
- Body uses the skill's template: Experience (user-voice sentence), Persona, Anchored-in (Stakeholders / Triggers / Outputs / Business Rules / Success Criteria), Parent capability link, Authoring pointer to `define-user-experience`, and a Related backlink.

```
gh issue create --title "story(ux): upload a photo or file — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): back up a device automatically — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): bulk import from an existing provider — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): browse and view own content — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): share content with another user or group — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): delete own content — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): pull an on-demand archive of own content — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): provision a new user — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): remove a user from the system — self-hosted-personal-media-storage" --body "..."
gh issue create --title "story(ux): take over as successor operator — self-hosted-personal-media-storage" --body "..."
```

(Full bodies, with HEREDOCs, are in `gh-commands.txt`.)

Next step per experience: invoke `define-user-experience` against the issue to produce one UX doc under `docs/content/capabilities/self-hosted-personal-media-storage/user-experiences/{slug}.md`.
