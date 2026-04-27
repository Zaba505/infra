---
title: "[0003] Tenant Packaging Form"
description: >
    Tenants ship an OCI image plus a small platform-defined manifest schema describing name, ports, env, secret references, declared resource needs, and (for migration jobs) a re-run-contract flag. The platform translates this into Kubernetes manifests at provision time.
type: docs
weight: 3
category: "strategic"
status: "accepted"
date: 2026-04-27
deciders: []
consulted: []
informed: []
---

**Parent capability:** [Self-Hosted Application Platform](../_index.md)
**Addresses requirements:** [TR-10](../tech-requirements.md#tr-10-provide-a-packaging-form-the-platform-accepts-for-all-tenant-components), [TR-11](../tech-requirements.md#tr-11-provide-a-secret-management-offering-tenants-can-register-secrets-with-and-reference-by-name), [TR-12](../tech-requirements.md#tr-12-provide-a-one-shot-migration-process-offering-that-runs-tenant-supplied-migration-jobs), [TR-22](../tech-requirements.md#tr-22-tracked-changes-and-immutability-for-all-platform-state-modifying-actions), [TR-24](../tech-requirements.md#tr-24-tenant-provisioning-must-run-only-through-the-platforms-existing-definitions), [TR-26](../tech-requirements.md#tr-26-tenants-declare-resource-needs-at-onboarding-and-on-every-modify-the-platform-admits-or-refuses-based-on-those-declarations), [TR-31](../tech-requirements.md#tr-31-migration-jobs-declare-a-re-run-contract-that-the-platform-records-and-respects), [TR-32](../tech-requirements.md#tr-32-per-tenant-authentication-and-isolation-strong-enough-that-no-tenant-or-its-capability-owner-via-the-observability-offering-can-read-another-tenants-data-or-signals)

## Context and Problem Statement

[ADR-0002](./0002-compute-substrate.md) chose Kubernetes as the compute substrate. That constrains but does not fully decide the tenant packaging form: K8s can accept anything from "an OCI image and we'll write the manifests for you" to "a full Helm chart" to "raw Deployment YAML." TR-10 demands one accepted form for tenant components, and that same form must be acceptable for migration-process artifacts (TR-12).

Several other TRs ride on this choice:

- TR-26 says tenants declare resource needs at onboarding and modify; the platform admits/refuses against the declaration. Where do those declarations live — in the artifact, in a side document, or in the issue body?
- TR-31 says migration jobs declare a re-run contract. Same question.
- TR-11 says tenants reference secrets by name. The packaging form must let them do that without exposing values.
- TR-22 / TR-24 say platform state is tracked, immutable, and provisioned only through definitions. The form's expressiveness directly determines how much of the K8s surface a tenant can reach into versus what the platform mediates.
- TR-32 says isolation must hold. A form that lets tenants ship NetworkPolicy / RBAC / Pod Security overrides hands them the keys to the isolation primitives.

The host-a-capability UX further treats artifact handoff as the moment the contract is accepted ([§4](../user-experiences/host-a-capability.md#4-hand-off-packaged-artifacts), [§Constraints — Tenants must accept the platform's contract](../user-experiences/host-a-capability.md#constraints-inherited-from-the-capability)): the design *is* the acceptance, and the artifact is what the design promised. The form must make that submission complete on its own, not require side-channel coordination on the issue thread.

## Decision Drivers

- **TR-10 — single accepted form.** Whatever is chosen must serve every tenant component and every migration job. Multiple forms break the rule.
- **TR-26 / TR-31 declarations live with the artifact.** Submission is acceptance; the declarations have to be inside what is submitted.
- **TR-32 isolation.** The form must not let tenants reach into the substrate's isolation primitives.
- **Capability tiebreaker — minimize friction at the contract boundary.** Capability owners are not asked to be Kubernetes experts (per the host-a-capability persona); the form should be readable and writable without K8s fluency.
- **TR-22 / TR-24 platform-controlled provisioning.** The platform translates the form into the substrate's manifests; the platform owns those manifests, not the tenant.
- **TR-33 (≤2 hr/week).** A form whose schema turns into a 50-field nightmare costs ongoing review effort. Smaller is better; growth is allowed only through the contract-change UX.

## Considered Options

### Option A — OCI image only

Capability owner hands off `registry/path:tag`. All manifest (Deployment, Service, NetworkPolicy, ResourceQuota, env, ports) comes from values the operator fills in on the platform side, derived from issue-thread discussion.

- **TR-10:** met — one artifact form (the image).
- **TR-26 / TR-31:** declarations live outside the artifact, in the issue thread or side documents. Submission is *not* acceptance — the operator must transcribe declarations into the platform definitions before provisioning.
- **TR-32:** strong — tenant cannot influence isolation primitives at all.
- **Friction:** non-trivial topology (multiple ports, special envs, sidecars) becomes back-and-forth on the issue thread. The host-a-capability "submission is acceptance" promise is weakened: an image alone is not enough to provision against.

### Option B — OCI image + Helm chart

Capability owner ships an image and a Helm chart parametrizable with platform-supplied values.

- **TR-10:** met — one artifact form (image + chart bundle).
- **TR-26 / TR-31:** declarations live in chart values; possible but loose — chart authors decide the schema, so the platform has to validate every chart's values against the contract separately.
- **TR-32:** weaker — charts can ship arbitrary K8s resources, including NetworkPolicy, RBAC, Pod overrides. The platform must filter or override what the chart produces, which is fragile.
- **Friction:** capability owners must learn Helm. The platform contract effectively includes "the part of Helm we accept" — a moving target as Helm itself evolves.
- **Cost over time:** Helm CLI / chart API stability is now part of the platform contract.

### Option C — OCI image + a small platform-defined manifest schema

Capability owner ships an image and a small declarative manifest in a schema the platform defines. Initial fields: tenant/component name; image reference; ports; env (with secret refs by name per TR-11); declared CPU, memory, persistent-storage needs (TR-26); for migration artifacts, a re-run-contract flag (TR-31); optional extra containers within strict caps. The platform translates this into Kubernetes resources (Deployment / Service / NetworkPolicy / ResourceQuota / Job for migrations) at provision time.

- **TR-10:** met — one artifact form (image + manifest), used for both tenant components and migration jobs.
- **TR-26 / TR-31:** declarations are *fields in the schema*, present in the artifact submitted. Submission really is acceptance.
- **TR-32:** strong — schema does not expose NetworkPolicy / RBAC / Pod Security; the platform sets these from per-tenant defaults.
- **TR-22 / TR-24:** platform owns the K8s manifests it generates from the schema, satisfying both.
- **Friction:** capability owners write a small YAML; they don't need K8s fluency.
- **Cost over time:** the schema must be maintained as offerings grow. Additions are contract changes, which is the UX #5 path — the cost is *visible*, not hidden.

### Option D — Raw Kubernetes manifests

Capability owner ships a tarball / set of manifests (`Deployment.yaml`, `Service.yaml`, etc.) plus the image references they use.

- **TR-10:** met — one form (manifests).
- **TR-26 / TR-31:** could be enforced by the platform reading specific resource fields, but the contract is "the K8s API" — a huge surface that shifts under us.
- **TR-32:** fails by default — owners can ship NetworkPolicy and RBAC, sidestepping the platform's isolation. Only fixable by extensive admission policy that effectively re-implements Option C's filtering.
- **Friction:** maximum K8s fluency demanded of capability owners.
- **Cost over time:** every K8s API change is a contract change, whether the platform wanted one or not.

## Decision Outcome

Chosen option: **Option C — OCI image + a small platform-defined manifest schema**.

This option is chosen because:

- It is the only option where TR-26 and TR-31 declarations live *inside the artifact submitted*, making the host-a-capability promise — "submission is acceptance" — true mechanically and not by convention.
- It keeps TR-32 isolation primitives entirely on the platform side. Tenants cannot ship NetworkPolicy, RBAC, or Pod Security overrides; the platform sets those from per-tenant defaults derived from the schema.
- The same artifact form serves tenant components and migration jobs (TR-10 / TR-12), with the migration-specific bits expressed as schema fields (re-run-contract flag, declared spike vs. steady-state) rather than a parallel pipeline.
- The schema is a bounded, owned artifact: when offerings grow and tenants need new fields, the path forward is the contract-change UX (UX #5) — cost is visible and tracked, not hidden in chart-author conventions or K8s API drift.
- Capability owners are not required to learn Kubernetes, consistent with the host-a-capability persona that treats the capability owner as a *customer* of the platform and not a contributor to it.

The schema starts minimal and grows by contract change. The exact wire format (YAML vs. JSON vs. some other) and field-set are deferred to the working schema document that lives alongside the platform definitions; this ADR commits to "the schema is small, platform-defined, and grows only through UX #5."

### Consequences

- **Good, because** the artifact + schema is self-sufficient at submission time; the operator's review (host-a-capability §2) reads the schema fields directly rather than reconstructing intent from issue-thread discussion.
- **Good, because** the same form covers tenant components and migration jobs, so the migration-process offering (TR-12) does not need a parallel packaging path.
- **Good, because** TR-32 isolation primitives are entirely platform-side; misconfiguration risk lives in one place rather than every tenant submission.
- **Good, because** schema growth is governed: new fields ship through UX #5 (platform-contract-change rollout), so capability owners are migrated, not surprised.
- **Bad, because** the platform now owns and maintains a translator from the schema to Kubernetes manifests. That translator is platform code (or a controller — see Realization) and is on the operator's maintenance budget (TR-33).
- **Bad, because** capability owners with truly non-standard topology (e.g. statefulsets with niche volume layouts, headless services for peer-to-peer, daemonsets) hit the schema's expressivity ceiling. The intended response per the *capability evolves with its tenants* rule is to consider whether the platform should grow an offering — not to widen the schema indefinitely.
- **Requires:** a translator from the schema to K8s manifests, owned by the platform. This is realized either as a thin tool inside ADR #12's definitions tooling, or as a Kubernetes controller in `services/`. The choice between those is deferred to ADR #12; both are compatible with this ADR. ADR #9 (secrets) must define how secret-name references in the schema resolve to in-pod values. ADR #10 (migration runner) consumes the same schema with the migration-specific fields populated.

### Realization

How this decision shows up in the repo:

- **A working schema document** lives alongside the platform definitions (e.g. `platform-contract/tenant-manifest.md` plus a machine-readable schema file). It is the source of truth for what the platform accepts; capability owners read it, the translator validates against it.
- **An OCI image registry** the platform reads from. Whether this is operator-hosted or a public registry is deferred to the registry-placement decision (likely folded into ADR #6 or a small follow-on); the schema references images by registry-qualified name.
- **A translator** — either a CLI/library used by ADR #12's tooling, or a Kubernetes controller running in-cluster — converts a tenant manifest into the per-namespace K8s resources (Deployment / Service / Job / ResourceQuota / etc.). NetworkPolicy, RBAC bindings, and ResourceQuota templates are platform-owned and applied alongside, not derived from the tenant manifest.
- **Migration jobs** use the same schema with `kind: migration-job` (or equivalent) and the re-run-contract flag (TR-31) plus declared spike (TR-13). The translator emits a `Job` rather than a `Deployment`.
- **The canary tenant (TR-20)** is itself an instance of this schema, kept alongside the platform definitions, exercising the translator end-to-end as part of every rebuild.
- **Per-tenant manifests in the platform definitions repo** consist of: the canonical tenant manifest (the one submitted by the capability owner), the per-tenant overrides the platform applies (resource quota expressing TR-26 declarations and TR-13 caps), and the namespace + isolation-primitive set the platform owns.

## Open Questions

- **Schema wire format (YAML / JSON / protobuf).** Deferred to the working schema document.
- **Translator placement — CLI tool vs. in-cluster controller.** Deferred to ADR #12 (definitions tooling); both are compatible with this ADR. A controller has the advantage of keeping the schema → manifest translation continuous (TR-22 / TR-24); a CLI is simpler. Pick when the tooling is chosen.
- **Image registry.** Where tenant images live. The schema references registry-qualified image names; *which* registry is a deployment concern, not a schema concern, and is decided alongside ADR #6 (network) or a small follow-on.
- **Schema ceiling and the "capability evolves with its tenants" path.** When a tenant need exceeds the schema, the right response per the capability rule is sometimes to grow an offering and sometimes to widen the schema. The judgement is per case and lives in the host-a-capability "new offering needed" branch ([§3b](../user-experiences/host-a-capability.md#3-resolution--one-of-three-branches)).
