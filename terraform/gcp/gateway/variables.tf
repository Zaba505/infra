variable "domain" {
  type = string
}

variable "cloud_run" {
  type = map(object({
    locations = list(string)
    paths     = list(string)
  }))
}

variable "default_service" {
  type = list(object({
    image = object({
      name = string
      tag  = string
    })
    location                    = string
    cpu_limit                   = optional(number, 1)
    memory_limit                = optional(string, "512Mi")
    max_instance_count          = optional(number, 10)
    max_concurrent_requests     = optional(number, 100)
    max_request_timeout_seconds = optional(number, 1)
  }))
}