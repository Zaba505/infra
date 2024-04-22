variable "project_id" {
  type = string
}

variable "cloudflare_api_token" {
  type = string
}

variable "ipv6_address" {
  type = string
}

variable "hosts" {
  type = list(string)
}

variable "ca_certificate_pems" {
  type = list(string)
}

variable "default_service" {
  type = object({
    name      = string
    locations = list(string)
  })
}

variable "cloud_run" {
  type = map(object({
    locations = list(string)
    paths     = list(string)
  }))
}