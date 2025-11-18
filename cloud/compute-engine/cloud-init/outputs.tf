output "instance_group_manager_id" {
  value       = google_compute_instance_group_manager.default.id
  description = "ID of the managed instance group"
}

output "instance_group_manager_name" {
  value       = google_compute_instance_group_manager.default.name
  description = "Name of the managed instance group"
}

output "instance_group_manager_self_link" {
  value       = google_compute_instance_group_manager.default.self_link
  description = "Self-link of the managed instance group"
}

output "instance_template_id" {
  value       = google_compute_instance_template.default.id
  description = "ID of the instance template"
}

output "instance_template_name" {
  value       = google_compute_instance_template.default.name
  description = "Name of the instance template"
}

output "instance_template_self_link" {
  value       = google_compute_instance_template.default.self_link
  description = "Self-link of the instance template"
}

output "health_check_id" {
  value       = google_compute_health_check.autohealing.id
  description = "ID of the health check"
}

output "health_check_self_link" {
  value       = google_compute_health_check.autohealing.self_link
  description = "Self-link of the health check"
}
