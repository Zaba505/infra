---
title: "UEFI HTTP Boot Endpoints"
type: docs
description: "Boot endpoints accessed by bare metal servers during network boot"
weight: 10
---

These endpoints are accessed by bare metal servers (HP DL360 Gen 9) during the UEFI HTTP boot process. All endpoints are accessed through the WireGuard VPN tunnel and use source IP validation for security.

## Boot Script Endpoint

### `GET /boot.ipxe`

Serves iPXE boot scripts customized for the requesting machine based on its MAC address.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mac` | string | Yes | MAC address of the requesting machine (format: `aa:bb:cc:dd:ee:ff`) |

**Request Example:**

```http
GET /boot.ipxe?mac=52:54:00:12:34:56 HTTP/1.1
Host: boot.internal
```

**Response Example (200 OK):**

```text
#!ipxe

# Boot configuration for node-01 (52:54:00:12:34:56)
# Profile: ubuntu-22.04-server
# Generated: 2025-11-19T06:00:00Z

kernel /assets/ubuntu-2204/kernel console=tty0 console=ttyS0 ip=dhcp
initrd /assets/ubuntu-2204/initrd
boot
```

**Response Headers:**

- `Content-Type: text/plain; charset=utf-8`
- `Cache-Control: no-cache, no-store, must-revalidate`

**Error Responses:**

| Status Code | Description | Example |
|-------------|-------------|---------|
| 400 Bad Request | Missing or invalid MAC address | `{"error": {"code": "INVALID_MAC_ADDRESS", "message": "MAC address must be in format aa:bb:cc:dd:ee:ff"}}` |
| 404 Not Found | No boot configuration found for MAC | `{"error": {"code": "MACHINE_NOT_CONFIGURED", "message": "No boot configuration found for MAC 52:54:00:12:34:56"}}` |
| 500 Internal Server Error | Database or template error | `{"error": {"code": "INTERNAL_ERROR", "message": "Failed to generate boot script"}}` |

**Boot Script Variables:**

The iPXE script may include the following dynamic values:

- Machine-specific kernel parameters
- Cloud-init data source URLs
- Asset download URLs
- Network configuration parameters

---

## Kernel Image Endpoint

### `GET /assets/{id}/kernel`

Streams kernel images from Cloud Storage for the boot process.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot image identifier (e.g., `ubuntu-2204`, `talos-v1.8`) |

**Request Example:**

```http
GET /assets/ubuntu-2204/kernel HTTP/1.1
Host: boot.internal
```

**Response Example (200 OK):**

Binary kernel image streamed from Cloud Storage.

**Response Headers:**

- `Content-Type: application/octet-stream`
- `Content-Length: 8388608` (actual kernel size in bytes)
- `Cache-Control: public, max-age=3600`
- `ETag: "abc123..."`

**Error Responses:**

| Status Code | Description | Example |
|-------------|-------------|---------|
| 404 Not Found | Kernel image not found | `{"error": {"code": "KERNEL_NOT_FOUND", "message": "Kernel image 'ubuntu-2204' not found"}}` |
| 500 Internal Server Error | Cloud Storage error | `{"error": {"code": "STORAGE_ERROR", "message": "Failed to retrieve kernel from storage"}}` |

**Performance Characteristics:**

- **Streaming**: File is streamed directly from Cloud Storage (no buffering in memory)
- **Target Latency**: < 100ms to first byte
- **Typical Size**: 8-15 MB for Linux kernels

---

## Initrd Image Endpoint

### `GET /assets/{id}/initrd`

Streams initial ramdisk (initrd) images from Cloud Storage for the boot process.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Boot image identifier (e.g., `ubuntu-2204`, `talos-v1.8`) |

**Request Example:**

```http
GET /assets/ubuntu-2204/initrd HTTP/1.1
Host: boot.internal
```

**Response Example (200 OK):**

Binary initrd image streamed from Cloud Storage.

**Response Headers:**

- `Content-Type: application/octet-stream`
- `Content-Length: 52428800` (actual initrd size in bytes)
- `Cache-Control: public, max-age=3600`
- `ETag: "def456..."`

**Error Responses:**

| Status Code | Description | Example |
|-------------|-------------|---------|
| 404 Not Found | Initrd image not found | `{"error": {"code": "INITRD_NOT_FOUND", "message": "Initrd image 'ubuntu-2204' not found"}}` |
| 500 Internal Server Error | Cloud Storage error | `{"error": {"code": "STORAGE_ERROR", "message": "Failed to retrieve initrd from storage"}}` |

**Performance Characteristics:**

- **Streaming**: File is streamed directly from Cloud Storage (no buffering in memory)
- **Target Latency**: < 100ms to first byte
- **Typical Size**: 50-150 MB for Linux initrd images

---

## Cloud-Init Configuration Endpoint

### `GET /cloud-init/{machine_id}`

Serves cloud-init configuration files customized for specific machines.

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `machine_id` | string | Yes | Machine identifier (hostname, UUID, or MAC address) |

**Request Example:**

```http
GET /cloud-init/node-01 HTTP/1.1
Host: boot.internal
```

**Response Example (200 OK):**

```yaml
#cloud-config

hostname: node-01
fqdn: node-01.homelab.local

users:
  - name: admin
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh_authorized_keys:
      - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample...

packages:
  - qemu-guest-agent
  - curl

runcmd:
  - systemctl enable qemu-guest-agent
  - systemctl start qemu-guest-agent

final_message: "Cloud-init complete for node-01 after $UPTIME seconds"
```

**Response Headers:**

- `Content-Type: text/cloud-config; charset=utf-8`
- `Cache-Control: no-cache, no-store, must-revalidate`

**Error Responses:**

| Status Code | Description | Example |
|-------------|-------------|---------|
| 404 Not Found | Machine configuration not found | `{"error": {"code": "MACHINE_NOT_FOUND", "message": "Cloud-init configuration for 'node-01' not found"}}` |
| 500 Internal Server Error | Template or storage error | `{"error": {"code": "INTERNAL_ERROR", "message": "Failed to generate cloud-init configuration"}}` |

**Cloud-Init Features:**

The cloud-init configuration supports:

- Hostname and FQDN configuration
- User account creation with SSH keys
- Package installation
- Custom commands (runcmd)
- Network configuration
- Disk partitioning and filesystem setup
- Service management

**Template Variables:**

Cloud-init templates can use machine-specific variables:

- `{{.Hostname}}` - Machine hostname
- `{{.FQDN}}` - Fully qualified domain name
- `{{.MACAddress}}` - Primary MAC address
- `{{.IPAddress}}` - Assigned IP address (if static)
- `{{.SSHKeys}}` - Authorized SSH public keys
- `{{.Metadata}}` - Custom machine metadata

---

## Security Considerations

### VPN Source IP Validation

All boot endpoints validate that requests originate from the WireGuard VPN subnet:

- **Allowed CIDR**: `10.x.x.0/24` (WireGuard VPN network)
- **Validation**: Performed at Cloud Run ingress or application layer
- **Rejection**: Requests from outside VPN return `403 Forbidden`

### Rate Limiting

To prevent abuse, boot endpoints are rate-limited:

- **Boot Script**: 10 requests/minute per MAC address
- **Asset Downloads**: 5 concurrent downloads per MAC address
- **Cloud-Init**: 10 requests/minute per machine_id

### Asset Integrity

Boot assets are validated for integrity:

- **Checksums**: SHA-256 checksums stored in Firestore
- **Verification**: Computed on upload, verified on download (optional)
- **ETag Headers**: Enable client-side caching and integrity checks

## Observability

All boot endpoint requests are instrumented with OpenTelemetry:

- **Metrics**: Request count, latency, error rate, bytes transferred
- **Traces**: End-to-end tracing from request to Cloud Storage retrieval
- **Logs**: Structured logs with MAC address, image ID, response status

**Key Metrics:**

- `boot_script_requests_total` - Total boot script requests
- `boot_script_latency_ms` - Boot script generation latency
- `asset_download_bytes_total` - Total bytes transferred for boot assets
- `asset_download_duration_ms` - Asset download duration
- `cloud_init_requests_total` - Cloud-init configuration requests
