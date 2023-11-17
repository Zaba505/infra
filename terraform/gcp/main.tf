terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

resource "google_artifact_registry_repository" "container_images" {
  description   = "Container images"
  repository_id = "docker-infra"
  location      = var.container-images-registry-location
  format        = "DOCKER"
  mode          = "STANDARD_REPOSITORY"
}

module "copy_machinemgmt_image" {
  source = "./modules/copy_container_image"

  source-image         = "ghcr.io/zaba505/infra/machinemgmt:${var.machine-image-service-image-tag}"
  destination-registry = "${google_artifact_registry_repository.container_images.location}-docker.pkg.dev/${var.gcp-project-id}/${google_artifact_registry_repository.container_images.repository_id}"
}

module "copy_lb_sink_image" {
  source = "./modules/copy_container_image"

  source-image         = "ghcr.io/zaba505/infra/lb-sink:${var.lb-sink-service-image-tag}"
  destination-registry = "${google_artifact_registry_repository.container_images.location}-docker.pkg.dev/${var.gcp-project-id}/${google_artifact_registry_repository.container_images.repository_id}"
}

module "storage" {
  source = "./modules/storage"

  boot-image-bucket-name     = var.boot-image-bucket-name
  boot-image-bucket-location = var.boot-image-bucket-location
}

module "machinemgmt" {
  source = "./modules/machinemgmt"
  depends_on = [
    module.copy_machinemgmt_image
  ]

  gcp-project-id = var.gcp-project-id

  boot-image-bucket-name = module.storage.bucket_name

  machine-image-service-account-id              = "machine-mgmt-sa"
  machine-image-service-image                   = module.copy_machinemgmt_image.destination-image
  machine-image-service-locations               = var.machine-image-service-locations
  machine-image-service-cpu-limit               = 1
  machine-image-service-memory-limit            = "512Mi"
  machine-image-service-env-vars                = var.machine-image-service-env-vars
  machine-image-service-max-instance-count      = var.machine-image-service-max-instance-count
  machine-image-service-max-concurrent-requests = var.machine-image-service-max-concurrent-requests
  machine-image-service-max-request-timeout     = var.machine-image-service-max-request-timeout
}

module "lb_sink_service" {
  source = "./modules/cloud_run_service"
  depends_on = [
    module.copy_lb_sink_image
  ]

  name        = "lb-sink"
  description = "Respond to all unmatched routes by the Load Balancer"

  image = module.copy_lb_sink_image.destination-image

  locations               = var.lb-sink-service-locations
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
    machinemgmt = {
      email = module.machinemgmt.service_account_email
    }
  }
}

module "gateway" {
  source = "./modules/gateway"
  depends_on = [
    module.lb_sink_service,
    module.machinemgmt
  ]

  domains = var.domains

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
        service_name = "vm-machine-image-service"
        locations    = module.machinemgmt.machine-image-service-locations
      }
    }
  }
}