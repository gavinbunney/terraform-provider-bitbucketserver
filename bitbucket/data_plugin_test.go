package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataPlugin(t *testing.T) {
	config := `
		data "bitbucketserver_plugin" "upm" {
			key = "com.atlassian.upm.atlassian-universal-plugin-manager-plugin"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "key", "com.atlassian.upm.atlassian-universal-plugin-manager-plugin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "enabled", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "enabled_by_default", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.upm", "version"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "name", "Atlassian Universal Plugin Manager Plugin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "description", "This is the plugin that provides the Atlassian Universal Plugin Manager."),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "user_installed", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "optional", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor_name", "Atlassian"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor_link", "http://www.atlassian.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor_marketplace_link", "http://www.atlassian.com"),
				),
			},
		},
	})
}
