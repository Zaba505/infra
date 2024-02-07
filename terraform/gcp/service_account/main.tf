terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }
  }
}

locals {
  cloud_storage_resource_name_condition = join(" || ", [
    for bucket in var.cloud_storage.buckets : "resource.name.startsWith('projects/_/buckets/${bucket}')"
  ])
}

resource "google_service_account" "this" {
  account_id   = var.name
  display_name = var.name
}

data "google_client_config" "default" {}

resource "google_project_iam_member" "cloud_trace" {
  count = var.cloud_trace ? 1 : 0

  project = data.google_client_config.default.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.this.email}"
}

resource "google_project_iam_member" "cloud_storage" {
  count = length(var.cloud_storage.buckets) > 0 ? 1 : 0

  project = data.google_client_config.default.project
  role    = "roles/run.serviceAgent"
  member  = "serviceAccount:${google_service_account.this.email}"

  condition {
    title      = "only_cloud_storage"
    expression = <<-EOL
      resource.service == 'storage.googleapis.com'
      && (resource.type == 'storage.googleapis.com/Bucket' || resource.type == 'storage.googleapis.com/Object)
      && (${local.cloud_storage_resource_name_condition})
    EOL
  }
}

resource "google_storage_bucket_access_control" "this" {
  for_each = toset(var.cloud_storage.buckets)

  bucket = each.value
  role   = "READER"
  entity = "user-${google_service_account.this.email}"
}

resource "google_storage_default_object_access_control" "this" {
  for_each = toset(var.cloud_storage.buckets)

  bucket = each.value
  role   = "READER"
  entity = "user-${google_service_account.this.email}"
}