variable "name" {
  type        = string
  description = "The name of the Firestore database. Must be unique within the project. Use '(default)' for the default database."
}

variable "location" {
  type        = string
  description = "The location of the database. Available locations: https://cloud.google.com/firestore/docs/locations"
}

variable "database_type" {
  type        = string
  description = "The type of the database. Possible values: FIRESTORE_NATIVE, DATASTORE_MODE"
  default     = "FIRESTORE_NATIVE"

  validation {
    condition     = contains(["FIRESTORE_NATIVE", "DATASTORE_MODE"], var.database_type)
    error_message = "database_type must be either FIRESTORE_NATIVE or DATASTORE_MODE"
  }
}

variable "concurrency_mode" {
  type        = string
  description = "The concurrency control mode to use for this database. Possible values: OPTIMISTIC, PESSIMISTIC"
  default     = "OPTIMISTIC"

  validation {
    condition     = contains(["OPTIMISTIC", "PESSIMISTIC"], var.concurrency_mode)
    error_message = "concurrency_mode must be either OPTIMISTIC or PESSIMISTIC"
  }
}

variable "app_engine_integration_mode" {
  type        = string
  description = "The App Engine integration mode to use for this database. Possible values: ENABLED, DISABLED"
  default     = "DISABLED"

  validation {
    condition     = contains(["ENABLED", "DISABLED"], var.app_engine_integration_mode)
    error_message = "app_engine_integration_mode must be either ENABLED or DISABLED"
  }
}

variable "deletion_policy" {
  type        = string
  description = "Deletion behavior for this database. Possible values: ABANDON (database not deleted), DELETE (database is deleted)"
  default     = "DELETE"

  validation {
    condition     = contains(["ABANDON", "DELETE"], var.deletion_policy)
    error_message = "deletion_policy must be either ABANDON or DELETE"
  }
}

variable "point_in_time_recovery_enabled" {
  type        = bool
  description = "Whether to enable point-in-time recovery for this database"
  default     = false
}

variable "service_account_users" {
  type        = list(string)
  description = "List of service account emails to grant roles/datastore.user (read/write) access"
  default     = []
}

variable "service_account_viewers" {
  type        = list(string)
  description = "List of service account emails to grant roles/datastore.viewer (read-only) access"
  default     = []
}
