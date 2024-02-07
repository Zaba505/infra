terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = ">= 5.6.0"
    }
  }
}

locals {
  default_service_name = "lb-sink-service"

  cloud_run_region_neg_configs = merge([
    for name, cfg in var.cloud_run : {
      for loc in cfg.locations : "${name}-${loc}-neg" => {
        location     = loc
        service_name = name
      }
    }
  ]...)
}

resource "google_compute_global_address" "ipv6" {
  name         = "global-gateway-ipv6"
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}

resource "google_compute_managed_ssl_certificate" "global_gateway" {
  name = "global-gateway-ssl-cert"

  managed {
    domains = [var.domain]
  }
}

module "default_service_account" {
  source = "../service_account"

  name        = "${local.default_service_name}-sa"
  cloud_trace = true
}

module "default_service" {
  source = "../cloud_run"

  for_each = { for v in var.default_service : v.location => v }

  name                  = local.default_service_name
  description           = "Service for sinking all unknown requests to"
  service_account_email = module.default_service_account.service_account_email
  location              = each.key

  image = {
    name = each.value.image.name
    tag  = each.value.image.tag
  }

  cpu_limit    = each.value.cpu_limit
  memory_limit = each.value.memory_limit
  env_vars = [
    {
      name  = "LOG_LEVEL"
      value = "ERROR"
    },
    {
      name  = "SERVICE_NAME"
      value = local.default_service_name
    },
    {
      name  = "SERVICE_VERSION",
      value = each.value.image.tag
    },
    {
      name  = "HTTP_PORT",
      value = "8080"
    }
  ]
  max_instance_count          = each.value.max_instance_count
  max_concurrent_requests     = each.value.max_concurrent_requests
  max_request_timeout_seconds = each.value.max_request_timeout_seconds
}

locals {
  default_service_locations = [for v in var.default_service : v.location]

  default_service_negs = {
    for loc in local.default_service_locations : loc => "${local.default_service_name}-${loc}-neg"
  }
}

resource "google_compute_region_network_endpoint_group" "default_service" {
  depends_on = [module.default_service]

  for_each = toset(local.default_service_locations)

  name                  = local.default_service_negs[each.value]
  network_endpoint_type = "SERVERLESS"
  region                = each.value

  cloud_run {
    service = local.default_service_name
  }
}

resource "google_compute_backend_service" "default_service" {
  name                  = local.default_service_name
  protocol              = "HTTPS"
  load_balancing_scheme = "EXTERNAL_MANAGED"

  dynamic "backend" {
    for_each = local.default_service_negs

    content {
      group = google_compute_region_network_endpoint_group.api[backend.value].id
    }
  }
}

resource "google_compute_region_network_endpoint_group" "cloud_run" {
  for_each = local.cloud_run_region_neg_configs

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.value.location

  cloud_run {
    service = each.value.service_name
  }
}

resource "google_compute_backend_service" "cloud_run" {
  for_each = var.cloud_run

  name                  = each.key
  protocol              = "HTTPS"
  load_balancing_scheme = "EXTERNAL_MANAGED"

  dynamic "backend" {
    for_each = toset(each.value.locations)

    content {
      group = google_compute_region_network_endpoint_group.cloud_run["${each.key}-${backend.value}-neg"].id
    }
  }
}

resource "google_compute_url_map" "apis" {
  name = "apis"

  default_service = google_compute_backend_service.default_service.id

  host_rule {
    hosts        = [var.domain]
    path_matcher = "apis"
  }

  path_matcher {
    name = "apis"

    default_service = google_compute_backend_service.default_service.id

    dynamic "path_rule" {
      for_each = var.cloud_run

      content {
        paths   = path_rule.value["paths"]
        service = google_compute_backend_service.cloud_run[path_rule.key].id
      }
    }
  }
}

resource "google_compute_target_https_proxy" "apis" {
  name             = "apis"
  url_map          = google_compute_url_map.apis.id
  ssl_certificates = [google_compute_managed_ssl_certificate.global_gateway.id]
}

resource "google_compute_global_forwarding_rule" "ipv6" {
  name                  = "apis-ipv6"
  ip_address            = google_compute_global_address.ipv6.id
  port_range            = "443"
  target                = google_compute_target_https_proxy.apis.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}