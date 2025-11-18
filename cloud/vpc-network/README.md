# VPC Network Module

Reusable Terraform module for provisioning VPC networks with subnets, firewall rules, and Cloud NAT on Google Cloud Platform.

This module supports common VPC patterns including subnet configuration, firewall rules for WireGuard (UDP/51820), boot server access (TCP/80, TCP/443), and Cloud NAT for outbound connectivity.

## Features

- VPC network with configurable routing mode (REGIONAL or GLOBAL)
- Multiple subnets with configurable CIDR ranges per region
- Secondary IP ranges for GKE/Cloud Run (optional)
- Flexible firewall rules:
  - Ingress rules (configurable protocol, port, source ranges)
  - Egress rules (configurable destinations)
  - Network tag-based targeting
  - Service account-based targeting
  - Configurable logging for troubleshooting
- Cloud Router and Cloud NAT for outbound connectivity
- Zero-downtime updates with `create_before_destroy` lifecycle

## Usage

### Basic VPC with Subnets

```hcl
module "vpc" {
  source = "./cloud/vpc-network"

  name         = "my-vpc"
  description  = "VPC for home lab infrastructure"
  routing_mode = "REGIONAL"

  subnets = {
    "us-central1-subnet" = {
      ip_cidr_range = "10.0.0.0/24"
      region        = "us-central1"
      description   = "Primary subnet in us-central1"
    }
  }
}
```

### Boot Server Network Pattern (ADR-0005)

Example configuration for network boot infrastructure with WireGuard VPN gateway:

```hcl
module "boot_network" {
  source = "./cloud/vpc-network"

  name         = "boot-server-vpc"
  description  = "VPC network for network boot infrastructure and WireGuard gateway"
  routing_mode = "REGIONAL"

  # Subnets for boot server and WireGuard gateway
  subnets = {
    "us-central1-boot" = {
      ip_cidr_range            = "10.128.0.0/20"
      region                   = "us-central1"
      description              = "Subnet for boot server and WireGuard gateway"
      private_ip_google_access = true
    }
  }

  # Firewall rules for WireGuard and boot server
  firewall_rules = [
    # Allow WireGuard UDP traffic from Internet
    {
      name          = "allow-wireguard-ingress"
      direction     = "INGRESS"
      description   = "Allow WireGuard VPN traffic (UDP/51820) from Internet"
      source_ranges = ["0.0.0.0/0"]
      target_tags   = ["wireguard-gateway"]
      allow = [{
        protocol = "udp"
        ports    = ["51820"]
      }]
      log_config = {
        metadata = "INCLUDE_ALL_METADATA"
      }
    },
    # Allow HTTP/HTTPS from WireGuard subnet to boot server
    {
      name          = "allow-boot-server-http"
      direction     = "INGRESS"
      description   = "Allow HTTP/HTTPS traffic to boot server from WireGuard clients"
      source_ranges = ["10.8.0.0/24"] # WireGuard VPN subnet
      target_tags   = ["boot-server"]
      allow = [
        {
          protocol = "tcp"
          ports    = ["80", "443"]
        }
      ]
      log_config = {
        metadata = "INCLUDE_ALL_METADATA"
      }
    },
    # Allow SSH for management
    {
      name          = "allow-ssh-iap"
      direction     = "INGRESS"
      description   = "Allow SSH via Identity-Aware Proxy"
      source_ranges = ["35.235.240.0/20"] # IAP source range
      target_tags   = ["ssh-enabled"]
      allow = [{
        protocol = "tcp"
        ports    = ["22"]
      }]
    },
    # Allow internal communication
    {
      name          = "allow-internal"
      direction     = "INGRESS"
      description   = "Allow all internal traffic within VPC"
      source_ranges = ["10.128.0.0/20"]
      allow = [{
        protocol = "all"
      }]
    },
    # Deny all other ingress by default (implicit, but can be explicit)
  ]

  # Cloud NAT for outbound connectivity (e.g., accessing Cloud Storage)
  enable_cloud_nat = true
  cloud_nat_configs = {
    "us-central1-nat" = {
      router_name    = "us-central1-router"
      nat_name       = "us-central1-nat"
      region         = "us-central1"
      enable_logging = true
      log_filter     = "ERRORS_ONLY"
    }
  }
}
```

### VPC with Secondary IP Ranges for GKE

```hcl
module "gke_vpc" {
  source = "./cloud/vpc-network"

  name         = "gke-vpc"
  description  = "VPC for GKE cluster"
  routing_mode = "REGIONAL"

  subnets = {
    "us-central1-gke" = {
      ip_cidr_range            = "10.0.0.0/24"
      region                   = "us-central1"
      description              = "GKE cluster subnet"
      private_ip_google_access = true
      
      # Secondary IP ranges for GKE pods and services
      secondary_ip_ranges = [
        {
          range_name    = "pods"
          ip_cidr_range = "10.1.0.0/16"
        },
        {
          range_name    = "services"
          ip_cidr_range = "10.2.0.0/16"
        }
      ]
    }
  }

  enable_cloud_nat = true
  cloud_nat_configs = {
    "us-central1-nat" = {
      router_name = "us-central1-router"
      nat_name    = "us-central1-nat"
      region      = "us-central1"
    }
  }
}
```

### Multi-Region VPC with Global Routing

```hcl
module "multi_region_vpc" {
  source = "./cloud/vpc-network"

  name         = "multi-region-vpc"
  description  = "Multi-region VPC with global routing"
  routing_mode = "GLOBAL"

  subnets = {
    "us-central1-subnet" = {
      ip_cidr_range = "10.0.0.0/24"
      region        = "us-central1"
    }
    "europe-west1-subnet" = {
      ip_cidr_range = "10.1.0.0/24"
      region        = "europe-west1"
    }
  }

  enable_cloud_nat = true
  cloud_nat_configs = {
    "us-central1-nat" = {
      router_name = "us-central1-router"
      nat_name    = "us-central1-nat"
      region      = "us-central1"
    }
    "europe-west1-nat" = {
      router_name = "europe-west1-router"
      nat_name    = "europe-west1-nat"
      region      = "europe-west1"
    }
  }
}
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| name | Name of the VPC network | `string` | n/a | yes |
| description | Description of the VPC network | `string` | `""` | no |
| routing_mode | Network-wide routing mode (REGIONAL or GLOBAL) | `string` | `"REGIONAL"` | no |
| subnets | Map of subnets to create | `map(object)` | n/a | yes |
| firewall_rules | List of firewall rules to create | `list(object)` | `[]` | no |
| enable_cloud_nat | Enable Cloud NAT for outbound connectivity | `bool` | `false` | no |
| cloud_nat_configs | Map of Cloud NAT configurations | `map(object)` | `{}` | no |

See [variables.tf](./variables.tf) for detailed variable descriptions and validation rules.

## Outputs

| Name | Description |
|------|-------------|
| vpc_id | ID of the VPC network |
| vpc_name | Name of the VPC network |
| vpc_self_link | Self-link of the VPC network |
| subnet_ids | Map of subnet names to their IDs |
| subnet_names | Map of subnet names to their full names |
| subnet_self_links | Map of subnet names to their self-links |
| subnet_cidr_ranges | Map of subnet names to their CIDR ranges |
| subnet_regions | Map of subnet names to their regions |
| cloud_router_ids | Map of Cloud Router configuration names to their IDs |
| cloud_nat_ids | Map of Cloud NAT configuration names to their IDs |
| firewall_rule_ids | Map of firewall rule names to their IDs |

## Requirements

- Terraform >= 1.0
- Google Cloud Provider 7.11.0

## Design Decisions

- **Zero-downtime updates**: All resources use `create_before_destroy` lifecycle for safe updates
- **Flexible firewall rules**: Supports both allow and deny rules with extensive configuration options
- **Network tags**: Firewall rules can target specific instances using network tags
- **Cloud NAT**: Optional per-region Cloud NAT for outbound connectivity
- **Private Google Access**: Enabled by default for subnets to access Google Cloud APIs
- **Logging**: Configurable firewall logging for troubleshooting

## Related Documentation

- [ADR-0005: Network Boot Infrastructure Implementation on Google Cloud](../../docs/content/r&d/adrs/0005-network-boot-infrastructure-gcp.md)
- [Google Cloud VPC Documentation](https://cloud.google.com/vpc/docs)
- [Cloud NAT Documentation](https://cloud.google.com/nat/docs)
- [VPC Firewall Rules](https://cloud.google.com/vpc/docs/firewalls)
