# Plan Implementation — self-hosted-personal-media-storage

## Gate result: STOPPING — preconditions not met

I cannot plan implementation tasks for this capability yet. Two related blockers:

### 1. Missing component design doc

The inventory in `tech-design.md` lists **five** components (note: the user prompt said four — the actual inventory also includes `media-bucket`):

1. `tenant-registry` — design doc present at `components/tenant-registry.md`
2. `photo-store` — design doc present at `components/photo-store.md`
3. `share-service` — design doc present at `components/share-service.md`
4. `firestore-tenants` — **design doc MISSING** (`components/firestore-tenants.md` does not exist)
5. `media-bucket` — design doc present at `components/media-bucket.md`

The `tech-design.md` audit trail links `firestore-tenants` to ADR-0001 and ADR-0003 and gives it ownership of the `tenants`, `share-grants`, and `share-audit` Firestore collections — substantive responsibilities that need a real component design before implementation can be sliced.

### 2. Open `story(component):` issue confirms it

`gh` shows `#824 story(component): firestore-tenants — self-hosted-personal-media-storage` is still **open**. An open `story(component):` issue against this capability is a hard signal the component design hasn't been authored. Per the skill's preconditions, this alone halts slicing.

## What needs to happen before re-invoking

1. Run `define-component-design` against issue **#824** to author `docs/content/capabilities/self-hosted-personal-media-storage/components/firestore-tenants.md`. The doc should cover at minimum:
   - Firestore database mode and location decision (per ADR-0001)
   - Collection schemas: `tenants`, `share-grants`, `share-audit`
   - Composite indexes required by the registry and share-service query patterns
   - IAM bindings for the `tenant-registry` and `share-service` service accounts
   - Backup / point-in-time-recovery posture (per ADR-0001 if it speaks to durability)
2. Close issue **#824** once the design lands.
3. Re-invoke `plan-implementation` for `self-hosted-personal-media-storage`.

At that point the gate passes (tech-design `reviewed_at: 2026-04-27` is current; no open `story(gap):` issues; all five component designs present) and I can slice tasks across the five components.

## Why I'm not slicing around it

`firestore-tenants` is the data-layer foundation that `tenant-registry` and `share-service` both depend on. Slicing implementation tasks for those services without a written `firestore-tenants` design would force the implementor to invent the schema and indexes mid-PR — exactly the ambiguity Step 11 is meant to surface, not paper over.

## No issues filed

Per skill discipline, I file no `story(impl):` issues until the manifest is approved, and I cannot produce a manifest with a component design missing. `outputs/intended_issues.jsonl` is intentionally empty for this run.
