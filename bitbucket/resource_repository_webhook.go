package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"strings"
)

type WebhookConfiguration struct {
	Secret string `json:"secret,omitempty"`
}

type Webhook struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	CreatedDate   jsonTime             `json:"createdDate,omitempty"`
	UpdatedDate   jsonTime             `json:"updatedDate,omitempty"`
	URL           string               `json:"url,omitempty"`
	Active        bool                 `json:"active,omitempty"`
	Events        []interface{}        `json:"events"`
	Configuration WebhookConfiguration `json:"configuration"`
}

type WebhookListResponse struct {
	Size       int       `json:"size,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	Start      int       `json:"start,omitempty"`
	IsLastPage bool      `json:"isLastPage,omitempty"`
	Values     []Webhook `json:"values"`
}

func resourceRepositoryWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryWebhookCreate,
		Update: resourceRepositoryWebhookUpdate,
		Read:   resourceRepositoryWebhookRead,
		Delete: resourceRepositoryWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"webhook_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"events": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"webhook_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRepositoryWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	id := d.Get("webhook_id").(int)
	webhook := newWebhookFromResource(d)

	request, err := json.Marshal(webhook)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%d",
		project,
		repository,
		id,
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	return resourceRepositoryWebhookRead(d, m)
}

func resourceRepositoryWebhookCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	webhook := newWebhookFromResource(d)

	request, err := json.Marshal(webhook)

	res, err := client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks",
		project,
		repository,
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	var webhookResponse Webhook

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &webhookResponse)

	if err != nil {
		return err
	}

	_ = d.Set("webhook_id", webhookResponse.ID)

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("name").(string)))
	return resourceRepositoryWebhookRead(d, m)
}

func resourceRepositoryWebhookRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("name", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/name`")
		}
	}

	webhookId := d.Get("webhook_id").(int)

	var err error

	if webhookId != 0 {
		err = getRepositoryWebhookFromId(d, m)
	} else {
		err = getRepositoryWebhookFromList(d, m)
	}

	if err != nil {
		return err
	}

	return nil
}

func resourceRepositoryWebhookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%d",
		d.Get("project").(string),
		d.Get("repository").(string),
		d.Get("webhook_id").(int)))

	return err
}

func newWebhookFromResource(d *schema.ResourceData) (Hook *Webhook) {
	configuration := &WebhookConfiguration{
		Secret: d.Get("secret").(string),
	}

	webhook := &Webhook{
		Name:          d.Get("name").(string),
		URL:           d.Get("webhook_url").(string),
		Active:        d.Get("active").(bool),
		Events:        d.Get("events").([]interface{}),
		Configuration: *configuration,
	}

	return webhook
}

func getRepositoryWebhookFromId(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	id := d.Get("webhook_id").(int)

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks/%d",
		project,
		repository,
		id,
	))

	if err != nil {
		return err
	}

	var webhook Webhook

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&webhook)

	if err != nil {
		return err
	}

	_ = d.Set("webhook_id", webhook.ID)
	_ = d.Set("webhook_url", webhook.URL)
	_ = d.Set("active", webhook.Active)
	_ = d.Set("events", webhook.Events)
	_ = d.Set("secret", webhook.Configuration.Secret)

	return nil
}

func getRepositoryWebhookFromList(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	name := d.Get("name").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks",
		project,
		repository,
	))

	if err != nil {
		return err
	}

	var webhookListResponse WebhookListResponse

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&webhookListResponse)

	if err != nil {
		return err
	}

	for _, webhook := range webhookListResponse.Values {
		if webhook.Name == name {
			_ = d.Set("webhook_id", webhook.ID)
			_ = d.Set("webhook_url", webhook.URL)
			_ = d.Set("active", webhook.Active)
			_ = d.Set("events", webhook.Events)
			_ = d.Set("secret", webhook.Configuration.Secret)
			return nil
		}
	}

	return fmt.Errorf("incorrect ID format, should match `project/repository/name`")
}
