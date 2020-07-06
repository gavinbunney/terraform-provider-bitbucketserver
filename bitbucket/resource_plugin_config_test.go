package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketPluginConfig(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
	config := `
	resource "bitbucketserver_plugin" "test" {
		key     = "de.codecentric.atlassian.oidc.bitbucket-oidc-plugin"
		version = "1.5.1"
	}

	resource "bitbucketserver_plugin_config" "test" {
		config_endpoint  = "/rest/oidc/1.0/config"
		values = "{\"state\":\"OPTIONAL\",\"autoLogin\":\"DISABLED\",\"restAuthSso\":false,\"disableWebSudo\":false,\"issuerUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf\",\"authUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/authorize\",\"tokenUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/token\",\"logoutUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/logout\",\"checkSessionIFrameUrl\":\"\",\"jwkSetUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/keys\",\"usernameClaim\":\"preferred_username\",\"clientId\":\"123123123123123123\",\"clientSecret\":\"12312312312312312312312123123123123123123123123\",\"additionalAuthReqParams\":{\"scope\":\"groups\"},\"ssoButtonText\":\"OpenID Connect SSO\",\"redirectUrl\":\"\",\"createUsers\":true,\"createGroups\":false,\"updateUserInfo\":true,\"groupsClaim\":\"groups\",\"updateGroups\":true,\"requireGroups\":false,\"defaultGroups\":[],\"additionalGroups\":[\"stash-users\"]}"
		depends_on = [ bitbucketserver_plugin.test ]
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_plugin_config.test", "config_endpoint", "/rest/oidc/1.0/config"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin_config.test", "values", "{\"state\":\"OPTIONAL\",\"autoLogin\":\"DISABLED\",\"restAuthSso\":false,\"disableWebSudo\":false,\"issuerUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf\",\"authUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/authorize\",\"tokenUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/token\",\"logoutUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/logout\",\"checkSessionIFrameUrl\":\"\",\"jwkSetUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/keys\",\"usernameClaim\":\"preferred_username\",\"clientId\":\"123123123123123123\",\"clientSecret\":\"12312312312312312312312123123123123123123123123\",\"additionalAuthReqParams\":{\"scope\":\"groups\"},\"ssoButtonText\":\"OpenID Connect SSO\",\"redirectUrl\":\"\",\"createUsers\":true,\"createGroups\":false,\"updateUserInfo\":true,\"groupsClaim\":\"groups\",\"updateGroups\":true,\"requireGroups\":false,\"defaultGroups\":[],\"additionalGroups\":[\"stash-users\"]}"),
				),
			},
		},
	})
}
