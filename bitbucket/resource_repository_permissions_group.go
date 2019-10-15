package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
	"strings"
)

func resourceRepositoryPermissionsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryPermissionsGroupCreate,
		Update: resourceRepositoryPermissionsGroupUpdate,
		Read:   resourceRepositoryPermissionsGroupRead,
		Delete: resourceRepositoryPermissionsGroupDelete,
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
			"group": {
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

func resourceRepositoryPermissionsGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/groups?permission=%s&name=%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("group").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceRepositoryPermissionsGroupRead(d, m)
}

func resourceRepositoryPermissionsGroupCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceRepositoryPermissionsGroupUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("group").(string)))
	return resourceRepositoryPermissionsGroupRead(d, m)
}

func resourceRepositoryPermissionsGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("group", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/group`")
		}
	}

	group := d.Get("group").(string)
	groups, err := readRepositoryPermissionsGroups(m, d.Get("project").(string), d.Get("repository").(string), group)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range groups {
		if g.Name == group {
			_ = d.Set("permission", g.Permission)
			break
		}
	}

	return nil
}

func resourceRepositoryPermissionsGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/groups?name=%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Get("group").(string)),
	))

	return err
}
