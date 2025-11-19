output "forwarding_rule_id" {
  description = "The ID of the internal forwarding rule"
  value       = google_compute_forwarding_rule.default.id
}

output "forwarding_rule_name" {
  description = "The name of the internal forwarding rule"
  value       = google_compute_forwarding_rule.default.name
}

output "internal_ip_address" {
  description = "The internal IP address of the load balancer"
  value       = google_compute_forwarding_rule.default.ip_address
}

output "url_map_id" {
  description = "The ID of the URL map"
  value       = google_compute_region_url_map.default.id
}

output "backend_service_ids" {
  description = "Map of Cloud Run service names to their backend service IDs"
  value = merge(
    { (var.default_service.name) = google_compute_backend_service.default_service.id },
    { for name, svc in google_compute_backend_service.cloud_run : name => svc.id }
  )
}

output "health_check_id" {
  description = "The ID of the health check"
  value       = google_compute_health_check.default.id
}
