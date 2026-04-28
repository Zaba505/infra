# Extracting Business Requirements: Shared File Vault

Happy to help — but let's pump the brakes for a second. Business requirements get extracted *from* a defined capability and its user experiences. Since you mentioned "we haven't written the capability doc yet," what we actually have right now is what's in your head. That's fine as a starting point, but I want to be honest about the order of operations so we don't end up with requirements that drift later.

I see two reasonable paths. Tell me which you want:

**Path A — Quick & dirty (recommended for "let's just go"):**
You answer the elicitation questions below, I synthesize a `business-requirements.md` directly from your answers. We skip the formal capability doc for now and circle back to write it later from the requirements (reverse-engineered). Faster, but the audit trail back to "the capability says X" will be weaker.

**Path B — Proper order:**
We do a lightweight capability sketch first (15 min of Q&A), then I extract BRs from that sketch with proper traceability links. Slower but the BRs trace cleanly to the capability + UXes, which matches how the rest of the repo's docs are structured.

---

## Elicitation questions (applies to either path)

To extract business requirements I need to understand the *business intent*, not the technical shape. Please answer in your own words — bullets are fine.

### 1. Purpose & value
- In one sentence, what is the shared file vault *for*?
- What problem does its absence cause today? Who feels that pain?
- If we shipped this and it worked perfectly, what would change for the business / household / users?

### 2. Actors & ownership
- Who uses it? (you only? family? guests? services / automation?)
- Who *owns* the data inside it? (single owner, per-file owner, shared pool?)
- Who can grant access to whom?

### 3. The "shared" part
- Shared between *whom*, exactly? People? Devices? Services? All three?
- Is sharing default-on (everything visible to all members) or default-off (explicit grants per file/folder)?
- Are there sharing scopes (read-only, read-write, time-bound links, public links)?

### 4. The "vault" part
- What does "vault" imply to you here? (encryption-at-rest? versioning? immutability? audit log? recovery from deletion?)
- Are there files that should be *more* protected than others (sensitivity tiers)?
- Retention: do files live forever, or is there a lifecycle (auto-expire, archive, hard-delete)?

### 5. Content & scale
- What kinds of files? (docs, photos, video, secrets/credentials, backups, mixed?)
- Rough volume — KBs, GBs, TBs?
- Largest single file you'd realistically put in?
- Read-heavy or write-heavy?

### 6. Access shape
- From where do people access it? (browser, mobile, CLI, mounted filesystem, API?)
- Online-only, or offline/sync expected?
- LAN-only, or accessible from the internet via your Cloudflare-fronted setup?

### 7. Trust & safety
- What's the worst-case scenario you want to prevent? (data loss? unauthorized access? ransomware? accidental deletion? leak of a shared link?)
- Any compliance/regulatory considerations, or is this purely personal?
- Do you need an audit trail of who-accessed-what-when?

### 8. Boundaries (what it is NOT)
- What are you explicitly *not* trying to build here? (e.g., "not a Dropbox competitor," "not a backup system," "not a media library")
- What's out of scope for v1?

### 9. Success criteria
- How will you know it's working? (concrete, measurable: "I stop using X," "family stops emailing PDFs around," "I can find any tax doc in <30s")
- What would make you call this a failure?

### 10. Non-functional gut-checks
- Availability expectation? (best-effort? "I'd be annoyed if it's down for an hour"? "must survive home internet outage"?)
- Recovery expectation if a disk dies? (acceptable to lose last 24h? zero data loss?)
- Anyone else depending on it being up, or just you?

---

## What I'll produce once you answer

A `business-requirements.md` with entries shaped like:

```
BR-01: The vault SHALL allow <actor> to <do thing> so that <business value>.
  Source: <capability section / UX step / your answer to Q#>
  Acceptance: <observable condition>
```

Each requirement will be:
- Phrased in business language (no tech choices)
- Traceable back to one of your answers (or to the capability doc, if we go Path B)
- Testable (someone could check whether it's met)
- Scoped to *what* and *why*, never *how*

Reply with **Path A** or **Path B**, then take a swing at the questions — answer what's obvious, mark `?` on what you haven't thought through, and we'll iterate on the gaps together.
