package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

type CloneUrl struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
}

type Repository struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	SCM         string `json:"scmId,omitempty"`
	Description string `json:"description,omitempty"`
	Forkable    bool   `json:"forkable,omitempty"`
	Public      bool   `json:"public,omitempty"`
	Links       struct {
		Clone []CloneUrl `json:"clone,omitempty"`
	} `json:"links,omitempty"`
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Update: resourceRepositoryUpdate,
		Read:   resourceRepositoryRead,
		Exists: resourceRepositoryExists,
		Delete: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "git",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forkable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"clone_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_https": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newRepositoryFromResource(d *schema.ResourceData) (Repo *Repository, Project string) {
	repo := &Repository{
		Name:        d.Get("name").(string),
		Slug:        d.Get("slug").(string),
		SCM:         d.Get("scm").(string),
		Description: d.Get("description").(string),
		Forkable:    d.Get("forkable").(bool),
		Public:      d.Get("public").(bool),
	}

	return repo, d.Get("project").(string)
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	repo, project := newRepositoryFromResource(d)

	bytedata, err := json.Marshal(repo)

	if err != nil {
		return err
	}

	repoSlug := determineSlug(d)

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	repo, project := newRepositoryFromResource(d)

	bytedata, err := json.Marshal(repo)

	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos",
		project,
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	d.SetId(string(fmt.Sprintf("%s/%s", project, repo.Name)))

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			d.Set("project", idparts[0])
			d.Set("slug", idparts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/slug`")
		}
	}

	repoSlug := determineSlug(d)
	project := d.Get("project").(string)

	client := m.(*BitbucketClient)
	repo_req, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	))

	if err != nil {
		return err
	}

	if repo_req.StatusCode == 200 {

		var repo Repository

		body, readerr := ioutil.ReadAll(repo_req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &repo)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("name", repo.Name)
		if repo.Slug != "" && repo.Name != repo.Slug {
			d.Set("slug", repo.Slug)
		}
		d.Set("scmId", repo.SCM)
		d.Set("description", repo.Description)
		d.Set("forkable", repo.Forkable)
		d.Set("public", repo.Public)

		for _, clone_url := range repo.Links.Clone {
			if clone_url.Name == "http" {
				d.Set("clone_https", clone_url.Href)
			} else {
				d.Set("clone_ssh", clone_url.Href)
			}
		}
	}

	return nil
}

func resourceRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)
	repoSlug := determineSlug(d)
	project := d.Get("project").(string)
	repo_req, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	))

	if err != nil {
		return false, fmt.Errorf("failed to get repository %s/%s from bitbucket: %+v", project, repoSlug, err)
	}

	if repo_req.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	repoSlug := determineSlug(d)
	project := d.Get("project").(string)
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	))

	return err
}

func determineSlug(d *schema.ResourceData) string {
	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}

	return repoSlug
}
