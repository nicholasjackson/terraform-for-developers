#!/usr/bin/env python
from constructs import Construct
from cdktf import App, TerraformStack, TerraformHclModule, TerraformOutput
from cdktf_cdktf_provider_digitalocean import DigitaloceanProvider

class MyStack(TerraformStack):
    def __init__(self, scope: Construct, ns: str):
        super().__init__(scope, ns)

        do = DigitaloceanProvider(self, "digitalocean")

        # define resources here
        do_module = TerraformHclModule(self, "app",
                source="../typescript/do_app/modules/stacks/static_app",
                variables={"name": "test"},
                providers=[do],
                )

        TerraformOutput(self, "uri", value=do_module.get_string("uri"))


app = App()
MyStack(app, "example")

app.synth()
