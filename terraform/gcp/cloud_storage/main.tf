terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }
  }
}

resource "google_storage_bucket" "boot_images" {
  name     = var.bucket-name
  location = var.bucket-location

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