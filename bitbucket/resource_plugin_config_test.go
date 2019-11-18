package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketPluginConfig(t *testing.T) {
	config := `
		resource "bitbucketserver_plugin_config" "test" {
			key = "oidc"
			values = "{\"state\":\"OPTIONAL\",\"autoLogin\":\"DISABLED\",\"restAuthSso\":false,\"disableWebSudo\":false,\"issuerUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf\",\"authUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/authorize\",\"tokenUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/token\",\"logoutUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/logout\",\"checkSessionIFrameUrl\":\"\",\"jwkSetUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/keys\",\"usernameClaim\":\"preferred_username\",\"clientId\":\"123123123123123123\",\"clientSecret\":\"12312312312312312312312123123123123123123123123\",\"additionalAuthReqParams\":{\"scope\":\"groups\"},\"ssoButtonText\":\"OpenID Connect SSO\",\"redirectUrl\":\"\",\"createUsers\":true,\"createGroups\":false,\"updateUserInfo\":true,\"groupsClaim\":\"groups\",\"updateGroups\":true,\"requireGroups\":false,\"defaultGroups\":[],\"additionalGroups\":[\"stash-users\"]}"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_plugin_config.test", "key", "oidc"),
					resource.TestCheckResourceAttr("bitbucketserver_plugin_config.test", "values", "{\"state\":\"OPTIONAL\",\"autoLogin\":\"DISABLED\",\"restAuthSso\":false,\"disableWebSudo\":false,\"issuerUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf\",\"authUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/authorize\",\"tokenUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/token\",\"logoutUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/logout\",\"checkSessionIFrameUrl\":\"\",\"jwkSetUrl\":\"https://example.oktapreview.com/oauth2/1234567890abcdf/v1/keys\",\"usernameClaim\":\"preferred_username\",\"clientId\":\"123123123123123123\",\"clientSecret\":\"12312312312312312312312123123123123123123123123\",\"additionalAuthReqParams\":{\"scope\":\"groups\"},\"ssoButtonText\":\"OpenID Connect SSO\",\"redirectUrl\":\"\",\"createUsers\":true,\"createGroups\":false,\"updateUserInfo\":true,\"groupsClaim\":\"groups\",\"updateGroups\":true,\"requireGroups\":false,\"defaultGroups\":[],\"additionalGroups\":[\"stash-users\"]}"),
				),
			},
		},
	})
}
