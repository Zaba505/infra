# Stage 2 — Enumerating Decisions

`tech-requirements.md` has `reviewed_at: 2026-04-26` and the file hasn't been modified since review. **Stage 1 gate is open — we're now entering Stage 2.**

Per the skill's discipline, I'm not going to start drafting ADRs yet. First we enumerate the **set of decisions the requirements force** and confirm the list. Bundling decisions or skipping straight to drafting produces tangled ADRs that can't be superseded independently.

## Proposed decision list

Below is the set of capability-scoped decisions I see in `tech-requirements.md`. Each is anchored in the TR(s) that force it. I have flagged ones that look cross-capability and should be deferred to a shared ADR rather than drafted here.

### Capability-scoped (candidates for `adrs/{NNNN}-...md`)

1. **Tenant compute isolation boundary** — addresses TR-01, TR-04
   What enforces "no tenant can observe another tenant's compute / traffic"? Per-tenant Cloud Run service vs. shared service with in-process tenant routing vs. per-tenant GKE namespace, etc. Also constrains how zero-downtime updates (TR-04) work, since the unit of update is the unit of isolation.

2. **Tenant state storage & partitioning** — addresses TR-01, TR-03, TR-05, TR-06
   Where does tenant state live, and how is it partitioned such that (a) no tenant can read another's, (b) per-tenant export for eviction is feasible, (c) bulk import for migration is feasible? Options include per-tenant Firestore database, shared Firestore with tenant-prefixed collections, per-tenant Postgres, etc.

3. **Tenant identity propagation across services** — addresses TR-01, TR-03
   How does a request carry "which tenant am I acting for" from the edge through every internal hop, in a way the isolation boundary can trust? Signed JWT claim, mTLS SPIFFE ID, header injected by an authenticating reverse proxy, etc.

4. **Platform contract versioning & concurrent-version hosting** — addresses TR-02
   How are contract versions identified, and how does the platform run two versions concurrently during a migration window? Semver vs. date-based vs. monotonic integer; routing strategy (URL prefix, header, separate deployment per version).

5. **Zero-downtime tenant update mechanism** — addresses TR-04
   Given the isolation boundary chosen in (1), how is an update applied without dropping in-flight requests? Blue/green with traffic-split, rolling with readiness gates, drain-and-replace, etc.

6. **Per-tenant observability scoping** — addresses TR-01, TR-03
   How are metrics/logs/traces tagged at emit time and filtered at query time so a tenant only sees their own? Cloud Logging log buckets per tenant vs. shared bucket with label-scoped views vs. dedicated Grafana org per tenant.

7. **Tenant data export mechanism** — addresses TR-05
   Format (the open question already surfaced: on-demand vs. continuously-available), packaging (single archive, per-collection files), delivery (signed GCS URL, push to tenant-supplied bucket), and the export window's enforcement.

8. **Tenant data import / migration mechanism** — addresses TR-06
   How does the platform accept pre-existing data idempotently with integrity verification? Idempotency key strategy, hashing/checksum scheme, partial-failure resume protocol.

### Likely cross-capability — flag and defer

- **TR-07 (Cloudflare → GCP path)** is itself an inherited shared decision. I would *not* draft a capability ADR for it. Any sub-decision purely about edge topology (e.g. mTLS trust anchor rotation, WireGuard mesh shape) belongs in `docs/content/r&d/adrs/`. If a sub-decision is capability-specific (e.g. how the platform's tenant router plugs into the existing edge), it can stay capability-scoped — but I want to call that out before drafting so we don't accidentally re-decide topology here.
- **Identity propagation (decision 3)** is borderline. If the repo already has a shared identity-propagation ADR, this becomes "how the platform consumes the existing standard" (capability-scoped). If not, the *standard itself* is shared and should be drafted in `docs/content/r&d/adrs/` first. Worth a quick check before I draft.

## What I need from you before drafting any ADR

1. **Confirm or correct the decision list above.** Anything missing? Anything that should be split further or merged?
2. **Resolve the two flag-and-defer items:** Is there an existing shared identity-propagation ADR I should cite as a constraint, or do we need to pause and draft one in `r&d/adrs/` first?
3. **Pick the order.** I'll draft them one at a time, mirror back the decision and rationale before writing the file, and only move on once you've marked it `accepted`. My suggested order — because later decisions depend on earlier ones — is: 1 → 2 → 3 → 5 → 4 → 6 → 7 → 8. Happy to reorder.

Once you confirm, I'll start ADR-0001 by mirroring back the problem statement, the TR citations, and at least two considered options before writing anything to `docs/content/capabilities/self-hosted-application-platform/adrs/0001-*.md`.

## Notes on preconditions

- I read `tech-requirements.md` end-to-end; the review gate is open (`reviewed_at: 2026-04-26`, file unchanged since).
- I did **not** invent requirements or smuggle decisions back into Stage 1. The two existing open questions in `tech-requirements.md` (export cadence, contract versioning scheme) are folded into decisions 7 and 4 respectively.
- I have not created `adrs/_index.md` yet; I'll create it together with ADR-0001 once you confirm the list.
</content>
</invoke>