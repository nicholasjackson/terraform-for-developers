// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, TerraformStack, TerraformOutput } from "cdktf";
import { DoStaticApp } from "./.gen/modules/stacks/do_static_app"
import {DigitaloceanProvider} from "@cdktf/provider-digitalocean"

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    // define resources here
    const provider = new DigitaloceanProvider(this,"digitalocean");

    const app = new DoStaticApp(this, "app", {name:"typescriptmodule", providers: [provider]});

    new TerraformOutput(this,"uri",{value: app.uriOutput}); 
  }
}

const app = new App();
new MyStack(app, "example");
app.synth();
