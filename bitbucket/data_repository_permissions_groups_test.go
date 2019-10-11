package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataRepositoryPermissionsGroups_simple(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key  = "%v"
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
		
		data "bitbucketserver_repository_permissions_groups" "test" {
			repository = bitbucketserver_repository_permissions_group.test.repository
			project    = bitbucketserver_repository_permissions_group.test.project
		}
	`, projectKey, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_groups.test", "groups.0.name", "stash-users"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_groups.test", "groups.0.permission", "REPO_WRITE"),
				),
			},
		},
	})
}

func TestAccBitbucketDataRepositoryPermissionsGroups_empty(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key  = "%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test" {
			name    = "test"
			project = bitbucketserver_project.test.key
		}

		data "bitbucketserver_repository_permissions_groups" "test" {
			repository = bitbucketserver_repository.test.name
			project    = bitbucketserver_repository.test.project
		}
	`, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_groups.test", "groups.#", "0"),
				),
			},
		},
	})
}
