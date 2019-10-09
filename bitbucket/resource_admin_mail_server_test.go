package bitbucket

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketAdminMailServer(t *testing.T) {
	testAccBitbucketAdminMailServerConfig := `
		resource "bitbucketserver_admin_mail_server" "test" {
			hostname = "mail.example.com"
			port = 465
			protocol = "SMTP"
			use_start_tls = true
			require_start_tls = true
			sender_address = "test@example.com"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketAdminMailServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketAdminMailServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketAdminMailServerExists("bitbucketserver_admin_mail_server.test"),
				),
			},
		},
	})
}

func testAccCheckBitbucketAdminMailServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	_, ok := s.RootModule().Resources["bitbucketserver_admin_mail_server.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_admin_mail_server.test")
	}

	response, _ := client.Get("/rest/api/1.0/admin/main-server")
	if response.StatusCode != 404 {
		return fmt.Errorf("mail-server configuration still exists")
	}

	return nil
}

func testAccCheckBitbucketAdminMailServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no mail ID is set")
		}
		return nil
	}
}
