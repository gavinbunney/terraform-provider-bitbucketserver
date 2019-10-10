package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketDataGroups_basic(t *testing.T) {
	config := `
		data "bitbucketserver_groups" "test" {
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
					resource.TestCheckResourceAttr("data.bitbucketserver_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_groups.test", "groups.0.name", "stash-users"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGroups_additional(t *testing.T) {
	config := `
		resource "bitbucketserver_group" "test" {
			name = "test-group"
		}

		data "bitbucketserver_groups" "test" {
			filter = bitbucketserver_group.test.name
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_groups.test", "groups.0.name", "test-group"),
				),
			},
		},
	})
}
