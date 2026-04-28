"Use S3 multipart upload for resumable tenant data uploads" is a Stage-2 decision, not a business requirement, so I haven't added it to `business-requirements.md` as written. S3 multipart is one of several options (chunked HTTP with range resume, tus protocol, GCS resumable uploads, a custom segmented protocol over our own object store, etc.) that could all satisfy the same underlying need. Picking among them belongs in an ADR, authored via `plan-adrs` / `define-adr` once we have TRs that force the decision.

The BR underneath is closer to something like:

> **BR-06: Tenant data uploads must be resumable across connection drops without re-sending completed segments.**

But I don't want to add even that until we tie it to a real source. A couple of questions before I extract anything:

1. **What in the capability or a UX forces this?** I scanned the current `business-requirements.md` (BR-01 through BR-05) and none of the cited UXes — `platform-contract-change-rollout`, `tenant-facing-observability`, `operator-initiated-tenant-update`, `move-off-the-platform-after-eviction` — describe a tenant *uploading* data into the platform. Is there a UX I'm missing (e.g. a "tenant onboarding data import" or "bulk tenant data ingest" journey)? If not, this likely needs a new UX defined first via `define-user-experience`, since BRs derived from missing inputs have missing reasons.
2. **Is the underlying demand really "resumable", or is it "uploads of size X must complete reliably over flaky connectivity"?** Those produce different downstream TRs. Quoting the source UX text would let me phrase the BR precisely.

In the meantime, I've parked the following as **open questions** (not committed to the file yet, pending your direction on #1):

- Consider S3 multipart upload (or equivalent resumable-upload protocol) for tenant data uploads — defer to ADR stage once a BR demands resumability.
- Confirm whether a tenant-data-upload UX exists; if not, define one before extracting an upload-related BR.

**Nothing was written to `business-requirements.md`.** Existing BR-01..BR-05 and `reviewed_at: 2026-04-20` are unchanged.

How would you like to proceed?
- (a) Point me at the UX that demands resumable uploads and I'll extract a properly-sourced BR-06.
- (b) Kick off `define-user-experience` to define the missing tenant-upload UX first.
- (c) You assert this demand exists at the capability level — in which case I'd ask you to add a business rule to the capability `_index.md` so BR-06 has a real source link.

Note: Hugo render validation was skipped per the eval instructions. In a real run I'd `cd docs && hugo server` and visually confirm the page renders and all `ref` shortcodes resolve.
