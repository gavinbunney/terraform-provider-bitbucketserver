package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataProjectPermissionsGroups_simple(t *testing.T) {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_project_permissions_group" "test" {
			project = bitbucketserver_project.test.key
			group = "stash-users"
			permission = "PROJECT_WRITE"
		}
		
		data "bitbucketserver_project_permissions_groups" "test" {
			project = bitbucketserver_project_permissions_group.test.project
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_groups.test", "groups.0.name", "stash-users"),
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_groups.test", "groups.0.permission", "PROJECT_WRITE"),
				),
			},
		},
	})
}

func TestAccBitbucketDataProjectPermissionsGroups_check_empty(t *testing.T) {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}	

		data "bitbucketserver_project_permissions_groups" "test" {
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
					resource.TestCheckResourceAttr("data.bitbucketserver_project_permissions_groups.test", "groups.#", "0"),
				),
			},
		},
	})
}
