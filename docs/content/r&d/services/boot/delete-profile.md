---
title: "DELETE /api/v1/boot/{machine_id}/profile"
type: docs
description: "Delete a machine's boot profile and its associated blobs"
weight: 23
---

Delete a machine's boot profile and its associated blobs.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Client as Admin Client
    participant Boot as Boot Service
    participant Storage as Cloud Storage
    participant DB as Firestore

    Client->>Boot: DELETE /api/v1/boot/{machine_id}/profile
    Boot->>DB: Get kernel_id and initrd_id
    DB-->>Boot: Blob IDs
    Boot->>Storage: DELETE gs://bucket/blobs/{kernel_id}
    Boot->>Storage: DELETE gs://bucket/blobs/{initrd_id}
    Boot->>DB: Delete boot profile
    Boot-->>Client: 204 No Content
```

## Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `machine_id` | string | Yes | Machine identifier (UUIDv7 format) |

**Example Request:**

```http
DELETE /api/v1/boot/018c7dbd-c000-7000-8000-fedcba987654/profile HTTP/1.1
Host: boot.example.com
```

## Response

**Response (204 No Content):**

Empty response body.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine not found or has no boot profile |
