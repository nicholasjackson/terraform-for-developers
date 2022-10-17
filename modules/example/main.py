#!/usr/bin/env python
from constructs import Construct
from cdktf import App, TerraformStack, TerraformHclModule

class MyStack(TerraformStack):
    def __init__(self, scope: Construct, ns: str):
        super().__init__(scope, ns)

        # define resources here
        TerraformHclModule(self, "app",
                source="../typescript/do_app",
                variables={"name": "test"}
                )


app = App()
MyStack(app, "example")

app.synth()
