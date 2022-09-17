# Authentication

Terraform requires authentication in order to interact with the DigitalOcean and CloudFlare APIs. The following section shows how to obtain
and set the tokens for each provider.

## Digital Ocean Token

Fetch your token from 

(https://cloud.digitalocean.com/account/api/tokens)[https://cloud.digitalocean.com/account/api/tokens]

![](images/do_1.jpg)

If you complete the details and then press the `Generate` button you will be taken back to the other screen and the token will
be shown.

![](images/do_2.jpg)

Copy this token and set it as an environment variable, Terraform will automatically read this variable and automatically
pass it with any API request. Using an environment ensures that the API key is not hardcoded into the config and accidentally
leaks into the public domain. Leaking your API key will allow third parties to create infrastructure in your DigitalOcean
account.

```shell
export DIGITALOCEAN_TOKEN="dop_v1_a2e2a9b30b499fe62c64101528753c46fd4db80c1edfd3f5df2485c86344b78c"
```

## Cloudflare Token

You can get a token from the following URL:

https://dash.cloudflare.com/profile/api-tokens

If you click on 

```shell
export CLOUDFLARE_ACCOUNT_ID="8abde2bb370dcd685cbd1a84ec091a68"
export CLOUDFLARE_API_TOKEN="-655W6itNZ9sCKBphNCsBxOETIHTDzzT7BHmSX9J"
```
