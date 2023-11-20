terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.3.0"
    }
  }
}

resource "google_privateca_ca_pool" "default" {
  name     = "global-private-ca"
  tier     = "ENTERPRISE"
  location = "us-east1"

  publishing_options {
    publish_ca_cert = true
    publish_crl     = true
  }

  issuance_policy {
    baseline_values {
      ca_options {
        is_ca                  = true
        max_issuer_path_length = 10
      }

      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign  = true
        }
        extended_key_usage {
          server_auth = true
          client_auth = true
        }
      }

      name_constraints {
        critical            = true
        permitted_dns_names = var.domains
      }
    }
  }
}

resource "google_privateca_certificate_authority" "root" {
  pool                                   = google_privateca_ca_pool.default.name
  certificate_authority_id               = "global-certificate-authority"
  location                               = "us-east1"
  deletion_protection                    = false
  skip_grace_period                      = true
  ignore_active_certificates_on_deletion = true

  config {
    subject_config {
      subject {
        organization = var.organization_name
        common_name  = "global-certificate-authority"
      }
      subject_alt_name {
        dns_names = var.domains
      }
    }

    x509_config {
      ca_options {
        is_ca                  = true
        max_issuer_path_length = 10
      }

      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign  = true
        }
        extended_key_usage {
          server_auth = true
          client_auth = true
        }
      }

      name_constraints {
        critical            = true
        permitted_dns_names = var.domains
      }
    }
  }
  lifetime = "31104000s"
  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }
}
