variable "boot-image-storage-bucket-name" {
  type = string
}

variable "boot-image-service-accounts" {
  type = map(object({
    email = string
  }))
}