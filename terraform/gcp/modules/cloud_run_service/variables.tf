variable "artifact_registry_id" {
  type = string
}

variable "access" {
  type = object({
    cloud_storage = optional(object({
      bucket_name = string
    }))
  })
  default = {}
}

variable "unauthenticated" {
  type    = bool
  default = false
}

variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = ""
}

variable "locations" {
  type = list(string)
}

variable "image" {
  type = object({
    name = string
    tag  = optional(string, "latest")
  })
}

variable "env_vars" {
  type = list(object({
    name  = string
    value = string
  }))
}

variable "cpu_limit" {
  type    = number
  default = 1
}

variable "memory_limit" {
  type    = string
  default = "512Mi"
}

variable "max_instance_count" {
  type    = number
  default = null

  validation {
    condition     = var.max_instance_count != null ? var.max_instance_count > 0 : true
    error_message = "Max instance count must be greater than zero."
  }
}

variable "max_request_timeout" {
  type        = number
  default     = 10
  description = "Max allowed time for an instance to respond to a request. A duration in seconds with up to nine fractional digits."

  validation {
    condition     = var.max_request_timeout > 0
    error_message = "Max request timeout must be greater than zero."
  }
}

variable "max_concurrent_requests" {
  type        = number
  default     = 10
  description = "Max concurrent requests per instance."

  validation {
    condition     = var.max_concurrent_requests > 0
    error_message = "Max concurrent requests must be greater than zero."
  }
}