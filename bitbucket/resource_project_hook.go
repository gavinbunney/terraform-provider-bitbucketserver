package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectHookCreate,
		Update: resourceProjectHookUpdate,
		Read:   resourceProjectHookRead,
		Delete: resourceProjectHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hook": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"settings": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceProjectHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ServerProvider).Client

	project := d.Get("project").(string)
	hook := d.Get("hook").(string)
	settings := d.Get("settings").(map[string]interface{})

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/settings/hooks/%s/enabled",
		project,
		hook,
	), bytes.NewBuffer(settingsJSON))

	if err != nil {
		return err
	}

	return resourceProjectHookRead(d, m)
}

func resourceProjectHookCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceProjectHookUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("project").(string), d.Get("hook").(string)))
	return resourceProjectHookRead(d, m)
}

func resourceProjectHookRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 2 {
			_ = d.Set("project", parts[0])
			_ = d.Set("hook", parts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/hook`")
		}
	}

	project := d.Get("project").(string)
	hook := d.Get("hook").(string)

	client := m.(*ServerProvider).Client
	resp, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/settings/hooks/%s/settings",
		project,
		hook,
	))

	if err != nil {
		return err
	}

	var settings map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&settings)
	if err != nil {
		return err
	}

	_ = d.Set("settings", settings)

	return nil
}

func resourceProjectHookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ServerProvider).Client
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/settings/hooks/%s/enabled",
		d.Get("project").(string),
		d.Get("hook").(string)))

	return err
}
