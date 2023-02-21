package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataRepositoryAAA(t *testing.T) {
	randKey := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-project-%v"
			description = "test project"
			public = "true"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "test-repo-%v"
			description = "test repo"
		}
		
		data "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			slug = bitbucketserver_repository.test.slug
		}
	`, randKey, randKey, randKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_repository.test", "project", fmt.Sprintf("TEST%v", randKey)),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository.test", "slug", fmt.Sprintf("test-repo-%v", randKey)),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository.test", "name", fmt.Sprintf("test-repo-%v", randKey)),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository.test", "public", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository.test", "description", "test repo")),
			},
		},
	})
}
