package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

type ProjectPermissionsGroup struct {
	Group struct {
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type PaginatedProjectPermissionsGroups struct {
	Values        []ProjectPermissionsGroup `json:"values,omitempty"`
	Size          int                       `json:"size,omitempty"`
	Limit         int                       `json:"limit,omitempty"`
	IsLastPage    bool                      `json:"isLastPage,omitempty"`
	Start         int                       `json:"start,omitempty"`
	NextPageStart int                       `json:"nextPageStart,omitempty"`
}

func dataSourceProjectPermissionsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectPermissionsGroupsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectPermissionsGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	resourceURL := fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups",
		d.Get("project").(string),
	)

	var projectGroups PaginatedProjectPermissionsGroups
	var terraformGroups []interface{}

	for {
		reviewersResponse, err := client.Get(resourceURL)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(reviewersResponse.Body)
		err = decoder.Decode(&projectGroups)
		if err != nil {
			return err
		}

		for _, group := range projectGroups.Values {
			g := make(map[string]interface{})
			g["name"] = group.Group.Name
			g["permission"] = group.Permission
			terraformGroups = append(terraformGroups, g)
		}

		if projectGroups.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups?start=%d",
				d.Get("project").(string),
				projectGroups.NextPageStart,
			)
			projectGroups = PaginatedProjectPermissionsGroups{}
		} else {
			break
		}
	}

	d.SetId(d.Get("project").(string))
	d.Set("groups", terraformGroups)
	return nil
}
