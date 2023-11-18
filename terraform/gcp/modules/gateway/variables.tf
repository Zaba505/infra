variable "domains" {
  type = list(string)
}

variable "default_service" {
  type = object({
    name      = string
    locations = list(string)
  })
}

variable "apis" {
  type = map(object({
    paths = list(string)

    cloud_run = optional(object({
      service_name = string
      locations    = list(string)
    }))
  }))
}