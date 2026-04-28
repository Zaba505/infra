# Refusing to draft ADR — review gate not passed

I can't draft the ADR for tenant state storage yet.

`docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` has `reviewed_at: null`. Per the `define-adr` preconditions, the review gate requires `reviewed_at:` to be an ISO date newer than the file's last modification time. The current TRs haven't been reviewed by a human, and an ADR sourced from unreviewed drafts is meaningless — the option analysis would be anchored to TR-NNs that may still change.

## What needs to happen first

1. A human reviews `tech-requirements.md` end-to-end, confirming each TR (TR-01 through TR-07) is correctly extracted from the capability and UX docs and is the right shape for forcing architectural decisions.
2. Update the frontmatter to set `reviewed_at:` to today's ISO date (e.g. `reviewed_at: 2026-04-27`) once the review is complete.
3. Re-invoke `define-adr` for the tenant state storage decision. With the gate green, I'll proceed with Step 1 (research surfacing) — for tenant state storage that will likely include questions like:
   - Is there an existing Firestore/Cloud SQL footprint we'd reuse, or are we provisioning fresh?
   - Does a prior shared ADR in `docs/content/r&d/adrs/` already constrain identifiers, regional placement, or backup posture?
   - What's the realistic cost/quota envelope for the candidate stores under TR-01 (isolation) and TR-05 (export)?

I'm not going to draft options against unreviewed TRs — that's exactly the failure mode the gate exists to prevent.
