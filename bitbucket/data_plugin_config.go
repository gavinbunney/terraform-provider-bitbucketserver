package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

type PluginConfig struct {
	ValidLicense bool            `json:"validLicense"`
	ValuesRaw    json.RawMessage `json:"values"`
	Values       string
}

func dataSourcePluginConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginConfigRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validlicense": {
				Type:     schema.TypeBool,
				Optional: true,
				Required: false,
				ForceNew: true,
			},
			"values": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				ForceNew: true,
			},
		},
	}
}

func dataSourcePluginConfigRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))

	pluginConfig, err := readPluginConfig(m, d.Id())
	if err != nil {
		return err
	}
	err = d.Set("validlicense", pluginConfig.ValidLicense)
	if err != nil {
		return err
	}
	err = d.Set("values", pluginConfig.Values)
	if err != nil {
		return err
	}
	return nil
}

func readPluginConfig(m interface{}, id string) (PluginConfig, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient
	resourceURL := fmt.Sprintf("/rest/%s/1.0/config", url.QueryEscape(id))

	resp, err := client.Get(resourceURL)
	if err != nil {
		return PluginConfig{}, err
	}
	var pluginConfig PluginConfig

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PluginConfig{}, err
	}

	err = json.Unmarshal(body, &pluginConfig)
	if err != nil {
		return PluginConfig{}, err
	}
	pluginConfig.Values = string(pluginConfig.ValuesRaw)
	return pluginConfig, nil
}
