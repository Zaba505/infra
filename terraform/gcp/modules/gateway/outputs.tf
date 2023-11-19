output "global_ipv4_address" {
  value = google_compute_global_address.ipv4.address
}

output "global_ipv6_address" {
  value = google_compute_global_address.ipv6.address
}