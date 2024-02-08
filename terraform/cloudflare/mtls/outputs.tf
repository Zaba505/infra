output "ca_certificate_pem" {
  value = cloudflare_origin_ca_certificate.instance.certificate
}

output "ca_private_key" {
  value = tls_private_key.instance.private_key_pem
}