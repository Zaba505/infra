terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.12.0"
    }
  }
}

# VPC Network
resource "google_compute_network" "vpc" {
  name                    = var.name
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
  description             = var.description

  lifecycle {
    create_before_destroy = true
  }
}

# Subnets
resource "google_compute_subnetwork" "subnets" {
  for_each = var.subnets

  name          = each.key
  ip_cidr_range = each.value.ip_cidr_range
  region        = each.value.region
  network       = google_compute_network.vpc.id
  description   = lookup(each.value, "description", "Subnet ${each.key}")

  # Secondary IP ranges for GKE/Cloud Run (optional)
  dynamic "secondary_ip_range" {
    for_each = lookup(each.value, "secondary_ip_ranges", [])
    content {
      range_name    = secondary_ip_range.value.range_name
      ip_cidr_range = secondary_ip_range.value.ip_cidr_range
    }
  }

  # Private Google Access for Cloud APIs
  private_ip_google_access = lookup(each.value, "private_ip_google_access", true)

  lifecycle {
    create_before_destroy = true
  }
}

# Firewall Rules - Ingress
resource "google_compute_firewall" "ingress" {
  for_each = { for rule in var.firewall_rules : rule.name => rule if rule.direction == "INGRESS" }

  name        = each.value.name
  network     = google_compute_network.vpc.id
  description = lookup(each.value, "description", "Ingress rule ${each.value.name}")
  direction   = "INGRESS"
  priority    = lookup(each.value, "priority", 1000)

  # Source configuration
  source_ranges = lookup(each.value, "source_ranges", [])
  source_tags   = lookup(each.value, "source_tags", [])

  # Target configuration
  target_tags             = lookup(each.value, "target_tags", [])
  target_service_accounts = lookup(each.value, "target_service_accounts", [])

  # Allow/Deny rules
  dynamic "allow" {
    for_each = lookup(each.value, "allow", [])
    content {
      protocol = allow.value.protocol
      ports    = lookup(allow.value, "ports", [])
    }
  }

  dynamic "deny" {
    for_each = lookup(each.value, "deny", [])
    content {
      protocol = deny.value.protocol
      ports    = lookup(deny.value, "ports", [])
    }
  }

  # Logging configuration
  dynamic "log_config" {
    for_each = lookup(each.value, "log_config", null) != null ? [each.value.log_config] : []
    content {
      metadata = log_config.value.metadata
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Firewall Rules - Egress
resource "google_compute_firewall" "egress" {
  for_each = { for rule in var.firewall_rules : rule.name => rule if rule.direction == "EGRESS" }

  name        = each.value.name
  network     = google_compute_network.vpc.id
  description = lookup(each.value, "description", "Egress rule ${each.value.name}")
  direction   = "EGRESS"
  priority    = lookup(each.value, "priority", 1000)

  # Destination configuration
  destination_ranges = lookup(each.value, "destination_ranges", [])

  # Target configuration
  target_tags             = lookup(each.value, "target_tags", [])
  target_service_accounts = lookup(each.value, "target_service_accounts", [])

  # Allow/Deny rules
  dynamic "allow" {
    for_each = lookup(each.value, "allow", [])
    content {
      protocol = allow.value.protocol
      ports    = lookup(allow.value, "ports", [])
    }
  }

  dynamic "deny" {
    for_each = lookup(each.value, "deny", [])
    content {
      protocol = deny.value.protocol
      ports    = lookup(deny.value, "ports", [])
    }
  }

  # Logging configuration
  dynamic "log_config" {
    for_each = lookup(each.value, "log_config", null) != null ? [each.value.log_config] : []
    content {
      metadata = log_config.value.metadata
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Cloud Router for Cloud NAT
resource "google_compute_router" "router" {
  for_each = var.enable_cloud_nat ? var.cloud_nat_configs : {}

  name    = each.value.router_name
  region  = each.value.region
  network = google_compute_network.vpc.id

  bgp {
    asn = lookup(each.value, "asn", 64514)
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Cloud NAT for outbound connectivity
resource "google_compute_router_nat" "nat" {
  for_each = var.enable_cloud_nat ? var.cloud_nat_configs : {}

  name   = each.value.nat_name
  router = google_compute_router.router[each.key].name
  region = each.value.region

  nat_ip_allocate_option             = lookup(each.value, "nat_ip_allocate_option", "AUTO_ONLY")
  source_subnetwork_ip_ranges_to_nat = lookup(each.value, "source_subnetwork_ip_ranges_to_nat", "ALL_SUBNETWORKS_ALL_IP_RANGES")

  # NAT IP addresses (if using MANUAL_ONLY)
  nat_ips = lookup(each.value, "nat_ips", [])

  # Logging configuration
  dynamic "log_config" {
    for_each = lookup(each.value, "enable_logging", false) ? [1] : []
    content {
      enable = true
      filter = lookup(each.value, "log_filter", "ERRORS_ONLY")
    }
  }

  # Min ports per VM
  min_ports_per_vm = lookup(each.value, "min_ports_per_vm", 64)

  lifecycle {
    create_before_destroy = true
  }
}
