That proposed BR-06 is a TR, not a BR — so I haven't added it.

"The system must checkpoint upload progress every 10 seconds to durable storage" is the technical translation: it names the mechanism (checkpoint to durable storage) and the cadence (every 10 seconds). Both belong in `define-technical-requirements`, derived from a BR — not in `business-requirements.md`.

The BR underneath is closer to something like:

> **BR-06: Uploads must survive a connection drop without losing user-perceived progress**
> A user whose connection drops mid-upload must be able to resume the upload without re-sending data they have already successfully transferred, and without the system telling them their prior progress was lost.

But before I add that, I need a source. Two questions:

1. **Where in the capability or a UX is "uploads" actually demanded?** I re-read `_index.md` and the five UX docs implied by BR-01–BR-05 (`platform-contract-change-rollout`, `tenant-facing-observability`, `operator-initiated-tenant-update`, `move-off-the-platform-after-eviction`) and don't see an upload flow that forces a resumability requirement. The only export-shaped UX is the eviction one (BR-05), and that's a one-shot data export, not an interactive upload. Is there a UX I'm missing — e.g. a tenant-onboarding or workload-deploy UX where the user uploads artifacts? If so, please point me to it (or run `define-user-experience` to add it) so the BR has provenance.
2. **Whose upload?** "Uploads" could mean tenant operators uploading workload images/artifacts to the platform, or end-users of a tenant's app uploading data through it (which would be the tenant's concern, not the platform's). The capability scope matters — the platform may not own the second case at all.

Once we agree on the source, I'll add the BR. In the meantime I'm capturing the user-supplied number as an open question for the TR stage:

> - Consider a 10-second checkpoint cadence to durable storage as the TR derived from the (yet-to-be-sourced) upload-resumability BR — surfaced during BR-06 discussion; defer to TR stage.

**No file changes made.** The existing `business-requirements.md` (BR-01 through BR-05, `reviewed_at: 2026-04-20`) is unchanged, so its review state is preserved. Once you confirm the source for the upload requirement, I'll re-extract — which will null out `reviewed_at` and you'll need to re-review before `define-technical-requirements` will run.

(Skipped Hugo render validation since no file changes were made.)
