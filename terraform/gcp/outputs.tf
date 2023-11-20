output "global_ipv4_address" {
  value = module.gateway.global_ipv4_address
}

output "global_ipv6_address" {
  value = module.gateway.global_ipv6_address
}

output "client_cacert_pem_certificates" {
  value = module.client_cacert[*].pem_certificate
}

output "client_cacert_pem_private_keys" {
  value = module.client_cacert[*].private_key_pem
}