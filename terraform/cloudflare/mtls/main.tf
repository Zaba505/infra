terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.5"
    }
  }
}

resource "tls_private_key" "instance" {
  algorithm = "RSA"
}

resource "tls_cert_request" "instance" {
  private_key_pem = tls_private_key.instance.private_key_pem

  subject {
    common_name = var.hostname
  }
}

resource "cloudflare_origin_ca_certificate" "instance" {
  csr                  = tls_cert_request.instance.cert_request_pem
  hostnames            = [var.hostname]
  request_type         = "origin-rsa"
  requested_validity   = 365
  min_days_for_renewal = 30
}