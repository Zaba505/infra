---
title: "Amazon Web Services Analysis"
type: docs
description: >
  Technical analysis of Amazon Web Services capabilities for hosting network boot infrastructure
---

This section contains detailed analysis of Amazon Web Services (AWS) for hosting the network boot server infrastructure, evaluating its support for TFTP, HTTP/HTTPS routing, and WireGuard VPN connectivity as required by ADR-0002.

## Overview

Amazon Web Services is Amazon's comprehensive cloud computing platform, offering compute, storage, networking, and managed services. This analysis focuses on AWS's capabilities to support the network boot architecture decided in [ADR-0002](../../adrs/0002-network-boot-architecture/).

## Key Services Evaluated

- **EC2**: Virtual machine instances for hosting boot server
- **VPN / VPC**: Network connectivity and VPN capabilities
- **Elastic Load Balancing**: Application and Network Load Balancers
- **NAT Gateway**: Network address translation for outbound connectivity
- **VPC**: Virtual Private Cloud networking and routing

## Documentation Sections

- [Network Boot Support](./network-boot/) - Analysis of TFTP, HTTP, and HTTPS routing capabilities
- [WireGuard Support](./wireguard/) - Evaluation of WireGuard VPN integration options
