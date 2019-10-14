package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

type PluginMarketplaceVersion struct {
	Version string `json:"name,omitempty"`
	Links   struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	Embedded struct {
		Artifact struct {
			Links struct {
				Self struct {
					Href string `json:"href,omitempty"`
				} `json:"self,omitempty"`
				Binary struct {
					Href string `json:"href,omitempty"`
				} `json:"binary,omitempty"`
			} `json:"_links,omitempty"`
		} `json:"artifact,omitempty"`
	} `json:"_embedded,omitempty"`
}

func (p *PluginMarketplaceVersion) Key() string {
	re, _ := regexp.Compile("/rest/2/addons/([a-zA-Z0-9.-]*)/versions/build/.*")
	values := re.FindStringSubmatch(p.Links.Self.Href)
	if len(values) > 0 {
		return values[1]
	}
	return ""
}

func (p *PluginMarketplaceVersion) Filename() string {
	ext := filepath.Ext(p.Embedded.Artifact.Links.Self.Href)
	return fmt.Sprintf("%s-%s%s", p.Key(), p.Version, ext)
}

func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginCreate,
		Update: resourcePluginUpdate,
		Read:   resourcePluginRead,
		Exists: resourcePluginExists,
		Delete: resourcePluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"license": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled_by_default": {
				Type:     schema.TypeBool,
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
			"applied_license": {
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

func resourcePluginCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*BitbucketServerProvider)

	key := d.Get("key").(string)
	version := d.Get("version").(string)

	marketplacePluginVersion, err := readMarketplacePluginVersion(key, version, provider)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(os.TempDir(), "*"+marketplacePluginVersion.Filename())
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	// download the plugin artifact for uploading to the API
	err = provider.MarketplaceClient.DownloadArtifact(marketplacePluginVersion.Embedded.Artifact.Links.Binary.Href, file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	// first get a token for interacting with the UPM
	resp, err := provider.BitbucketClient.Get("/rest/plugins/1.0/?os_authType=basic")
	if err != nil {
		return err
	}
	upmToken := resp.Header.Get("upm-token")

	// now we can use the token to upload the downloaded marketplace file to bitbucket
	_, err = provider.BitbucketClient.PostFileUpload("/rest/plugins/1.0/?token="+upmToken, nil, "plugin", file.Name())
	if err != nil {
		return err
	}

	d.SetId(key)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate),
		func() *resource.RetryError {
			exists, err := resourcePluginExists(d, m)
			if exists == false || err != nil {
				return resource.RetryableError(fmt.Errorf("Waiting for plugin installation to finish..."))
			} else {
				return nil
			}
		})
	if err != nil {
		return err
	}

	// need to also run an update loop to set enabled flags and license details
	err = resource.Retry(d.Timeout(schema.TimeoutCreate),
		func() *resource.RetryError {
			err := resourcePluginUpdate(d, m)
			if err != nil {
				return resource.RetryableError(fmt.Errorf("Waiting for plugin updates to finish..."))
			} else {
				return nil
			}
		})
	if err != nil {
		return err
	}

	return nil
}

func resourcePluginUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	key := d.Get("key").(string)

	if d.IsNewResource() || d.HasChange("enabled") {
		var plugin Plugin
		req, err := client.Do("GET", fmt.Sprintf("/rest/plugins/1.0/%s-key?os_authType=basic", key), nil, "application/vnd.atl.plugins.plugin+json")
		if err != nil {
			return nil
		}

		body, readErr := ioutil.ReadAll(req.Body)
		if readErr != nil {
			return readErr
		}

		decodeErr := json.Unmarshal(body, &plugin)
		if decodeErr != nil {
			return decodeErr
		}

		plugin.Enabled = d.Get("enabled").(bool)
		bytedata, err := json.Marshal(plugin)
		_, err = client.Do("PUT", fmt.Sprintf("/rest/plugins/1.0/%s-key?os_authType=basic", key), bytes.NewBuffer(bytedata), "application/vnd.atl.plugins.plugin+json")
		if err != nil {
			return err
		}
	}

	if d.IsNewResource() || d.HasChange("license") {
		license := d.Get("license").(string)
		if license != "" {
			licenseJson := map[string]string{"rawLicense": license}
			bytedata, err := json.Marshal(licenseJson)
			if err != nil {
				return err
			}

			req, err := client.Do("PUT", fmt.Sprintf("/rest/plugins/1.0/%s-key/license?os_authType=basic", key), bytes.NewBuffer(bytedata), "application/vnd.atl.plugins+json")

			// ignore 400 errors as this happens if the license is already applied
			if req == nil || (err != nil && req != nil && req.StatusCode != 400) {
				return err
			}
		} else {
			_, err := client.Do("DELETE", fmt.Sprintf("/rest/plugins/1.0/%s-key/license?os_authType=basic", key), nil, "application/vnd.atl.plugins+json")
			if err != nil {
				return err
			}
		}
	}

	return resourcePluginRead(d, m)
}

func resourcePluginRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		_ = d.Set("key", id)
	}

	client := m.(*BitbucketServerProvider).BitbucketClient
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
	_ = d.Set("applied_license", license)

	return nil
}

func resourcePluginExists(d *schema.ResourceData, m interface{}) (bool, error) {
	var key = ""
	id := d.Id()
	if id != "" {
		key = id
	} else {
		key = d.Get("key").(string)
	}

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/plugins/1.0/%s-key",
		key,
	))

	if err != nil {
		return false, fmt.Errorf("failed to get plugin %s from bitbucket: %+v", key, err)
	}

	if req.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourcePluginDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/plugins/1.0/%s-key",
		d.Get("key").(string),
	))

	return err
}

func readMarketplacePluginVersion(key string, version string, provider *BitbucketServerProvider) (*PluginMarketplaceVersion, error) {
	marketplaceRequest, err := provider.MarketplaceClient.Get(fmt.Sprintf("/rest/2/addons/%s/versions/name/%s", key, version))
	if err != nil {
		return nil, err
	}

	var marketplaceVersion PluginMarketplaceVersion

	body, readerr := ioutil.ReadAll(marketplaceRequest.Body)
	if readerr != nil {
		return nil, readerr
	}

	decodeerr := json.Unmarshal(body, &marketplaceVersion)
	if decodeerr != nil {
		return nil, decodeerr
	}

	return &marketplaceVersion, nil
}
