# Response

Quick heads-up before I made the edit: the slot you asked for (TR-04) was already taken by the operator-initiated-tenant-update requirement, and the thing you want to add ("use Postgres with logical replication for cross-region failover") reads to me more like a design decision than a technical requirement. I went ahead and made the change you asked for, but want to flag the framing so you can decide whether to refine it.

## What I did

1. Inserted a new **TR-04: Tenant metadata storage uses Postgres with logical replication for cross-region failover** at the requested slot.
2. Renumbered the existing TR-04..TR-07 down by one (now TR-05..TR-08). Their content is unchanged.
3. Added a short inline note on the new TR-04 calling out that it encodes a tech choice rather than a capability-driven need, and recommending the underlying requirement get captured separately with the choice moved to an ADR.

Updated file: `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`

## Why I'd push back on the framing

A technical requirement should describe *what must be true* about the system, sourced from a capability rule, business rule, or UX. "Use Postgres" and "use logical replication" are *how* you intend to satisfy something. The actual requirements hiding behind that sentence are probably:

- **Availability / RTO**: "Tenant metadata must remain available across the loss of a single region within N minutes."
- **Durability / RPO**: "On regional failover, no more than N seconds of committed tenant-metadata writes may be lost."
- **Source**: which is this — a capability business rule? a UX (operator-initiated failover, disaster recovery)? a prior shared decision? Right now it has none, which is the same smell as "use Postgres" — there is no traceable *why*.

If you split it that way, then:

- The **requirement** survives even if you later swap Postgres for CockroachDB, Spanner, or anything else.
- The **decision** ("Postgres + logical replication, primary in region A, warm standby in region B, manual promotion") lives in an ADR under this capability, where its alternatives, trade-offs, and consequences can be argued properly. Logical replication in particular has well-known sharp edges (no DDL replication, sequence handling, slot bloat on standby outage, conflict handling on failback) that deserve an ADR's "Consequences" section, not a one-line requirement.

## Recommended next step

Want me to:

1. Replace the new TR-04 with an availability/RPO-style requirement (you give me the numbers), and
2. Draft an ADR under this capability capturing the Postgres + logical-replication choice?

That keeps the requirements doc clean and gives the technology choice the scrutiny it deserves.
