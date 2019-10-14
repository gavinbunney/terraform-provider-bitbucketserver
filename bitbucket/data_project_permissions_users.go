package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedProjectPermissionsUsersValue struct {
	User struct {
		Name         string `json:"name,omitempty"`
		EmailAddress string `json:"emailAddress,omitempty"`
		DisplayName  string `json:"displayName,omitempty"`
		Active       bool   `json:"active,omitempty"`
	} `json:"user,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type ProjectPermissionsUser struct {
	Name         string
	EmailAddress string
	DisplayName  string
	Active       bool
	Permission   string
}

type PaginatedProjectPermissionsUsers struct {
	Values        []PaginatedProjectPermissionsUsersValue `json:"values,omitempty"`
	Size          int                                     `json:"size,omitempty"`
	Limit         int                                     `json:"limit,omitempty"`
	IsLastPage    bool                                    `json:"isLastPage,omitempty"`
	Start         int                                     `json:"start,omitempty"`
	NextPageStart int                                     `json:"nextPageStart,omitempty"`
}

func dataSourceProjectPermissionsUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectPermissionsUsersRead,

		Schema: map[string]*schema.Schema{
			"project": {
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

func dataSourceProjectPermissionsUsersRead(d *schema.ResourceData, m interface{}) error {
	users, err := readProjectPermissionsUsers(m, d.Get("project").(string), d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId(d.Get("project").(string))

	var terraformUsers []interface{}
	for _, group := range users {
		g := make(map[string]interface{})
		g["name"] = group.Name
		g["email_address"] = group.EmailAddress
		g["display_name"] = group.DisplayName
		g["active"] = group.Active
		g["permission"] = group.Permission
		terraformUsers = append(terraformUsers, g)
	}

	_ = d.Set("users", terraformUsers)
	return nil
}

func readProjectPermissionsUsers(m interface{}, project string, filter string) ([]ProjectPermissionsUser, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/users",
		project,
	)

	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var projectUsers PaginatedProjectPermissionsUsers
	var users []ProjectPermissionsUser

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&projectUsers)
		if err != nil {
			return nil, err
		}

		for _, user := range projectUsers.Values {
			g := ProjectPermissionsUser{
				Name:         user.User.Name,
				EmailAddress: user.User.EmailAddress,
				DisplayName:  user.User.DisplayName,
				Active:       user.User.Active,
				Permission:   user.Permission,
			}
			users = append(users, g)
		}

		if projectUsers.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/permissions/users?start=%d",
				project,
				projectUsers.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			projectUsers = PaginatedProjectPermissionsUsers{}
		} else {
			break
		}
	}

	return users, nil
}
