package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketProject(t *testing.T) {
	var repo Repository

	testAccBitbucketProjectConfig := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketProjectExists("bitbucketserver_project.test", &repo),
				),
			},
		},
	})
}

func testAccCheckBitbucketProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucketserver_project.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_project.test")
	}

	response, _ := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s", rs.Primary.Attributes["key"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("project still exists")
	}

	return nil
}

func testAccCheckBitbucketProjectExists(n string, repository *Repository) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no project ID is set")
		}
		return nil
	}
}
