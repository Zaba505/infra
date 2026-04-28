Stopping — I won't enumerate experiences yet.

I located the capability doc at `docs/content/capabilities/self-hosted-personal-media-storage.md` and inspected its frontmatter. The frontmatter has no `reviewed_at` field at all:

```yaml
---
title: "Self-Hosted Personal Media Storage"
description: >
    A self-operated place where the operator, family, and friends can store and share personal media under their own control.
type: docs
weight: 10
---
```

The `reviewed_at` field is the gate for this skill. Without it, I have no signal that the current contents of the doc have been human-reviewed, and planning experiences against an unreviewed capability would source journeys from a draft — exactly what the gate exists to prevent.

**To proceed, do one of the following and then re-invoke me:**

1. Read the capability doc end-to-end, confirm it reflects the capability you want to plan experiences against, and add `reviewed_at: 2026-04-28` to the frontmatter, then re-invoke me. **Or,**
2. Tell me verbally that you've reviewed the current contents and I'll record your confirmation by setting `reviewed_at: 2026-04-28` on your behalf before proceeding.

If on review you find the capability doc needs changes first, the right next step is `define-capability`, not this skill — update the doc, then set `reviewed_at` to today's ISO date, then come back here.

No issues were filed and no list of experiences was proposed; the gate must be satisfied first.
