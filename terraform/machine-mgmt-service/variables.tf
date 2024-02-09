variable "gcp_locations" {
  type = list(string)
}

variable "destination_registries" {
  type = map(string)
}

variable "image_tag" {
  type = string
}

variable "cpu_limit" {
  type    = number
  default = 1
}

variable "memory_limit" {
  type    = string
  default = "512Mi"
}

variable "max_instance_count" {
  type    = number
  default = 1
}

variable "max_concurrent_requests" {
  type    = number
  default = 10
}

variable "max_request_timeout_seconds" {
  type    = number
  default = 5
}