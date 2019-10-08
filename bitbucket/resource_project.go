package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type Project struct {
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Update: resourceProjectUpdate,
		Read:   resourceProjectRead,
		Exists: resourceProjectExists,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"avatar": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func newProjectFromResource(d *schema.ResourceData) *Project {
	project := &Project{
		Name:        d.Get("name").(string),
		Key:         d.Get("key").(string),
		Description: d.Get("description").(string),
		Public:      d.Get("public").(bool),
		Avatar:      d.Get("avatar").(string),
	}

	return project
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	bytedata, err := json.Marshal(project)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s",
		project.Key,
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	bytedata, err := json.Marshal(project)

	if err != nil {
		return err
	}

	_, err = client.Post("/rest/api/1.0/projects", bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	d.SetId(project.Key)

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		d.Set("key", id)
	}

	project := d.Get("key").(string)

	client := m.(*BitbucketClient)
	project_req, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s",
		project,
	))

	if err != nil {
		return err
	}

	if project_req.StatusCode == 200 {

		var project Project

		body, readerr := ioutil.ReadAll(project_req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &project)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("name", project.Name)
		d.Set("key", project.Key)
		d.Set("description", project.Description)
		d.Set("public", project.Public)
		d.Set("avatar", project.Avatar)
	}

	return nil
}

func resourceProjectExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)
	project := d.Get("key").(string)
	repo_req, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s",
		project,
	))

	if err != nil {
		return false, fmt.Errorf("failed to get project %s from bitbucket: %+v", project, err)
	}

	if repo_req.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	project := d.Get("key").(string)
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s",
		project,
	))

	return err
}
