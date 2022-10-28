#!/usr/bin/env python
from os import path
from constructs import Construct
from cdktf import App, TerraformStack, TerraformVariable, TerraformOutput, Fn
from cdktf_cdktf_provider_digitalocean import App as DOApp, DigitaloceanProvider
from cdktf_cdktf_provider_cloudflare import CloudflareProvider, DataCloudflareZone, Record, WorkerScript, WorkerRoute


class MyStack(TerraformStack):
    def __init__(self, scope: Construct, ns: str):
        super().__init__(scope, ns)

        DigitaloceanProvider(self, "digitalocean")
        CloudflareProvider(self, "cloudflare")
        
        domain = TerraformVariable(self, "domain",
                default="demo.gs")

        subdomain = TerraformVariable(self, "subdomain",
                default="mame-python")

        cloudflareZone = DataCloudflareZone(self, "cloudflare_domain", 
                name=domain.string_value) 

        app = DOApp(self, "static_site_example",
                spec={
                    "name": "static-site-" + subdomain.string_value,
                    "region": "ams",

                    "static_site": [{
                        "name": subdomain.string_value,
                        "sourceDir": "/src",

                        "github": {
                            "repo": "nicholasjackson/mame-wasm",
                            "deploy_on_push": True,
                            "branch": "main"
                        }
                    }],
                })

        record = Record(self,"mame",
                    zone_id=cloudflareZone.zone_id,
                    name=subdomain.string_value,
                    value="192.0.2.1",
                    type="A",
                    proxied=True)

        script = WorkerScript(self, "redirect_script", 
                    name="proxy-" + subdomain.string_value,
                    content=Fn.templatefile(path.abspath("./workers/proxy.js"), {"hostname": Fn.trimprefix(app.live_url, "https://")}))

        WorkerRoute(self, "proxy_route",
                    zone_id=cloudflareZone.zone_id,
                    pattern=record.hostname + "/*",
                    script_name=script.name)

        TerraformOutput(self, "cloudflare_zone",
                value=cloudflareZone.zone_id)
        
        TerraformOutput(self, "cloudflare_url",
                value="https://" + record.hostname)
        
        TerraformOutput(self, "digitalocean_zone",
                value=app.live_url)
        
        

app = App()
MyStack(app, "cdktf-python")

app.synth()
