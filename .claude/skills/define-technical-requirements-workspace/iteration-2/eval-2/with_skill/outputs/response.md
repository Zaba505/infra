I've extracted the technical requirements for the **Self-Hosted Application Platform** capability to `tech-requirements.md` in this evaluation's outputs directory (NOT under `docs/`, per eval instructions).

## Summary

- **39 TRs** added (TR-01 through TR-39), each citing the BR-NN(s) it derives from. Multi-sourced TRs (the most load-bearing ones) cite all the BRs that force them — for example, TR-09 (administrative isolation) cites BR-04, BR-19, and BR-34; TR-22 (tenant-facing observability) cites BR-17 and BR-18.
- **0 preserved / 0 flagged** — this is a fresh extraction, no prior `tech-requirements.md` existed.
- `reviewed_at: null` in frontmatter — this is a fresh doc, awaiting your review.

## How the BRs mapped to TRs

I walked the entire BR list (BR-01 through BR-54). Most BRs produced one TR; some BRs are bundled into a single TR where they force the same constraint from different angles (BR-01/13/23 → TR-01 uniform tenant runtime; BR-02/35/36/42 → TR-02 cross-environment reproducibility; BR-04/19/34 → TR-09 admin isolation). A handful of BRs did NOT produce a TR because they are governance/relationship commitments that don't translate into a technical constraint:

- **BR-11** (default response to fall-behind is "bring current, don't evict") — operator-relationship policy, no technical surface forced beyond what TR-32/33 already encode.
- **BR-12** (default response to unmet need is "evolve the platform") — posture commitment; the technical surface (new-offering branch joining update flow) is captured by TR-34.
- **BR-43** (skill development must not influence buy-vs-build) — decision-discipline rule, not a system constraint.
- **BR-49** (no per-tenant deadline extensions) — coordination policy; the technical surface is just "deadline is a single global field on the rollout record," already implicit in TR-16/TR-17.

I called these out implicitly by their absence; if you'd prefer explicit "no TR forced" rows, I can add them.

## Open Questions parked for the ADR stage

9 items, including:
- Drift-detection mechanism (TR-06)
- Backup/DR standard numbers (TR-27) — flagged as a BR-doc gap to fill before TR-27 is complete
- Deeper backup-tier policy after eviction retention (TR-29 / TR-31) — carried forward from the BR doc's own open question
- Tenant-facing observability signal bundle composition (TR-22)
- Migration concurrency capacity bound (TR-20)
- Cost-attribution granularity (TR-39)
- Explicit `{#anchor-id}` annotations missing on capability and UX section headings — only `#br-NN` deep-links resolve cleanly today; section-level anchors on `_index.md` and UX pages still need to be added in source

## What I deliberately did NOT do

- **No technology choices.** TRs say "one-shot job runner with provision-run-observe-teardown lifecycle," not "use Kubernetes Jobs"; "secret-management surface scoped per tenant," not "use GCP Secret Manager"; "concurrent old/new contract forms during rollout window," not "use feature flags." Those are `plan-adrs` / `define-adr` work.
- **No restated BRs.** Each TR translates the BR into the constraint the technical solution is forced to meet (mechanisms, surfaces, thresholds in implementation-neutral terms).
- **Did not modify anything under `docs/`** per eval instructions.
- **Skipped Hugo render validation** per eval instructions. In a real run, the next step would be `cd docs && hugo server` to confirm every `ref` resolves and the page renders. **Note that several `ref` links in this doc point at section anchors (e.g. `business-requirements.md#br-29`) that depend on the source BR doc carrying explicit `{#br-NN}` heading anchors — see the Open Questions section.**

## Next step for the human

> Review and edit `tech-requirements.md` before we proceed. When you're done, set `reviewed_at:` in the frontmatter to today's ISO date. The next skill — `plan-adrs` — won't enumerate decisions until that's set.
