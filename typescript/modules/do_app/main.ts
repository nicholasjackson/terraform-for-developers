import { Construct } from "constructs";
import { App } from "cdktf";
import { App as DOApp } from "@cdktf/provider-digitalocean"
import {TFModuleStack, TFModuleVariable, ProviderRequirement, TFModuleOutput} from "cdktf-tf-module-stack"

class DoStaticApp extends TFModuleStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new ProviderRequirement(this, "digitalocean",">=2.0.0", "digitalocean/digitalocean")

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
    
    new TFModuleOutput(this, "uri", {value: app.liveUrl, description: "Live URL for the DigitalOcean application"}) 

  }
}

const app = new App();
new DoStaticApp(app, "do_static_app");
app.synth();
