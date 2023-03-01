package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/codingconcepts/env"
)

// Ensure the config can be set from json or environment variables
type config struct {
	CloudFlareEnabled bool   `json:"cloudflare_enabled" env:"VAR_cloudflare_enabled"`
	Name              string `json:"name" env:"VAR_name"`
	Domain            string `json:"domain" env:"VAR_domain"`
	Region            string `json:"region" env:"VAR_region"`
}

func loadConfig(path string) (*config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %s, error: %s", path, err)
	}

	c := &config{}
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
