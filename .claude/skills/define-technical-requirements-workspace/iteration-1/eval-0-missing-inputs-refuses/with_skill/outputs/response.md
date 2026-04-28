I can't extract technical requirements for "shared file vault" yet — and I want to explain why rather than just barrel ahead, because the discipline at this step is what makes everything downstream trustworthy.

## Why I'm stopping

The `define-technical-requirements` skill has a hard precondition: it extracts TRs **from** an already-written capability doc and its user-experience docs. Each TR has to link back to a specific section of those docs so the chain of evidence is auditable later (every ADR cites TR-NN, every tech-design line traces TR → ADR → component).

I checked the repo:

- `docs/content/capabilities/` contains only `self-hosted-application-platform/` and `self-hosted-personal-media-storage.md`.
- There is **no** `docs/content/capabilities/shared-file-vault/` page bundle.
- There are therefore no UX docs under `.../shared-file-vault/user-experiences/` either.

The skill is explicit about this case:

> If the capability or any UX docs are missing, **stop and route the user to `define-capability` / `define-user-experience` first.** Tech requirements derived from missing inputs are tech requirements with missing reasons. Refuse to proceed; do not invent the missing inputs in your head.

"I know what I want — let's just go" is exactly the failure mode the skill is designed to prevent. If I write TRs from your verbal description, every one of them will have a broken `Source:` link (or worse, a fabricated one), `reviewed_at` will gate `plan-adrs` against a doc that has no real foundation, and the first ADR that cites `TR-03` will be citing something I made up on the spot. That is the unjustified-decisions trap the lifecycle exists to close.

## What to do instead — in order

1. **Run `define-capability`** for "shared file vault." That conversation captures stakeholders, business rules, success criteria, out-of-scope items, etc., and produces `docs/content/capabilities/shared-file-vault/_index.md` as a page bundle.
2. **Run `define-user-experience`** once per distinct user journey you have in mind (e.g. "upload a file to the vault", "share a vault file with another user", "revoke access", "recover a deleted file" — whatever your mental model includes). Each produces a UX doc under `docs/content/capabilities/shared-file-vault/user-experiences/`.
3. **Then come back and invoke this skill.** I'll read the capability and every UX end-to-end, skim `docs/content/r&d/adrs/` for prior shared decisions that constrain you (cloud provider, network topology, error format, ID standard), note the repo's house patterns from `CLAUDE.md`, and extract a numbered, source-linked `tech-requirements.md`.

## What I will *not* do right now

- I will not draft a placeholder capability doc from your description. That's `define-capability`'s job, and it's a real conversation, not a fill-in-the-blank.
- I will not draft "provisional" TRs for you to massage later. Provisional TRs without sources are just opinions, and once they're on the page they tend to ossify.
- I will not propose technologies or storage approaches. Even after the capability and UX exist, this skill explicitly forbids picking technologies — that belongs to `plan-adrs` / `define-adr`.

## A small offer

If it's helpful, when you're ready to start `define-capability`, you can dump your "I know what I want" brain into the chat as raw notes — stakeholders, what users do, what's in/out of scope, any non-negotiables — and that skill will turn it into a proper capability doc through structured questions. That's the right place for the energy you've got right now.

Want to kick off `define-capability` for "shared file vault"?
