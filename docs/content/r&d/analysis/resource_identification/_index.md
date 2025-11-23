---
title: "Resource Identifiers Analysis"
weight: 50
---

# Resource Identifiers: UUID, ULID, and Snowflake ID

This analysis compares three popular distributed identifier strategies for use in modern systems: UUID (particularly v4 and v7), ULID, and Snowflake ID. The comparison focuses on three critical aspects:

1. **URI Safety**: Can they be used directly in URLs without encoding?
2. **Database Performance**: Storage size and index performance implications
3. **Generation Model**: Centralized vs decentralized generation

## Quick Comparison Table

| Aspect | UUID v4 | UUID v7 | ULID | Snowflake ID |
|--------|---------|---------|------|--------------|
| **Size (binary)** | 16 bytes | 16 bytes | 16 bytes | 8 bytes |
| **Size (string)** | 36 chars | 36 chars | 26 chars | 18-19 chars |
| **URI Safe** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Time-Ordered** | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes |
| **Decentralized** | ✅ Yes | ✅ Yes | ✅ Yes | ⚠️ Mostly |
| **Index Performance** | ⚠️ Poor | ✅ Good | ✅ Good | ✅ Excellent |
| **Standardized** | ✅ RFC 9562 | ✅ RFC 9562 | ❌ Spec only | ❌ Pattern |
| **Database Support** | ✅ Native | 🆕 Limited | ❌ Custom | ❌ Custom |

## Detailed Analyses

- [UUID Analysis](uuid/)
- [ULID Analysis](ulid/)
- [Snowflake ID Analysis](snowflake/)

## Decision Guide

### Use UUID v7 when:
- ✅ RFC standardization is important
- ✅ Native database support is desired (PostgreSQL 18+)
- ✅ You need URN compatibility (`urn:uuid:...`)
- ✅ You want official vendor support and tooling

### Use ULID when:
- ✅ Human readability is valued (Crockford Base32)
- ✅ You prefer compact string representation (26 vs 36 chars)
- ✅ Lexicographic sorting is important
- ✅ You want case-insensitive identifiers

### Use Snowflake ID when:
- ✅ Storage efficiency is critical (8 vs 16 bytes)
- ✅ Numeric IDs are required
- ✅ You can manage worker ID allocation
- ✅ Maximum database performance is needed
- ✅ You have a fixed number of generator nodes (<1024)

### Use UUID v4 when:
- ✅ Maximum randomness is required
- ✅ Session tokens or one-time IDs
- ✅ Time ordering is unimportant
- ❌ Avoid for database primary keys

## Modern Recommendations (2024-2025)

**For new projects with database primary keys:**

1. **First choice: UUID v7 or ULID**
   - Both offer excellent performance with time-ordering
   - UUID v7: Better standardization and tooling
   - ULID: Better readability and compact format

2. **Storage-constrained systems: Snowflake ID**
   - 50% smaller than UUID/ULID
   - Best database performance
   - Requires worker ID coordination

3. **Legacy compatibility: UUID v4**
   - Only if required by existing systems
   - Significant performance penalty for databases

**Avoid entirely:**
- UUID v1: Privacy concerns (leaks MAC address)
- UUID v6: Superseded by v7
- Auto-increment integers: Not distributed-system safe

## Key Insights

### URI Safety
All three identifier types are completely safe for direct use in URIs without percent-encoding:
- **UUID**: Hexadecimal + hyphens (RFC 3986 unreserved characters)
- **ULID**: Crockford Base32 alphabet (no confusing characters)
- **Snowflake**: Decimal integers (0-9 only)

### Database Performance

The critical factor is **sequential vs random insertion**:

**Random insertion (UUID v4):**
- Causes B-tree page splits throughout the index
- Results in fragmentation and bloat
- Poor cache utilization
- 2-5× slower than sequential

**Sequential insertion (UUID v7, ULID, Snowflake):**
- Appends to end of B-tree
- Minimal page splits
- Better cache locality
- Comparable to auto-increment integers

**Storage comparison:**
```
Snowflake ID:  8 bytes  (baseline)
UUID/ULID:    16 bytes  (2× larger)
UUID string:  36 bytes  (4.5× larger)
```

### Generation Models

**Fully decentralized (no coordination):**
- UUID v4: Pure randomness
- UUID v7: Timestamp + randomness
- ULID: Timestamp + randomness

**Minimal coordination (worker ID only):**
- Snowflake ID: Requires unique worker ID per generator
  - One-time configuration
  - Supports 1,024 workers (10 bits)
  - Challenge: Auto-scaling environments

## Performance Benchmarks

From recent studies (2024-2025):

**PostgreSQL INSERT operations:**
- Snowflake ID: ~34,000 ops/sec
- UUID v7: ~34,000 ops/sec (33% faster than v4)
- ULID: ~34,000 ops/sec (comparable to v7)
- UUID v4: ~25,000 ops/sec

**Index fragmentation (PostgreSQL):**
- UUID v4: 85% larger indexes, 54% larger tables
- UUID v7/ULID: Minimal fragmentation
- Snowflake: Minimal fragmentation, 50% smaller indexes

**Write-Ahead Log (WAL) generation:**
- UUID v7: 50% reduction vs UUID v4
- Sequential IDs reduce database write amplification

## Collision Resistance

All three approaches provide exceptional collision resistance:

**UUID v4:**
- 122 bits of randomness
- Need ~2.7 × 10¹⁸ IDs for 50% collision probability

**UUID v7:**
- 48-bit timestamp + 74-bit random
- Negligible collision risk even at millions per millisecond

**ULID:**
- 48-bit timestamp + 80-bit random
- 1.21 × 10²⁴ unique IDs per millisecond possible

**Snowflake ID:**
- Mathematical uniqueness guarantee
- No collisions possible if worker IDs are unique
- 4,096 IDs per millisecond per worker

## Implementation Considerations

### PostgreSQL
```sql
-- UUID v7 (PostgreSQL 18+)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7()
);

-- UUID v7 (PostgreSQL <18 with extension)
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- Use custom function or library

-- ULID (custom type or text/bytea)
CREATE TABLE users (
    id BYTEA PRIMARY KEY DEFAULT ulid_generate()
);

-- Snowflake (bigint)
CREATE TABLE users (
    id BIGINT PRIMARY KEY DEFAULT snowflake_generate()
);
```

### MySQL
```sql
-- UUID (binary storage recommended)
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID()))
);

-- Snowflake
CREATE TABLE users (
    id BIGINT PRIMARY KEY
);
```

### Application-Level Generation

**Advantages:**
- No database dependency
- Works with any database
- Consistent across different storage systems
- Better control over implementation

**Disadvantages:**
- Requires library/code maintenance
- Clock synchronization considerations
- Worker ID management (Snowflake only)

## Migration Strategies

### Moving from Auto-Increment

**Considerations:**
- Foreign key updates required
- Index rebuilds may be needed
- Application code changes
- Dual-write period during migration

**Recommended approach:**
1. Add new ID column alongside existing
2. Generate IDs for existing rows
3. Update foreign keys progressively
4. Migrate application code
5. Remove old ID column

### Moving from UUID v4 to v7/ULID

**Benefits:**
- Same storage size (16 bytes)
- Can keep existing IDs
- Only new records use v7/ULID
- Gradual performance improvement

## Security Considerations

### Information Leakage

**UUID v4:**
- ✅ Reveals nothing (pure randomness)

**UUID v7 / ULID:**
- ⚠️ Reveals creation timestamp (usually acceptable)
- ⚠️ May reveal approximate volume (via sequence patterns)

**Snowflake ID:**
- ⚠️ Reveals exact creation time (41-bit timestamp)
- ⚠️ Reveals which worker generated it
- ⚠️ Reveals sequence count within millisecond

### Enumeration Attacks

**Random IDs (UUID v4):**
- ✅ Resistant to enumeration
- Guessing next ID is infeasible

**Sequential IDs (v7, ULID, Snowflake):**
- ⚠️ Predictable patterns
- Can estimate next ID value
- **Mitigation**: Use authentication/authorization, don't rely on ID secrecy

### Recommendation
Never rely on ID unpredictability as a security mechanism. Always use proper authentication and authorization regardless of ID type.

## Conclusion

The landscape of distributed identifiers has evolved significantly:

**2010-2020:** UUID v4 was the default distributed identifier despite performance issues

**2020-2024:** Community alternatives (ULID, Snowflake) gained popularity for performance

**2024+:** UUID v7 (RFC 9562) provides standardized time-ordered IDs with vendor support

For most modern applications, **UUID v7 or ULID** represent the optimal balance of performance, standardization, and operational simplicity. **Snowflake IDs** remain compelling for storage-constrained systems where the 8-byte size and numeric format provide tangible benefits.

The days of suffering UUID v4's random insertion penalty for database primary keys are over—time-ordered identifiers are now the recommended default.
