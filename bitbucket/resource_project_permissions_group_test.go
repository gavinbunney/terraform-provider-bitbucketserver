package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceProjectPermissionsGroup(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_project_permissions_group" "test" {
			project = bitbucketserver_project.test.key
			group = "stash-users"
			permission = "PROJECT_WRITE"
		}
	`, projectKey)

	configModified := strings.ReplaceAll(config, "PROJECT_WRITE", "PROJECT_READ")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "id", projectKey+"/stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "group", "stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "permission", "PROJECT_WRITE"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "id", projectKey+"/stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "group", "stash-users"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_group.test", "permission", "PROJECT_READ"),
				),
			},
		},
	})
}
