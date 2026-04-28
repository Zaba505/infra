Added BR-06 to `business-requirements.md` for the self-hosted-application-platform capability.

A few things to flag for your review:

1. **BR-06 mixes outcome with mechanism.** As written, "checkpoint upload progress every 10 seconds to durable storage" is a technical mechanism, not a pure business outcome. The business need is "uploads survive a connection drop and can be resumed." The cadence (10 seconds) and the storage choice (durable storage) are implementation details that normally belong in a technical requirement or an ADR. I recorded BR-06 verbatim per your request, but I added a note in the "Why this is a requirement" section and an entry under Open Questions so this gets revisited at the TR stage.

2. **No source link.** Every other BR points back to a section of the capability doc or a specific UX. BR-06 has no such anchor — I marked the source as "User request (no linked capability or UX provided)" and added an Open Question to confirm which UX it derives from. If there is an existing UX (e.g. a tenant-upload journey) it should be linked; if not, that UX likely needs to be defined before BR-06 is review-ready.

3. **Scope check.** "Uploads" isn't a concept that appears in BR-01 through BR-05 or in the surrounding requirements. Worth confirming this belongs to the self-hosted-application-platform capability rather than a separate capability (e.g. a tenant-data-ingestion or artifact-upload capability).

The file was saved to the requested output path; nothing under `docs/` was modified.
