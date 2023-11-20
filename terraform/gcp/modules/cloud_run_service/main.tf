terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

locals {
  artifact_registry_locations = {
    for loc in var.locations : loc =>
    startswith(loc, "us") || startswith(loc, "europe") || startswith(loc, "asia") ? split("-", loc)[0] : loc
  }
}

data "google_client_config" "default" {}

module "copy_image_to_gcr" {
  for_each = toset(var.locations)

  source = "../copy_container_image"

  source-image         = "${var.image.name}:${var.image.tag}"
  destination-registry = "${local.artifact_registry_locations[each.value]}-docker.pkg.dev/${data.google_client_config.default.project}/${var.artifact_registry_id}"
}

resource "google_service_account" "api" {
  account_id   = "${var.name}-sa"
  display_name = "${var.name}-sa"
}

resource "google_project_iam_member" "cloud_trace" {
  project = data.google_client_config.default.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.api.email}"
}

resource "google_project_iam_member" "cloud_storage" {
  count = var.access.cloud_storage != null ? 1 : 0

  project = data.google_client_config.default.project
  role    = "roles/run.serviceAgent"
  member  = "serviceAccount:${google_service_account.api.email}"

  condition {
    title      = "only_cloud_storage"
    expression = <<-EOL
      resource.service == 'storage.googleapis.com'
      && (resource.type == 'storage.googleapis.com/Bucket' || resource.type == 'storage.googleapis.com/Object')
      && resource.name.startsWith('projects/_/buckets/${var.access.cloud_storage.bucket_name}')
    EOL
  }
}

resource "google_cloud_run_v2_service" "api" {
  for_each = toset(var.locations)

  name        = var.name
  description = var.description

  location = each.value
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    service_account = google_service_account.api.email

    containers {
      image = module.copy_image_to_gcr[each.value].destination-image

      resources {
        limits = {
          cpu    = var.cpu_limit
          memory = var.memory_limit
        }
        cpu_idle = false
      }

      dynamic "env" {
        for_each = var.env_vars
        content {
          name  = env.value["name"]
          value = env.value["value"]
        }
      }

      ports {
        container_port = one([for env_var in var.env_vars : env_var.value if env_var.name == "HTTP_PORT"])
      }

      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 1
        period_seconds        = 10
        failure_threshold     = 3

        http_get {
          path = "/health/startup"
        }
      }

      liveness_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 1
        period_seconds        = 10
        failure_threshold     = 3

        http_get {
          path = "/health/liveness"
        }
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = var.max_instance_count
    }

    timeout                          = "${var.max_request_timeout}s"
    max_instance_request_concurrency = var.max_concurrent_requests
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}
