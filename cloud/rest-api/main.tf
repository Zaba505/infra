terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.11.0"
    }
  }
}

locals {
  default_env = {
    "LOG_LEVEL"       = "INFO",
    "SERVICE_NAME"    = var.name
    "SERVICE_VERSION" = var.image.tag
    "HTTP_PORT"       = "8080"
  }

  // since var.env appears later in the args,
  // then any keys in var.env will override the
  // values in local.default_env if the keys match
  envs = merge(local.default_env, var.env)
}

data "google_project" "default" {}

resource "google_cloud_run_v2_service" "rest_api" {
  name        = var.name
  description = var.description

  location = var.location
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    service_account = var.service_account_email

    containers {
      image = "${var.location}-docker.pkg.dev/${data.google_project.default.name}/${var.image.artifact_registry_name}/${var.image.name}:${var.image.tag}"

      resources {
        limits = {
          cpu    = var.cpu_limit
          memory = var.memory_limit
        }
        cpu_idle = false
      }

      dynamic "env" {
        for_each = local.envs
        content {
          name  = env.key
          value = env.value
        }
      }

      ports {
        container_port = coalesce(
          one([for k, v in local.envs : v if k == "HTTP_PORT"]),
          8080
        )
      }

      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 30
        period_seconds        = 10
        failure_threshold     = 3

        http_get {
          path = "/health/startup"
        }
      }

      liveness_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 30
        period_seconds        = 10
        failure_threshold     = 3

        http_get {
          path = "/health/liveness"
        }
      }
    }

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    timeout                          = "${var.max_request_timeout_seconds}s"
    max_instance_request_concurrency = var.max_concurrent_requests
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}

resource "google_cloud_run_v2_service_iam_binding" "default" {
  count = var.unsecured ? 1 : 0

  location = google_cloud_run_v2_service.rest_api.location
  name     = google_cloud_run_v2_service.rest_api.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}