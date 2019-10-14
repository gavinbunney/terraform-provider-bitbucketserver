package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

type UserGroup struct {
	User  string
	Group string
}

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Delete: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func newUserGroupFromResource(d *schema.ResourceData) *UserGroup {
	userGroup := &UserGroup{
		User:  d.Get("user").(string),
		Group: d.Get("group").(string),
	}

	return userGroup
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	type UserGroupRequest struct {
		User   string   `json:"user,omitempty"`
		Groups []string `json:"groups,omitempty"`
	}

	request := &UserGroupRequest{
		User: d.Get("user").(string),
		Groups: []string{
			d.Get("group").(string),
		},
	}

	bytedata, err := json.Marshal(request)

	if err != nil {
		return err
	}

	_, err = client.Post("/rest/api/1.0/admin/users/add-groups", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", request.User, request.Groups[0]))

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 2 {
			_ = d.Set("user", parts[0])
			_ = d.Set("group", parts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `user/group`")
		}
	}

	userGroup := newUserGroupFromResource(d)

	groupUsers, err := readGroupUsers(m, userGroup.Group, userGroup.User)
	if err != nil {
		return err
	}

	// API only filters but we need to find an exact match
	for _, g := range groupUsers {
		if g.Name == userGroup.User {
			return nil
		}
	}

	return fmt.Errorf("unable to find a matching user %s in group %s", userGroup.User, userGroup.Group)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {

	userGroup := newUserGroupFromResource(d)

	client := m.(*BitbucketServerProvider).BitbucketClient

	type RemoveRequest struct {
		User  string `json:"context,omitempty"`
		Group string `json:"itemName,omitempty"`
	}

	removeRequest := &RemoveRequest{
		User:  userGroup.User,
		Group: userGroup.Group,
	}

	bytedata, err := json.Marshal(removeRequest)
	if err != nil {
		return err
	}

	_, err = client.Post("/rest/api/1.0/admin/users/remove-group", bytes.NewBuffer(bytedata))

	return err
}
