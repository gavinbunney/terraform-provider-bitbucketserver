package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
)

func resourceGlobalPermissionsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalPermissionsGroupCreate,
		Update: resourceGlobalPermissionsGroupUpdate,
		Read:   resourceGlobalPermissionsGroupRead,
		Delete: resourceGlobalPermissionsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group": {
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

func resourceGlobalPermissionsGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/admin/permissions/groups?permission=%s&name=%s",
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("group").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceGlobalPermissionsGroupRead(d, m)
}

func resourceGlobalPermissionsGroupCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceGlobalPermissionsGroupUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(d.Get("group").(string))
	return resourceGlobalPermissionsGroupRead(d, m)
}

func resourceGlobalPermissionsGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		_ = d.Set("group", id)
	}

	group := d.Get("group").(string)
	groups, err := readGlobalPermissionsGroups(m, group)
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

func resourceGlobalPermissionsGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/admin/permissions/groups?name=%s",
		url.QueryEscape(d.Get("group").(string)),
	))

	return err
}
