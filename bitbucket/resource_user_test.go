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

func TestAccBitbucketUser(t *testing.T) {
	userRand := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_user" "test" {
			name = "admin%v"
			display_name = "Admin %v"
			email_address = "admin%v@example.com"
		}
	`, userRand, userRand, userRand)

	configModified := strings.ReplaceAll(config, "Admin ", "Admin Updated ")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "name", "admin"+userRand),
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "display_name", "Admin "+userRand),
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "email_address", "admin"+userRand+"@example.com"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user.test", "initial_password"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "name", "admin"+userRand),
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "display_name", "Admin Updated "+userRand),
					resource.TestCheckResourceAttr("bitbucketserver_user.test", "email_address", "admin"+userRand+"@example.com"),
					resource.TestCheckResourceAttrSet("bitbucketserver_user.test", "initial_password"),
				),
			},
		},
	})
}

func testAccCheckBitbucketUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketServerProvider).BitbucketClient
	rs, ok := s.RootModule().Resources["bitbucketserver_user.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_user.test")
	}

	response, _ := client.Get(fmt.Sprintf("/rest/api/1.0/users/%s", rs.Primary.Attributes["name"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("user still exists")
	}

	return nil
}
