# Implementation plan: self-hosted-personal-media-storage

## Scope check

Tech design and all five component design docs are present and `reviewed_at: 2026-04-27`. No deferred items. Safe to slice.

Stub state confirms no open `story(gap):` or `story(component):` issues, so we are not blocked by upstream gaps.

## Slicing principles applied

- One issue per task; each task ships independently and is reviewable on its own.
- Foundational Terraform modules first (they define contracts the services depend on).
- Each service decomposes into: skeleton -> backend client + endpoint interface -> endpoints. This matches the `services/machine` reference shape (main.go -> app/ -> endpoint/ + service/).
- Cross-service runtime dependencies (e.g. photo-store calling tenant-registry) are stubbed at the interface boundary so tasks don't serialize.
- Final integration task validates the TRs end to end.

## Task breakdown (13 tasks)

| # | Task | Type | Depends on |
|---|------|------|-----------|
| 1 | `cloud/firestore-tenants` module (schema) | terraform | — |
| 2 | `cloud/media-bucket` module | terraform | — |
| 3 | `services/tenant-registry` skeleton | go | — |
| 4 | tenant-registry Firestore client + endpoint interface | go | 3 |
| 5 | tenant-registry CRUD endpoints | go | 4 |
| 6 | `services/photo-store` skeleton | go | — |
| 7 | photo-store GCS client + tenant validator | go | 6 |
| 8 | photo-store media endpoints | go | 7 |
| 9 | `services/share-service` skeleton | go | — |
| 10 | share-service Firestore clients (grants + append-only audit) | go | 9 |
| 11 | share-service share endpoints + signed-URL | go | 10 |
| 12 | Cloud Run + LB wiring (Terraform) | terraform | 1, 2, 5, 8, 11 |
| 13 | E2E integration test (TR-01, TR-03, TR-06, TR-07) | test | 5, 8, 11 |

TR coverage check: TR-01 (1, 2, 8), TR-02 (2), TR-03 (11), TR-04 (1, 3, 5), TR-05 (1, 5), TR-06 (11, 13), TR-07 (1, 10, 11). All covered.

Issues logged to `outputs/intended_issues.jsonl` (13 entries) per the no-`gh` policy.

## Note on the per-task plan request

The `plan-implementation` skill's contract is "task breakdown plus issues — not the per-task plan; that belongs to Step 12 (`/begin-task`)." The user explicitly asked to fold the schema task's per-task plan in here ("save me from running Step 12 later"), so it follows below. Flagging the deviation so a future review of this thread isn't surprised.

---

## Schema task — per-task plan (Task #1: `cloud/firestore-tenants`)

**Issue:** "story(task): scaffold cloud/firestore-tenants module (schema)"
**Component spec:** `docs/content/capabilities/self-hosted-personal-media-storage/components/firestore-tenants.md`
**Established by:** ADR-0001, ADR-0003
**TRs satisfied:** TR-01, TR-04, TR-05, TR-07

### Files to create

All paths are repo-relative.

| Path | Purpose |
|------|---------|
| `cloud/firestore-tenants/main.tf` | `google_firestore_database` (Native mode) + `google_firestore_index` resources for the three collections, with `lifecycle { create_before_destroy = true }`. |
| `cloud/firestore-tenants/variables.tf` | `project_id`, `location` (default `nam5`), `tenant_registry_service_account`, `share_service_service_account`. |
| `cloud/firestore-tenants/outputs.tf` | `database_name`, `collections` map (`tenants`, `share_grants`, `share_audit`). |
| `cloud/firestore-tenants/iam.tf` | `google_project_iam_member` / `google_firestore_*` IAM: tenant-registry SA -> rw on `tenants`; share-service SA -> rw on `share-grants`, append-only on `share-audit` (datastore.entities.create only — no update/delete). |
| `cloud/firestore-tenants/versions.tf` | `terraform { required_providers { google = "7.11.0" } }` to match the rest of `cloud/`. |
| `cloud/firestore-tenants/README.md` | Inputs/outputs table; example `module` block; note about `create_before_destroy`. |

Indexes to declare (driven by component spec query paths):
- `tenants`: composite on `(state ASC, updated_at DESC)` — supports list-by-state.
- `share_grants`: composite on `(grantee_id ASC, expires_at DESC)` — supports grantee-side lookups; second composite on `(tenant_id ASC, state ASC, expires_at DESC)` — supports list-grants-for-tenant.
- `share_audit`: composite on `(grant_id ASC, at DESC)` — supports per-grant audit replay.

### Tests to write

Terraform modules in this repo aren't unit-tested (CI runs `terraform fmt -recursive -check` only). For this task add:

1. `cloud/firestore-tenants/examples/basic/main.tf` — minimal example invocation that `terraform validate` can run against. CI `fmt` will keep it clean.
2. Local validation steps (run before pushing):
   - `terraform -chdir=cloud/firestore-tenants/examples/basic init -backend=false`
   - `terraform -chdir=cloud/firestore-tenants/examples/basic validate`
   - `terraform fmt -recursive -check cloud/firestore-tenants`
3. (Optional, defer if scope creeps) `tflint` config aligned with the rest of `cloud/`.

No Go tests for this task.

### Migration command

Firestore Native databases aren't created with `terraform migrate`; this is a fresh provision. The "exact migration command" is the apply against the platform project:

```bash
# from repo root, against the GCP project hosting this capability
terraform -chdir=cloud/firestore-tenants/examples/basic init
terraform -chdir=cloud/firestore-tenants/examples/basic plan \
  -var "project_id=$GCP_PROJECT_ID" \
  -var "location=nam5" \
  -var "tenant_registry_service_account=tenant-registry-sa@$GCP_PROJECT_ID.iam.gserviceaccount.com" \
  -var "share_service_service_account=share-service-sa@$GCP_PROJECT_ID.iam.gserviceaccount.com" \
  -out tfplan
terraform -chdir=cloud/firestore-tenants/examples/basic apply tfplan
```

When the module is consumed from the real root module (Task #12), the equivalent `terraform apply` runs from that root. The repo has no automated apply, so this is a manual operator step (`CLAUDE.md`: "No automated apply; manual deployment required").

### Definition of done

- All six files present; `terraform fmt -recursive -check` is clean.
- `terraform validate` passes for `examples/basic`.
- Outputs documented in README match `outputs.tf`.
- IAM bindings restrict share-service SA to `datastore.entities.create` on `share-audit` (no update/delete) — verify by inspecting the rendered plan.
- PR linked to the task issue with commit `feat(issue-N): scaffold firestore-tenants module` on a `story/issue-N/firestore-tenants` branch.

### Risks / watch-outs

- Firestore database name is project-scoped and effectively immutable; pick `(default)` unless ADR-0001 specified otherwise — re-check ADR-0001 before apply.
- `create_before_destroy` on `google_firestore_database` is a no-op for the database itself (single per project); the lifecycle directive applies meaningfully to the indexes and IAM.
- Append-only audit IAM must be enforced via custom role or careful predefined-role choice; `roles/datastore.user` is too broad.

---

## Files

- Plan + per-task plan: `/home/carson/github.com/Zaba505/infra/.claude/skills/plan-implementation-workspace/iteration-1/eval-3/without_skill/outputs/response.md`
- Intended issues (13): `/home/carson/github.com/Zaba505/infra/.claude/skills/plan-implementation-workspace/iteration-1/eval-3/without_skill/outputs/intended_issues.jsonl`
