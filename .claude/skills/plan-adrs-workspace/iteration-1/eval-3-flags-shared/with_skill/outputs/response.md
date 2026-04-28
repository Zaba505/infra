# Response

That ADR is a shared decision, not a capability-scoped one — it touches both Cloudflare topology (where TLS terminates, edge cert lifecycle) and inter-service identity (how every platform service authenticates to every other). Both ripple across every capability hosted on the platform, not just `self-hosted-application-platform`. Shared decisions like this belong in `docs/content/r&d/adrs/` and go through a separate flow, not this skill.

There's a second issue: "let's go with mTLS terminated at Cloudflare edge with 30-day rotation" is the answer. The `plan-adrs` skill only enumerates decisions to be made — options and rationale are `define-adr`'s job. Even if this were capability-scoped, I wouldn't file an issue that pre-commits to mTLS-at-the-edge-with-30-day-rotation; that's the kind of choice the ADR itself needs to weigh against alternatives (mTLS terminated inside the cluster, SPIFFE/SPIRE workload identity, GCP IAM-based service-to-service auth, etc.) with the trade-offs documented.

So I'm stopping without filing anything. Two things to decide:

1. **Defer the shared ADR** and tell me which *capability-scoped* decisions you want to plan for `self-hosted-application-platform` against its TRs — I'll re-read `tech-requirements.md` and propose a list grounded in TR-NNs.
2. **Pause this skill** and handle the shared inter-service-auth ADR via the shared-ADR flow under `docs/content/r&d/adrs/`. When you do, frame it as "decide how platform services authenticate to each other," not "adopt mTLS at the edge" — let the ADR record the option set and the rationale.

Which way do you want to go?

## Gate state (for the record)

I did not run the review-gate check on `tech-requirements.md` because the request never reached the enumeration step — the proposed ADR is out of scope for this skill regardless of gate state. If you choose option 1, I'll check `reviewed_at` and the file mtime before proposing a list.

## Artifacts

- `gh-invocations.txt` — empty (no issues filed); contains a note explaining why.
