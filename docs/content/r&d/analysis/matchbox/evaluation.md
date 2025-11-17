---
title: "Use Case Evaluation"
type: docs
weight: 5
description: "Evaluation of Matchbox for specific use cases and comparison with alternatives"
---

# Matchbox Use Case Evaluation

Analysis of Matchbox's suitability for various use cases, strengths, limitations, and comparison with alternative provisioning solutions.

## Use Case Fit Analysis

### ✅ Ideal Use Cases

#### 1. Bare-Metal Kubernetes Clusters

**Scenario:** Provisioning 10-1000 physical servers for Kubernetes nodes

**Why Matchbox Excels:**
- Ignition-native (perfect for Fedora CoreOS/Flatcar)
- Declarative machine provisioning via Terraform
- Label-based matching (region, role, hardware type)
- Integration with Typhoon Kubernetes distribution
- Minimal OS surface (immutable, container-optimized)

**Example workflow:**
```hcl
resource "matchbox_profile" "k8s_controller" {
  name   = "k8s-controller"
  kernel = "/assets/fedora-coreos/.../kernel"
  raw_ignition = data.ct_config.controller.rendered
}

resource "matchbox_group" "controllers" {
  profile = matchbox_profile.k8s_controller.name
  selector = {
    role = "controller"
  }
}
```

**Alternatives considered:**
- **Cloud-init + netboot.xyz**: Less declarative, no native Ignition support
- **Foreman**: Heavier, more complex for container-centric workloads
- **Metal³**: Kubernetes-native but requires existing cluster

**Verdict:** ⭐⭐⭐⭐⭐ Matchbox is purpose-built for this

---

#### 2. Lab/Development Environments

**Scenario:** Rapid PXE boot testing with QEMU/KVM VMs or homelab servers

**Why Matchbox Excels:**
- Quick setup (binary + dnsmasq container)
- No DHCP infrastructure required (proxy-DHCP mode)
- Localhost deployment (no external dependencies)
- Fast iteration (change configs, re-PXE)
- Included examples and scripts

**Example setup:**
```bash
# Start Matchbox locally
docker run -d --net=host -v /var/lib/matchbox:/var/lib/matchbox \
  quay.io/poseidon/matchbox:latest -address=0.0.0.0:8080

# Start dnsmasq on same host
docker run -d --net=host --cap-add=NET_ADMIN \
  quay.io/poseidon/dnsmasq ...
```

**Alternatives considered:**
- **netboot.xyz**: Great for manual OS selection, no automation
- **PiXE server**: Simpler but less flexible matching logic
- **Manual iPXE scripts**: No dynamic matching, manual maintenance

**Verdict:** ⭐⭐⭐⭐⭐ Minimal setup, maximum flexibility

---

#### 3. Edge/Remote Site Provisioning

**Scenario:** Provision machines at 10+ remote datacenters or edge locations

**Why Matchbox Excels:**
- Lightweight (single binary, ~20MB)
- Declarative region-based matching
- Centralized config management (Terraform)
- Can run on minimal hardware (ARM support)
- HTTP-based (works over WAN with reverse proxy)

**Architecture:**
```
Central Matchbox (via Terraform)
  ↓ gRPC API
Regional Matchbox Instances (read-only cache)
  ↓ HTTP
Edge Machines (PXE boot)
```

**Label-based routing:**
```json
{
  "selector": {
    "region": "us-west",
    "site": "pdx-1"
  },
  "metadata": {
    "ntp_servers": ["10.100.1.1", "10.100.1.2"]
  }
}
```

**Alternatives considered:**
- **Foreman**: Requires more resources per site
- **Ansible + netboot**: No declarative PXE boot, post-install only
- **Cloud-init datasources**: Requires cloud metadata service per site

**Verdict:** ⭐⭐⭐⭐☆ Good fit, but consider caching strategy for WAN

---

### ⚠️ Moderate Fit Use Cases

#### 4. Multi-Tenant Bare-Metal Cloud

**Scenario:** Provide bare-metal-as-a-service to multiple customers

**Matchbox challenges:**
- No built-in multi-tenancy (single namespace)
- No RBAC (gRPC API is all-or-nothing with client certs)
- No customer self-service portal

**Workarounds:**
- Deploy separate Matchbox per tenant (isolation via separate instances)
- Proxy gRPC API with custom RBAC layer
- Use group selectors with customer IDs

**Better alternatives:**
- **Metal³** (Kubernetes-native, better multi-tenancy)
- **OpenStack Ironic** (purpose-built for bare-metal cloud)
- **MAAS** (Ubuntu-specific, has RBAC)

**Verdict:** ⭐⭐☆☆☆ Possible but architecturally challenging

---

#### 5. Heterogeneous OS Provisioning

**Scenario:** Need to provision Fedora CoreOS, Ubuntu, RHEL, Windows

**Matchbox challenges:**
- Designed for Ignition-based OSes (FCOS, Flatcar, RHCOS)
- No native support for Kickstart (RHEL/CentOS)
- No support for Preseed (Ubuntu/Debian)
- No Windows unattend.xml support

**What works:**
- Fedora CoreOS ✅
- Flatcar Linux ✅
- RHEL CoreOS ✅
- Container Linux (deprecated but supported) ✅

**What requires workarounds:**
- RHEL/CentOS: Possible via generic configs + Kickstart URLs, but not native
- Ubuntu: Can PXE boot and point to autoinstall ISO, but loses Matchbox templating benefits
- Debian: Similar to Ubuntu
- Windows: Not supported (different PXE boot mechanisms)

**Better alternatives for heterogeneous environments:**
- **Foreman** (supports Kickstart, Preseed, unattend.xml)
- **MAAS** (Ubuntu-centric but extensible)
- **Cobbler** (older but supports many OS types)

**Verdict:** ⭐⭐☆☆☆ Stick to Ignition-based OSes or use different tool

---

### ❌ Poor Fit Use Cases

#### 6. Windows PXE Boot

**Why Matchbox doesn't fit:**
- No WinPE support
- No unattend.xml rendering
- Different PXE boot chain (WDS/SCCM model)

**Recommendation:** Use Microsoft WDS or SCCM

**Verdict:** ⭐☆☆☆☆ Not designed for this

---

#### 7. BIOS/Firmware Updates

**Why Matchbox doesn't fit:**
- Focused on OS provisioning, not firmware
- No vendor-specific tooling (Dell iDRAC, HP iLO integration)

**Recommendation:** Use vendor tools or Ansible with ipmi/redfish modules

**Verdict:** ⭐☆☆☆☆ Out of scope

---

## Strengths

### 1. Ignition-First Design
- Native support for modern immutable OSes
- Declarative, atomic provisioning (no config drift)
- First-boot partition/filesystem setup

### 2. Label-Based Matching
- Flexible machine classification (MAC, UUID, region, role, custom)
- Most-specific-match algorithm (override defaults per machine)
- Query params for dynamic attributes

### 3. Terraform Integration
- Declarative infrastructure as code
- Plan before apply (preview changes)
- State tracking for auditability
- Rich templating (ct_config provider for Butane)

### 4. Minimal Dependencies
- Single static binary (~20MB)
- No database required (FileStore default)
- No built-in DHCP/TFTP (separation of concerns)
- Container-ready (OCI image available)

### 5. HTTP-Centric
- Faster downloads than TFTP (iPXE via HTTP)
- Proxy/CDN friendly for asset distribution
- Standard web tooling (curl, load balancers, Ingress)

### 6. Production-Ready
- Used by Typhoon Kubernetes (battle-tested)
- Clear upgrade path (SemVer releases)
- OpenPGP signature support for config integrity

## Limitations

### 1. No Multi-Tenancy
- Single namespace (all groups/profiles global)
- No RBAC on gRPC API (client cert = full access)
- Requires separate instances per tenant

### 2. Ignition-Only Focus
- Cloud-Config deprecated (legacy support only)
- No native Kickstart/Preseed/unattend.xml
- Limits OS choice to CoreOS family

### 3. Storage Constraints
- FileStore doesn't scale to 10,000+ profiles
- No built-in HA storage (requires NFS or custom backend)
- Kubernetes deployment needs RWX PersistentVolume

### 4. No Machine Discovery
- Doesn't detect new machines (passive service)
- No inventory management (use external CMDB)
- No hardware introspection (use Ironic for that)

### 5. Limited Observability
- No built-in metrics (Prometheus integration requires reverse proxy)
- Logs are minimal (request logging only)
- No audit trail for gRPC API changes (use Terraform state)

### 6. TFTP Still Required
- Legacy BIOS PXE needs TFTP for chainloading to iPXE
- Can't fully eliminate TFTP unless all machines have native iPXE

## Comparison with Alternatives

### vs. Foreman

| Feature | Matchbox | Foreman |
|---------|----------|---------|
| **OS Support** | Ignition-based | Kickstart, Preseed, AutoYaST, etc. |
| **Complexity** | Low (single binary) | High (Rails app, DB, Puppet/Ansible) |
| **Config Model** | Declarative (Ignition) | Imperative (post-install scripts) |
| **API** | HTTP + gRPC | REST API |
| **UI** | None (API-only) | Full web UI |
| **Terraform** | Native provider | Community modules |
| **Use Case** | Container-centric infra | Traditional Linux servers |

**When to choose Matchbox:** CoreOS-based Kubernetes clusters, minimal infrastructure  
**When to choose Foreman:** Heterogeneous OS, need web UI, traditional config mgmt

---

### vs. Metal³

| Feature | Matchbox | Metal³ |
|---------|----------|--------|
| **Platform** | Standalone | Kubernetes-native (operator) |
| **Bootstrap** | Can bootstrap k8s cluster | Needs existing k8s cluster |
| **Machine Lifecycle** | Provision only | Provision + decommission + reprovision |
| **Hardware Introspection** | No (labels passed manually) | Yes (via Ironic) |
| **Multi-tenancy** | No | Yes (via k8s namespaces) |
| **Complexity** | Low | High (requires Ironic, DHCP, etc.) |

**When to choose Matchbox:** Greenfield bare-metal, no existing k8s  
**When to choose Metal³:** Existing k8s, need hardware mgmt lifecycle

---

### vs. Cobbler

| Feature | Matchbox | Cobbler |
|---------|----------|---------|
| **Age** | Modern (2016+) | Legacy (2008+) |
| **Config Format** | Ignition (declarative) | Kickstart/Preseed (imperative) |
| **Templating** | Go templates (minimal) | Cheetah templates (extensive) |
| **Python** | Go (static binary) | Python (requires interpreter) |
| **DHCP Management** | External | Can manage DHCP |
| **Maintenance** | Active (Poseidon) | Low activity |

**When to choose Matchbox:** Modern immutable OSes, container workloads  
**When to choose Cobbler:** Legacy infra, need DHCP management, heterogeneous OS

---

### vs. MAAS (Ubuntu)

| Feature | Matchbox | MAAS |
|---------|----------|------|
| **OS Support** | CoreOS family | Ubuntu (primary), others (limited) |
| **IPAM** | No (external DHCP) | Built-in IPAM |
| **Power Mgmt** | No (manual or scripts) | Built-in (IPMI, AMT, etc.) |
| **UI** | No | Full web UI |
| **Declarative** | Yes (Terraform) | Limited (CLI mostly) |
| **Cloud Integration** | No | Yes (libvirt, LXD, VM hosts) |

**When to choose Matchbox:** Non-Ubuntu, Kubernetes, minimal dependencies  
**When to choose MAAS:** Ubuntu-centric, need power mgmt, cloud integration

---

### vs. netboot.xyz

| Feature | Matchbox | netboot.xyz |
|---------|----------|-------------|
| **Purpose** | Automated provisioning | Manual OS selection menu |
| **Automation** | Full (API-driven) | None (interactive menu) |
| **Customization** | Per-machine configs | Global menu |
| **Ignition** | Native support | No |
| **Complexity** | Medium | Very low |

**When to choose Matchbox:** Automated fleet provisioning  
**When to choose netboot.xyz:** Ad-hoc OS installation, homelab

---

## Decision Matrix

Use this table to evaluate Matchbox for your use case:

| Requirement | Weight | Matchbox Score | Notes |
|-------------|--------|----------------|-------|
| **Ignition/CoreOS support** | High | ⭐⭐⭐⭐⭐ | Native, first-class |
| **Heterogeneous OS** | High | ⭐⭐☆☆☆ | Limited to Ignition OSes |
| **Declarative provisioning** | Medium | ⭐⭐⭐⭐⭐ | Terraform native |
| **Multi-tenancy** | Medium | ⭐☆☆☆☆ | Requires separate instances |
| **Web UI** | Medium | ☆☆☆☆☆ | No UI (API-only) |
| **Ease of deployment** | Medium | ⭐⭐⭐⭐☆ | Binary or container, minimal deps |
| **Scalability** | Medium | ⭐⭐⭐☆☆ | FileStore limits, need shared storage for HA |
| **Hardware mgmt** | Low | ☆☆☆☆☆ | No power mgmt, no introspection |
| **Cost** | Low | ⭐⭐⭐⭐⭐ | Open source, Apache 2.0 |

**Scoring:**
- ⭐⭐⭐⭐⭐ Excellent
- ⭐⭐⭐⭐☆ Good
- ⭐⭐⭐☆☆ Adequate
- ⭐⭐☆☆☆ Limited
- ⭐☆☆☆☆ Poor
- ☆☆☆☆☆ Not supported

## Recommendations

### Choose Matchbox if:
1. ✅ Provisioning Fedora CoreOS, Flatcar, or RHEL CoreOS
2. ✅ Building bare-metal Kubernetes clusters
3. ✅ Prefer declarative infrastructure (Terraform)
4. ✅ Want minimal dependencies (single binary)
5. ✅ Need flexible label-based machine matching
6. ✅ Have homogeneous OS requirements (all Ignition-based)

### Avoid Matchbox if:
1. ❌ Need multi-OS support (Windows, traditional Linux)
2. ❌ Require web UI for operations teams
3. ❌ Need built-in hardware management (power, BIOS config)
4. ❌ Have strict multi-tenancy requirements
5. ❌ Need automated hardware discovery/introspection

### Hybrid Approaches

**Pattern 1: Matchbox + Ansible**
- Matchbox: Initial OS provisioning
- Ansible: Post-boot configuration, app deployment
- Works well for stateful services on bare-metal

**Pattern 2: Matchbox + Metal³**
- Matchbox: Bootstrap initial k8s cluster
- Metal³: Ongoing cluster node lifecycle management
- Gradual migration from Matchbox to Metal³

**Pattern 3: Matchbox + Terraform + External Secrets**
- Matchbox: Base OS + minimal config
- Ignition: Fetch secrets from Vault/GCP Secret Manager
- Terraform: Orchestrate end-to-end provisioning

## Conclusion

Matchbox is a **purpose-built, minimalist network boot service** optimized for modern immutable operating systems (Ignition-based). It excels in container-centric bare-metal environments, particularly for Kubernetes clusters built with Fedora CoreOS or Flatcar Linux.

**Best fit:** Organizations adopting immutable infrastructure patterns, container orchestration, and declarative provisioning workflows.

**Not ideal for:** Heterogeneous OS environments, multi-tenant bare-metal clouds, or teams requiring extensive web UI and built-in hardware management.

For home labs and development, Matchbox offers an excellent balance of simplicity and power. For production Kubernetes deployments, it's a proven, battle-tested solution (via Typhoon). For complex enterprise provisioning with mixed OS requirements, consider Foreman or MAAS instead.
