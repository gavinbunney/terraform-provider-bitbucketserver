package bitbucket

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forkable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	slug := d.Get("slug").(string)
	project := d.Get("project").(string)
	d.SetId(fmt.Sprintf("%s/%s", project, slug))
	return resourceRepositoryRead(d, m)
}
