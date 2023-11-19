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