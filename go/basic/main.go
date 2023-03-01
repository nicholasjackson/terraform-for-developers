package main

import (
	"fmt"
	"path/filepath"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	cloudflare "github.com/hashicorp/cdktf-provider-cloudflare-go/cloudflare/v2"
	digitalocean "github.com/hashicorp/cdktf-provider-digitalocean-go/digitalocean/v2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

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
	c, err := loadConfig(absPath)
	if err != nil {
		fmt.Printf("Error reading config: %s", err)
		return
	}

	app := cdktf.NewApp(nil)

	NewMyStack(app, "cdktf", c)

	app.Synth()
}
