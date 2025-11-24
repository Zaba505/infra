---
title: "PUT /api/v1/machines/{id}"
type: docs
description: "Update a machine's hardware profile"
weight: 23
---

Update a machine's hardware profile.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Client as Admin Client
    participant API as Machine Service
    participant DB as Firestore

    Client->>API: PUT /api/v1/machines/{id}
    API->>DB: Update machine profile
    DB-->>API: Machine updated
    API-->>Client: 200 OK (updated profile)
```

## Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Machine identifier (UUIDv7 format) |

**Request Body:**

Full machine profile (same structure as POST /api/v1/machines):

```json
{
  "cpus": [
    {
      "manufacturer": "Intel",
      "clock_frequency": 2400000000,
      "cores": 8
    }
  ],
  "memory_modules": [
    {
      "size": 17179869184
    },
    {
      "size": 17179869184
    }
  ],
  "accelerators": [],
  "nics": [
    {
      "mac": "52:54:00:12:34:56"
    }
  ],
  "drives": [
    {
      "capacity": 500107862016
    }
  ]
}
```

## Response

**Response (200 OK):**

Full machine profile with updated fields:

```json
{
  "id": "018c7dbd-c000-7000-8000-fedcba987654",
  "cpus": [
    {
      "manufacturer": "Intel",
      "clock_frequency": 2400000000,
      "cores": 8
    }
  ],
  "memory_modules": [
    {
      "size": 17179869184
    },
    {
      "size": 17179869184
    }
  ],
  "accelerators": [],
  "nics": [
    {
      "mac": "52:54:00:12:34:56"
    }
  ],
  "drives": [
    {
      "capacity": 500107862016
    }
  ]
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine with specified ID not found |
| 400 Bad Request | Invalid request body |
