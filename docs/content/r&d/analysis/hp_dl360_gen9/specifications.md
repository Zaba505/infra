---
type: docs
title: "Hardware Specifications"
weight: 2
description: >
  Detailed hardware specifications and configuration options for HP ProLiant DL360 Gen9
---

## System Overview

The HP ProLiant DL360 Gen9 is a dual-socket 1U rack server designed for data center and enterprise deployments, also popular in home lab environments due to its performance and manageability.

**Generation**: Gen9 (2014-2017 product cycle)  
**Form Factor**: 1U rack-mountable (19-inch standard rack)  
**Dimensions**: 43.46 x 67.31 x 4.29 cm (17.1 x 26.5 x 1.69 in)

## Processor Support

### Supported CPU Families

The DL360 Gen9 supports Intel Xeon E5-2600 v3 and v4 series processors:

- **E5-2600 v3** (Haswell-EP): Released Q3 2014
  - Process: 22nm
  - Cores: 4-18 per socket
  - TDP: 55W-145W
  - Max Memory Speed: DDR4-2133

- **E5-2600 v4** (Broadwell-EP): Released Q1 2016
  - Process: 14nm
  - Cores: 4-22 per socket
  - TDP: 55W-145W
  - Max Memory Speed: DDR4-2400

### Popular CPU Options

**Value**: E5-2620 v3/v4 (6 cores, 15MB cache, 85W)  
**Balanced**: E5-2650 v3/v4 (10-12 cores, 25-30MB cache, 105W)  
**Performance**: E5-2680 v3/v4 (12-14 cores, 30-35MB cache, 120W)  
**High Core Count**: E5-2699 v4 (22 cores, 55MB cache, 145W)

### Configuration Options

- **Single Processor**: One CPU socket populated (budget option)
- **Dual Processor**: Both sockets populated (full performance)

**Note**: Memory and I/O performance scales with processor count. Single-CPU configuration limits memory channels and PCIe lanes.

## Memory Architecture

### Memory Specifications

- **Type**: DDR4 RDIMM or LRDIMM
- **Speed**: DDR4-2133 (v3) or DDR4-2400 (v4)
- **Slots**: 24 DIMM slots (12 per processor)
- **Maximum Capacity**: 
  - 768GB with 32GB RDIMMs
  - 1.5TB with 64GB LRDIMMs (v4 processors)
- **Minimum**: 8GB (1x 8GB DIMM)

### Memory Configuration Rules

- **Channels per CPU**: 4 channels, 3 DIMMs per channel
- **Population**: Populate channels evenly for optimal bandwidth
- **Mixing**: Do not mix RDIMM and LRDIMM types
- **Speed**: All DIMMs run at speed of slowest DIMM

### Recommended Configurations

**Basic Home Lab** (Single CPU):
- 4x 16GB = 64GB (one DIMM per channel on both memory boards)

**Standard** (Dual CPU):
- 8x 16GB = 128GB (one DIMM per channel)
- 12x 16GB = 192GB (two DIMMs per channel on primary channels)

**High Capacity** (Dual CPU):
- 24x 32GB = 768GB (all slots populated, RDIMM)

**Performance Priority**: Populate all channels before adding second DIMM per channel

## Storage Options

### Drive Bay Configurations

The DL360 Gen9 offers multiple drive bay configurations:

1. **8 SFF (2.5-inch)**: Most common configuration
2. **10 SFF**: Extended bay version
3. **4 LFF (3.5-inch)**: Less common in 1U form factor

### Drive Types Supported

- **SAS**: 12Gb/s, 6Gb/s (enterprise-grade)
- **SATA**: 6Gb/s, 3Gb/s (value option)
- **SSD**: SAS/SATA SSD, NVMe (with appropriate controller)

### Storage Controllers

**Smart Array Controllers** (HPE proprietary RAID):
- **P440ar**: Entry-level, 2GB FBWC (Flash-Backed Write Cache), RAID 0/1/5/6/10
- **P840ar**: High-performance, 4GB FBWC, RAID 0/1/5/6/10/50/60
- **P440**: PCIe card version, 2GB FBWC
- **P840**: PCIe card version, 4GB FBWC

**HBA Mode** (non-RAID pass-through):
- Smart Array controllers in HBA mode for software RAID (ZFS, mdadm)
- Limited support; check firmware version

**Alternative Controllers**:
- LSI/Broadcom HBA controllers in PCIe slots
- H240ar (12Gb/s HBA mode)

### Boot Drive Options

For network-focused deployments:
- **Minimal Local Storage**: 2x SSD in RAID 1 for hypervisor/OS
- **USB/SD Boot**: iLO supports USB boot, SD card (internal USB)
- **Diskless**: Pure network boot (subject of network-boot.md)

## Network Connectivity

### Integrated FlexibleLOM

The DL360 Gen9 includes a FlexibleLOM slot for swappable network adapters:

**Common FlexibleLOM Options**:
- **HPE 366FLR**: 4x 1GbE (Broadcom BCM5719)
  - Most common, good for general use
  - Supports PXE, UEFI network boot, SR-IOV
  
- **HPE 560FLR-SFP+**: 2x 10GbE SFP+ (Intel X710)
  - High performance, fiber or DAC
  - Supports PXE, UEFI boot, SR-IOV, RDMA (RoCE)

- **HPE 361i**: 2x 1GbE (Intel I350)
  - Entry-level, good driver support

### PCIe Expansion Slots

**Slot Configuration**:
- **Slot 1**: PCIe 3.0 x16 (low-profile)
- **Slot 2**: PCIe 3.0 x8 (low-profile)
- **Slot 3**: PCIe 3.0 x8 (low-profile) - optional, depends on riser

**Network Card Options**:
- Intel X520/X710 (10GbE)
- Mellanox ConnectX-3/ConnectX-4 (10/25/40GbE, InfiniBand)
- Broadcom NetXtreme (1/10/25GbE)

**Note**: Ensure cards are low-profile for 1U chassis compatibility

## Power Supply

### PSU Options

- **500W**: Single PSU, non-redundant (not recommended)
- **800W**: Common, supports dual CPU + moderate expansion
- **1400W**: High-power, dual CPU with high TDP + GPUs
- **Redundancy**: 1+1 redundant hot-plug recommended

### Power Configuration

- **Platinum Efficiency**: 94%+ at 50% load
- **Hot-Plug**: Replace without powering down
- **Auto-Switching**: 100-240V AC, 50/60Hz

**Home Lab Power Draw** (typical):
- Idle (dual E5-2650 v3, 128GB RAM): 100-130W
- Load: 200-350W depending on CPU and drive configuration

### Power Management

- **HPE Dynamic Power Capping**: Limit max power via iLO
- **Collaborative Power**: Share power budget across chassis in blade environments
- **Energy Efficient Ethernet (EEE)**: Reduce NIC power during low utilization

## Cooling and Acoustics

### Fan Configuration

- **6x Hot-Plug Fans**: Front-mounted, redundant (N+1)
- **Variable Speed**: Controlled by System ROM based on thermal sensors
- **iLO Management**: Monitor fan speed, temperature via iLO

### Thermal Management

- **Temperature Range**: 10-35°C (50-95°F) operating
- **Altitude**: Up to 3,050m (10,000 ft) at reduced temperature
- **Airflow**: Front-to-back, ensure clear intake and exhaust

### Noise Level

- **Idle**: ~45 dBA (quiet for 1U server)
- **Load**: 55-70 dBA depending on thermal demand
- **Home Lab Consideration**: Audible but acceptable in dedicated space; louder than desktop workstation

**Noise Reduction**:
- Run lower TDP CPUs (e.g., E5-2620 series)
- Maintain ambient temperature <25°C
- Ensure adequate airflow (not in enclosed cabinet without ventilation)

## Management - iLO 4

### iLO 4 Features

The Integrated Lights-Out 4 (iLO 4) provides out-of-band management:

- **Web Interface**: HTTPS management console
- **Remote Console**: HTML5 or Java-based KVM
- **Virtual Media**: Mount ISOs/images remotely
- **Power Control**: Power on/off, reset, cold boot
- **Monitoring**: Sensors, event logs, hardware health
- **Alerting**: Email alerts, SNMP traps, syslog
- **Scripting**: RESTful API (Redfish standard)

### iLO Licensing

- **iLO Standard** (included): Basic management, remote console
- **iLO Advanced** (license required): 
  - Virtual media
  - Remote console performance improvements
  - Directory integration (LDAP/AD)
  - Graphical remote console
- **iLO Advanced Premium** (license required):
  - Insight Remote Support
  - Federation
  - Jitter smoothing

**Home Lab**: iLO Advanced license highly recommended for virtual media and full remote console features

### iLO Network Configuration

- **Dedicated iLO Port**: Separate 1GbE management port (recommended)
- **Shared LOM**: Share FlexibleLOM port with OS (not recommended for isolation)

**Security**: Isolate iLO on dedicated management VLAN, disable if not needed

## BIOS and Firmware

### System ROM (BIOS/UEFI)

- **Firmware Type**: UEFI 2.31 or later
- **Boot Modes**: UEFI, Legacy BIOS, or hybrid
- **Configuration**: RBSU (ROM-Based Setup Utility) accessible via F9

### Firmware Update Methods

1. **Service Pack for ProLiant (SPP)**: Comprehensive bundle of all firmware
2. **iLO Online Flash**: Update via web interface
3. **Online ROM Flash**: Linux utility for online updates
4. **USB Flash**: Boot from USB with firmware update utility

**Recommended Practice**: Update to latest SPP for security patches and feature improvements

### Secure Boot

- **UEFI Secure Boot**: Supported, validates boot loader signatures
- **TPM**: Optional Trusted Platform Module 1.2 or 2.0
- **Boot Order Protection**: Prevent unauthorized boot device changes

## Expansion and Modularity

### GPU Support

Limited GPU support due to 1U form factor and power constraints:
- **Low-Profile GPUs**: Nvidia T4, AMD Instinct MI25 (may require custom cooling)
- **Power**: Consider 1400W PSU for high-power GPUs
- **Not Ideal**: For GPU-heavy workloads, consider 2U+ servers (e.g., DL380 Gen9)

### USB Ports

- **Front**: 1x USB 3.0
- **Rear**: 2x USB 3.0
- **Internal**: 1x USB 2.0 (for SD/USB boot device)

### Serial Port

- Rear serial port for legacy console access
- Useful for network equipment serial console, debug

## Home Lab Considerations

### Pros for Home Lab

1. **Density**: 1U form factor saves rack space
2. **iLO Management**: Enterprise remote management without KVM
3. **Network Boot**: Excellent PXE/UEFI boot support (see network-boot.md)
4. **Serviceability**: Hot-swap drives, PSU, fans
5. **Documentation**: Extensive HPE documentation and community support
6. **Parts Availability**: Common on secondary market, affordable

### Cons for Home Lab

1. **Noise**: Louder than tower servers or workstations
2. **Power**: Higher idle power than consumer hardware (100-130W idle)
3. **1U Limitations**: Limited GPU, PCIe expansion vs 2U/4U chassis
4. **Firmware**: Requires HPE account for SPP downloads (free but registration required)

### Recommended Home Lab Configuration

**Budget (~$500-800 used)**:
- Dual E5-2620 v3 or v4 (6 cores each, 85W TDP)
- 128GB RAM (8x 16GB DDR4)
- 2x SSD (boot), 4-6x HDD/SSD (data)
- HPE 366FLR (4x 1GbE)
- Dual 500W or 800W PSU (redundant)
- iLO Advanced license

**Performance (~$1000-1500 used)**:
- Dual E5-2680 v4 (14 cores each, 120W TDP)
- 256GB RAM (16x 16GB DDR4)
- 2x NVMe SSD (boot/cache), 6-8x SSD (data)
- HPE 560FLR-SFP+ (2x 10GbE) + PCIe 4x1GbE card
- Dual 800W PSU
- iLO Advanced license

## Comparison with Other Generations

### vs Gen8 (Previous)

**Gen9 Advantages**:
- DDR4 vs DDR3 (lower power, higher capacity)
- Better UEFI support and HTTP boot
- Newer processor architecture (Haswell/Broadwell vs Sandy Bridge/Ivy Bridge)
- iLO 4 vs iLO 3 (better HTML5 console)

**Gen8 Advantages**:
- Lower cost on secondary market
- Adequate for light workloads

### vs Gen10 (Next)

**Gen10 Advantages**:
- Newer CPUs (Skylake-SP/Cascade Lake)
- More PCIe lanes
- Better UEFI firmware and security features
- DDR4-2666/2933 support

**Gen9 Advantages**:
- Lower cost (mature product cycle)
- Excellent value for performance/dollar
- Still well-supported by modern OS and firmware

## Technical Resources

- **QuickSpecs**: HPE ProLiant DL360 Gen9 Server QuickSpecs
- **User Guide**: HPE ProLiant DL360 Gen9 Server User Guide
- **Maintenance and Service Guide**: Detailed disassembly and part replacement
- **Firmware Downloads**: HPE Support Portal (requires free account)

## Summary

The HP ProLiant DL360 Gen9 remains an excellent choice for home labs and small deployments in 2024-2025. Its balance of performance (dual Xeon v4, 768GB RAM capacity), manageability (iLO 4), and network boot capabilities make it particularly well-suited for virtualization, container hosting, and infrastructure automation workflows. While not the latest generation, it offers strong value with robust firmware support and wide secondary market availability.

**Best For**: 
- Virtualization hosts (ESXi, Proxmox, Hyper-V)
- Kubernetes/container platforms
- Network boot/diskless deployments
- Storage servers (with appropriate controller)
- General compute workloads

**Avoid For**:
- GPU-intensive workloads (1U constraints)
- Noise-sensitive environments (unless isolated)
- Extreme low-power requirements (100W+ idle)
