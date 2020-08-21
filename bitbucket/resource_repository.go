package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CloneURL ...
type CloneURL struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
}

// Repository ...
type Repository struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Forkable    bool   `json:"forkable,omitempty"`
	Public      bool   `json:"public,omitempty"`
	Links       struct {
		Clone []CloneURL `json:"clone,omitempty"`
	} `json:"links,omitempty"`
}

// RepositoryForkProject ...
type RepositoryForkProject struct {
	Key string `json:"key,omitempty"`
}

// RepositoryFork ...
type RepositoryFork struct {
	Name    string                `json:"name,omitempty"`
	Project RepositoryForkProject `json:"project,omitempty"`
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
				ForceNew: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"fork_repository_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"fork_repository_slug": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_git_lfs": {
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

func newRepositoryFromResource(d *schema.ResourceData) (Repo *Repository) {
	repo := &Repository{
		Name:        d.Get("name").(string),
		Slug:        d.Get("slug").(string),
		Description: d.Get("description").(string),
		Forkable:    d.Get("forkable").(bool),
		Public:      d.Get("public").(bool),
	}

	return repo
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ServerProvider).Client
	project := d.Get("project").(string)
	repo := newRepositoryFromResource(d)

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

	err = handleRepositoryGitLFSChanges(client, project, repoSlug, d)
	if err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ServerProvider).Client

	project := d.Get("project").(string)
	repoSlug := determineSlug(d)
	name := d.Get("name").(string)

	forkProject := d.Get("fork_repository_project").(string)
	forkRepo := d.Get("fork_repository_slug").(string)
	if (forkProject != "" && forkRepo == "") || (forkRepo != "" && forkProject == "") {
		return fmt.Errorf("both fork_repository_project and fork_repository_slug need to be specified when forking an existing repository")
	}

	if forkProject != "" {
		err := createNewRepositoryFromFork(client, d, project, repoSlug, forkProject, forkRepo)
		if err != nil {
			return err
		}
	} else {
		err := createNewRepository(client, d, project)
		if err != nil {
			return err
		}
	}

	d.SetId(string(fmt.Sprintf("%s/%s", project, name)))

	err := handleRepositoryGitLFSChanges(client, project, repoSlug, d)
	if err != nil {
		return err
	}

	if forkProject != "" {
		// after forking a repository, run the update loop to update any names/descriptions etc of the forked repo
		return resourceRepositoryUpdate(d, m)
	}
	return resourceRepositoryRead(d, m)
}

func createNewRepository(client *Client, d *schema.ResourceData, project string) error {
	repo := newRepositoryFromResource(d)
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

	return nil
}

func createNewRepositoryFromFork(client *Client, d *schema.ResourceData, project string, repository string, forkProject string, forkRepository string) error {
	requestBody := &RepositoryFork{
		Name: repository,
		Project: RepositoryForkProject{
			Key: forkProject,
		},
	}

	bytedata, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s", forkProject, forkRepository), bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	return nil
}

func handleRepositoryGitLFSChanges(client *Client, project string, repoSlug string, d *schema.ResourceData) error {
	enableGitLFS := d.Get("enable_git_lfs").(bool)
	if (d.IsNewResource() && enableGitLFS) || d.HasChange("enable_git_lfs") {
		if enableGitLFS {
			_, err := client.Put(fmt.Sprintf("/rest/git-lfs/admin/projects/%s/repos/%s/enabled",
				project,
				repoSlug,
			), nil)

			if err != nil {
				return err
			}
		} else {
			_, err := client.Delete(fmt.Sprintf("/rest/git-lfs/admin/projects/%s/repos/%s/enabled",
				project,
				repoSlug,
			))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			_ = d.Set("project", idparts[0])
			_ = d.Set("slug", idparts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/slug`")
		}
	}

	repoSlug := determineSlug(d)
	project := d.Get("project").(string)

	client := m.(*ServerProvider).Client
	repoReq, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	))

	if err != nil {
		return err
	}

	if repoReq.StatusCode == 200 {

		var repo Repository

		body, readerr := ioutil.ReadAll(repoReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &repo)
		if decodeerr != nil {
			return decodeerr
		}

		_ = d.Set("name", repo.Name)
		if repo.Slug != "" && repo.Name != repo.Slug {
			_ = d.Set("slug", repo.Slug)
		}
		_ = d.Set("description", repo.Description)
		_ = d.Set("forkable", repo.Forkable)
		_ = d.Set("public", repo.Public)

		for _, cloneURL := range repo.Links.Clone {
			if cloneURL.Name == "http" {
				_ = d.Set("clone_https", cloneURL.Href)
			} else {
				_ = d.Set("clone_ssh", cloneURL.Href)
			}
		}

		gifLFS, err := client.Get(fmt.Sprintf("/rest/git-lfs/admin/projects/%s/repos/%s/enabled",
			project,
			repoSlug,
		))
		_ = d.Set("enable_git_lfs", err == nil && gifLFS.StatusCode == 200)
	}

	return nil
}

func resourceRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {

	var project = ""
	var repoSlug = ""
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			project = idparts[0]
			repoSlug = idparts[1]
		} else {
			return false, fmt.Errorf("incorrect ID format, should match `project/slug`")
		}
	}

	client := m.(*ServerProvider).Client
	repoReq, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s",
		project,
		repoSlug,
	))

	if err != nil {
		return false, fmt.Errorf("failed to get repository %s/%s from bitbucket: %+v", project, repoSlug, err)
	}

	if repoReq.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	repoSlug := determineSlug(d)
	project := d.Get("project").(string)
	client := m.(*ServerProvider).Client
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
