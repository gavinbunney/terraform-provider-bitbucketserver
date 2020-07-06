package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketPlugin_install(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
	config := `
		resource "bitbucketserver_plugin" "test" {
			key = "com.plugin.commitgraph.commitgraph"
			version = "5.3.3"
			license = "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "key", "com.plugin.commitgraph.commitgraph"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "enabled", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "version", "5.3.3"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "license", "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "enabled_by_default", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "name", "Charts and Graphs for Bitbucket Server"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "description", "Gain insight into Bitbucket with charts and graphs that help you visualize user contributions, commits, and team activity."),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "user_installed", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "optional", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "vendor.name", "Mohami"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "vendor.link", "https://mohami.io"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "vendor.marketplace_link", "https://mohami.io"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.valid", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.evaluation", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.nearly_expired", "true"),
					resource.TestCheckResourceAttrSet("bitbucketserver_plugin.test", "applied_license.0.maintenance_expiry_date"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.maintenance_expired", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.license_type", "DEVELOPER"),
					resource.TestCheckResourceAttrSet("bitbucketserver_plugin.test", "applied_license.0.expiry_date"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.raw_license", "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.renewable", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.organization_name", "Atlassian"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.enterprise", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.data_center", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.subscription", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.active", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.auto_renewal", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.upgradable", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.crossgradeable", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin.test", "applied_license.0.purchase_past_server_cutoff_date", "true"),
				),
			},
		},
	})
}
