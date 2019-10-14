package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedProjectPermissionsGroupsValue struct {
	Group struct {
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type ProjectPermissionsGroup struct {
	Name       string
	Permission string
}

type PaginatedProjectPermissionsGroups struct {
	Values        []PaginatedProjectPermissionsGroupsValue `json:"values,omitempty"`
	Size          int                                      `json:"size,omitempty"`
	Limit         int                                      `json:"limit,omitempty"`
	IsLastPage    bool                                     `json:"isLastPage,omitempty"`
	Start         int                                      `json:"start,omitempty"`
	NextPageStart int                                      `json:"nextPageStart,omitempty"`
}

func dataSourceProjectPermissionsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectPermissionsGroupsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
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
	groups, err := readProjectPermissionsGroups(m, d.Get("project").(string), d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId(d.Get("project").(string))

	var terraformGroups []interface{}
	for _, group := range groups {
		g := make(map[string]interface{})
		g["name"] = group.Name
		g["permission"] = group.Permission
		terraformGroups = append(terraformGroups, g)
	}

	_ = d.Set("groups", terraformGroups)
	return nil
}

func readProjectPermissionsGroups(m interface{}, project string, filter string) ([]ProjectPermissionsGroup, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups",
		project,
	)

	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var projectGroups PaginatedProjectPermissionsGroups
	var groups []ProjectPermissionsGroup

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&projectGroups)
		if err != nil {
			return nil, err
		}

		for _, group := range projectGroups.Values {
			g := ProjectPermissionsGroup{
				Name:       group.Group.Name,
				Permission: group.Permission,
			}
			groups = append(groups, g)
		}

		if projectGroups.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/groups?start=%d",
				project,
				projectGroups.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			projectGroups = PaginatedProjectPermissionsGroups{}
		} else {
			break
		}
	}

	return groups, nil
}
