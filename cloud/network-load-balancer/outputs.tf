output "external_ip" {
  description = "External IP address of the network load balancer"
  value = (
    var.external_ip_address != null
    ? data.google_compute_address.external_ip[0].address
    : google_compute_address.external_ip[0].address
  )
}

output "external_ip_name" {
  description = "Name of the external IP address resource"
  value = (
    var.external_ip_address != null
    ? var.external_ip_address
    : google_compute_address.external_ip[0].name
  )
}

output "backend_service_id" {
  description = "ID of the backend service"
  value       = google_compute_region_backend_service.default.id
}

output "backend_service_name" {
  description = "Name of the backend service"
  value       = google_compute_region_backend_service.default.name
}

output "tcp_forwarding_rule_id" {
  description = "ID of the TCP forwarding rule (if created)"
  value       = try(google_compute_forwarding_rule.tcp[0].id, null)
}

output "udp_forwarding_rule_id" {
  description = "ID of the UDP forwarding rule (if created)"
  value       = try(google_compute_forwarding_rule.udp[0].id, null)
}

output "health_check_id" {
  description = "ID of the health check"
  value       = google_compute_health_check.default.id
}
