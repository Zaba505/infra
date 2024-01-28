terraform {
  required_providers {
    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.5"
    }

    acme = {
      source  = "vancluever/acme"
      version = ">= 2.19"
    }
  }
}

resource "tls_private_key" "reg_private_key" {
  algorithm = "ED25519"
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.reg_private_key.private_key_pem
  email_address   = var.email_address
}

resource "acme_certificate" "certificate" {
  account_key_pem           = acme_registration.reg.account_key_pem
  common_name               = "www.${var.domain_name}"
  subject_alternative_names = var.subject_alternative_names

  dns_challenge {
    provider = "cloudflare"
  }
}