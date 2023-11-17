output "service_account_email" {
  value = google_service_account.api.email
}

output "name" {
  value = var.name
}

output "locations" {
  value = var.locations
}