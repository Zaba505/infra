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

output "health_check_ids" {
  description = "Map of health check IDs by instance group index"
  value       = { for idx, hc in google_compute_region_health_check.instance_group : idx => hc.id }
}

