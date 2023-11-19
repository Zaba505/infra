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
  zone_id = data.cloudflare_zone.default.id
  name    = var.record_name
  value   = var.ipv4_address
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ipv6" {
  count = var.ipv6_address != null ? 1 : 0

  zone_id = data.cloudflare_zone.default.id
  name    = var.record_name
  value   = var.ipv6_address
  type    = "AAAA"
  proxied = true
}