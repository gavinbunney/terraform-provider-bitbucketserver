package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketPluginConfig_basic(t *testing.T) {
	config := `
		data "bitbucketserver_plugin_config" "test" {
			key = "oidc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin_config.test", "id", "oidc"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin_config.test", "validlicense", "true"),
				),
			},
		},
	})
}
