package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
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

		resource "bitbucketserver_user" "mreynolds" {
		  name          = "mreynolds"
		  display_name  = "Malcolm Reynolds"
		  email_address = "browncoat@example.com"
		}

		resource "bitbucketserver_project_permissions_user" "test" {
			project = bitbucketserver_project.test.key
			user = bitbucketserver_user.mreynolds.name
			permission = "PROJECT_READ"
		}
	`, projectKey)

	configModified := strings.ReplaceAll(config, "PROJECT_READ", "PROJECT_WRITE")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "id", projectKey+"/mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "permission", "PROJECT_READ"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "id", projectKey+"/mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "user", "mreynolds"),
					resource.TestCheckResourceAttr("bitbucketserver_project_permissions_user.test", "permission", "PROJECT_WRITE"),
				),
			},
		},
	})
}
