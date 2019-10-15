package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketBanner_basic(t *testing.T) {
	testAccBitbucketBannerConfig := `
		resource "bitbucketserver_banner" "test" {
			message = "Test Banner\n*bold*"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketBannerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "message", "Test Banner\n*bold*"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "enabled", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "audience", "ALL"),
				),
			},
		},
	})
}

func TestAccBitbucketBanner_authenticated(t *testing.T) {
	testAccBitbucketBannerConfig := `
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
				Config: testAccBitbucketBannerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "message", "Test Banner\n*bold*"),
					resource.TestCheckResourceAttr("bitbucketserver_banner.test", "audience", "AUTHENTICATED"),
				),
			},
		},
	})
}
