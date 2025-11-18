variable "name" {
  description = "Name of the network load balancer and associated resources"
  type        = string
}

variable "region" {
  description = "GCP region for the regional network load balancer"
  type        = string
}

variable "protocols" {
  description = "List of protocols to support (TCP, UDP, or both). WireGuard uses UDP only."
  type        = list(string)
  default     = ["UDP"]

  validation {
    condition     = length(var.protocols) > 0
    error_message = "At least one protocol (TCP or UDP) must be specified."
  }

  validation {
    condition     = alltrue([for p in var.protocols : contains(["TCP", "UDP"], p)])
    error_message = "Protocols must be either 'TCP' or 'UDP'."
  }
}

variable "port_range" {
  description = "Port range for the forwarding rules. For a single port, set start and end to the same value."
  type = object({
    start = number
    end   = number
  })

  validation {
    condition     = var.port_range.start <= var.port_range.end
    error_message = "Port range start must be less than or equal to end."
  }

  validation {
    condition     = var.port_range.start >= 1 && var.port_range.end <= 65535
    error_message = "Port range must be between 1 and 65535."
  }
}

variable "external_ip_address" {
  description = "Name of an existing external IP address to use. Required to prevent IP changes that break DNS records."
  type        = string
}

variable "network_tier" {
  description = "Network tier for the external IP and forwarding rules (PREMIUM or STANDARD)"
  type        = string
  default     = "PREMIUM"

  validation {
    condition     = contains(["PREMIUM", "STANDARD"], var.network_tier)
    error_message = "Network tier must be either 'PREMIUM' or 'STANDARD'."
  }
}

variable "backend_protocol" {
  description = "Protocol used for communication between the load balancer and backends (TCP or UDP)"
  type        = string
  default     = "TCP"

  validation {
    condition     = contains(["TCP", "UDP"], var.backend_protocol)
    error_message = "Backend protocol must be either 'TCP' or 'UDP'."
  }
}

variable "backend_timeout_sec" {
  description = "Backend service timeout in seconds"
  type        = number
  default     = 30

  validation {
    condition     = var.backend_timeout_sec >= 1 && var.backend_timeout_sec <= 2147483647
    error_message = "Backend timeout must be between 1 and 2147483647 seconds."
  }
}

variable "instance_groups" {
  description = "List of instance group backends for the load balancer with their health check configurations"
  type = list(object({
    instance_group               = string
    balancing_mode               = string
    capacity_scaler              = optional(number)
    max_connections              = optional(number)
    max_connections_per_instance = optional(number)
    max_rate                     = optional(number)
    max_rate_per_instance        = optional(number)
    max_utilization              = optional(number)
    health_check = object({
      protocol            = string
      port                = number
      request_path        = optional(string)
      check_interval_sec  = optional(number, 10)
      timeout_sec         = optional(number, 5)
      healthy_threshold   = optional(number, 2)
      unhealthy_threshold = optional(number, 3)
    })
  }))
  default = []

  validation {
    condition = alltrue([
      for ig in var.instance_groups :
      contains(["TCP", "HTTP", "HTTPS"], ig.health_check.protocol)
    ])
    error_message = "Health check protocol must be 'TCP', 'HTTP', or 'HTTPS' for all instance groups."
  }

  validation {
    condition = alltrue([
      for ig in var.instance_groups :
      ig.health_check.protocol == "TCP" ||
      (ig.health_check.request_path != null && ig.health_check.request_path != "")
    ])
    error_message = "Health check request_path is required for HTTP and HTTPS protocols."
  }
}

