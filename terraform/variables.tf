variable "email_address" {
  type = string
}

variable "domain_zone" {
  type = string
}

variable "gcp_locations" {
  type = list(string)
}

variable "default_service" {
  type = object({
    image_tag                   = optional(string, "latest")
    cpu_limit                   = optional(number, 1)
    memory_limit                = optional(string, "512Mi")
    max_instance_count          = optional(number, 10)
    max_concurrent_requests     = optional(number, 100)
    max_request_timeout_seconds = optional(number, 1)
  })
}

variable "machine_mgmt_service" {
  type = object({
    image_tag               = optional(string, "latest")
    cpu_limit               = optional(number, 1)
    memory_limit            = optional(string, "512Mi")
    max_instance_count      = optional(number, 1)
    max_concurrent_requests = optional(number, 10)
  })
}