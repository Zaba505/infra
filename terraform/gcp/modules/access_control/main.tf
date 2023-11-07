terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

resource "google_storage_bucket_access_control" "boot_images" {
  for_each = toset(var.boot-image-service-account-emails)

  bucket = var.boot-image-storage-bucket-name
  role   = "READER"
  entity = "user-${each.key}"
}

resource "google_storage_default_object_access_control" "boot_images" {
  for_each = toset(var.boot-image-service-account-emails)

  bucket = var.boot-image-storage-bucket-name
  role   = "READER"
  entity = "user-${each.key}"
}