# Response

I'm going to push back on adding this as written.

"Use Postgres for tenant metadata storage with logical replication for cross-region failover" is a **decision**, not a technical requirement. It picks a specific technology (Postgres) and a specific replication mechanism (logical replication) out of several options that could all satisfy some underlying constraint. That kind of choice belongs in an ADR (`plan-adrs` / `define-adr`), not in `tech-requirements.md`. If we let it in here, we erase the alternatives the ADR is supposed to weigh — and any later ADR citing this TR would be circular ("we chose Postgres because TR-04 says use Postgres").

A second issue: TR-04 is already taken. The current file has TR-01 through TR-07 populated. TR identifiers are append-only forever, so a new requirement extracted from this conversation would be TR-08, not TR-04.

A third issue: I can't actually extract a real TR from your message, because you didn't cite a source in the capability doc or any UX doc. Every TR in this file must link back to something that *forces* it. The skill is "extract, don't design" — if nothing in the capability or a UX implies the constraint, it isn't a requirement yet.

## What I'd need to turn this into a real TR

The requirement *underneath* your decision is probably something like:

> Tenant metadata must be queryable by tenant ID with strong consistency, and must survive a single-region failure within {N} minutes / with {M} bound on data loss.

But I don't want to invent that. To add it properly we'd need:

1. **A source.** Which UX or capability business rule demands strong-consistent tenant-metadata lookups? Which one demands cross-region survival, and on what RTO/RPO budget? Quote the line.
2. **The actual constraint, not the mechanism.** "Survives single-region failure within N minutes with at most M seconds of data loss" is a requirement. "Logical replication" is one of several ways to achieve it.

If you can point me at the source UX/section, I'll append it as TR-08 (or higher) on a re-extraction pass.

## What I'm doing with your suggestion

I'm parking it as an open question for the ADR stage, exactly as the skill prescribes. I have **not** modified `tech-requirements.md` — no TR was added, nothing was renumbered, and `reviewed_at` is unchanged.

Open question to carry into `plan-adrs`:

> Consider Postgres + logical replication for tenant metadata storage with cross-region failover. Pending: identification of the underlying TR (consistency model + RTO/RPO) sourced from the capability or a UX.

## Next step for you

Either:
- Point me at the UX/capability section that forces a consistency + cross-region-survival constraint on tenant metadata, and I'll extract a properly-sourced TR-08; or
- Leave the Postgres idea parked as an open question and we'll revisit it in `plan-adrs` once the underlying TR exists.

Whichever you pick, `tech-requirements.md` is unchanged, and `reviewed_at` stays at its current value.
