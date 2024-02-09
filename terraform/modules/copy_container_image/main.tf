terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

locals {
  source-image-parts      = split(":", var.source-image)
  source-image-name-parts = split("/", local.source-image-parts[0])
  source-image-name       = element(local.source-image-name-parts, length(local.source-image-name-parts) - 1)
  source-image-tag        = length(local.source-image-parts) != 2 ? "latest" : local.source-image-parts[1]

  destination-image-name = "${var.destination-registry}/${local.source-image-name}"
}

resource "docker_image" "source" {
  name = var.source-image
}

resource "docker_tag" "destination" {
  source_image = docker_image.source.name

  target_image = "${local.destination-image-name}:${local.source-image-tag}"
}

resource "docker_registry_image" "destination" {
  depends_on = [docker_tag.destination]

  name = "${local.destination-image-name}:${local.source-image-tag}"
}