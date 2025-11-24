---
title: "[0006] Universal Resource Identifier Standard"
description: >
    Select a consistent resource identifier format to be used across all types of resources in the system for tracking, linking, and reference.
type: docs
weight: 6
category: "strategic"
status: "accepted"
date: 2025-11-23
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions (frameworks, auth strategies, cross-cutting patterns)
- user-journey: Solutions for specific user journey problems (feature implementation approaches)
- api-design: API endpoint design decisions (pagination, filtering, bulk operations)
-->

## Context and Problem Statement

As the infrastructure and services grow, we need a consistent way to identify resources across different systems, databases, APIs, and logs. Resources include servers, VMs, containers, services, configurations, and any other trackable entities. Without a standardized identifier format, we risk inconsistencies in references, difficulty correlating logs and traces, and challenges in cross-system integration.

What identifier format should we standardize on for all resources to ensure uniqueness, readability, and compatibility across our entire infrastructure?

<!-- This is an optional element. Feel free to remove. -->
## Decision Drivers

* Global uniqueness across all resources and systems
* Human readability for debugging and operational tasks
* Compatibility with various systems (databases, URLs, file systems, logs)
* Sortability and ability to extract creation time information
* Low collision probability without coordination
* Performance considerations for generation and storage
* Industry standard practices and ecosystem tooling support

## Considered Options

* UUIDv4 (Random UUID)
* UUIDv7 (Time-ordered UUID)
* ULID (Universally Unique Lexicographically Sortable Identifier)
* Custom sequential IDs
* Snowflake IDs

## Decision Outcome

Chosen option: "UUIDv7 (Time-ordered UUID)", because it provides the best balance of standardization, performance, and functionality. UUIDv7 maintains full UUID compatibility while addressing the major weaknesses of UUIDv4 (poor database index performance and lack of temporal ordering). As a new IETF standard (RFC 9562), it has growing ecosystem support and is becoming available in standard libraries, making it a future-proof choice.

<!-- This is an optional element. Feel free to remove. -->
### Consequences

* Good, because time-ordered IDs improve database index performance and reduce fragmentation
* Good, because sortability by creation time simplifies debugging and operational tasks
* Good, because extractable timestamps provide valuable metadata without additional storage
* Good, because full UUID compatibility ensures broad system and tool support
* Good, because no coordination needed between distributed systems for generation
* Good, because growing standard library support reduces dependency burden
* Neutral, because migration from existing UUIDv4 identifiers will need to be managed incrementally
* Bad, because not all systems have UUIDv7 support yet (requires library updates or polyfills)
* Bad, because IDs are still 36 characters, less compact than alternatives like ULID

<!-- This is an optional element. Feel free to remove. -->
### Confirmation

Implementation compliance will be confirmed by:
* Code reviews ensuring all new resource types use the standardized identifier format
* Linting rules or code generation templates that enforce the identifier format
* Documentation updates reflecting the standard identifier format
* Database schema reviews to verify proper column types and indexing

<!-- This is an optional element. Feel free to remove. -->
## Pros and Cons of the Options

### UUIDv4 (Random UUID)

Standard 128-bit random identifier following RFC 4122, format: `550e8400-e29b-41d4-a716-446655440000`

* Good, because widely supported across all languages, databases, and systems
* Good, because cryptographically random with extremely low collision probability
* Good, because no coordination needed between systems for generation
* Good, because standard library support in most languages
* Neutral, because 36 characters (with hyphens) or 32 hex characters
* Bad, because not time-ordered, leading to poor database index performance
* Bad, because not sortable by creation time
* Bad, because difficult for humans to read or verify
* Bad, because provides no temporal information

### UUIDv7 (Time-ordered UUID)

Latest UUID standard (RFC 9562) with millisecond timestamp prefix, format: `018c7dbd-9265-7000-8000-123456789abc`

* Good, because maintains UUID compatibility while adding time-ordering
* Good, because better database index performance than UUIDv4 due to ordering
* Good, because sortable by creation time (millisecond precision)
* Good, because extractable timestamp for debugging
* Good, because growing ecosystem support and standard library adoption
* Good, because globally unique without coordination
* Neutral, because same 36/32 character length as UUIDv4
* Neutral, because newer standard, not yet universally supported in all systems
* Bad, because still not particularly human-readable

### ULID (Universally Unique Lexicographically Sortable Identifier)

26-character base32 encoded identifier with millisecond timestamp, format: `01ARZ3NDEKTSV4RRFFQ69G5FAV`

* Good, because lexicographically sortable by creation time
* Good, because shorter than UUID (26 chars vs 36) making it more readable
* Good, because case-insensitive base32 encoding avoids ambiguous characters
* Good, because extractable millisecond timestamp
* Good, because URL-safe without escaping
* Good, because better database index performance due to ordering
* Neutral, because requires library support (not in standard libraries yet)
* Neutral, because growing adoption but not as universal as UUID
* Bad, because not an official IETF standard (though widely used)
* Bad, because 26 characters still not trivially human-memorable

### Custom sequential IDs

System-specific sequential numbering (e.g., `SRV-00001`, `VM-12345`)

* Good, because very human-readable and memorable
* Good, because short and compact
* Good, because naturally sortable
* Bad, because requires centralized coordination for uniqueness
* Bad, because difficult to merge or sync across distributed systems
* Bad, because exposes information about total resource count
* Bad, because not globally unique without namespace management
* Bad, because potential security concerns (predictability, enumeration)

### Snowflake IDs

64-bit IDs with timestamp, datacenter, and sequence components (Twitter Snowflake pattern)

* Good, because compact 64-bit integer format
* Good, because time-ordered and sortable
* Good, because efficient storage and indexing
* Good, because high performance generation
* Neutral, because requires coordination for datacenter/worker IDs
* Bad, because not human-readable (large integers)
* Bad, because requires infrastructure for ID generation service
* Bad, because 64-bit limit may be reached for very high-volume systems
* Bad, because less portable across systems than string-based identifiers

<!-- This is an optional element. Feel free to remove. -->
## More Information

* UUIDv7 specification: RFC 9562 - https://datatracker.ietf.org/doc/html/rfc9562
* UUIDv4 specification: RFC 4122 - https://datatracker.ietf.org/doc/html/rfc4122
* ULID specification: https://github.com/ulid/spec
* Considerations for database indexing performance with different ID types
* Impact on API design, URL structure, and client implementation
* Migration strategy for existing resources using different identifier formats
