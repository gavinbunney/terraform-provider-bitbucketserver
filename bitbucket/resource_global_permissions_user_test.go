package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceGlobalPermissionsUser(t *testing.T) {
	user := fmt.Sprintf("test-%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_user" "test" {
		  name          = "%v"
		  display_name  = "Test User"
		  email_address = "test@example.com"
		}

		resource "bitbucketserver_global_permissions_user" "test" {
			user = bitbucketserver_user.test.name
			permission = "SYS_ADMIN"
		}
	`, user)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_user.test", "id", user),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_user.test", "user", user),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_user.test", "permission", "SYS_ADMIN"),
				),
			},
		},
	})
}
