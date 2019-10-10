package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceProjectPermissionsUser(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_project_permissions_user" "test" {
			project = bitbucketserver_project.test.key
			user = "admin2"
			permission = "PROJECT_READ"
		}
	`, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "id", projectKey+"/admin2"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "user", "admin2"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "permission", "PROJECT_READ"),
				),
			},
		},
	})
}
