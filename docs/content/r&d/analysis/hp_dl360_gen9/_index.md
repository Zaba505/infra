---
title: "HP ProLiant DL360 Gen9 Analysis"
type: docs
description: >
  Technical analysis of HP ProLiant DL360 Gen9 server capabilities with focus on network boot support
---

This section contains detailed analysis of the HP ProLiant DL360 Gen9 server platform, including hardware specifications, network boot capabilities, and configuration guidance for home lab deployments.

## Overview

The HP ProLiant DL360 Gen9 is a 1U rack-mountable server released by HPE as part of their Generation 9 (Gen9) product line, introduced in 2014. It's a popular choice for home labs due to its balance of performance, density, and relative power efficiency compared to earlier generations.

## Key Features

- **Form Factor**: 1U rack-mountable
- **Processor Support**: Dual Intel Xeon E5-2600 v3/v4 processors (Haswell/Broadwell)
- **Memory**: Up to 768GB DDR4 RAM (24 DIMM slots)
- **Storage**: Flexible SFF/LFF drive configurations
- **Network**: Integrated quad-port 1GbE or 10GbE FlexibleLOM options
- **Management**: iLO 4 (Integrated Lights-Out) with remote KVM and virtual media
- **Boot Options**: UEFI and Legacy BIOS support with extensive network boot capabilities

## Documentation Sections

- [Network Boot Capabilities](./network-boot/) - Detailed analysis of PXE, iPXE, and UEFI HTTP boot support
- [Hardware Specifications](./specifications/) - Complete hardware configuration details
- [Configuration Guide](./configuration/) - Setup and optimization recommendations
