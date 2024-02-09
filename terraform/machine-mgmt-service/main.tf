terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.2"
    }

    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }

    random = {
      source  = "hashicorp/random"
      version = ">= 3.6.0"
    }
  }
}

resource "random_uuid" "machine_boot_image_bucket_name" {
  for_each = toset(var.gcp_locations)
}

resource "google_storage_bucket" "boot_images" {
  for_each = toset(var.gcp_locations)

  name     = random_uuid.machine_boot_image_bucket_name[each.value].result
  location = each.value

  force_destroy            = true
  public_access_prevention = "enforced"

  autoclass {
    enabled = true
  }

  versioning {
    enabled = true
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      days_since_noncurrent_time = 7
    }
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      num_newer_versions = 1
    }
  }
}

module "machine_mgmt_service_sa" {
  source = "../modules/gcp/service_account"

  name = "machine-mgmt-service-sa"

  cloud_trace = true

  cloud_storage = {
    buckets = { for loc in var.gcp_locations : loc => google_storage_bucket.boot_images[loc].name }
  }
}

module "copy_machine_mgmt_image_to_artifact_registry" {
  source = "../modules/copy_container_image"

  for_each = var.destination_registries

  source-image         = "ghcr.io/zaba505/infra/machinemgmt:${var.image_tag}"
  destination-registry = each.value
}

module "machine_mgmt_service" {
  source = "../modules/gcp/cloud_run"

  for_each = toset(var.gcp_locations)

  name                  = "machine-mgmt-service"
  description           = "Service for fetching machine boot images"
  service_account_email = module.machine_mgmt_service_sa.service_account_email
  location              = each.value

  image = {
    name = module.copy_machine_mgmt_image_to_artifact_registry[each.value].destination-image-name
    tag  = var.image_tag
  }

  cpu_limit    = var.cpu_limit
  memory_limit = var.memory_limit
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
      value = var.image_tag
    },
    {
      name  = "HTTP_PORT"
      value = "8080"
    },
    {
      name  = "IMAGE_BUCKET_NAME"
      value = google_storage_bucket.boot_images[each.value].name
    }
  ]
  max_instance_count          = var.max_instance_count
  max_concurrent_requests     = var.max_concurrent_requests
  max_request_timeout_seconds = var.max_request_timeout_seconds
}