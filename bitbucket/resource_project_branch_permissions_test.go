package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceBranchPermission_requiredArgumentsOnly(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForRepositoryBasedTests(projectKey) + `
	resource "bitbucketserver_project_branch_permissions" "test" {
		project          = bitbucketserver_project.test.key
		repository       = bitbucketserver_repository.test.slug
		ref_pattern      = "refs/heads/master"
		type             = "pull-request-only"
	}`

	resourceName := "bitbucketserver_project_branch_permissions.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo|refs/heads/master|pull-request-only", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "ref_pattern", "refs/heads/master"),
					resource.TestCheckResourceAttr(resourceName, "type", "pull-request-only"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_id"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceBranchPermission_allArguments(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForRepositoryBasedTests(projectKey) + `
	resource "bitbucketserver_group" "test" {
		name = bitbucketserver_project.test.name
	}

	resource "bitbucketserver_group" "test_2" {
		name = format("%s-%s", bitbucketserver_project.test.name, "2")
	}

	resource "bitbucketserver_project_branch_permissions" "test" {
		project          = bitbucketserver_project.test.key
		repository       = bitbucketserver_repository.test.slug
		ref_pattern      = "refs/heads/master"
		type             = "pull-request-only"
		exception_users  = ["admin"]
		exception_groups = [bitbucketserver_group.test.name, bitbucketserver_group.test_2.name]
	}

	resource "bitbucketserver_project_branch_permissions" "test_2" {
		project     = bitbucketserver_project.test.key
		repository  = bitbucketserver_repository.test.slug
		ref_pattern = "refs/heads/master"
		type        = "no-deletes"
		depends_on  = [bitbucketserver_project_branch_permissions.test]
	}
	`

	resourceName := "bitbucketserver_project_branch_permissions.test"
	resourceName2 := "bitbucketserver_project_branch_permissions.test_2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo|refs/heads/master|pull-request-only", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "ref_pattern", "refs/heads/master"),
					resource.TestCheckResourceAttr(resourceName, "type", "pull-request-only"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_id"),
					resource.TestCheckResourceAttr(resourceName, "exception_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "exception_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "exception_groups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "exception_groups.0", fmt.Sprintf("test-project-%s", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "exception_groups.1", fmt.Sprintf("test-project-%s-2", projectKey)),

					resource.TestCheckResourceAttr(resourceName2, "type", "no-deletes"),
					resource.TestCheckResourceAttrSet(resourceName2, "permission_id"),
				),
			},
		},
	})
}
