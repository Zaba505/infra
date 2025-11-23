---
title: "Boot Service"
type: docs
description: "UEFI HTTP boot endpoints and boot profile management"
weight: 10
---

The Boot Service is a custom Go microservice that provides UEFI HTTP boot endpoints for bare metal servers and manages boot profiles. It serves boot scripts, streams kernel/initrd assets, and handles boot profile administration (kernel/initrd upload, storage, and lifecycle management).

## Architecture Overview

The Boot Service is deployed on GCP Cloud Run and accessed through a WireGuard VPN tunnel from bare metal servers. It integrates with:

- **Machine Service**: Retrieves machine hardware profiles by MAC address
- **Cloud Storage**: Stores and retrieves kernel/initrd blobs
- **Firestore**: Stores boot profile metadata
- **Cloud Monitoring**: OpenTelemetry observability with distributed tracing

## Related Documentation

- [Machine Service](../machine-mgmt/) - Machine hardware profile management
- [ADR-0005: Network Boot Infrastructure Implementation on Google Cloud](../../adrs/0005-network-boot-infrastructure-gcp/) - Architecture decision and design rationale
- [ADR-0002: Network Boot Architecture](../../adrs/0002-network-boot-architecture/) - Overall network boot strategy

## API Categories

1. **[UEFI HTTP Boot Endpoints](./uefi-boot-endpoints/)** - Accessed by bare metal servers during boot process (via WireGuard VPN)
2. **[Admin API](./admin-api/)** - Boot profile management endpoints for administrators
3. **[Health Check Endpoints](./health-checks/)** - Standard Cloud Run health endpoints

## Security Model

### VPN-Based Access Control

Since HP DL360 Gen 9 servers do not support client-side TLS certificates for UEFI HTTP boot, all boot traffic is secured via WireGuard VPN:

- **Boot Endpoints**: Only accessible through WireGuard tunnel (source IP validation)
- **Transport Security**: WireGuard provides mutual authentication and encryption

### Authentication Methods

- **UEFI Boot Endpoints**: VPN source IP validation (bare metal servers)
- **Health Checks**: Unauthenticated (used by Cloud Run for liveness/startup probes)

## Common Patterns

### Error Responses

All API endpoints follow a consistent error response format:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Machine with MAC address aa:bb:cc:dd:ee:ff not found",
    "details": {
      "mac_address": "aa:bb:cc:dd:ee:ff"
    }
  }
}
```

### Standard HTTP Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `204 No Content` - Successful deletion
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error

### Content Types

- `application/json` - JSON responses (admin API)
- `text/plain` - iPXE boot scripts
- `application/octet-stream` - Binary boot assets (kernel, initrd)
- `text/cloud-config` - Cloud-init configuration files
