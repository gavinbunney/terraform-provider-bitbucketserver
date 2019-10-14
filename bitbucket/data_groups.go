package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedGroupsValue struct {
	Name string `json:"name,omitempty"`
}

type PaginatedGroups struct {
	Values        []PaginatedGroupsValue `json:"values,omitempty"`
	Size          int                    `json:"size,omitempty"`
	Limit         int                    `json:"limit,omitempty"`
	IsLastPage    bool                   `json:"isLastPage,omitempty"`
	Start         int                    `json:"start,omitempty"`
	NextPageStart int                    `json:"nextPageStart,omitempty"`
}

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupsRead,

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
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	groups, err := readGroups(m, d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId("groups")

	var terraformGroups []interface{}
	for _, group := range groups {
		g := make(map[string]interface{})
		g["name"] = group
		terraformGroups = append(terraformGroups, g)
	}

	_ = d.Set("groups", terraformGroups)
	return nil
}

func readGroups(m interface{}, filter string) ([]string, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := "/rest/api/1.0/admin/groups"
	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var paginatedGroups PaginatedGroups
	var groups []string

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&paginatedGroups)
		if err != nil {
			return nil, err
		}

		for _, group := range paginatedGroups.Values {
			groups = append(groups, group.Name)
		}

		if paginatedGroups.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/admin/groups?start=%d",
				paginatedGroups.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			paginatedGroups = PaginatedGroups{}
		} else {
			break
		}
	}

	return groups, nil
}
