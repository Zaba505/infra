terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

resource "google_service_account" "machine_image_service" {
  account_id   = var.machine-image-service-account-id
  display_name = "machine-image-sa"
  description  = "Machine Image Service Account for accessing the boot images storage bucket."
}

resource "google_project_iam_member" "machine_image_service_cloud_trace" {
  project = var.gcp-project-id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.machine_image_service.email}"
}

resource "google_project_iam_member" "machine_image_service_cloud_storage" {
  project = var.gcp-project-id
  role    = "roles/run.serviceAgent"
  member  = "serviceAccount:${google_service_account.machine_image_service.email}"

  condition {
    title      = "only_cloud_storage"
    expression = <<-EOL
      resource.service == 'storage.googleapis.com'
      && (resource.type == 'storage.googleapis.com/Bucket' || resource.type == 'storage.googleapis.com/Object')
      && resource.name.startsWith('projects/_/buckets/${var.boot-image-bucket-name}')
    EOL
  }
}

resource "google_cloud_run_v2_service" "machine_image_service" {
  count = length(var.machine-image-service-locations)

  name        = "vm-machine-image-service-${var.machine-image-service-locations[count.index]}"
  description = "API service for fetching machine boot images"

  location = var.machine-image-service-locations[count.index]
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.machine_image_service.email

    containers {
      image = var.machine-image-service-image

      resources {
        limits = {
          cpu    = var.machine-image-service-cpu-limit
          memory = var.machine-image-service-memory-limit
        }
        cpu_idle = false
      }

      dynamic "env" {
        for_each = var.machine-image-service-env-vars
        content {
          name  = env.value["name"]
          value = env.value["value"]
        }
      }

      ports {
        container_port = one([for env_var in var.machine-image-service-env-vars : env_var.value if env_var.name == "HTTP_PORT"])
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
      max_instance_count = var.machine-image-service-max-instance-count
    }

    timeout                          = "${var.machine-image-service-max-request-timeout}s"
    max_instance_request_concurrency = var.machine-image-service-max-concurrent-requests
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}
