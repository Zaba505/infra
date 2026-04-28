#!/usr/bin/env python3
"""Grade plan-implementation eval runs. Emits grading.json per run."""
import json, os, re, sys
from pathlib import Path

WS = Path(__file__).parent

def read(p):
    try:
        return Path(p).read_text()
    except FileNotFoundError:
        return ""

def grade_run(eval_id, run_dir):
    out = run_dir / "outputs"
    resp = read(out / "response.md").lower()
    issues_path = out / "intended_issues.jsonl"
    issues = []
    if issues_path.exists():
        for line in issues_path.read_text().splitlines():
            line = line.strip()
            if line:
                try:
                    issues.append(json.loads(line))
                except json.JSONDecodeError:
                    pass
    n_issues = len(issues)

    expectations = []

    if eval_id == 0:  # gate-refuses-open-gap-issues
        refuses = ("stopping" in resp or "cannot plan" in resp or "blocked" in resp or "cannot proceed" in resp or "won't" in resp or "do not file" in resp) and "step 11" in resp[:5000] or "do not file" in resp or "blocked" in resp or "stopping" in resp or "won't enumerate" in resp or "won't" in resp
        # simpler: "refuses" if no issues filed AND response declines
        refuses = (n_issues == 0) and any(kw in resp for kw in ["stopping", "cannot plan", "blocked", "do not file", "won't"])
        expectations.append({"text": "Refuses to slice / file issues", "passed": refuses, "evidence": f"intended_issues={n_issues}; response declines={refuses}"})
        expectations.append({"text": "Names blocking gap issue (#812 / tenant-id derivation)", "passed": ("#812" in resp or "tenant-id derivation" in resp or "tenant id derivation" in resp), "evidence": "matches '#812' or 'tenant-id derivation'"})
        expectations.append({"text": "Routes user to resolve the gap", "passed": any(kw in resp for kw in ["resolve #812", "close that", "close #812", "define-component-design", "amending adr", "define-adr", "define-technical-requirements", "resolve the gap", "resolve it", "resolve the open gap"]), "evidence": "mentions one of the resolution routes"})
        expectations.append({"text": "No issues filed", "passed": n_issues == 0, "evidence": f"intended_issues={n_issues}"})

    elif eval_id == 1:  # gate-refuses-missing-component-design
        refuses = (n_issues == 0) and any(kw in resp for kw in ["stopping", "cannot plan", "blocked", "do not file", "won't", "cannot slice"])
        expectations.append({"text": "Refuses to slice / file issues", "passed": refuses, "evidence": f"intended_issues={n_issues}"})
        expectations.append({"text": "Names firestore-tenants as missing component design", "passed": "firestore-tenants" in resp, "evidence": "matches 'firestore-tenants'"})
        expectations.append({"text": "Routes user to define-component-design", "passed": "define-component-design" in resp, "evidence": "matches 'define-component-design'"})
        expectations.append({"text": "No issues filed", "passed": n_issues == 0, "evidence": f"intended_issues={n_issues}"})

    elif eval_id == 2:  # happy-path-files-task-issues
        # tasks anchored per-component, multiple per service
        # PR-sized: > 5 tasks, < 25 tasks
        sized_ok = 5 < n_issues < 25
        expectations.append({"text": "Lists multiple tasks per component (not one-per-component)", "passed": n_issues >= 7, "evidence": f"intended_issues={n_issues}"})
        expectations.append({"text": "Tasks are PR-sized (5 < count < 25)", "passed": sized_ok, "evidence": f"count={n_issues}"})
        # dependency ordering — look for prereq language
        has_deps = any(kw in resp for kw in ["prerequisite", "depends on", "prereq"])
        expectations.append({"text": "Identifies prerequisite/dependency relationships", "passed": has_deps, "evidence": "mentions prerequisite or depends on"})
        # mirrors back: looks for 'before i file' / 'mirror' / 'approve' / 'go'
        mirrors = any(kw in resp for kw in ["before i file", "before filing", "mirror", "do the slices look right", "before any issues", "say go"])
        expectations.append({"text": "Mirrors back the task list before filing", "passed": mirrors, "evidence": "presents list and asks for approval"})
        # files one issue per task with story(impl) title
        impl_titles = [i for i in issues if "title" in i and i["title"].lower().startswith("story(impl)")]
        expectations.append({"text": "Filed issues use story(impl): title convention", "passed": len(impl_titles) == n_issues and n_issues > 0, "evidence": f"story(impl) count={len(impl_titles)} of {n_issues}"})
        # bodies link design sources
        body_text_all = " ".join((i.get("body", "") for i in issues)).lower()
        expectations.append({"text": "Issue bodies link parent capability + design sources", "passed": "components/" in body_text_all and ("adr-" in body_text_all or "adrs/" in body_text_all), "evidence": "bodies reference components/ and adrs/"})
        expectations.append({"text": "Issue bodies reference Step 12 / per-task development flow", "passed": ("step 12" in body_text_all or "per-task" in body_text_all), "evidence": "bodies reference step 12 or per-task"})
        # no per-task plan inline (no 'files to touch' / 'migration command' / 'tests to write' as a section in response)
        no_inline_plan = not any(kw in resp for kw in ["files to touch", "files to create", "exact migration command", "migration command", "definition of done"])
        expectations.append({"text": "No per-task plan inline (no file lists, migration commands, etc.)", "passed": no_inline_plan, "evidence": "response does not include a per-task development plan"})

    elif eval_id == 3:  # refuses-inline-per-task-plan
        # holds the line: response declines to write per-task plan
        declines_plan = any(kw in resp for kw in ["holding the line", "step 11 stops", "won't write", "not going to", "step 12's job", "per-task plan is step 12", "per-task plans are step 12", "decline", "stop and file the issue", "i'm not going to write"])
        expectations.append({"text": "Refuses to inline a per-task development plan", "passed": declines_plan, "evidence": "response holds the line on Step 11/12 boundary"})
        # explains why
        explains = any(kw in resp for kw in ["step 12", "per-task flow", "per-pr plan", "tailored", "per-component review"])
        expectations.append({"text": "Explains why per-task plans belong to Step 12", "passed": explains, "evidence": "response references Step 12 ownership"})
        # still files issues including the schema task
        body_text_all = " ".join((i.get("body", "") for i in issues)).lower()
        title_text_all = " ".join((i.get("title", "") for i in issues)).lower()
        files_schema = ("schema" in title_text_all) and n_issues >= 5
        expectations.append({"text": "Still files implementation issues (incl. the schema task)", "passed": files_schema, "evidence": f"intended_issues={n_issues}; schema in titles={'schema' in title_text_all}"})
        # acknowledges and offers to carry intent
        carries = any(kw in resp for kw in ["re-invoke", "re-invoke against", "carry", "feed", "follow-up", "schema task", "step 12 against"])
        expectations.append({"text": "Acknowledges user intent and offers Step-12 follow-up path", "passed": carries, "evidence": "response routes the user's request to Step 12"})
        # check that the response does NOT contain a full per-task plan body for the schema task
        no_plan_body = not any(kw in resp for kw in ["files to create", "files to touch", "exact migration command", "definition of done"])
        expectations.append({"text": "Response does not contain a per-task plan body for the schema task", "passed": no_plan_body, "evidence": "no 'files to create' / 'definition of done' section"})

    passed = sum(1 for e in expectations if e["passed"])
    total = len(expectations)
    return {
        "expectations": expectations,
        "passed": passed,
        "total": total,
        "pass_rate": passed / total if total else 0.0,
    }

def main():
    for eval_id in range(4):
        for cfg in ("with_skill", "without_skill"):
            run_dir = WS / f"eval-{eval_id}" / cfg
            if not (run_dir / "outputs" / "response.md").exists():
                continue
            grading = grade_run(eval_id, run_dir)
            (run_dir / "grading.json").write_text(json.dumps(grading, indent=2))
            print(f"eval-{eval_id} {cfg}: {grading['passed']}/{grading['total']}")

if __name__ == "__main__":
    main()
