output "service_account_email" {
  value = google_service_account.api.email
}

output "service_id" {
  value = google_cloud_run_v2_service.api.id
}
