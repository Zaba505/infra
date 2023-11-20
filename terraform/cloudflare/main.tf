terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.0"
    }
  }
}

data "cloudflare_zone" "default" {
  name = var.domain_name
}

resource "cloudflare_record" "ipv4" {
  for_each = var.records

  zone_id = data.cloudflare_zone.default.id
  name    = each.key
  value   = each.value.ipv4
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ipv6" {
  for_each = var.records

  zone_id = data.cloudflare_zone.default.id
  name    = each.key
  value   = each.value.ipv6
  type    = "AAAA"
  proxied = true
}

resource "cloudflare_authenticated_origin_pulls_certificate" "per_hostname" {
  depends_on = [
    cloudflare_record.ipv4,
    cloudflare_record.ipv6
  ]
  for_each = var.records

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
  for_each = var.records

  zone_id  = data.cloudflare_zone.default.id
  enabled  = true
  hostname = "${each.key}.${var.domain_name}"
}