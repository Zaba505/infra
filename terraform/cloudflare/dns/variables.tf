variable "domain_name" {
  type = string
}

variable "records" {
  type = map(object({
    ipv4 = optional(object({
      address = string
    }))

    ipv6 = optional(object({
      address = string
    }))

    certificate = optional(object({
      pem         = string
      private_key = string
    }))
  }))
}