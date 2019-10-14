package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedGroupUsersValue struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
}

type GroupUser struct {
	Name         string
	EmailAddress string
	DisplayName  string
	Active       bool
}

type PaginatedGroupUsers struct {
	Values        []PaginatedGroupUsersValue `json:"values,omitempty"`
	Size          int                        `json:"size,omitempty"`
	Limit         int                        `json:"limit,omitempty"`
	IsLastPage    bool                       `json:"isLastPage,omitempty"`
	Start         int                        `json:"start,omitempty"`
	NextPageStart int                        `json:"nextPageStart,omitempty"`
}

func dataSourceGroupUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupUsersRead,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupUsersRead(d *schema.ResourceData, m interface{}) error {
	users, err := readGroupUsers(m, d.Get("group").(string), d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId(d.Get("group").(string))

	var terraformUsers []interface{}
	for _, group := range users {
		g := make(map[string]interface{})
		g["name"] = group.Name
		g["email_address"] = group.EmailAddress
		g["display_name"] = group.DisplayName
		g["active"] = group.Active
		terraformUsers = append(terraformUsers, g)
	}

	_ = d.Set("users", terraformUsers)
	return nil
}

func readGroupUsers(m interface{}, group string, filter string) ([]GroupUser, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := fmt.Sprintf("/rest/api/1.0/admin/groups/more-members?context=%s",
		url.QueryEscape(group),
	)

	if filter != "" {
		resourceURL += "&filter=" + url.QueryEscape(filter)
	}

	var groupUsers PaginatedGroupUsers
	var users []GroupUser

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&groupUsers)
		if err != nil {
			return nil, err
		}

		for _, user := range groupUsers.Values {
			g := GroupUser{
				Name:         user.Name,
				EmailAddress: user.EmailAddress,
				DisplayName:  user.DisplayName,
				Active:       user.Active,
			}
			users = append(users, g)
		}

		if groupUsers.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/users?start=%d",
				group,
				groupUsers.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			groupUsers = PaginatedGroupUsers{}
		} else {
			break
		}
	}

	return users, nil
}
