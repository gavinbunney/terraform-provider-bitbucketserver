package bitbucket

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type PluginConfig struct {
	ValuesRaw json.RawMessage `json:"values"`
	Values    string
}

func resourcePluginConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginConfigCreateOrUpdate,
		Read:   resourcePluginConfigRead,
		Delete: resourcePluginConfigDelete,
		Update: resourcePluginConfigCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"config_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"values": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourcePluginConfigCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	configEndpoint := d.Get("config_endpoint").(string)
	values := d.Get("values").(string)
	config := []byte(values)
	payload := &bytes.Buffer{}
	_, err := payload.Write(config)
	if err != nil {
		panic(err)
	}

	_, err = client.Put(configEndpoint, payload)
	if err != nil {
		return err
	}

	d.SetId(configEndpoint)

	return resourcePluginConfigRead(d, m)
}

func resourcePluginConfigRead(d *schema.ResourceData, m interface{}) error {
	err := d.Set("config_endpoint", d.Id())
	if err != nil {
		return err
	}

	configEndpoint := d.Get("config_endpoint").(string)

	pluginConfig, err := readPluginConfig(m, configEndpoint)
	if err != nil {
		return err
	}

	err = d.Set("values", pluginConfig.Values)
	if err != nil {
		return err
	}

	return nil
}

func resourcePluginConfigDelete(d *schema.ResourceData, m interface{}) error {
	// Delete is a no-op
	return nil
}

func readPluginConfig(m interface{}, configEndpoint string) (PluginConfig, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient
	resp, err := client.Get(configEndpoint)
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
