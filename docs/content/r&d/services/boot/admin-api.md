---
title: "HTTP REST Admin API"
type: docs
description: "Management API for boot images, machines, and profiles"
weight: 20
---

The Admin API provides HTTP REST endpoints for managing boot images, machine mappings, and boot profiles. All endpoints require GCP IAM authentication.

## Authentication

All admin API endpoints require authentication using one of the following methods:

- **Service Account Token**: `Authorization: Bearer <service-account-token>`
- **User OAuth Token**: `Authorization: Bearer <user-oauth-token>`
- **IAM Authentication**: Validated against GCP IAM policies

**Required IAM Permissions:**

- `bootserver.images.create` - Upload boot images
- `bootserver.images.read` - List and retrieve boot images
- `bootserver.images.delete` - Delete boot images
- `bootserver.machines.create` - Register machines
- `bootserver.machines.read` - List and retrieve machine configurations
- `bootserver.machines.update` - Update machine mappings
- `bootserver.machines.delete` - Delete machine registrations
- `bootserver.profiles.create` - Create boot profiles
- `bootserver.profiles.read` - List and retrieve profiles
- `bootserver.profiles.update` - Update boot profiles
- `bootserver.profiles.delete` - Delete boot profiles

---

## Boot Image Management

### `POST /api/v1/images`

Upload a new boot image (kernel, initrd, and metadata).

**Request Body:**

```json
{
  "id": "ubuntu-2204",
  "name": "Ubuntu 22.04 LTS Server",
  "version": "22.04.3",
  "kernel": {
    "url": "gs://boot-images/kernels/ubuntu-22.04.3-kernel.img",
    "sha256": "a1b2c3d4e5f6..."
  },
  "initrd": {
    "url": "gs://boot-images/initrd/ubuntu-22.04.3-initrd.img",
    "sha256": "f6e5d4c3b2a1..."
  },
  "metadata": {
    "os": "ubuntu",
    "os_version": "22.04.3",
    "architecture": "x86_64",
    "tags": ["lts", "server"]
  }
}
```

**Request Headers:**

- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**Response (201 Created):**

```json
{
  "id": "ubuntu-2204",
  "name": "Ubuntu 22.04 LTS Server",
  "version": "22.04.3",
  "kernel": {
    "url": "gs://boot-images/kernels/ubuntu-22.04.3-kernel.img",
    "sha256": "a1b2c3d4e5f6...",
    "size_bytes": 8388608
  },
  "initrd": {
    "url": "gs://boot-images/initrd/ubuntu-22.04.3-initrd.img",
    "sha256": "f6e5d4c3b2a1...",
    "size_bytes": 52428800
  },
  "metadata": {
    "os": "ubuntu",
    "os_version": "22.04.3",
    "architecture": "x86_64",
    "tags": ["lts", "server"]
  },
  "created_at": "2025-11-19T06:00:00Z",
  "created_by": "admin@example.com"
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | Invalid request body or missing required fields |
| 401 Unauthorized | Missing or invalid authentication |
| 403 Forbidden | Insufficient permissions |
| 409 Conflict | Image with the same ID already exists |
| 422 Unprocessable Entity | Validation error (invalid SHA256, unreachable URL) |

---

### `GET /api/v1/images`

List all boot images.

**Query Parameters:**

| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `page` | integer | No | Page number (1-indexed) | 1 |
| `per_page` | integer | No | Results per page (1-100) | 20 |
| `os` | string | No | Filter by operating system | - |
| `architecture` | string | No | Filter by architecture | - |
| `tags` | string | No | Filter by tags (comma-separated) | - |

**Request Example:**

```http
GET /api/v1/images?os=ubuntu&architecture=x86_64&page=1&per_page=20 HTTP/1.1
Host: boot.example.com
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
{
  "images": [
    {
      "id": "ubuntu-2204",
      "name": "Ubuntu 22.04 LTS Server",
      "version": "22.04.3",
      "kernel": {
        "url": "gs://boot-images/kernels/ubuntu-22.04.3-kernel.img",
        "sha256": "a1b2c3d4e5f6...",
        "size_bytes": 8388608
      },
      "initrd": {
        "url": "gs://boot-images/initrd/ubuntu-22.04.3-initrd.img",
        "sha256": "f6e5d4c3b2a1...",
        "size_bytes": 52428800
      },
      "metadata": {
        "os": "ubuntu",
        "os_version": "22.04.3",
        "architecture": "x86_64",
        "tags": ["lts", "server"]
      },
      "created_at": "2025-11-19T06:00:00Z",
      "created_by": "admin@example.com"
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "per_page": 20,
    "total_pages": 1
  }
}
```

---

### `GET /api/v1/images/{id}`

Retrieve a specific boot image by ID.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot image identifier |

**Request Example:**

```http
GET /api/v1/images/ubuntu-2204 HTTP/1.1
Host: boot.example.com
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
{
  "id": "ubuntu-2204",
  "name": "Ubuntu 22.04 LTS Server",
  "version": "22.04.3",
  "kernel": {
    "url": "gs://boot-images/kernels/ubuntu-22.04.3-kernel.img",
    "sha256": "a1b2c3d4e5f6...",
    "size_bytes": 8388608
  },
  "initrd": {
    "url": "gs://boot-images/initrd/ubuntu-22.04.3-initrd.img",
    "sha256": "f6e5d4c3b2a1...",
    "size_bytes": 52428800
  },
  "metadata": {
    "os": "ubuntu",
    "os_version": "22.04.3",
    "architecture": "x86_64",
    "tags": ["lts", "server"]
  },
  "created_at": "2025-11-19T06:00:00Z",
  "created_by": "admin@example.com"
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Image with specified ID not found |

---

### `DELETE /api/v1/images/{id}`

Delete a boot image.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot image identifier |

**Request Example:**

```http
DELETE /api/v1/images/ubuntu-2204 HTTP/1.1
Host: boot.example.com
Authorization: Bearer <token>
```

**Response (204 No Content):**

Empty response body.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Image with specified ID not found |
| 409 Conflict | Image is currently in use by one or more machines |

---

## Machine Mapping Management

### `POST /api/v1/machines`

Register a new machine and map it to a boot profile.

**Request Body:**

```json
{
  "mac_address": "52:54:00:12:34:56",
  "hostname": "node-01",
  "profile_id": "ubuntu-server-base",
  "metadata": {
    "datacenter": "homelab",
    "rack": "A1",
    "role": "compute"
  },
  "network": {
    "ip_address": "10.0.1.10",
    "netmask": "255.255.255.0",
    "gateway": "10.0.1.1",
    "dns_servers": ["10.0.1.1", "8.8.8.8"]
  }
}
```

**Response (201 Created):**

```json
{
  "mac_address": "52:54:00:12:34:56",
  "hostname": "node-01",
  "profile_id": "ubuntu-server-base",
  "metadata": {
    "datacenter": "homelab",
    "rack": "A1",
    "role": "compute"
  },
  "network": {
    "ip_address": "10.0.1.10",
    "netmask": "255.255.255.0",
    "gateway": "10.0.1.1",
    "dns_servers": ["10.0.1.1", "8.8.8.8"]
  },
  "created_at": "2025-11-19T06:00:00Z",
  "updated_at": "2025-11-19T06:00:00Z"
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | Invalid MAC address format or missing required fields |
| 409 Conflict | Machine with the same MAC address already exists |
| 422 Unprocessable Entity | Referenced profile_id does not exist |

---

### `GET /api/v1/machines`

List all registered machines.

**Query Parameters:**

| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `page` | integer | No | Page number (1-indexed) | 1 |
| `per_page` | integer | No | Results per page (1-100) | 20 |
| `profile_id` | string | No | Filter by boot profile | - |
| `role` | string | No | Filter by machine role | - |

**Response (200 OK):**

```json
{
  "machines": [
    {
      "mac_address": "52:54:00:12:34:56",
      "hostname": "node-01",
      "profile_id": "ubuntu-server-base",
      "metadata": {
        "datacenter": "homelab",
        "rack": "A1",
        "role": "compute"
      },
      "network": {
        "ip_address": "10.0.1.10",
        "netmask": "255.255.255.0",
        "gateway": "10.0.1.1",
        "dns_servers": ["10.0.1.1", "8.8.8.8"]
      },
      "created_at": "2025-11-19T06:00:00Z",
      "updated_at": "2025-11-19T06:00:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "per_page": 20,
    "total_pages": 1
  }
}
```

---

### `GET /api/v1/machines/{mac}`

Retrieve a specific machine configuration by MAC address.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mac` | string | Yes | MAC address (format: `aa:bb:cc:dd:ee:ff`) |

**Response (200 OK):**

Same as POST response.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine with specified MAC address not found |

---

### `PUT /api/v1/machines/{mac}`

Update a machine's configuration.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mac` | string | Yes | MAC address (format: `aa:bb:cc:dd:ee:ff`) |

**Request Body:**

```json
{
  "profile_id": "ubuntu-server-v2",
  "metadata": {
    "datacenter": "homelab",
    "rack": "A1",
    "role": "storage"
  }
}
```

**Response (200 OK):**

Full machine configuration with updated fields.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine with specified MAC address not found |
| 422 Unprocessable Entity | Referenced profile_id does not exist |

---

### `DELETE /api/v1/machines/{mac}`

Delete a machine registration.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mac` | string | Yes | MAC address (format: `aa:bb:cc:dd:ee:ff`) |

**Response (204 No Content):**

Empty response body.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine with specified MAC address not found |

---

## Boot Profile Management

### `POST /api/v1/profiles`

Create a new boot profile.

**Request Body:**

```json
{
  "id": "ubuntu-server-base",
  "name": "Ubuntu Server Base Profile",
  "image_id": "ubuntu-2204",
  "kernel_args": [
    "console=tty0",
    "console=ttyS0",
    "ip=dhcp"
  ],
  "cloud_init_template": "ubuntu-base.yaml",
  "metadata": {
    "description": "Base Ubuntu server configuration with minimal packages",
    "tags": ["base", "minimal"]
  }
}
```

**Response (201 Created):**

```json
{
  "id": "ubuntu-server-base",
  "name": "Ubuntu Server Base Profile",
  "image_id": "ubuntu-2204",
  "kernel_args": [
    "console=tty0",
    "console=ttyS0",
    "ip=dhcp"
  ],
  "cloud_init_template": "ubuntu-base.yaml",
  "metadata": {
    "description": "Base Ubuntu server configuration with minimal packages",
    "tags": ["base", "minimal"]
  },
  "created_at": "2025-11-19T06:00:00Z",
  "updated_at": "2025-11-19T06:00:00Z"
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 409 Conflict | Profile with the same ID already exists |
| 422 Unprocessable Entity | Referenced image_id does not exist |

---

### `GET /api/v1/profiles`

List all boot profiles.

**Query Parameters:**

| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `page` | integer | No | Page number (1-indexed) | 1 |
| `per_page` | integer | No | Results per page (1-100) | 20 |
| `image_id` | string | No | Filter by boot image | - |

**Response (200 OK):**

```json
{
  "profiles": [
    {
      "id": "ubuntu-server-base",
      "name": "Ubuntu Server Base Profile",
      "image_id": "ubuntu-2204",
      "kernel_args": [
        "console=tty0",
        "console=ttyS0",
        "ip=dhcp"
      ],
      "cloud_init_template": "ubuntu-base.yaml",
      "metadata": {
        "description": "Base Ubuntu server configuration with minimal packages",
        "tags": ["base", "minimal"]
      },
      "created_at": "2025-11-19T06:00:00Z",
      "updated_at": "2025-11-19T06:00:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "per_page": 20,
    "total_pages": 1
  }
}
```

---

### `GET /api/v1/profiles/{id}`

Retrieve a specific boot profile.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot profile identifier |

**Response (200 OK):**

Same as POST response.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Profile with specified ID not found |

---

### `PUT /api/v1/profiles/{id}`

Update a boot profile.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot profile identifier |

**Request Body:**

```json
{
  "kernel_args": [
    "console=tty0",
    "console=ttyS0",
    "ip=dhcp",
    "net.ifnames=0"
  ],
  "metadata": {
    "description": "Updated Ubuntu server configuration",
    "tags": ["base", "minimal", "updated"]
  }
}
```

**Response (200 OK):**

Full profile configuration with updated fields.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Profile with specified ID not found |
| 422 Unprocessable Entity | Referenced image_id does not exist (if updated) |

---

### `DELETE /api/v1/profiles/{id}`

Delete a boot profile.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot profile identifier |

**Response (204 No Content):**

Empty response body.

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Profile with specified ID not found |
| 409 Conflict | Profile is currently assigned to one or more machines |

---

## Machine Rollback

### `POST /api/v1/machines/{mac}/rollback`

Rollback a machine to its previous boot profile.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mac` | string | Yes | MAC address (format: `aa:bb:cc:dd:ee:ff`) |

**Request Body:**

```json
{
  "reason": "Failed upgrade to new kernel version"
}
```

**Response (200 OK):**

```json
{
  "mac_address": "52:54:00:12:34:56",
  "hostname": "node-01",
  "profile_id": "ubuntu-server-base",
  "previous_profile_id": "ubuntu-server-v2",
  "rollback_reason": "Failed upgrade to new kernel version",
  "rollback_at": "2025-11-19T06:30:00Z",
  "rollback_by": "admin@example.com"
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | Machine with specified MAC address not found |
| 409 Conflict | No previous profile available for rollback |

**Rollback History:**

The system maintains a history of profile changes to enable rollback:

- Up to 10 previous profile assignments per machine
- Rollback can be performed multiple times (limited by history depth)
- History includes timestamp, user, and reason for each change

---

## Data Models

### Boot Image

```typescript
interface BootImage {
  id: string;              // Unique identifier (e.g., "ubuntu-2204")
  name: string;            // Human-readable name
  version: string;         // Semantic version
  kernel: {
    url: string;           // Cloud Storage URL
    sha256: string;        // SHA-256 checksum
    size_bytes: number;    // File size in bytes
  };
  initrd: {
    url: string;           // Cloud Storage URL
    sha256: string;        // SHA-256 checksum
    size_bytes: number;    // File size in bytes
  };
  metadata: {
    os: string;            // Operating system (ubuntu, fedora, talos)
    os_version: string;    // OS version
    architecture: string;  // CPU architecture (x86_64, arm64)
    tags: string[];        // Custom tags
  };
  created_at: string;      // ISO 8601 timestamp
  created_by: string;      // User or service account email
}
```

### Machine

```typescript
interface Machine {
  mac_address: string;     // MAC address (primary key)
  hostname: string;        // Machine hostname
  profile_id: string;      // Reference to boot profile
  metadata: {
    datacenter?: string;
    rack?: string;
    role?: string;
    [key: string]: any;    // Custom metadata
  };
  network: {
    ip_address?: string;   // Static IP (optional)
    netmask?: string;
    gateway?: string;
    dns_servers?: string[];
  };
  created_at: string;      // ISO 8601 timestamp
  updated_at: string;      // ISO 8601 timestamp
}
```

### Boot Profile

```typescript
interface BootProfile {
  id: string;              // Unique identifier
  name: string;            // Human-readable name
  image_id: string;        // Reference to boot image
  kernel_args: string[];   // Kernel command-line arguments
  cloud_init_template: string; // Cloud-init template filename
  metadata: {
    description?: string;
    tags?: string[];
    [key: string]: any;    // Custom metadata
  };
  created_at: string;      // ISO 8601 timestamp
  updated_at: string;      // ISO 8601 timestamp
}
```

---

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

---

## Versioning

The Admin API uses URL versioning (`/api/v1/`):

- **Current Version**: v1
- **Deprecation Policy**: Minimum 6 months notice before version deprecation
- **Version Header**: `X-API-Version: v1` included in all responses
