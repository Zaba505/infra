# External Network Load Balancer Module
#
# This module creates a regional external Network Load Balancer on GCP for exposing
# services to the internet. It supports both TCP and UDP protocols.
#
# Example Usage (WireGuard VPN Gateway):
#
#   resource "google_compute_address" "wireguard_ip" {
#     name         = "wireguard-ip"
#     region       = "us-central1"
#     address_type = "EXTERNAL"
#   }
#
#   module "wireguard_nlb" {
#     source = "./cloud/network-load-balancer"
#
#     name                = "wireguard-gateway"
#     region              = "us-central1"
#     external_ip_address = google_compute_address.wireguard_ip.name
#
#     protocols = ["UDP"]
#     port_range = {
#       start = 51820
#       end   = 51820
#     }
#
#     instance_groups = [{
#       instance_group = google_compute_instance_group.wireguard.id
#       balancing_mode = "CONNECTION"
#       health_check = {
#         protocol           = "TCP"
#         port               = 51820
#         check_interval_sec = 10
#         timeout_sec        = 5
#       }
#     }]
#   }
#
# Example Usage (HTTP Health Check):
#
#   resource "google_compute_address" "http_ip" {
#     name         = "http-ip"
#     region       = "us-central1"
#     address_type = "EXTERNAL"
#   }
#
#   module "http_nlb" {
#     source = "./cloud/network-load-balancer"
#
#     name                = "http-service"
#     region              = "us-central1"
#     external_ip_address = google_compute_address.http_ip.name
#
#     protocols = ["TCP"]
#     port_range = {
#       start = 80
#       end   = 80
#     }
#
#     instance_groups = [{
#       instance_group = google_compute_instance_group.http_servers.id
#       balancing_mode = "UTILIZATION"
#       max_utilization = 0.8
#       health_check = {
#         protocol           = "HTTP"
#         port               = 8080
#         request_path       = "/health"
#         check_interval_sec = 10
#         timeout_sec        = 5
#       }
#     }]
#   }

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.11.0"
    }
  }
}

# Reference existing external IP address (required)
data "google_compute_address" "external_ip" {
  name   = var.external_ip_address
  region = var.region
}

# Health checks for each instance group
resource "google_compute_region_health_check" "instance_group" {
  for_each = { for idx, ig in var.instance_groups : idx => ig }

  name   = "${var.name}-health-check-${each.key}"
  region = var.region

  check_interval_sec  = each.value.health_check.check_interval_sec
  timeout_sec         = each.value.health_check.timeout_sec
  healthy_threshold   = each.value.health_check.healthy_threshold
  unhealthy_threshold = each.value.health_check.unhealthy_threshold

  dynamic "tcp_health_check" {
    for_each = each.value.health_check.protocol == "TCP" ? [1] : []
    content {
      port = each.value.health_check.port
    }
  }

  dynamic "http_health_check" {
    for_each = each.value.health_check.protocol == "HTTP" ? [1] : []
    content {
      port         = each.value.health_check.port
      request_path = each.value.health_check.request_path
    }
  }

  dynamic "https_health_check" {
    for_each = each.value.health_check.protocol == "HTTPS" ? [1] : []
    content {
      port         = each.value.health_check.port
      request_path = each.value.health_check.request_path
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Backend service for the network load balancer
resource "google_compute_region_backend_service" "default" {
  name                  = var.name
  region                = var.region
  protocol              = var.backend_protocol
  load_balancing_scheme = "EXTERNAL"
  timeout_sec           = var.backend_timeout_sec
  health_checks         = [for hc in google_compute_region_health_check.instance_group : hc.id]

  dynamic "backend" {
    for_each = { for idx, ig in var.instance_groups : idx => ig }

    content {
      group                        = backend.value.instance_group
      balancing_mode               = backend.value.balancing_mode
      capacity_scaler              = try(backend.value.capacity_scaler, null)
      max_connections              = try(backend.value.max_connections, null)
      max_connections_per_instance = try(backend.value.max_connections_per_instance, null)
      max_rate                     = try(backend.value.max_rate, null)
      max_rate_per_instance        = try(backend.value.max_rate_per_instance, null)
      max_utilization              = try(backend.value.max_utilization, null)
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Forwarding rule for TCP traffic
resource "google_compute_forwarding_rule" "tcp" {
  count = contains(var.protocols, "TCP") ? 1 : 0

  name                  = "${var.name}-tcp"
  region                = var.region
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  backend_service       = google_compute_region_backend_service.default.id
  network_tier          = var.network_tier
  port_range            = var.port_range.start == var.port_range.end ? tostring(var.port_range.start) : join("-", [var.port_range.start, var.port_range.end])
  ip_address            = data.google_compute_address.external_ip.address

  lifecycle {
    create_before_destroy = true
  }
}

# Forwarding rule for UDP traffic
resource "google_compute_forwarding_rule" "udp" {
  count = contains(var.protocols, "UDP") ? 1 : 0

  name                  = "${var.name}-udp"
  region                = var.region
  ip_protocol           = "UDP"
  load_balancing_scheme = "EXTERNAL"
  backend_service       = google_compute_region_backend_service.default.id
  network_tier          = var.network_tier
  port_range            = var.port_range.start == var.port_range.end ? tostring(var.port_range.start) : join("-", [var.port_range.start, var.port_range.end])
  ip_address            = data.google_compute_address.external_ip.address

  lifecycle {
    create_before_destroy = true
  }
}
