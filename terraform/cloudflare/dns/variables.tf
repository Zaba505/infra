variable "domain_name" {
  type = string
}

variable "records" {
  type = map(object({
    ipv4 = optional(object({
      address     = string
      certificate = optional(string)
      private_key = optional(string)
    }))

    ipv6 = optional(object({
      address     = string
      certificate = optional(string)
      private_key = optional(string)
    }))
  }))
}