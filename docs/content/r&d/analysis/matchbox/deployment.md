---
title: "Deployment Patterns"
type: docs
weight: 4
description: "Matchbox deployment options and operational considerations"
---

# Matchbox Deployment Patterns

Analysis of deployment architectures, installation methods, and operational considerations for running Matchbox in production.

## Deployment Architectures

### Single-Host Deployment

```
┌─────────────────────────────────────────────────────┐
│           Provisioning Host                         │
│  ┌─────────────┐        ┌─────────────┐            │
│  │  Matchbox   │        │  dnsmasq    │            │
│  │  :8080 HTTP │        │  DHCP/TFTP  │            │
│  │  :8081 gRPC │        │  :67,:69    │            │
│  └─────────────┘        └─────────────┘            │
│         │                      │                    │
│         └──────────┬───────────┘                    │
│                    │                                │
│  /var/lib/matchbox/                                 │
│  ├── groups/                                        │
│  ├── profiles/                                      │
│  ├── ignition/                                      │
│  └── assets/                                        │
└─────────────────────────────────────────────────────┘
              │
              │ Network
              ▼
     ┌──────────────┐
     │ PXE Clients  │
     └──────────────┘
```

**Use case:** Lab, development, small deployments (<50 machines)

**Pros:**
- Simple setup
- Single service to manage
- Minimal resource requirements

**Cons:**
- Single point of failure
- No scalability
- Downtime during updates

### HA Deployment (Multiple Matchbox Instances)

```
┌─────────────────────────────────────────────────────┐
│              Load Balancer (Ingress/HAProxy)        │
│           :8080 HTTP        :8081 gRPC              │
└─────────────────────────────────────────────────────┘
       │                              │
       ├─────────────┬────────────────┤
       ▼             ▼                ▼
┌──────────┐  ┌──────────┐    ┌──────────┐
│Matchbox 1│  │Matchbox 2│    │Matchbox N│
│ (Pod/VM) │  │ (Pod/VM) │    │ (Pod/VM) │
└──────────┘  └──────────┘    └──────────┘
       │             │                │
       └─────────────┴────────────────┘
                     │
                     ▼
         ┌────────────────────────┐
         │  Shared Storage        │
         │  /var/lib/matchbox     │
         │  (NFS, PV, ConfigMap)  │
         └────────────────────────┘
```

**Use case:** Production, datacenter-scale (100+ machines)

**Pros:**
- High availability (no single point of failure)
- Rolling updates (zero downtime)
- Load distribution

**Cons:**
- Complex storage (shared volume or etcd backend)
- More infrastructure required

**Storage options:**
1. **Kubernetes PersistentVolume** (RWX mode)
2. **NFS share** mounted on multiple hosts
3. **Custom etcd-backed Store** (requires custom implementation)
4. **Git-sync sidecar** (read-only, periodic pull)

### Kubernetes Deployment

```
┌─────────────────────────────────────────────────────┐
│              Ingress Controller                     │
│  matchbox.example.com → Service matchbox:8080       │
│  matchbox-rpc.example.com → Service matchbox:8081   │
└─────────────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────┐
│          Service: matchbox (ClusterIP)              │
│            ports: 8080/TCP, 8081/TCP                │
└─────────────────────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│  Pod: matchbox  │     │  Pod: matchbox  │
│  replicas: 2+   │     │  replicas: 2+   │
└─────────────────┘     └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     ▼
┌─────────────────────────────────────────────────────┐
│    PersistentVolumeClaim: matchbox-data             │
│    /var/lib/matchbox (RWX mode)                     │
└─────────────────────────────────────────────────────┘
```

**Manifest structure:**
```
contrib/k8s/
├── matchbox-deployment.yaml  # Deployment + replicas
├── matchbox-service.yaml     # Service (8080, 8081)
├── matchbox-ingress.yaml     # Ingress (HTTP + gRPC TLS)
└── matchbox-pvc.yaml         # PersistentVolumeClaim
```

**Key configurations:**

1. **Secret for gRPC TLS:**
   ```bash
   kubectl create secret generic matchbox-rpc \
     --from-file=ca.crt \
     --from-file=server.crt \
     --from-file=server.key
   ```

2. **Ingress for gRPC (TLS passthrough):**
   ```yaml
   metadata:
     annotations:
       nginx.ingress.kubernetes.io/ssl-passthrough: "true"
       nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
   ```

3. **Volume mount:**
   ```yaml
   volumes:
     - name: data
       persistentVolumeClaim:
         claimName: matchbox-data
   volumeMounts:
     - name: data
       mountPath: /var/lib/matchbox
   ```

**Use case:** Cloud-native deployments, Kubernetes-based infrastructure

**Pros:**
- Native Kubernetes primitives (Deployments, Services, Ingress)
- Rolling updates via Deployment strategy
- Easy scaling (`kubectl scale`)
- Health checks + auto-restart

**Cons:**
- Requires RWX PersistentVolume or shared storage
- Ingress TLS configuration complexity (gRPC passthrough)
- Cluster dependency (can't provision cluster bootstrap nodes)

⚠️ **Bootstrap problem:** Kubernetes-hosted Matchbox can't PXE boot its own cluster nodes (chicken-and-egg). Use external Matchbox for initial cluster bootstrap, then migrate.

## Installation Methods

### 1. Binary Installation (systemd)

**Recommended for:** Bare-metal hosts, VMs, traditional Linux servers

**Steps:**

1. **Download and verify:**
   ```bash
   wget https://github.com/poseidon/matchbox/releases/download/v0.10.0/matchbox-v0.10.0-linux-amd64.tar.gz
   wget https://github.com/poseidon/matchbox/releases/download/v0.10.0/matchbox-v0.10.0-linux-amd64.tar.gz.asc
   gpg --verify matchbox-v0.10.0-linux-amd64.tar.gz.asc
   ```

2. **Extract and install:**
   ```bash
   tar xzf matchbox-v0.10.0-linux-amd64.tar.gz
   sudo cp matchbox-v0.10.0-linux-amd64/matchbox /usr/local/bin/
   ```

3. **Create user and directories:**
   ```bash
   sudo useradd -U matchbox
   sudo mkdir -p /var/lib/matchbox/{assets,groups,profiles,ignition}
   sudo chown -R matchbox:matchbox /var/lib/matchbox
   ```

4. **Install systemd unit:**
   ```bash
   sudo cp contrib/systemd/matchbox.service /etc/systemd/system/
   ```

5. **Configure via systemd dropin:**
   ```bash
   sudo systemctl edit matchbox
   ```
   ```ini
   [Service]
   Environment="MATCHBOX_ADDRESS=0.0.0.0:8080"
   Environment="MATCHBOX_RPC_ADDRESS=0.0.0.0:8081"
   Environment="MATCHBOX_LOG_LEVEL=debug"
   ```

6. **Start service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start matchbox
   sudo systemctl enable matchbox
   ```

**Pros:**
- Direct control over service
- Easy log access (`journalctl -u matchbox`)
- Native OS integration

**Cons:**
- Manual updates required
- OS dependency (package compatibility)

### 2. Container Deployment (Docker/Podman)

**Recommended for:** Docker hosts, quick testing, immutable infrastructure

**Docker:**
```bash
mkdir -p /var/lib/matchbox/assets
docker run -d --name matchbox \
  --net=host \
  -v /var/lib/matchbox:/var/lib/matchbox:Z \
  -v /etc/matchbox:/etc/matchbox:Z,ro \
  quay.io/poseidon/matchbox:v0.10.0 \
  -address=0.0.0.0:8080 \
  -rpc-address=0.0.0.0:8081 \
  -log-level=debug
```

**Podman:**
```bash
podman run -d --name matchbox \
  --net=host \
  -v /var/lib/matchbox:/var/lib/matchbox:Z \
  -v /etc/matchbox:/etc/matchbox:Z,ro \
  quay.io/poseidon/matchbox:v0.10.0 \
  -address=0.0.0.0:8080 \
  -rpc-address=0.0.0.0:8081 \
  -log-level=debug
```

**Volume mounts:**
- `/var/lib/matchbox` - Data directory (groups, profiles, configs, assets)
- `/etc/matchbox` - TLS certificates (ca.crt, server.crt, server.key)

**Network mode:**
- `--net=host` - Required for DHCP/TFTP interaction on same host
- Bridge mode possible if Matchbox is on separate host from dnsmasq

**Pros:**
- Immutable deployments
- Easy updates (pull new image)
- Portable across hosts

**Cons:**
- Volume management complexity
- SELinux considerations (`:Z` flag)

### 3. Kubernetes Deployment

**Recommended for:** Kubernetes environments, cloud platforms

**Quick start:**
```bash
# Create TLS secret for gRPC
kubectl create secret generic matchbox-rpc \
  --from-file=ca.crt=~/.matchbox/ca.crt \
  --from-file=server.crt=~/.matchbox/server.crt \
  --from-file=server.key=~/.matchbox/server.key

# Deploy manifests
kubectl apply -R -f contrib/k8s/

# Check status
kubectl get pods -l app=matchbox
kubectl get svc matchbox
kubectl get ingress matchbox matchbox-rpc
```

**Persistence options:**

**Option 1: emptyDir (ephemeral, dev only):**
```yaml
volumes:
  - name: data
    emptyDir: {}
```

**Option 2: PersistentVolumeClaim (production):**
```yaml
volumes:
  - name: data
    persistentVolumeClaim:
      claimName: matchbox-data
```

**Option 3: ConfigMap (static configs):**
```yaml
volumes:
  - name: groups
    configMap:
      name: matchbox-groups
  - name: profiles
    configMap:
      name: matchbox-profiles
```

**Option 4: Git-sync sidecar (GitOps):**
```yaml
initContainers:
  - name: git-sync
    image: k8s.gcr.io/git-sync:v3.6.3
    env:
      - name: GIT_SYNC_REPO
        value: https://github.com/example/matchbox-configs
      - name: GIT_SYNC_DEST
        value: /var/lib/matchbox
    volumeMounts:
      - name: data
        mountPath: /var/lib/matchbox
```

**Pros:**
- Native k8s features (scaling, health checks, rolling updates)
- Ingress integration
- GitOps workflows

**Cons:**
- Complexity (Ingress, PVC, TLS)
- Can't bootstrap own cluster

## Network Boot Environment Setup

Matchbox requires separate DHCP/TFTP/DNS services. Options:

### Option 1: dnsmasq Container (Quickest)

**Use case:** Lab, testing, environments without existing DHCP

**Full DHCP + TFTP + DNS:**
```bash
docker run -d --name dnsmasq \
  --cap-add=NET_ADMIN \
  --net=host \
  quay.io/poseidon/dnsmasq:latest \
  -d -q \
  --dhcp-range=192.168.1.3,192.168.1.254,30m \
  --enable-tftp \
  --tftp-root=/var/lib/tftpboot \
  --dhcp-match=set:bios,option:client-arch,0 \
  --dhcp-boot=tag:bios,undionly.kpxe \
  --dhcp-match=set:efi64,option:client-arch,9 \
  --dhcp-boot=tag:efi64,ipxe.efi \
  --dhcp-userclass=set:ipxe,iPXE \
  --dhcp-boot=tag:ipxe,http://matchbox.example.com:8080/boot.ipxe \
  --address=/matchbox.example.com/192.168.1.2 \
  --log-queries \
  --log-dhcp
```

**Proxy DHCP (alongside existing DHCP):**
```bash
docker run -d --name dnsmasq \
  --cap-add=NET_ADMIN \
  --net=host \
  quay.io/poseidon/dnsmasq:latest \
  -d -q \
  --dhcp-range=192.168.1.1,proxy,255.255.255.0 \
  --enable-tftp \
  --tftp-root=/var/lib/tftpboot \
  --dhcp-userclass=set:ipxe,iPXE \
  --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe \
  --pxe-service=tag:ipxe,x86PC,"iPXE",http://matchbox.example.com:8080/boot.ipxe \
  --log-queries \
  --log-dhcp
```

**Included files:** `undionly.kpxe`, `ipxe.efi`, `grub.efi` (bundled in image)

### Option 2: Existing DHCP/TFTP Infrastructure

**Use case:** Enterprise environments with network admin policies

**Required DHCP options (ISC DHCP example):**
```
subnet 192.168.1.0 netmask 255.255.255.0 {
  range 192.168.1.10 192.168.1.250;
  
  # BIOS clients
  if option architecture-type = 00:00 {
    filename "undionly.kpxe";
  }
  # UEFI clients
  elsif option architecture-type = 00:09 {
    filename "ipxe.efi";
  }
  # iPXE clients
  elsif exists user-class and option user-class = "iPXE" {
    filename "http://matchbox.example.com:8080/boot.ipxe";
  }
  
  next-server 192.168.1.100;  # TFTP server IP
}
```

**TFTP files (place in tftp root):**
- Download from http://boot.ipxe.org/undionly.kpxe
- Download from http://boot.ipxe.org/ipxe.efi

### Option 3: iPXE-only (No PXE Chainload)

**Use case:** Modern hardware with native iPXE firmware

**DHCP config (simpler):**
```
filename "http://matchbox.example.com:8080/boot.ipxe";
```

**No TFTP server needed** (iPXE fetches directly via HTTP)

**Limitation:** Doesn't support legacy BIOS with basic PXE ROM

## TLS Certificate Setup

gRPC API requires TLS client certificates for authentication.

### Option 1: Provided cert-gen Script

```bash
cd scripts/tls
export SAN=DNS.1:matchbox.example.com,IP.1:192.168.1.100
./cert-gen
```

**Generates:**
- `ca.crt` - Self-signed CA
- `server.crt`, `server.key` - Server credentials
- `client.crt`, `client.key` - Client credentials (for Terraform)

**Install server certs:**
```bash
sudo mkdir -p /etc/matchbox
sudo cp ca.crt server.crt server.key /etc/matchbox/
sudo chown -R matchbox:matchbox /etc/matchbox
```

**Save client certs for Terraform:**
```bash
mkdir -p ~/.matchbox
cp client.crt client.key ca.crt ~/.matchbox/
```

### Option 2: Corporate PKI

**Preferred for production:** Use organization's certificate authority

**Requirements:**
- Server cert with SAN: `DNS:matchbox.example.com`
- Client cert issued by same CA
- CA cert for validation

**Matchbox flags:**
```
-ca-file=/etc/matchbox/ca.crt
-cert-file=/etc/matchbox/server.crt
-key-file=/etc/matchbox/server.key
```

**Terraform provider config:**
```hcl
provider "matchbox" {
  endpoint    = "matchbox.example.com:8081"
  client_cert = file("/path/to/client.crt")
  client_key  = file("/path/to/client.key")
  ca          = file("/path/to/ca.crt")
}
```

### Option 3: Let's Encrypt (HTTP API only)

**Note:** gRPC requires client cert auth (incompatible with Let's Encrypt)

**Use case:** TLS for HTTP endpoints only (read-only API)

**Matchbox flags:**
```
-web-ssl=true
-web-cert-file=/etc/letsencrypt/live/matchbox.example.com/fullchain.pem
-web-key-file=/etc/letsencrypt/live/matchbox.example.com/privkey.pem
```

**Limitation:** Still need self-signed certs for gRPC API

## Configuration Flags

### Core Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-address` | `127.0.0.1:8080` | HTTP API listen address |
| `-rpc-address` | `` | gRPC API listen address (empty = disabled) |
| `-data-path` | `/var/lib/matchbox` | Data directory (FileStore) |
| `-assets-path` | `/var/lib/matchbox/assets` | Static assets directory |
| `-log-level` | `info` | Logging level (debug, info, warn, error) |

### TLS Flags (gRPC)

| Flag | Default | Description |
|------|---------|-------------|
| `-ca-file` | `/etc/matchbox/ca.crt` | CA certificate for client verification |
| `-cert-file` | `/etc/matchbox/server.crt` | Server TLS certificate |
| `-key-file` | `/etc/matchbox/server.key` | Server TLS private key |

### TLS Flags (HTTP, optional)

| Flag | Default | Description |
|------|---------|-------------|
| `-web-ssl` | `false` | Enable TLS for HTTP API |
| `-web-cert-file` | `` | HTTP server TLS certificate |
| `-web-key-file` | `` | HTTP server TLS private key |

### Environment Variables

All flags can be set via environment variables with `MATCHBOX_` prefix:

```bash
export MATCHBOX_ADDRESS=0.0.0.0:8080
export MATCHBOX_RPC_ADDRESS=0.0.0.0:8081
export MATCHBOX_LOG_LEVEL=debug
export MATCHBOX_DATA_PATH=/custom/path
```

## Operational Considerations

### Firewall Configuration

**Matchbox host:**
```bash
firewall-cmd --permanent --add-port=8080/tcp  # HTTP API
firewall-cmd --permanent --add-port=8081/tcp  # gRPC API
firewall-cmd --reload
```

**dnsmasq host (if separate):**
```bash
firewall-cmd --permanent --add-service=dhcp
firewall-cmd --permanent --add-service=tftp
firewall-cmd --permanent --add-service=dns  # optional
firewall-cmd --reload
```

### Monitoring

**Health check endpoints:**
```bash
# HTTP API
curl http://matchbox.example.com:8080
# Should return: matchbox

# gRPC API
openssl s_client -connect matchbox.example.com:8081 \
  -CAfile ~/.matchbox/ca.crt \
  -cert ~/.matchbox/client.crt \
  -key ~/.matchbox/client.key
```

**Prometheus metrics:** Not built-in; consider adding reverse proxy (e.g., nginx) with metrics exporter

**Logs (systemd):**
```bash
journalctl -u matchbox -f
```

**Logs (container):**
```bash
docker logs -f matchbox
```

### Backup Strategy

**What to backup:**
1. `/var/lib/matchbox/{groups,profiles,ignition}` - Configs
2. `/etc/matchbox/*.{crt,key}` - TLS certificates
3. Terraform state (if using Terraform provider)

**Backup command:**
```bash
tar -czf matchbox-backup-$(date +%F).tar.gz \
  /var/lib/matchbox/{groups,profiles,ignition} \
  /etc/matchbox
```

**Restore:**
```bash
tar -xzf matchbox-backup-YYYY-MM-DD.tar.gz -C /
sudo chown -R matchbox:matchbox /var/lib/matchbox
sudo systemctl restart matchbox
```

**GitOps approach:** Store configs in git repository for versioning and auditability

### Updates

**Binary deployment:**
```bash
# Download new version
wget https://github.com/poseidon/matchbox/releases/download/vX.Y.Z/matchbox-vX.Y.Z-linux-amd64.tar.gz
tar xzf matchbox-vX.Y.Z-linux-amd64.tar.gz

# Replace binary
sudo systemctl stop matchbox
sudo cp matchbox-vX.Y.Z-linux-amd64/matchbox /usr/local/bin/
sudo systemctl start matchbox
```

**Container deployment:**
```bash
docker pull quay.io/poseidon/matchbox:vX.Y.Z
docker stop matchbox
docker rm matchbox
docker run -d --name matchbox ... quay.io/poseidon/matchbox:vX.Y.Z ...
```

**Kubernetes deployment:**
```bash
kubectl set image deployment/matchbox matchbox=quay.io/poseidon/matchbox:vX.Y.Z
kubectl rollout status deployment/matchbox
```

### Scaling Considerations

**Vertical scaling (single instance):**
- CPU: Minimal (config rendering is lightweight)
- Memory: ~50MB base + asset cache
- Disk: Depends on cached assets (100MB - 10GB+)

**Horizontal scaling (multiple instances):**
- Stateless HTTP API (load balance round-robin)
- Shared storage required (RWX PV, NFS, or custom backend)
- gRPC API can be load-balanced with gRPC-aware LB

**Asset serving optimization:**
- Use CDN or cache proxy for remote assets
- Local asset caching for <100 machines
- Dedicated HTTP server (nginx) for large deployments (1000+ machines)

### Security Best Practices

1. **Don't store secrets in Ignition configs**
   - Use Ignition `files.source` with auth headers to fetch from Vault
   - Or provision minimal config, fetch secrets post-boot

2. **Network segmentation**
   - Provision VLAN isolated from production
   - Firewall rules: only allow provisioning traffic

3. **gRPC API access control**
   - Client cert authentication (mandatory)
   - Restrict cert issuance to authorized personnel/systems
   - Rotate certs periodically

4. **Audit logging**
   - Version control groups/profiles (git)
   - Log gRPC API changes (Terraform state tracking)
   - Monitor HTTP endpoint access

5. **Validate configs before deployment**
   - `fcct --strict` for Butane configs
   - Terraform plan before apply
   - Test in dev environment first

## Troubleshooting

### Common Issues

**1. Machines not PXE booting:**
```bash
# Check DHCP responses
tcpdump -i eth0 port 67 and port 68

# Verify TFTP files
ls -la /var/lib/tftpboot/
curl tftp://192.168.1.100/undionly.kpxe

# Check Matchbox accessibility
curl http://matchbox.example.com:8080/boot.ipxe
```

**2. 404 Not Found on /ignition:**
```bash
# Test group matching
curl 'http://matchbox.example.com:8080/ignition?mac=52:54:00:89:d8:10'

# Check group exists
ls -la /var/lib/matchbox/groups/

# Check profile referenced by group exists
ls -la /var/lib/matchbox/profiles/

# Verify ignition_id file exists
ls -la /var/lib/matchbox/ignition/
```

**3. gRPC connection refused (Terraform):**
```bash
# Test TLS connection
openssl s_client -connect matchbox.example.com:8081 \
  -CAfile ~/.matchbox/ca.crt \
  -cert ~/.matchbox/client.crt \
  -key ~/.matchbox/client.key

# Check Matchbox gRPC is listening
sudo ss -tlnp | grep 8081

# Verify firewall
sudo firewall-cmd --list-ports
```

**4. Ignition config validation errors:**
```bash
# Validate Butane locally
podman run -i --rm quay.io/coreos/fcct:release --strict < config.yaml

# Fetch rendered Ignition
curl 'http://matchbox.example.com:8080/ignition?mac=...' | jq .

# Validate Ignition spec
curl 'http://matchbox.example.com:8080/ignition?mac=...' | \
  podman run -i --rm quay.io/coreos/ignition-validate:latest
```

## Summary

Matchbox deployment considerations:

- **Architecture:** Single-host (dev/lab) vs HA (production) vs Kubernetes
- **Installation:** Binary (systemd), container (Docker/Podman), or Kubernetes manifests
- **Network boot:** dnsmasq container (quick), existing infrastructure (enterprise), or iPXE-only (modern)
- **TLS:** Self-signed (dev), corporate PKI (production), Let's Encrypt (HTTP only)
- **Scaling:** Vertical (simple) vs horizontal (requires shared storage)
- **Security:** Client cert auth, network segmentation, no secrets in configs
- **Operations:** Backup configs, GitOps workflow, monitoring/logging

**Recommendation for production:**
- HA deployment (2+ instances) with load balancer
- Shared storage (NFS or RWX PV on Kubernetes)
- Corporate PKI for TLS certificates
- GitOps workflow (Terraform + git-controlled configs)
- Network segmentation (dedicated provisioning VLAN)
- Prometheus/Grafana monitoring
