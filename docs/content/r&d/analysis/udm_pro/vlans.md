---
type: docs
title: "UDM Pro VLAN Configuration & Capabilities"
linkTitle: "VLANs"
description: >
  Detailed analysis of VLAN support on the Ubiquiti Dream Machine Pro, including port-based VLAN assignment and VPN integration.
---

## Overview

The **Ubiquiti Dream Machine Pro (UDM Pro)** provides robust VLAN support through native 802.1Q tagging, enabling network segmentation for security, performance, and organizational purposes. This document covers VLAN configuration capabilities, port assignments, and VPN integration.

## VLAN Fundamentals on UDM Pro

### Supported Standards
- **802.1Q VLAN Tagging**: Full support for standard VLAN tagging
- **VLAN Range**: IDs 1-4094 (standard IEEE 802.1Q range)
- **Maximum VLANs**: Up to 32 networks/VLANs per device
- **Native VLAN**: Configurable per port (default: VLAN 1)

### VLAN Types

**Corporate Network**
- Default network type for general-purpose VLANs
- Provides DHCP, inter-VLAN routing, and firewall capabilities
- Can enable/disable guest policies, IGMP snooping, and multicast DNS

**Guest Network**
- Isolated network with internet-only access
- Automatic firewall rules preventing access to other VLANs
- Captive portal support for guest authentication

**IoT Network**
- Optimized for IoT devices with device isolation
- Prevents lateral movement between IoT devices
- Allows communication with controller/gateway only

## Port-Based VLAN Assignment

### Per-Port VLAN Configuration

The UDM Pro's **8x 1 Gbps LAN ports** and **SFP/SFP+ ports** support flexible VLAN assignment:

**Configuration Options per Port:**
1. **Native VLAN/Untagged VLAN**: The default VLAN for untagged traffic on the port
2. **Tagged VLANs**: Multiple VLANs that can pass through the port with 802.1Q tags
3. **Port Profile**: Pre-configured VLAN assignments that can be applied to ports

### Port Profile Types

**All**: Port accepts all VLANs (trunk mode)
- Passes all configured VLANs with tags
- Used for connecting managed switches or access points
- Native VLAN for untagged traffic

**Specific VLANs**: Port limited to selected VLANs
- Choose which VLANs are allowed (tagged)
- Set native/untagged VLAN
- Used for controlled trunk links

**Single VLAN**: Access port mode
- Port carries only one VLAN (untagged)
- All traffic on this port belongs to specified VLAN
- Used for end devices (PCs, servers, printers)

### Configuration Steps

**Via UniFi Controller GUI:**

1. **Create Port Profile**:
   - Navigate to **Settings** → **Profiles** → **Port Manager**
   - Click **Create New Port Profile**
   - Select profile type (All, LAN, or Custom)
   - Configure VLAN settings:
     - **Native VLAN/Network**: Untagged VLAN
     - **Tagged VLANs**: Select allowed VLANs (for trunk mode)
   - Enable/disable settings: PoE, Storm Control, Port Isolation

2. **Assign Profile to Ports**:
   - Navigate to **UniFi Devices** → Select **UDM Pro**
   - Go to **Ports** tab
   - For each LAN port (1-8) or SFP port:
     - Click port to edit
     - Select **Port Profile** from dropdown
     - Apply changes

3. **Quick Port Assignment** (Alternative):
   - **Settings** → **Networks** → Select VLAN
   - Under **Port Manager**, assign specific ports to this network
   - Ports become access ports for this VLAN

### Example Port Layout

```
UDM Pro Port Assignment Example:

Port 1: Native VLAN 10 (Management) - Access Mode
        └── Use: Ansible control server

Port 2: Native VLAN 20 (Kubernetes) - Access Mode
        └── Use: K8s master node

Port 3: Native VLAN 30 (Storage) - Access Mode
        └── Use: NAS/SAN device

Port 4: Native VLAN 1, Tagged: 10,20,30,40 - Trunk Mode
        └── Use: Managed switch uplink

Port 5-7: Native VLAN 40 (DMZ) - Access Mode
          └── Use: Public-facing servers

Port 8: Native VLAN 1 (Default/Untagged) - Access Mode
        └── Use: Management laptop (temporary)

SFP+: Native VLAN 1, Tagged: All - Trunk Mode
      └── Use: 10G uplink to core switch
```

## VLAN Features and Capabilities

### Inter-VLAN Routing

**Enabled by Default:**
- Hardware-accelerated routing between VLANs
- Wire-speed performance (8 Gbps backplane)
- Routing decisions made at Layer 3

**Firewall Control:**
- Default behavior: Allow all inter-VLAN traffic
- Recommended: Create explicit allow/deny rules per VLAN pair
- Granular control: Protocol, port, source/destination filtering

**Example Firewall Rules:**
```
Rule 1: Allow Management (VLAN 10) → All VLANs
        Source: 192.168.10.0/24
        Destination: Any
        Action: Accept

Rule 2: Allow K8s (VLAN 20) → Storage (VLAN 30) - NFS only
        Source: 192.168.20.0/24
        Destination: 192.168.30.0/24
        Ports: 2049 (NFS), 111 (Portmapper)
        Action: Accept

Rule 3: Block IoT (VLAN 50) → All Private Networks
        Source: 192.168.50.0/24
        Destination: 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12
        Action: Drop

Rule 4 (Implicit): Default Deny Between VLANs
        Source: Any
        Destination: Any
        Action: Drop
```

### DHCP per VLAN

**Each VLAN can have its own DHCP server:**
- Independent IP ranges per VLAN
- Separate DHCP options (DNS, gateway, NTP, domain)
- Static DHCP reservations per VLAN
- PXE boot options (Option 66/67) per network

**Configuration:**
- **Settings** → **Networks** → Select VLAN
- **DHCP** section:
  - Enable DHCP server
  - Define IP range (e.g., 192.168.10.100-192.168.10.254)
  - Set lease time
  - Configure gateway (usually UDM Pro's IP on this VLAN)
  - Add custom DHCP options

**Example DHCP Configuration:**
```
VLAN 10 (Management):
  Subnet: 192.168.10.0/24
  Gateway: 192.168.10.1 (UDM Pro)
  DHCP Range: 192.168.10.100-192.168.10.200
  DNS: 192.168.10.10 (local DNS server)
  TFTP Server (Option 66): 192.168.10.16
  Boot Filename (Option 67): pxelinux.0

VLAN 20 (Kubernetes):
  Subnet: 192.168.20.0/24
  Gateway: 192.168.20.1 (UDM Pro)
  DHCP Range: 192.168.20.50-192.168.20.99
  DNS: 8.8.8.8, 8.8.4.4
  Domain Name: k8s.lab.local
```

### VLAN Isolation

**Guest Portal Isolation:**
- Guest networks auto-configured with isolation rules
- Prevents access to RFC1918 private networks
- Internet-only access by default

**Manual Isolation (Firewall Rules):**
- Create LAN In rules to block inter-VLAN traffic
- Use groups for easier management of multiple VLANs
- Apply port isolation for additional security

**Device Isolation (IoT Networks):**
- Prevents devices on same VLAN from communicating
- Only controller/gateway access allowed
- Use for untrusted IoT devices (cameras, smart home)

## VPN and VLAN Integration

### Site-to-Site VPN VLAN Assignment

**✅ VLANs CAN be assigned to site-to-site VPN connections:**

**WireGuard VPN:**
- Configure remote subnet to map to specific local VLAN
- Example: GCP subnet 10.128.0.0/20 → routed through VLAN 10
- Routing table automatically updated
- Firewall rules apply to VPN traffic

**IPsec Site-to-Site:**
- Specify local networks (can select specific VLANs)
- Remote networks configured in tunnel settings
- Multiple VLANs can traverse single VPN tunnel
- Perfect Forward Secrecy supported

**Configuration Steps:**
1. **Settings** → **VPN** → **Site-to-Site VPN**
2. **Create New** VPN tunnel (WireGuard or IPsec)
3. Under **Local Networks**, select VLANs to include:
   - Option 1: Select "All" networks
   - Option 2: Choose specific VLANs (e.g., VLAN 10, 20 only)
4. Configure **Remote Networks** (cloud provider subnets)
5. Set encryption parameters and pre-shared keys
6. **Create Firewall Rules** for VPN traffic:
   - Allow specific VLAN → VPN tunnel
   - Control which VLANs can reach remote networks

**Example Site-to-Site Config:**
```
Home Lab → GCP WireGuard VPN

Local Networks:
  - VLAN 10 (Management): 192.168.10.0/24
  - VLAN 20 (Kubernetes): 192.168.20.0/24

Remote Networks:
  - GCP VPC: 10.128.0.0/20

Firewall Rules:
  - Allow VLAN 10 → GCP VPC (all protocols)
  - Allow VLAN 20 → GCP VPC (HTTPS, kubectl API only)
  - Block all other VLANs from VPN tunnel
```

### Remote Access VPN VLAN Assignment

**✅ VLANs CAN be assigned to remote access VPN clients:**

**L2TP/IPsec Remote Access:**
- VPN clients land on a specific VLAN
- Default: All clients in same VPN subnet
- Firewall rules control VLAN access from VPN

**OpenVPN Remote Access (via UniFi Network Application addon):**
- Not natively built into UDM Pro
- Requires UniFi Network Application 6.0+
- Can route VPN clients to specific VLAN

**Teleport VPN (UniFi's solution):**
- Built-in remote access VPN
- Clients route through UDM Pro
- Can access specific VLANs based on firewall rules
- Layer 3 routing to VLANs

**Configuration:**
1. **Settings** → **VPN** → **Remote Access**
2. Enable **L2TP** or configure **Teleport**
3. Set **VPN Network** (e.g., 192.168.100.0/24)
4. **Advanced**:
   - Enable access to specific VLANs
   - By default, VPN network is treated as separate VLAN
5. **Firewall Rules** to allow VPN → VLANs:
   - Source: VPN network (192.168.100.0/24)
   - Destination: VLAN 10, VLAN 20 (or specific resources)
   - Action: Accept

**Example Remote Access Config:**
```
Remote VPN Users → Home Lab Access

VPN Network: 192.168.100.0/24
VPN Gateway: 192.168.100.1 (UDM Pro)

Firewall Rules:
  Rule 1: Allow VPN → Management VLAN (admin users)
          Source: 192.168.100.0/24
          Dest: 192.168.10.0/24
          Ports: SSH (22), HTTPS (443)
  
  Rule 2: Allow VPN → Kubernetes VLAN (developers)
          Source: 192.168.100.0/24
          Dest: 192.168.20.0/24
          Ports: kubectl (6443), app ports (8080-8090)
  
  Rule 3: Block VPN → Storage VLAN (security)
          Source: 192.168.100.0/24
          Dest: 192.168.30.0/24
          Action: Drop
```

### VPN VLAN Routing Limitations

**Current Limitations:**
- Cannot assign individual VPN clients to different VLANs dynamically
- No VLAN assignment based on user identity (all clients in same VPN network)
- RADIUS integration does not support per-user VLAN assignment for VPN
- For per-user VLAN control, use firewall rules based on source IP

**Workarounds:**
- Use firewall rules with VPN client IP ranges for granular access
- Deploy separate VPN tunnels for different access levels
- Use RADIUS for authentication + firewall rules for authorization

## VLAN Best Practices for Home Lab

### Network Segmentation Strategy

**Recommended VLAN Layout:**

```
VLAN 1:   Default/Management (UDM Pro access)
VLAN 10:  Infrastructure Management (Ansible, PXE, monitoring)
VLAN 20:  Kubernetes Cluster (control plane + workers)
VLAN 30:  Storage Network (NFS, iSCSI, object storage)
VLAN 40:  DMZ/Public Services (exposed to internet via Cloudflare)
VLAN 50:  IoT Devices (isolated smart home devices)
VLAN 60:  Guest Network (visitor WiFi, untrusted devices)
VLAN 100: VPN Remote Access (remote admin/dev access)
```

### Firewall Policy Design

**Default Deny Approach:**
1. Create explicit allow rules for necessary traffic
2. Set implicit deny for all inter-VLAN traffic
3. Log dropped packets for troubleshooting

**Rule Order (top to bottom):**
1. Management VLAN → All (with source IP restrictions)
2. Kubernetes → Storage (specific ports)
3. DMZ → Internet (outbound only)
4. VPN → Specific VLANs (based on role)
5. All → Internet (NAT)
6. Block RFC1918 from DMZ
7. Drop all (implicit)

### Performance Optimization

**VLAN Routing Performance:**
- Inter-VLAN routing is hardware-accelerated
- No performance penalty for multiple VLANs
- Use VLAN tagging on trunk ports to reduce switch load

**Multicast and Broadcast Control:**
- Enable IGMP snooping per VLAN for multicast efficiency
- Disable multicast DNS (mDNS) between VLANs if not needed
- Use multicast routing for cross-VLAN multicast (advanced)

## Advanced VLAN Features

### VLAN-Specific Services

**DNS per VLAN:**
- Configure different DNS servers per VLAN via DHCP
- Example: Management VLAN uses local DNS, DMZ uses public DNS

**NTP per VLAN:**
- DHCP Option 42 for NTP server
- Different time sources per network segment

**Domain Name per VLAN:**
- DHCP Option 15 for domain name
- Useful for split-horizon DNS setups

### VLAN Tagging on WiFi

**UniFi WiFi Integration:**
- Each WiFi SSID can map to a specific VLAN
- Multiple SSIDs on same AP → different VLANs
- Seamless VLAN tagging for wireless clients

**Configuration:**
- Create WiFi network in UniFi Controller
- Assign VLAN ID to SSID
- Client traffic automatically tagged

### VLAN Monitoring and Troubleshooting

**Traffic Statistics:**
- Per-VLAN bandwidth usage visible in UniFi Controller
- Deep Packet Inspection (DPI) provides application-level stats
- Export data for analysis in external tools

**Debugging Tools:**
- Port mirroring for packet capture
- Flow logs for traffic analysis
- Firewall logs show inter-VLAN blocks

**Common Issues:**
1. **VLAN not working**: Check port profile assignment and native VLAN config
2. **No inter-VLAN routing**: Verify firewall rules aren't blocking traffic
3. **DHCP not working on VLAN**: Ensure DHCP server enabled on that network
4. **VPN can't reach VLAN**: Check VPN local networks include the VLAN

## Summary

### VLAN Port Assignment: ✅ YES
The UDM Pro **fully supports port-based VLAN assignment**:
- Individual ports can be assigned to specific VLANs (access mode)
- Ports can carry multiple tagged VLANs (trunk mode)
- Native/untagged VLAN configurable per port
- Port profiles simplify configuration across multiple devices

### VPN VLAN Assignment: ✅ YES
VLANs **can be assigned to VPN connections**:
- **Site-to-Site VPN**: Select which VLANs traverse the tunnel
- **Remote Access VPN**: VPN clients route to specific VLANs via firewall rules
- **Routing Control**: Full control over which VLANs are accessible via VPN
- **Limitations**: No per-user VLAN assignment; use firewall rules for granular access

### Key Capabilities
- Up to 32 VLANs supported
- Hardware-accelerated inter-VLAN routing
- Per-VLAN DHCP, DNS, and firewall policies
- Full integration with UniFi WiFi for SSID-to-VLAN mapping
- Flexible port profiles for easy configuration
- VPN integration for both site-to-site and remote access scenarios
