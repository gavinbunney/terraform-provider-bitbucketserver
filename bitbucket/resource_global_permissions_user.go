package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
)

func resourceGlobalPermissionsUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalPermissionsUserCreate,
		Update: resourceGlobalPermissionsUserUpdate,
		Read:   resourceGlobalPermissionsUserRead,
		Delete: resourceGlobalPermissionsUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LICENSED_USER", "PROJECT_CREATE", "ADMIN", "SYS_ADMIN"}, false),
			},
		},
	}
}

func resourceGlobalPermissionsUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/admin/permissions/users?permission=%s&name=%s",
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("user").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceGlobalPermissionsUserRead(d, m)
}

func resourceGlobalPermissionsUserCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceGlobalPermissionsUserUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(d.Get("user").(string))
	return resourceGlobalPermissionsUserRead(d, m)
}

func resourceGlobalPermissionsUserRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		_ = d.Set("user", id)
	}

	user := d.Get("user").(string)
	users, err := readGlobalPermissionsUsers(m, user)
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

func resourceGlobalPermissionsUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/admin/permissions/users?name=%s",
		url.QueryEscape(d.Get("user").(string)),
	))

	return err
}
