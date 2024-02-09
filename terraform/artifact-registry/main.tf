terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }
  }
}

data "google_client_config" "default" {}

resource "google_artifact_registry_repository" "docker" {
  for_each = toset(var.gcp_locations)

  format        = "DOCKER"
  repository_id = "docker-infra"
  description   = "internal infra container images"
  location      = each.value
  mode          = "STANDARD_REPOSITORY"
}