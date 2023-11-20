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

    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

locals {
  artifact_registry_locations = values({
    for loc in var.locations : loc =>
    startswith(loc, "us") || startswith(loc, "europe") || startswith(loc, "asia") ? split("-", loc)[0] : loc
  })
}

resource "google_artifact_registry_repository" "container_images" {
  for_each = toset(local.artifact_registry_locations)

  description   = "Container images"
  repository_id = "docker-infra"
  location      = each.value
  format        = "DOCKER"
  mode          = "STANDARD_REPOSITORY"
}

module "storage" {
  source = "./modules/storage"

  boot-image-bucket-name     = var.boot-image-bucket-name
  boot-image-bucket-location = var.boot-image-bucket-location
}

module "machine_image_service" {
  source = "./modules/cloud_run_service"
  depends_on = [
    google_artifact_registry_repository.container_images
  ]

  artifact_registry_id = "docker-infra"

  access = {
    cloud_storage = {
      bucket_name = module.storage.bucket_name
    }
  }

  name        = "machine-image-service"
  description = "API service for fetching machine boot images"

  image = {
    name = "ghcr.io/zaba505/infra/machinemgmt"
    tag  = var.machine-image-service-image-tag
  }

  locations               = var.locations
  cpu_limit               = 1
  memory_limit            = "512Mi"
  env_vars                = var.machine-image-service-env-vars
  max_instance_count      = var.machine-image-service-max-instance-count
  max_concurrent_requests = var.machine-image-service-max-concurrent-requests
  max_request_timeout     = var.machine-image-service-max-request-timeout
}

module "lb_sink_service" {
  source = "./modules/cloud_run_service"
  depends_on = [
    google_artifact_registry_repository.container_images
  ]

  artifact_registry_id = "docker-infra"

  name        = "lb-sink-service"
  description = "Respond to all unmatched routes by the Load Balancer"

  # this service is unauthenticated so people don't know that the request
  # is even making it to a service. The service will immediately return a 503
  unauthenticated = true

  image = {
    name = "ghcr.io/zaba505/infra/lb-sink"
    tag  = var.lb-sink-service-image-tag
  }

  locations               = var.locations
  cpu_limit               = 1
  memory_limit            = "512Mi"
  env_vars                = var.lb-sink-service-env-vars
  max_instance_count      = var.lb-sink-service-max-instance-count
  max_concurrent_requests = var.lb-sink-service-max-concurrent-requests
  max_request_timeout     = var.lb-sink-service-max-request-timeout
}

module "access_control" {
  source = "./modules/access_control"

  boot-image-storage-bucket-name = module.storage.bucket_name
  boot-image-service-accounts = {
    machine_image_service = {
      email = module.machine_image_service.service_account_email
    }
  }
}

module "rootca" {
  source = "./modules/rootca"

  organization_name = var.organization_name
  domains           = var.domains
}

module "lb_cacert" {
  source = "./modules/cacert"

  name                     = "global-lb-cacert"
  organization             = var.organization_name
  common_name              = ""
  location                 = "us-east1"
  algorithm                = "RSA"
  lifetime                 = "31104000s"
  privateca_pool_name      = module.rootca.ca_pool_name
  certificate_authority_id = module.rootca.certificate_authority_id
}

module "client_cacert" {
  source   = "./modules/cacert"
  for_each = var.clients

  name                     = each.key
  organization             = var.organization_name
  common_name              = var.common_name
  location                 = "us-east1"
  algorithm                = each.value.algorithm
  lifetime                 = each.value.lifetime
  privateca_pool_name      = module.rootca.ca_pool_name
  certificate_authority_id = module.rootca.certificate_authority_id
}

module "gateway" {
  source = "./modules/gateway"
  depends_on = [
    module.lb_sink_service,
    module.machine_image_service
  ]

  domains = var.domains

  root_pem_certificate = module.lb_cacert.pem_certificate

  default_service = {
    name      = module.lb_sink_service.name
    locations = module.lb_sink_service.locations
  }

  apis = {
    "machine-image-service" = {
      paths = [
        "/bootstrap/image"
      ]
      cloud_run = {
        service_name = module.machine_image_service.name
        locations    = module.machine_image_service.locations
      }
    }
  }
}