# Response

I drafted the ADR, but I want to flag a structural concern with the proposed decision before you consider it final — so I marked the ADR `proposed` rather than `accepted`.

## What I drafted

`docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0001-inter-service-authentication.md` — a MADR-format ADR following the conventions in `docs/content/r&d/adrs/0007-standard-api-error-response.md`.

## The concern

Your question was: *how do all services authenticate to each other?* The proposed answer (mTLS terminated at the Cloudflare edge, 30-day rotation) doesn't actually answer that question — it answers a different one (how does Cloudflare authenticate to the origin).

Specifically:

- **"Terminated at the edge" loses caller context.** Once TLS terminates at Cloudflare, the origin sees a connection from Cloudflare, not from the originating service. Service B cannot cryptographically verify the request came from service A.
- **TR-01 (tenant isolation) is not satisfied.** A compromised or hostile tenant workload inside the platform faces no cryptographic barrier to impersonating a platform service if east-west auth is undefined.
- **30 days is premature.** The right rotation cadence is a property of the issuer (SPIRE defaults to hours; ACME is typically 90 days; an internal CA could be anything). Pinning the cadence before picking the issuer puts the cart before the horse.
- **It depends on Q-02 and Q-04.** The compute substrate (Q-02) determines which workload-identity primitives are even available, and the identity-service decision (Q-04) determines the issuer/verifier. Both are listed as open questions in `tech-requirements.md` and should be decided first.

## What the ADR recommends instead

- **Defer** this ADR (Option D) until Q-02 and Q-04 land.
- Record a **preference for Option B** over the original proposal: keep edge mTLS for north-south (which is what `cloud/mtls/cloudflare-gcp/` already does), and add workload-identity-issued mTLS for east-west once the substrate is chosen.
- Also consider **Option C** (SPIFFE JWT-SVID / OIDC tokens) — tokens survive proxy termination, which sidesteps the edge-termination problem entirely.

## If you disagree

If you want to push through with the original proposal as-stated (Option A), the ADR documents it clearly as one of the considered options with explicit pros/cons, so it's a one-line edit to flip the decision outcome and change `status: proposed` to `accepted`. But I'd want you to at least see the trade-off written down before doing that.

## Files written

- `/home/carson/github.com/Zaba505/infra/.claude/worktrees/agent-a4f15ba964ad8502b/docs/content/capabilities/self-hosted-application-platform/tech-design/adrs/0001-inter-service-authentication.md`
- `/home/carson/github.com/Zaba505/infra/.claude/worktrees/agent-a4f15ba964ad8502b/docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (was already at the destination from the setup step)
