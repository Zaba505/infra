---
title: "Machine Service"
type: docs
description: "Service for managing machine hardware profiles"
weight: 20
---

The Machine Service is a REST API that manages machine hardware profiles for the network boot infrastructure. It stores machine specifications (CPUs, memory, NICs, drives, accelerators) in Firestore and is queried by the Boot Service during boot operations and by administrators for configuration management.

## Architecture

The service is responsible for:

- **Machine Profile Management**: Creating, listing, retrieving, updating, and deleting machine hardware profiles
- **Hardware Specification Storage**: Storing detailed hardware specifications in Firestore
- **Machine Lookup**: Providing machine profile queries by ID or NIC MAC address

## Components

- **Firestore**: Stores machine hardware profiles
- **REST API**: HTTP endpoints for machine profile management

## Clients

The service is consumed by:

1. **Boot Service**: Queries machine profiles by MAC address during boot operations
2. **Admin Tools**: CLI or web interfaces for managing machine inventory
3. **Monitoring Systems**: Hardware inventory and asset management tools

## Deployment

- **Platform**: GCP Cloud Run
- **Scaling**: Automatic scaling based on request load
- **Availability**: Min instances = 1 for low-latency responses
- **Region**: Same region as Boot Service for minimal latency

## API Endpoints

### Machine Management

- [POST /api/v1/machines](./post-machines/) - Register a new machine with hardware specifications
- [GET /api/v1/machines](./get-machines/) - List all registered machines
- [GET /api/v1/machines/{id}](./get-machine/) - Retrieve a specific machine by ID
- [PUT /api/v1/machines/{id}](./put-machine/) - Update a machine's hardware profile
- [DELETE /api/v1/machines/{id}](./delete-machine/) - Delete a machine registration

## Rate Limiting

Admin API endpoints are rate-limited to prevent abuse:

- **Per User/Service Account**: 100 requests/minute
- **Per IP Address**: 300 requests/minute
- **Global**: 1000 requests/minute

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1700000000
```

When rate limit is exceeded, API returns `429 Too Many Requests`:

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again in 30 seconds.",
    "details": {
      "retry_after": 30
    }
  }
}
```

## Versioning

The Admin API uses URL versioning (`/api/v1/`):

- **Current Version**: v1
- **Deprecation Policy**: Minimum 6 months notice before version deprecation
- **Version Header**: `X-API-Version: v1` included in all responses
