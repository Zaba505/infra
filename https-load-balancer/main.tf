terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.10.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "6.10.0"
    }
  }
}

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
}

resource "google_compute_backend_service" "default_service" {
  name                  = var.default_service.name
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

locals {
  hosts_to_cloud_run_services = transpose({ for name, cr in var.cloud_run : name => cr.hosts })
}

resource "google_compute_url_map" "https" {
  name = "https"

  default_service = google_compute_backend_service.default_service.id

  dynamic "host_rule" {
    for_each = local.hosts_to_cloud_run_services

    content {
      hosts        = [host_rule.key]
      path_matcher = host_rule.key
    }
  }

  dynamic "path_matcher" {
    for_each = local.hosts_to_cloud_run_services

    content {
      name = path_matcher.key

      default_service = google_compute_backend_service.default_service.id

      dynamic "path_rule" {
        for_each = path_matcher.value

        content {
          service = google_compute_backend_service.cloud_run[path_rule.key].id

          paths = var.cloud_run[path_rule.key].paths
        }
      }
    }
  }
}

locals {
  trust_anchor_secrets = toset(var.trust_anchor_secrets)
}

data "google_secret_manager_secret_version_access" "trust_anchor" {
  for_each = local.trust_anchor_secrets

  secret = each.value.secret
  version = each.value.version
}

resource "google_certificate_manager_trust_config" "lb_https" {
  provider = google-beta

  name     = "${var.name}-trust-config"
  location = "global"

  trust_stores {
    dynamic "trust_anchors" {
      for_each = local.trust_anchor_secrets

      content {
        pem_certificate = data.google_secret_manager_secret_version_access.trust_anchor[trust_anchors.key].secret_data
      }
    }
  }
}

resource "google_network_security_server_tls_policy" "lb_https" {
  provider = google-beta

  name       = "${var.name}-tls-policy"
  location   = "global"
  allow_open = false
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = google_certificate_manager_trust_config.lb_https.id
  }
}

locals {
  server_certificate_secrets = toset(var.server_certificate_secrets)
}

data "google_secret_manager_secret_version_access" "server_certificate" {
  for_each = local.server_certificate_secrets

  secret = each.value.certificate_secret
  version = each.value.certificate_version
}

data "google_secret_manager_secret_version_access" "server_private_key" {
  for_each = local.server_certificate_secrets

  secret = each.value.private_key_secret
  version = each.value.private_key_version
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
resource "google_compute_ssl_certificate" "lb_https" {
  for_each = local.server_certificate_secrets

  name_prefix = "${each.value.certificate_secret}-"

  certificate = data.google_secret_manager_secret_version_access.server_certificate[each.key].secret_data
  private_key = data.google_secret_manager_secret_version_access.server_private_key[each.key].secret_data

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_target_https_proxy" "lb_https" {
  provider = google-beta

  name              = var.name
  url_map           = google_compute_url_map.https.id
  ssl_certificates  = google_compute_ssl_certificate.lb_https[*].id
  server_tls_policy = google_network_security_server_tls_policy.lb_https.id
}

locals {
  ip_addresses = toset(var.ip_addresses)
}

data "google_compute_global_address" "ip" {
  for_each = local.ip_addresses

  name = each.value.name
}

resource "google_compute_global_forwarding_rule" "lb_https_ipv6" {
  for_each = local.ip_addresses

  name                  = var.name
  ip_address            = data.google_compute_global_address.ip[each.key].address
  port_range            = "443"
  target                = google_compute_target_https_proxy.lb_https.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}