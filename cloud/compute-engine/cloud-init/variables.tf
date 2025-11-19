variable "name" {
  type        = string
  description = "Name of the instance and managed instance group"
}

variable "description" {
  type        = string
  description = "Description of the instance template"
  default     = "Compute Engine instance with cloud-init support"
}

variable "zones" {
  type        = list(string)
  description = "GCP zones for the instance groups (e.g., [\"us-central1-a\", \"us-central1-b\"])"
}

variable "machine_type" {
  type        = string
  description = "Machine type for the instance (e.g., e2-micro, e2-small, e2-medium)"
  default     = "e2-micro"
}

variable "instance_count" {
  type        = number
  description = "Number of instances to create in the managed instance group"
  default     = 1
}

variable "boot_disk" {
  type = object({
    image   = string
    size_gb = number
    type    = string
  })
  description = "Boot disk configuration including image, size, and type"
  default = {
    image   = "ubuntu-os-cloud/ubuntu-2204-lts"
    size_gb = 10
    type    = "pd-standard"
  }
}

variable "network" {
  type = object({
    vpc         = string
    subnet      = string
    external_ip = bool
  })
  description = "Network configuration including VPC, subnet, and external IP settings"
}

variable "service_account_scopes" {
  type        = list(string)
  description = "Service account scopes for the instance"
  default = [
    "https://www.googleapis.com/auth/cloud-platform",
  ]
}

variable "service_account_roles" {
  type        = list(string)
  description = "IAM roles to grant to the service account"
  default     = []
}

variable "network_tags" {
  type        = list(string)
  description = "Network tags for firewall rules"
  default     = []
}

variable "cloud_init_config" {
  type        = string
  description = "Cloud-init user-data configuration (YAML format). Used directly if cloud_init_secret is not provided."
  default     = ""
}

variable "cloud_init_secret" {
  type = object({
    name    = string
    version = string
  })
  description = "GCP Secret Manager secret containing cloud-init configuration. If provided, this takes precedence over cloud_init_config."
  default     = null
}

variable "additional_metadata" {
  type        = map(string)
  description = "Additional metadata for the instance"
  default     = {}
}

variable "labels" {
  type        = map(string)
  description = "Labels to apply to the instance template and instances"
  default     = {}
}

variable "health_check" {
  type = object({
    type                = string
    port                = number
    request_path        = string
    check_interval_sec  = number
    timeout_sec         = number
    healthy_threshold   = number
    unhealthy_threshold = number
    initial_delay_sec   = number
  })
  description = "Health check configuration for auto-healing. Type can be 'http', 'https', or 'tcp'."
  default = {
    type                = "tcp"
    port                = 22
    request_path        = "/"
    check_interval_sec  = 30
    timeout_sec         = 10
    healthy_threshold   = 2
    unhealthy_threshold = 3
    initial_delay_sec   = 300
  }
}
