package bitbucket

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceGlobalPermissionsGroup(t *testing.T) {
	groupName := fmt.Sprintf("test-group-%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	config := fmt.Sprintf(`
		resource "bitbucketserver_group" "test" {
			name = "%v"
		}

		resource "bitbucketserver_global_permissions_group" "test" {
			group = bitbucketserver_group.test.name
			permission = "ADMIN"
		}
	`, groupName)

	configModified := strings.ReplaceAll(config, "ADMIN", "LICENSED_USER")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "id", groupName),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "group", groupName),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "permission", "ADMIN"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "id", groupName),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "group", groupName),
					resource.TestCheckResourceAttr("bitbucketserver_global_permissions_group.test", "permission", "LICENSED_USER"),
				),
			},
		},
	})
}
