package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketResourceUserGroup_basic(t *testing.T) {
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
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user_group.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_user_group.test", "group", "stash-users"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceUserGroup_new_group(t *testing.T) {
	config := `
		resource "bitbucketserver_group" "test" {
			name = "test-group"
		}

		resource "bitbucketserver_user" "mreynolds" {
		  name          = "mreynolds"
		  display_name  = "Malcolm Reynolds"
		  email_address = "browncoat@example.com"
		}

		resource "bitbucketserver_user_group" "test" {
			user = bitbucketserver_user.mreynolds.name
			group = bitbucketserver_group.test.name
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user_group.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_user_group.test", "group", "test-group"),
				),
			},
		},
	})
}
