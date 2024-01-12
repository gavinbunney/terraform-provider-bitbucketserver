package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bitbucket.Provider})
}
