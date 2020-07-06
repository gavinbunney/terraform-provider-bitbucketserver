package bitbucket

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"import_if_exists": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	groupName := d.Get("name").(string)
	importIfExists := d.Get("import_if_exists").(bool)
	var newResource = true
	_, err := client.Post(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s", url.QueryEscape(groupName)), nil)
	if err != nil {
		if importIfExists && strings.Contains(err.Error(), "API Error: 409") {
			newResource = false
		} else {
			return err
		}
	}

	if newResource {
		d.MarkNewResource()
	}
	d.SetId(groupName)

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		_ = d.Set("name", id)
	}

	groupName := d.Get("name").(string)

	groupMatches, err := readGroups(m, groupName)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range groupMatches {
		if g == groupName {
			return nil
		}
	}

	return fmt.Errorf("unable to find a matching group %s", groupName)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	// The only attribute in the schema that does not have "ForceNew: true" is "import_if_exists",
	// so we are not actually updating any groups in Bitbucket, we just need to read and return.
	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	groupName := d.Get("name").(string)
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s",
		url.QueryEscape(groupName),
	))

	return err
}
