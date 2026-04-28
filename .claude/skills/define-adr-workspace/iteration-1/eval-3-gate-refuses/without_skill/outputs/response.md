# ADR drafting refused — review gate not satisfied

**Decision requested:** tenant state storage for `self-hosted-application-platform`
**Tech-requirements file:** `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md`

## Why I stopped

The `define-adr` skill enforces a hard precondition: `tech-requirements.md` must carry a `reviewed_at:` ISO date that is **newer than the file's last modification time**. The current frontmatter says:

```yaml
reviewed_at: null
```

That means the technical requirements have not yet been human-reviewed. Drafting an ADR against unreviewed TRs would produce a decision sourced from drafts — every option's pros/cons would be anchored to TR-NN citations that may still change, and the ADR's audit value (showing *why* the decision was made against a stable constraint set) would be void.

This gate is not a formality. ADR options are required to be evaluated explicitly against TR-NNs (e.g. "Option A satisfies TR-01 and TR-04 but partially fails TR-05"). If the TRs are still in flux, those evaluations are unstable, and a future engineer trying to supersede this ADR will not be able to trust the rejected-options analysis.

## What needs to happen before I can draft this ADR

1. A human reviews `tech-requirements.md` end-to-end — confirms each TR's wording, source link, and that the list is complete for the capability.
2. The reviewer sets `reviewed_at:` in the frontmatter to today's ISO date (e.g. `reviewed_at: 2026-04-27`) **as part of the same change** that finalizes the review, so the timestamp is newer than any subsequent edit.
3. Re-invoke `define-adr` for the tenant-state-storage decision. At that point I will:
   - Read every TR the decision touches (likely TR-01, TR-03, TR-04, TR-05, TR-06, TR-07 based on a skim).
   - Surface research tasks before proposing options (existing Firestore footprint, regional failover semantics, prior shared ADRs on identifiers and Cloudflare→GCP topology).
   - Draft ≥2 options with pros/cons phrased in TR-NN terms.
   - Mirror options back and let you pick — `status: accepted` is set only after explicit human selection.

## What I did not do

- I did not draft options.
- I did not pick a storage technology.
- I did not create an ADR file under `docs/content/capabilities/self-hosted-application-platform/adrs/`.
- I did not invent TR review status or backdate the gate.

## Suggested next step

Review `tech-requirements.md`, set `reviewed_at: 2026-04-27` in its frontmatter, commit, then re-run: *"Draft the ADR for tenant state storage for the self-hosted-application-platform capability."*
