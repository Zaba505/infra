terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.43.0"
    }

    google = {
      source  = "hashicorp/google"
      version = "6.5.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "6.4.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "4.0.6"
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

resource "google_certificate_manager_trust_config" "lb_https" {
  provider = google-beta

  name     = "${var.name}-trust-config"
  location = "global"

  trust_stores {
    dynamic "trust_anchors" {
      for_each = toset(var.ca_certificate_pems)

      content {
        pem_certificate = trust_anchors.value
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

resource "tls_private_key" "instance" {
  algorithm = "RSA"
}

resource "tls_cert_request" "instance" {
  private_key_pem = tls_private_key.instance.private_key_pem

  subject {
    # TODO: figure out to properly handle multiple host names
    common_name = local.hosts_to_cloud_run_services[0]
  }
}

resource "cloudflare_origin_ca_certificate" "lb_https" {
  # TODO: figure out to properly handle multiple host names and multiple CSRs
  csr                  = tls_cert_request.instance.cert_request_pem
  hostnames            = keys(local.hosts_to_cloud_run_services)
  request_type         = "origin-rsa"
  requested_validity   = 365
  min_days_for_renewal = 30
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
  name_prefix = "${var.name}-ssl-cert-"

  certificate = cloudflare_origin_ca_certificate.lb_https.certificate
  private_key = tls_private_key.instance.private_key_pem

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_target_https_proxy" "lb_https" {
  provider = google-beta

  name              = var.name
  url_map           = google_compute_url_map.https.id
  ssl_certificates  = [google_compute_ssl_certificate.lb_https.id]
  server_tls_policy = google_network_security_server_tls_policy.lb_https.id
}

resource "google_compute_global_address" "ipv6" {
  name         = "${var.name}-ipv6-addr"
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}

resource "google_compute_global_forwarding_rule" "lb_https_ipv6" {
  name                  = var.name
  ip_address            = google_compute_global_address.ipv6.address
  port_range            = "443"
  target                = google_compute_target_https_proxy.lb_https.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}