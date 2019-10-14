package bitbucket

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type ApplicationProperties struct {
	Version     string `json:"version,omitempty"`
	BuildNumber string `json:"buildNumber,omitempty"`
	BuildDate   string `json:"buildDate,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func dataSourceApplicationProperties() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationPropertiesRead,

		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationPropertiesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get("/rest/api/1.0/application-properties")

	if err != nil {
		return err
	}

	if req.StatusCode == 200 {

		var applicationProperties ApplicationProperties

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &applicationProperties)
		if decodeerr != nil {
			return decodeerr
		}

		d.SetId(applicationProperties.Version)
		_ = d.Set("version", applicationProperties.Version)
		_ = d.Set("build_number", applicationProperties.BuildNumber)
		_ = d.Set("build_date", applicationProperties.BuildDate)
		_ = d.Set("display_name", applicationProperties.DisplayName)
	}

	return nil
}
