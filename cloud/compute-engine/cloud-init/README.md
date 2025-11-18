# Compute Engine with Cloud-Init Module

This Terraform module provisions Google Compute Engine VM instances with cloud-init support. It creates a managed instance group with auto-healing capabilities, making it ideal for deploying containerized services like WireGuard VPN gateway.

## Features

- **Cloud-init Support**: User-data configuration via cloud-init for automated instance setup
- **Secret Manager Integration**: Fetch cloud-init configuration securely from GCP Secret Manager
- **Managed Instance Group**: Auto-healing and restart policies for high availability
- **Health Checks**: Configurable HTTP/HTTPS/TCP health checks for instance monitoring
- **Zero-Downtime Updates**: Uses `create_before_destroy` lifecycle for seamless updates
- **Flexible Configuration**: Customizable machine type, boot disk, network settings
- **Service Account Binding**: IAM role integration for GCP service access

## Usage

### Basic Example

```hcl
module "wireguard_gateway" {
  source = "../../cloud/compute-engine/cloud-init"

  name        = "wireguard-gateway"
  description = "WireGuard VPN Gateway Instance"
  zone        = "us-central1-a"

  machine_type   = "e2-micro"
  instance_count = 1

  boot_disk = {
    image   = "cos-cloud/cos-stable"  # Container-Optimized OS
    size_gb = 10
    type    = "pd-standard"
  }

  network = {
    vpc         = "default"
    subnet      = "default"
    external_ip = true  # Required for VPN gateway
  }

  service_account_email = "wireguard-sa@project.iam.gserviceaccount.com"

  network_tags = ["wireguard", "vpn-gateway"]

  cloud_init_config = <<-EOT
    #cloud-config
    write_files:
      - path: /etc/systemd/system/wireguard.service
        permissions: '0644'
        content: |
          [Unit]
          Description=WireGuard VPN
          After=docker.service
          Requires=docker.service

          [Service]
          ExecStartPre=-/usr/bin/docker rm -f wireguard
          ExecStart=/usr/bin/docker run --name wireguard \
            --cap-add=NET_ADMIN \
            --cap-add=SYS_MODULE \
            --sysctl="net.ipv4.conf.all.src_valid_mark=1" \
            -e PUID=1000 \
            -e PGID=1000 \
            -e TZ=America/Chicago \
            -p 51820:51820/udp \
            -v /etc/wireguard:/config \
            --restart unless-stopped \
            linuxserver/wireguard
          ExecStop=/usr/bin/docker stop wireguard

          [Install]
          WantedBy=multi-user.target

    runcmd:
      - systemctl daemon-reload
      - systemctl enable wireguard.service
      - systemctl start wireguard.service
  EOT

  health_check = {
    type                = "tcp"
    port                = 51820  # WireGuard port
    request_path        = "/"
    check_interval_sec  = 30
    timeout_sec         = 10
    healthy_threshold   = 2
    unhealthy_threshold = 3
    initial_delay_sec   = 120  # Allow time for Docker to pull image and start
  }

  labels = {
    service     = "wireguard"
    environment = "production"
  }
}
```

### WireGuard with Secret Manager

```hcl
# Store WireGuard configuration in Secret Manager
resource "google_secret_manager_secret" "wireguard_cloud_init" {
  secret_id = "wireguard-cloud-init"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "wireguard_cloud_init" {
  secret = google_secret_manager_secret.wireguard_cloud_init.id

  secret_data = <<-EOT
    #cloud-config
    write_files:
      - path: /etc/wireguard/wg0.conf
        permissions: '0600'
        content: |
          [Interface]
          PrivateKey = ${var.wireguard_private_key}
          Address = 10.100.0.1/24
          ListenPort = 51820

          [Peer]
          PublicKey = ${var.wireguard_peer_public_key}
          AllowedIPs = 10.100.0.2/32

      - path: /etc/systemd/system/wireguard.service
        permissions: '0644'
        content: |
          [Unit]
          Description=WireGuard VPN
          After=docker.service
          Requires=docker.service

          [Service]
          ExecStartPre=-/usr/bin/docker rm -f wireguard
          ExecStart=/usr/bin/docker run --name wireguard \
            --cap-add=NET_ADMIN \
            --cap-add=SYS_MODULE \
            --sysctl="net.ipv4.conf.all.src_valid_mark=1" \
            -e PUID=1000 \
            -e PGID=1000 \
            -v /etc/wireguard:/config \
            -p 51820:51820/udp \
            --restart unless-stopped \
            linuxserver/wireguard
          ExecStop=/usr/bin/docker stop wireguard

          [Install]
          WantedBy=multi-user.target

    runcmd:
      - systemctl daemon-reload
      - systemctl enable wireguard.service
      - systemctl start wireguard.service
  EOT
}

module "wireguard_gateway" {
  source = "../../cloud/compute-engine/cloud-init"

  name        = "wireguard-gateway"
  description = "WireGuard VPN Gateway with cloud-init from Secret Manager"
  zone        = "us-central1-a"

  machine_type   = "e2-micro"
  instance_count = 1

  boot_disk = {
    image   = "cos-cloud/cos-stable"
    size_gb = 10
    type    = "pd-standard"
  }

  network = {
    vpc         = "default"
    subnet      = "default"
    external_ip = true
  }

  service_account_email = "wireguard-sa@project.iam.gserviceaccount.com"

  network_tags = ["wireguard", "vpn-gateway"]

  # Fetch cloud-init from Secret Manager
  cloud_init_secret = {
    name    = google_secret_manager_secret.wireguard_cloud_init.secret_id
    version = "latest"
  }

  health_check = {
    type                = "tcp"
    port                = 51820
    request_path        = "/"
    check_interval_sec  = 30
    timeout_sec         = 10
    healthy_threshold   = 2
    unhealthy_threshold = 3
    initial_delay_sec   = 120
  }
}

# Grant the service account permission to read the secret
resource "google_secret_manager_secret_iam_member" "wireguard_sa" {
  secret_id = google_secret_manager_secret.wireguard_cloud_init.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:wireguard-sa@project.iam.gserviceaccount.com"
}
```

### HTTP Health Check Example

For services with HTTP endpoints:

```hcl
module "web_service" {
  source = "../../cloud/compute-engine/cloud-init"

  name = "web-service"
  zone = "us-central1-a"

  machine_type = "e2-small"

  boot_disk = {
    image   = "ubuntu-os-cloud/ubuntu-2204-lts"
    size_gb = 20
    type    = "pd-standard"
  }

  network = {
    vpc         = "default"
    subnet      = "default"
    external_ip = false
  }

  service_account_email = "web-service-sa@project.iam.gserviceaccount.com"

  network_tags = ["web-service"]

  cloud_init_config = <<-EOT
    #cloud-config
    packages:
      - docker.io

    runcmd:
      - systemctl enable docker
      - systemctl start docker
      - docker run -d -p 8080:8080 --name web-app my-web-app:latest
  EOT

  health_check = {
    type                = "http"
    port                = 8080
    request_path        = "/health"
    check_interval_sec  = 10
    timeout_sec         = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
    initial_delay_sec   = 60
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| name | Name of the instance and managed instance group | `string` | n/a | yes |
| zone | GCP zone for the instance group | `string` | n/a | yes |
| network | Network configuration including VPC, subnet, and external IP | `object` | n/a | yes |
| service_account_email | Service account email for the instance | `string` | n/a | yes |
| description | Description of the instance template | `string` | `"Compute Engine instance with cloud-init support"` | no |
| machine_type | Machine type for the instance | `string` | `"e2-micro"` | no |
| instance_count | Number of instances in the managed instance group | `number` | `1` | no |
| boot_disk | Boot disk configuration | `object` | `{image = "ubuntu-os-cloud/ubuntu-2204-lts", size_gb = 10, type = "pd-standard"}` | no |
| service_account_scopes | Service account scopes | `list(string)` | `["https://www.googleapis.com/auth/cloud-platform"]` | no |
| network_tags | Network tags for firewall rules | `list(string)` | `[]` | no |
| cloud_init_config | Cloud-init user-data configuration (YAML) | `string` | `""` | no |
| cloud_init_secret | Secret Manager secret containing cloud-init config | `object` | `null` | no |
| additional_metadata | Additional metadata for the instance | `map(string)` | `{}` | no |
| labels | Labels to apply to resources | `map(string)` | `{}` | no |
| health_check | Health check configuration for auto-healing | `object` | See variables.tf | no |

## Outputs

| Name | Description |
|------|-------------|
| instance_group_manager_id | ID of the managed instance group |
| instance_group_manager_name | Name of the managed instance group |
| instance_group_manager_self_link | Self-link of the managed instance group |
| instance_template_id | ID of the instance template |
| instance_template_name | Name of the instance template |
| instance_template_self_link | Self-link of the instance template |
| health_check_id | ID of the health check |
| health_check_self_link | Self-link of the health check |

## Architecture Patterns

### WireGuard VPN Gateway Deployment

This module implements the architecture described in [ADR-0005: Network Boot Infrastructure Implementation on Google Cloud](../../../docs/content/r&d/adrs/0005-network-boot-infrastructure-gcp.md).

The WireGuard gateway serves as the secure tunnel between the home lab and GCP:

```
Home Lab Servers → UDM Pro → WireGuard VPN → GCP Services
```

Key characteristics:
- **Container-Optimized OS**: Uses Google's cos-cloud image for minimal attack surface
- **Cloud-init**: Automates WireGuard container deployment on boot
- **Auto-healing**: Managed instance group restarts failed instances
- **UDP Health Checks**: Cannot directly check UDP port, so uses TCP or SSH
- **Secret Manager**: Stores WireGuard private keys and peer configurations

### Zero-Downtime Updates

The module uses `lifecycle { create_before_destroy = true }` on both the instance template and managed instance group to enable zero-downtime updates:

1. New instance template created with updated configuration
2. Managed instance group updated to reference new template
3. New instances created before old instances are destroyed
4. Health checks ensure new instances are healthy before traffic switches

## Requirements

- Terraform >= 1.0
- GCP provider 7.11.0
- GCP project with Compute Engine API enabled
- Service account with appropriate IAM roles:
  - `roles/compute.instanceAdmin.v1` for instance management
  - `roles/iam.serviceAccountUser` for service account attachment
  - `roles/secretmanager.secretAccessor` if using Secret Manager

## Related Documentation

- [ADR-0005: Network Boot Infrastructure Implementation](../../../docs/content/r&d/adrs/0005-network-boot-infrastructure-gcp.md)
- [GCP Compute Engine Instance Templates](https://cloud.google.com/compute/docs/instance-templates)
- [Cloud-init Documentation](https://cloudinit.readthedocs.io/)
- [WireGuard VPN](https://www.wireguard.com/)
