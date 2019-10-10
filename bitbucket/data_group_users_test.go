package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketDataGroupUsers_check_default_included(t *testing.T) {
	config := `
		data "bitbucketserver_group_users" "test" {
			group = "stash-users"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.active", "true"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGroupUsers_additional(t *testing.T) {
	config := `
		resource "bitbucketserver_user" "mreynolds" {
		  name          = "mreynolds"
		  display_name  = "Malcolm Reynolds"
		  email_address = "browncoat@example.com"
		}

		resource "bitbucketserver_user_group" "test" {
			user = bitbucketserver_user.mreynolds.name
			group = "stash-users"
		}

		data "bitbucketserver_group_users" "test" {
			group = bitbucketserver_user_group.test.group
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.#", "2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.1.name", "mreynolds"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.1.display_name", "Malcolm Reynolds"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.1.email_address", "browncoat@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_group_users.test", "users.1.active", "true"),
				),
			},
		},
	})
}
