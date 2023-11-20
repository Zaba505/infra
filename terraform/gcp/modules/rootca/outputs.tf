output "ca_pool_name" {
  value = google_privateca_ca_pool.default.name
}

output "certificate_authority_id" {
  value = google_privateca_certificate_authority.root.certificate_authority_id
}