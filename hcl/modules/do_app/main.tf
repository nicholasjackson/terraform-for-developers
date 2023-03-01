terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
}


variable "name" {
  default = "mame-hcl"
}

resource "digitalocean_app" "static_site_example" {
  spec {
    name   = "static-site-${var.name}"
    region = "ams"

    static_site {
      name       = var.name
      source_dir = "/src"

      github {
        repo           = "nicholasjackson/mame-wasm"
        deploy_on_push = true
        branch         = "main"
      }
    }
  }
}

output "digitalocean_url" {
  value = digitalocean_app.static_site_example.live_url
}
