package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketRepository_basic(t *testing.T) {
	var repo Repository

	testAccBitbucketRepositoryConfig := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "test-repo-for-repository-test"
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
				),
			},
		},
	})
}

func TestAccBitbucketRepository_namewithspaces(t *testing.T) {
	var repo Repository

	testAccBitbucketRepositoryConfig := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "Test Repo For Repository Test"
			slug = "test-repo-for-repository-test"
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucketserver_repository.test_repo"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_repository.test_repo")
	}

	response, _ := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s", rs.Primary.Attributes["project"], rs.Primary.Attributes["slug"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("repository still exists")
	}

	return nil
}

func testAccCheckBitbucketRepositoryExists(n string, repository *Repository) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no repository ID is set")
		}
		return nil
	}
}
