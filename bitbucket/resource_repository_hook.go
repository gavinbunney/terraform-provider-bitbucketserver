package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceRepositoryHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryHookCreate,
		Update: resourceRepositoryHookUpdate,
		Read:   resourceRepositoryHookRead,
		Delete: resourceRepositoryHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
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

func resourceRepositoryHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	hook := d.Get("hook").(string)
	settings := d.Get("settings").(map[string]interface{})

	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/hooks/%s/enabled",
		project,
		repository,
		hook,
	), bytes.NewBuffer(settingsJson))

	if err != nil {
		return err
	}

	return resourceRepositoryHookRead(d, m)
}

func resourceRepositoryHookCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceRepositoryHookUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("hook").(string)))
	return resourceRepositoryHookRead(d, m)
}

func resourceRepositoryHookRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("hook", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/hook`")
		}
	}

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	hook := d.Get("hook").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient
	resp, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/hooks/%s/settings",
		project,
		repository,
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

func resourceRepositoryHookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/hooks/%s/enabled",
		d.Get("project").(string),
		d.Get("repository").(string),
		d.Get("hook").(string)))

	return err
}
