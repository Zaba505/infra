variable "name" {
  type = string
}

variable "ip_addresses" {
  type = list(object({
    name = string
  }))
}

variable "trust_anchor_secrets" {
  type = list(object({
    secret  = string
    version = string
  }))
}

variable "server_certificate_secrets" {
  type = map(object({
    certificate_secret  = string
    certificate_version = string
    private_key_secret  = string
    private_key_version = string
  }))
}

variable "default_service" {
  type = object({
    name      = string
    locations = list(string)
  })
}

variable "cloud_run" {
  type = map(object({
    hosts     = list(string)
    locations = list(string)
    paths     = list(string)
  }))
}