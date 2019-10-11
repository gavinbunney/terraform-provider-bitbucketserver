package bitbucket

import (
	"net/http"
	"strings"

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
			"bitbucketserver_application_properties":     dataSourceApplicationProperties(),
			"bitbucketserver_global_permissions_groups":  dataSourceGlobalPermissionsGroups(),
			"bitbucketserver_global_permissions_users":   dataSourceGlobalPermissionsUsers(),
			"bitbucketserver_groups":                     dataSourceGroups(),
			"bitbucketserver_group_users":                dataSourceGroupUsers(),
			"bitbucketserver_project_permissions_groups": dataSourceProjectPermissionsGroups(),
			"bitbucketserver_project_permissions_users":  dataSourceProjectPermissionsUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"bitbucketserver_group":                     resourceGroup(),
			"bitbucketserver_license":                   resourceLicense(),
			"bitbucketserver_mail_server":               resourceMailServer(),
			"bitbucketserver_project":                   resourceProject(),
			"bitbucketserver_project_permissions_group": resourceProjectPermissionsGroup(),
			"bitbucketserver_project_permissions_user":  resourceProjectPermissionsUser(),
			"bitbucketserver_repository":                resourceRepository(),
			"bitbucketserver_user":                      resourceUser(),
			"bitbucketserver_user_group":                resourceUserGroup(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	serverSanitized := d.Get("server").(string)
	if strings.HasSuffix(serverSanitized, "/") {
		serverSanitized = serverSanitized[0 : len(serverSanitized)-1]
	}

	client := &BitbucketClient{
		Server:     serverSanitized,
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
