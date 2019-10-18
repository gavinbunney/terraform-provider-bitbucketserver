package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketBanner_basic(t *testing.T) {
	config := `
		resource "bitbucketserver_banner" "test" {
			message = "Test Banner\n*bold*"
		}`

	configModified := `
		resource "bitbucketserver_banner" "test" {
			message = "Test Banner changed"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "message", "Test Banner\n*bold*"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "enabled", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "audience", "ALL"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "message", "Test Banner changed"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "enabled", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "audience", "ALL"),
				),
			},
		},
	})
}

func TestAccBitbucketBanner_authenticated(t *testing.T) {
	config := `
		resource "bitbucketserver_banner" "test" {
			message  = "Test Banner\n*bold*"
			audience = "AUTHENTICATED"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "message", "Test Banner\n*bold*"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "audience", "AUTHENTICATED"),
				),
			},
		},
	})
}
