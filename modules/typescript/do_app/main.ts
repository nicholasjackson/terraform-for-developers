import { Construct } from "constructs";
import { App, TerraformOutput } from "cdktf";
import { App as DOApp } from "@cdktf/provider-digitalocean"
import {TFModuleStack, TFModuleVariable, ProviderRequirement} from "cdktf-tf-module-stack"

class StaticApp extends TFModuleStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new ProviderRequirement(this, "digitalocean")

    const appName = new TFModuleVariable(this, "name", {
      default: "module",
      description: "Name of the DigitalOcean static site"
    })
    
    const region = new TFModuleVariable(this, "region", {
      default: "ams",
      description: "Region to deploy the app to"
    })

    // define resources here
    const app = new DOApp(this, "static_site", {
      spec: {
        name: appName.value,
        region: region.value,

        staticSite: [{
          name: appName.value,
          sourceDir: "/src",

          github: {
            repo: "nicholasjackson/mame-wasm",
            deployOnPush: true,
            branch: "main"
          }
        }]
      }
    });
    
    new TerraformOutput(this, "uri", {value: app.fqn}) 

  }
}

const app = new App();
new StaticApp(app, "static_app");
app.synth();
