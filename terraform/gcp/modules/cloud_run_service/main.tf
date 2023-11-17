terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

resource "google_service_account" "api" {
  account_id   = "${var.name}-sa"
  display_name = "${var.name}-sa"
  description  = "Machine Image Service Account for accessing the boot images storage bucket."
}

resource "google_cloud_run_v2_service" "api" {
  count = length(var.locations)

  name        = var.name
  description = var.description

  location = var.locations[count.index]
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.api.email

    containers {
      image = var.image

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