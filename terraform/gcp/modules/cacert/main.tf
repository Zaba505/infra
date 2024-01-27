terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.6.0"
    }

    random = {
      source  = "hashicorp/random"
      version = ">= 3.5.1"
    }
  }
}

resource "tls_private_key" "default" {
  algorithm = var.algorithm
}

resource "tls_cert_request" "default" {
  private_key_pem = tls_private_key.default.private_key_pem

  subject {
    common_name  = var.common_name
    organization = var.organization
  }
}

resource "random_id" "default" {
  byte_length = 8
}

resource "google_privateca_certificate" "default" {
  pool                  = var.privateca_pool_name
  location              = var.location
  certificate_authority = var.certificate_authority_id
  lifetime              = var.lifetime
  name                  = "${var.name}-${random_id.default.hex}"
  pem_csr               = tls_cert_request.default.cert_request_pem
}