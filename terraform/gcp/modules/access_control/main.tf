terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

resource "google_storage_bucket_access_control" "boot_images" {
  for_each = var.boot-image-service-accounts

  bucket = var.boot-image-storage-bucket-name
  role   = "READER"
  entity = "user-${each.value.email}"
}

resource "google_storage_default_object_access_control" "boot_images" {
  for_each = var.boot-image-service-accounts

  bucket = var.boot-image-storage-bucket-name
  role   = "READER"
  entity = "user-${each.value.email}"
}