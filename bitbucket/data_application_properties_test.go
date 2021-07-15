package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataApplicationProperties(t *testing.T) {
	config := `
		data "bitbucketserver_application_properties" "main" {}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "version", "6.10.0"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "build_number", "6010000"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "build_date", "1578626301641"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "display_name", "Bitbucket"),
				),
			},
		},
	})
}
