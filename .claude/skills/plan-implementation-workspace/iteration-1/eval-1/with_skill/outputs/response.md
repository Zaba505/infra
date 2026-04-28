# Plan implementation — self-hosted-personal-media-storage

**Stopping — preconditions not satisfied. No tasks sliced, no issues filed.**

## Gate state

- Capability page bundle: found at `docs/content/capabilities/self-hosted-personal-media-storage/_index.md`.
- `tech-design.md`: present.
- Open `story(gap):` issues for this capability: none.
- Open `story(component):` issues for this capability: **#824 `story(component): firestore-tenants — self-hosted-personal-media-storage`** (open).
- Component inventory in `tech-design.md` lists **5** components (not 4 as the prompt suggested):
  1. `tenant-registry` service
  2. `photo-store` service
  3. `share-service`
  4. `firestore-tenants` module
  5. `media-bucket` module
- Component design docs on disk under `components/`:
  - `tenant-registry.md` — present
  - `photo-store.md` — present
  - `share-service.md` — present
  - `media-bucket.md` — present
  - `firestore-tenants.md` — **missing**

## Why I'm stopping

Per Step 11's preconditions, every component named in the inventory must have a design doc authored via `define-component-design` before tasks can be sliced. `firestore-tenants` is in the inventory but has no design doc on disk, and the open `story(component): firestore-tenants` issue (#824) confirms the design work hasn't landed yet.

Slicing tasks now would either (a) fabricate a Terraform-module task without grounding it in the actual schema, indexes, IAM, and provisioning decisions the component design will pin down, or (b) silently drop a component that the tech design's audit trail depends on. Both hide risk. The line in the skill is sharp on purpose: tasks must trace to a real component design.

Note also that the prompt said "tech-design.md lists four components" — the file actually lists five. The two `module`-typed components (`firestore-tenants`, `media-bucket`) are easy to overlook when scanning, but they're load-bearing: the data layer for tenants and the storage layer for media respectively. `media-bucket.md` exists; `firestore-tenants.md` does not.

## What to do next

1. Run `define-component-design` against issue **#824** to author `docs/content/capabilities/self-hosted-personal-media-storage/components/firestore-tenants.md`. The doc should pin down the Firestore collection shape for tenant records, indexes, IAM bindings to the `tenant-registry` service account, and the Terraform-module surface under `cloud/firestore-tenants/` (per repo convention in `CLAUDE.md`).
2. Once that design doc is on disk and #824 is closed, re-invoke `plan-implementation` for `self-hosted-personal-media-storage`. At that point all five components will have designs, and I can read them end-to-end, slice tasks along their seams (schemas, endpoints, modules, deployment wiring), order by hard dependency, mirror the proposed list back for your approval, and then file one `story(impl):` issue per task.

## Issues filed

None. `outputs/intended_issues.jsonl` is intentionally empty — Step 11 does not file issues until the precondition gate passes and the human approves the sliced task list.
