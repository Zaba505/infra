terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

locals {
  api_region_neg_configs = merge([
    for name, api in var.apis : {
      for region in api.cloud_run.locations : "${name}-${region}-neg" => {
        region       = region
        service_name = api.cloud_run.service_name
      }
    }
  ]...)
}

resource "google_compute_global_address" "ipv4" {
  name         = "global-gateway-ipv4"
  ip_version   = "IPV4"
  address_type = "EXTERNAL"
}

resource "google_compute_global_address" "ipv6" {
  name         = "global-gateway-ipv6"
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}

resource "google_compute_managed_ssl_certificate" "global_gateway" {
  name = "global-gateway-ssl-cert"

  managed {
    domains = var.domains
  }
}

resource "google_compute_region_network_endpoint_group" "default" {
  for_each = { for loc in var.default_service.locations : "${var.default_service.name}-${loc}-neg" => loc }

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.value

  cloud_run {
    service = var.default_service.name
  }
}

resource "google_compute_backend_service" "default" {
  name                  = var.default_service.name
  protocol              = "HTTPS"
  load_balancing_scheme = "EXTERNAL_MANAGED"

  dynamic "backend" {
    for_each = toset(var.default_service.locations)

    content {
      group = google_compute_region_network_endpoint_group.default["${var.default_service.name}-${backend.value}-neg"].id
    }
  }
}

resource "google_compute_region_network_endpoint_group" "api" {
  for_each = local.api_region_neg_configs

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.value.region

  cloud_run {
    service = each.value.service_name
  }
}

resource "google_compute_backend_service" "api" {
  for_each = var.apis

  name                  = each.key
  protocol              = "HTTPS"
  load_balancing_scheme = "EXTERNAL_MANAGED"

  dynamic "backend" {
    for_each = each.value.cloud_run.locations

    content {
      group = google_compute_region_network_endpoint_group.api["${each.key}-${backend.value}-neg"].id
    }
  }
}

resource "google_compute_url_map" "apis" {
  name = "apis"

  default_service = google_compute_backend_service.default.id

  host_rule {
    hosts        = var.domains
    path_matcher = "apis"
  }

  path_matcher {
    name = "apis"

    default_service = google_compute_backend_service.default.id

    dynamic "path_rule" {
      for_each = var.apis

      content {
        paths   = path_rule.value["paths"]
        service = google_compute_backend_service.api[path_rule.key].id
      }
    }
  }
}

resource "google_compute_target_https_proxy" "apis" {
  name             = "apis"
  url_map          = google_compute_url_map.apis.id
  ssl_certificates = [google_compute_managed_ssl_certificate.global_gateway.id]
}

resource "google_compute_global_forwarding_rule" "ipv4" {
  name                  = "apis"
  ip_address            = google_compute_global_address.ipv4.id
  port_range            = "443"
  target                = google_compute_target_https_proxy.apis.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}

resource "google_compute_global_forwarding_rule" "ipv6" {
  name                  = "apis"
  ip_address            = google_compute_global_address.ipv6.id
  port_range            = "443"
  target                = google_compute_target_https_proxy.apis.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}