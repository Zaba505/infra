variable "hostnames" {
  type = list(string)
}

variable "private_key" {
  type = object({
    algorithm = string
    rsa_bits  = optional(number)
  })

  default = {
    algorithm = "RSA"
    rsa_bits  = 2048
  }
}

variable "days_valid_for" {
  type    = number
  default = 365
}
