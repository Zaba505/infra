output "global_ipv4_address" {
  value = module.gateway.global_ipv4_address
}

output "global_ipv6_address" {
  value = module.gateway.global_ipv6_address
}

output "client_cacert_pem_certificates" {
  value = { for key, mod in module.client_cacert : key => mod.pem_certificate }
}

output "client_cacert_pem_private_keys" {
  value = { for key, mod in module.client_cacert : key => mod.private_key_pem }
}