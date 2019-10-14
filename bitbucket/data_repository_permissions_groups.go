package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedRepositoryPermissionsGroupsValue struct {
	Group struct {
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type RepositoryPermissionsGroup struct {
	Name       string
	Permission string
}

type PaginatedRepositoryPermissionsGroups struct {
	Values        []PaginatedRepositoryPermissionsGroupsValue `json:"values,omitempty"`
	Size          int                                         `json:"size,omitempty"`
	Limit         int                                         `json:"limit,omitempty"`
	IsLastPage    bool                                        `json:"isLastPage,omitempty"`
	Start         int                                         `json:"start,omitempty"`
	NextPageStart int                                         `json:"nextPageStart,omitempty"`
}

func dataSourceRepositoryPermissionsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRepositoryPermissionsGroupsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
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

func dataSourceRepositoryPermissionsGroupsRead(d *schema.ResourceData, m interface{}) error {
	groups, err := readRepositoryPermissionsGroups(m, d.Get("project").(string), d.Get("repository").(string), d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("project").(string), d.Get("repository").(string)))

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

func readRepositoryPermissionsGroups(m interface{}, project string, repository string, filter string) ([]RepositoryPermissionsGroup, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/groups",
		url.QueryEscape(project),
		url.QueryEscape(repository),
	)

	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var projectGroups PaginatedRepositoryPermissionsGroups
	var groups []RepositoryPermissionsGroup

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
			g := RepositoryPermissionsGroup{
				Name:       group.Group.Name,
				Permission: group.Permission,
			}
			groups = append(groups, g)
		}

		if projectGroups.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/permissions/groups?start=%d",
				url.QueryEscape(project),
				url.QueryEscape(repository),
				projectGroups.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			projectGroups = PaginatedRepositoryPermissionsGroups{}
		} else {
			break
		}
	}

	return groups, nil
}
