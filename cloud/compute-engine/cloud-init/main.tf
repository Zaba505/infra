terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.12.0"
    }
  }
}

data "google_client_config" "default" {}

# Create service account for instances
resource "google_service_account" "this" {
  account_id   = var.name
  display_name = var.name
}

# Grant IAM roles to the service account
resource "google_project_iam_member" "this" {
  for_each = toset(var.service_account_roles)

  project = data.google_client_config.default.project
  role    = each.value
  member  = "serviceAccount:${google_service_account.this.email}"
}

# Fetch cloud-init configuration from Secret Manager if provided
data "google_secret_manager_secret_version_access" "cloud_init" {
  count = var.cloud_init_secret != null ? 1 : 0

  secret  = var.cloud_init_secret.name
  version = var.cloud_init_secret.version
}

locals {
  # Use cloud-init from Secret Manager if provided, otherwise use the direct value
  cloud_init_config = var.cloud_init_secret != null ? data.google_secret_manager_secret_version_access.cloud_init[0].secret_data : var.cloud_init_config
}

# Instance template with cloud-init support
resource "google_compute_instance_template" "default" {
  name_prefix = "${var.name}-"
  description = var.description

  machine_type = var.machine_type

  # Boot disk configuration
  disk {
    source_image = var.boot_disk.image
    disk_size_gb = var.boot_disk.size_gb
    disk_type    = var.boot_disk.type
    auto_delete  = true
    boot         = true
  }

  # Network configuration
  network_interface {
    network    = var.network.vpc
    subnetwork = var.network.subnet

    # External IP configuration
    dynamic "access_config" {
      for_each = var.network.external_ip ? [1] : []
      content {
        network_tier = "PREMIUM"
      }
    }
  }

  # Service account configuration
  service_account {
    email  = google_service_account.this.email
    scopes = var.service_account_scopes
  }

  # Network tags for firewall rules
  tags = var.network_tags

  # Cloud-init user-data metadata
  metadata = merge(
    {
      "user-data" = local.cloud_init_config
    },
    var.additional_metadata
  )

  # Lifecycle management for zero-downtime updates
  lifecycle {
    create_before_destroy = true
  }

  # Labels for resource organization
  labels = var.labels
}

# Health check for auto-healing
resource "google_compute_health_check" "autohealing" {
  name                = "${var.name}-autohealing"
  check_interval_sec  = var.health_check.check_interval_sec
  timeout_sec         = var.health_check.timeout_sec
  healthy_threshold   = var.health_check.healthy_threshold
  unhealthy_threshold = var.health_check.unhealthy_threshold

  dynamic "http_health_check" {
    for_each = var.health_check.type == "http" ? [1] : []
    content {
      port         = var.health_check.port
      request_path = var.health_check.request_path
    }
  }

  dynamic "https_health_check" {
    for_each = var.health_check.type == "https" ? [1] : []
    content {
      port         = var.health_check.port
      request_path = var.health_check.request_path
    }
  }

  dynamic "tcp_health_check" {
    for_each = var.health_check.type == "tcp" ? [1] : []
    content {
      port = var.health_check.port
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Managed instance group for auto-healing and restart policies
resource "google_compute_instance_group_manager" "default" {
  for_each = toset(var.zones)

  name               = "${var.name}-${each.value}"
  base_instance_name = var.name
  zone               = each.value

  version {
    instance_template = google_compute_instance_template.default.id
  }

  target_size = var.instance_count

  # Auto-healing configuration
  auto_healing_policies {
    health_check      = google_compute_health_check.autohealing.id
    initial_delay_sec = var.health_check.initial_delay_sec
  }

  # Update policy for zero-downtime updates
  update_policy {
    type                           = "PROACTIVE"
    minimal_action                 = "REPLACE"
    max_surge_fixed                = 1
    max_unavailable_fixed          = 0
    replacement_method             = "SUBSTITUTE"
    most_disruptive_allowed_action = "REPLACE"
  }

  lifecycle {
    create_before_destroy = true
  }
}
