package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataProjectAAA(t *testing.T) {
	projectKey := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-project-%v"
			description = "test project"
			public = "true"
		}

		data "bitbucketserver_project" "test" {
			key = bitbucketserver_project.test.key
		}
	`, projectKey, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project.test", "name", fmt.Sprintf("test-project-%v", projectKey)),
					resource.TestCheckResourceAttr("data.bitbucketserver_project.test", "key", fmt.Sprintf("TEST%v", projectKey)),
					resource.TestCheckResourceAttr("data.bitbucketserver_project.test", "description", "test project"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project.test", "public", "true"),
				),
			},
		},
	})
}
