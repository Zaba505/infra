variable "hostnames" {
  type = list(string)
}

variable "days_valid_for" {
  type    = number
  default = 180
}