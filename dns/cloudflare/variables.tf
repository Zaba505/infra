variable "cloudflare_api_token" {
  type = string
}

variable "domain_name" {
  type = string
}

variable "records" {
  type = map(object({
    enable_mtls = bool

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