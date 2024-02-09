output "private_key_pem" {
  value = acme_certificate.certificate.private_key_pem
}

output "certificate_pem" {
  value = acme_certificate.certificate.certificate_pem
}

output "issuer_pem" {
  value = acme_certificate.certificate.issuer_pem
}

output "full_certificate_pem" {
  value = "${acme_certificate.certificate.certificate_pem}${acme_certificate.certificate.issuer_pem}"
}

output "certificate_not_after" {
  value = acme_certificate.certificate.certificate_not_after
}