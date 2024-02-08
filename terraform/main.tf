terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.19.0"
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

    random = {
      source  = "hashicorp/random"
      version = ">= 3.6.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.5"
    }
  }
}

data "google_client_config" "default" {}

resource "google_artifact_registry_repository" "docker" {
  for_each = toset(var.gcp_locations)

  format        = "DOCKER"
  repository_id = "docker-infra"
  description   = "container images"
  location      = each.value
  mode          = "STANDARD_REPOSITORY"
}

locals {
  destination_registries = {
    for loc in var.gcp_locations : loc => "${loc}-docker.pkg.dev/${data.google_client_config.default.project}/${google_artifact_registry_repository.docker[loc].name}"
  }
}

resource "random_uuid" "machine_boot_image_bucket_name" {
  for_each = toset(var.gcp_locations)
}

module "machine_image_bucket" {
  source = "./gcp/cloud_storage"

  for_each = toset(var.gcp_locations)

  bucket-name     = random_uuid.machine_boot_image_bucket_name[each.value].result
  bucket-location = each.value
}

module "machine_mgmt_service_sa" {
  source = "./gcp/service_account"
  depends_on = [
    module.machine_image_bucket
  ]

  name = "machine-mgmt-service-sa"

  cloud_trace = true

  cloud_storage = {
    buckets = { for loc in var.gcp_locations : loc => random_uuid.machine_boot_image_bucket_name[loc].result }
  }
}

module "copy_machine_mgmt_image_to_artifact_registry" {
  source = "./copy_container_image"

  for_each = local.destination_registries

  source-image         = "ghcr.io/zaba505/infra/machinemgmt:${var.machine_mgmt_service.image_tag}"
  destination-registry = each.value
}

module "machine_mgmt_service" {
  source = "./gcp/cloud_run"

  for_each = toset(var.gcp_locations)

  name                  = "machine-mgmt-service"
  description           = "Service for fetching machine boot images"
  service_account_email = module.machine_mgmt_service_sa.service_account_email
  location              = each.value

  image = {
    name = module.copy_machine_mgmt_image_to_artifact_registry[each.value].destination-image-name
    tag  = var.machine_mgmt_service.image_tag
  }

  cpu_limit    = var.machine_mgmt_service.cpu_limit
  memory_limit = var.machine_mgmt_service.memory_limit
  env_vars = [
    {
      name  = "LOG_LEVEL"
      value = "INFO"
    },
    {
      name  = "SERVICE_NAME"
      value = "machine-mgmt-service"
    },
    {
      name  = "SERVICE_VERSION"
      value = var.machine_mgmt_service.image_tag
    },
    {
      name  = "HTTP_PORT"
      value = "8080"
    },
    {
      name  = "IMAGE_BUCKET_NAME"
      value = random_uuid.machine_boot_image_bucket_name[each.value].result
    }
  ]
  max_instance_count      = var.machine_mgmt_service.max_instance_count
  max_concurrent_requests = var.machine_mgmt_service.max_concurrent_requests
}

module "copy_default_service_image_to_artifact_registry" {
  source = "./copy_container_image"

  for_each = local.destination_registries

  source-image         = "ghcr.io/zaba505/infra/lb-sink:${var.default_service.image_tag}"
  destination-registry = each.value
}

module "origin_ca_cert" {
  source = "./cloudflare/mtls"

  hostname = "machine.${var.domain_zone}"
}

data "cloudflare_origin_ca_root_certificate" "rsa" {
  algorithm = "rsa"
}

module "gateway" {
  source = "./gcp/gateway"

  domain = "machine.${var.domain_zone}"

  ca_certificate_pem = data.cloudflare_origin_ca_root_certificate.rsa.cert_pem

  lb_certificate = {
    pem         = module.origin_ca_cert.ca_certificate_pem
    private_key = module.origin_ca_cert.ca_private_key
  }

  default_service = [
    for loc, reg in local.destination_registries : {
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
  ]

  cloud_run = {
    "machinemgmt-service" = {
      locations = var.gcp_locations
      paths = [
        "/bootstrap/image"
      ]
    }
  }
}

module "client_ca_cert" {
  source = "./cloudflare/mtls"

  hostname = "machine.${var.domain_zone}"
}

module "cloudflare_dns" {
  source = "./cloudflare/dns"

  domain_name = var.domain_zone

  records = {
    machine = {
      // Only use IPV6 when proxying requests from Cloudflare to GCP
      ipv6 = {
        address = module.gateway.global_ipv6_address
      }

      certificate = {
        pem         = module.client_ca_cert.ca_certificate_pem
        private_key = module.client_ca_cert.ca_private_key
      }
    }
  }
}