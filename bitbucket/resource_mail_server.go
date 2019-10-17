package bitbucket

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type MailConfiguration struct {
	Hostname        string `json:"hostname,omitempty"`
	Port            int    `json:"port,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	UseStartTLS     bool   `json:"use-start-tls,omitempty"`
	RequireStartTLS bool   `json:"require-start-tls,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	SenderAddress   string `json:"sender-address,omitempty"`
}

func resourceMailServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceMailServerCreate,
		Update: resourceMailServerUpdate,
		Read:   resourceMailServerRead,
		Delete: resourceMailServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"use_start_tls": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"require_start_tls": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"sender_address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func newMailConfigurationFromResource(d *schema.ResourceData) *MailConfiguration {
	mailConfiguration := &MailConfiguration{
		Hostname:        d.Get("hostname").(string),
		Port:            d.Get("port").(int),
		Protocol:        d.Get("protocol").(string),
		UseStartTLS:     d.Get("use_start_tls").(bool),
		RequireStartTLS: d.Get("require_start_tls").(bool),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
		SenderAddress:   d.Get("sender_address").(string),
	}

	return mailConfiguration
}

func resourceMailServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	mailConfiguration := newMailConfigurationFromResource(d)

	bytedata, err := json.Marshal(mailConfiguration)

	if err != nil {
		return err
	}

	_, err = client.Put("/rest/api/1.0/admin/mail-server", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	d.SetId(mailConfiguration.Hostname)

	return resourceMailServerRead(d, m)
}

func resourceMailServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceMailServerUpdate(d, m)
}

func resourceMailServerRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get("/rest/api/1.0/admin/mail-server")

	if err != nil {
		return err
	}

	if req.StatusCode == 200 {

		var mailConfiguration MailConfiguration

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &mailConfiguration)
		if decodeerr != nil {
			return decodeerr
		}

		_ = d.Set("hostname", mailConfiguration.Hostname)
		_ = d.Set("port", mailConfiguration.Port)
		_ = d.Set("protocol", mailConfiguration.Protocol)
		_ = d.Set("use_start_tls", mailConfiguration.UseStartTLS)
		_ = d.Set("require_start_tls", mailConfiguration.RequireStartTLS)
		_ = d.Set("username", mailConfiguration.Username)
		_ = d.Set("sender_address", mailConfiguration.SenderAddress)
	}

	return nil
}

func resourceMailServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete("/rest/api/1.0/admin/mail-server")
	return err
}
