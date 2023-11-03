variable "container-images-registry-location" {
  type = string
}

variable "machine-image-service-account-id" {
  type    = string
  default = "machine-image-sa"

  validation {
    condition     = can(regex("^[a-z]([-a-z0-9]*[a-z0-9])$", var.machine-image-service-account-id))
    error_message = "The Service Account ID must comply with RFC1035"
  }
}

variable "machine-image-service-location" {
  type = string
}

variable "machine-image-service-image" {
  type = string
}

variable "machine-image-service-cpu-limit" {
  type    = number
  default = 1
}

variable "machine-image-service-memory-limit" {
  type    = string
  default = "512Mi"
}

variable "machine-image-service-max-instance-count" {
  type    = number
  default = null

  validation {
    condition     = var.machine-image-service-max-instance-count != null ? var.machine-image-service-max-instance-count > 0 : true
    error_message = "Max instance count must be greater than zero."
  }
}

variable "machine-image-service-max-request-timeout" {
  type        = number
  default     = 10
  description = "Max allowed time for an instance to respond to a request. A duration in seconds with up to nine fractional digits."

  validation {
    condition     = var.machine-image-service-max-request-timeout > 0
    error_message = "Max request timeout must be greater than zero."
  }
}

variable "machine-image-service-max-concurrent-requests" {
  type        = number
  default     = 10
  description = "Max concurrent requests per instance."

  validation {
    condition     = var.machine-image-service-max-concurrent-requests > 0
    error_message = "Max concurrent requests must be greater than zero."
  }
}

variable "boot-image-bucket-name" {
  type = string
}

variable "boot-image-bucket-location" {
  type = string
}