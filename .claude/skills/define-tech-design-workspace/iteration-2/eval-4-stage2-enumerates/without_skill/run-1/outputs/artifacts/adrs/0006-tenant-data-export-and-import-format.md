---
title: "[0006] Tenant Data Export and Migration Import Format"
description: >
    Choose a single portable archive format and integrity scheme that serves both tenant data export (eviction / move-off) and migration import (incoming tenants).
type: docs
weight: 6
category: "strategic"
status: "proposed"
date: 2026-04-26
deciders: [operator]
consulted: []
informed: []
---

## Context and Problem Statement

[TR-05] requires that an evicted tenant be able to retrieve all of their data in a portable format within a defined export window. [TR-06] requires that new tenants be able to migrate existing data into the platform without loss or corruption — idempotently, with verifiable integrity.

These are two ends of one pipe: the *format* on the wire, the *integrity* check, and the *idempotency* property are the same problem viewed from opposite sides. This ADR picks the format and the integrity scheme that serve both directions.

## Decision Drivers

* TR-05: portable format, bounded export window
* TR-06: idempotent import (no silent loss, no duplication on retry), verifiable integrity
* The 30-day post-eviction grace window (NFR-06 long-form) — exports must be self-contained; the platform may not be reachable when the tenant restores
* NFR-02: maintenance budget — one format, not one-per-tenant
* Tiebreaker: vendor independence > minimizing operator effort — the format must not require the platform to read it back

## Considered Options

* **Per-tenant custom format (each tenant defines their own dump shape)** — maximum flexibility, no convergence.
* **Database-native dumps (pg_dump, mysqldump, etc., per backing store)** — tenant gets the raw dumps of whatever they were running on.
* **Tar archive with manifest + per-file SHA-256 + content-addressed object layout** — single format across tenants; importer keys on content hash for idempotency.

## Decision Outcome

Chosen option: **tar archive with a top-level `manifest.json` listing every file's relative path, size, and SHA-256, plus a content-addressed `objects/` layout where the file body is stored under its own hash**.

Export produces:

```
<tenant-id>-<export-timestamp>.tar
├── manifest.json          # { files: [{path, size, sha256}, …], total_bytes, schema_version }
├── manifest.json.sig      # operator-key signature over manifest.json
└── objects/
    ├── ab/cd…ef           # content-addressed; path under objects/ is the SHA-256 split for fs friendliness
    └── …
```

Import (TR-06) is idempotent because the importer reads `manifest.json`, computes the SHA-256 of each incoming object, and skips any object whose hash already exists in the target store. A retry of a partially-completed import resumes from the first missing hash. Verifiable integrity (TR-06) is the manifest signature plus the per-file hashes. Portability (TR-05) is satisfied because `tar` + `sha256sum` are universal — an evicted tenant does not need any platform-supplied tool to verify or use the archive.

Database-native dumps were rejected because they leak the platform's storage choices to the tenant (anti-portability) and because each backing store would need a different importer. Per-tenant custom formats were rejected because they violate NFR-02 the moment there are two tenants.

### Consequences

* Good, because TR-05 portability is satisfied with universal tools (`tar`, `sha256sum`, `openssl`)
* Good, because TR-06 idempotency is structural — content-addressed storage *is* dedup
* Good, because the importer and exporter share the same manifest schema; one piece of code, two directions
* Good, because the manifest signature (TR-06: verifiable integrity) lets a re-imported archive be authenticated against the operator's public key without trusting the transport
* Neutral, because tenant-specific data semantics (e.g. database schema) live *inside* the objects; the platform treats the contents opaquely
* Bad, because `objects/` layout means an evicted tenant who wants to *use* their data must understand the manifest to reconstruct paths — the manifest doc is part of the export
* Bad, because content-addressing inflates very-many-small-files exports modestly; mitigated by tar's natural batching

### Confirmation

* The export tool is part of the platform's tenant-lifecycle code; exercised by the canary (REQ-18 long-form) — canary export → verify manifest signature → import to a fresh tenant slot → byte-equal check
* The import tool is exercised by `migrate-existing-data` (TR-06's UX); the issue template includes a manifest-validation step
* A documented procedure for an evicted tenant to verify their archive without platform access lives in the move-off-the-platform-after-eviction UX doc
* Schema versioning of `manifest.json` follows the contract-versioning scheme from ADR-0003 (`schema_version: "2026-04"`)

## Pros and Cons of the Options

### Per-tenant custom format

* Good, because each tenant gets exactly the dump shape they want
* Bad, because operator maintains N export/import pipelines — fails NFR-02
* Bad, because no shared verification story
* Bad, because new tenants increase operator surface linearly

### Database-native dumps (pg_dump, etc.)

* Good, because the dump is a first-class citizen for whatever DB the tenant uses
* Bad, because it leaks the platform's storage choices into the tenant's archive — anti-portability (TR-05)
* Bad, because the platform must include the matching version of the DB tooling in the export, or the tenant must
* Bad, because the importer must understand every supported DB's dump format

### Tar archive with manifest + SHA-256 + content-addressed objects

* Good, because portable with universal tools
* Good, because content addressing gives idempotency for free
* Good, because manifest signing gives verifiable integrity
* Good, because one format covers both directions
* Neutral, because tenants needing to *use* the archive must read the manifest (documented)
* Bad, because content-addressed layout has modest small-file overhead

## More Information

* The export window (TR-05) is the 30-day post-eviction grace window from NFR-06 long-form — the export tool remains usable for that window
* The signing key is the operator's platform key, the same one used for ingress signing in ADR-0004; key rotation procedures cover both
* Encryption at rest of the archive is out of scope here — tenants who need it wrap the tar themselves; the platform does not assume tenant-side key management
