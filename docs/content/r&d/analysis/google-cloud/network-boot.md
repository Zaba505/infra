---
title: "GCP Network Boot Protocol Support"
type: docs
description: >
  Analysis of Google Cloud Platform's support for TFTP, HTTP, and HTTPS routing for network boot infrastructure
---

# Network Boot Protocol Support on Google Cloud Platform

This document analyzes GCP's capabilities for hosting network boot infrastructure, specifically focusing on TFTP, HTTP, and HTTPS protocol support.

## TFTP (Trivial File Transfer Protocol) Support

### Native Support

**Status**: ❌ **Not natively supported by Cloud Load Balancing**

GCP's Cloud Load Balancing services (Application Load Balancer, Network Load Balancer) do **not** support TFTP protocol natively. TFTP operates on UDP port 69 and has unique protocol requirements that are not compatible with GCP's load balancing services.

### Implementation Options

#### Option 1: Direct VM Access (Recommended for VPN Scenario)

Since ADR-0002 specifies a VPN-based architecture, TFTP can be served directly from a Compute Engine VM without load balancing:

- **Approach**: Run TFTP server (e.g., `tftpd-hpa`, `dnsmasq`) on a Compute Engine VM
- **Access**: Home lab connects via VPN tunnel to the VM's private IP
- **Routing**: VPC firewall rules allow UDP/69 from VPN subnet
- **Pros**:
  - Simple implementation
  - No need for load balancing (single boot server sufficient)
  - TFTP traffic encrypted through VPN tunnel
  - Direct VM-to-client communication
- **Cons**:
  - Single point of failure (no load balancing/HA)
  - Manual failover required if VM fails

#### Option 2: Network Load Balancer (NLB) Passthrough

While NLB doesn't parse TFTP protocol, it can forward UDP traffic:

- **Approach**: Configure Network Load Balancer for UDP/69 passthrough
- **Limitations**:
  - No protocol-aware health checks for TFTP
  - Health checks would use TCP or HTTP on alternate port
  - Adds complexity without significant benefit for single boot server
- **Use Case**: Only relevant for multi-region HA deployment (overkill for home lab)

### TFTP Security Considerations

- **Encryption**: TFTP protocol itself is unencrypted, but VPN tunnel provides encryption
- **Firewall Rules**: Restrict UDP/69 to VPN subnet only (no public access)
- **File Access Control**: Configure TFTP server with restricted file access
- **Read-Only Mode**: Deploy TFTP server in read-only mode to prevent uploads

## HTTP Support

### Native Support

**Status**: ✅ **Fully supported**

GCP provides comprehensive HTTP support through multiple services:

#### Cloud Load Balancing - Application Load Balancer

- **Protocol Support**: HTTP/1.1, HTTP/2, HTTP/3 (QUIC)
- **Port**: Any port (typically 80 for HTTP)
- **Routing**: URL-based routing, host-based routing, path-based routing
- **Health Checks**: HTTP health checks with configurable paths
- **SSL Offloading**: Can terminate SSL at load balancer and use HTTP backend
- **Backend**: Compute Engine VMs, instance groups, Cloud Run, GKE

#### Compute Engine Direct Access

For VPN scenario, HTTP can be served directly from VM:

- **Approach**: Run HTTP server (nginx, Apache, custom service) on Compute Engine VM
- **Access**: Home lab accesses via VPN tunnel to private IP
- **Firewall**: VPC firewall rules allow TCP/80 from VPN subnet
- **Pros**: Simpler than load balancer for single boot server

### HTTP Boot Flow for Network Boot

1. **PXE → TFTP**: Initial bootloader (iPXE) loaded via TFTP
2. **iPXE → HTTP**: iPXE chainloads boot files via HTTP from same server
3. **Kernel/Initrd**: Large boot files served efficiently over HTTP

### Performance Considerations

- **Connection Pooling**: HTTP/1.1 keep-alive reduces connection overhead
- **Compression**: gzip compression for text-based boot configs
- **Caching**: Cloud CDN can cache boot files for faster delivery
- **TCP Optimization**: GCP's network optimized for low-latency TCP

## HTTPS Support

### Native Support

**Status**: ✅ **Fully supported with advanced features**

GCP provides enterprise-grade HTTPS support:

#### Cloud Load Balancing - Application Load Balancer

- **Protocol Support**: HTTPS/1.1, HTTP/2 over TLS, HTTP/3 with QUIC
- **SSL/TLS Termination**: Terminate SSL at load balancer
- **Certificate Management**:
  - Google-managed SSL certificates (automatic renewal)
  - Self-managed certificates (bring your own)
  - Certificate Map for multiple domains
- **TLS Versions**: TLS 1.0, 1.1, 1.2, 1.3 (configurable minimum version)
- **Cipher Suites**: Modern, compatible, or custom cipher suites
- **mTLS Support**: Mutual TLS authentication (client certificates)

#### Certificate Manager

- **Managed Certificates**: Automatic provisioning and renewal via Let's Encrypt integration
- **Private CA**: Integration with Google Cloud Certificate Authority Service
- **Certificate Maps**: Route different domains to different backends based on SNI
- **Certificate Monitoring**: Automatic alerts before expiration

### HTTPS for Network Boot

#### Use Case

Modern UEFI firmware and iPXE support HTTPS boot:

- **iPXE HTTPS**: iPXE compiled with `DOWNLOAD_PROTO_HTTPS` can fetch over HTTPS
- **UEFI HTTP Boot**: UEFI firmware natively supports HTTP/HTTPS boot (RFC 3720 iSCSI boot)
- **Security**: Boot file integrity verified via HTTPS chain of trust

#### Implementation on GCP

1. **Certificate Provisioning**:
   - Use Google-managed certificate for public domain (if boot server has public DNS)
   - Use self-signed certificate for VPN-only access (add to iPXE trust store)
   - Use private CA for internal PKI

2. **Load Balancer Configuration**:
   - HTTPS frontend (port 443)
   - Backend service to Compute Engine VM running boot server
   - SSL policy with TLS 1.2+ minimum

3. **Alternative: Direct VM HTTPS**:
   - Run nginx/Apache with TLS on Compute Engine VM
   - Access via VPN tunnel to private IP with HTTPS
   - Simpler setup for VPN-only scenario

### mTLS Support for Enhanced Security

GCP's Application Load Balancer supports mutual TLS authentication:

- **Client Certificates**: Require client certificates for additional authentication
- **Certificate Validation**: Validate client certificates against trusted CA
- **Use Case**: Ensure only authorized home lab servers can access boot files
- **Integration**: Combine with VPN for defense-in-depth

## Routing and Load Balancing Capabilities

### VPC Routing

- **Custom Routes**: Define routes to direct traffic through VPN gateway
- **Route Priority**: Configure route priorities for failover scenarios
- **BGP Support**: Dynamic routing with Cloud Router (for advanced VPN setups)

### Firewall Rules

- **Ingress/Egress Rules**: Fine-grained control over traffic
- **Source/Destination Filters**: IP ranges, tags, service accounts
- **Protocol Filtering**: Allow specific protocols (UDP/69, TCP/80, TCP/443)
- **VPN Subnet Restriction**: Limit access to VPN-connected home lab subnet

### Cloud Armor (Optional)

For additional security if boot server has public access:

- **DDoS Protection**: Layer 3/4 DDoS mitigation
- **WAF Rules**: Application-level filtering
- **IP Allowlisting**: Restrict to known public IPs
- **Rate Limiting**: Prevent abuse

## Cost Implications

### Network Egress Costs

- **VPN Traffic**: Egress to VPN endpoint charged at standard internet egress rates
- **Intra-Region**: Free for traffic within same region
- **Boot File Sizes**: Typical kernel + initrd = 50-200MB per boot
- **Monthly Estimate**: 10 boots/month × 150MB = 1.5GB ≈ $0.18/month (US egress)

### Load Balancing Costs

- **Application Load Balancer**: ~$0.025/hour + $0.008 per LCU-hour
- **Network Load Balancer**: ~$0.025/hour + data processing charges
- **For VPN Scenario**: Load balancer likely unnecessary (single VM sufficient)

### Compute Costs

- **e2-micro Instance**: ~$6-7/month (suitable for boot server)
- **f1-micro Instance**: ~$4-5/month (even smaller, might suffice)
- **Reserved/Committed Use**: Discounts for long-term commitment

## Comparison with Requirements

| Requirement | GCP Support | Implementation |
|------------|-------------|----------------|
| TFTP | ⚠️ Via VM, not LB | Direct VM access via VPN |
| HTTP | ✅ Full support | VM or ALB |
| HTTPS | ✅ Full support | VM or ALB with Certificate Manager |
| VPN Integration | ✅ Native VPN | Cloud VPN or self-managed WireGuard |
| Load Balancing | ✅ ALB, NLB | Optional for HA |
| Certificate Mgmt | ✅ Managed certs | Certificate Manager |
| Cost Efficiency | ✅ Low-cost VMs | e2-micro sufficient |

## Recommendations

### For VPN-Based Architecture (per ADR-0002)

1. **Compute Engine VM**: Deploy single e2-micro VM with:
   - TFTP server (`tftpd-hpa` or `dnsmasq`)
   - HTTP server (nginx or simple Python HTTP server)
   - Optional HTTPS with self-signed certificate

2. **VPN Tunnel**: Connect home lab to GCP via:
   - Cloud VPN (IPsec) - easier setup, higher cost
   - Self-managed WireGuard on Compute Engine - lower cost, more control

3. **VPC Firewall**: Restrict access to:
   - UDP/69 (TFTP) from VPN subnet only
   - TCP/80 (HTTP) from VPN subnet only
   - TCP/443 (HTTPS) from VPN subnet only

4. **No Load Balancer**: For home lab scale, direct VM access is sufficient

5. **Health Monitoring**: Use Cloud Monitoring for VM and service health

### If HA Required (Future Enhancement)

- Deploy multi-zone VMs with Network Load Balancer
- Use Cloud Storage as backend for boot files with VM serving as cache
- Implement failover automation with Cloud Functions

## References

- [GCP Cloud Load Balancing Documentation](https://cloud.google.com/load-balancing/docs)
- [GCP Certificate Manager](https://cloud.google.com/certificate-manager/docs)
- [GCP Cloud VPN](https://cloud.google.com/network-connectivity/docs/vpn)
- [iPXE HTTPS Boot](https://ipxe.org/crypto)
- [UEFI HTTP Boot](https://uefi.org/specs/UEFI/2.10/24_Network_Protocols.html#http-boot)
