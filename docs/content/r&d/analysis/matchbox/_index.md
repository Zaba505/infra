---
title: "Matchbox Analysis"
type: docs
description: "Analysis of Matchbox network boot service capabilities and architecture"
---

# Matchbox Network Boot Analysis

This section contains a comprehensive analysis of [Matchbox](https://matchbox.psdn.io/), a network boot service for provisioning bare-metal machines.

## Overview

Matchbox is an HTTP and gRPC service developed by Poseidon that automates bare-metal machine provisioning through network booting. It matches machines to configuration profiles based on hardware attributes and serves boot configurations, kernel images, and provisioning configs.

**Primary Repository**: [poseidon/matchbox](https://github.com/poseidon/matchbox)  
**Documentation**: https://matchbox.psdn.io/  
**License**: Apache 2.0  

## Key Features

- **Network Boot Support**: iPXE, PXELINUX, GRUB2 chainloading
- **OS Provisioning**: Fedora CoreOS, Flatcar Linux, RHEL CoreOS
- **Configuration Management**: Ignition v3.x configs, Butane transpilation
- **Machine Matching**: Label-based matching (MAC, UUID, hostname, serial, custom)
- **API**: Read-only HTTP API + authenticated gRPC API
- **Asset Serving**: Local caching of OS images for faster deployment
- **Templating**: Go template support for dynamic configuration

## Use Cases

1. **Bare-metal Kubernetes clusters** - Provision CoreOS nodes for k8s
2. **Lab/development environments** - Quick PXE boot for testing
3. **Datacenter provisioning** - Automate OS installation across fleets
4. **Immutable infrastructure** - Declarative machine provisioning via Terraform

## Analysis Contents

- [Network Boot Architecture](./network-boot/) - Deep dive into PXE, iPXE, and GRUB support
- [Configuration Model](./configuration/) - Profiles, Groups, and templating system
- [Deployment Patterns](./deployment/) - Installation options and operational considerations

## Quick Architecture

```
┌─────────────┐
│   Machine   │ PXE Boot
│  (BIOS/UEFI)│───┐
└─────────────┘   │
                  │
┌─────────────┐   │ DHCP/TFTP
│   dnsmasq   │◄──┘ (chainload to iPXE)
│  DHCP+TFTP  │
└─────────────┘
       │
       │ HTTP
       ▼
┌─────────────────────────┐
│      Matchbox           │
│  ┌──────────────────┐   │
│  │  HTTP Endpoints  │   │ /boot.ipxe, /ignition
│  └──────────────────┘   │
│  ┌──────────────────┐   │
│  │   gRPC API       │   │ Terraform provider
│  └──────────────────┘   │
│  ┌──────────────────┐   │
│  │ Profile/Group    │   │ Match machines
│  │   Matcher        │   │ to configs
│  └──────────────────┘   │
└─────────────────────────┘
```

## Technology Stack

- **Language**: Go
- **Config Formats**: Ignition JSON, Butane YAML
- **Boot Protocols**: PXE, iPXE, GRUB2
- **APIs**: HTTP (read-only), gRPC (authenticated)
- **Deployment**: Binary, container (Podman/Docker), Kubernetes

## Integration Points

- **Terraform**: `terraform-provider-matchbox` for declarative provisioning
- **Ignition/Butane**: CoreOS provisioning configs
- **dnsmasq**: Reference DHCP/TFTP/DNS implementation (`quay.io/poseidon/dnsmasq`)
- **Asset sources**: Can serve local or remote (HTTPS) OS images
