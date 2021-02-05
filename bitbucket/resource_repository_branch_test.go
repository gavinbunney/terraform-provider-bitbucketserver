package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceRepositoryBranch_simple(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "repo"
		}

		resource "bitbucketserver_repository_branch" "test" {
			project = bitbucketserver_project.test.key
			repository = bitbucketserver_repository.test.slug
			branch_name = "test-branch" 
			source_ref = "refs/head/master"
		}
	`, projectKey, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_respository_branch.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_respository_branch.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_respository_branch.test", "branch_name", "test-branch"),
				),
			},
		},
	})
}
