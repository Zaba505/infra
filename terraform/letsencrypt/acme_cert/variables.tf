variable "account_private_key_pem" {
  type      = string
  sensitive = true
}

variable "common_name" {
  type = string
}

variable "subject_alternative_names" {
  type    = list(string)
  default = []
}

variable "dns_challenge" {
  type = object({
    cloudflare = optional(object({}))
  })
}

