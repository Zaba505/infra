---
title: "Configuration Model"
type: docs
description: "Analysis of Matchbox's profile, group, and templating system"
---

# Matchbox Configuration Model

Matchbox uses a flexible configuration model based on **Profiles** (what to provision) and **Groups** (which machines get which profile), with support for templating and metadata.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Matchbox Store                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  Profiles  │  │   Groups   │  │   Assets   │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│        │               │                                    │
│        │               │                                    │
│        ▼               ▼                                    │
│  ┌─────────────────────────────────────┐                   │
│  │       Matcher Engine                │                   │
│  │  (Label-based group selection)      │                   │
│  └─────────────────────────────────────┘                   │
│                    │                                        │
│                    ▼                                        │
│  ┌─────────────────────────────────────┐                   │
│  │    Template Renderer                │                   │
│  │  (Go templates + metadata)          │                   │
│  └─────────────────────────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
            Rendered Config (iPXE, Ignition, etc.)
```

## Data Directory Structure

Matchbox uses a FileStore (default) that reads from `-data-path` (default: `/var/lib/matchbox`):

```
/var/lib/matchbox/
├── groups/              # Machine group definitions (JSON)
│   ├── default.json
│   ├── node1.json
│   └── us-west.json
├── profiles/            # Profile definitions (JSON)
│   ├── worker.json
│   ├── controller.json
│   └── etcd.json
├── ignition/            # Ignition configs (.ign) or Butane (.yaml)
│   ├── worker.ign
│   ├── controller.ign
│   └── butane-example.yaml
├── cloud/               # Cloud-Config templates (DEPRECATED)
│   └── legacy.yaml.tmpl
├── generic/             # Arbitrary config templates
│   ├── setup.cfg
│   └── metadata.yaml.tmpl
└── assets/              # Static files (kernel, initrd)
    ├── fedora-coreos/
    └── flatcar/
```

**Version control:** Poseidon recommends keeping `/var/lib/matchbox` under git for auditability and rollback.

## Profiles

Profiles define **what to provision**: network boot settings (kernel, initrd, args) and config references (Ignition, Cloud-Config, generic).

### Profile Schema

```json
{
  "id": "worker",
  "name": "Fedora CoreOS Worker Node",
  "boot": {
    "kernel": "/assets/fedora-coreos/36.20220906.3.2/fedora-coreos-36.20220906.3.2-live-kernel-x86_64",
    "initrd": [
      "--name main /assets/fedora-coreos/36.20220906.3.2/fedora-coreos-36.20220906.3.2-live-initramfs.x86_64.img"
    ],
    "args": [
      "initrd=main",
      "coreos.live.rootfs_url=http://matchbox.example.com:8080/assets/fedora-coreos/36.20220906.3.2/fedora-coreos-36.20220906.3.2-live-rootfs.x86_64.img",
      "coreos.inst.install_dev=/dev/sda",
      "coreos.inst.ignition_url=http://matchbox.example.com:8080/ignition?uuid=${uuid}&mac=${mac:hexhyp}"
    ]
  },
  "ignition_id": "worker.ign",
  "cloud_id": "",
  "generic_id": ""
}
```

### Profile Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | ✅ | Unique profile identifier (referenced by groups) |
| `name` | string | ❌ | Human-readable description |
| `boot` | object | ❌ | Network boot configuration |
| `boot.kernel` | string | ❌ | Kernel URL (HTTP/HTTPS or /assets path) |
| `boot.initrd` | array | ❌ | Initrd URLs (can specify `--name` for multi-initrd) |
| `boot.args` | array | ❌ | Kernel command-line arguments |
| `ignition_id` | string | ❌ | Ignition/Butane config filename in `ignition/` |
| `cloud_id` | string | ❌ | Cloud-Config filename in `cloud/` (deprecated) |
| `generic_id` | string | ❌ | Generic config filename in `generic/` |

### Boot Configuration Patterns

#### Pattern 1: Live PXE (RAM-based, ephemeral)

Boot and run OS entirely from RAM, no disk install:

```json
{
  "boot": {
    "kernel": "/assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-kernel-x86_64",
    "initrd": [
      "--name main /assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-initramfs.x86_64.img"
    ],
    "args": [
      "initrd=main",
      "coreos.live.rootfs_url=http://matchbox/assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-rootfs.x86_64.img",
      "ignition.config.url=http://matchbox/ignition?uuid=${uuid}&mac=${mac:hexhyp}"
    ]
  }
}
```

**Use case:** Diskless workers, testing, ephemeral compute

#### Pattern 2: Disk Install (persistent)

PXE boot live image, install to disk, reboot to disk:

```json
{
  "boot": {
    "kernel": "/assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-kernel-x86_64",
    "initrd": [
      "--name main /assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-initramfs.x86_64.img"
    ],
    "args": [
      "initrd=main",
      "coreos.live.rootfs_url=http://matchbox/assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-rootfs.x86_64.img",
      "coreos.inst.install_dev=/dev/sda",
      "coreos.inst.ignition_url=http://matchbox/ignition?uuid=${uuid}&mac=${mac:hexhyp}"
    ]
  }
}
```

**Key difference:** `coreos.inst.install_dev` triggers disk install before reboot

#### Pattern 3: Multi-initrd (layered)

Multiple initrds can be loaded (e.g., base + drivers):

```json
{
  "initrd": [
    "--name main /assets/fedora-coreos/VERSION/fedora-coreos-VERSION-live-initramfs.x86_64.img",
    "--name drivers /assets/drivers/custom-drivers.img"
  ],
  "args": [
    "initrd=main,drivers",
    "..."
  ]
}
```

### Config References

#### Ignition Configs

**Direct Ignition (.ign files):**
```json
{
  "ignition_id": "worker.ign"
}
```

File: `/var/lib/matchbox/ignition/worker.ign`
```json
{
  "ignition": { "version": "3.3.0" },
  "systemd": {
    "units": [{
      "name": "example.service",
      "enabled": true,
      "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello\n\n[Install]\nWantedBy=multi-user.target"
    }]
  }
}
```

**Butane Configs (transpiled to Ignition):**
```json
{
  "ignition_id": "worker.yaml"
}
```

File: `/var/lib/matchbox/ignition/worker.yaml`
```yaml
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-ed25519 AAAA...
systemd:
  units:
    - name: etcd.service
      enabled: true
```

**Matchbox automatically:**
1. Detects Butane format (file doesn't end in `.ign` or `.ignition`)
2. Transpiles Butane → Ignition using embedded library
3. Renders templates with group metadata
4. Serves as Ignition v3.3.0

#### Generic Configs

For non-Ignition configs (scripts, YAML, arbitrary data):

```json
{
  "generic_id": "setup-script.sh.tmpl"
}
```

File: `/var/lib/matchbox/generic/setup-script.sh.tmpl`
```bash
#!/bin/bash
# Rendered with group metadata
NODE_NAME={{.node_name}}
CLUSTER_ID={{.cluster_id}}
echo "Provisioning ${NODE_NAME} in cluster ${CLUSTER_ID}"
```

**Access via:** `GET /generic?uuid=...&mac=...`

## Groups

Groups match machines to profiles using **selectors** (label matching) and provide **metadata** for template rendering.

### Group Schema

```json
{
  "id": "node1-worker",
  "name": "Worker Node 1",
  "profile": "worker",
  "selector": {
    "mac": "52:54:00:89:d8:10",
    "uuid": "550e8400-e29b-41d4-a716-446655440000"
  },
  "metadata": {
    "node_name": "worker-01",
    "cluster_id": "prod-cluster",
    "etcd_endpoints": "https://10.0.1.10:2379,https://10.0.1.11:2379",
    "ssh_authorized_keys": [
      "ssh-ed25519 AAAA...",
      "ssh-rsa AAAA..."
    ]
  }
}
```

### Group Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | ✅ | Unique group identifier |
| `name` | string | ❌ | Human-readable description |
| `profile` | string | ✅ | Profile ID to apply |
| `selector` | object | ❌ | Label match criteria (omit for default group) |
| `metadata` | object | ❌ | Key-value data for template rendering |

### Selector Matching

**Reserved selectors** (automatically populated from machine attributes):

| Selector | Source | Example | Normalized |
|----------|--------|---------|------------|
| `uuid` | SMBIOS UUID | `550e8400-e29b-41d4-a716-446655440000` | Lowercase |
| `mac` | Primary NIC MAC | `52:54:00:89:d8:10` | Colon-separated |
| `hostname` | Network hostname | `node1.example.com` | As reported |
| `serial` | Hardware serial | `VMware-42 1a...` | As reported |

**Custom selectors** (passed as query params):
```json
{
  "selector": {
    "region": "us-west",
    "environment": "production",
    "rack": "A23"
  }
}
```

**Matching request:** `/ipxe?mac=52:54:00:89:d8:10&region=us-west&environment=production&rack=A23`

**Matching logic:**
1. All selector key-value pairs must match request labels (AND logic)
2. Most specific group wins (most selector matches)
3. If multiple groups have same specificity, first match wins (undefined order)
4. Groups with no selectors = default group (matches all)

### Default Groups

Group with empty `selector` matches all machines:

```json
{
  "id": "default-worker",
  "name": "Default Worker",
  "profile": "worker",
  "metadata": {
    "environment": "dev"
  }
}
```

⚠️ **Warning:** Avoid multiple default groups (non-deterministic matching)

### Example: Region-based Matching

**Group 1: US-West Workers**
```json
{
  "id": "us-west-workers",
  "profile": "worker",
  "selector": {
    "region": "us-west"
  },
  "metadata": {
    "etcd_endpoints": "https://etcd-usw.example.com:2379"
  }
}
```

**Group 2: EU Workers**
```json
{
  "id": "eu-workers",
  "profile": "worker",
  "selector": {
    "region": "eu"
  },
  "metadata": {
    "etcd_endpoints": "https://etcd-eu.example.com:2379"
  }
}
```

**Group 3: Specific Machine Override**
```json
{
  "id": "node-special",
  "profile": "controller",
  "selector": {
    "mac": "52:54:00:89:d8:10",
    "region": "us-west"
  },
  "metadata": {
    "role": "controller"
  }
}
```

**Matching precedence:**
- Machine with `mac=52:54:00:89:d8:10&region=us-west` → `node-special` (2 selectors)
- Machine with `region=us-west` → `us-west-workers` (1 selector)
- Machine with `region=eu` → `eu-workers` (1 selector)

## Templating System

Matchbox uses Go's `text/template` for rendering configs with group metadata.

### Template Context

Available variables in Ignition/Butane/Cloud-Config/generic templates:

```go
// Group metadata (all keys from group.metadata)
{{.node_name}}
{{.cluster_id}}
{{.etcd_endpoints}}

// Group selectors (normalized)
{{.mac}}      // e.g., "52:54:00:89:d8:10"
{{.uuid}}     // e.g., "550e8400-..."
{{.region}}   // Custom selector

// Request query params (raw)
{{.request.query.mac}}     // As passed in URL
{{.request.query.foo}}     // Custom query param
{{.request.raw_query}}     // Full query string

// Special functions
{{if index . "ssh_authorized_keys"}}  // Check if key exists
{{range $element := .ssh_authorized_keys}}  // Iterate arrays
```

### Example: Templated Butane Config

**Group metadata:**
```json
{
  "metadata": {
    "node_name": "worker-01",
    "ssh_authorized_keys": [
      "ssh-ed25519 AAA...",
      "ssh-rsa BBB..."
    ],
    "ntp_servers": ["time1.google.com", "time2.google.com"]
  }
}
```

**Butane template:** `/var/lib/matchbox/ignition/worker.yaml`
```yaml
variant: fcos
version: 1.5.0

storage:
  files:
    - path: /etc/hostname
      mode: 0644
      contents:
        inline: {{.node_name}}

    - path: /etc/systemd/timesyncd.conf
      mode: 0644
      contents:
        inline: |
          [Time]
          {{range $server := .ntp_servers}}
          NTP={{$server}}
          {{end}}

{{if index . "ssh_authorized_keys"}}
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        {{range $key := .ssh_authorized_keys}}
        - {{$key}}
        {{end}}
{{end}}
```

**Rendered Ignition (simplified):**
```json
{
  "ignition": {"version": "3.3.0"},
  "storage": {
    "files": [
      {
        "path": "/etc/hostname",
        "contents": {"source": "data:,worker-01"},
        "mode": 420
      },
      {
        "path": "/etc/systemd/timesyncd.conf",
        "contents": {"source": "data:,%5BTime%5D%0ANTP%3Dtime1.google.com%0ANTP%3Dtime2.google.com"},
        "mode": 420
      }
    ]
  },
  "passwd": {
    "users": [{
      "name": "core",
      "sshAuthorizedKeys": ["ssh-ed25519 AAA...", "ssh-rsa BBB..."]
    }]
  }
}
```

### Template Best Practices

1. **Prefer external rendering:** Use Terraform + `ct_config` provider for complex templates
2. **Validate Butane:** Use `strict: true` in Terraform or `fcct --strict`
3. **Escape carefully:** Go templates use `{{}}`, Butane uses YAML - mind the interaction
4. **Test rendering:** Request `/ignition?mac=...` directly to inspect output
5. **Version control:** Keep templates + groups in git for auditability

### Reserved Metadata Keys

**Warning:** `.request` is reserved for query param access. Group metadata with `"request": {...}` will be overwritten.

**Reserved keys:**
- `request.query.*` - Query parameters
- `request.raw_query` - Raw query string

## API Integration

### HTTP Endpoints (Read-only)

| Endpoint | Purpose | Template Context |
|----------|---------|------------------|
| `/ipxe` | iPXE boot script | Profile `boot` section |
| `/grub` | GRUB config | Profile `boot` section |
| `/ignition` | Ignition config | Group metadata + selectors + query |
| `/cloud` | Cloud-Config (deprecated) | Group metadata + selectors + query |
| `/generic` | Generic config | Group metadata + selectors + query |
| `/metadata` | Key-value env format | Group metadata + selectors + query |

**Example metadata endpoint response:**
```
GET /metadata?mac=52:54:00:89:d8:10&foo=bar

NODE_NAME=worker-01
CLUSTER_ID=prod
MAC=52:54:00:89:d8:10
REQUEST_QUERY_MAC=52:54:00:89:d8:10
REQUEST_QUERY_FOO=bar
REQUEST_RAW_QUERY=mac=52:54:00:89:d8:10&foo=bar
```

### gRPC API (Authenticated, mutable)

Used by `terraform-provider-matchbox` for declarative infrastructure:

**Terraform example:**
```hcl
provider "matchbox" {
  endpoint    = "matchbox.example.com:8081"
  client_cert = file("~/.matchbox/client.crt")
  client_key  = file("~/.matchbox/client.key")
  ca          = file("~/.matchbox/ca.crt")
}

resource "matchbox_profile" "worker" {
  name   = "worker"
  kernel = "/assets/fedora-coreos/.../kernel"
  initrd = ["--name main /assets/fedora-coreos/.../initramfs.img"]
  args   = [
    "initrd=main",
    "coreos.inst.install_dev=/dev/sda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}"
  ]
  raw_ignition = data.ct_config.worker.rendered
}

resource "matchbox_group" "node1" {
  name    = "node1"
  profile = matchbox_profile.worker.name
  selector = {
    mac = "52:54:00:89:d8:10"
  }
  metadata = {
    node_name = "worker-01"
  }
}
```

**Operations:**
- `CreateProfile`, `GetProfile`, `UpdateProfile`, `DeleteProfile`
- `CreateGroup`, `GetGroup`, `UpdateGroup`, `DeleteGroup`

**TLS client authentication required** (see deployment docs)

## Configuration Workflow

### Recommended: Terraform + External Configs

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Write Butane configs (YAML)                             │
│    - worker.yaml, controller.yaml                          │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Terraform ct_config transpiles Butane → Ignition        │
│    data "ct_config" "worker" {                             │
│      content = file("worker.yaml")                         │
│      strict  = true                                        │
│    }                                                        │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Terraform creates profiles + groups in Matchbox         │
│    matchbox_profile.worker → gRPC CreateProfile()          │
│    matchbox_group.node1 → gRPC CreateGroup()               │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Machine PXE boots, queries Matchbox                     │
│    GET /ipxe?mac=... → matches group → returns profile     │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Ignition fetches rendered config                        │
│    GET /ignition?mac=... → Matchbox returns Ignition       │
└─────────────────────────────────────────────────────────────┘
```

**Benefits:**
- Rich Terraform templating (loops, conditionals, external data sources)
- Butane validation before deployment
- Declarative infrastructure (can `terraform plan` before apply)
- Version control workflow (git + CI/CD)

### Alternative: Manual FileStore

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Create profile JSON manually                            │
│    /var/lib/matchbox/profiles/worker.json                  │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Create group JSON manually                              │
│    /var/lib/matchbox/groups/node1.json                     │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Write Ignition/Butane config                            │
│    /var/lib/matchbox/ignition/worker.ign                   │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Restart matchbox (to reload FileStore)                  │
│    systemctl restart matchbox                              │
└─────────────────────────────────────────────────────────────┘
```

**Drawbacks:**
- Manual file management
- No validation before deployment
- Requires matchbox restart to pick up changes
- Error-prone for large fleets

## Storage Backends

### FileStore (Default)

**Config:** `-data-path=/var/lib/matchbox`

**Pros:**
- Simple file-based storage
- Easy to version control (git)
- Human-readable JSON

**Cons:**
- Requires file system access
- Manual reload for gRPC-created resources

### Custom Store (Extensible)

Matchbox's `Store` interface allows custom backends:

```go
type Store interface {
  ProfileGet(id string) (*Profile, error)
  GroupGet(id string) (*Group, error)
  IgnitionGet(name string) (string, error)
  // ... other methods
}
```

**Potential custom stores:**
- etcd backend (for HA Matchbox)
- Database backend (PostgreSQL, MySQL)
- S3/object storage backend

**Note:** Not officially provided by Matchbox project; requires custom implementation

## Security Considerations

1. **gRPC API authentication:** Requires TLS client certificates
   - `ca.crt` - CA that signed client certs
   - `server.crt`/`server.key` - Server TLS identity
   - `client.crt`/`client.key` - Client credentials (Terraform)

2. **HTTP endpoints are read-only:** No auth, machines fetch configs
   - Do NOT put secrets in Ignition configs
   - Use external secret stores (Vault, GCP Secret Manager)
   - Reference secrets via Ignition `files.source` with auth headers

3. **Network segmentation:** Matchbox on provisioning VLAN, isolate from production

4. **Config validation:** Validate Ignition/Butane before deployment to avoid boot failures

5. **Audit logging:** Version control groups/profiles; log gRPC API changes

## Operational Tips

1. **Test groups with curl:**
   ```bash
   curl 'http://matchbox.example.com:8080/ignition?mac=52:54:00:89:d8:10'
   ```

2. **List profiles:**
   ```bash
   ls -la /var/lib/matchbox/profiles/
   ```

3. **Validate Butane:**
   ```bash
   podman run -i --rm quay.io/coreos/fcct:release --strict < worker.yaml
   ```

4. **Check group matching:**
   ```bash
   # Default group (no selectors)
   curl http://matchbox.example.com:8080/ignition
   
   # Specific machine
   curl 'http://matchbox.example.com:8080/ignition?mac=52:54:00:89:d8:10&uuid=550e8400-e29b-41d4-a716-446655440000'
   ```

5. **Backup configs:**
   ```bash
   tar -czf matchbox-backup-$(date +%F).tar.gz /var/lib/matchbox/{groups,profiles,ignition}
   ```

## Summary

Matchbox's configuration model provides:

- **Separation of concerns:** Profiles (what) vs Groups (who/where)
- **Flexible matching:** Label-based, multi-attribute, custom selectors
- **Template support:** Go templates for dynamic configs (but prefer external rendering)
- **API-driven:** Terraform integration for GitOps workflows
- **Storage options:** FileStore (simple) or custom backends (extensible)
- **OS-agnostic:** Works with any Ignition-based distro (FCOS, Flatcar, RHCOS)

**Best practice:** Use Terraform + external Butane configs for production; manual FileStore for labs/development.
