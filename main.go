package main

import (
	"github.com/andrein/terraform-provider-jsonnet/jsonnet"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return jsonnet.Provider()
		},
	})
}