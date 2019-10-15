package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
	"strings"
)

func resourceProjectPermissionsUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionsUserCreate,
		Update: resourceProjectPermissionsUserUpdate,
		Read:   resourceProjectPermissionsUserRead,
		Delete: resourceProjectPermissionsUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
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
				ValidateFunc: validation.StringInSlice([]string{"PROJECT_READ", "PROJECT_WRITE", "PROJECT_ADMIN"}, false),
			},
		},
	}
}

func resourceProjectPermissionsUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/users?permission=%s&name=%s",
		d.Get("project").(string),
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("user").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceProjectPermissionsUserRead(d, m)
}

func resourceProjectPermissionsUserCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceProjectPermissionsUserUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("project").(string), d.Get("user").(string)))
	return resourceProjectPermissionsUserRead(d, m)
}

func resourceProjectPermissionsUserRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 2 {
			d.Set("project", parts[0])
			d.Set("user", parts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/user`")
		}
	}

	user := d.Get("user").(string)
	users, err := readProjectPermissionsUsers(m, d.Get("project").(string), user)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range users {
		if g.Name == user {
			d.Set("permission", g.Permission)
			break
		}
	}

	return nil
}

func resourceProjectPermissionsUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/users?name=%s",
		d.Get("project").(string),
		url.QueryEscape(d.Get("user").(string)),
	))

	return err
}
