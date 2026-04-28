#!/usr/bin/env python3
"""Grade plan-implementation eval runs. Emits grading.json per run.

Each assertion produces an `evidence` string built from what was actually
observed in the response — matched substrings, counts, booleans — so the
evidence reflects the actual pass/fail decision rather than a static template.
"""
import json
from pathlib import Path

WS = Path(__file__).parent


def read(p: Path) -> str:
    try:
        return p.read_text()
    except FileNotFoundError:
        return ""


def first_match(haystack: str, needles: list[str]) -> str | None:
    for n in needles:
        if n in haystack:
            return n
    return None


def all_matches(haystack: str, needles: list[str]) -> list[str]:
    return [n for n in needles if n in haystack]


def grade_run(eval_id: int, run_dir: Path) -> dict:
    out = run_dir / "outputs"
    resp = read(out / "response.md").lower()
    issues_path = out / "intended_issues.jsonl"
    issues: list[dict] = []
    if issues_path.exists():
        for line in issues_path.read_text().splitlines():
            line = line.strip()
            if line:
                try:
                    issues.append(json.loads(line))
                except json.JSONDecodeError:
                    pass
    n_issues = len(issues)
    expectations: list[dict] = []

    if eval_id == 0:
        # gate-refuses-open-gap-issues
        decline_markers = ["stopping", "cannot plan", "blocked", "do not file", "won't"]
        decline_hits = all_matches(resp, decline_markers)
        refuses = (n_issues == 0) and bool(decline_hits)
        expectations.append({
            "text": "Refuses to slice / file issues",
            "passed": refuses,
            "evidence": (
                f"intended_issues={n_issues}; decline markers matched={decline_hits or 'none'}"
            ),
        })

        gap_markers = ["#812", "tenant-id derivation", "tenant id derivation"]
        gap_hits = all_matches(resp, gap_markers)
        expectations.append({
            "text": "Names blocking gap issue (#812 / tenant-id derivation)",
            "passed": bool(gap_hits),
            "evidence": f"matched={gap_hits or 'none'}",
        })

        route_markers = [
            "resolve #812", "close that", "close #812", "define-component-design",
            "amending adr", "define-adr", "define-technical-requirements",
            "resolve the gap", "resolve it", "resolve the open gap",
        ]
        route_hits = all_matches(resp, route_markers)
        expectations.append({
            "text": "Routes user to resolve the gap",
            "passed": bool(route_hits),
            "evidence": f"resolution-route markers matched={route_hits or 'none'}",
        })

        expectations.append({
            "text": "No issues filed",
            "passed": n_issues == 0,
            "evidence": f"intended_issues={n_issues}",
        })

    elif eval_id == 1:
        # gate-refuses-missing-component-design
        decline_markers = ["stopping", "cannot plan", "blocked", "do not file", "won't", "cannot slice"]
        decline_hits = all_matches(resp, decline_markers)
        refuses = (n_issues == 0) and bool(decline_hits)
        expectations.append({
            "text": "Refuses to slice / file issues",
            "passed": refuses,
            "evidence": f"intended_issues={n_issues}; decline markers matched={decline_hits or 'none'}",
        })

        names_component = "firestore-tenants" in resp
        expectations.append({
            "text": "Names firestore-tenants as missing component design",
            "passed": names_component,
            "evidence": f"'firestore-tenants' present in response={names_component}",
        })

        names_skill = "define-component-design" in resp
        expectations.append({
            "text": "Routes user to define-component-design",
            "passed": names_skill,
            "evidence": f"'define-component-design' present in response={names_skill}",
        })

        expectations.append({
            "text": "No issues filed",
            "passed": n_issues == 0,
            "evidence": f"intended_issues={n_issues}",
        })

    elif eval_id == 2:
        # happy-path-files-task-issues
        expectations.append({
            "text": "Lists multiple tasks per component (not one-per-component)",
            "passed": n_issues >= 7,
            "evidence": f"intended_issues={n_issues} (threshold>=7)",
        })

        sized_ok = 5 < n_issues < 25
        expectations.append({
            "text": "Tasks are PR-sized (5 < count < 25)",
            "passed": sized_ok,
            "evidence": f"intended_issues={n_issues}",
        })

        dep_markers = ["prerequisite", "depends on", "prereq"]
        dep_hits = all_matches(resp, dep_markers)
        expectations.append({
            "text": "Identifies prerequisite/dependency relationships",
            "passed": bool(dep_hits),
            "evidence": f"dependency markers matched={dep_hits or 'none'}",
        })

        mirror_markers = [
            "before i file", "before filing", "mirror", "do the slices look right",
            "before any issues", "say go",
        ]
        mirror_hits = all_matches(resp, mirror_markers)
        expectations.append({
            "text": "Mirrors back the task list before filing",
            "passed": bool(mirror_hits),
            "evidence": f"mirror markers matched={mirror_hits or 'none'}",
        })

        impl_titles = [i for i in issues if i.get("title", "").lower().startswith("story(impl)")]
        impl_ok = (n_issues > 0) and (len(impl_titles) == n_issues)
        expectations.append({
            "text": "Filed issues use story(impl): title convention",
            "passed": impl_ok,
            "evidence": f"story(impl) titles={len(impl_titles)} of {n_issues}",
        })

        body_text_all = " ".join((i.get("body", "") for i in issues)).lower()
        has_components_link = "components/" in body_text_all
        has_adrs_link = "adr-" in body_text_all or "adrs/" in body_text_all
        expectations.append({
            "text": "Issue bodies link parent capability + design sources",
            "passed": has_components_link and has_adrs_link,
            "evidence": f"bodies reference components/={has_components_link}, adrs/={has_adrs_link}",
        })

        has_step12 = "step 12" in body_text_all or "per-task" in body_text_all
        expectations.append({
            "text": "Issue bodies reference Step 12 / per-task development flow",
            "passed": has_step12,
            "evidence": f"bodies reference 'step 12' or 'per-task'={has_step12}",
        })

        plan_markers = [
            "files to touch", "files to create", "exact migration command",
            "migration command", "definition of done",
        ]
        plan_hits = all_matches(resp, plan_markers)
        expectations.append({
            "text": "No per-task plan inline (no file lists, migration commands, etc.)",
            "passed": not plan_hits,
            "evidence": f"per-task-plan markers matched in response={plan_hits or 'none'}",
        })

    elif eval_id == 3:
        # refuses-inline-per-task-plan
        decline_markers = [
            "holding the line", "step 11 stops", "won't write", "not going to",
            "step 12's job", "per-task plan is step 12", "per-task plans are step 12",
            "i'm not going to write",
        ]
        decline_hits = all_matches(resp, decline_markers)
        expectations.append({
            "text": "Refuses to inline a per-task development plan",
            "passed": bool(decline_hits),
            "evidence": f"refusal markers matched={decline_hits or 'none'}",
        })

        explain_markers = ["step 12", "per-task flow", "per-pr plan", "tailored", "per-component review"]
        explain_hits = all_matches(resp, explain_markers)
        expectations.append({
            "text": "Explains why per-task plans belong to Step 12",
            "passed": bool(explain_hits),
            "evidence": f"explanation markers matched={explain_hits or 'none'}",
        })

        title_text_all = " ".join((i.get("title", "") for i in issues)).lower()
        files_schema = ("schema" in title_text_all) and (n_issues >= 5)
        expectations.append({
            "text": "Still files implementation issues (incl. the schema task)",
            "passed": files_schema,
            "evidence": f"intended_issues={n_issues}; 'schema' in any title={'schema' in title_text_all}",
        })

        carry_markers = [
            "re-invoke", "re-invoke against", "carry", "feed", "follow-up",
            "schema task", "step 12 against",
        ]
        carry_hits = all_matches(resp, carry_markers)
        expectations.append({
            "text": "Acknowledges user intent and offers Step-12 follow-up path",
            "passed": bool(carry_hits),
            "evidence": f"follow-up markers matched={carry_hits or 'none'}",
        })

        plan_body_markers = ["files to create", "files to touch", "exact migration command", "definition of done"]
        plan_body_hits = all_matches(resp, plan_body_markers)
        no_plan_body = not plan_body_hits
        if not no_plan_body and "no migration command" in resp:
            note = " (note: 'exact migration command' appears in a refutation, not a per-task plan body — known grader false positive)"
        else:
            note = ""
        expectations.append({
            "text": "Response does not contain a per-task plan body for the schema task",
            "passed": no_plan_body,
            "evidence": f"per-task-plan-body markers matched={plan_body_hits or 'none'}{note}",
        })

    passed = sum(1 for e in expectations if e["passed"])
    total = len(expectations)
    return {
        "expectations": expectations,
        "passed": passed,
        "total": total,
        "pass_rate": round(passed / total, 4) if total else 0.0,
    }


def main() -> None:
    for eval_id in range(4):
        for cfg in ("with_skill", "without_skill"):
            run_dir = WS / f"eval-{eval_id}" / cfg
            if not (run_dir / "outputs" / "response.md").exists():
                continue
            grading = grade_run(eval_id, run_dir)
            (run_dir / "grading.json").write_text(json.dumps(grading, indent=2) + "\n")
            print(f"eval-{eval_id} {cfg}: {grading['passed']}/{grading['total']}")


if __name__ == "__main__":
    main()
