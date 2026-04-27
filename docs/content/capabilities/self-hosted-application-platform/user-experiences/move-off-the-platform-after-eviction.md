---
title: "Move Off the Platform After Eviction"
description: >
    A capability owner whose capability has been evicted gets their data out cleanly and walks away with no obligations and nothing left to access on the platform before the retention window closes.
type: docs
---

> **One-line definition:** A capability owner whose capability has been evicted gets their data out cleanly and walks away with no obligations and nothing left to access on the platform before the retention window closes.

**Parent capability:** [Self-Hosted Application Platform](../_index.md)

## Persona

The actor here is a **capability owner** whose capability has been evicted — a *Primary actor (initiator)* from the parent capability's Stakeholders, on the way out. As elsewhere in this capability's UX docs, the role is treated as separate from the operator's even though today both hats are worn by the same person.

- **Role:** Capability owner. The party who originally onboarded a capability onto the platform via `host-a-capability`, has been hosting it for some period, and is now being removed.
- **Context they come from:** The parting is **amicable**. Eviction was triggered by a divergence the platform legitimately cannot meet — specialized hardware, regulatory constraints, an availability target stronger than the platform offers — *not* by a missed deadline in the `operator-initiated-tenant-update` flow. Negotiation over the eviction date has already happened upstream, before this UX begins. The capability owner accepts that they are leaving and has agreed to the date.
- **What they care about here:** A clean exit. By the eviction date their capability is fully off the platform, their data is in their hands in a portable form they can verify, and nothing of theirs lingers on the platform afterward. They are *not* asking the platform to help them figure out where to run next — that is their problem to solve.

## Goal

> "By the time the platform is finished with my capability, I have my data, I know it's complete, and I have nothing left to chase down here."

## Entry Point

The capability owner arrives at this experience because the operator has filed an **eviction issue** against the infra repo tagging them. The issue contains exactly:

- The eviction date (already negotiated upstream — not up for renegotiation in this journey).
- The reason for eviction (so it is on the record and the parting stays amicable).
- A link to the platform's export tooling, with documentation on how to use it and what the export shape looks like for their tenant.

That is all the issue carries. The capability owner's state of mind is "the date is set, I know where the export tool is, I have a window of time to get my data out and walk away cleanly."

## Journey

The journey runs in three phases keyed off the eviction date: a pre-eviction window where the tenant is still live, the eviction date itself when compute and network resources go away, and a 30-day grace window where data is held in an export-only, read-only state before being permanently removed.

### Phase A — Before the eviction date (tenant still live)

#### 1. Read the eviction issue and the export documentation

The capability owner reads the issue, follows the link to the export tooling, and reads its documentation. They learn what the export will produce — file layout, formats, what is included, what is not — and roughly how long an export of their dataset will take to run. No back-and-forth with the operator is expected here; the issue and the docs are meant to be self-sufficient.

#### 2. Notify their own end users

The capability owner tells *their* end users that the capability is going away on the eviction date — separately from the platform, on whatever channel they use with their users. The platform plays no role here; end users of a tenant capability are not visible to the platform and the platform does not communicate with them. (See *No direct end-user access to the platform* in Constraints.)

#### 3. Run the export and verify it themselves

The capability owner kicks off the export using the platform's export tool. What they perceive is an archive of their tenant's data, produced for them to download then and there, plus a checksum/hash and total size in bytes that the platform produces alongside it. **Validation that the export is complete and correct is the capability owner's responsibility, not the platform's.** Only the capability owner knows their data well enough to say "yes, this is all of it and it is intact." The platform offers checksum/hash and total size as the ceiling of what it can verify on the capability owner's behalf — anything beyond that (record counts, schema integrity, business invariants) is theirs.

#### 4. (Optional) Run the export iteratively

Because end users may still be writing data while the tenant is live, an export taken in Phase A is not necessarily the *final* export. The capability owner may run multiple exports across Phase A — one early to validate that the tooling produces something usable, another later to capture more recent writes. Whether they do this is their call; the platform supports it because the export tool simply runs whenever invoked. Each run is ephemeral: if they want to keep an export, they download it when it is produced. The platform does not keep a history of prior exports around for them.

### Phase B — The eviction date

#### 5. Compute and network resources are torn down; the tenant stops serving

On the eviction date the operator deprovisions the tenant's compute, network, and other live resources. From the capability owner's seat: the tenant is no longer reachable by their end users. The data persists, but only in an export-only, read-only state — no further writes can occur, by anyone. A comment is posted on the eviction issue confirming the cutover and the start of the 30-day retention window.

What the capability owner perceives: the issue gets a status comment, and they now know their dataset is frozen. If they had not finished extracting data before this point, they still have 30 days — but the dataset they extract from now on is the *final* one.

### Phase C — Post-eviction (30-day retention window)

#### 6. Run the export of record (if not already taken)

In Phase C the export tool still works, but now against a stable, read-only snapshot. For capability owners with more data than they could extract during Phase A, or for those who deliberately deferred to avoid racing live writes, this is when the **definitive** export is pulled. They re-run the same export tool, get back an archive plus checksum/hash and size, download it, and validate it the same way they validated in Phase A.

For capability owners who already pulled what they needed in Phase A, Phase C is a safety net — "I forgot a thing, let me grab it" — rather than the main event.

#### 7. Walk away

Once the capability owner is satisfied they have everything, they comment on the issue indicating they are done. The operator closes the issue. After 30 days from the eviction date, the platform stops offering any tenant-accessible copy of their data regardless of whether the capability owner ever closed the loop. Residual copies may still exist only inside the platform's normal backup-retention machinery and are not accessible back to the tenant. There is no "are you sure?" — the 30-day clock is hard.

### Flow Diagram

```mermaid
flowchart TD
    Start([Eviction issue filed by operator<br/>date already negotiated]) --> Read[Read issue + export tooling docs]
    Read --> Notify[Notify own end users<br/>off-platform]
    Notify --> ExportLive[Run export against live tenant<br/>verify checksum / size / contents]
    ExportLive --> Iter{More writes expected<br/>before eviction date?}
    Iter -->|Yes| ExportLive
    Iter -->|No| Wait[Wait for eviction date]
    Wait --> Cutover[Eviction date:<br/>compute/network torn down<br/>data → read-only<br/>comment posted on issue]
    Cutover --> Phase C{Need more data<br/>from final snapshot?}
    Phase C -->|Yes| ExportFinal[Run export against frozen snapshot<br/>verify the same way]
    Phase C -->|No, already complete| Done
    ExportFinal --> Done[Comment 'done' on issue;<br/>operator closes it]
    Done --> Permanent([30 days post-eviction:<br/>no tenant-accessible copy remains])
```

## Success

When the journey ends cleanly, the capability owner walks away with:

- A verified, complete archive of their tenant's data, sized and checksummed by the platform, validated by them.
- A clear paper trail on the eviction issue showing the date, the reason, and confirmation that they pulled what they needed.
- Nothing left to chase down on the platform. After the 30-day window the platform offers no tenant-accessible copy of their data; any deeper residual backup copies are outside the tenant's reach and simply age out on the platform's normal retention schedule.
- An amicable ending. The operator filed the issue, the platform held the data the agreed amount of time, and the capability owner left under their own power. The relationship is intact for whatever comes next.

## Edge Cases & Failure Modes

- **Capability owner asks for more time after the eviction date.** Hard wall. The negotiation over the eviction date happened upstream of this journey; once that date is set, it is the date. The 30-day post-eviction retention is the only post-date slack and it is fixed.
- **Export takes longer than 30 days to actually run on a very large dataset.** Same hard wall — the capability owner had Phase A *plus* 30 days of Phase C to extract; if that is not enough, they had advance warning during eviction-date negotiation and should have raised it then. The platform does not extend the retention window for slow extracts.
- **Export comes back wrong** (checksum mismatch, missing files, corruption visible to the capability owner). The capability owner reports the problem on the eviction issue so that thread remains the coordination record. This is the *one* exception to the 30-day hard wall: if the failure is shown to be in the platform's export tooling or its data hosting, the operator pauses the permanent-removal clock for that tenant until the platform-side issue is resolved and a clean export has been produced. No separate restoration SLA is promised in this UX; the issue stays open until the capability owner can pull a clean export. Failures rooted in the capability owner's own validation steps do not pause the clock.
- **Export tooling does not exist for this tenant's data shape at the time of eviction.** Cannot happen by design — export tooling is a **core platform feature**, present for every kind of data the platform hosts. If a hole is discovered, that is itself a platform bug, handled the same way as the previous bullet (eviction issue remains open, retention clock paused).
- **Capability owner ignores the issue entirely and never extracts anything.** No special handling. The 30-day clock runs, tenant-accessible data is removed, the issue is closed by the operator. The capability owner may have made themselves whole through other means (their own backups, accepting the loss); the platform does not chase them.
- **End users keep hitting the tenant after the eviction date.** They get whatever connection failure the underlying infra produces. The capability owner is responsible for having warned their end users; the platform does not present a "this tenant has been retired" page or otherwise communicate with end users — end users belong to the capability, and from the platform's seat, *the capability is the end user*.
- **Capability owner wants to come back later** (re-onboard the same capability after the divergence is resolved). That is a *new* `host-a-capability` journey, not a continuation of this one. It is not blocked, but nothing about this UX preserves state to make it easier.

## Constraints Inherited from the Capability

This UX must respect the following items from the parent capability — by name:

- **Eviction is allowed when needs and capabilities diverge.** This UX is the operationalization of the *amicable* form of that rule: the divergence is real (specialized hardware, regulatory constraint, availability target the platform cannot meet) and the parting is mutual. The fall-behind variant of eviction is handled separately via `operator-initiated-tenant-update`.
- **No direct end-user access to the platform.** End users of the tenant capability are not visible to the platform and are not communicated with by the platform during eviction. Notification of end users is purely the capability owner's responsibility.
- **Operator succession — on-demand exportable archives.** The same export mechanism that the parent capability promises for operator-succession scenarios is what powers this journey. Export tooling is therefore not bespoke to eviction; it is a core platform feature that exists at all times for every tenant. This UX simply consumes it.
- **Operator-only operation.** The capability owner has no administrative access during this journey. Everything they do — running exports, leaving comments — is done through the same surfaces an end-state non-operator has. The operator is the one who deprovisions resources and closes the issue.
- **Affected parties (end users of the tenant capability).** End users feel this journey indirectly: their access to the capability ends on the eviction date. The platform does not surface this to them — the *capability owner* does, separately, on their own channels.
- **KPI: 2-hr/week operator maintenance budget.** Implication: this journey must not require the operator to do bespoke per-tenant work. The export tool is generic and runs on demand; the operator's only routine touchpoints are filing the issue, posting the cutover comment, and closing the issue at the end. A tenant whose eviction would require custom export work is itself a sign the platform's export tooling has a gap that needs fixing — handled as a platform bug, not as an operator-effort overrun.
- **KPI: 1-hour reproducibility.** Implication: the data formats produced by the export tool, and the way they relate to the platform's definitions, should be expressible as part of the platform itself, not as snowflake per-tenant logic. (Standing the platform up should not require remembering "and here is the special export path for tenant X.")

## Out of Scope

- **The eviction-decision journey itself.** *Why* the operator decided to evict, and the conversation that established the eviction date, happens before this UX. By the time this UX begins, the issue is filed, the date is set, and both parties have agreed.
- **The fall-behind eviction path.** Eviction triggered by a missed extended date in `operator-initiated-tenant-update` is a different shape (less amicable, possibly compressed timelines). It enters a separate journey not covered here, even though the *mechanics* of getting data out via the export tool may overlap.
- **Helping the capability owner figure out where to run next.** The platform does not point at alternative hosts, port the capability's runtime, or assist with migration. The export tool produces data; the rest is the capability owner's problem.
- **Application/runtime/configuration migration tooling.** Only *data* export is provided. Capability code, container images, configuration, secrets management at the destination — none of this is the platform's concern.
- **Re-onboarding the same capability later.** If the capability owner wants to come back, that is a fresh `host-a-capability` journey with no special path inherited from having previously been here.
- **Operator's side of this journey.** This UX is written from the capability owner's seat. The operator's experience (filing the issue, deprovisioning on the date, posting the cutover comment, closing the issue, watching the 30-day clock) is captured here as a responder, not as a separate document.

## Open Questions

_None at this time._
