# Cannot proceed — `business-requirements.md` is unreviewed

I can't extract technical requirements yet. The skill's preconditions require `business-requirements.md` to have been human-reviewed before any TRs are derived from it, and right now its frontmatter shows:

```yaml
reviewed_at: null
```

at `docs/content/capabilities/self-hosted-application-platform/business-requirements.md`.

## Why this matters

TRs are the technical translation of BRs — every TR has to cite the BR-NN it derives from. If I extract TRs from a BR list that hasn't been signed off, every downstream artifact (TRs, ADRs, components in the eventual tech-design) inherits the un-reviewed reasoning. The whole BR → TR → ADR → component evidence chain stops being trustworthy at its root. The review gate exists precisely to prevent that.

## What I need you to do

1. Open `docs/content/capabilities/self-hosted-application-platform/business-requirements.md` and read it end-to-end. There are 54 BRs in there — go through each one, edit anything that is wrong, missing, or imprecise, and remove anything that doesn't belong.
2. When you're satisfied, set the frontmatter to today's ISO date:
   ```yaml
   reviewed_at: 2026-04-28
   ```
   (or tell me you've reviewed it verbally and I'll record the date on your behalf — but I do need an explicit confirmation; I won't set it silently).
3. Re-invoke this skill. I'll then re-read the BR doc, the capability `_index.md`, every UX under `user-experiences/`, and the shared ADRs under `docs/content/r&d/adrs/`, and produce `tech-requirements.md` with each TR citing its BR (and any framing UX / capability section / inherited ADR).

## What I will not do

- I will not extract TRs against the current BR doc and assume "you'll review later". The gate is upstream of TR extraction by design.
- I will not flip `reviewed_at` myself based on file age or guesswork. Only an explicit human "yes, I reviewed it" (or a date you set) clears the gate.
- I will not invent BRs in my head to fill gaps I notice while reading the capability/UX docs. If something is missing, it should be added back in `extract-business-requirements`, not silently in TR extraction.

## When you come back

Once `reviewed_at` is set to a date newer than the file's last modification time, I will:
- Read the capability doc, every UX, and the shared ADRs for inherited constraints.
- Walk the BR list and translate each BR into one or more TRs (or zero, if already covered by an inherited constraint).
- Cite BR-NN as the primary source on every TR, plus capability/UX/ADR context where it shapes the translation.
- Use Hugo `ref` shortcodes for every internal link, and call out any BR/UX/capability headings that are missing explicit `{#anchor}` annotations needed for deep-linking.
- Park any premature decisions you mention as **Open Questions** for the later `plan-adrs` step rather than letting them in as TRs.
- Save to `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` with `reviewed_at: null` so you can review the freshly-extracted TRs.

Once that file exists, **you'll need to review it and set its own `reviewed_at`** before `plan-adrs` will run — same gate, one stage downstream.

---

**Note on Hugo render validation:** per the eval setup, I am skipping the `cd docs && hugo server` rendered-page validation step that the skill normally requires after writing `tech-requirements.md`. That step would still apply on a real run.
