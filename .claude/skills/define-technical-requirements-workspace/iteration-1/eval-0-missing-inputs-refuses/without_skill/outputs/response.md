# Can't extract tech requirements yet — the capability doc doesn't exist

I want to help, but I have to push back on "let's just go" here. Extracting technical requirements only works if there's a written capability (and ideally one or more user experiences) to extract them *from*. That source-of-truth is what makes each requirement traceable and reviewable — otherwise I'd just be inventing requirements from my own assumptions about what a "shared file vault" means to you, and we'd be debugging that gap for weeks.

## What I checked

In this repo, capabilities live under `docs/content/capabilities/`. Today I see:

- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/_index.md`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-application-platform/`
- `/home/carson/github.com/Zaba505/infra/docs/content/capabilities/self-hosted-personal-media-storage.md`

There is no `shared-file-vault` capability doc. There are also no UX docs for it. So there is nothing to extract from.

The repo's own workflow (see the `define-capability`, `define-user-experience`, and `define-technical-requirements` skills) is explicit about the order:

1. `define-capability` — write the business-level capability doc
2. `define-user-experience` — define one or more end-to-end user journeys against that capability
3. `define-technical-requirements` — extract a `tech-requirements.md` whose entries each link back to a line in the capability doc or a UX

Step 3 is what you're asking for, but steps 1 and 2 are the inputs. Without them, "tech requirements" becomes a wishlist instead of a traceable artifact.

## What I'd like to do instead

I'd suggest we run this in two passes — happy to drive both back-to-back in the same session:

1. Run `define-capability` for "shared file vault." You said you know what you want; great — that skill is built to elicit it via conversation and produce the capability doc under `docs/content/capabilities/shared-file-vault/`.
2. Run `define-user-experience` for at least one journey (e.g. "a household member uploads a document and shares it with another member"). One UX is enough to start; we can add more later.
3. *Then* run `define-technical-requirements`, and every requirement we surface will point back to a specific line in (1) or (2).

If you want to skip ahead anyway — for example, jot informal notes I can later reshape into a real capability doc — say so explicitly and I'll capture them as a scratch doc clearly marked as not-yet-a-capability. I just don't want to write `tech-requirements.md` against a phantom source.

Want to start with `define-capability` now? If yes, just say "go" and tell me roughly who the vault is for (you only? household? guests?) and I'll take it from there.
