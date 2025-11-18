output "database_name" {
  description = "The name of the Firestore database"
  value       = google_firestore_database.this.name
}

output "database_id" {
  description = "The ID of the Firestore database"
  value       = google_firestore_database.this.id
}

output "database_location" {
  description = "The location of the Firestore database"
  value       = google_firestore_database.this.location_id
}

output "database_type" {
  description = "The type of the Firestore database (FIRESTORE_NATIVE or DATASTORE_MODE)"
  value       = google_firestore_database.this.type
}

output "database_uid" {
  description = "The unique identifier of the Firestore database"
  value       = google_firestore_database.this.uid
}
