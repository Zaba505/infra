# Response

Drafted ADR `0008-platform-inter-service-authentication.md` for the self-hosted-application-platform capability. The decision recorded matches what you specified: mTLS terminated at the Cloudflare edge with 30-day client-cert rotation. The ADR is tied back to TR-01 (tenant isolation) and TR-07 (Cloudflare → GCP path) from the reviewed `tech-requirements.md`, and it reuses the existing `cloud/mtls/cloudflare-gcp/` Terraform module pattern.

## What's in the ADR

- **Status:** `proposed` (date 2026-04-27, no deciders filled — please add yourself before accepting).
- **Category:** `strategic` (cross-cutting auth pattern across all platform services).
- **Considered options (4):** chosen edge-terminated mTLS + 3 alternatives (per-service mTLS, internal-JWT, network-only) with pros/cons each.
- **Confirmation:** Terraform-defined certs, observable rotation job, integration test for unauthenticated calls being rejected at the edge, code review against bespoke auth.
- **Consequences flagged honestly:** Cloudflare becomes a dependency for *internal* auth (not just ingress), and rotation automation must be built before 30 days is a meaningful cadence — these are real costs, not just upside.

## Things to flag before you accept it

1. **Scope question — capability ADR vs shared/strategic ADR.** "How all platform services authenticate to each other" is a cross-cutting decision that will likely affect *every future capability* hosted on the platform, not just self-hosted-application-platform itself. Per the repo's convention, `docs/content/r&d/adrs/` is where shared/strategic decisions live (0001–0007 are all there). I placed the new ADR in `outputs/` as you requested, but if you want this to bind future capabilities, consider moving it under `docs/content/r&d/adrs/` as `0008-...md` rather than under the capability folder. I left the numbering (`0008`) and category (`strategic`) consistent with that location so the move is a one-step `mv`.
2. **TR-07 already calls out a "prior shared decision"** for the Cloudflare → GCP path. This new ADR builds on that and arguably *should* live next to it.
3. **The forwarded-identity-header trust assumption** is the subtlest consequence. The destination GCP service trusts a header set by Cloudflare; this is only safe because the WireGuard/GCP path is non-tenant-reachable. Worth confirming that's documented somewhere durable (the ADR notes it, but the *enforcement* lives in network config).
4. **30-day rotation only pays off if rotation is automated end-to-end.** The ADR is explicit about this in Consequences and Confirmation, but it's worth knowing that until that automation exists, the effective rotation cadence will be "whenever the operator remembers", which is not 30 days.

## Files

- `/home/carson/github.com/Zaba505/infra/.claude/skills/define-adr-workspace/iteration-1/eval-0-shared-decision/without_skill/outputs/0008-platform-inter-service-authentication.md`
