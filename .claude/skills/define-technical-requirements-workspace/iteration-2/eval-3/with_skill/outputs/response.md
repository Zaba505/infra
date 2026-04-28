I'm going to push back on this one before touching `tech-requirements.md`.

What you've proposed — "use Postgres for tenant metadata storage with logical replication for cross-region failover" — is a Stage-2 decision, not a TR. Postgres is one of several stores that could hold tenant metadata, and logical replication is one of several mechanisms that could give cross-region failover. If I file that as TR-08, I erase the alternatives that the eventual ADR is supposed to weigh, and I break the rule that TRs translate constraints without choosing technology.

The TR underneath is a constraint pair: there must be tenant metadata storage with defined query/consistency properties, and it must survive a single-region failure within a bounded recovery objective. But before I write that, I need to source it — every TR has to cite a BR (or a prior shared ADR / inherited constraint).

Walking the BR list, I don't see a clean source for this:

- **BR-13** says the platform must provide compute, persistent storage, and network reachability *to each tenant*. That's tenant-facing storage for tenant data, not platform-internal tenant metadata.
- **BR-16** demands backup and DR for *tenant data*, again not platform-internal metadata.
- **BR-2** (rebuildable from definitions) and **BR-42** (successor can run standup identically) imply the platform's internal state is either reproducible from definitions or recoverable, but neither names "tenant metadata" or a cross-region failover RPO/RTO.
- No BR I can find establishes a cross-region availability target for platform-internal state.

So I think one of two things is true, and I'd like you to pick:

1. **There's a missing BR.** Something like "the platform must remain operational for hosted tenants through the loss of a single hosting region" or "the operator's view of which tenants exist and their state must survive a single-region failure within {N} minutes." If that's the real demand, route back to `extract-business-requirements` and add it; once it's reviewed I can derive the TR cleanly. (Note: the platform's stated availability posture in the BRs is fairly modest — BR-8 explicitly says tenants needing stronger guarantees than the platform offers must host elsewhere — so a cross-region failover BR would be a real expansion of the platform's contract, worth the operator's deliberate sign-off.)

2. **It's actually scaffolding for an ADR**, not a TR at all. In which case I'd capture "consider Postgres + logical replication for tenant metadata cross-region failover" as an Open Question on `tech-requirements.md` so it surfaces during `plan-adrs`, and we'd let the underlying TR (whatever it ends up being) get written when the BR exists.

I have **not** modified `tech-requirements.md`. The current file still has TR-01 through TR-07, `reviewed_at: 2026-04-20`. Once you tell me which path to take — add a BR, or park as an open question — I'll proceed.

A couple of incidental notes from reading the live files:

- **TR-06's source link is already broken** (`migrate-existing-data.md#a-section-that-no-longer-exists`). Per the skill, that should carry a `> ⚠️ source no longer resolves — human review` flag. I haven't touched it because you asked me to add TR-08, not re-extract, but it's worth fixing on the next pass.
- The existing TRs use raw markdown paths rather than Hugo `{{< ref >}}` shortcodes. The skill calls for `ref` on every internal link. Again, scope of this turn is just TR-08, but flagging it.
- **Hugo render validation skipped** as instructed for this eval — no `hugo server` was run.
