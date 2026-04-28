# Plan: Operator Stands Up the Hosting Platform — for `self-hosted-personal-media-storage`

**Capability reviewed:** `.claude/skills/plan-experiences/evals/fixtures/self-hosted-personal-media-storage.md` (reviewed_at: 2026-04-28)

---

## Important: scope concern with the requested journey

The task asks to plan an experience for "the operator standing up the hosting platform from scratch." That phrasing maps more cleanly to the sibling capability `self-hosted-application-platform`, which already owns a UX at `docs/content/capabilities/self-hosted-application-platform/user-experiences/stand-up-the-platform.md`. The fixture for `self-hosted-personal-media-storage` is explicit:

> **Out of Scope — General-purpose self-hosted application hosting** (e.g. running a Minecraft server, hosting other apps for the operator's circle). This is a separate capability and will be defined in its own document.

So a literal "stand up the hosting platform" UX does not belong here. Two reasonable interpretations:

1. **Misplacement** — the requested journey should live under `self-hosted-application-platform` (likely already does). No new UX is needed for `self-hosted-personal-media-storage`.
2. **Reframing** — the operator does need a way to *initially deploy the media-storage service itself* (provision storage, install the application, configure auth, invite the first user). That is in scope for this capability and is a valid UX.

The plan below covers interpretation (2): an operator-facing "first-time deploy of the media-storage service" experience. The proposed issue body explicitly flags the scope concern and asks the human to confirm before authoring proceeds.

---

## Proposed experience (one issue)

### UX-1: Operator deploys the media-storage service for the first time

- **Primary actor:** The operator (sole accountable party).
- **Trigger:** The operator decides to start running the capability and has not yet deployed it.
- **Goal:** A running, reachable instance of the media-storage service exists, the operator holds working admin credentials, the successor handoff is sealed, and the operator can invite the first user.
- **Why it's in scope:** The capability's outcomes (privacy, longevity, control) cannot be delivered until the system exists; the fixture's "Operator succession" rule (sealed successor credentials, runbook) implies a deploy step where those artifacts are produced. This UX covers the operator-only path to "the system is running" — it does NOT cover hosting other applications, which is out of scope.
- **Out of scope for this UX:**
  - Inviting/provisioning end users (separate UX: user provisioning).
  - Routine operations, upgrades, capacity expansion (separate UX: ongoing operation).
  - Off-site backup setup (related rule but a distinct UX).
  - Anything the `self-hosted-application-platform` capability covers — if that capability is used as the substrate, this UX *consumes* it rather than re-implementing it.
- **Key business rules in play:**
  - Closed user set — deploy must end in a state where only the operator has access.
  - Private by default — initial config must default to private; the operator must not be able to read user content even after deploy.
  - Operator succession — successor credentials must be produced, sealed, and not used for routine operation.
  - Lost credentials = lost data — the deploy flow must make this trade-off visible to the operator and confirm the operator has stored their own credentials safely.
- **Success signal:** Operator can authenticate, the privacy posture is verifiable, the successor handoff is sealed, and a runbook exists for the successor.

---

## GitHub commands (printed, NOT executed)

The following `gh issue create` command would file the planned UX issue. It is printed only.

```
gh issue create \
  --title "UX: Operator deploys the media-storage service for the first time (self-hosted-personal-media-storage)" \
  --label "user-experience,capability:self-hosted-personal-media-storage" \
  --body "$(cat <<'EOF'
## Context

Capability: `self-hosted-personal-media-storage`
Source: `.claude/skills/plan-experiences/evals/fixtures/self-hosted-personal-media-storage.md` (reviewed_at: 2026-04-28)

## Scope concern — confirm before authoring

The originating request asked for "the operator standing up the hosting platform from scratch." That phrasing overlaps with the sibling capability `self-hosted-application-platform`, which already owns a `stand-up-the-platform` UX. The fixture for *this* capability marks "general-purpose self-hosted application hosting" as out of scope.

This issue reframes the journey as **deploying the media-storage service itself for the first time** — which is in scope. Before authoring, please confirm:

- [ ] The intent is the media-storage service deploy (not the general application platform).
- [ ] If the media-storage service is built on top of `self-hosted-application-platform`, this UX *consumes* that capability's stand-up UX rather than duplicating it.
- [ ] Otherwise, close this issue and direct the work to `self-hosted-application-platform`'s existing UX.

## Proposed UX

**Title:** Operator deploys the media-storage service for the first time

**Primary actor:** Operator (sole accountable party)

**Trigger:** Operator decides to start running the capability; nothing is deployed yet.

**Goal state:** A running, reachable media-storage service exists; operator holds working admin credentials; successor credentials are sealed and escrowed; a runbook exists; the operator is ready to invite the first user.

**Key business rules in play:**
- Closed user set — deploy ends with only the operator having access.
- Private by default — operator cannot read user content post-deploy.
- Operator succession — successor credentials produced, sealed, not used routinely.
- Lost credentials = lost data — flow surfaces this trade-off to the operator.

**Out of scope for this UX:**
- Inviting/provisioning end users.
- Routine operation, upgrades, capacity expansion.
- Off-site backup configuration.
- General-purpose application hosting (separate capability).

**Success signal:** Operator authenticates; privacy posture is verifiable; successor handoff is sealed; runbook exists.

## Next step

Author this UX via the `define-user-experience` skill once the scope concern above is resolved.
EOF
)"
```

---

## Summary of decisions

- **One UX proposed**, not multiple — the operator-deploy journey is a single coherent end-to-end flow.
- **Scope concern flagged** in the issue body so the human can redirect to `self-hosted-application-platform` if that's actually where the work belongs.
- **No `gh issue create` was executed.** The command above is printed for the human to run (or discard) after confirming scope.
