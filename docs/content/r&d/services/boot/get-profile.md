---
title: "GET /api/v1/boot/{machine_id}/profile"
type: docs
description: "Retrieve the active boot profile for a specific machine"
weight: 21
---

Retrieve the active boot profile for a specific machine.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Client as Admin Client
    participant Boot as Boot Service
    participant DB as Firestore

    Client->>Boot: GET /api/v1/boot/{machine_id}/profile
    Boot->>DB: Query active boot profile for machine
    DB-->>Boot: Boot profile
    Boot-->>Client: 200 OK (boot profile)
```

## Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `machine_id` | string | Yes | Machine identifier (UUIDv7 format) |

**Example Request:**

```http
GET /api/v1/boot/018c7dbd-c000-7000-8000-fedcba987654/profile HTTP/1.1
Host: boot.example.com
```

## Response

**Response (200 OK):**

```json
{
  "id": "018c7dbd-a000-7000-8000-abcdef123456",
  "machine_id": "018c7dbd-c000-7000-8000-fedcba987654",
  "kernel": {
    "id": "018c7dbd-b100-7000-8000-123456789abc",
    "args": ["console=tty0", "console=ttyS0", "ip=dhcp"]
  },
  "initrd": {
    "id": "018c7dbd-b200-7000-8000-987654321fed"
  }
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine not found or has no boot profile |
