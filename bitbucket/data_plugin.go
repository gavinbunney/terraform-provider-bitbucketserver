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

type PluginLicense struct {
	Valid                        bool     `json:"valid,omitempty"`
	Evaluation                   bool     `json:"evaluation,omitempty"`
	NearlyExpired                bool     `json:"nearlyExpired,omitempty"`
	MaintenanceExpiryDate        jsonTime `json:"maintenanceExpiryDate,omitempty"`
	MaintenanceExpired           bool     `json:"maintenanceExpired,omitempty"`
	LicenseType                  string   `json:"licenseType,omitempty"`
	ExpiryDate                   jsonTime `json:"expiryDate,omitempty"`
	RawLicense                   string   `json:"rawLicense,omitempty"`
	Renewable                    bool     `json:"renewable,omitempty"`
	OrganizationName             string   `json:"organizationName,omitempty"`
	ContactEmail                 string   `json:"contactEmail,omitempty"`
	Enterprise                   bool     `json:"enterprise,omitempty"`
	DataCenter                   bool     `json:"dataCenter,omitempty"`
	Subscription                 bool     `json:"subscription,omitempty"`
	Active                       bool     `json:"active,omitempty"`
	AutoRenewal                  bool     `json:"autoRenewal,omitempty"`
	Upgradable                   bool     `json:"upgradable,omitempty"`
	Crossgradeable               bool     `json:"crossgradeable,omitempty"`
	PurchasePastServerCutoffDate bool     `json:"purchasePastServerCutoffDate,omitempty"`
	SupportEntitlementNumber     string   `json:"supportEntitlementNumber,omitempty"`
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
			"vendor": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"marketplace_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"license": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"valid": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"evaluation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nearly_expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintenance_expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintenance_expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raw_license": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"renewable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"data_center": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"subscription": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renewal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"upgradable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crossgradeable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"purchase_past_server_cutoff_date": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"support_entitlement_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

	vendor := map[string]string{
		"name":             plugin.Vendor.Name,
		"link":             plugin.Vendor.Link,
		"marketplace_link": plugin.Vendor.MarketplaceLink,
	}
	_ = d.Set("vendor", vendor)

	// Hit the license API to get license details

	req, err = client.Get(fmt.Sprintf("/rest/plugins/1.0/%s-key/license", d.Get("key").(string)))
	if err != nil {
		return err
	}

	var pluginLicense PluginLicense

	body, readErr = ioutil.ReadAll(req.Body)
	if readErr != nil {
		return readErr
	}

	decodeErr = json.Unmarshal(body, &pluginLicense)
	if decodeErr != nil {
		return decodeErr
	}

	license := [1]map[string]interface{}{{
		"valid":                            pluginLicense.Valid,
		"evaluation":                       pluginLicense.Evaluation,
		"nearly_expired":                   pluginLicense.NearlyExpired,
		"maintenance_expiry_date":          pluginLicense.MaintenanceExpiryDate.String(),
		"maintenance_expired":              pluginLicense.MaintenanceExpired,
		"license_type":                     pluginLicense.LicenseType,
		"expiry_date":                      pluginLicense.ExpiryDate.String(),
		"raw_license":                      pluginLicense.RawLicense,
		"renewable":                        pluginLicense.Renewable,
		"organization_name":                pluginLicense.OrganizationName,
		"contact_email":                    pluginLicense.ContactEmail,
		"enterprise":                       pluginLicense.Enterprise,
		"data_center":                      pluginLicense.DataCenter,
		"subscription":                     pluginLicense.Subscription,
		"active":                           pluginLicense.Active,
		"auto_renewal":                     pluginLicense.AutoRenewal,
		"upgradable":                       pluginLicense.Upgradable,
		"crossgradeable":                   pluginLicense.Crossgradeable,
		"purchase_past_server_cutoff_date": pluginLicense.PurchasePastServerCutoffDate,
		"support_entitlement_number":       pluginLicense.SupportEntitlementNumber,
	}}
	_ = d.Set("license", license)

	return nil
}
