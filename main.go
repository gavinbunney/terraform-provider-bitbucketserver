package main

import (
	"github.com/gavinbunney/terraform-provider-bitbucketserver/bitbucket"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return bitbucket.Provider()
		},
	})
}
