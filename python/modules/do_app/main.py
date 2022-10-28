#!/usr/bin/env python
from constructs import Construct
from cdktf import App
from cdktf_tf_module_stack import TFModuleStack, ProviderRequirement, TFModuleOutput, TFModuleVariable
from cdktf_cdktf_provider_digitalocean import App as DOApp

class DOStaticApp(TFModuleStack):
    def __init__(self, scope: Construct, id: str):
        super().__init__(scope, id)

        # define resources here
        ProviderRequirement(self,"digitalocean", ">=2.0.0","digitalocean/digitalocean")
    
        name = TFModuleVariable(self,"name", default="mame",description="Name of the applicaton to create in DigitalOcean") 
        region = TFModuleVariable(self,"region", default="ams",description="DgitalOcean region to deploy the app to") 

        app = DOApp(self, "static_site_example",
                spec={
                    "name": "static-site-" + name.string_value,
                    "region": region.string_value,

                    "static_site": [{
                        "name": name.string_value,
                        "sourceDir": "/src",

                        "github": {
                            "repo": "nicholasjackson/mame-wasm",
                            "deploy_on_push": True,
                            "branch": "main"
                        }
                    }],
                })

        TFModuleOutput(self,"uri",value=app.live_url, description="Live URL for the application")

app = App()
DOStaticApp(app, "do_static_app")

app.synth()
