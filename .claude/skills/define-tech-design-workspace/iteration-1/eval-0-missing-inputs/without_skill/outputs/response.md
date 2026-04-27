Before we dive into tech design, I want to flag something and make sure we're aligned.

Per the repo conventions (and just good practice for this codebase), capability docs under `docs/content/capabilities/` are the source of truth that tech designs hang off of. Jumping straight to tech without one tends to produce a design that's solving the wrong problem — or solving the right problem at the wrong layer — because the user-facing experience, tenancy model, and platform contracts haven't been pinned down yet.

A few concrete things I'd want nailed down before I can give you a tech design worth the bits it's stored in:

1. **Who is the tenant / user?** Is this a per-tenant vault on the self-hosted platform (fits the `issue-697` story arc you're on), a shared-across-tenants vault, or a single-user personal thing? The blast radius and auth model change completely.
2. **"Shared" with whom?** Shared between users in one tenant? Shared across tenants? Shared via public links? Each one is a different system.
3. **File semantics** — immutable blobs, mutable files, versioned? Max size (KB, MB, GB, TB)? Streaming reads or whole-file? Do we need rename/move, or is it content-addressed?
4. **Access control** — per-file ACLs, per-folder, capability tokens, signed URLs? Audit requirements?
5. **Durability / retention** — backups, soft-delete window, legal hold, geographic constraints?
6. **Encryption** — at-rest only, client-side E2E, per-tenant KMS keys?
7. **Existing platform fit** — does this live behind the same Cloudflare → home-lab → GCP path as everything else? Is the storage backend GCS, a home-lab disk, or both with tiering?
8. **Scale envelope** — number of files, number of concurrent readers, hot vs cold access patterns. (Even rough numbers; this is a home lab, but it changes whether we reach for object storage vs a real filesystem vs a DB-backed blob store.)

My recommendation: let's spend 15 minutes writing a thin capability doc first — even just bullet points covering the user experience, the tenancy model, and the platform contract — following the same shape as the recent `issue-697` capability docs you've been landing. Then the tech design basically writes itself and we won't have to redo it when one of the above answers surprises us.

If you'd rather just start sketching, I can do that too — but tell me your answers to (1), (2), (3), and (4) at minimum and I'll draft a tech design ADR (MADR 4.0.0, matching `docs/content/r&d/adrs/`) on that basis, with explicit "ASSUMPTION" callouts wherever I'm guessing. Your call.

Which way do you want to go?
