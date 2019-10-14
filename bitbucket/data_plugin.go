package bitbucket

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_by_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"marketplace_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"applied_license": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"valid": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"evaluation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nearly_expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintenance_expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintenance_expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raw_license": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"renewable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"data_center": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"subscription": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renewal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"upgradable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crossgradeable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"purchase_past_server_cutoff_date": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"support_entitlement_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePluginRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))
	return resourcePluginRead(d, m)
}
