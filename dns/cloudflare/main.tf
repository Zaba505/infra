terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.45.0"
    }
  }
}

locals {
  a_records = {
    for name, record in var.records : "${name}-${record.ipv4.address}" => {
      name = name
      address = record.ipv4.address
    }
    if record.ipv4 != null
  }

  aaaa_records = {
    for name, record in var.records : "${name}-${record.ipv6.address}" => {
      name = name
      address = record.ipv6.address
    } if record.ipv6 != null
  }
}

data "cloudflare_zone" "default" {
  name = var.domain_name
}

resource "cloudflare_record" "ipv4" {
  for_each = local.a_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.value.name
  value   = each.value.address
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ipv6" {
  for_each = local.aaaa_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.value.name
  value   = each.value.address
  type    = "AAAA"
  proxied = true
}

locals {
  mtls_proxy_records = toset([for name, record in var.records: name if record.authenticated_origin_pulls_enabled])
}

resource "cloudflare_authenticated_origin_pulls" "per_hostname" {
  depends_on = [
    cloudflare_record.ipv4,
    cloudflare_record.ipv6
  ]

  for_each = local.mtls_proxy_records

  zone_id  = data.cloudflare_zone.default.id
  enabled = true
  hostname = "${each.value}.${var.domain_name}"
}