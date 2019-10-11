package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataRepositoryPermissionsUsers_basic(t *testing.T) {
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

		resource "bitbucketserver_repository_permissions_user" "test" {
			repository = bitbucketserver_repository.test.name
			project    = bitbucketserver_project.test.key
			user       = "admin"
			permission = "REPO_WRITE"
		}

		data "bitbucketserver_repository_permissions_users" "test" {
			repository = bitbucketserver_repository_permissions_user.test.repository
			project    = bitbucketserver_repository_permissions_user.test.project
		}
	`, projectKey, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.0.permission", "REPO_WRITE"),
				),
			},
		},
	})
}

func TestAccBitbucketDataRepositoryPermissionsUsers_empty(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test" {
			name    = "test"
			project = bitbucketserver_project.test.key
		}

		data "bitbucketserver_repository_permissions_users" "test" {
			repository = bitbucketserver_repository.test.name
			project = bitbucketserver_project.test.key
		}
	`, projectKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_repository_permissions_users.test", "users.#", "0"),
				),
			},
		},
	})
}
