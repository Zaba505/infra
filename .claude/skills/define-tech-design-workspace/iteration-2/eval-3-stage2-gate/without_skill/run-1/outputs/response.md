Great — `tech-requirements.md` is in place with TR-01 through TR-07. Before we start drafting ADRs, I want to flag a few things and then suggest a sequence so we don't end up with a tangle of half-decided documents.

## A couple of observations on the tech-requirements doc

Looking at the file, I think it's a solid first pass but worth a quick review pass before we start citing TR numbers in ADRs:

1. **`reviewed_at: null` in the frontmatter.** Once you've done a read-through and made any edits you want, set this to today's date (2026-04-26). Treat the TR list as append-only from that point — once an ADR cites `TR-04`, renumbering breaks provenance.
2. **TR-07 ("all inter-service comms traverse Cloudflare → GCP")** is sourced from "prior shared decision" rather than the capability or a UX. That's fine, but it's worth linking to the actual repo-level architecture note (CLAUDE.md or a root ADR) so future readers can trace it.
3. **Each TR could use a one-line "why this is a requirement, not a decision"** — i.e. what makes it non-negotiable at this layer. That'll save us re-litigating it inside ADRs later.

If you agree, I'd tighten those three things first, then move on.

## Proposed ADR sequence

Each TR doesn't need its own ADR — several decisions cut across multiple TRs. Here's how I'd group them, ordered so earlier ADRs unblock later ones:

| # | ADR | Addresses | Why first |
|---|-----|-----------|-----------|
| 1 | **Tenant isolation model** (namespace-per-tenant vs cluster-per-tenant vs VM-per-tenant) | TR-01, TR-03 | Foundational; shapes data, network, observability decisions |
| 2 | **Tenant identity & data scoping primitive** (how a tenant ID is propagated and enforced at every read/write) | TR-01, TR-03, TR-05, TR-06 | Every other ADR assumes this primitive exists |
| 3 | **Platform-contract versioning & coexistence** (semver? N and N-1? how long is the migration window?) | TR-02 | Drives rollout mechanics; needs to settle before we design the rollout pipeline |
| 4 | **Zero-downtime tenant update mechanism** (blue/green, rolling, surge — and how it interacts with the contract version above) | TR-04 | Depends on ADR-3 |
| 5 | **Tenant data import/export format & integrity protocol** (one ADR covering both directions, since they share the portable-format and integrity-check decisions) | TR-05, TR-06 | Independent of 1–4; can be drafted in parallel |
| 6 | **Per-tenant observability storage & query path** (how isolation from ADR-1 is enforced specifically for metrics/logs/traces) | TR-03 | Depends on ADR-1 and ADR-2 |

TR-07 (Cloudflare → GCP topology) is an inherited constraint, not something we should re-decide here — I'd reference it in each ADR's context rather than give it its own ADR.

## Where the ADRs should live

Two reasonable options:

- **`docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/`** — capability-scoped, doesn't pollute the global ADR space. Good if these decisions only bind this capability.
- **`docs/content/r&d/adrs/`** (alongside 0001–0007) — global, numbered into the repo-wide sequence. Good if any of these decisions will bind other capabilities.

My recommendation: capability-scoped for ADRs 1, 4, 5, 6 (tenant-internal mechanics); repo-wide for ADRs 2 and 3 (tenant identity and contract versioning are likely to be reused).

## What I need from you to proceed

1. Confirm or adjust the **TR review** (frontmatter date + the two tweaks above).
2. Confirm or adjust the **ADR grouping and sequence**.
3. Confirm the **location split** (capability-scoped vs repo-wide).
4. Pick which ADR you want to draft first — I'd suggest **ADR-1 (tenant isolation model)** since it unblocks the most.

Once you give the green light, I'll draft them one at a time in MADR 4.0.0 format (matching the style of `0007-standard-api-error-response.md`), with options laid out and a recommended decision rather than a prescribed one — your call on which option wins.
