package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketProject(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-repo-for-repository-test"
			description = "My description"
			avatar = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg=="
		}
	`, projectKey)

	configModified := strings.ReplaceAll(config, "My description", "My updated description")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketProjectExists("bitbucketserver_project.test"),
					resource.TestCheckResourceAttr("bitbucketserver_project.test", "key", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project.test", "description", "My description"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketProjectExists("bitbucketserver_project.test"),
					resource.TestCheckResourceAttr("bitbucketserver_project.test", "key", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_project.test", "description", "My updated description"),
				),
			},
		},
	})
}

func testAccCheckBitbucketProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketServerProvider).BitbucketClient
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

func testAccCheckBitbucketProjectExists(n string) resource.TestCheckFunc {
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
