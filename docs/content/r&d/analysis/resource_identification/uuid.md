---
title: "UUID Analysis"
type: docs
weight: 1
---

# UUID (Universally Unique Identifier)

## Overview

UUIDs are 128-bit identifiers standardized in RFC 9562 (May 2024), which obsoletes the previous RFC 4122. The latest specification introduces three new versions (v6, v7, v8) while maintaining backward compatibility with existing versions.

## URI Safety

### ✅ Fully URI-Safe

UUIDs are inherently safe for use in URIs without any encoding required.

**Standard format:**
```
550e8400-e29b-41d4-a716-446655440000
```

**Characteristics:**
- 36 characters: 32 hexadecimal digits + 4 hyphens
- Character set: `a-f`, `0-9`, `-`
- All characters are in RFC 3986 §2.3 unreserved set
- Case-insensitive (lowercase recommended per RFC 9562)

**Usage in URIs:**
```
/api/users/550e8400-e29b-41d4-a716-446655440000
?id=550e8400-e29b-41d4-a716-446655440000
urn:uuid:550e8400-e29b-41d4-a716-446655440000
```

**Alternative encodings:**
- Base64 URL-safe: 22 characters (optimization, not required)
- Base62: Similar length, avoids `+` and `/`
- These are for compactness, not safety

## Database Storage and Performance

### Storage Size

**Binary format:**
- **16 bytes (128 bits)** - canonical storage format
- Defined in RFC 9562

**String format:**
- 36 characters (`CHAR(36)`)
- Actual storage: 36-40 bytes depending on database encoding

**Storage comparison:**

| Format | Size | Overhead |
|--------|------|----------|
| Binary (`BINARY(16)`) | 16 bytes | baseline |
| String (`CHAR(36)`) | 36 bytes | 2.25× |
| String (`VARCHAR(36)`) | 38-40 bytes | ~2.5× |

### Database-Specific Implementations

**PostgreSQL:**
```sql
-- Use native UUID type (16 bytes internally)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

-- PostgreSQL 18+ supports UUIDv7
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7()
);
```

**Performance impact:**
- Native `UUID` type: 16 bytes
- Text storage: Tables 54% larger, indexes 85% larger

**MySQL:**
```sql
-- Use BINARY(16) with conversion functions
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID()))
);

-- Retrieve with conversion
SELECT BIN_TO_UUID(id) as id FROM users;
```

**SQL Server:**
```sql
CREATE TABLE users (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID()
);
```
- Note: `NEWSEQUENTIALID()` generates sequential UUIDs, not `NEWID()` which is random

### Index Performance

#### The UUID v4 Problem

**Random insertion issues:**

1. **Page splits:** New UUIDs insert at arbitrary positions in B-tree
2. **Fragmentation:** Index becomes scattered across non-contiguous pages
3. **Wasted space:** Page splits leave gaps throughout index
4. **Cache inefficiency:** Poor locality leads to more cache misses
5. **Write amplification:** More disk I/O per insert

**Measured impact:**
- Constant page splits during INSERT operations
- Index bloat (more pages for same data)
- 2-5× slower than sequential IDs
- Degraded SELECT performance

#### The UUID v7 Solution

**Sequential insertion benefits:**

1. **Append-only writes:** New entries go to end of index
2. **Minimal page splits:** Only last page splits when full
3. **Low fragmentation:** Index remains mostly contiguous
4. **Better caching:** Sequential access patterns
5. **Reduced I/O:** Fewer disk operations

**Measured improvements:**
- 2-5× faster insert performance vs v4
- 50% reduction in Write-Ahead Log (WAL) rate
- Fewer page splits comparable to auto-increment
- Better storage efficiency

### Binary vs String Storage

**Index size comparison (PostgreSQL):**

| Storage Type | Table Size | Index Size |
|--------------|------------|------------|
| Binary (UUID) | 100% (baseline) | 100% (baseline) |
| String (TEXT) | 154% | 185% |

**Why binary is faster:**
- Smaller indexes (fewer pages)
- Better cache utilization
- Faster CPU comparisons (128-bit integers)
- Reduced I/O (less data transfer)

## Generation Approach

### ✅ Fully Decentralized

One of UUID's core design goals is **decentralized generation without coordination**. Multiple systems can generate UUIDs independently without collision risk.

### UUID Version Comparison

#### UUID v1 - Time-based + MAC Address

**Structure:**
```
Timestamp (60 bits) + Clock Sequence (14 bits) + MAC Address (48 bits)
```

**Generation:**
- Timestamp: 100-nanosecond intervals since Oct 15, 1582
- Node ID: System's MAC address
- Clock sequence: Random value to prevent duplicates

**Pros:**
- Sequential (sorts chronologically)
- Very low collision risk
- Decentralized

**Cons:**
- ❌ **Privacy concern:** Leaks MAC address (physical location)
- ❌ Timestamp not in sortable byte order
- ❌ Modern systems avoid for security reasons

**Use case:** Legacy systems only (prefer v7)

#### UUID v4 - Random

**Structure:**
```
122 random bits + 6 version/variant bits
```

**Generation:**
- Entirely random (cryptographically secure RNG recommended)
- No coordination needed
- No sequential ordering

**Pros:**
- ✅ Maximum privacy (no identifying information)
- ✅ Simplest to generate
- ✅ Works offline
- ✅ Truly decentralized

**Cons:**
- ❌ **Poor database performance:** Random insertion causes fragmentation
- ❌ No time information
- ❌ Higher collision probability (still astronomically low)

**Collision probability:**
- 122 bits of entropy
- Need ~2.7 × 10¹⁸ UUIDs for 50% collision chance
- In practice: negligible

**Use cases:**
- Session IDs
- One-time tokens
- Non-database identifiers
- When pure randomness is desired

#### UUID v6 - Reordered Time-based

**Structure:**
```
Timestamp (60 bits, big-endian) + Clock Sequence + Node ID
```

**Generation:**
- Like v1 but timestamp bytes reordered for sorting
- Maintains MAC address (privacy concern)

**Pros:**
- Sortable (better than v1)
- Sequential insertion performance

**Cons:**
- ❌ Still leaks MAC address
- ❌ **Superseded by v7:** RFC 9562 recommends v7

**Use case:** None - v7 is better

#### UUID v7 - Time-ordered + Random ⭐ RECOMMENDED

**Structure:**
```
Unix Timestamp (48 bits, millisecond) + Random (74 bits)
```

**Generation:**
- Top 48 bits: Unix epoch milliseconds
- Bottom 74 bits: Random data
- No MAC address
- Monotonically increasing

**Pros:**
- ✅ **Excellent database performance:** Sequential inserts
- ✅ **Privacy-preserving:** No MAC address
- ✅ **Sortable:** Natural time ordering
- ✅ **Decentralized:** No coordination needed
- ✅ **Random component:** Prevents collisions from multiple nodes

**Performance measured:**
- 2-5× faster inserts than v4
- 50% reduction in WAL rate
- Minimal page splits
- Better cache locality

**Cons:**
- ⚠️ Exposes creation timestamp (usually acceptable)
- Slightly more complex than v4

**Use cases:**
- **Database primary keys** (optimal choice)
- Distributed systems
- Event IDs with time ordering
- Modern applications (default recommendation)

### Decentralization Requirements

**No central service required for any version:**

```go
// Example: Independent generation
// Node A
uuid1 := uuid.NewV7() // 0191e1a6-8b2c-7890-abcd-123456789abc

// Node B (same time)
uuid2 := uuid.NewV7() // 0191e1a6-8b2c-7890-xyz1-987654321def
```

**How v7 avoids collisions:**
1. **Time component:** Millisecond precision provides separation
2. **Random component:** 74 bits prevents same-millisecond collisions
3. **No coordination:** Each node generates independently

**Collision risk (UUID v7):**
- Within same millisecond: 2⁷⁴ unique values possible
- Even at 1 billion IDs per millisecond: negligible collision risk

## Version Selection Guide

```
┌─────────────────────────────────────────────────────┐
│ Which UUID Version?                                  │
├─────────────────────────────────────────────────────┤
│                                                      │
│  Database Primary Key? ──YES──> UUID v7            │
│         │                                            │
│         NO                                           │
│         │                                            │
│  Need time ordering? ──YES──> UUID v7              │
│         │                                            │
│         NO                                           │
│         │                                            │
│  Need pure randomness? ──YES──> UUID v4            │
│                                                      │
│  ❌ Avoid: v1 (privacy), v6 (superseded)           │
└─────────────────────────────────────────────────────┘
```

## Modern Recommendations (2024-2025)

**For new projects:**

1. **Default choice: UUID v7**
   - Best performance
   - Decentralized generation
   - No privacy concerns
   - Sortable

2. **Special cases: UUID v4**
   - Explicit randomness needed
   - Non-database contexts
   - Legacy compatibility

3. **Avoid: v1, v6**
   - v1: Privacy issues (MAC address)
   - v6: v7 is better in every way

## Recent Developments

### RFC 9562 (May 2024)
- Obsoletes RFC 4122
- Introduces v6, v7, v8
- Recommends v7 for database keys

### PostgreSQL 18 (2025)
- Native `gen_uuid_v7()` function
- Solves B-tree fragmentation
- Built-in time-ordered UUID generation

### Industry Adoption
- Buildkite: "Goodbye to sequential integers, hello UUIDv7"
- Cloud providers adding native support
- Database vendors implementing optimizations

## Summary

| Aspect | UUID v4 | UUID v7 |
|--------|---------|---------|
| **Storage** | 16 bytes binary | 16 bytes binary |
| **Generation** | Fully random | Time + random |
| **Decentralized** | ✅ Yes | ✅ Yes |
| **Coordination** | ❌ No | ❌ No |
| **URI safe** | ✅ Yes | ✅ Yes |
| **DB inserts** | ⚠️ Slow (random) | ✅ Fast (sequential) |
| **Fragmentation** | ⚠️ High | ✅ Low |
| **Page splits** | ⚠️ Frequent | ✅ Minimal |
| **Sortable** | ❌ No | ✅ Yes (by time) |
| **Privacy** | ✅ Maximum | ✅ Good |
| **Best for** | Tokens, session IDs | Database keys |

## Key Takeaways

1. **Always use binary storage** in databases (16 bytes vs 36-40 bytes)
2. **UUID v7 is the modern default** for database primary keys
3. **UUID v4 still useful** for session tokens and random IDs
4. **No coordination required** - all versions are fully decentralized
5. **URI-safe by design** - use directly in URLs without encoding
6. **RFC standardized** - wide vendor support and tooling available
