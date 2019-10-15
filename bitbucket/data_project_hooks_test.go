package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataProjectHooks_simple(t *testing.T) {
	projectKey := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-project-%v"
		}

		data "bitbucketserver_project_hooks" "test" {
			project = bitbucketserver_project.test.key
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.#", "8"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.key", "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:all-approvers-merge-check"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.name", "All reviewers approve"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.type", "PRE_PULL_REQUEST_MERGE"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.description", "Require all reviewers to approve the pull request."),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.version", "6.7.0"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.#", "2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.0", "PROJECT"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.1", "REPOSITORY"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.configured", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_type", "PROJECT"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_project_hooks.test", "hooks.0.scope_resource_id"),
				),
			},
		},
	})
}

func TestAccBitbucketDataProjectHooks_type(t *testing.T) {
	projectKey := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-project-%v"
		}

		data "bitbucketserver_project_hooks" "test" {
			project = bitbucketserver_project.test.key
            type    = "PRE_RECEIVE"
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.#", "3"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.key", "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.name", "Reject Force Push"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.type", "PRE_RECEIVE"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.description", "Reject all force pushes (git push --force) to this repository"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.version", "6.7.0"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.#", "2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.0", "PROJECT"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_types.1", "REPOSITORY"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.configured", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_hooks.test", "hooks.0.scope_type", "PROJECT"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_project_hooks.test", "hooks.0.scope_resource_id"),
				),
			},
		},
	})
}
