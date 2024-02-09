output "destination_registries" {
  value = {
    for loc in var.gcp_locations : loc => "${loc}-docker.pkg.dev/${data.google_client_config.default.project}/${google_artifact_registry_repository.docker[loc].name}"
  }
}