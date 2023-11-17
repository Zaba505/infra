variable "domains" {
  type = list(string)
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