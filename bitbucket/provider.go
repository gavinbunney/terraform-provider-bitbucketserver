package bitbucket

import (
	"net/http"
	"strings"

	"github.com/gavinbunney/terraform-provider-bitbucketserver/bitbucket/marketplace"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_SERVER", nil),
			},
			"username": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
			},
		},
		ConfigureFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucketserver_application_properties":        dataSourceApplicationProperties(),
			"bitbucketserver_cluster":                       dataSourceCluster(),
			"bitbucketserver_global_permissions_groups":     dataSourceGlobalPermissionsGroups(),
			"bitbucketserver_global_permissions_users":      dataSourceGlobalPermissionsUsers(),
			"bitbucketserver_groups":                        dataSourceGroups(),
			"bitbucketserver_group_users":                   dataSourceGroupUsers(),
			"bitbucketserver_plugin":                        dataSourcePlugin(),
			"bitbucketserver_project_hooks":                 dataSourceProjectHooks(),
			"bitbucketserver_project_permissions_groups":    dataSourceProjectPermissionsGroups(),
			"bitbucketserver_project_permissions_users":     dataSourceProjectPermissionsUsers(),
			"bitbucketserver_repository_hooks":              dataSourceRepositoryHooks(),
			"bitbucketserver_repository_permissions_groups": dataSourceRepositoryPermissionsGroups(),
			"bitbucketserver_repository_permissions_users":  dataSourceRepositoryPermissionsUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"bitbucketserver_banner":                       resourceBanner(),
			"bitbucketserver_global_permissions_group":     resourceGlobalPermissionsGroup(),
			"bitbucketserver_global_permissions_user":      resourceGlobalPermissionsUser(),
			"bitbucketserver_group":                        resourceGroup(),
			"bitbucketserver_license":                      resourceLicense(),
			"bitbucketserver_mail_server":                  resourceMailServer(),
			"bitbucketserver_plugin":                       resourcePlugin(),
			"bitbucketserver_plugin_config":                resourcePluginConfig(),
			"bitbucketserver_project":                      resourceProject(),
			"bitbucketserver_project_hook":                 resourceProjectHook(),
			"bitbucketserver_project_permissions_group":    resourceProjectPermissionsGroup(),
			"bitbucketserver_project_permissions_user":     resourceProjectPermissionsUser(),
			"bitbucketserver_repository":                   resourceRepository(),
			"bitbucketserver_repository_hook":              resourceRepositoryHook(),
			"bitbucketserver_repository_permissions_group": resourceRepositoryPermissionsGroup(),
			"bitbucketserver_repository_permissions_user":  resourceRepositoryPermissionsUser(),
			"bitbucketserver_user":                         resourceUser(),
			"bitbucketserver_user_access_token":            resourceUserAccessToken(),
			"bitbucketserver_user_group":                   resourceUserGroup(),
		},
	}
}

type BitbucketServerProvider struct {
	BitbucketClient   *BitbucketClient
	MarketplaceClient *marketplace.Client
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	serverSanitized := d.Get("server").(string)
	if strings.HasSuffix(serverSanitized, "/") {
		serverSanitized = serverSanitized[0 : len(serverSanitized)-1]
	}

	b := &BitbucketClient{
		Server:     serverSanitized,
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		HTTPClient: &http.Client{},
	}

	m := &marketplace.Client{
		HTTPClient: &http.Client{},
	}

	return &BitbucketServerProvider{
		BitbucketClient:   b,
		MarketplaceClient: m,
	}, nil
}
