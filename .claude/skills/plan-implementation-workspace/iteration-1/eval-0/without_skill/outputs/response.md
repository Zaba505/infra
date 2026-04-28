# Plan implementation — self-hosted-personal-media-storage

## Verdict: blocked. Do not file implementation issues yet.

The tech design is composed and `reviewed_at: 2026-04-27` is current, but there is an open `story(gap):` issue against this capability:

- **#812 story(gap): tenant-id derivation — self-hosted-personal-media-storage** (open)

Implementation planning requires a complete design. An open gap means at least one piece of the contract between components is undecided, so any task breakdown filed now would either:

1. encode an assumption about how tenant IDs are derived that the gap may overturn, or
2. omit the work the gap will eventually require.

Either way the issues filed would be wrong. We do not file `story(task):` issues against an incomplete design.

There is also an internal inconsistency worth flagging: `tech-design.md` ends with **"Deferred / Open: None."** while issue #812 is open against the same capability. Whoever resolves #812 should also reconcile the tech-design's Deferred/Open section so the document and the issue tracker agree.

## What "tenant-id derivation" touches

Skimming the design (`docs/content/capabilities/self-hosted-personal-media-storage/`), tenant ID is load-bearing across every component:

- **tenant-registry** — keys the `tenants` Firestore collection (ADR-0001).
- **photo-store** — prefixes GCS object paths by tenant ID for per-tenant isolation (ADR-0002).
- **share-service** — joins share grants to a tenant when minting signed URLs (ADR-0003).
- **firestore-tenants** module — the collection schema depends on the ID shape (string? UUID? caller-supplied?).
- **media-bucket** module — IAM conditions and prefix conventions depend on the ID shape.

Because all five components consume the tenant ID, deriving it is a cross-cutting decision, not a single-component concern. Filing per-task issues before it lands would scatter the rework.

## Recommended next steps (in order)

1. **Resolve #812.** Decide where tenant IDs come from (platform onboarding claim? registry-minted UUID? caller-supplied slug with uniqueness check?) and how they propagate (header? mTLS SAN? token claim?). If the decision is architectural, it should land as a new ADR (likely ADR-0004) via `plan-adrs` / `define-adr`; if it's purely an interface choice, it can be folded into the affected component design docs.
2. **Reconcile `tech-design.md`.** Update the **Deferred / Open** section to reference #812 until it closes, then flip back to "None" once the gap is resolved and the relevant component designs / ADRs are updated. Bump `reviewed_at`.
3. **Re-run implementation planning** once #812 is closed and the design is consistent. At that point the breakdown will likely be: one Terraform task per `cloud/` module (`firestore-tenants`, `media-bucket`), one bootstrap task per Go service (`tenant-registry`, `photo-store`, `share-service`) following the `services/machine/` shape, and per-endpoint tasks layered on top — but pinning that down is premature today.

## Issues filed

None. `outputs/intended_issues.jsonl` was not created — filing implementation issues now would be incorrect per the rationale above.

## Key paths

- `/home/carson/github.com/Zaba505/infra/.claude/skills/plan-implementation-workspace/iteration-1/eval-0/without_skill/staged-repo/docs/content/capabilities/self-hosted-personal-media-storage/tech-design.md` — the composed design; Deferred/Open section needs reconciling with #812.
- `/home/carson/github.com/Zaba505/infra/.claude/skills/plan-implementation-workspace/iteration-1/eval-0/without_skill/staged-repo/docs/content/capabilities/self-hosted-personal-media-storage/adrs/` — three accepted ADRs; tenant-id derivation is not among them.
- `/home/carson/github.com/Zaba505/infra/.claude/skills/plan-implementation-workspace/iteration-1/eval-0/without_skill/staged-repo/docs/content/capabilities/self-hosted-personal-media-storage/components/` — five component design docs that all depend on the tenant-id contract.
