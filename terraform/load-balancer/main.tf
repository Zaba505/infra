terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.0"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.2"
    }

    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = ">= 5.6.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.5"
    }
  }
}

module "copy_default_service_image_to_artifact_registry" {
  source = "../modules/copy_container_image"

  for_each = var.destination_registries

  source-image         = "ghcr.io/zaba505/infra/lb-sink:${var.default_service.image_tag}"
  destination-registry = each.value
}

locals {
  default_service_name = "lb-sink-service"
  default_service_location_cfgs = {
    for loc, reg in var.destination_registries : loc => {
      image = {
        name = module.copy_default_service_image_to_artifact_registry[loc].destination-image-name
        tag  = var.default_service.image_tag
      }
      location                    = loc
      cpu_limit                   = var.default_service.cpu_limit
      memory_limit                = var.default_service.memory_limit
      max_instance_count          = var.default_service.max_instance_count
      max_concurrent_requests     = var.default_service.max_concurrent_requests
      max_request_timeout_seconds = var.default_service.max_request_timeout_seconds
    }
  }

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
  source = "../modules/gcp/service_account"

  name        = "${local.default_service_name}-sa"
  cloud_trace = true
}

module "default_service" {
  source = "../modules/gcp/cloud_run"

  for_each = local.default_service_location_cfgs

  name                  = local.default_service_name
  description           = "Service for sinking all unknown requests to"
  service_account_email = module.default_service_account.service_account_email
  location              = each.key
  unsecured             = true

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
  default_service_locations = keys(local.default_service_location_cfgs)

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

locals {
  hostname = "machine.${var.domain_zone}"
}

resource "google_compute_url_map" "https" {
  name = "https"

  default_service = google_compute_backend_service.default_service.id

  host_rule {
    hosts        = [local.hostname]
    path_matcher = "https"
  }

  path_matcher {
    name = "https"

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

resource "google_certificate_manager_trust_config" "lb_https" {
  provider = google-beta

  name     = "lb-https-trust-config"
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

  name       = "lb-https-tls-policy"
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
    common_name = local.hostname
  }
}

resource "cloudflare_origin_ca_certificate" "lb_https" {
  csr                  = tls_cert_request.instance.cert_request_pem
  hostnames            = [local.hostname]
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
  name_prefix = "lb-https-ssl-cert-"

  certificate = cloudflare_origin_ca_certificate.lb_https.certificate
  private_key = tls_private_key.instance.private_key_pem

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_target_https_proxy" "lb_https" {
  provider = google-beta

  name              = "lb-https"
  url_map           = google_compute_url_map.https.id
  ssl_certificates  = [google_compute_ssl_certificate.lb_https.id]
  server_tls_policy = google_network_security_server_tls_policy.lb_https.id
}

resource "google_compute_global_address" "lb_https_ipv6" {
  name         = "lb-https-ipv6"
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}

resource "google_compute_global_forwarding_rule" "lb_https_ipv6" {
  name                  = "lb-https-ipv6"
  ip_address            = google_compute_global_address.lb_https_ipv6.id
  port_range            = "443"
  target                = google_compute_target_https_proxy.lb_https.id
  load_balancing_scheme = "EXTERNAL_MANAGED"
}