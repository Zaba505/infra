terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.11.0"
    }
  }
}

data "google_project" "default" {}

resource "google_firestore_database" "this" {
  project     = data.google_project.default.name
  name        = var.name
  location_id = var.location

  type                        = var.database_type
  concurrency_mode            = var.concurrency_mode
  app_engine_integration_mode = var.app_engine_integration_mode
  deletion_policy             = var.deletion_policy

  point_in_time_recovery_enablement = var.point_in_time_recovery_enabled ? "POINT_IN_TIME_RECOVERY_ENABLED" : "POINT_IN_TIME_RECOVERY_DISABLED"

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_project_iam_member" "datastore_user" {
  for_each = toset(var.service_account_users)

  project = data.google_project.default.name
  role    = "roles/datastore.user"
  member  = "serviceAccount:${each.value}"

  condition {
    title       = "datastore_${var.name}_access"
    description = "Access to Firestore database ${var.name}"
    expression  = "resource.name.startsWith('projects/${data.google_project.default.name}/databases/${var.name}')"
  }
}

resource "google_project_iam_member" "datastore_viewer" {
  for_each = toset(var.service_account_viewers)

  project = data.google_project.default.name
  role    = "roles/datastore.viewer"
  member  = "serviceAccount:${each.value}"

  condition {
    title       = "datastore_${var.name}_viewer"
    description = "Read-only access to Firestore database ${var.name}"
    expression  = "resource.name.startsWith('projects/${data.google_project.default.name}/databases/${var.name}')"
  }
}
