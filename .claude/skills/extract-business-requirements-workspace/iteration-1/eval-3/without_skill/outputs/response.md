I didn't add BR-06 as requested — the proposed text reads as a technical decision, not a business requirement, and adding it would break the pattern the existing BRs follow.

## Why I pushed back

"Use S3 multipart upload for resumable tenant data uploads" specifies:

- **A technology** (S3 multipart upload) — a specific AWS API.
- **A mechanism** (multipart upload chunks).

Compare to the existing BRs in the file:

- BR-01 — "no tenant can ever observe another tenant's state" (tenant-facing invariant)
- BR-02 — "existing tenants must keep working when the platform contract evolves" (tenant-facing outcome)
- BR-03 — "a tenant must be able to query observability data for their own workloads" (tenant-facing capability)
- BR-04 — "operator-initiated tenant updates must be invisible to the tenant's end users" (tenant-facing outcome)
- BR-05 — "an evicted tenant must be able to leave with all of their data" (tenant-facing commitment)

Each is a *what the tenant experiences / what the platform owes the tenant*, deliberately silent on *how*. Each "Why this is a requirement, not a TR or decision" note explicitly defers implementation choices (versioning scheme, rolling deploys, export mechanism) to TR/ADR stages.

"Use S3 multipart upload" is the kind of statement those notes warn against — it picks the implementation up front.

## Two ways forward

**Option A — re-frame as a business requirement.** The underlying tenant need is probably something like:

> **BR-06: Tenants must be able to upload large data sets without losing progress on transient failure**
> **Source:** (which UX or capability section is driving this? — I need a link)
> **Requirement:** A tenant uploading data to the platform must be able to resume an interrupted upload without re-transmitting bytes already accepted, up to a defined object-size ceiling.
> **Why this is a requirement, not a TR or decision:** Resumability is a tenant-perceived outcome; the transport mechanism (S3 multipart, tus, chunked HTTP, etc.) is a TR/ADR question.

If you confirm the source UX/section and the framing, I'll add this version.

**Option B — record it where it belongs.** "Use S3 multipart upload" is an architectural decision. It belongs in an ADR (with options weighed: S3 multipart vs. tus vs. chunked PUT vs. pre-signed resumable URLs) under the capability's `decisions/` directory, traced back to a TR like "TR-NN: support resumable uploads of objects up to N GB", which in turn traces back to the BR in Option A.

## What I'd like from you

1. Which UX or capability section motivates resumable uploads? (BRs require a `Source:` link — there's no such UX listed in the current file.)
2. Confirm Option A wording (or supply your own tenant-outcome framing) and I'll add BR-06.
3. Separately, do you want me to file an ADR-planning issue for the multipart-vs-alternatives decision?

I have **not** modified `business-requirements.md`.
