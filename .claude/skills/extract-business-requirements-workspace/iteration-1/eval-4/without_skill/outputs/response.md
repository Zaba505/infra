Re-extracted `business-requirements.md` for the self-hosted-application-platform capability.

**What changed vs. the existing file:**

- **BR-05 source link fixed.** The previous link pointed to `move-off-the-platform-after-eviction.md#a-section-that-no-longer-exists`, which no longer resolves in the current UX doc. The corresponding section in the live UX is now `## Edge Cases & Failure Modes`, so the link is updated to `move-off-the-platform-after-eviction.md#edge-cases--failure-modes` and the source label is updated to `§Edge Cases & Failure Modes` to match the section title verbatim.
- **`reviewed_at` bumped** from `2026-04-20` to `2026-04-28` to reflect this re-extraction.
- **Anchors made explicit** on BR-03 (`#journey`) and BR-04 (`#success`). The previous file used `§Journey` / `§Success` in the link text but did not include the anchor fragment in the URL itself; per the cross-linking convention (section links require explicit anchors), they now do.

**No requirements were added or removed.** BR-01 through BR-05 are unchanged in substance; the capability's Business Rules and the five UX docs still imply exactly the same business requirements they did on 2026-04-20. The two open questions are also unchanged.

**Verification done:**

- All five source links resolve against the current docs:
  - `_index.md#business-rules` -> exists (Business Rules & Constraints).
  - `platform-contract-change-rollout.md` -> exists.
  - `tenant-facing-observability.md#journey` -> exists.
  - `operator-initiated-tenant-update.md#success` -> exists.
  - `move-off-the-platform-after-eviction.md#edge-cases--failure-modes` -> exists.

The file is ready for human review before the next stage (technical requirements).
