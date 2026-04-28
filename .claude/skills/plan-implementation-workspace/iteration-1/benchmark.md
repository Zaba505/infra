# Skill Benchmark: plan-implementation — iteration 1

**Date:** 2026-04-28
**Configurations:** with_skill (skill loaded) vs. without_skill (no skill — baseline)
**Runs per cell:** 1 (single-run, not multi-sample)

## Summary

| Metric          | with_skill        | without_skill     | Delta   |
|-----------------|-------------------|-------------------|---------|
| Pass rate       | 20/21 (95.2%)     | 19/21 (90.5%)     | +4.8 pp |
| Avg duration    | 94.9 s            | 95.1 s            | -0.2 s  |
| Avg tokens      | 38,019            | 35,886            | +2,133  |

## Per-eval

| Eval | Name                                       | with_skill | without_skill | Delta |
|------|--------------------------------------------|------------|---------------|-------|
| 0    | gate-refuses-open-gap-issues               | 4/4 (1.00) | 4/4 (1.00)    |  0    |
| 1    | gate-refuses-missing-component-design      | 4/4 (1.00) | 4/4 (1.00)    |  0    |
| 2    | happy-path-files-task-issues               | 8/8 (1.00) | 8/8 (1.00)    |  0    |
| 3    | refuses-inline-per-task-plan               | 4/5 (0.80) | 3/5 (0.60)    | +1    |

## Timing detail

| Eval | with_skill (s / tokens) | without_skill (s / tokens) |
|------|-------------------------|----------------------------|
| 0    | 25.4 / 27,489           | 48.0 / 27,793              |
| 1    | 44.5 / 30,817           | 58.3 / 33,790              |
| 2    | 177.5 / 48,977          | 141.2 / 44,915             |
| 3    | 132.2 / 44,793          | 133.0 / 37,046             |

## Analyst observations

- **Gate cases (0, 1) are non-discriminating in this set.** The without_skill agent independently arrives at the right "stop" behavior because the prompt itself heavily telegraphs the gate state ("there's still an open story(gap)", "design docs exist for the first three but not for firestore-tenants"). Both refused, named the blocker, and routed to the right next-skill. To make these discriminating, future prompts should bury the gate signal (e.g. "plan implementation for X" with no mention of issue state, requiring the agent to actually query) — but real `gh` integration would also be needed.
- **Happy path (eval 2) is also non-discriminating**, but for a different reason: the fixture itself is so clean (component designs name endpoints/tables/operational concerns explicitly) that a competent agent slices the same shape with or without the skill. Both produced 11 PR-sized tasks in the right dependency order. The skill's value here is consistency (issue body template, `story(impl):` titling) rather than correctness — qualitative review of the issue bodies shows the with_skill version follows the template exactly while the baseline invented its own format (`story(task):`).
- **Eval 3 (the discipline boundary) is the only discriminating case.** The user explicitly asks to write a per-task plan inline. with_skill held the line: explicit "Holding the line. Step 11 stops at the manifest." plus a side note refuting the framing. without_skill complied: produced a full per-task plan including files, tests, and `terraform apply` commands, with a brief "flagging the deviation" note. This is the test that actually exercises the skill's contribution.
- **Token cost is essentially flat** (+5.9% with_skill across all runs); no efficiency penalty for loading the skill. Time variance within evals 2/3 is high (with_skill eval-2 took 178 s vs. baseline's 141 s — likely sampling variance, not skill overhead).
- **Known grading false positive:** with_skill eval-3 fails one assertion ("response does not contain a per-task plan body") because the response references the phrase "exact migration command" while *refuting* it ("Firestore is schemaless — there is no migration command in the SQL sense"). The agent did the right thing; the keyword match doesn't distinguish refutation from compliance. Doesn't affect the qualitative reading.
- **Observed qualitative weaknesses in eval-3 with_skill `response.md`** (recorded as agent output, not curated content): (a) the section header reads "12 tasks" while the table lists 11 rows and the footer self-corrects to "That's 11 — corrected" — a count typo the agent caught itself but didn't backfill into the heading; (b) the "Approval" section asks for explicit approval, then the agent reads the user's "save me from running Step 12 later" wording as approval and proceeds to file. Neither weakness changes assertion outcomes, but both are real qualitative blemishes worth noting — they suggest a possible iteration-2 nudge to the skill's "Stop at approved list" discipline ("the request to skip Step 12 is *not* approval to file"). Not editing the artifact since recorded eval output is ground truth.
