terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.34.0"
    }
  }
}

resource "google_artifact_registry_repository" "docker" {
  format        = "DOCKER"
  repository_id = var.name
  description   = var.description
  location      = var.location
  mode          = "STANDARD_REPOSITORY"
}