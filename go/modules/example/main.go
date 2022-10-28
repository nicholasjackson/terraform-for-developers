package main

import (
	"cdk.tf/go/stack/generated/stacks/do_static_app"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-digitalocean-go/digitalocean/v2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	// The code that defines your stack goes here
	digitalocean.NewDigitaloceanProvider(stack, jsii.String("digitalocean"), &digitalocean.DigitaloceanProviderConfig{})
  app := do_static_app.NewDoStaticApp(stack, jsii.String("app"),
    &do_static_app.DoStaticAppOptions{
      Name: "test",
    },
  )

  cdktf.NewTerraformOutput(stack, jsii.String("uri"), &cdktf.TerraformOutputConfig{Value: app.UriOutput()})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "example")

	app.Synth()
}
