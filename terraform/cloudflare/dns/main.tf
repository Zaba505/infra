terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.0"
    }
  }
}

locals {
  a_records = { for name, record in var.records : name => record.ipv4 if record.ipv4 != null }

  aaaa_records = { for name, record in var.records : name => record.ipv6 if record.ipv6 != null }

  secured_records = { for name, record in var.records : name => {
    certificate = record.certificate,
    private_key = record.private_key
  } if record.certificate != null }
}

data "cloudflare_zone" "default" {
  name = var.domain_name
}

resource "cloudflare_record" "ipv4" {
  for_each = local.a_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.key
  value   = each.value.address
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ipv6" {
  for_each = local.aaaa_records

  zone_id = data.cloudflare_zone.default.id
  name    = each.key
  value   = each.value.address
  type    = "AAAA"
  proxied = true
}

resource "cloudflare_authenticated_origin_pulls_certificate" "per_hostname" {
  depends_on = [
    cloudflare_record.ipv4,
    cloudflare_record.ipv6
  ]
  for_each = local.secured_records

  zone_id = data.cloudflare_zone.default.id
  type    = "per-hostname"

  certificate = each.value.certificate
  private_key = each.value.private_key
}

resource "cloudflare_authenticated_origin_pulls" "per_hostname" {
  depends_on = [
    cloudflare_record.ipv4,
    cloudflare_record.ipv6
  ]
  for_each = local.secured_records

  zone_id  = data.cloudflare_zone.default.id
  enabled  = true
  hostname = "${each.key}.${var.domain_name}"
}