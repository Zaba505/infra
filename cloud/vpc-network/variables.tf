variable "name" {
  description = "Name of the VPC network"
  type        = string
}

variable "description" {
  description = "Description of the VPC network"
  type        = string
  default     = ""
}

variable "routing_mode" {
  description = "Network-wide routing mode. Valid values: REGIONAL, GLOBAL"
  type        = string
  default     = "REGIONAL"

  validation {
    condition     = contains(["REGIONAL", "GLOBAL"], var.routing_mode)
    error_message = "routing_mode must be either REGIONAL or GLOBAL"
  }
}

variable "subnets" {
  description = <<-EOT
    Map of subnets to create. Each subnet requires:
    - ip_cidr_range: CIDR range for the subnet
    - region: GCP region for the subnet
    Optional fields:
    - description: Subnet description
    - private_ip_google_access: Enable private Google access (default: true)
    - secondary_ip_ranges: List of secondary IP ranges for GKE/Cloud Run
      - range_name: Name of the secondary range
      - ip_cidr_range: CIDR range for the secondary range
  EOT
  type = map(object({
    ip_cidr_range            = string
    region                   = string
    description              = optional(string)
    private_ip_google_access = optional(bool)
    secondary_ip_ranges = optional(list(object({
      range_name    = string
      ip_cidr_range = string
    })))
  }))
}

variable "firewall_rules" {
  description = <<-EOT
    List of firewall rules to create. Each rule requires:
    - name: Name of the firewall rule
    - direction: INGRESS or EGRESS
    - allow or deny: List of protocol/port combinations
      - protocol: Protocol (tcp, udp, icmp, esp, ah, sctp, ipip, all)
      - ports: List of ports or port ranges (e.g., ["80", "443", "8080-8090"])
    
    For INGRESS rules:
    - source_ranges: List of source CIDR ranges (optional)
    - source_tags: List of source network tags (optional)
    
    For EGRESS rules:
    - destination_ranges: List of destination CIDR ranges (optional)
    
    Common optional fields:
    - description: Rule description
    - priority: Rule priority (default: 1000, lower is higher priority)
    - target_tags: List of target network tags
    - target_service_accounts: List of target service accounts
    - log_config: Logging configuration
      - metadata: EXCLUDE_ALL_METADATA, INCLUDE_ALL_METADATA
    
    Example:
    [
      {
        name      = "allow-wireguard"
        direction = "INGRESS"
        source_ranges = ["0.0.0.0/0"]
        target_tags = ["wireguard-gateway"]
        allow = [{
          protocol = "udp"
          ports    = ["51820"]
        }]
        log_config = {
          metadata = "INCLUDE_ALL_METADATA"
        }
      }
    ]
  EOT
  type = list(object({
    name                    = string
    direction               = string
    description             = optional(string)
    priority                = optional(number)
    source_ranges           = optional(list(string))
    source_tags             = optional(list(string))
    destination_ranges      = optional(list(string))
    target_tags             = optional(list(string))
    target_service_accounts = optional(list(string))
    allow = optional(list(object({
      protocol = string
      ports    = optional(list(string))
    })))
    deny = optional(list(object({
      protocol = string
      ports    = optional(list(string))
    })))
    log_config = optional(object({
      metadata = string
    }))
  }))
  default = []

  validation {
    condition = alltrue([
      for rule in var.firewall_rules : contains(["INGRESS", "EGRESS"], rule.direction)
    ])
    error_message = "firewall_rules direction must be either INGRESS or EGRESS"
  }
}

variable "enable_cloud_nat" {
  description = "Enable Cloud NAT for outbound connectivity"
  type        = bool
  default     = false
}

variable "cloud_nat_configs" {
  description = <<-EOT
    Map of Cloud NAT configurations. Each configuration requires:
    - router_name: Name of the Cloud Router
    - nat_name: Name of the Cloud NAT
    - region: GCP region for Cloud Router and NAT
    
    Optional fields:
    - asn: BGP ASN for the Cloud Router (default: 64514)
    - nat_ip_allocate_option: AUTO_ONLY or MANUAL_ONLY (default: AUTO_ONLY)
    - source_subnetwork_ip_ranges_to_nat: ALL_SUBNETWORKS_ALL_IP_RANGES, ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES, LIST_OF_SUBNETWORKS (default: ALL_SUBNETWORKS_ALL_IP_RANGES)
    - nat_ips: List of NAT IP addresses (required if nat_ip_allocate_option is MANUAL_ONLY)
    - enable_logging: Enable NAT logging (default: false)
    - log_filter: ERRORS_ONLY, TRANSLATIONS_ONLY, ALL (default: ERRORS_ONLY)
    - min_ports_per_vm: Minimum ports per VM (default: 64)
    
    Example:
    {
      "us-central1-nat" = {
        router_name = "us-central1-router"
        nat_name    = "us-central1-nat"
        region      = "us-central1"
        enable_logging = true
        log_filter  = "ERRORS_ONLY"
      }
    }
  EOT
  type = map(object({
    router_name                        = string
    nat_name                           = string
    region                             = string
    asn                                = optional(number)
    nat_ip_allocate_option             = optional(string)
    source_subnetwork_ip_ranges_to_nat = optional(string)
    nat_ips                            = optional(list(string))
    enable_logging                     = optional(bool)
    log_filter                         = optional(string)
    min_ports_per_vm                   = optional(number)
  }))
  default = {}
}
