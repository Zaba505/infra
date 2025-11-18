# Internal Application Load Balancer Module

This Terraform module provisions an internal Application Load Balancer for Cloud Run services on Google Cloud Platform. The load balancer is VPC-internal and uses the `INTERNAL_MANAGED` scheme, making it accessible only within the VPC network (e.g., through WireGuard VPN).

## Features

- **Internal Load Balancing**: Creates VPC-internal load balancer (not exposed to the internet)
- **Serverless NEGs**: Automatic integration with Cloud Run services across multiple regions
- **HTTP/HTTPS Support**: Configurable protocol (HTTP by default, optional HTTPS with SSL certificates)
- **URL-based Routing**: Host rules and path matchers for advanced traffic routing
- **Health Checks**: Configurable health checks for Cloud Run backends
- **Configurable Timeouts**: Backend timeout configuration for large file downloads (e.g., kernel/initrd images)
- **Zero-downtime Updates**: Uses `create_before_destroy` lifecycle for seamless updates

## Usage

### Basic Example (HTTP)

```hcl
module "boot_server_lb" {
  source = "./cloud/internal-application-load-balancer"

  name       = "boot-server-lb"
  network    = google_compute_network.vpc.self_link
  subnetwork = google_compute_subnetwork.subnet.self_link
  region     = "us-central1"

  default_service = {
    name      = "boot-server"
    locations = ["us-central1", "us-east1"]
  }

  backend_timeout_seconds = 600  # 10 minutes for large file downloads
}
```

### Advanced Example (HTTPS with Multiple Services)

```hcl
module "internal_lb" {
  source = "./cloud/internal-application-load-balancer"

  name       = "internal-app-lb"
  network    = google_compute_network.vpc.self_link
  subnetwork = google_compute_subnetwork.subnet.self_link
  region     = "us-central1"

  enable_https      = true
  ssl_certificates  = [google_compute_region_ssl_certificate.internal.self_link]

  default_service = {
    name      = "api-gateway"
    locations = ["us-central1"]
  }

  cloud_run = {
    user-service = {
      hosts     = ["api.internal"]
      locations = ["us-central1", "us-east1"]
      paths     = ["/users/*", "/auth/*"]
    }
    data-service = {
      hosts     = ["api.internal"]
      locations = ["us-central1"]
      paths     = ["/data/*"]
    }
  }

  health_check = {
    check_interval_sec  = 10
    timeout_sec         = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    request_path        = "/health/liveness"
  }
}
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| `name` | Name of the internal load balancer | `string` | - | yes |
| `network` | VPC network self-link for the internal load balancer | `string` | - | yes |
| `subnetwork` | VPC subnetwork self-link for the internal load balancer | `string` | - | yes |
| `region` | Region for the internal forwarding rule | `string` | - | yes |
| `default_service` | Default Cloud Run service configuration (name and deployment regions) | `object` | - | yes |
| `cloud_run` | Map of Cloud Run services with host rules, deployment regions, and path matchers | `map(object)` | `{}` | no |
| `backend_timeout_seconds` | Backend service timeout in seconds | `number` | `300` | no |
| `enable_https` | Enable HTTPS (requires ssl_certificates variable) | `bool` | `false` | no |
| `ssl_certificates` | List of SSL certificate self-links for HTTPS | `list(string)` | `[]` | no |
| `health_check` | Health check configuration for backend services | `object` | See below | no |

### Default Health Check Configuration

```hcl
{
  check_interval_sec  = 10
  timeout_sec         = 5
  healthy_threshold   = 2
  unhealthy_threshold = 2
  request_path        = "/health/liveness"
}
```

## Outputs

| Name | Description |
|------|-------------|
| `forwarding_rule_id` | The ID of the internal forwarding rule |
| `forwarding_rule_name` | The name of the internal forwarding rule |
| `internal_ip_address` | The internal IP address of the load balancer |
| `url_map_id` | The ID of the URL map |
| `backend_service_ids` | Map of Cloud Run service names to their backend service IDs |
| `health_check_id` | The ID of the health check |

## Architecture

This module creates the following resources:

1. **Regional Network Endpoint Groups (NEGs)**: Serverless NEGs for each Cloud Run service in each deployment region
2. **Health Check**: HTTP health check for backend services
3. **Backend Services**: One per Cloud Run service, configured with the NEGs
4. **URL Map**: Routes traffic based on host rules and path matchers
5. **Target Proxy**: Regional HTTP or HTTPS proxy (depending on `enable_https`)
6. **Forwarding Rule**: Internal forwarding rule with `INTERNAL_MANAGED` scheme

## Requirements

- GCP Provider version: `7.11.0`
- Cloud Run services must be deployed with `ingress = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"`
- Cloud Run services must expose health check endpoints (default: `/health/liveness`)
- VPC network and subnetwork must exist before applying this module

## Notes

- The load balancer is **internal only** - not accessible from the internet
- Backend timeout is configurable to support large file downloads (e.g., boot images)
- All resources use `create_before_destroy` lifecycle for zero-downtime updates
- Health checks default to port 80 and `/health/liveness` path
- HTTPS support is optional and requires providing SSL certificates separately

## Related Modules

- `cloud/rest-api`: Creates Cloud Run services compatible with this load balancer
- `cloud/https-load-balancer`: External load balancer with mTLS support (different use case)
