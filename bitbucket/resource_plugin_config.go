package bitbucket

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

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
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validlicense": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
				ForceNew: false,
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
	key := d.Get("key").(string)
	values := d.Get("values").(string)
	config := []byte(values)
	payload := &bytes.Buffer{}
	_, err := payload.Write(config)
	if err != nil {
		panic(err)
	}

	_, err = client.Put(fmt.Sprintf("/rest/%s/1.0/config", url.QueryEscape(key)), payload)
	if err != nil {
		return err
	}

	d.SetId(key)

	return resourcePluginConfigRead(d, m)
}

func resourcePluginConfigRead(d *schema.ResourceData, m interface{}) error {
	err := d.Set("key", d.Id())
	if err != nil {
		return err
	}

	key := d.Get("key").(string)

	pluginConfig, err := readPluginConfig(m, key)
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

func resourcePluginConfigDelete(d *schema.ResourceData, m interface{}) error {
	// Delete is a no-op
	return nil
}
