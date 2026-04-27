# Response

We're still in **Stage 1** — I can't start drafting ADRs yet, and I want to be straight with you about why before we lose time.

I read the existing `tech-requirements.md` for the self-hosted-application-platform capability. It's in good shape — seven TRs (TR-01 through TR-07) covering tenant isolation, contract-version rollout, per-tenant observability, zero-downtime updates, evicted-tenant data export, lossless data import, and the inherited Cloudflare → GCP topology constraint, plus two open questions on the export model and contract versioning scheme.

**The blocker:** the frontmatter still has `reviewed_at: null`. The Stage 2 gate is explicit — I won't propose ADRs until a human has reviewed the requirements and `reviewed_at` is set to a date newer than the file's last modification. This gate exists for a reason: ADRs cite TR-NN identifiers as their chain of evidence, and if we draft them against unreviewed requirements we'll either (a) rationalize decisions for requirements that turn out to be wrong, or (b) silently break ADR provenance later when you fix the requirements.

**What I need from you to unblock Stage 2:**

1. Read `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` end-to-end. Edit anything that's wrong, missing, or muddled. Pay particular attention to:
   - Whether TR-01 through TR-07 actually cover everything the capability `_index.md` and the six UX docs demand. Stand-up-the-platform and host-a-capability are both UXes that I don't see directly cited — are their requirements subsumed in the existing TRs, or are we missing some?
   - The two open questions at the bottom — those need resolving (or explicit deferral to specific Stage 2 ADRs) before we can enumerate decisions.
2. Set `reviewed_at:` in the frontmatter to today's ISO date (2026-04-26), or tell me to set it on your behalf as a record of your verbal sign-off.

**Once the gate is open, here's the plan for Stage 2.** I won't write any ADRs until you confirm this decision list, but to give you a preview of what we'd be enumerating from the current TRs:

- Tenant compute & isolation boundary (forced by TR-01, TR-04, TR-07) — what runs tenant workloads and how isolation is enforced.
- Tenant state storage & isolation (forced by TR-01, TR-05, TR-06) — where tenant data lives and how export/import work.
- Platform contract versioning & coexistence (forced by TR-02) — resolves the open question on versioning scheme.
- Zero-downtime update mechanism (forced by TR-04) — rolling vs. blue/green vs. tenant-coordinated.
- Per-tenant observability pipeline (forced by TR-03) — how telemetry is segregated and queried.
- Tenant data export mechanism (forced by TR-05) — resolves the open question on on-demand vs. continuous.
- Tenant data import / migration mechanism (forced by TR-06) — idempotency and integrity verification.

A couple of those (notably anything that further constrains the Cloudflare → GCP path) may turn out to be cross-capability — if so, I'll flag and stop rather than draft them as capability-scoped ADRs.

**Note on the worktree:** your setup step asked me to copy the fixture into the worktree at `.claude/worktrees/agent-a35e8b35bc35c7c11/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`, but that worktree's `docs/content/capabilities/` only contains the flat `self-hosted-application-platform.md` — there's no page-bundle directory to drop the file into, and the `cp` was sandbox-denied so I couldn't run it. The reviewed copy of `tech-requirements.md` that I worked from lives in the main repo at `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`, which is page-bundle form and matches the fixture content. If you specifically need the worktree populated, let me know and I'll either migrate the flat file to bundle form there or re-run the copy with permissions.

Set `reviewed_at` and ping me, and I'll come back with the proposed decision list for your sign-off before writing any ADRs.
