package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRepositoryWebHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryWebHookCreate,
		Update: resourceRepositoryWebHookUpdate,
		Read:   resourceRepositoryWebHookRead,
		Delete: resourceRepositoryWebHookDelete,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"events": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

type hookPostBody struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Events []string `json:"events"`
	URL    string   `json:"url"`
	Active bool     `json:"active"`
}

func resourceRepositoryWebHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	name := d.Get("name").(string)
	url := d.Get("url").(string)
	active := d.Get("active").(bool)
	events := d.Get("events").([]interface{})

	eventsStrings := []string{}
	for _, event := range events {
		eventsStrings = append(eventsStrings, event.(string))
	}

	settingsJson, err := json.Marshal(hookPostBody{
		ID:     strings.Split(d.Id(), "/")[2],
		Name:   name,
		Events: eventsStrings,
		URL:    url,
		Active: active,
	})
	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%s",
		project,
		repository,
		strings.Split(d.Id(), "/")[2],
	), bytes.NewBuffer(settingsJson))
	if err != nil {
		return err
	}

	return resourceRepositoryWebHookRead(d, m)
}

func resourceRepositoryWebHookCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	name := d.Get("name").(string)
	url := d.Get("url").(string)
	active := d.Get("active").(bool)
	events := d.Get("events").([]interface{})

	eventsStrings := []string{}
	for _, event := range events {
		eventsStrings = append(eventsStrings, event.(string))
	}

	settingsJson, err := json.Marshal(hookPostBody{
		Name:   name,
		Events: eventsStrings,
		URL:    url,
		Active: active,
	})
	if err != nil {
		return err
	}

	response, err := client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks",
		project,
		repository,
	), bytes.NewBuffer(settingsJson))

	if err != nil {
		return err
	}

	location, err := response.Location()
	if err != nil {
		return err
	}

	tokens := strings.Split(location.String(), "/")
	id := tokens[len(tokens)-1]

	d.SetId(project + "/" + repository + "/" + id)

	return resourceRepositoryWebHookRead(d, m)
}

func resourceRepositoryWebHookRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("id", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/id`")
		}
	}

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient
	resp, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%s",
		project,
		repository,
		strings.Split(d.Id(), "/")[2],
	))
	if err != nil {
		return err
	}

	_ = resp

	return nil
}

func resourceRepositoryWebHookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%s",
		d.Get("project").(string),
		d.Get("repository").(string),
		strings.Split(d.Id(), "/")[2]),
	)

	return err
}
