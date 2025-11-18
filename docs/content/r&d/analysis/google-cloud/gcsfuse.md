---
title: "Cloud Storage FUSE (gcsfuse)"
description: Analysis of Google Cloud Storage FUSE for mounting GCS buckets as local filesystems in network boot infrastructure
type: docs
weight: 3
---

## Overview

Cloud Storage FUSE (gcsfuse) is a FUSE-based filesystem adapter that allows Google Cloud Storage (GCS) buckets to be mounted and accessed as local filesystems on Linux systems. This enables applications to interact with object storage using standard filesystem operations (open, read, write, etc.) rather than requiring GCS-specific APIs.

**Project**: [GoogleCloudPlatform/gcsfuse](https://github.com/GoogleCloudPlatform/gcsfuse)
**License**: Apache 2.0
**Status**: Generally Available (GA)
**Latest Version**: v2.x (as of 2024)

## How gcsfuse Works

gcsfuse translates filesystem operations into GCS API calls:

1. **Mount Operation**: `gcsfuse bucket-name /mount/point` maps a GCS bucket to a local directory
2. **Directory Structure**: Interprets `/` in object names as directory separators
3. **File Operations**: Translates `read()`, `write()`, `open()`, etc. into GCS API requests
4. **Metadata**: Maintains file attributes (size, modification time) via GCS metadata
5. **Caching**: Optional stat, type, list, and file caching to reduce API calls

**Example**:
- GCS object: `gs://boot-assets/kernels/talos-v1.6.0.img`
- Mounted path: `/mnt/boot-assets/kernels/talos-v1.6.0.img`

## Relevance to Network Boot Infrastructure

In the context of [ADR-0005 Network Boot Infrastructure](../../adrs/0005-network-boot-infrastructure-gcp.md), gcsfuse offers a potential approach for serving boot assets from Cloud Storage without custom integration code.

### Potential Use Cases

1. **Boot Asset Storage**: Mount `gs://boot-assets/` to `/var/lib/boot-server/assets/`
2. **Configuration Sync**: Access boot profiles and machine mappings from GCS as local files
3. **Matchbox Integration**: Mount GCS bucket to `/var/lib/matchbox/` for assets/profiles/groups
4. **Simplified Development**: Eliminate custom Cloud Storage SDK integration in boot server code

### Architecture Pattern

```
┌─────────────────────────┐
│   Boot Server Process   │
│  (Cloud Run/Compute)    │
└───────────┬─────────────┘
            │ filesystem operations
            │ (read, open, stat)
            ▼
┌─────────────────────────┐
│   gcsfuse mount point   │
│   /var/lib/boot-assets  │
└───────────┬─────────────┘
            │ FUSE layer
            │ (translates to GCS API)
            ▼
┌─────────────────────────┐
│  Cloud Storage Bucket   │
│   gs://boot-assets/     │
└─────────────────────────┘
```

## Performance Characteristics

### Latency

- **Much higher latency than local filesystem**: Every operation requires GCS API call(s)
- **No default caching**: Without caching enabled, every read re-fetches from GCS
- **Network round-trip**: Minimum ~10-50ms latency per operation (depending on region)

### Throughput

**Single Large File**:
- Read: ~4.1 MiB/s (individual file), up to 63.3 MiB/s (archive files)
- Write: Comparable to `gsutil cp` for large files
- **With parallel downloads**: Up to 9x faster for single-threaded reads of large files

**Small Files**:
- Poor performance for random I/O on small files
- Bulk operations on many small files create significant bottlenecks
- `ls` on directories with thousands of objects can take minutes

**Concurrent Access**:
- Performance degrades significantly with parallel readers (8 instances: ~30 hours vs 16 minutes with local data)
- Not recommended for high-concurrency scenarios (web servers, NAS)

### Performance Improvements (Recent Features)

1. **Streaming Writes** (default): Upload data directly to GCS as written
   - Up to 40% faster for large sequential writes
   - Reduces local disk usage (no staging file)

2. **Parallel Downloads**: Download large files using multiple workers
   - Up to 9x faster model load times
   - Best for single-threaded reads of large files

3. **File Cache**: Cache file contents locally (Local SSD, Persistent Disk, or tmpfs)
   - Up to 2.3x faster training time (AI/ML workloads)
   - Up to 3.4x higher throughput
   - Requires explicit cache directory configuration

4. **Metadata Cache**: Cache stat, type, and list operations
   - Stat and type caches enabled by default
   - Configurable TTL (default: 60s, set `-1` for unlimited)

## Caching Configuration

gcsfuse provides four types of caching:

### 1. Stat Cache

Caches file attributes (size, modification time, existence).

```bash
# Enable with unlimited size and TTL
gcsfuse \
  --stat-cache-max-size-mb=-1 \
  --metadata-cache-ttl-secs=-1 \
  bucket-name /mount/point
```

**Use case**: Reduces API calls for repeated `stat()` operations (e.g., checking file existence).

### 2. Type Cache

Caches file vs directory type information.

```bash
gcsfuse \
  --type-cache-max-size-mb=-1 \
  --metadata-cache-ttl-secs=-1 \
  bucket-name /mount/point
```

**Use case**: Speeds up directory traversal and `ls` operations.

### 3. List Cache

Caches directory listing results.

```bash
gcsfuse \
  --max-conns-per-host=100 \
  --metadata-cache-ttl-secs=-1 \
  bucket-name /mount/point
```

**Use case**: Improves performance for applications that repeatedly list directory contents.

### 4. File Cache

Caches actual file contents locally.

```bash
gcsfuse \
  --file-cache-max-size-mb=-1 \
  --cache-dir=/mnt/local-ssd \
  --file-cache-cache-file-for-range-read=true \
  --file-cache-enable-parallel-downloads=true \
  bucket-name /mount/point
```

**Use case**: Essential for AI/ML training, repeated reads of large files.

**Recommended cache storage**:
- **Local SSD**: Fastest, but ephemeral (data lost on restart)
- **Persistent Disk**: Persistent but slower than Local SSD
- **tmpfs** (RAM disk): Fastest but limited by memory

### Production Configuration Example

```yaml
# config.yaml for gcsfuse
metadata-cache:
  ttl-secs: -1  # Never expire (use only if bucket is read-only or single-writer)
  stat-cache-max-size-mb: -1
  type-cache-max-size-mb: -1

file-cache:
  max-size-mb: -1  # Unlimited (limited by disk space)
  cache-file-for-range-read: true
  enable-parallel-downloads: true
  parallel-downloads-per-file: 16
  download-chunk-size-mb: 50

write:
  create-empty-file: false  # Streaming writes (default)

logging:
  severity: info
  format: json
```

```bash
gcsfuse --config-file=config.yaml boot-assets /mnt/boot-assets
```

## Limitations and Considerations

### Filesystem Semantics

gcsfuse provides **approximate POSIX semantics** but is not fully POSIX-compliant:

- **No atomic rename**: Rename operations are copy-then-delete (not atomic)
- **No hard links**: GCS doesn't support hard links
- **No file locking**: `flock()` is a no-op
- **Limited permissions**: GCS has simpler ACLs than POSIX permissions
- **No sparse files**: Writes always materialize full file content

### Performance Anti-Patterns

❌ **Avoid**:
- Serving web content or acting as NAS (concurrent connections)
- Random I/O on many small files (image datasets, text corpora)
- Reading during ML training loops (download first, then train)
- High-concurrency workloads (multiple parallel readers/writers)

✅ **Good for**:
- Sequential reads of large files (models, checkpoints, kernels)
- Infrequent writes of entire files
- Read-mostly workloads with caching enabled
- Single-writer scenarios

### Consistency Trade-offs

**With caching enabled**:
- Stale reads possible if cache TTL > 0 and external modifications occur
- Safe only for:
  - Read-only buckets
  - Single-writer, single-mount scenarios
  - Workloads tolerant of eventual consistency

**Without caching**:
- Strong consistency (every read fetches latest from GCS)
- Much slower performance

### Resource Requirements

- **Disk space**: File cache and streaming writes require local storage
  - File cache: Size of cached files (can be large for ML datasets)
  - Streaming writes: Temporary staging (proportional to concurrent writes)
- **Memory**: Metadata caches consume RAM
- **File handles**: Can exceed system limits with high concurrency
- **Network bandwidth**: All data transfers via GCS API

## Installation

### On Compute Engine (Container-Optimized OS)

```bash
# Install gcsfuse (Container-Optimized OS doesn't include package managers)
export GCSFUSE_VERSION=2.x.x
curl -L -O https://github.com/GoogleCloudPlatform/gcsfuse/releases/download/v${GCSFUSE_VERSION}/gcsfuse_${GCSFUSE_VERSION}_amd64.deb
sudo dpkg -i gcsfuse_${GCSFUSE_VERSION}_amd64.deb
```

### On Debian/Ubuntu

```bash
export GCSFUSE_REPO=gcsfuse-`lsb_release -c -s`
echo "deb https://packages.cloud.google.com/apt $GCSFUSE_REPO main" | sudo tee /etc/apt/sources.list.d/gcsfuse.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo apt-get update
sudo apt-get install gcsfuse
```

### In Docker/Cloud Run

```dockerfile
FROM ubuntu:22.04

# Install gcsfuse
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    lsb-release \
  && export GCSFUSE_REPO=gcsfuse-$(lsb_release -c -s) \
  && echo "deb https://packages.cloud.google.com/apt $GCSFUSE_REPO main" | tee /etc/apt/sources.list.d/gcsfuse.list \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
  && apt-get update \
  && apt-get install -y gcsfuse \
  && rm -rf /var/lib/apt/lists/*

# Create mount point
RUN mkdir -p /mnt/boot-assets

# Mount gcsfuse at startup
CMD gcsfuse --foreground boot-assets /mnt/boot-assets & \
    /usr/local/bin/boot-server
```

**Note**: Cloud Run **does not support FUSE filesystems** (requires privileged mode). gcsfuse only works on Compute Engine or GKE.

## Network Boot Infrastructure Evaluation

### Applicability to ADR-0005

Based on the analysis, gcsfuse is **not recommended** for the network boot infrastructure for the following reasons:

#### ❌ Cloud Run Incompatibility

- gcsfuse requires FUSE kernel module and privileged containers
- Cloud Run does not support FUSE or privileged mode
- ADR-0005 prefers Cloud Run deployment (HTTP-only boot enables serverless)
- **Impact**: Blocks Cloud Run deployment, forcing Compute Engine VM

#### ❌ Boot Latency Requirements

- Boot file requests target < 100ms latency (ADR-0005 confirmation criteria)
- gcsfuse adds 10-50ms+ latency per operation (network round-trips)
- Kernel/initrd downloads are latency-sensitive (network boot timeout)
- **Impact**: May exceed boot timeout thresholds

#### ❌ No Caching for Read-Write Workloads

- Boot server needs to write new assets and read existing ones
- File cache with unlimited TTL requires read-only or single-writer assumption
- Multiple boot server instances (autoscaling) violate single-writer constraint
- **Impact**: Either accept stale reads or disable caching (slow)

#### ❌ Small File Performance

- Machine mapping configs, boot scripts, profiles are small files (KB range)
- gcsfuse performs poorly on small, random I/O
- `ls` operations on directories with many profiles can be slow
- **Impact**: Slow boot configuration lookups

#### ✅ Alternative: Direct Cloud Storage SDK

Using `cloud.google.com/go/storage` SDK directly offers:

- **Lower latency**: Direct API calls without FUSE overhead
- **Cloud Run compatible**: No kernel module or privileged mode required
- **Better control**: Explicit caching, parallel downloads, streaming
- **Simpler deployment**: No mount management, no FUSE dependencies
- **Cost**: Similar API call costs to gcsfuse

**Recommended approach** (from ADR-0005):
```go
// Custom boot server using Cloud Storage SDK
storage := storage.NewClient(ctx)
bucket := storage.Bucket("boot-assets")

// Stream kernel to boot client
obj := bucket.Object("kernels/talos-v1.6.0.img")
reader, _ := obj.NewReader(ctx)
defer reader.Close()
io.Copy(w, reader)  // Stream to HTTP response
```

### When gcsfuse MIGHT Be Useful

Despite the above limitations, gcsfuse could be considered for:

1. **Matchbox on Compute Engine**:
   - Matchbox expects filesystem paths for assets (`/var/lib/matchbox/assets/`)
   - Compute Engine VM supports FUSE
   - Read-heavy workload (boot assets rarely change)
   - Could mount `gs://boot-assets/` to `/var/lib/matchbox/assets/` with file cache

2. **Development/Testing**:
   - Quick prototyping without writing Cloud Storage integration
   - Local development with production bucket access
   - Not recommended for production deployment

3. **Low-Throughput Scenarios**:
   - Home lab scale (< 10 boots/hour)
   - File cache enabled with Local SSD
   - Single Compute Engine VM (not autoscaled)

**Configuration for Matchbox + gcsfuse**:

```bash
#!/bin/bash
# Mount boot assets for Matchbox

BUCKET="boot-assets"
MOUNT_POINT="/var/lib/matchbox/assets"
CACHE_DIR="/mnt/disks/local-ssd/gcsfuse-cache"

mkdir -p "$MOUNT_POINT" "$CACHE_DIR"

gcsfuse \
  --stat-cache-max-size-mb=-1 \
  --type-cache-max-size-mb=-1 \
  --metadata-cache-ttl-secs=-1 \
  --file-cache-max-size-mb=-1 \
  --cache-dir="$CACHE_DIR" \
  --file-cache-cache-file-for-range-read=true \
  --file-cache-enable-parallel-downloads=true \
  --implicit-dirs \
  --foreground \
  "$BUCKET" "$MOUNT_POINT"
```

## Monitoring and Troubleshooting

### Metrics

gcsfuse exposes Prometheus metrics:

```bash
gcsfuse --prometheus --prometheus-port=9101 bucket /mnt/point
```

**Key metrics**:
- `gcs_read_count`: Number of GCS read operations
- `gcs_write_count`: Number of GCS write operations
- `gcs_read_bytes`: Bytes read from GCS
- `gcs_write_bytes`: Bytes written to GCS
- `fs_ops_count`: Filesystem operations by type (open, read, write, etc.)
- `fs_ops_error_count`: Filesystem operation errors

### Logging

```bash
# JSON logging for Cloud Logging integration
gcsfuse --log-format=json --log-file=/var/log/gcsfuse.log bucket /mnt/point
```

### Common Issues

**Issue**: `ls` on large directories takes minutes

**Solution**:
- Enable list caching with `--metadata-cache-ttl-secs=-1`
- Reduce directory depth (flatten object hierarchy)
- Consider prefix-based filtering instead of full listings

**Issue**: Stale reads after external bucket modifications

**Solution**:
- Reduce `--metadata-cache-ttl-secs` (default 60s)
- Disable caching entirely for strong consistency
- Use versioned object names (immutable assets)

**Issue**: `Transport endpoint is not connected` errors

**Solution**:
- Unmount cleanly before remounting: `fusermount -u /mnt/point`
- Check GCS bucket permissions (IAM roles)
- Verify network connectivity to `storage.googleapis.com`

**Issue**: High memory usage

**Solution**:
- Limit metadata cache sizes: `--stat-cache-max-size-mb=1024`
- Disable file cache if not needed
- Monitor with `--prometheus` metrics

## Comparison to Alternatives

### gcsfuse vs Direct Cloud Storage SDK

| Aspect | gcsfuse | Cloud Storage SDK |
|--------|---------|-------------------|
| **Latency** | Higher (FUSE overhead + GCS API) | Lower (direct GCS API) |
| **Cloud Run** | ❌ Not supported | ✅ Fully supported |
| **Development Effort** | Low (standard filesystem code) | Medium (SDK integration) |
| **Performance** | Slower (filesystem abstraction) | Faster (optimized for use case) |
| **Caching** | Built-in (stat, type, list, file) | Manual (application-level) |
| **Streaming** | Automatic | Explicit (`io.Copy`) |
| **Dependencies** | FUSE kernel module, privileged mode | None (pure Go library) |

**Recommendation**: Use Cloud Storage SDK directly for production network boot infrastructure.

### gcsfuse vs rsync/gsutil Sync

**Periodic sync pattern**:
```bash
# Sync bucket to local disk every 5 minutes
*/5 * * * * gsutil -m rsync -r gs://boot-assets /var/lib/boot-assets
```

| Aspect | gcsfuse | rsync/gsutil sync |
|--------|---------|-------------------|
| **Consistency** | Eventual (with caching) | Strong (within sync interval) |
| **Disk Usage** | Minimal (file cache optional) | Full copy of assets |
| **Latency** | GCS API per request | Local disk (fast) |
| **Sync Lag** | Real-time (no caching) or TTL | Sync interval (minutes) |
| **Deployment** | Requires FUSE | Simple cron job |

**Recommendation**: For read-heavy, infrequent-write workloads on Compute Engine, rsync/gsutil sync is simpler and faster than gcsfuse.

## Conclusion

Cloud Storage FUSE (gcsfuse) provides a convenient filesystem abstraction over GCS buckets, but **is not recommended for the network boot infrastructure** due to:

1. Cloud Run incompatibility (requires FUSE kernel module)
2. Added latency (FUSE overhead + network round-trips)
3. Poor performance for small files and concurrent access
4. Caching trade-offs (consistency vs performance)

**Recommended alternatives**:
- **Custom Boot Server**: Direct Cloud Storage SDK integration (`cloud.google.com/go/storage`)
- **Matchbox on Compute Engine**: rsync/gsutil sync to local disk
- **Cloud Run Deployment**: Direct SDK (no gcsfuse possible)

gcsfuse may be useful for **development/testing** or **Matchbox prototyping** on Compute Engine, but production deployments should use direct SDK integration or periodic sync for optimal performance and Cloud Run compatibility.

## References

- [Cloud Storage FUSE Documentation](https://cloud.google.com/storage/docs/cloud-storage-fuse/overview)
- [gcsfuse GitHub Repository](https://github.com/GoogleCloudPlatform/gcsfuse)
- [Performance Tuning Best Practices](https://cloud.google.com/storage/docs/gcsfuse-performance-and-best-practices)
- [gcsfuse Semantics](https://github.com/GoogleCloudPlatform/gcsfuse/blob/master/docs/semantics.md)
- [ADR-0005: Network Boot Infrastructure Implementation on Google Cloud](../../adrs/0005-network-boot-infrastructure-gcp.md)
