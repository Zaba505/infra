variable "name" {
  type        = string
  description = "Name of the internal load balancer"
}

variable "network" {
  type        = string
  description = "VPC network self-link for the internal load balancer"
}

variable "subnetwork" {
  type        = string
  description = "VPC subnetwork self-link for the internal load balancer"
}

variable "region" {
  type        = string
  description = "Region for the internal forwarding rule"
}

variable "default_service" {
  type = object({
    name      = string
    locations = list(string)
  })
  description = "Default Cloud Run service configuration (name and deployment regions)"
}

variable "cloud_run" {
  type = map(object({
    hosts     = list(string)
    locations = list(string)
    paths     = list(string)
  }))
  description = "Map of Cloud Run services with host rules, deployment regions, and path matchers"
  default     = {}
}

variable "backend_timeout_seconds" {
  type        = number
  description = "Backend service timeout in seconds (for large file downloads like kernel/initrd)"
  default     = 300
}

variable "enable_https" {
  type        = bool
  description = "Enable HTTPS (requires ssl_certificates variable). If false, uses HTTP only."
  default     = false
}

variable "ssl_certificates" {
  type        = list(string)
  description = "List of SSL certificate self-links for HTTPS. Required if enable_https is true."
  default     = []
}

variable "health_check" {
  type = object({
    check_interval_sec  = number
    timeout_sec         = number
    healthy_threshold   = number
    unhealthy_threshold = number
    request_path        = string
  })
  description = "Health check configuration for backend services"
  default = {
    check_interval_sec  = 10
    timeout_sec         = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    request_path        = "/health/liveness"
  }
}
