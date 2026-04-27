I can't start a tech design for "shared file vault" yet — it isn't a defined capability in this repo.

I checked `docs/content/capabilities/` and only two capabilities exist:

- `self-hosted-application-platform`
- `self-hosted-personal-media-storage`

There is no `shared-file-vault/` page bundle, no `_index.md` for it, and therefore no UX docs under `user-experiences/`. The `define-tech-design` skill has a hard precondition that the parent capability doc and every UX doc must already exist before Stage 1 (extracting technical requirements) can begin. The reasoning: tech requirements are *extracted* from capability rules and UX steps with source links back to them. With no capability or UX docs to cite, every "requirement" I produced would be invented, every ADR would be unjustified, and the whole audit chain (TR → ADR → tech-design) would be fiction.

"I know what I want — let's just go" is exactly the failure mode the skill is built to refuse. If we skip ahead, you'll get a confidently-wrong design.

Here's the path forward — please pick one:

1. **Define the capability first (recommended).** Run the `define-capability` skill for `shared-file-vault`. That gives us stakeholders, business rules, success criteria, and the out-of-scope list. Then run `define-user-experience` once per user journey (upload, share, retrieve, revoke, etc.). Then come back here and I'll extract `tech-requirements.md` from real sources.

2. **It's actually one of the existing capabilities.** If "shared file vault" is really a slice of `self-hosted-personal-media-storage` (or a new UX on top of it), tell me and we'll either add a UX there or split out a new capability — but that's still a `define-capability` / `define-user-experience` call, not this skill.

3. **You want a throwaway sketch, not a tracked design.** If you just want me to brainstorm architecture options informally — no `tech-requirements.md`, no ADRs, nothing committed under `docs/content/capabilities/` — say so explicitly and I'll do that as a scratch conversation outside this skill. But nothing from that sketch will be authoritative or linkable later.

Which do you want?
