terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.13.0"
    }
  }
}

locals {
  a_records = merge([
    for name, record in var.records : {
      for ipv4 in record.ipv4 : "${name}-${ipv4.address}" => {
        name    = name
        address = ipv4.address
      }
    }
    if record.ipv4 != null
  ]...)

  aaaa_records = merge([
    for name, record in var.records : {
      for ipv6 in record.ipv6 : "${name}-${ipv6.address}" => {
        name    = name
        address = ipv6.address
      }
    }
    if record.ipv6 != null
  ]...)
}

data "cloudflare_zone" "default" {
  name = var.domain_name
}

resource "cloudflare_record" "a" {
  for_each = local.a_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.value.name
  value   = each.value.address
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "aaaa" {
  for_each = local.aaaa_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.value.name
  value   = each.value.address
  type    = "AAAA"
  proxied = true
}

locals {
  mtls_proxy_records = toset([for name, record in var.records : name if record.authenticated_origin_pulls_enabled])
}

resource "cloudflare_authenticated_origin_pulls" "per_hostname" {
  depends_on = [
    cloudflare_record.a,
    cloudflare_record.aaaa
  ]

  for_each = local.mtls_proxy_records

  zone_id  = data.cloudflare_zone.default.id
  enabled  = true
  hostname = "${each.value}.${var.domain_name}"
}