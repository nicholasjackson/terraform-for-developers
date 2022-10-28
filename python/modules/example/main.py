#!/usr/bin/env python
from constructs import Construct
from cdktf import App, TerraformStack, TerraformOutput
from cdktf_cdktf_provider_digitalocean import DigitaloceanProvider
from imports.stacks import DoStaticApp

class MyStack(TerraformStack):
    def __init__(self, scope: Construct, ns: str):
        super().__init__(scope, ns)

        do = DigitaloceanProvider(self, "digitalocean")

        # define resources here
        #do_module = TerraformHclModule(self, "app",
        #        source="../typescript/do_app/modules/stacks/static_app",
        #        variables={"name": "test"},
        #        providers=[do],
        #        )
        
        #TerraformOutput(self, "uri", value=do_module.get_string("uri"))

        do_module = DoStaticApp(self,"mame",name="mytest",providers=[do]) 
        TerraformOutput(self, "uri", value=do_module.uri_output)



app = App()
MyStack(app, "example")

app.synth()
