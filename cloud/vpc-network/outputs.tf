output "vpc_id" {
  description = "ID of the VPC network"
  value       = google_compute_network.vpc.id
}

output "vpc_name" {
  description = "Name of the VPC network"
  value       = google_compute_network.vpc.name
}

output "vpc_self_link" {
  description = "Self-link of the VPC network"
  value       = google_compute_network.vpc.self_link
}

output "subnet_ids" {
  description = "Map of subnet names to their IDs"
  value       = { for k, v in google_compute_subnetwork.subnets : k => v.id }
}

output "subnet_names" {
  description = "Map of subnet names to their full names"
  value       = { for k, v in google_compute_subnetwork.subnets : k => v.name }
}

output "subnet_self_links" {
  description = "Map of subnet names to their self-links"
  value       = { for k, v in google_compute_subnetwork.subnets : k => v.self_link }
}

output "subnet_cidr_ranges" {
  description = "Map of subnet names to their CIDR ranges"
  value       = { for k, v in google_compute_subnetwork.subnets : k => v.ip_cidr_range }
}

output "subnet_regions" {
  description = "Map of subnet names to their regions"
  value       = { for k, v in google_compute_subnetwork.subnets : k => v.region }
}

output "cloud_router_ids" {
  description = "Map of Cloud Router configuration names to their IDs"
  value       = { for k, v in google_compute_router.router : k => v.id }
}

output "cloud_nat_ids" {
  description = "Map of Cloud NAT configuration names to their IDs"
  value       = { for k, v in google_compute_router_nat.nat : k => v.id }
}

output "firewall_rule_ids" {
  description = "Map of firewall rule names to their IDs (combined ingress and egress)"
  value = merge(
    { for k, v in google_compute_firewall.ingress : k => v.id },
    { for k, v in google_compute_firewall.egress : k => v.id }
  )
}
