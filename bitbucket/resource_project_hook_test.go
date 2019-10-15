package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceProjectHook_simple(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_project_hook" "test" {
			project = bitbucketserver_project.test.key
			hook = "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook" 
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_project_hook.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project_hook.test", "hook", "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook"),
					resource.TestCheckResourceAttr("bitbucketserver_project_hook.test", "settings.%", "0"),
				),
			},
		},
	})
}
