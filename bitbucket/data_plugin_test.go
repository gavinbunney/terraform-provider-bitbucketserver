package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataPlugin_notifyer(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
	config := `
        resource "bitbucketserver_plugin" "test" {
			key = "nl.stefankohler.stash.stash-notification-plugin"
			version = "4.5.1"
			license = "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"
		}

		data "bitbucketserver_plugin" "test" {
			key = bitbucketserver_plugin.test.key
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "key", "nl.stefankohler.stash.stash-notification-plugin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.test", "enabled_by_default"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.test", "version"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "name", "Notifyr - Notification for Bitbucket"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "description", "Watch your repositories, branches, and tags and receive email notifications on changes."),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "user_installed", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "optional", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "vendor.name", "ASK Software"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "vendor.link", "http://www.stefankohler.nl/"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "vendor.marketplace_link", "http://www.stefankohler.nl/"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.valid", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.evaluation", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.nearly_expired", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.test", "applied_license.0.maintenance_expiry_date"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.maintenance_expired", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.license_type", "DEVELOPER"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.test", "applied_license.0.expiry_date"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.raw_license", "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.renewable", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.organization_name", "Atlassian"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.enterprise", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.data_center", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.subscription", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.active", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.auto_renewal", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.upgradable", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.crossgradeable", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.test", "applied_license.0.purchase_past_server_cutoff_date", "true"),
				),
			},
		},
	})
}

func TestAccBitbucketDataPlugin_upm(t *testing.T) {
	config := `
		data "bitbucketserver_plugin" "upm" {
			key = "com.atlassian.upm.atlassian-universal-plugin-manager-plugin"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "key", "com.atlassian.upm.atlassian-universal-plugin-manager-plugin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "enabled", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "enabled_by_default", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_plugin.upm", "version"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "name", "Atlassian Universal Plugin Manager Plugin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "description", "This is the plugin that provides the Atlassian Universal Plugin Manager."),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "user_installed", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "optional", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor.name", "Atlassian"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor.link", "http://www.atlassian.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "vendor.marketplace_link", "http://www.atlassian.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "applied_license.0.valid", "false"),
					resource.TestCheckResourceAttr("data.bitbucketserver_plugin.upm", "applied_license.0.active", "false"),
				),
			},
		},
	})
}
