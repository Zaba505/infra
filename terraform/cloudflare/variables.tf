variable "domain_name" {
  type = string
}

variable "records" {
  type = map(object({
    ipv4 = string
    ipv6 = string
  }))
}