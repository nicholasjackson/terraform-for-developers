package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"path/filepath"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/codingconcepts/env"
	cloudflare "github.com/hashicorp/cdktf-provider-cloudflare-go/cloudflare/v2"
	digitalocean "github.com/hashicorp/cdktf-provider-digitalocean-go/digitalocean/v2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

// Ensure the config can be set from json or environment variables
type config struct {
  CloudFlareEnabled bool   `json:"cloudflare_enabled" env:"VAR_cloudflare_enabled"`
  Name              string `json:"name" env:"VAR_name"`
  Domain            string `json:"domain" env:"VAR_domain"`
  Region            string `json:"region" env:"VAR_region"`
}

func loadConfig(path string) (*config,error) {
  d,err := ioutil.ReadFile(path)
  if err != nil {
    return nil, fmt.Errorf("unable to read config file: %s, error: %s", path, err)
  }

  c:=&config{}
  err = json.Unmarshal(d, c)
  if err != nil {
    return nil, fmt.Errorf("unable to process config as json, file: %s, error: %s", path, err)
  }

  // attempt to process environment variables
  err = env.Set(c)
  if err != nil {
    return nil, fmt.Errorf("unable to process environment variables, error: %s", err)
  }

  return c, nil
}

func NewMyStack(scope constructs.Construct, id string, c *config) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	digitalocean.NewDigitaloceanProvider(stack, jsii.String("digitalocean"), &digitalocean.DigitaloceanProviderConfig{})
	cloudflare.NewCloudflareProvider(stack, jsii.String("cloudflare"), &cloudflare.CloudflareProviderConfig{})

  app := createDOApp(stack, c.Name, c.Region)

  if c.CloudFlareEnabled {
    createCloudflare(stack, c.Domain, c.Region, *app.LiveUrl())
  }

	return stack
}

func createDOApp(stack cdktf.TerraformStack, name, region string) digitalocean.App {
  app := digitalocean.NewApp(stack, jsii.String("static_site_example"), &digitalocean.AppConfig{
		Spec: &digitalocean.AppSpec{
			Name:   jsii.String(fmt.Sprintf("static-site-%s", name)),
			Region: jsii.String(region),
			StaticSite: []*digitalocean.AppSpecStaticSite{
				{
					Name:      jsii.String(name),
					SourceDir: jsii.String("/src"),

					Github: &digitalocean.AppSpecStaticSiteGithub{
						Repo:         jsii.String("nicholasjackson/mame-wasm"),
						DeployOnPush: jsii.Bool(true),
						Branch:       jsii.String("main"),
					},
				},
			},
		},
	})
	
  cdktf.NewTerraformOutput(stack, jsii.String("digitalocean_url"), &cdktf.TerraformOutputConfig{
		Value: app.LiveUrl(),
	})

  return app
} 

func createCloudflare(stack cdktf.TerraformStack, domain, appName, doURL string) {
	cloudFlareZone := cloudflare.NewDataCloudflareZone(stack, jsii.String("zone_domain"), &cloudflare.DataCloudflareZoneConfig{
		Name: jsii.String(domain),
	})

	record := cloudflare.NewRecord(stack, jsii.String("mame"), &cloudflare.RecordConfig{
		ZoneId:  cloudFlareZone.ZoneId(),
		Name:    &appName,
		Value:   jsii.String("192.0.2.1"),
		Type:    jsii.String("A"),
		Proxied: jsii.Bool(true),
	})

	templatePath, _ := filepath.Abs("./workers/proxy.js")

	script := cloudflare.NewWorkerScript(stack, jsii.String("redirect_script"), &cloudflare.WorkerScriptConfig{
		Name: jsii.String(fmt.Sprintf("proxy-%s", appName)),
		Content: cdktf.Fn_Templatefile(
			jsii.String(templatePath),
			map[string]string{
				"hostname": *cdktf.Fn_Trimprefix(&doURL, jsii.String("https://")),
			},
		),
	})

	cloudflare.NewWorkerRoute(stack, jsii.String("proxy_route"), &cloudflare.WorkerRouteConfig{
		ZoneId:     cloudFlareZone.ZoneId(),
		Pattern:    jsii.String(fmt.Sprintf("%s/*", *record.Hostname())),
		ScriptName: script.Name(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("cloudflare_zone"), &cdktf.TerraformOutputConfig{
		Value: cloudFlareZone.ZoneId(),
	})
	
  cdktf.NewTerraformOutput(stack, jsii.String("cloudflare_url"), &cdktf.TerraformOutputConfig{
		Value: jsii.String(fmt.Sprintf("https://%s", *record.Hostname())),
	})
}

func main() {
  absPath, _ := filepath.Abs("./config.json")
  // load the config
  c,err := loadConfig(absPath)
  if err != nil {
    fmt.Printf("Error reading config: %s", err)
    return
  }

	app := cdktf.NewApp(nil)

	NewMyStack(app, "cdktf", c)

	app.Synth()
}
