terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.40.0"
    }
  }
}

resource "google_compute_global_address" "ipv6" {
  name         = var.name
  ip_version   = "IPV6"
  address_type = "EXTERNAL"
}