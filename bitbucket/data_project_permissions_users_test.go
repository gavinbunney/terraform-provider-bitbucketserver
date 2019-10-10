package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataProjectPermissionsUsers_check_creator_included(t *testing.T) {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}	

		data "bitbucketserver_project_permissions_users" "test" {
			project = bitbucketserver_project.test.key
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.permission", "PROJECT_ADMIN"),
				),
			},
		},
	})
}

func TestAccBitbucketDataProjectPermissionsUsers_additional(t *testing.T) {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_project_permissions_user" "test" {
			project = bitbucketserver_project.test.key
			user = "admin2"
			permission = "PROJECT_WRITE"
		}
		
		data "bitbucketserver_project_permissions_users" "test" {
			project = bitbucketserver_project_permissions_user.test.project
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.#", "2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.display_name", "Admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.0.permission", "PROJECT_ADMIN"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.1.name", "admin2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.1.display_name", "Admin 2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.1.email_address", "admin2@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.1.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_users.test", "users.1.permission", "PROJECT_WRITE"),
				),
			},
		},
	})
}
