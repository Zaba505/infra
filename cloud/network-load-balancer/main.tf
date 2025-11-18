# External Network Load Balancer Module
#
# This module creates a regional external Network Load Balancer on GCP for exposing
# services to the internet. It supports both TCP and UDP protocols (required for WireGuard).
#
# Example Usage (WireGuard VPN Gateway):
#
#   module "wireguard_nlb" {
#     source = "./cloud/network-load-balancer"
#
#     name   = "wireguard-gateway"
#     region = "us-central1"
#
#     protocols = ["TCP", "UDP"]
#     port_range = {
#       start = 51820
#       end   = 51820
#     }
#
#     instance_groups = [{
#       instance_group = google_compute_instance_group.wireguard.id
#       balancing_mode = "CONNECTION"
#     }]
#
#     health_check = {
#       protocol           = "TCP"
#       port               = 51820
#       check_interval_sec = 10
#       timeout_sec        = 5
#     }
#   }
#
# Example Usage (Static External IP):
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
#     protocols = ["TCP", "UDP"]
#     port_range = {
#       start = 51820
#       end   = 51820
#     }
#
#     instance_groups = [{
#       instance_group = google_compute_instance_group.wireguard.id
#       balancing_mode = "CONNECTION"
#     }]
#   }
#
# Example Usage (HTTP Health Check):
#
#   module "http_nlb" {
#     source = "./cloud/network-load-balancer"
#
#     name   = "http-service"
#     region = "us-central1"
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
#     }]
#
#     health_check = {
#       protocol           = "HTTP"
#       port               = 8080
#       request_path       = "/health"
#       check_interval_sec = 10
#       timeout_sec        = 5
#     }
#   }

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.11.0"
    }
  }
}

# Optional: Create or use existing external IP address
resource "google_compute_address" "external_ip" {
  count = var.external_ip_address == null ? 1 : 0

  name         = "${var.name}-ip"
  region       = var.region
  address_type = "EXTERNAL"
  network_tier = var.network_tier

  lifecycle {
    create_before_destroy = true
  }
}

# Reference existing external IP if provided
data "google_compute_address" "external_ip" {
  count = var.external_ip_address != null ? 1 : 0

  name   = var.external_ip_address
  region = var.region
}

# Health check for backend instances
resource "google_compute_health_check" "default" {
  name                = "${var.name}-health-check"
  check_interval_sec  = var.health_check.check_interval_sec
  timeout_sec         = var.health_check.timeout_sec
  healthy_threshold   = var.health_check.healthy_threshold
  unhealthy_threshold = var.health_check.unhealthy_threshold

  dynamic "tcp_health_check" {
    for_each = var.health_check.protocol == "TCP" ? [1] : []
    content {
      port = var.health_check.port
    }
  }

  dynamic "http_health_check" {
    for_each = var.health_check.protocol == "HTTP" ? [1] : []
    content {
      port         = var.health_check.port
      request_path = var.health_check.request_path
    }
  }

  dynamic "https_health_check" {
    for_each = var.health_check.protocol == "HTTPS" ? [1] : []
    content {
      port         = var.health_check.port
      request_path = var.health_check.request_path
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
  health_checks         = [google_compute_health_check.default.id]

  dynamic "backend" {
    for_each = var.instance_groups

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
  port_range            = join("-", [var.port_range.start, var.port_range.end])
  ip_address = (
    var.external_ip_address != null
    ? data.google_compute_address.external_ip[0].address
    : google_compute_address.external_ip[0].address
  )

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
  port_range            = join("-", [var.port_range.start, var.port_range.end])
  ip_address = (
    var.external_ip_address != null
    ? data.google_compute_address.external_ip[0].address
    : google_compute_address.external_ip[0].address
  )

  lifecycle {
    create_before_destroy = true
  }
}
