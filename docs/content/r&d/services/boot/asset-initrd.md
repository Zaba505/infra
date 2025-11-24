---
title: "GET /asset/{boot_profile_id}/initrd"
type: docs
description: "Streams initial ramdisk images from Cloud Storage for the boot process"
weight: 12
---

Streams initial ramdisk (initrd) images from Cloud Storage for the boot process. This endpoint is accessed by bare metal servers during UEFI HTTP boot through the WireGuard VPN tunnel.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Client as Bare Metal Server
    participant Boot as Boot Service
    participant Storage as Cloud Storage
    participant DB as Firestore

    Client->>Boot: GET /asset/018c7dbd-a1b2-7000-8000-987654321def/initrd
    Boot->>Boot: Validate UUIDv7 format
    Boot->>DB: Query boot profile by ID
    DB-->>Boot: Boot profile (initrd_id)
    Boot->>Storage: GET gs://bucket/blobs/{initrd_id}
    Storage-->>Boot: Initrd data stream
    Boot-->>Client: 200 OK (initrd stream)
```

## Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `boot_profile_id` | string (UUIDv7) | Yes | Boot profile identifier (UUIDv7 format: `018c7dbd-a1b2-7000-8000-987654321def`) |

**Request Example:**

```http
GET /asset/018c7dbd-a1b2-7000-8000-987654321def/initrd HTTP/1.1
Host: boot.internal
```

## Response

**Response Example (200 OK):**

Binary initrd image streamed from Cloud Storage.

**Response Headers:**

- `Content-Type: application/octet-stream`
- `Content-Length: 52428800` (actual initrd size in bytes)
- `Cache-Control: public, max-age=3600`
- `ETag: "def456..."`

**Error Responses:**

All error responses follow RFC 7807 Problem Details format (see [ADR-0007](../../adrs/0007-standard-api-error-response/)) with `Content-Type: application/problem+json`.

**404 Not Found** - Initrd image not found:

```json
{
  "type": "https://api.example.com/errors/initrd-not-found",
  "title": "Initrd Not Found",
  "status": 404,
  "detail": "Initrd image not found for boot profile 018c7dbd-a1b2-7000-8000-987654321def",
  "instance": "/asset/018c7dbd-a1b2-7000-8000-987654321def/initrd",
  "boot_profile_id": "018c7dbd-a1b2-7000-8000-987654321def"
}
```

**500 Internal Server Error** - Cloud Storage error:

```json
{
  "type": "https://api.example.com/errors/storage-error",
  "title": "Storage Error",
  "status": 500,
  "detail": "Failed to retrieve initrd from storage due to an internal error",
  "instance": "/asset/018c7dbd-a1b2-7000-8000-987654321def/initrd"
}
```

## Performance Characteristics

- **Streaming**: File is streamed directly from Cloud Storage (no buffering in memory)
- **Target Latency**: < 100ms to first byte
- **Typical Size**: 50-150 MB for Linux initrd images

## Security Considerations

### VPN Source IP Validation

All boot endpoints validate that requests originate from the WireGuard VPN subnet:

- **Allowed CIDR**: `10.x.x.0/24` (WireGuard VPN network)
- **Validation**: Performed at Cloud Run ingress or application layer
- **Rejection**: Requests from outside VPN return `403 Forbidden`

### Rate Limiting

To prevent abuse, asset download endpoints are rate-limited:

- **Asset Downloads**: 5 concurrent downloads per MAC address

### Asset Integrity

Boot assets are validated for integrity:

- **Checksums**: SHA-256 checksums stored in Firestore
- **Verification**: Computed on upload, verified on download (optional)
- **ETag Headers**: Enable client-side caching and integrity checks

## Observability

All boot endpoint requests are instrumented with OpenTelemetry following HTTP semantic conventions:

- **Metrics**: OpenTelemetry HTTP server metrics
  - `http.server.request.duration` - Request duration histogram
  - `http.server.response.body.size` - Response body size (tracks bytes transferred)
- **Traces**: End-to-end tracing from request to Cloud Storage retrieval
  - HTTP server span captures request details (method, route, status code)
  - Child spans for database queries and Cloud Storage operations
- **Logs**: Structured logs with boot profile ID, initrd ID, response status
