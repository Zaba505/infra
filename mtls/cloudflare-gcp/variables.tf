variable "hostnames" {
  type = list(string)
}

variable "days_valid_for" {
  type    = number
  default = 365
}

variable "service_account_email" {
  type = string
}