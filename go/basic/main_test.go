package main

import (
	"testing"

	"github.com/hashicorp/cdktf-provider-digitalocean-go/digitalocean/v2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

// The tests below are example tests, you can find more information at
// https://cdk.tf/testing

var cfg = &config{
	CloudFlareEnabled: true,
	Name:              "testing",
}

var (
	fulltest = true
	stack    = NewMyStack(cdktf.Testing_App(nil), "stack", cfg)
	synth    = cdktf.Testing_Synth(stack, &fulltest)
)

func TestShouldContainContainer(t *testing.T) {
	s := stack
	_ = s

	result := cdktf.Testing_ToHaveResource(synth, digitalocean.App_TfResourceType())
	if !*result {
		t.Error("expected digital ocean app resource")
	}
}

func TestShouldHaveNameFromConfig(t *testing.T) {
	properties := map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "static-site-testing",
		},
	}

	result := cdktf.Testing_ToHaveResourceWithProperties(synth, digitalocean.App_TfResourceType(), &properties)

	if !*result {
		t.Error("expected digital ocean app resource, to have properties", *synth)
	}
}

func TestCheckValidity(t *testing.T) {
	assertion := cdktf.Testing_ToBeValidTerraform(cdktf.Testing_FullSynth(stack))

	if !*assertion {
		t.Error("invalid terraform config")
	}
}
