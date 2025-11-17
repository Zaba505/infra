---
type: docs
title: "Network Boot Capabilities"
description: >
  Comprehensive analysis of network boot support on HP ProLiant DL360 Gen9
---

## Overview

The HP ProLiant DL360 Gen9 provides robust network boot capabilities through multiple protocols and firmware interfaces. This makes it particularly well-suited for diskless deployments, automated provisioning, and infrastructure-as-code workflows.

## Supported Network Boot Protocols

### PXE (Preboot Execution Environment)

The DL360 Gen9 fully supports PXE boot via both legacy BIOS and UEFI firmware modes:

- **Legacy BIOS PXE**: Traditional PXE implementation using TFTP
  - Protocol: PXEv2 (PXE 2.1)
  - Network Stack: IPv4 only in legacy mode
  - Boot files: `pxelinux.0`, `undionly.kpxe`, or custom NBP
  - DHCP options: Standard options 66 (TFTP server) and 67 (boot filename)

- **UEFI PXE**: Modern UEFI network boot implementation
  - Protocol: PXEv2 with UEFI extensions
  - Network Stack: IPv4 and IPv6 support
  - Boot files: `bootx64.efi`, `grubx64.efi`, `shimx64.efi`
  - Architecture: x64 (EFI BC)
  - DHCP Architecture ID: 0x0007 (EFI BC) or 0x0009 (EFI x86-64)

### iPXE Support

The DL360 Gen9 can boot iPXE, enabling advanced features:

- **Chainloading**: Boot standard PXE, then chainload iPXE for enhanced capabilities
- **HTTP/HTTPS Boot**: Download kernels and images over HTTP(S) instead of TFTP
- **SAN Boot**: iSCSI and AoE (ATA over Ethernet) support
- **Scripting**: Conditional boot logic and dynamic configuration
- **Embedded Scripts**: iPXE can be compiled with embedded boot scripts

**Implementation Methods**:
1. Chainload from standard PXE: DHCP points to `undionly.kpxe` or `ipxe.efi`
2. Flash iPXE to FlexibleLOM option ROM (advanced, requires care)
3. Boot iPXE from USB, then continue network boot

### UEFI HTTP Boot

Native UEFI HTTP boot is supported on Gen9 servers with recent firmware:

- **Protocol**: RFC 7230 HTTP/1.1
- **Requirements**: 
  - UEFI firmware version 2.40 or later (check via iLO)
  - DHCP option 60 (vendor class identifier) = "HTTPClient"
  - DHCP option 67 pointing to HTTP(S) URL
- **Advantages**:
  - No TFTP server required
  - Faster transfers than TFTP
  - Support for HTTPS with certificate validation
  - Better suited for large images (kernels, initramfs)
- **Limitations**:
  - UEFI mode only (not available in legacy BIOS)
  - Requires DHCP server with HTTP URL support

### HTTP(S) Boot Configuration

For UEFI HTTP boot on DL360 Gen9:

```dhcpd
# Example ISC DHCP configuration for UEFI HTTP boot
class "httpclients" {
    match if substring(option vendor-class-identifier, 0, 10) = "HTTPClient";
}

pool {
    allow members of "httpclients";
    option vendor-class-identifier "HTTPClient";
    # Point to HTTP boot URI
    filename "http://boot.example.com/boot/efi/bootx64.efi";
}
```

## Network Interface Options

The DL360 Gen9 supports multiple network adapter configurations for boot:

### FlexibleLOM (LOM = LAN on Motherboard)

HPE FlexibleLOM slot supports:
- **HPE 366FLR**: Quad-port 1GbE (Broadcom BCM5719)
- **HPE 560FLR-SFP+**: Dual-port 10GbE (Intel X710)
- **HPE 361i**: Dual-port 1GbE (Intel I350)

All FlexibleLOM adapters support PXE and UEFI network boot. The option ROM can be configured via BIOS/UEFI settings.

### PCIe Network Adapters

Standard PCIe network cards with PXE/UEFI boot ROM support:
- Intel X520, X710 series (10GbE)
- Broadcom NetXtreme series
- Mellanox ConnectX-3/4 (with appropriate firmware)

**Boot Priority**: Configure via System ROM > Network Boot Options to select which NIC boots first.

## Firmware Configuration

### Accessing Boot Configuration

1. **RBSU (ROM-Based Setup Utility)**: Press F9 during POST
2. **iLO 4 Remote Console**: Access via network, then virtual F9
3. **UEFI System Utilities**: Modern interface for UEFI firmware settings

### Key Settings

Navigate to: **System Configuration > BIOS/Platform Configuration (RBSU) > Network Boot Options**

- **Network Boot**: Enable/Disable
- **Boot Mode**: UEFI or Legacy BIOS
- **IPv4/IPv6**: Enable protocol support
- **Boot Retry**: Number of attempts before falling back to next boot device
- **Boot Order**: Prioritize network boot in boot sequence

### Per-NIC Configuration

In RBSU > Network Options:
- **Option ROM**: Enable/Disable per adapter
- **Link Speed**: Force speed/duplex or auto-negotiate
- **VLAN**: VLAN tagging for boot (if supported by DHCP/PXE environment)
- **PXE Menu**: Enable interactive PXE menu (Ctrl+S during PXE boot)

## iLO 4 Integration

The DL360 Gen9's iLO 4 provides additional network boot features:

### Virtual Media Network Boot

- Mount ISO images remotely via iLO Virtual Media
- Boot from network-attached ISO without physical media
- Useful for OS installation or diagnostics

**Workflow**:
1. Upload ISO to HTTP/HTTPS server or use SMB/NFS share
2. iLO Remote Console > Virtual Devices > Image File CD-ROM/DVD
3. Set boot order to prioritize virtual optical drive
4. Reboot server

### Scripted Deployment via iLO

iLO 4 RESTful API allows:
- Setting one-time boot to network via API call
- Automating PXE boot for provisioning pipelines
- Integration with tools like Terraform, Ansible

Example using iLO RESTful API:
```bash
curl -k -u admin:password -X PATCH \
  https://ilo-hostname/redfish/v1/Systems/1/ \
  -d '{"Boot":{"BootSourceOverrideTarget":"Pxe","BootSourceOverrideEnabled":"Once"}}'
```

## Boot Process Flow

### Legacy BIOS PXE Boot

1. Server powers on, initializes NICs
2. NIC sends DHCPDISCOVER with PXE vendor options
3. DHCP server responds with IP, TFTP server (option 66), boot file (option 67)
4. NIC downloads NBP (Network Bootstrap Program) via TFTP
5. NBP executes (e.g., pxelinux.0 loads syslinux menu)
6. User selects boot target or automated script continues
7. Kernel and initramfs download and boot

### UEFI PXE Boot

1. UEFI firmware initializes network stack
2. UEFI PXE driver sends DHCPv4/v6 DISCOVER
3. DHCP responds with boot file (e.g., `bootx64.efi`)
4. UEFI downloads boot file via TFTP
5. UEFI loads and executes boot loader (GRUB2, systemd-boot, iPXE)
6. Boot loader may download additional files (kernel, initrd, config)
7. OS boots

### UEFI HTTP Boot

1. UEFI firmware with HTTP Boot support enabled
2. DHCP request includes "HTTPClient" vendor class
3. DHCP responds with HTTP(S) URL in option 67
4. UEFI HTTP client downloads boot file over HTTP(S)
5. Execution continues as with UEFI PXE

## Performance Considerations

### TFTP vs HTTP

- **TFTP**: Slow for large files (typical: 1-5 MB/s)
  - Use for small boot loaders only
  - Chainload to iPXE or HTTP boot for better performance
- **HTTP**: 10-100x faster depending on network and server
  - Recommended for kernels, initramfs, live OS images
  - iPXE or UEFI HTTP boot required

### Network Speed Impact

DL360 Gen9 boot performance by NIC speed:
- **1GbE**: Adequate for most PXE deployments (100-125 MB/s theoretical max)
- **10GbE**: Significant improvement for large image downloads (1-2 GB/s)
- **Bonding/Teaming**: Not typically used for boot (single NIC boots)

**Recommendation**: For production diskless nodes or frequent re-provisioning, 10GbE with HTTP boot provides best performance.

## Common Use Cases

### 1. Automated OS Provisioning

Boot into installer via PXE:
- **Kickstart** (RHEL/CentOS/Rocky)
- **Preseed** (Debian/Ubuntu)
- **Ignition** (Fedora CoreOS, Flatcar)

### 2. Diskless Boot

Boot OS entirely from network/RAM:
- **Network root**: NFS or iSCSI root filesystem
- **Overlay**: Persistent storage via network overlay
- **Stateless**: Boot identical image, no local state

### 3. Rescue and Diagnostics

Boot live environments:
- **SystemRescue**
- **Clonezilla**
- **Memtest86+**
- **Hardware diagnostics** (HPE Service Pack for ProLiant)

### 4. Kubernetes/Container Hosts

PXE boot immutable OS images:
- **Talos Linux**: API-driven, diskless k8s nodes
- **Flatcar Container Linux**: Automated updates
- **k3OS**: Lightweight k8s OS

## Troubleshooting

### PXE Boot Fails

**Symptoms**: "PXE-E51: No DHCP or proxy DHCP offers received" or timeout

**Checks**:
1. Verify NIC link light and switch port status
2. Confirm DHCP server is responding (check DHCP logs)
3. Ensure DHCP options 66 and 67 are set correctly
4. Test TFTP server accessibility (`tftp -i <server> GET <file>`)
5. Check BIOS/UEFI network boot is enabled
6. Verify boot order prioritizes network boot
7. Disable Secure Boot if using unsigned boot files

### UEFI Network Boot Not Available

**Symptoms**: Network boot option missing in UEFI boot menu

**Resolution**:
1. Enter RBSU (F9), navigate to Network Options
2. Ensure at least one NIC has "Option ROM" enabled
3. Verify Boot Mode is set to UEFI (not Legacy)
4. Update System ROM to latest version if option is missing
5. Some FlexibleLOM cards require firmware update for UEFI boot support

### HTTP Boot Fails

**Symptoms**: UEFI HTTP boot option present but fails to download

**Checks**:
1. Verify firmware version supports HTTP boot (>=2.40)
2. Ensure DHCP option 67 contains valid HTTP(S) URL
3. Test URL accessibility from another client
4. Check DNS resolution if using hostname in URL
5. For HTTPS: Verify certificate is trusted (or disable cert validation in test)

### Slow PXE Boot

**Symptoms**: Boot process takes minutes instead of seconds

**Optimizations**:
1. Switch from TFTP to HTTP (chainload iPXE or use UEFI HTTP boot)
2. Increase TFTP server block size (`tftp-hpa --blocksize 1468`)
3. Tune DHCP response times (reduce lease query delays)
4. Use local network segment for boot server (avoid WAN/VPN)
5. Enable NIC interrupt coalescing in BIOS for 10GbE

## Security Considerations

### Secure Boot

DL360 Gen9 supports UEFI Secure Boot:
- Validates signed boot loaders (shim, GRUB, kernel)
- Prevents unsigned code execution during boot
- Required for some compliance scenarios

**Configuration**: RBSU > Boot Options > Secure Boot = Enabled

**Implications for Network Boot**:
- Must use signed boot loaders (e.g., shim.efi signed by Microsoft/vendor)
- Custom kernels require signing or disabling Secure Boot
- iPXE must be signed or chainloaded from signed shim

### Network Security

**Risks**:
- PXE/TFTP is unencrypted and unauthenticated
- Attacker on network can serve malicious boot images
- DHCP spoofing can redirect to malicious boot server

**Mitigations**:
1. **Network Segmentation**: Isolate PXE boot to management VLAN
2. **DHCP Snooping**: Prevent rogue DHCP servers on switch
3. **HTTPS Boot**: Use UEFI HTTP boot with TLS and certificate validation
4. **iPXE with HTTPS**: Chainload iPXE, then use HTTPS for all downloads
5. **Signed Images**: Use Secure Boot with signed boot chain
6. **802.1X**: Require network authentication before DHCP (complex for PXE)

### iLO Security

- Change default iLO password immediately
- Use TLS for iLO web interface and API
- Restrict iLO network access (firewall, separate VLAN)
- Disable iLO Virtual Media if not needed
- Enable iLO Security Override for extra security during boot

## Firmware and Driver Resources

### Required Firmware Versions

For optimal network boot support:
- **System ROM**: v2.60 or later (latest recommended)
- **iLO 4 Firmware**: v2.80 or later
- **NIC Firmware**: Latest for specific FlexibleLOM/PCIe card

Check current versions: iLO web interface > Information > Firmware Information

### Updating Firmware

Methods:
1. **HPE Service Pack for ProLiant (SPP)**: Comprehensive update bundle
   - Boot from SPP ISO (via iLO Virtual Media or USB)
   - Runs Smart Update Manager (SUM) in Linux environment
   - Updates all firmware, drivers, system ROM automatically

2. **iLO Web Interface**: Individual component updates
   - System ROM: Administration > Firmware > Update Firmware
   - Upload .fwpkg or .bin files from HPE support site

3. **Online Flash Component**: Linux Online ROM Flash utility
   - Install `hp-firmware-*` packages
   - Run updates while OS is running (requires reboot to apply)

**Download Source**: https://support.hpe.com/connect/s/product?language=en_US&kmpmoid=1010026910 (requires HPE Passport account, free registration)

## Best Practices

1. **Use UEFI Mode**: Better security, IPv6 support, larger disk support
2. **Enable HTTP Boot**: Faster and more reliable than TFTP for large files
3. **Chainload iPXE**: Flexibility of iPXE with standard PXE infrastructure
4. **Update Firmware**: Keep System ROM and iLO current for bug fixes and features
5. **Isolate Boot Network**: Use dedicated management VLAN for PXE/provisioning
6. **Test Failover**: Configure multiple DHCP servers and boot mirrors for redundancy
7. **Document Configuration**: Record BIOS settings, DHCP config, and boot infrastructure
8. **Monitor iLO Logs**: Track boot failures and hardware issues via iLO event log

## References

- HPE ProLiant DL360 Gen9 Server User Guide
- HPE UEFI System Utilities User Guide
- iLO 4 User Guide (firmware version 2.80)
- Intel PXE Specification v2.1
- UEFI Specification v2.8 (HTTP Boot)
- iPXE Documentation: https://ipxe.org/

## Conclusion

The HP ProLiant DL360 Gen9 provides enterprise-grade network boot capabilities suitable for both traditional PXE deployments and modern UEFI HTTP boot scenarios. Its flexible configuration options, mature firmware support, and iLO integration make it an excellent platform for automated provisioning, diskless computing, and infrastructure-as-code workflows in home lab environments.

For home lab use, the recommended configuration is:
- UEFI boot mode with Secure Boot disabled (unless required)
- iPXE chainloading for flexibility and HTTP performance
- iLO 4 configured for remote management and scripted provisioning
- Latest firmware for stability and feature support
