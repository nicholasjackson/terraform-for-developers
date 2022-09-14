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

data "cloudflare_zone" "domain" {
  name = var.cloudflare_domain
}

resource "digitalocean_app" "static_site_example" {
  spec {
    name   = "static-site-mame"
    region = "ams"

    static_site {
      name       = "mame-sf2"
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
  zone_id = data.cloudflare_zone.domain.zone_id
  name    = "mame"
  value   = "192.0.2.1"
  type    = "A"
  proxied = true
}

resource "cloudflare_worker_script" "redirect_script" {
  name    = "proxy-3343424sfdksjfsf"
  content = templatefile("./workers/proxy.js", { hostname = trimprefix(digitalocean_app.static_site_example.live_url, "https://") })
}

resource "cloudflare_worker_route" "proxy_route" {
  zone_id     = data.cloudflare_zone.domain.zone_id
  pattern     = "${cloudflare_record.mame.hostname}/*"
  script_name = cloudflare_worker_script.redirect_script.name
}

output "cloudflare_zone" {
  value = data.cloudflare_zone.domain.zone_id
}

output "do_url" {
  value = digitalocean_app.static_site_example.live_url
}
