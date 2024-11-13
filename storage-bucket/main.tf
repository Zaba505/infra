terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.11.1"
    }
  }
}

resource "google_storage_bucket" "this" {
  name     = var.name
  location = var.location

  force_destroy            = var.force_destroy
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