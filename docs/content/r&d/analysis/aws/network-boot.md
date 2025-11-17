---
title: "AWS Network Boot Protocol Support"
type: docs
description: >
  Analysis of Amazon Web Services support for TFTP, HTTP, and HTTPS routing for network boot infrastructure
---

# Network Boot Protocol Support on Amazon Web Services

This document analyzes AWS's capabilities for hosting network boot infrastructure, specifically focusing on TFTP, HTTP, and HTTPS protocol support.

## TFTP (Trivial File Transfer Protocol) Support

### Native Support

**Status**: ❌ **Not natively supported by Elastic Load Balancing**

AWS's Elastic Load Balancing services do **not** support TFTP protocol natively:

- **Application Load Balancer (ALB)**: HTTP/HTTPS only (Layer 7)
- **Network Load Balancer (NLB)**: TCP/UDP support, but **not TFTP-aware**
- **Classic Load Balancer**: Deprecated, similar limitations

TFTP operates on UDP port 69 with unique protocol semantics (variable block sizes, retransmissions, port negotiation) that standard load balancers cannot parse.

### Implementation Options

#### Option 1: Direct EC2 Instance Access (Recommended for VPN Scenario)

Since ADR-0002 specifies a VPN-based architecture, TFTP can be served directly from an EC2 instance:

- **Approach**: Run TFTP server (e.g., `tftpd-hpa`, `dnsmasq`) on an EC2 instance
- **Access**: Home lab connects via VPN tunnel to instance's private IP
- **Security Group**: Allow UDP/69 from VPN subnet/security group
- **Pros**:
  - Simple implementation
  - No load balancer needed (single boot server sufficient for home lab)
  - TFTP traffic encrypted through VPN tunnel
  - Direct instance-to-client communication
- **Cons**:
  - Single point of failure (no HA)
  - Manual failover if instance fails

#### Option 2: Network Load Balancer (NLB) UDP Passthrough

While NLB doesn't understand TFTP protocol, it can forward UDP traffic:

- **Approach**: Configure NLB to forward UDP/69 to target group
- **Limitations**:
  - No TFTP-specific health checks
  - Health checks would use TCP or different protocol
  - Adds cost and complexity without significant benefit for single server
- **Use Case**: Only relevant for multi-AZ HA deployment (overkill for home lab)

### TFTP Security Considerations

- **Encryption**: TFTP itself is unencrypted, but VPN tunnel provides encryption
- **Security Groups**: Restrict UDP/69 to VPN security group or CIDR only
- **File Access Control**: Configure TFTP server with restricted file access
- **Read-Only Mode**: Deploy TFTP server in read-only mode to prevent uploads

## HTTP Support

### Native Support

**Status**: ✅ **Fully supported**

AWS provides comprehensive HTTP support through multiple services:

#### Elastic Load Balancing - Application Load Balancer

- **Protocol Support**: HTTP/1.1, HTTP/2, HTTP/3 (preview)
- **Port**: Any port (typically 80 for HTTP)
- **Routing**: Path-based, host-based, query string, header-based routing
- **Health Checks**: HTTP health checks with configurable paths and response codes
- **SSL Offloading**: Terminate SSL at ALB and use HTTP to backend
- **Backend**: EC2 instances, ECS, EKS, Lambda

#### EC2 Direct Access

For VPN scenario, HTTP can be served directly from EC2 instance:

- **Approach**: Run HTTP server (nginx, Apache, custom service) on EC2
- **Access**: Home lab accesses via VPN tunnel to private IP
- **Security Group**: Allow TCP/80 from VPN security group
- **Pros**: Simpler than ALB for single boot server

### HTTP Boot Flow for Network Boot

1. **PXE → TFTP**: Initial bootloader (iPXE) loaded via TFTP
2. **iPXE → HTTP**: iPXE chainloads kernel/initrd via HTTP
3. **Kernel/Initrd**: Large boot files served efficiently over HTTP

### Performance Considerations

- **Connection Pooling**: HTTP/1.1 keep-alive reduces connection overhead
- **Compression**: gzip compression for text-based configs
- **CloudFront**: Optional CDN for caching boot files (probably overkill for VPN scenario)
- **TCP Optimization**: AWS network optimized for low-latency TCP

## HTTPS Support

### Native Support

**Status**: ✅ **Fully supported with advanced features**

AWS provides enterprise-grade HTTPS support:

#### Elastic Load Balancing - Application Load Balancer

- **Protocol Support**: HTTPS/1.1, HTTP/2 over TLS, HTTP/3 (preview)
- **SSL/TLS Termination**: Terminate SSL at ALB
- **Certificate Management**:
  - AWS Certificate Manager (ACM) - free SSL certificates with automatic renewal
  - Import custom certificates
  - Integration with private CA via ACM Private CA
- **TLS Versions**: TLS 1.0, 1.1, 1.2, 1.3 (configurable via security policy)
- **Cipher Suites**: Predefined security policies (modern, compatible, legacy)
- **SNI Support**: Multiple certificates on single load balancer

#### AWS Certificate Manager (ACM)

- **Free Certificates**: No cost for public SSL certificates used with AWS services
- **Automatic Renewal**: ACM automatically renews certificates before expiration
- **Private CA**: ACM Private CA for internal PKI (additional cost)
- **Integration**: Native integration with ALB, CloudFront, API Gateway

### HTTPS for Network Boot

#### Use Case

Modern UEFI firmware and iPXE support HTTPS boot:

- **iPXE HTTPS**: iPXE compiled with `DOWNLOAD_PROTO_HTTPS` can fetch over HTTPS
- **UEFI HTTP Boot**: UEFI firmware natively supports HTTP/HTTPS boot
- **Security**: Boot file integrity verified via HTTPS chain of trust

#### Implementation on AWS

1. **Certificate Provisioning**:
   - Use ACM certificate for public domain (free, auto-renewed)
   - Use self-signed certificate for VPN-only access (add to iPXE trust store)
   - Use ACM Private CA for internal PKI ($400/month - expensive for home lab)

2. **ALB Configuration**:
   - HTTPS listener on port 443
   - Target group pointing to EC2 boot server
   - Security policy with TLS 1.2+ minimum

3. **Alternative: Direct EC2 HTTPS**:
   - Run nginx/Apache with TLS on EC2 instance
   - Access via VPN tunnel to private IP with HTTPS
   - Simpler setup for VPN-only scenario
   - Use Let's Encrypt or self-signed certificate

### Mutual TLS (mTLS) Support

AWS ALB supports mutual TLS authentication (as of 2022):

- **Client Certificates**: Require client certificates for authentication
- **Trust Store**: Upload trusted CA certificates to ALB
- **Use Case**: Ensure only authorized home lab servers can access boot files
- **Integration**: Combine with VPN for defense-in-depth
- **Passthrough Mode**: ALB can pass client cert to backend for validation

## Routing and Load Balancing Capabilities

### VPC Routing

- **Route Tables**: Define routes to direct traffic through VPN gateway
- **Route Propagation**: BGP route propagation for VPN connections
- **Transit Gateway**: Advanced multi-VPC/VPN routing (overkill for home lab)

### Security Groups

- **Stateful Firewall**: Automatic return traffic handling
- **Ingress/Egress Rules**: Fine-grained control by protocol, port, source/destination
- **Security Group Chaining**: Reference security groups in rules (elegant for VPN setup)
- **VPN Subnet Restriction**: Allow traffic only from VPN-connected subnet

### Network ACLs (Optional)

- **Stateless Firewall**: Subnet-level access control
- **Defense in Depth**: Additional layer beyond security groups
- **Use Case**: Probably unnecessary for simple VPN boot server

## Cost Implications

### Data Transfer Costs

- **VPN Traffic**: Data transfer through VPN gateway charged at standard rates
- **Intra-Region**: Free for traffic within same region/VPC
- **Boot File Sizes**: Typical kernel + initrd = 50-200MB per boot
- **Monthly Estimate**: 10 boots/month × 150MB = 1.5GB ≈ $0.14/month (US East egress)

### Load Balancing Costs

- **Application Load Balancer**: ~$0.0225/hour + $0.008 per LCU-hour (~$16-20/month minimum)
- **Network Load Balancer**: ~$0.0225/hour + $0.006 per NLCU-hour (~$16-18/month minimum)
- **For VPN Scenario**: Load balancer unnecessary (single EC2 instance sufficient)

### Compute Costs

- **t3.micro Instance**: ~$7.50/month (on-demand pricing, US East)
- **t4g.micro Instance**: ~$6.00/month (ARM-based, cheaper, sufficient for boot server)
- **Reserved Instances**: Up to 72% savings with 1-year or 3-year commitment
- **Savings Plans**: Flexible discounts for consistent compute usage

### ACM Certificate Costs

- **Public Certificates**: **Free** when used with AWS services
- **Private CA**: $400/month (too expensive for home lab)

## Comparison with Requirements

| Requirement | AWS Support | Implementation |
|------------|-------------|----------------|
| TFTP | ⚠️ Via EC2, not ELB | Direct EC2 access via VPN |
| HTTP | ✅ Full support | EC2 or ALB |
| HTTPS | ✅ Full support | EC2 or ALB with ACM |
| VPN Integration | ✅ Native VPN | Site-to-Site VPN or self-managed |
| Load Balancing | ✅ ALB, NLB | Optional for HA |
| Certificate Mgmt | ✅ ACM (free) | Automatic renewal |
| Cost Efficiency | ✅ Low-cost instances | t4g.micro sufficient |

## Recommendations

### For VPN-Based Architecture (per ADR-0002)

1. **EC2 Instance**: Deploy single t4g.micro or t3.micro instance with:
   - TFTP server (`tftpd-hpa` or `dnsmasq`)
   - HTTP server (nginx or simple Python HTTP server)
   - Optional HTTPS with Let's Encrypt or self-signed certificate

2. **VPN Connection**: Connect home lab to AWS via:
   - Site-to-Site VPN (IPsec) - managed service, higher cost (~$36/month)
   - Self-managed WireGuard on EC2 - lower cost, more control

3. **Security Groups**: Restrict access to:
   - UDP/69 (TFTP) from VPN security group only
   - TCP/80 (HTTP) from VPN security group only
   - TCP/443 (HTTPS) from VPN security group only

4. **No Load Balancer**: For home lab scale, direct EC2 access is sufficient

5. **Health Monitoring**: Use CloudWatch for instance and service health

### If HA Required (Future Enhancement)

- Deploy multi-AZ EC2 instances with Network Load Balancer
- Use S3 as backend for boot files with EC2 serving as cache
- Implement auto-recovery with Auto Scaling Group (min=max=1)

## References

- [AWS Elastic Load Balancing Documentation](https://docs.aws.amazon.com/elasticloadbalancing/)
- [AWS Certificate Manager](https://docs.aws.amazon.com/acm/)
- [AWS Site-to-Site VPN](https://docs.aws.amazon.com/vpn/)
- [iPXE HTTPS Boot](https://ipxe.org/crypto)
- [UEFI HTTP Boot Specification](https://uefi.org/specs/UEFI/2.10/24_Network_Protocols.html#http-boot)
