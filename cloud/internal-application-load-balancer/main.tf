terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.11.0"
    }
  }
}

# Create serverless NEGs for the default service in each region
locals {
  default_service_negs = {
    for loc in var.default_service.locations : loc => "${var.default_service.name}-${loc}-neg"
  }
}

resource "google_compute_region_network_endpoint_group" "default_service" {
  for_each = local.default_service_negs

  name                  = each.value
  network_endpoint_type = "SERVERLESS"
  region                = each.key

  cloud_run {
    service = var.default_service.name
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create health check for backend services
resource "google_compute_health_check" "default" {
  name = "${var.name}-health-check"

  timeout_sec         = var.health_check.timeout_sec
  check_interval_sec  = var.health_check.check_interval_sec
  healthy_threshold   = var.health_check.healthy_threshold
  unhealthy_threshold = var.health_check.unhealthy_threshold

  http_health_check {
    port         = 80
    request_path = var.health_check.request_path
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create backend service for the default service
resource "google_compute_backend_service" "default_service" {
  name                  = var.default_service.name
  protocol              = var.enable_https ? "HTTPS" : "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = var.backend_timeout_seconds
  health_checks         = [google_compute_health_check.default.id]

  dynamic "backend" {
    for_each = local.default_service_negs

    content {
      group = google_compute_region_network_endpoint_group.default_service[backend.key].id
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create serverless NEGs for additional Cloud Run services
locals {
  cloud_run_region_neg_configs = merge([
    for name, cfg in var.cloud_run : {
      for loc in cfg.locations : "${name}-${loc}-neg" => {
        location     = loc
        service_name = name
      }
    }
  ]...)
}

resource "google_compute_region_network_endpoint_group" "cloud_run" {
  for_each = local.cloud_run_region_neg_configs

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.value.location

  cloud_run {
    service = each.value.service_name
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create backend services for additional Cloud Run services
resource "google_compute_backend_service" "cloud_run" {
  for_each = var.cloud_run

  name                  = each.key
  protocol              = var.enable_https ? "HTTPS" : "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = var.backend_timeout_seconds
  health_checks         = [google_compute_health_check.default.id]

  dynamic "backend" {
    for_each = toset(each.value.locations)

    content {
      group = google_compute_region_network_endpoint_group.cloud_run["${each.key}-${backend.value}-neg"].id
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create URL map with host rules and path matchers
locals {
  hosts_to_cloud_run_services = transpose({ for name, cr in var.cloud_run : name => cr.hosts })
}

resource "google_compute_region_url_map" "default" {
  name   = var.name
  region = var.region

  default_service = google_compute_backend_service.default_service.id

  dynamic "host_rule" {
    for_each = local.hosts_to_cloud_run_services

    content {
      hosts        = [host_rule.key]
      path_matcher = replace(host_rule.key, ".", "-")
    }
  }

  dynamic "path_matcher" {
    for_each = local.hosts_to_cloud_run_services

    content {
      name = replace(path_matcher.key, ".", "-")

      default_service = google_compute_backend_service.default_service.id

      dynamic "path_rule" {
        for_each = path_matcher.value

        content {
          service = google_compute_backend_service.cloud_run[path_rule.value].id

          paths = var.cloud_run[path_rule.value].paths
        }
      }
    }
  }
}

# Create target proxy (HTTP or HTTPS based on configuration)
resource "google_compute_region_target_http_proxy" "default" {
  count = var.enable_https ? 0 : 1

  name    = var.name
  region  = var.region
  url_map = google_compute_region_url_map.default.id
}

resource "google_compute_region_target_https_proxy" "default" {
  count = var.enable_https ? 1 : 0

  name             = var.name
  region           = var.region
  url_map          = google_compute_region_url_map.default.id
  ssl_certificates = var.ssl_certificates
}

# Create internal forwarding rule
resource "google_compute_forwarding_rule" "default" {
  name = var.name

  region                = var.region
  network               = var.network
  subnetwork            = var.subnetwork
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = var.enable_https ? "443" : "80"
  ip_protocol           = "TCP"
  target = var.enable_https ? (
    google_compute_region_target_https_proxy.default[0].id
    ) : (
    google_compute_region_target_http_proxy.default[0].id
  )
}
