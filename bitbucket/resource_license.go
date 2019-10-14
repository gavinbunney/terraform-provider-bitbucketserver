package bitbucket

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type License struct {
	License string `json:"license,omitempty"`
}

type LicenseResponse struct {
	License                  string   `json:"license,omitempty"`
	CreationDate             jsonTime `json:"creationDate,omitempty"`
	PurchaseDate             jsonTime `json:"purchaseDate,omitempty"`
	ExpiryDate               jsonTime `json:"expiryDate,omitempty"`
	MaintenanceExpiryDate    jsonTime `json:"maintenanceExpiryDate,omitempty"`
	GracePeriodEndDate       jsonTime `json:"gracePeriodEndDate,omitempty"`
	MaximumNumberOfUsers     int      `json:"maximumNumberOfUsers,omitempty"`
	UnlimitedUsers           bool     `json:"unlimitedNumberOfUsers,omitempty"`
	ServerId                 string   `json:"serverId,omitempty"`
	SupportEntitlementNumber string   `json:"supportEntitlementNumber,omitempty"`
}

func resourceLicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceLicenseCreate,
		Update: resourceLicenseUpdate,
		Read:   resourceLicenseRead,
		Delete: resourceLicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"license": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchase_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"grace_period_end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maximum_users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unlimited_users": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_entitlement_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newLicenseFromResource(d *schema.ResourceData) *License {
	license := &License{
		License: d.Get("license").(string),
	}

	return license
}

func resourceLicenseUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	license := newLicenseFromResource(d)

	bytedata, err := json.Marshal(license)

	if err != nil {
		return err
	}

	_, err = client.Post("/rest/api/1.0/admin/license", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%x", sha256.Sum256([]byte(license.License))))
	return resourceLicenseRead(d, m)
}

func resourceLicenseCreate(d *schema.ResourceData, m interface{}) error {
	return resourceLicenseUpdate(d, m)
}

func resourceLicenseRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get("/rest/api/1.0/admin/license")

	if err != nil {
		return err
	}

	if req.StatusCode == 200 {

		var license LicenseResponse

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &license)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("license", license.License)
		d.Set("creation_date", license.CreationDate.String())
		d.Set("purchase_date", license.PurchaseDate.String())
		d.Set("expiry_date", license.ExpiryDate.String())
		d.Set("maintenance_expiry_date", license.MaintenanceExpiryDate.String())
		d.Set("grace_period_end_date", license.GracePeriodEndDate.String())
		d.Set("maximum_users", license.MaximumNumberOfUsers)
		d.Set("unlimited_users", license.UnlimitedUsers)
		d.Set("server_id", license.ServerId)
		d.Set("support_entitlement_number", license.SupportEntitlementNumber)
	}

	return nil
}

func resourceLicenseDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete("/rest/api/1.0/admin/mail-server")
	return err
}
