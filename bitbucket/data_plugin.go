package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type Plugin struct {
	Key              string `json:"key,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	EnabledByDefault bool   `json:"enabledByDefault,omitempty"`
	Version          string `json:"version,omitempty"`
	Description      string `json:"description,omitempty"`
	Name             string `json:"name,omitempty"`
	UserInstalled    bool   `json:"userInstalled,omitempty"`
	Optional         bool   `json:"optional,omitempty"`
	Vendor           struct {
		Name            string `json:"name,omitempty"`
		MarketplaceLink string `json:"marketplaceLink,omitempty"`
		Link            string `json:"link,omitempty"`
	} `json:"vendor,omitempty"`
}

func dataSourcePlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_by_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vendor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor_marketplace_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePluginRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	req, err := client.Get(fmt.Sprintf("/rest/plugins/1.0/%s-key", d.Get("key").(string)))
	if err != nil {
		return err
	}

	if req.StatusCode == 200 {

		var plugin Plugin

		body, readErr := ioutil.ReadAll(req.Body)
		if readErr != nil {
			return readErr
		}

		decodeErr := json.Unmarshal(body, &plugin)
		if decodeErr != nil {
			return decodeErr
		}

		d.SetId(plugin.Key)
		_ = d.Set("enabled", plugin.Enabled)
		_ = d.Set("enabled_by_default", plugin.EnabledByDefault)
		_ = d.Set("version", plugin.Version)
		_ = d.Set("description", plugin.Description)
		_ = d.Set("name", plugin.Name)
		_ = d.Set("user_installed", plugin.UserInstalled)
		_ = d.Set("optional", plugin.Optional)
		_ = d.Set("vendor_name", plugin.Vendor.Name)
		_ = d.Set("vendor_link", plugin.Vendor.Link)
		_ = d.Set("vendor_marketplace_link", plugin.Vendor.MarketplaceLink)
	}

	return nil
}
