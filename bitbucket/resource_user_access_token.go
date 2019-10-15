package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type AccessTokenRequest struct {
	Name        string        `json:"name,omitempty"`
	Permissions []interface{} `json:"permissions,omitempty"`
}

type AccessTokenResponse struct {
	Id                string   `json:"id,omitempty"`
	CreatedDate       jsonTime `json:"createdDate,omitempty"`
	LastAuthenticated jsonTime `json:"lastAuthenticated,omitempty"`
	Name              string   `json:"name,omitempty"`
	Permissions       []string `json:"permissions,omitempty"`
	Token             string   `json:"token,omitempty"`
}

func resourceUserAccessToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAccessTokenCreate,
		Update: resourceUserAccessTokenUpdate,
		Read:   resourceUserAccessTokenRead,
		Exists: resourceUserAccessTokenExists,
		Delete: resourceUserAccessTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_authenticated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceUserAccessTokenCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	accessTokenRequest := &AccessTokenRequest{
		Name:        d.Get("name").(string),
		Permissions: d.Get("permissions").([]interface{}),
	}

	byteData, err := json.Marshal(accessTokenRequest)
	if err != nil {
		return err
	}

	res, err := client.Put(fmt.Sprintf("/rest/access-tokens/1.0/users/%s",
		d.Get("user").(string),
	), bytes.NewBuffer(byteData))

	if err != nil {
		return err
	}

	var accessTokenResponse AccessTokenResponse

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	decodeErr := json.Unmarshal(body, &accessTokenResponse)
	if decodeErr != nil {
		return decodeErr
	}

	d.SetId(accessTokenResponse.Id)
	_ = d.Set("access_token", accessTokenResponse.Token)

	return resourceUserAccessTokenRead(d, m)
}

func resourceUserAccessTokenUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	accessTokenRequest := &AccessTokenRequest{
		Name:        d.Get("name").(string),
		Permissions: d.Get("permissions").([]interface{}),
	}

	byteData, err := json.Marshal(accessTokenRequest)
	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("/rest/access-tokens/1.0/users/%s/%s",
		d.Get("user").(string),
		d.Id(),
	), bytes.NewBuffer(byteData))

	if err != nil {
		return err
	}

	return resourceUserAccessTokenRead(d, m)
}

func resourceUserAccessTokenRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*BitbucketServerProvider).BitbucketClient
	res, err := client.Get(fmt.Sprintf("/rest/access-tokens/1.0/users/%s/%s",
		d.Get("user").(string),
		d.Id(),
	))

	if err != nil {
		return err
	}

	var accessTokenResponse AccessTokenResponse

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	decodeErr := json.Unmarshal(body, &accessTokenResponse)
	if decodeErr != nil {
		return decodeErr
	}

	_ = d.Set("name", accessTokenResponse.Name)
	_ = d.Set("created_date", accessTokenResponse.CreatedDate.String())
	_ = d.Set("last_authenticated", accessTokenResponse.LastAuthenticated.String())

	return nil
}

func resourceUserAccessTokenExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/access-tokens/1.0/users/%s/%s",
		d.Get("user").(string),
		d.Id(),
	))

	if err != nil {
		return false, fmt.Errorf("failed to get access token %s for user %s from bitbucket: %+v", d.Id(), d.Get("user").(string), err)
	}

	if req.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceUserAccessTokenDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/access-tokens/1.0/users/%s/%s",
		d.Get("user").(string),
		d.Id(),
	))

	return err
}
