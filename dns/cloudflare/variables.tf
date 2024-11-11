variable "domain_name" {
  type = string
}

variable "records" {
  type = map(object({
    authenticated_origin_pulls_enabled = bool

    ipv4 = optional(list(object({
      address = string
    })))

    ipv6 = optional(list(object({
      address = string
    })))
  }))
}