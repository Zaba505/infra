output "pem_certificate" {
  value = google_privateca_certificate.default.pem_certificate
}

output "private_key_pem" {
  value = tls_private_key.default.private_key_pem
}