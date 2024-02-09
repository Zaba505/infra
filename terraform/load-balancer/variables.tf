variable "domain_zone" {
  type = string
}

variable "ca_certificate_pems" {
  type = list(string)
}

variable "destination_registries" {
  type = map(string)
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

variable "cloud_run" {
  type = map(object({
    locations = list(string)
    paths     = list(string)
  }))
}