---
title: "Universally Unique Lexicographically Sortable Identifier (ULID) Analysis"
type: docs
weight: 2
---

## Overview

ULID is a community-driven specification for unique identifiers that combine the decentralized generation of UUIDs with the performance benefits of time-ordered sequential IDs. Created as an alternative to UUID v4's poor database performance, ULID predates UUID v7 but shares similar design goals.

## URI Safety

### ✅ Completely URI-Safe

ULIDs are designed with URI usage as a primary consideration.

**Character set:**
- Uses Crockford's Base32 alphabet
- Characters: `0123456789ABCDEFGHJKMNPQRSTVWXYZ`
- Excluded: `I`, `L`, `O`, `U` (avoid confusion and potential abuse)
- 32 unique characters

**Format:**
```
01ARZ3NDEKTSV4RRFFQ69G5FAV
```

**Characteristics:**
- **26 characters** (10 timestamp + 16 randomness)
- **No hyphens** (unlike UUID's 36 chars with hyphens)
- **Case-insensitive** (can be normalized)
- **More compact** than UUID string representation

**Advantages over UUID:**
- Shorter (26 vs 36 characters)
- No special characters required
- More human-readable
- Case-insensitive (easier to communicate verbally)

**Usage in URIs:**
```
/api/users/01ARZ3NDEKTSV4RRFFQ69G5FAV
?id=01ARZ3NDEKTSV4RRFFQ69G5FAV
```

## Database Storage and Performance

### Storage Size

**Binary representation:**
- **128 bits = 16 bytes**
- Same as UUID

**String representation:**
- **26 characters**
- As UTF-8 string: 26 bytes minimum
- As MySQL `CHAR(26)` with `utf8mb4`: 72 bytes
- **Recommendation:** Store as binary (16 bytes) for optimal efficiency

**Storage comparison:**

| Format | Size | Efficiency |
|--------|------|------------|
| Binary (`BYTEA`/`BINARY(16)`) | 16 bytes | Optimal |
| String (`CHAR(26)`) | 26+ bytes | 1.6× larger |
| UUID string (`CHAR(36)`) | 36+ bytes | 2.25× larger |

### Index Performance

ULIDs provide **significant performance advantages** over random identifiers:

#### B-tree Index Benefits

**Sequential insertion pattern:**
- ✅ Dramatically reduces page splits vs UUID v4
- ✅ Minimizes write amplification
- ✅ Improves cache utilization
- ✅ Reduces I/O operations
- ✅ Prevents index fragmentation and bloat

**Recent benchmarks (PostgreSQL, 2024-2025):**

| ID Type | Ops/Second | Latency | Index Size |
|---------|-----------|---------|------------|
| ULID (bytea) | ~34,000 | 58 μs | Baseline |
| UUID v7 | ~34,000 | 58 μs | Similar |
| UUID v4 | ~25,000 | 85 μs | 85% larger |

**Key findings:**
- ULID performance comparable to or slightly better than UUID v7
- 33% faster than UUID v4
- Significantly more stable performance (lower variance)

### Lexicographic Sorting Benefits

**Chronological ordering:**
- ULIDs sort lexicographically in timestamp order
- No need for additional timestamp indexes
- Natural time-based ordering

**Query optimization benefits:**
```sql
-- Time-range queries are efficient
SELECT * FROM events
WHERE event_id >= '01ARZ3NDEK000000000000000'
  AND event_id <= '01ARZ3NDEKZZZZZZZZZZZZZZ';
```

**Advantages:**
- Efficient range queries on time-based data
- Simplified debugging (IDs reveal creation time)
- Better query planner optimization
- Natural partitioning by time ranges

### Impact on Page Splits and Fragmentation

**Dramatically reduced fragmentation compared to UUID v4:**

**UUID v4 problems:**
- Excessive page splits even before pages are full
- Random writes throughout B-tree structure
- Index bloat increases size on disk
- Temporally related rows spread across index

**ULID advantages:**
- Inserts at end of B-tree
- Minimizes splits to only last page
- Sequential writes optimize for append-heavy workloads
- Reduced index maintenance overhead

**Storage efficiency:**
- Less wasted space from partial pages
- More compact indexes
- Better compression ratios
- Lower storage costs for write-heavy applications

### Sequential Nature and Timestamp Ordering

**48-bit timestamp component:**
- Millisecond precision Unix timestamp
- Representation until year **10889 AD**
- High-order bits ensure chronological insertion
- Enables time-based partitioning strategies

**Performance characteristics:**
- New records naturally fall at end of B-tree
- Predictable insertion patterns
- Optimizes for sequential writes
- Reduces fragmentation over time

## Generation Approach

### ✅ Fully Decentralized

ULIDs can be generated in a completely decentralized manner with no coordination required.

**No centralized service needed:**
- Each system/node generates independently
- Only requires system clock access
- Cryptographically secure random number generator (CSPRNG)
- No network coordination overhead

### Structure: Timestamp + Randomness

**128 bits total:**

```
 01AN4Z07BY      79KA1307SR9X4MV3
|----------|    |----------------|
 Timestamp          Randomness
   48bits             80bits
```

**Timestamp component (48 bits):**
- Milliseconds since Unix epoch
- First 10 characters in encoded form
- Provides temporal ordering

**Randomness component (80 bits):**
- Cryptographically secure random value
- Remaining 16 characters
- Ensures uniqueness within same millisecond

**Binary encoding:**
- Most Significant Byte first (network byte order)
- Each component encoded as octets
- Total: 16 octets (bytes)

### Collision Resistance

**Extremely high collision resistance:**

- **1.21 × 10²⁴ unique IDs per millisecond** (2⁸⁰ possible values)
- Collision probability is practically zero
- Even in distributed systems, likelihood of collision is exceedingly low

**Example scale:**
- Would need to generate **trillions of IDs per millisecond** to see collisions
- Far exceeds any practical generation rate
- Safe for production at any realistic scale

### Monotonicity Guarantees

#### Standard Generation (Non-Monotonic)

**Default behavior:**
- Each ULID uses fresh random 80 bits
- Sortable by timestamp (millisecond precision)
- No guarantee of order within same millisecond

#### Monotonic Mode (Optional)

**Algorithm:**
1. If timestamp same as previous: increment previous random component
2. If timestamp advanced: generate fresh random component
3. If overflow (2⁸⁰ increments): wait for next millisecond or fail

**Benefits:**
- ✅ Guarantees strict ordering even at sub-millisecond generation
- ✅ Better collision resistance through sequential randomness
- ✅ Maintains sortability within same timestamp

**Trade-offs:**
- ⚠️ Leaks information about IDs generated within same millisecond
- ⚠️ Potential security concern: enables enumeration attacks
- ⚠️ Can overflow if > 2⁸⁰ IDs generated in one millisecond (theoretical only)

**Collision probability in monotonic mode:**
- Actually reduces collision risk
- Incrementing creates number groups less likely to collide
- Safe to use in production systems

## Comparison to UUID v7

Both ULID and UUID v7 solve similar problems with different approaches:

| Aspect | ULID | UUID v7 |
|--------|------|---------|
| **Size** | 16 bytes | 16 bytes |
| **Timestamp bits** | 48 | 48 |
| **Random bits** | 80 | 74 |
| **String format** | 26 chars (Base32) | 36 chars (hex + hyphens) |
| **Standardization** | Community spec | RFC 9562 (official) |
| **DB support** | Custom | Native (PostgreSQL 18+) |
| **Readability** | Better (Base32) | Standard (hex) |
| **Case sensitivity** | Insensitive | Insensitive |
| **Hyphens** | None | 4 hyphens |

**ULID advantages:**
- More compact string representation (26 vs 36)
- Slightly more random bits (80 vs 74)
- Better human readability (Crockford Base32)
- No hyphens (simpler to handle)

**UUID v7 advantages:**
- Official RFC standardization
- Growing native database support
- URN namespace compatibility (`urn:uuid:...`)
- Wider vendor tooling support

## 2024-2025 Landscape

**Current state:**
- UUID v7 (RFC 9562, 2024) now offers similar benefits with standardization
- ULID remains compelling for human readability and compact representation
- Both vastly superior to UUID v4 for database performance
- Choice often: standardization (v7) vs. readability (ULID)

**Industry adoption:**
- incident.io uses ULIDs for all identifiers
- Various startups prefer ULID for API design
- UUID v7 gaining traction in enterprise systems

## Use Cases

**ULIDs are excellent for:**
- ✅ Database primary keys (especially write-heavy workloads)
- ✅ Distributed systems requiring decentralized ID generation
- ✅ Applications needing URI-safe identifiers
- ✅ Systems benefiting from time-ordered IDs
- ✅ Scenarios requiring human-readable identifiers
- ✅ APIs where compact IDs are valued

**Consider alternatives when:**
- ⚠️ Strict RFC/ISO standardization required (use UUID v7)
- ⚠️ Native database support is priority (UUID v7 has better tooling)
- ⚠️ Absolute minimal storage (auto-increment or Snowflake)
- ⚠️ High-security scenarios sensitive to timing information leakage

## Implementation Examples

### PostgreSQL

```sql
-- Store as bytea for optimal performance
CREATE TABLE events (
    event_id BYTEA PRIMARY KEY DEFAULT ulid_generate(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    data JSONB
);

-- Custom function needed (no native support)
CREATE OR REPLACE FUNCTION ulid_generate()
RETURNS BYTEA AS $$
    -- Implementation using pgcrypto or external library
$$ LANGUAGE plpgsql;
```

### MySQL

```sql
-- Store as BINARY(16)
CREATE TABLE events (
    event_id BINARY(16) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data JSON
);

-- Generate in application layer
```

### Application-Level Generation

**Go example:**
```go
import "github.com/oklog/ulid/v2"

// Standard generation
id := ulid.Make()
fmt.Println(id.String()) // 01ARZ3NDEKTSV4RRFFQ69G5FAV

// Monotonic generation
entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
```

## Go Library Support

### ✅ oklog/ulid Library

The canonical Go library for ULIDs is [`github.com/oklog/ulid/v2`](https://github.com/oklog/ulid), which provides full ULID specification support with both standard and monotonic generation modes.

**Installation:**
```bash
go get github.com/oklog/ulid/v2
```

**Usage examples:**

```go
import (
    "crypto/rand"
    "github.com/oklog/ulid/v2"
)

// Simple generation with default entropy
id := ulid.Make()
fmt.Println(id.String()) // e.g., 01ARZ3NDEKTSV4RRFFQ69G5FAV

// Monotonic generation for strict ordering
entropy := ulid.Monotonic(rand.Reader, 0)
id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
```

## Summary

ULID represents an excellent choice for modern distributed systems:

**Key strengths:**
1. **Fully decentralized** - no coordination required
2. **URI-safe and compact** - 26 characters, no special chars
3. **Excellent database performance** - time-ordered, minimal fragmentation
4. **Human-readable** - Crockford Base32 alphabet
5. **High collision resistance** - 1.21 × 10²⁴ IDs per millisecond

**Key considerations:**
1. Not officially standardized (community spec)
2. Requires custom database functions (no native support)
3. Exposes creation timestamp (like UUID v7)
4. Slightly more complex than UUID v4 generation

**Bottom line:**
ULID is an excellent choice when you value compact, human-readable identifiers and don't require strict RFC compliance. For official standardization, UUID v7 offers similar performance with growing vendor support.
