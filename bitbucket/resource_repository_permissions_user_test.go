package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceRepositoryPermissionsUser(t *testing.T) {
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

		resource "bitbucketserver_user" "mreynolds" {
			name          = "mreynolds"
			display_name  = "Malcolm Reynolds"
			email_address = "browncoat@example.com"
		}

		resource "bitbucketserver_repository_permissions_user" "test" {
			project    = bitbucketserver_project.test.key
			repository = bitbucketserver_repository.test.name
			user       = bitbucketserver_user.mreynolds.name
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
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "id", projectKey+"/test/mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "repository", "test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "permission", "REPO_WRITE"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "id", projectKey+"/test/mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "repository", "test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_permissions_user.test", "permission", "REPO_READ"),
				),
			},
		},
	})
}
