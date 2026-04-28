# ADR Plan — self-hosted-application-platform

Source: `docs/content/capabilities/self-hosted-application-platform/tech-requirements.md` (reviewed_at: 2026-04-26, unchanged).

Below is the enumeration of architectural decisions forced by the technical requirements. Each will be authored individually via `define-adr` after its issue is filed.

## Enumerated ADRs

### ADR-1: Tenant isolation model (data + compute boundary)
- **Drives:** TR-01, TR-03 (data-scope enforcement)
- **Decision:** Choose the isolation primitive (e.g., per-tenant GCP project vs. namespace-per-tenant on shared cluster vs. dedicated VM/runtime per tenant) and the data-layer scoping mechanism (per-tenant DB vs. row-level scope with enforced predicates).
- **Why now:** Every other decision (observability scoping, update strategy, export) inherits from this boundary.

### ADR-2: Platform contract versioning and concurrent-version support
- **Drives:** TR-02
- **Decision:** Choose the contract-versioning scheme (semver of contract bundle, version-pinned tenant manifests, deprecation window length) and the runtime mechanism that lets N versions run concurrently (router/dispatch by version, parallel control planes, etc.).

### ADR-3: Per-tenant observability pipeline and query scoping
- **Drives:** TR-03, partially TR-01
- **Decision:** Choose the telemetry stack (e.g., GCP Cloud Logging/Monitoring with per-tenant log buckets vs. self-hosted Loki/Tempo/Mimir with tenant labels) and the query-time scoping mechanism that prevents cross-tenant reads.

### ADR-4: Zero-downtime tenant update strategy
- **Drives:** TR-04
- **Decision:** Choose the update mechanism for online tenant workloads (blue/green via Cloud Run revisions, rolling with surge, traffic-split percentages) and the readiness/health contract tenants must implement.

### ADR-5: Tenant data export format and eviction/export window mechanics
- **Drives:** TR-05
- **Decision:** Choose the export format (portable schema — JSON/Parquet/protobuf bundle), the delivery channel (signed GCS URL vs. tenant-pull API), and the export-window lifecycle (trigger, duration, retention after eviction).

### ADR-6: Tenant data import / migration ingestion mechanism
- **Drives:** TR-06
- **Decision:** Choose the import API (batch upload vs. streaming), the idempotency key/fingerprint scheme, and the integrity-verification approach (per-record checksum, manifest hash) that guarantees no silent loss or duplicate-on-retry.

### ADR-7: Network topology conformance for platform & tenant traffic
- **Drives:** TR-07
- **Decision:** Codify how platform-internal and tenant-bound traffic is routed through Cloudflare → GCP (with WireGuard to home lab), including egress controls and the enforcement point that rejects non-conforming paths.

---

Total: 7 ADRs, one per architectural decision forced by the TRs. TR-01 is partially covered in ADR-3 because observability scoping inherits the isolation boundary set in ADR-1; the isolation primitive itself lives in ADR-1.

`gh-invocations.txt` contains the issue-creation commands (one per ADR). Not executed — written for review.
