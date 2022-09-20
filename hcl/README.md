# Installing providers

To download and install the providers needed for this example run the following command.

```hcl
terraform init
```

# Creating resources

You can then run the command `terraform apply` to create the resources, all resources should use the free tiers of DigitalOcean and Cloudflare.

```shell
➜ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with
the following symbols:
  + create

Terraform will perform the following actions:

  # digitalocean_app.static_site_example will be created
  + resource "digitalocean_app" "static_site_example" {
      + active_deployment_id = (known after apply)
      + created_at           = (known after apply)
      + default_ingress      = (known after apply)
      + id                   = (known after apply)
      + live_url             = (known after apply)
      + updated_at           = (known after apply)
      + urn                  = (known after apply)

      + spec {
          + domains = (known after apply)
          + name    = "static-site-hcl"
          + region  = "ams"

          + domain {
              + name     = (known after apply)
              + type     = (known after apply)
              + wildcard = (known after apply)
              + zone     = (known after apply)
            }

          + static_site {
              + name       = "hcl"
              + source_dir = "/src"

              + github {
                  + branch         = "main"
                  + deploy_on_push = true
                  + repo           = "nicholasjackson/mame-wasm"
                }

              + routes {
                  + path                 = (known after apply)
                  + preserve_path_prefix = (known after apply)
                }
            }
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

digitalocean_app.static_site_example: Creating...
digitalocean_app.static_site_example: Still creating... [10s elapsed]
digitalocean_app.static_site_example: Still creating... [20s elapsed]
digitalocean_app.static_site_example: Still creating... [30s elapsed]
digitalocean_app.static_site_example: Still creating... [40s elapsed]
digitalocean_app.static_site_example: Creation complete after 44s [id=8b30ec0a-6384-415e-85ed-5c0d601e85c6]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

# Cleanup

To remove any resouces created by this example, run the following command

```shell
terraform destroy
```
