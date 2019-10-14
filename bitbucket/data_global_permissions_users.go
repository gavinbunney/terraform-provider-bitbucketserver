package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

type PaginatedGlobalPermissionsUsersValue struct {
	User struct {
		Name         string `json:"name,omitempty"`
		EmailAddress string `json:"emailAddress,omitempty"`
		DisplayName  string `json:"displayName,omitempty"`
		Active       bool   `json:"active,omitempty"`
	} `json:"user,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type GlobalPermissionsUser struct {
	Name         string
	EmailAddress string
	DisplayName  string
	Active       bool
	Permission   string
}

type PaginatedGlobalPermissionsUsers struct {
	Values        []PaginatedGlobalPermissionsUsersValue `json:"values,omitempty"`
	Size          int                                    `json:"size,omitempty"`
	Limit         int                                    `json:"limit,omitempty"`
	IsLastPage    bool                                   `json:"isLastPage,omitempty"`
	Start         int                                    `json:"start,omitempty"`
	NextPageStart int                                    `json:"nextPageStart,omitempty"`
}

func dataSourceGlobalPermissionsUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlobalPermissionsUsersRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceGlobalPermissionsUsersRead(d *schema.ResourceData, m interface{}) error {
	users, err := readGlobalPermissionsUsers(m, d.Get("filter").(string))
	if err != nil {
		return err
	}

	d.SetId("global-permissions-users")

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

func readGlobalPermissionsUsers(m interface{}, filter string) ([]GlobalPermissionsUser, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := "/rest/api/1.0/admin/permissions/users"

	if filter != "" {
		resourceURL += "?filter=" + url.QueryEscape(filter)
	}

	var globalUsers PaginatedGlobalPermissionsUsers
	var users []GlobalPermissionsUser

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&globalUsers)
		if err != nil {
			return nil, err
		}

		for _, user := range globalUsers.Values {
			g := GlobalPermissionsUser{
				Name:         user.User.Name,
				EmailAddress: user.User.EmailAddress,
				DisplayName:  user.User.DisplayName,
				Active:       user.User.Active,
				Permission:   user.Permission,
			}
			users = append(users, g)
		}

		if globalUsers.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/admin/permissions/users?start=%d",
				globalUsers.NextPageStart,
			)

			if filter != "" {
				resourceURL += "&filter=" + url.QueryEscape(filter)
			}

			globalUsers = PaginatedGlobalPermissionsUsers{}
		} else {
			break
		}
	}

	return users, nil
}
