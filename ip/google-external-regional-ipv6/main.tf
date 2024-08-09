terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.25.0"
    }
  }
}

resource "google_compute_address" "ipv6" {
  name         = var.name
  ip_version   = "IPV6"
  address_type = "EXTERNAL"

  # As of now, PREMIUM must be set since regional
  # external IPv6 addresses are not supported for
  # the STANDARD tier.
  network_tier = "PREMIUM"
}