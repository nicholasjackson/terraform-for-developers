import * as path from "path";
import { Construct } from "constructs";
import { App, TerraformStack, TerraformVariable, TerraformOutput, Fn} from "cdktf";
import {DigitaloceanProvider, App as DOApp} from "@cdktf/provider-digitalocean";
import {CloudflareProvider, Record, DataCloudflareZone, WorkerScript, WorkerRoute } from "@cdktf/provider-cloudflare";

class MyStack extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);


    new DigitaloceanProvider(this, "digitalocean", {})
    new CloudflareProvider(this, "cloudflare", {})
    
    const domain = new TerraformVariable(this, "domain", {
      default: "demo.gs"
    });
    
    console.log("vars: %s", domain.stringValue);

    const subdomain = new TerraformVariable(this, "subdomain", {
      default: "mame-ts"
    });

    const zone = new DataCloudflareZone(this, "zone", {
      name: domain.value
    });

    const app = new DOApp(this, "static_site_example", {
      spec: {
        name: "static-site-" + subdomain.value,
        region: "ams",

        staticSite: [{
          name: subdomain.value,
          sourceDir: "/src",

          github: {
            repo: "nicholasjackson/mame-wasm",
            deployOnPush: true,
            branch: "main"
          }
        }]
      }
    });

    const record = new Record(this, "mame", {
      zoneId: zone.zoneId,
      name: subdomain.value,
      value: "192.0.2.1",
      type: "A",
      proxied: true
    });

    const scriptPath = path.resolve(__dirname, "./workers/proxy.js")

    const script = new WorkerScript(this, "redirect_script", {
      name: "proxy-" + subdomain.value,
      content: Fn.templatefile(scriptPath, {hostname: Fn.trimprefix(app.liveUrl, "https://")})
    });

    new WorkerRoute(this, "proxy_route", {
      zoneId: zone.zoneId,
      pattern: record.hostname + "/*",
      scriptName: script.name
    });

    new TerraformOutput(this, "cloudflare_zone", {
      value: zone.zoneId
    })
    
    new TerraformOutput(this, "cloudflare_url", {
      value: "https://" + record.hostname
    })
    
    new TerraformOutput(this, "digitalocean_url", {
      value: app.liveUrl
    })
  }
}

const app = new App();
new MyStack(app, "cdktf-ts");
app.synth();
