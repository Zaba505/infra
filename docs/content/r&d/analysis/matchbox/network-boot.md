---
title: "Network Boot Support"
type: docs
weight: 2
description: "Detailed analysis of Matchbox's network boot capabilities"
---

# Network Boot Support in Matchbox

Matchbox provides comprehensive network boot support for bare-metal provisioning, supporting multiple boot firmware types and protocols.

## Overview

Matchbox serves as an HTTP entrypoint for network-booted machines but **does not implement DHCP, TFTP, or DNS services itself**. Instead, it integrates with existing network infrastructure (or companion services like dnsmasq) to provide a complete PXE boot solution.

## Boot Protocol Support

### 1. PXE (Preboot Execution Environment)

**Legacy BIOS support via chainloading to iPXE:**

```
Machine BIOS → DHCP (gets TFTP server) → TFTP (gets undionly.kpxe) 
→ iPXE firmware → HTTP (Matchbox /boot.ipxe)
```

**Key characteristics:**
- Requires TFTP server to serve `undionly.kpxe` (iPXE bootloader)
- Chainloads from legacy PXE ROM to modern iPXE
- Supports older hardware with basic PXE firmware
- TFTP only used for initial iPXE bootstrap; subsequent downloads via HTTP

### 2. iPXE (Enhanced PXE)

**Primary boot method supported by Matchbox:**

```
iPXE Client → DHCP (gets boot script URL) → HTTP (Matchbox endpoints)
→ Kernel/initrd download → Boot with Ignition config
```

**Endpoints served by Matchbox:**

| Endpoint | Purpose |
|----------|---------|
| `/boot.ipxe` | Static script that gathers machine attributes (UUID, MAC, hostname, serial) |
| `/ipxe?<labels>` | Rendered iPXE script with kernel, initrd, and boot args for matched machine |
| `/assets/` | Optional local caching of kernel/initrd images |

**Example iPXE flow:**

1. Machine boots with iPXE firmware
2. DHCP response points to `http://matchbox.example.com:8080/boot.ipxe`
3. iPXE fetches `/boot.ipxe`:
   ```
   #!ipxe
   chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
   ```
4. iPXE makes request to `/ipxe?uuid=...&mac=...` with machine attributes
5. Matchbox matches machine to group/profile and renders iPXE script:
   ```
   #!ipxe
   kernel /assets/coreos/VERSION/coreos_production_pxe.vmlinuz \
     coreos.config.url=http://matchbox.foo:8080/ignition?uuid=${uuid}&mac=${mac:hexhyp} \
     coreos.first_boot=1
   initrd /assets/coreos/VERSION/coreos_production_pxe_image.cpio.gz
   boot
   ```

**Advantages:**
- HTTP downloads (faster than TFTP)
- Scriptable boot logic
- Can fetch configs from HTTP endpoints
- Supports HTTPS (if compiled with TLS support)

### 3. GRUB2

**UEFI firmware support:**

```
UEFI Firmware → DHCP (gets GRUB bootloader) → TFTP (grub.efi)
→ GRUB → HTTP (Matchbox /grub endpoint)
```

**Matchbox endpoint:** `/grub?<labels>`

**Example GRUB config rendered by Matchbox:**
```
default=0
timeout=1
menuentry "CoreOS" {
  echo "Loading kernel"
  linuxefi "(http;matchbox.foo:8080)/assets/coreos/VERSION/coreos_production_pxe.vmlinuz" \
    "coreos.config.url=http://matchbox.foo:8080/ignition" "coreos.first_boot"
  echo "Loading initrd"
  initrdefi "(http;matchbox.foo:8080)/assets/coreos/VERSION/coreos_production_pxe_image.cpio.gz"
}
```

**Use case:**
- UEFI systems that prefer GRUB over iPXE
- Environments with existing GRUB network boot infrastructure

### 4. PXELINUX (Legacy, via TFTP)

While not a primary Matchbox target, PXELINUX clients can be configured to chainload iPXE:

```
# /var/lib/tftpboot/pxelinux.cfg/default
timeout 10
default iPXE
LABEL iPXE
KERNEL ipxe.lkrn
APPEND dhcp && chain http://matchbox.example.com:8080/boot.ipxe
```

## DHCP Configuration Patterns

Matchbox supports two DHCP deployment models:

### Pattern 1: PXE-Enabled DHCP

Full DHCP server provides IP allocation + PXE boot options.

**Example dnsmasq configuration:**

```ini
dhcp-range=192.168.1.1,192.168.1.254,30m
enable-tftp
tftp-root=/var/lib/tftpboot

# Legacy BIOS → chainload to iPXE
dhcp-match=set:bios,option:client-arch,0
dhcp-boot=tag:bios,undionly.kpxe

# UEFI → iPXE
dhcp-match=set:efi32,option:client-arch,6
dhcp-boot=tag:efi32,ipxe.efi
dhcp-match=set:efi64,option:client-arch,9
dhcp-boot=tag:efi64,ipxe.efi

# iPXE clients → Matchbox
dhcp-userclass=set:ipxe,iPXE
dhcp-boot=tag:ipxe,http://matchbox.example.com:8080/boot.ipxe

# DNS for Matchbox
address=/matchbox.example.com/192.168.1.100
```

**Client architecture detection:**
- Option 93 (`client-arch`): Identifies BIOS (0), UEFI32 (6), UEFI64 (9)
- User class: Detects iPXE clients to skip TFTP chainloading

### Pattern 2: Proxy DHCP

Runs alongside existing DHCP server; provides only boot options (no IP allocation).

**Example dnsmasq proxy-DHCP:**

```ini
dhcp-range=192.168.1.1,proxy,255.255.255.0
enable-tftp
tftp-root=/var/lib/tftpboot

# Chainload legacy PXE to iPXE
pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe
# iPXE clients → Matchbox
dhcp-userclass=set:ipxe,iPXE
pxe-service=tag:ipxe,x86PC,"iPXE",http://matchbox.example.com:8080/boot.ipxe
```

**Benefits:**
- Non-invasive: doesn't replace existing DHCP
- PXE clients receive merged responses from both DHCP servers
- Ideal for environments where main DHCP cannot be modified

## Network Boot Flow (Complete)

### Scenario: BIOS machine with legacy PXE firmware

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. Machine powers on, BIOS set to network boot                  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 2. NIC PXE firmware broadcasts DHCPDISCOVER (PXEClient)          │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 3. DHCP/proxyDHCP responds with:                                 │
│    - IP address (if full DHCP)                                   │
│    - Next-server: TFTP server IP                                 │
│    - Filename: undionly.kpxe (based on arch=0)                   │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 4. PXE firmware downloads undionly.kpxe via TFTP                 │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 5. Execute iPXE (undionly.kpxe)                                  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 6. iPXE requests DHCP again, identifies as iPXE (user-class)     │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 7. DHCP responds with boot URL (not TFTP):                       │
│    http://matchbox.example.com:8080/boot.ipxe                    │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 8. iPXE fetches /boot.ipxe via HTTP:                             │
│    #!ipxe                                                        │
│    chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&...                 │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 9. iPXE chains to /ipxe?uuid=XXX&mac=YYY (introspected labels)   │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 10. Matchbox matches machine to group/profile                    │
│     - Finds most specific group matching labels                  │
│     - Retrieves profile (kernel, initrd, args, configs)          │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 11. Matchbox renders iPXE script with:                           │
│     - kernel URL (local asset or remote HTTPS)                   │
│     - initrd URL                                                 │
│     - kernel args (including ignition.config.url)                │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 12. iPXE downloads kernel + initrd (HTTP/HTTPS)                  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 13. iPXE boots kernel with specified args                        │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 14. Fedora CoreOS/Flatcar boots, Ignition runs                   │
│     - Fetches /ignition?uuid=XXX&mac=YYY from Matchbox           │
│     - Matchbox renders Ignition config with group metadata       │
│     - Ignition partitions disk, writes files, creates users      │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│ 15. System reboots (if disk install), boots from disk            │
└──────────────────────────────────────────────────────────────────┘
```

## Asset Serving

Matchbox can serve static assets (kernel, initrd images) from a local directory to reduce bandwidth and increase speed:

**Asset directory structure:**
```
/var/lib/matchbox/assets/
├── fedora-coreos/
│   └── 36.20220906.3.2/
│       ├── fedora-coreos-36.20220906.3.2-live-kernel-x86_64
│       ├── fedora-coreos-36.20220906.3.2-live-initramfs.x86_64.img
│       └── fedora-coreos-36.20220906.3.2-live-rootfs.x86_64.img
└── flatcar/
    └── 3227.2.0/
        ├── flatcar_production_pxe.vmlinuz
        ├── flatcar_production_pxe_image.cpio.gz
        └── version.txt
```

**HTTP endpoint:** `http://matchbox.example.com:8080/assets/`

**Scripts provided:**
- `scripts/get-fedora-coreos` - Download/verify Fedora CoreOS images
- `scripts/get-flatcar` - Download/verify Flatcar Linux images

**Profile reference:**
```json
{
  "boot": {
    "kernel": "/assets/fedora-coreos/36.20220906.3.2/fedora-coreos-36.20220906.3.2-live-kernel-x86_64",
    "initrd": ["--name main /assets/fedora-coreos/36.20220906.3.2/fedora-coreos-36.20220906.3.2-live-initramfs.x86_64.img"]
  }
}
```

**Alternative:** Profiles can reference remote HTTPS URLs (requires iPXE compiled with TLS support):
```json
{
  "kernel": "https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/36.20220906.3.2/x86_64/fedora-coreos-36.20220906.3.2-live-kernel-x86_64"
}
```

## OS Support

### Fedora CoreOS

**Boot types:**
1. **Live PXE** (RAM-only, ephemeral)
2. **Install to disk** (persistent, recommended)

**Required kernel args:**
- `coreos.inst.install_dev=/dev/sda` - Target disk for install
- `coreos.inst.ignition_url=http://matchbox/ignition?uuid=${uuid}&mac=${mac:hexhyp}` - Provisioning config
- `coreos.live.rootfs_url=...` - Root filesystem image

**Ignition fetch:** During first boot, `ignition.service` fetches config from Matchbox

### Flatcar Linux

**Boot types:**
1. **Live PXE** (RAM-only)
2. **Install to disk**

**Required kernel args:**
- `flatcar.first_boot=yes` - Marks first boot
- `flatcar.config.url=http://matchbox/ignition?uuid=${uuid}&mac=${mac:hexhyp}` - Ignition config URL
- `flatcar.autologin` - Auto-login to console (optional, dev/debug)

**Ignition support:** Flatcar uses Ignition v3.x for provisioning

### RHEL CoreOS

Supported as it uses Ignition like Fedora CoreOS. Requires Red Hat-specific image sources.

## Machine Matching & Labels

Matchbox matches machines to profiles using labels extracted during boot:

### Reserved Label Selectors

| Label | Source | Example | Normalized |
|-------|--------|---------|------------|
| `uuid` | SMBIOS UUID | `550e8400-e29b-41d4-a716-446655440000` | Lowercase |
| `mac` | NIC MAC address | `52:54:00:89:d8:10` | Normalized to colons |
| `hostname` | Network boot program | `node1.example.com` | As-is |
| `serial` | Hardware serial | `VMware-42 1a...` | As-is |

### Custom Labels

Groups can match on arbitrary labels passed as query params:
```
/ipxe?mac=52:54:00:89:d8:10&region=us-west&env=prod
```

**Matching precedence:** Most specific group wins (most selector matches)

## Firmware Compatibility

| Firmware Type | Client Arch | Boot File | Protocol | Matchbox Support |
|---------------|-------------|-----------|----------|------------------|
| BIOS (legacy PXE) | 0 | `undionly.kpxe` → iPXE | TFTP → HTTP | ✅ Via chainload |
| UEFI 32-bit | 6 | `ipxe.efi` | TFTP → HTTP | ✅ |
| UEFI (BIOS compat) | 7 | `ipxe.efi` | TFTP → HTTP | ✅ |
| UEFI 64-bit | 9 | `ipxe.efi` | TFTP → HTTP | ✅ |
| Native iPXE | - | N/A | HTTP | ✅ Direct |
| GRUB (UEFI) | - | `grub.efi` | TFTP → HTTP | ✅ `/grub` endpoint |

## Network Requirements

**Firewall rules on Matchbox host:**
```bash
# HTTP API (read-only)
firewall-cmd --add-port=8080/tcp --permanent

# gRPC API (authenticated, Terraform)
firewall-cmd --add-port=8081/tcp --permanent
```

**DNS requirement:**
- `matchbox.example.com` must resolve to Matchbox server IP
- Can be configured in dnsmasq, corporate DNS, or `/etc/hosts` on DHCP server

**DHCP/TFTP host (if using dnsmasq):**
```bash
firewall-cmd --add-service=dhcp --permanent
firewall-cmd --add-service=tftp --permanent
firewall-cmd --add-service=dns --permanent  # optional
```

## Troubleshooting Tips

1. **Verify Matchbox endpoints:**
   ```bash
   curl http://matchbox.example.com:8080
   # Should return: matchbox
   
   curl http://matchbox.example.com:8080/boot.ipxe
   # Should return iPXE script
   ```

2. **Test machine matching:**
   ```bash
   curl 'http://matchbox.example.com:8080/ipxe?mac=52:54:00:89:d8:10'
   # Should return rendered iPXE script with kernel/initrd
   ```

3. **Check TFTP files:**
   ```bash
   ls -la /var/lib/tftpboot/
   # Should contain: undionly.kpxe, ipxe.efi, grub.efi
   ```

4. **Verify DHCP responses:**
   ```bash
   tcpdump -i eth0 -n port 67 and port 68
   # Watch for DHCP offers with PXE options
   ```

5. **iPXE console debugging:**
   - Press Ctrl+B during iPXE boot to enter console
   - Commands: `dhcp`, `ifstat`, `show net0/ip`, `chain http://...`

## Limitations

1. **HTTPS support:** iPXE must be compiled with crypto support (larger binary, ~80KB vs ~45KB)
2. **TFTP dependency:** Legacy PXE requires TFTP for initial chainload (can't skip)
3. **No DHCP/TFTP built-in:** Must use external services or dnsmasq container
4. **Boot firmware variations:** Some vendor PXE implementations have quirks
5. **SecureBoot:** iPXE and GRUB must be signed (or SecureBoot disabled)

## Reference Implementation: dnsmasq Container

Matchbox project provides `quay.io/poseidon/dnsmasq` with:
- Pre-configured DHCP/TFTP/DNS service
- Bundled `ipxe.efi`, `undionly.kpxe`, `grub.efi`
- Example configs for PXE-DHCP and proxy-DHCP modes

**Quick start (full DHCP):**
```bash
docker run --rm --cap-add=NET_ADMIN --net=host quay.io/poseidon/dnsmasq \
  -d -q \
  --dhcp-range=192.168.1.3,192.168.1.254 \
  --enable-tftp --tftp-root=/var/lib/tftpboot \
  --dhcp-match=set:bios,option:client-arch,0 \
  --dhcp-boot=tag:bios,undionly.kpxe \
  --dhcp-match=set:efi64,option:client-arch,9 \
  --dhcp-boot=tag:efi64,ipxe.efi \
  --dhcp-userclass=set:ipxe,iPXE \
  --dhcp-boot=tag:ipxe,http://matchbox.example.com:8080/boot.ipxe \
  --address=/matchbox.example.com/192.168.1.2 \
  --log-queries --log-dhcp
```

**Quick start (proxy-DHCP):**
```bash
docker run --rm --cap-add=NET_ADMIN --net=host quay.io/poseidon/dnsmasq \
  -d -q \
  --dhcp-range=192.168.1.1,proxy,255.255.255.0 \
  --enable-tftp --tftp-root=/var/lib/tftpboot \
  --dhcp-userclass=set:ipxe,iPXE \
  --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe \
  --pxe-service=tag:ipxe,x86PC,"iPXE",http://matchbox.example.com:8080/boot.ipxe \
  --log-queries --log-dhcp
```

## Summary

Matchbox provides robust network boot support through:
- **Protocol flexibility:** iPXE (primary), GRUB2, legacy PXE (via chainload)
- **Firmware compatibility:** BIOS and UEFI
- **Modern approach:** HTTP-based with optional local asset caching
- **Clean separation:** Matchbox handles config rendering; external services handle DHCP/TFTP
- **Production-ready:** Used by Typhoon Kubernetes distributions for bare-metal provisioning
