output "service_account_email" {
  value       = google_service_account.this.email
  description = "Email of the created service account"
}

output "service_account_id" {
  value       = google_service_account.this.id
  description = "ID of the created service account"
}

output "instance_group_managers" {
  value = {
    for zone, igm in google_compute_instance_group_manager.default : zone => {
      id        = igm.id
      name      = igm.name
      self_link = igm.self_link
    }
  }
  description = "Map of zone to instance group manager details"
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
