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

  for_each = local.default_service_negs

  name                  = each.value
  network_endpoint_type = "SERVERLESS"
  region                = each.key

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
      group = google_compute_region_network_endpoint_group.default_service[backend.key].id
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

resource "google_certificate_manager_trust_config" "instance" {
  provider = google-beta

  name     = "global-gateway-trust-config"
  location = "global"

  trust_stores {
    trust_anchors {
      pem_certificate = var.ca_certificate_pem
    }
    intermediate_cas {
      pem_certificate = var.ca_certificate_pem
    }
  }
}

resource "google_network_security_server_tls_policy" "instance" {
  provider = google-beta

  name       = "global-gateway-tls-policy"
  location   = "global"
  allow_open = false
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = google_certificate_manager_trust_config.instance.id
  }
}

// Using with Target HTTPS Proxies
//
// SSL certificates cannot be updated after creation. In order to apply
// the specified configuration, Terraform will destroy the existing
// resource and create a replacement. To effectively use an SSL
// certificate resource with a Target HTTPS Proxy resource, it's
// recommended to specify create_before_destroy in a lifecycle block.
// Either omit the Instance Template name attribute, specify a partial
// name with name_prefix, or use random_id resource. Example:
resource "google_compute_ssl_certificate" "global_gateway" {
  name_prefix = "global-gateway-ssl-cert-"

  certificate = var.lb_certificate.pem
  private_key = var.lb_certificate.private_key

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_target_https_proxy" "instance" {
  name              = "apis"
  url_map           = google_compute_url_map.apis.id
  ssl_certificates  = [google_compute_ssl_certificate.global_gateway.id]
  server_tls_policy = google_network_security_server_tls_policy.instance.id
}

resource "google_compute_global_address" "ipv6" {
  name         = "global-gateway-ipv6"
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}

resource "google_compute_global_forwarding_rule" "ipv6" {
  name                  = "apis-ipv6"
  ip_address            = google_compute_global_address.ipv6.id
  port_range            = "443"
  target                = google_compute_target_https_proxy.apis.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}