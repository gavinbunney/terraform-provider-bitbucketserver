package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedGlobalPermissionsGroupsValue struct {
	Group struct {
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type GlobalPermissionsGroup struct {
	Name       string
	Permission string
}

type PaginatedGlobalPermissionsGroups struct {
	Values        []PaginatedGlobalPermissionsGroupsValue `json:"values,omitempty"`
	Size          int                                     `json:"size,omitempty"`
	Limit         int                                     `json:"limit,omitempty"`
	IsLastPage    bool                                    `json:"isLastPage,omitempty"`
	Start         int                                     `json:"start,omitempty"`
	NextPageStart int                                     `json:"nextPageStart,omitempty"`
}

func dataSourceGlobalPermissionsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlobalPermissionsGroupsRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceGlobalPermissionsGroupsRead(d *schema.ResourceData, m interface{}) error {
	groups, err := readGlobalPermissionsGroups(m, d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId("global-permissions-groups")

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

func readGlobalPermissionsGroups(m interface{}, filter string) ([]GlobalPermissionsGroup, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := "/rest/api/1.0/admin/permissions/groups"

	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var groupGroups PaginatedGlobalPermissionsGroups
	var groups []GlobalPermissionsGroup

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&groupGroups)
		if err != nil {
			return nil, err
		}

		for _, group := range groupGroups.Values {
			g := GlobalPermissionsGroup{
				Name:       group.Group.Name,
				Permission: group.Permission,
			}
			groups = append(groups, g)
		}

		if groupGroups.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/admin/permissions/groups?start=%d",
				groupGroups.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			groupGroups = PaginatedGlobalPermissionsGroups{}
		} else {
			break
		}
	}

	return groups, nil
}
