package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceRepositoryPermissionsGroup(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "%v"
		}

		resource "bitbucketserver_repository" "test" {
			name    = "test"
			project = bitbucketserver_project.test.key
		}

		resource "bitbucketserver_repository_permissions_group" "test" {
			repository = bitbucketserver_repository.test.name
			project    = bitbucketserver_repository.test.project
			group      = "stash-users"
			permission = "REPO_WRITE"
		}
	`, projectKey, projectKey)

	configModified := strings.ReplaceAll(config, "REPO_WRITE", "REPO_READ")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "id", projectKey+"/test/stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "repository", "test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "group", "stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "permission", "REPO_WRITE"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "id", projectKey+"/test/stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "repository", "test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "group", "stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_group.test", "permission", "REPO_READ"),
				),
			},
		},
	})
}
