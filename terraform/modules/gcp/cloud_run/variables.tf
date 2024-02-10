variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = ""
}

variable "service_account_email" {
  type = string
}

variable "unsecured" {
  type    = bool
  default = false
}

variable "location" {
  type = string
}

variable "image" {
  type = object({
    name = string
    tag  = string
  })
}

variable "env_vars" {
  type = list(object({
    name  = string
    value = string
  }))
}

variable "cpu_limit" {
  type    = number
  default = 1
}

variable "memory_limit" {
  type    = string
  default = "512Mi"
}

variable "min_instance_count" {
  type    = number
  default = 0
}

variable "max_instance_count" {
  type = number
}

variable "max_request_timeout_seconds" {
  type    = number
  default = 5
}

variable "max_concurrent_requests" {
  type = number
}