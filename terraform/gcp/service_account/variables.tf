variable "name" {
  type = string
}

variable "cloud_trace" {
  type    = bool
  default = false
}

variable "cloud_storage" {
  type = object({
    buckets = list(string)
  })
  default = {
    buckets = []
  }
}