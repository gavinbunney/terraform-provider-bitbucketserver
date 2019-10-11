package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketDataGlobalPermissionsGroups_basic(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_groups" "test" {
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.0.name", "stash-users"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.0.permission", "LICENSED_USER"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGlobalPermissionsGroups_filter(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_groups" "test" {
			filter = "stash-users"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.0.name", "stash-users"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.0.permission", "LICENSED_USER"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGlobalPermissionsGroups_filter_no_match(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_groups" "test" {
			filter = "stashing"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_groups.test", "groups.#", "0"),
				),
			},
		},
	})
}
