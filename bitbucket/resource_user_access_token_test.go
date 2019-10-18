package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccBitbucketUserAccessToken(t *testing.T) {
	config := `
		resource "bitbucketserver_user_access_token" "test" {
            user = "admin"
			name = "my-token"
			permissions = ["REPO_READ", "PROJECT_WRITE"]
		}
	`

	configModified := strings.ReplaceAll(config, "my-token", "my-updated-token")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketUserAccessTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "name", "my-token"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.#", "2"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.0", "REPO_READ"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.1", "PROJECT_WRITE"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "created_date"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "last_authenticated"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "access_token"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "name", "my-updated-token"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.#", "2"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.0", "REPO_READ"),
					resource.TestCheckResourceAttr("bitbucketserver_user_access_token.test", "permissions.1", "PROJECT_WRITE"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "created_date"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "last_authenticated"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user_access_token.test", "access_token"),
				),
			},
		},
	})
}

func testAccCheckBitbucketUserAccessTokenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketServerProvider).BitbucketClient
	rs, ok := s.RootModule().Resources["bitbucketserver_user_access_token.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_user_access_token.test")
	}

	response, _ := client.Get(fmt.Sprintf("/rest/access-tokens/1.0/users/%s/%s", rs.Primary.Attributes["user"], rs.Primary.ID))

	if response.StatusCode != 404 {
		return fmt.Errorf("access token still exists")
	}

	return nil
}
