package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/url"
	"sort"
)

type PaginatedProjectHooksValue struct {
	Details struct {
		Key         string   `json:"key,omitempty"`
		Name        string   `json:"name,omitempty"`
		Type        string   `json:"type,omitempty"`
		Description string   `json:"description,omitempty"`
		Version     string   `json:"version,omitempty"`
		ScopeTypes  []string `json:"scopeTypes,omitempty"`
	} `json:"details,omitempty"`
	Enabled    bool `json:"enabled,omitempty"`
	Configured bool `json:"configured,omitempty"`
	Scope      struct {
		Type       string `json:"type,omitempty"`
		ResourceId int    `json:"resourceId,omitempty"`
	} `json:"scope,omitempty"`
}

type ProjectHook struct {
	Key             string
	Name            string
	Type            string
	Description     string
	Version         string
	ScopeTypes      []string
	Enabled         bool
	Configured      bool
	ScopeType       string
	ScopeResourceId int
}

type PaginatedProjectHooks struct {
	Values        []PaginatedProjectHooksValue `json:"values,omitempty"`
	Size          int                          `json:"size,omitempty"`
	Limit         int                          `json:"limit,omitempty"`
	IsLastPage    bool                         `json:"isLastPage,omitempty"`
	Start         int                          `json:"start,omitempty"`
	NextPageStart int                          `json:"nextPageStart,omitempty"`
}

func dataSourceProjectHooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectHooksRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRE_RECEIVE", "POST_RECEIVE"}, false),
			},
			"hooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_types": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"configured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"scope_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_resource_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectHooksRead(d *schema.ResourceData, m interface{}) error {
	hooks, err := readProjectHooks(m, d.Get("project").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	d.SetId(d.Get("project").(string))

	var terraformHooks []interface{}
	for _, hook := range hooks {
		h := make(map[string]interface{})
		h["key"] = hook.Key
		h["name"] = hook.Name
		h["type"] = hook.Type
		h["description"] = hook.Description
		h["version"] = hook.Version
		h["scope_types"] = hook.ScopeTypes
		h["enabled"] = hook.Enabled
		h["configured"] = hook.Configured
		h["scope_type"] = hook.ScopeType
		h["scope_resource_id"] = hook.ScopeResourceId
		terraformHooks = append(terraformHooks, h)
	}

	_ = d.Set("hooks", terraformHooks)
	return nil
}

func readProjectHooks(m interface{}, project string, typeFilter string) ([]ProjectHook, error) {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resourceURL := fmt.Sprintf("/rest/api/1.0/projects/%s/settings/hooks",
		project,
	)

	if typeFilter != "" {
		resourceURL += "?type=" + url.QueryEscape(typeFilter)
	}

	var projectHooks PaginatedProjectHooks
	var hooks []ProjectHook

	for {
		resp, err := client.Get(resourceURL)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&projectHooks)
		if err != nil {
			return nil, err
		}

		for _, hook := range projectHooks.Values {
			sort.Strings(hook.Details.ScopeTypes)
			h := ProjectHook{
				Key:             hook.Details.Key,
				Name:            hook.Details.Name,
				Type:            hook.Details.Type,
				Description:     hook.Details.Description,
				Version:         hook.Details.Version,
				ScopeTypes:      hook.Details.ScopeTypes,
				Enabled:         hook.Enabled,
				Configured:      hook.Configured,
				ScopeType:       hook.Scope.Type,
				ScopeResourceId: hook.Scope.ResourceId,
			}
			hooks = append(hooks, h)
		}

		if projectHooks.IsLastPage == false {
			resourceURL = fmt.Sprintf("/rest/api/1.0/projects/%s/settings/hooks?start=%d",
				project,
				projectHooks.NextPageStart,
			)

			if typeFilter != "" {
				resourceURL += "&type=" + url.QueryEscape(typeFilter)
			}

			projectHooks = PaginatedProjectHooks{}
		} else {
			break
		}
	}

	return hooks, nil
}
