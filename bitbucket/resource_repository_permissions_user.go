package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
	"strings"
)

func resourceRepositoryPermissionsUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryPermissionsUserCreate,
		Update: resourceRepositoryPermissionsUserUpdate,
		Read:   resourceRepositoryPermissionsUserRead,
		Delete: resourceRepositoryPermissionsUserDelete,
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
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"REPO_READ", "REPO_WRITE", "REPO_ADMIN"}, false),
			},
		},
	}
}

func resourceRepositoryPermissionsUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/users?permission=%s&name=%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("user").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceRepositoryPermissionsUserRead(d, m)
}

func resourceRepositoryPermissionsUserCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceRepositoryPermissionsUserUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("user").(string)))
	return resourceRepositoryPermissionsUserRead(d, m)
}

func resourceRepositoryPermissionsUserRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("user", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/user`")
		}
	}

	user := d.Get("user").(string)
	users, err := readRepositoryPermissionsUsers(m, d.Get("project").(string), d.Get("repository").(string), user)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range users {
		if g.Name == user {
			_ = d.Set("permission", g.Permission)
			break
		}
	}

	return nil
}

func resourceRepositoryPermissionsUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/users?name=%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Get("user").(string)),
	))

	return err
}
