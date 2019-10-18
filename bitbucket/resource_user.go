package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"io/ioutil"
	"math/rand"
	"net/url"
	"time"
)

type User struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

type UserUpdate struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"email,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Update: resourceUserUpdate,
		Read:   resourceUserRead,
		Exists: resourceUserExists,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(5, 128),
			},
			"initial_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

const passwordCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"0123456789" +
	"@^*_-[]"

func generateUserPassword(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = passwordCharset[seededRand.Intn(len(passwordCharset))]
	}
	return string(b)
}

func newUserFromResource(d *schema.ResourceData) *User {
	user := &User{
		Name:         d.Get("name").(string),
		EmailAddress: d.Get("email_address").(string),
		DisplayName:  d.Get("display_name").(string),
	}

	return user
}

func newUserUpdateFromResource(d *schema.ResourceData) *UserUpdate {
	user := &UserUpdate{
		Name:         d.Get("name").(string),
		EmailAddress: d.Get("email_address").(string),
		DisplayName:  d.Get("display_name").(string),
	}

	return user
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	user := newUserUpdateFromResource(d)

	bytedata, err := json.Marshal(user)

	if err != nil {
		return err
	}

	_, err = client.Put("/rest/api/1.0/admin/users", bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	user := newUserFromResource(d)

	passwordLength := d.Get("password_length").(int)
	initialPassword := generateUserPassword(passwordLength)
	d.Set("initial_password", initialPassword)

	_, err := client.Post(fmt.Sprintf("/rest/api/1.0/admin/users?name=%s&password=%s&displayName=%s&emailAddress=%s",
		url.QueryEscape(user.Name),
		url.QueryEscape(initialPassword),
		url.QueryEscape(user.DisplayName),
		url.QueryEscape(user.EmailAddress),
	), nil)

	if err != nil {
		return err
	}

	d.SetId(user.Name)

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		d.Set("name", id)
	}

	name := d.Get("name").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/api/1.0/users/%s",
		url.QueryEscape(name),
	))

	if err != nil {
		return err
	}

	if req.StatusCode == 200 {

		var user User

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &user)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("name", user.Name)
		d.Set("email_address", user.EmailAddress)
		d.Set("display_name", user.DisplayName)
	}

	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	var name = ""
	id := d.Id()
	if id != "" {
		name = id
	} else {
		name = d.Get("name").(string)
	}

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/api/1.0/users/%s",
		url.QueryEscape(name),
	))

	if err != nil {
		return false, fmt.Errorf("failed to get user %s from bitbucket: %+v", name, err)
	}

	if req.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/admin/users?name=%s",
		url.QueryEscape(name),
	))

	return err
}
