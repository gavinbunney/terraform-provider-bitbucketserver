package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
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
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	groupName := d.Get("name").(string)
	_, err := client.Post(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s", url.QueryEscape(groupName)), nil)
	if err != nil {
		return err
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

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	groupName := d.Get("name").(string)
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s",
		url.QueryEscape(groupName),
	))

	return err
}
