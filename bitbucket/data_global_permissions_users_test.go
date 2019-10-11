package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketDataGlobalPermissionsUsers_basic(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_users" "test" {
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.permission", "SYS_ADMIN"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGlobalPermissionsUsers_filter(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_users" "test" {
			filter = "admin"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.0.permission", "SYS_ADMIN"),
				),
			},
		},
	})
}

func TestAccBitbucketDataGlobalPermissionsUsers_filter_no_match(t *testing.T) {
	config := `
		data "bitbucketserver_global_permissions_users" "test" {
			filter = "admining-the-country-side"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_global_permissions_users.test", "users.#", "0"),
				),
			},
		},
	})
}
