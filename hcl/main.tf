terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 3.0"
    }
  }
}

provider "digitalocean" {
}

provider "cloudflare" {
}

variable "cloudflare_domain" {
  default = "demo.gs"
}

variable "cloudflare_subdomain" {
  default = "mame-hcl"
}

data "cloudflare_zone" "domain" {
  count = var.cloudflare_domain == "" ? 0 : 1
  name  = var.cloudflare_domain
}

resource "digitalocean_app" "static_site_example" {
  spec {
    name   = "static-site-${var.cloudflare_subdomain}"
    region = "ams"

    static_site {
      name       = var.cloudflare_subdomain
      source_dir = "/src"

      github {
        repo           = "nicholasjackson/mame-wasm"
        deploy_on_push = true
        branch         = "main"
      }
    }
  }
}

resource "cloudflare_record" "mame" {
  count = var.cloudflare_domain == "" ? 0 : 1

  zone_id = data.cloudflare_zone.domain.0.zone_id
  name    = var.cloudflare_subdomain
  value   = "192.0.2.1"
  type    = "A"
  proxied = true
}

resource "cloudflare_worker_script" "redirect_script" {
  count = var.cloudflare_domain == "" ? 0 : 1

  name    = "proxy-${var.cloudflare_subdomain}"
  content = templatefile("./workers/proxy.js", { hostname = trimprefix(digitalocean_app.static_site_example.live_url, "https://") })
}

resource "cloudflare_worker_route" "proxy_route" {
  count = var.cloudflare_domain == "" ? 0 : 1

  zone_id     = data.cloudflare_zone.domain.zone_id
  pattern     = "${cloudflare_record.mame.0.hostname}/*"
  script_name = cloudflare_worker_script.redirect_script.0.name
}

output "cloudflare_zone" {
  value = data.cloudflare_zone.domain.zone_id
}

output "cloudflare_url" {
  value = "https://${cloudflare_record.mame.0.hostname}"
}

output "digitalocean_url" {
  value = digitalocean_app.static_site_example.live_url
}

