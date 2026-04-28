# Can't extract tech requirements yet — preconditions aren't met

I hear you, but the `define-technical-requirements` skill is **Step 6** of the capability lifecycle, and it has hard preconditions I can't skip. Here's the chain it expects to already exist before I can do anything useful:

1. **Capability doc** at `docs/content/capabilities/shared-file-vault/_index.md` — defines stakeholders, business rules, success criteria, out-of-scope. **Doesn't exist yet.**
2. **One or more UX docs** under `docs/content/capabilities/shared-file-vault/user-experiences/` — define how users actually accomplish things with the vault. **Don't exist yet.**
3. **Reviewed `business-requirements.md`** at `docs/content/capabilities/shared-file-vault/business-requirements.md` — the authoritative input every TR must cite (`TR-NN` derives from `BR-NN`). Its `reviewed_at` frontmatter must be a date newer than the file's last modification. **Doesn't exist yet.**

The skill is explicit on this:

> If `business-requirements.md` is missing or unreviewed, **stop and route the user to `extract-business-requirements`** … If the capability or any UX docs are missing, route to `define-capability` / `define-user-experience` first. Tech requirements derived from missing inputs are tech requirements with missing reasons. Refuse to proceed; do not invent the missing inputs in your head.

So if I just "went" off your verbal description, every TR I produced would be unsourceable — there'd be no `BR-NN` to cite and no UX section to frame the technical translation. The whole point of this doc is the audit chain `BR → TR → ADR → component`. Skipping the BRs breaks the chain at the root, and downstream `plan-adrs` will refuse to run anyway because the gate (`reviewed_at` on a doc that doesn't exist) can never be satisfied.

## The actual path forward

Run these three skills in order, then come back to me:

1. **`/define-capability`** — we'll talk through what the shared file vault is at the business level: who uses it, what business rules govern it, what success looks like, what's explicitly out of scope. Output: `docs/content/capabilities/shared-file-vault/_index.md`.
2. **`/define-user-experience`** — for each end-to-end user journey the vault enables (e.g. "share a file with another tenant", "revoke access", "audit who downloaded what"), we define one UX doc. You'll likely want 2–4 of these. Output: `docs/content/capabilities/shared-file-vault/user-experiences/{ux-name}.md`.
3. **`/extract-business-requirements`** — walks the capability + UXs and produces `business-requirements.md`. **You then review and edit it**, set `reviewed_at` to today's ISO date, and we're cleared to come back here.

That sounds like a lot, but the conversations are cheap — most of what's in your head will fall out in the first one. The upside: by the time we get to TRs, every constraint we write will have a written business reason behind it, and every later ADR will be picking among options that all satisfy a sourced TR. That's the whole reason the lifecycle is staged.

## What I'm not doing

- I'm **not** writing a `tech-requirements.md` from your verbal description. (Would create unsourced TRs and poison the audit chain.)
- I'm **not** modifying anything under `docs/` (per your instructions for this run).
- I'm **not** guessing at the capability shape. The skill is explicit: "do not invent the missing inputs in your head."

When you're ready, kick off `/define-capability` and we'll start there. If you'd rather I help draft a one-paragraph "north star" you can paste into that conversation to get it rolling, happy to do that — just say the word.
