terraform {
  required_providers {
    acme = {
      source  = "vancluever/acme"
      version = ">= 2.19"
    }
  }
}

provider "acme" {
  server_url = "https://acme-v02.api.letsencrypt.org/directory"
}

locals {
  dns_challenge_provider = var.dns_challenge.cloudflare != null ? "cloudflare" : ""
}

resource "acme_certificate" "certificate" {
  account_key_pem           = var.account_private_key_pem
  common_name               = var.common_name
  subject_alternative_names = var.subject_alternative_names

  dns_challenge {
    provider = local.dns_challenge_provider
  }
}