---
type: docs
title: "Configuration Guide"
weight: 3
description: >
  Setup, optimization, and configuration recommendations for HP ProLiant DL360 Gen9 in home lab environments
---

## Initial Setup

### Hardware Assembly

1. **Install Processors**:
   - Use thermal paste (HPE thermal grease recommended)
   - Align CPU carefully with socket (LGA 2011-3)
   - Secure heatsink with proper torque (hand-tighten screws in cross pattern)
   - Install both CPUs for dual-socket configuration

2. **Install Memory**:
   - Populate channels evenly (see Memory Configuration below)
   - Seat DIMMs firmly until retention clips engage
   - Verify all DIMMs recognized in POST

3. **Install Storage**:
   - Insert drives into hot-swap caddies
   - Label drives clearly for identification
   - Configure RAID controller (see Storage Configuration below)

4. **Install Network Cards**:
   - FlexibleLOM: Slide into dedicated slot until seated
   - PCIe cards: Ensure low-profile brackets, secure with screw
   - Note MAC addresses for DHCP reservations

5. **Connect Power**:
   - Install PSUs (both for redundancy)
   - Connect power cords
   - Verify PSU LEDs indicate proper operation

6. **Initial Power-On**:
   - Press power button
   - Monitor POST on screen or via iLO remote console
   - Address any POST errors before proceeding

## iLO 4 Initial Configuration

### Physical iLO Connection

1. Connect Ethernet cable to dedicated iLO port (not FlexibleLOM)
2. Default iLO IP: Obtains via DHCP, or use temporary address via RBSU
3. Check DHCP server logs for iLO MAC and assigned IP

### First Login

1. Access iLO web interface: `https://<ilo-ip>`
2. Default credentials:
   - Username: `Administrator`
   - Password: On label on server pull-out tab (or rear label)
3. **Immediately change default password** (Administration > Access Settings)

### Essential iLO Settings

**Network Configuration** (Administration > Network):
- Set static IP or DHCP reservation
- Configure DNS servers
- Set hostname (e.g., `ilo-dl360-01`)
- Enable SNTP time sync

**Security** (Administration > Security):
- Enforce HTTPS only (disable HTTP)
- Configure SSH key authentication if using CLI
- Set strong password policy
- Enable iLO Security features

**Access** (Administration > Access Settings):
- Configure iLO username/password for automation
- Create additional user accounts (separation of duties)
- Set session timeout (default: 30 minutes)

**Date and Time** (Administration > Date and Time):
- Set NTP servers for accurate timestamps
- Configure timezone

**Licenses** (Administration > Licensing):
- Install iLO Advanced license key (required for full virtual media)
- License can be purchased or acquired from secondary market

### iLO Firmware Update

Before production use, update iLO to latest version:

1. Download latest iLO 4 firmware from HPE Support Portal
2. Administration > Firmware > Update Firmware
3. Upload `.bin` file, apply update
4. iLO will reboot automatically (system stays running)

## System ROM (BIOS/UEFI) Configuration

### Accessing RBSU

- **Local**: Press F9 during POST
- **Remote**: iLO Remote Console > Power > Momentary Press > Press F9 when prompted

### Boot Mode Selection

**System Configuration > BIOS/Platform Configuration (RBSU) > Boot Mode**:

- **UEFI Mode** (recommended for modern OS):
  - Supports GPT partitions (>2TB disks)
  - Required for Secure Boot
  - Better UEFI HTTP boot support
  - IPv6 PXE boot support

- **Legacy BIOS Mode**:
  - For older OS or compatibility
  - MBR partition tables only
  - Traditional PXE boot

**Recommendation**: Use UEFI Mode unless legacy compatibility required

### Boot Order Configuration

**System Configuration > BIOS/Platform Configuration (RBSU) > Boot Options > UEFI Boot Order**:

Recommended order for network boot deployment:
1. **Network Boot**: FlexibleLOM or PCIe NIC
2. **Internal Storage**: RAID controller or disk
3. **Virtual Media**: iLO virtual CD/DVD (for installation media)
4. **USB**: For rescue/recovery

**Enable Network Boot**:
- System Configuration > BIOS/Platform Configuration (RBSU) > Network Options > Network Boot
- Set to "Enabled"

### Performance and Power Settings

**System Configuration > BIOS/Platform Configuration (RBSU) > Power Management**:

- **Power Regulator Mode**:
  - **HP Dynamic Power Savings**: Balanced power/performance (recommended for home lab)
  - **HP Static High Performance**: Maximum performance, higher power draw
  - **HP Static Low Power**: Minimize power, reduced performance
  - **OS Control**: Let OS manage (e.g., Linux cpufreq)

- **Collaborative Power Control**: Disabled (for standalone servers)
- **Minimum Processor Idle Power Core C-State**: C6 (lower idle power)
- **Energy/Performance Bias**: Balanced Performance (or Maximum Performance for compute workloads)

**Recommendation**: Start with "Dynamic Power Savings" and adjust based on workload

### Memory Configuration

**Optimal Population** (dual-CPU configuration):

For maximum performance, populate all channels before adding second DIMM per channel:

**64GB** (8x 8GB):
- CPU1: Slots 1, 4, 7, 10 and CPU2: Slots 1, 4, 7, 10
- Result: 4 channels per CPU, 1 DIMM per channel

**128GB** (8x 16GB):
- Same as above with 16GB DIMMs

**192GB** (12x 16GB):
- CPU1: Slots 1, 4, 7, 10, 2, 5 and CPU2: Slots 1, 4, 7, 10, 2, 5
- Result: 4 channels per CPU, some with 2 DIMMs per channel

**768GB** (24x 32GB):
- All slots populated

**Check Configuration**: RBSU > System Information > Memory Information

### Processor Options

**System Configuration > BIOS/Platform Configuration (RBSU) > Processor Options**:

- **Intel Hyperthreading**: Enabled (recommended for most workloads)
  - Doubles logical cores (e.g., 12-core CPU shows as 24 cores)
  - Benefits most virtualization and multi-threaded workloads
  - Disable only for specific security compliance (e.g., some cloud providers)

- **Intel Virtualization Technology (VT-x)**: Enabled (required for hypervisors)
- **Intel VT-d (IOMMU)**: Enabled (required for PCI passthrough, SR-IOV)

- **Turbo Boost**: Enabled (allows CPU to exceed base clock)
- **Cores Enabled**: All (or reduce to lower power/heat if needed)

### Integrated Devices

**System Configuration > BIOS/Platform Configuration (RBSU) > System Options > Integrated Devices**:

- **Embedded SATA Controller**: Enabled (if using SATA drives)
- **Embedded RAID Controller**: Enabled (for Smart Array controllers)
- **SR-IOV**: Enabled (if using virtual network interfaces with VMs)

### Network Controller Options

For each NIC (FlexibleLOM, PCIe):

**System Configuration > BIOS/Platform Configuration (RBSU) > Network Options > [Adapter]**:

- **Network Boot**: Enabled (for network boot on that NIC)
- **PXE/iSCSI**: Select PXE for standard network boot
- **Link Speed**: Auto-Negotiation (recommended) or force 1G/10G
- **IPv4**: Enabled (for IPv4 PXE boot)
- **IPv6**: Enabled (if using IPv6 PXE boot)

**Boot Order**: Configure which NIC boots first if multiple are enabled

### Secure Boot Configuration

**System Configuration > BIOS/Platform Configuration (RBSU) > Boot Options > Secure Boot**:

- **Secure Boot**: Disabled (for unsigned boot loaders, custom kernels)
- **Secure Boot**: Enabled (for signed boot loaders, Windows, some Linux distros)

**Note**: If using PXE with unsigned images (e.g., custom iPXE), Secure Boot must be disabled

### Firmware Updates

Update System ROM to latest version:

1. **Via iLO**:
   - iLO web > Administration > Firmware > Update Firmware
   - Upload System ROM `.fwpkg` or `.bin` file
   - Server reboots automatically to apply

2. **Via Service Pack for ProLiant (SPP)**:
   - Download SPP ISO from HPE Support Portal
   - Mount via iLO Virtual Media
   - Boot server from SPP ISO
   - Smart Update Manager (SUM) runs in Linux environment
   - Select components to update (System ROM, iLO, controller firmware, NIC firmware)
   - Apply updates, reboot

**Recommendation**: Use SPP for comprehensive updates on initial setup, then iLO for individual component updates

## Storage Configuration

### Smart Array Controller Setup

#### Access Smart Array Configuration

- **During POST**: Press F5 when "Smart Array Configuration Utility" message appears
- **Via RBSU**: System Configuration > BIOS/Platform Configuration (RBSU) > System Options > ROM-Based Setup Utility > Smart Array Configuration

#### Create RAID Arrays

1. **Delete Existing Arrays** (if reconfiguring):
   - Select controller > Configuration > Delete Array
   - Confirm deletion (data loss warning)

2. **Create New Array**:
   - Select controller > Configuration > Create Array
   - Select physical drives to include
   - Choose RAID level:
     - **RAID 0**: Striping, no redundancy (maximum performance, maximum capacity)
     - **RAID 1**: Mirroring (redundancy, half capacity, good for boot drives)
     - **RAID 5**: Striping + parity (redundancy, n-1 capacity, balanced)
     - **RAID 6**: Striping + double parity (dual-drive failure tolerance, n-2 capacity)
     - **RAID 10**: Mirror + stripe (high performance + redundancy, half capacity)
   - Configure spare drives (hot spares for automatic rebuild)
   - Create logical drive
   - Set bootable flag if boot drive

3. **Recommended Configurations**:
   - **Boot/OS**: 2x SSD in RAID 1 (redundancy, fast boot)
   - **Data (performance)**: 4-6x SSD in RAID 10 (fast, redundant)
   - **Data (capacity)**: 4-8x HDD in RAID 6 (capacity, dual-drive tolerance)

#### Controller Settings

- **Cache Settings**:
  - **Write Cache**: Enabled (requires battery/flash-backed cache)
  - **Read Cache**: Enabled
  - **No-Battery Write Cache**: Disabled (data safety) or Enabled (performance, risk)

- **Rebuild Priority**: Medium or High (faster rebuild, may impact performance)
- **Surface Scan Delay**: 3-7 days (periodic integrity check)

### HBA Mode (Non-RAID)

For software RAID (ZFS, mdadm, Ceph):

1. Access Smart Array Configuration (F5 during POST)
2. Controller > Configuration > Enable HBA Mode
3. Confirm (RAID arrays will be deleted)
4. Reboot

**Note**: Not all Smart Array controllers support HBA mode. Check compatibility. Alternative: Use separate LSI HBA in PCIe slot.

## Network Configuration for Boot

### DHCP Server Setup

For PXE/UEFI network boot, configure DHCP server with appropriate options:

**ISC DHCP Example** (`/etc/dhcp/dhcpd.conf`):

```dhcpd
# Define subnet
subnet 192.168.10.0 netmask 255.255.255.0 {
    range 192.168.10.100 192.168.10.200;
    option routers 192.168.10.1;
    option domain-name-servers 192.168.10.1;
    
    # PXE boot options
    next-server 192.168.10.5;  # TFTP server IP
    
    # Differentiate UEFI vs BIOS
    if exists user-class and option user-class = "iPXE" {
        # iPXE boot script
        filename "http://boot.example.com/boot.ipxe";
    } elsif option arch = 00:07 or option arch = 00:09 {
        # UEFI (x86-64)
        filename "bootx64.efi";
    } else {
        # Legacy BIOS
        filename "undionly.kpxe";
    }
}

# Static reservation for DL360
host dl360-01 {
    hardware ethernet xx:xx:xx:xx:xx:xx;  # FlexibleLOM MAC
    fixed-address 192.168.10.50;
    option host-name "dl360-01";
}
```

### FlexibleLOM Configuration

Configure FlexibleLOM NIC for network boot:

1. RBSU > Network Options > FlexibleLOM
2. Enable "Network Boot"
3. Select PXE or iSCSI
4. Configure IPv4/IPv6 as needed
5. Set as first boot device in boot order

### Multi-NIC Boot Priority

If multiple NICs have network boot enabled:

1. RBSU > Network Options > Network Boot Order
2. Drag/drop to prioritize NIC boot order
3. First NIC in list attempts boot first

**Recommendation**: Enable network boot on one NIC (typically FlexibleLOM port 1) to avoid confusion

## Operating System Installation

### Traditional Installation (Virtual Media)

1. Download OS ISO (e.g., Ubuntu Server, ESXi, Proxmox)
2. Upload ISO to HTTP/HTTPS server or local file
3. iLO Remote Console > Virtual Devices > Image File CD-ROM/DVD
4. Browse to ISO location, click "Insert Media"
5. Set boot order to prioritize virtual media
6. Reboot server, boot from virtual CD/DVD
7. Proceed with OS installation

### Network Installation (PXE)

See [Network Boot Capabilities](./network-boot/) for detailed PXE/UEFI boot setup

Quick workflow:
1. Configure DHCP server with PXE options
2. Setup TFTP server with boot files
3. Enable network boot in BIOS
4. Reboot, server PXE boots
5. Select OS installer from PXE menu
6. Automated installation proceeds (Kickstart/Preseed/Ignition)

## Optimization for Specific Workloads

### Virtualization (ESXi, Proxmox, Hyper-V)

**BIOS Settings**:
- Hyperthreading: Enabled
- VT-x: Enabled
- VT-d: Enabled
- Power Management: Dynamic or OS Control
- Turbo Boost: Enabled

**Hardware**:
- Maximum memory (384GB+ recommended)
- Fast storage (SSD RAID 10 for VM storage)
- 10GbE networking for VM traffic

**Configuration**:
- Pass through NICs to VMs (SR-IOV or PCI passthrough)
- Use storage controller in HBA mode for direct disk access to VM storage (ZFS, Ceph)

### Kubernetes/Container Platforms

**BIOS Settings**:
- Hyperthreading: Enabled
- VT-x/VT-d: Enabled (for nested virtualization, kata containers)
- Power Management: Dynamic or High Performance

**Hardware**:
- 128GB+ RAM for multi-tenant workloads
- Fast local NVMe/SSD for container image cache and ephemeral storage
- 10GbE for pod networking

**OS Recommendations**:
- Talos Linux: Network-bootable, immutable k8s OS
- Flatcar Container Linux: Auto-updating, minimal OS
- Ubuntu Server: Broad compatibility, snap/docker native

### Storage Server (NAS, SAN)

**BIOS Settings**:
- Disable Hyperthreading (slight performance improvement for ZFS)
- VT-d: Enabled (if passing through HBA to VM)
- Power Management: High Performance

**Hardware**:
- Maximum drive bays (8-10 SFF)
- HBA mode or separate LSI HBA controller
- 10GbE or bonded 1GbE for network storage traffic
- ECC memory (critical for ZFS)

**Software**:
- TrueNAS SCALE (Linux-based, k8s apps)
- OpenMediaVault (Debian-based, plugins)
- Ubuntu + ZFS (custom setup)

### Compute/HPC Workloads

**BIOS Settings**:
- Hyperthreading: Depends on workload (test both)
- Turbo Boost: Enabled
- Power Management: Maximum Performance
- C-States: Disabled (reduce latency)

**Hardware**:
- High core count CPUs (E5-2680 v4, 2690 v4)
- Maximum memory bandwidth (populate all channels)
- Fast local scratch storage (NVMe)

## Monitoring and Maintenance

### iLO Health Monitoring

**Information > System Information**:
- CPU temperature and status
- Memory status
- Drive status (via controller)
- Fan speeds
- PSU status
- Overall system health LED status

**Alerting** (Administration > Alerting):
- Configure email alerts for:
  - Fan failures
  - Temperature warnings
  - Drive failures
  - Memory errors
  - PSU failures
- Set up SNMP traps for integration with monitoring systems (Nagios, Zabbix, Prometheus)

### Integrated Management Log (IML)

**Information > Integrated Management Log**:
- View hardware events and errors
- Filter by severity (Informational, Caution, Critical)
- Export log for troubleshooting

**Regular Checks**:
- Review IML weekly for early warning signs
- Address caution-level events before they become critical

### Firmware Update Cadence

**Recommendation**:
- **iLO**: Update quarterly or when security advisories released
- **System ROM**: Update annually or for bug fixes
- **Storage Controller**: Update when issues arise or annually
- **NIC Firmware**: Update when issues arise

**Method**: Use SPP for annual comprehensive updates, iLO web interface for individual component updates

### Physical Maintenance

**Monthly**:
- Check fan noise (increased noise may indicate clogged air filters or failing fan)
- Verify PSU and drive LEDs (no amber lights)
- Check iLO for alerts

**Quarterly**:
- Clean air filters (if accessible, depends on rack airflow)
- Verify backup of iLO configuration
- Test iLO Virtual Media functionality

**Annually**:
- Update all firmware via SPP
- Verify RAID battery/flash-backed cache status
- Review and update BIOS settings as workload evolves

## Troubleshooting Common Issues

### Server Won't Power On

1. Check PSU power cords connected
2. Verify PSU LEDs indicate power
3. Press iLO power button via web interface
4. Check iLO IML for power-related errors
5. Reseat PSUs, check for blown fuses

### POST Errors

**Memory Errors**:
- Reseat memory DIMMs
- Test with minimal configuration (1 DIMM per CPU)
- Replace failing DIMMs identified in POST

**CPU Errors**:
- Verify heatsink properly seated
- Check thermal paste application
- Reseat CPU (careful with pins)

**Drive Errors**:
- Check drive connection to caddy
- Verify controller recognizes drive
- Replace failing drive

### No Network Boot

See [Network Boot Troubleshooting](./network-boot/#troubleshooting) for detailed diagnostics

Quick checks:
1. Verify NIC link light
2. Confirm network boot enabled in BIOS
3. Check DHCP server logs for PXE request
4. Test TFTP server accessibility

### iLO Not Accessible

1. Check physical Ethernet connection to iLO port
2. Verify switch port active
3. Reset iLO: Press and hold iLO NMI button (rear) for 5 seconds
4. Factory reset iLO via jumper (see maintenance guide)
5. Check iLO firmware version, update if outdated

### High Fan Noise

1. Check ambient temperature (<25°C recommended)
2. Verify airflow not blocked (front/rear clearance)
3. Clean dust from intake (compressed air)
4. Check iLO temperature sensors for elevated temps
5. Lower CPU TDP if temperatures excessive (lower power CPUs)
6. Verify all fans operational (replace failed fans)

## Security Hardening

### iLO Security

1. **Change Default Credentials**: Immediately on first boot
2. **Disable Unused Services**: SSH, IPMI if not needed
3. **Use HTTPS Only**: Disable HTTP (Administration > Network > HTTP Port)
4. **Network Isolation**: Dedicated management VLAN, firewall iLO access
5. **Update Firmware**: Apply security patches promptly
6. **Account Management**: Use separate accounts, least privilege

### BIOS/UEFI Security

1. **BIOS Password**: Set administrator password (RBSU > System Options > BIOS Admin Password)
2. **Secure Boot**: Enable if using signed boot loaders
3. **Boot Order Lock**: Prevent unauthorized boot device changes
4. **TPM**: Enable if using BitLocker or LUKS disk encryption

### Operating System Security

1. **Minimal Installation**: Install only required packages
2. **Firewall**: Enable host firewall (iptables, firewalld, ufw)
3. **SSH Hardening**: Key-based auth, disable password auth, non-standard port
4. **Automatic Updates**: Enable for security patches
5. **Monitoring**: Deploy intrusion detection (fail2ban, OSSEC)

## Conclusion

Proper configuration of the HP ProLiant DL360 Gen9 ensures optimal performance, reliability, and manageability for home lab and production deployments. The combination of UEFI boot capabilities, iLO remote management, and flexible hardware configuration makes the DL360 Gen9 a versatile platform for virtualization, containerization, storage, and compute workloads.

Key takeaways:
- Update firmware early (iLO, System ROM, controllers)
- Configure iLO for remote management and monitoring
- Choose boot mode (UEFI recommended) and configure network boot appropriately
- Optimize BIOS settings for specific workload (virtualization, storage, compute)
- Implement security hardening (iLO, BIOS, OS)
- Establish monitoring and maintenance schedule

For network boot-specific configuration, refer to the [Network Boot Capabilities](./network-boot/) guide.
