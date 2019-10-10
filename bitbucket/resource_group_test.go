package bitbucket

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBitbucketResourceGroup_basic(t *testing.T) {
	config := `
		resource "bitbucketserver_group" "test" {
			name = "test-group"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_group.test", "name", "test-group"),
				),
			},
		},
	})
}
