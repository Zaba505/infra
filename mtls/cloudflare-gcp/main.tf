terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.45.0"
    }

    google = {
      source  = "hashicorp/google"
      version = "6.10.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "4.0.6"
    }

    http = {
      source  = "hashicorp/http"
      version = "3.4.5"
    }
  }
}

resource "tls_private_key" "origin" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "tls_cert_request" "origin" {
  private_key_pem = tls_private_key.origin.private_key_pem

  dynamic "subject" {
    for_each = var.hostnames

    content {
      common_name = each.value
    }
  }
}

resource "cloudflare_origin_ca_certificate" "origin" {
  csr                = tls_cert_request.origin.cert_request_pem
  hostnames          = var.hostnames
  request_type       = "origin-rsa"
  requested_validity = var.days_valid_for
}

resource "google_secret_manager_secret" "origin_private_key" {
  secret_id = "cloudflare-origin-private-key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "origin_private_key" {
  secret = google_secret_manager_secret.origin_private_key.id

  secret_data = tls_private_key.origin.private_key_pem
}

resource "google_secret_manager_secret" "origin_certificate" {
  secret_id = "cloudflare-origin-certificate"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "origin_certificate" {
  secret = google_secret_manager_secret.origin_certificate.id

  secret_data = cloudflare_origin_ca_certificate.origin.certificate
}

resource "google_secret_manager_secret" "authenticated_origin_pull_ca" {
  secret_id = "cloudflare-authenticated-origin-pull-ca"

  replication {
    auto {}
  }
}

data "http" "authenticated_origin_pulls_ca_trust_anchor" {
  url = "https://developers.cloudflare.com/ssl/static/authenticated_origin_pull_ca.pem"
}

resource "google_secret_manager_secret_version" "authenticated_origin_pull_ca" {
  secret = google_secret_manager_secret.authenticated_origin_pull_ca.id

  secret_data = data.http.authenticated_origin_pulls_ca_trust_anchor.response_body
}