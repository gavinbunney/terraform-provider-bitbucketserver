package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
	"strings"
)

func resourceProjectPermissionsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionsGroupCreate,
		Update: resourceProjectPermissionsGroupUpdate,
		Read:   resourceProjectPermissionsGroupRead,
		Delete: resourceProjectPermissionsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
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
				ValidateFunc: validation.StringInSlice([]string{"PROJECT_READ", "PROJECT_WRITE", "PROJECT_ADMIN"}, false),
			},
		},
	}
}

func resourceProjectPermissionsGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups?permission=%s&name=%s",
		d.Get("project").(string),
		url.QueryEscape(d.Get("permission").(string)),
		url.QueryEscape(d.Get("group").(string)),
	), nil)

	if err != nil {
		return err
	}

	return resourceProjectPermissionsGroupRead(d, m)
}

func resourceProjectPermissionsGroupCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceProjectPermissionsGroupUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("project").(string), d.Get("group").(string)))
	return resourceProjectPermissionsGroupRead(d, m)
}

func resourceProjectPermissionsGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 2 {
			d.Set("project", parts[0])
			d.Set("group", parts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/group`")
		}
	}

	group := d.Get("group").(string)
	groups, err := readProjectPermissionsGroups(m, d.Get("project").(string), group)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range groups {
		if g.Name == group {
			d.Set("permission", g.Permission)
			break
		}
	}

	return nil
}

func resourceProjectPermissionsGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups?name=%s",
		d.Get("project").(string),
		url.QueryEscape(d.Get("group").(string)),
	))

	return err
}
