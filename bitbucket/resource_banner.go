package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"io/ioutil"
)

type Banner struct {
	Message  string `json:"message,omitempty"`
	Audience string `json:"audience,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

func resourceBanner() *schema.Resource {
	return &schema.Resource{
		Create: resourceBannerCreate,
		Update: resourceBannerUpdate,
		Read:   resourceBannerRead,
		Exists: resourceBannerExists,
		Delete: resourceBannerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audience": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALL",
				ValidateFunc: validation.StringInSlice([]string{"ALL", "AUTHENTICATED"}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func newBannerFromResource(d *schema.ResourceData) *Banner {
	banner := &Banner{
		Message:  d.Get("message").(string),
		Audience: d.Get("audience").(string),
		Enabled:  d.Get("enabled").(bool),
	}

	return banner
}

func resourceBannerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	banner := newBannerFromResource(d)

	bytedata, err := json.Marshal(banner)

	if err != nil {
		return err
	}

	_, err = client.Put("/rest/api/1.0/admin/banner", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	d.SetId("banner")
	return resourceBannerRead(d, m)
}

func resourceBannerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceBannerUpdate(d, m)
}

func resourceBannerRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get("/rest/api/1.0/admin/banner")

	if err != nil {
		return err
	}

	var banner Banner

	body, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		return readErr
	}

	decodeErr := json.Unmarshal(body, &banner)
	if decodeErr != nil {
		return decodeErr
	}

	_ = d.Set("message", banner.Message)
	_ = d.Set("audience", banner.Audience)
	_ = d.Set("enabled", banner.Enabled)

	return nil
}

func resourceBannerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient
	repoReq, err := client.Get("/rest/api/1.0/admin/banner")
	if err != nil {
		return false, fmt.Errorf("failed to get banner from bitbucket: %+v", err)
	}

	if repoReq.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceBannerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete("/rest/api/1.0/admin/banner")
	return err
}
