# Tech design composed

I composed the tech design at
`docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md`.

## How I treated the three ADRs

Only ADR-0001 (tenant state storage) is `accepted`. ADR-0002 (contract
versioning) and ADR-0003 (tenant eviction export) are still `proposed`. A
tech-design document is supposed to describe the intended final state, so
presenting proposed decisions as if they were settled would mislead the
reader. I took the middle path:

- The component map and the per-requirement sections **mark proposed
  decisions explicitly** (status badge + "*(proposed)*" inline) and note
  that the alternatives in the ADR remain on the table.
- An **Open items** section flags everything that still needs to land
  before this design is fully committed, including the two proposed ADRs
  and the two requirements (TR-03, TR-06) that have no ADR yet.
- A **Requirement → ADR coverage** table makes the gaps visible at a
  glance.

If you'd rather I write the document as if 0002 and 0003 were already
accepted, say the word and I'll re-compose — but I'd recommend accepting
those ADRs first and then re-running this composition, so the design
document and the ADRs stay in sync.

## Files written

- `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (fixture)
- `docs/content/capabilities/self-hosted-application-platform/adrs/_index.md` (fixture)
- `docs/content/capabilities/self-hosted-application-platform/adrs/0001-tenant-state-storage.md` (fixture)
- `docs/content/capabilities/self-hosted-application-platform/adrs/0002-contract-versioning.md` (fixture)
- `docs/content/capabilities/self-hosted-application-platform/adrs/0003-tenant-eviction-export.md` (fixture)
- `docs/content/capabilities/self-hosted-application-platform/tech-design/_index.md` (composed)
