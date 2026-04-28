I can't extract business requirements yet — the skill's preconditions aren't met.

This skill (`extract-business-requirements`) is **Step 5** of the capability development lifecycle. It only works once the upstream artifacts exist, because every BR has to link back — via Hugo `ref` shortcodes — to a specific source in the capability doc or a UX doc. Without those sources, any "requirements" I produce would be invented rather than extracted, and the downstream chain (BR → TR → ADR → component) would have nothing to stand on.

Specifically, the skill requires:

1. A capability page bundle at `docs/content/capabilities/shared-file-vault/_index.md` — read end-to-end for stakeholders, business rules, success criteria, out-of-scope.
2. At least one UX doc under `docs/content/capabilities/shared-file-vault/user-experiences/` — because BRs must serve all the user journeys, not just the one in your head.

Neither exists yet for "shared file vault." I have to refuse to proceed; the skill is explicit that "business requirements derived from missing inputs are business requirements with missing reasons."

## What to do instead

Run the upstream skills first, in order:

1. **`define-capability`** — I'll guide you through scoping the shared-file-vault capability: what the business does, stakeholders, business rules, success criteria, what's explicitly out of scope. This produces `docs/content/capabilities/shared-file-vault/_index.md`.
2. **`define-user-experience`** — for each meaningful user journey in the vault (e.g. "share a file with another tenant", "revoke access", "audit who has accessed what", "recover a deleted file" — whatever you actually want). One UX doc per journey under `user-experiences/`.
3. **Then come back here** and I'll extract `business-requirements.md` against the real, written sources, with every BR carrying a `ref`-linked citation.

I know it feels like extra ceremony when you "already know what you want" — but the value of this pipeline is exactly that the BRs are auditable derivations from written intent, not transcribed from a conversation that no one can review later. The capability and UX docs are also where you'll surface stakeholders, edge cases, and out-of-scope decisions you probably haven't fully thought through yet; skipping them tends to produce BRs that quietly bake in unstated assumptions.

Want to start with `define-capability` for shared-file-vault now? If so, tell me the one-paragraph pitch — what the vault is for, who uses it, and what problem it solves — and I'll drive the rest of the elicitation from there.
