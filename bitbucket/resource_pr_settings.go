package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// The struct represents this JSON payload:
// https://docs.atlassian.com/bitbucket-server/rest/7.17.0/bitbucket-rest.html#idp375
type PrSettings struct {
	RequiredApprovers        int         `json:"requiredApprovers"`
	RequiredSuccessfulBuilds int         `json:"requiredSuccessfulBuilds"`
	RequiredAllApprovers     bool        `json:"requiredAllApprovers,omitempty"`
	RequiredAllTasksComplete bool        `json:"requiredAllTasksComplete,omitempty"`
	NoNeedsWork              bool        `json:"needsWork"`
	MergeConfig              MergeConfig `json:"mergeConfig,omitempty"`
}

type MergeConfig struct {
	DefaultStrategy   MergeStrategy   `json:"defaultStrategy,omitempty"`
	EnabledStrategies []MergeStrategy `json:"strategies,omitempty"`
	CommitSummaries   int             `json:"commitSummaries"`
	Type              string          `json:"type,omitempty"`
}

type MergeStrategy struct {
	Id          string `json:"id,omitempty"`
	Flag        string `json:"flag,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
}

type DeleteMergeConfig struct {
	MergeConfig struct {
	} `json:"mergeConfig"`
}

func resourcePrSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrSettingsCreate,
		Read:   resourcePrSettingsRead,
		Update: resourcePrSettingsCreate,
		Delete: resourcePrSettingsDelete,
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
			"required_approvers": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"required_successful_builds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"required_all_approvers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"required_all_tasks_complete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"no_needs_work_status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"merge_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_strategy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"no-ff", "ff", "ff-only", "rebase-no-ff", "rebase-ff-only", "squash", "squash-ff-only"}, false),
						},
						"enabled_strategies": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"commit_summaries": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  20,
						},
					},
				},
			},
		},
	}
}

func resourcePrSettingsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	settings := newPrSettingsFromResource(d)

	bytedata, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/pull-requests",
		d.Get("project").(string),
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		fmt.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%s|%s", d.Get("project").(string), d.Get("repository").(string)))

	return resourcePrSettingsRead(d, m)
}

func newPrSettingsFromResource(d *schema.ResourceData) *PrSettings {
	settings := &PrSettings{
		RequiredApprovers:        d.Get("required_approvers").(int),
		RequiredSuccessfulBuilds: d.Get("required_successful_builds").(int),
		RequiredAllApprovers:     d.Get("required_all_approvers").(bool),
		RequiredAllTasksComplete: d.Get("required_all_tasks_complete").(bool),
		NoNeedsWork:              d.Get("no_needs_work_status").(bool),
		MergeConfig:              expandMergeConfig(d.Get("merge_config").([]interface{})),
	}

	return settings
}

func resourcePrSettingsRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	if id != "" {
		idparts := strings.Split(id, "|")
		if len(idparts) == 2 {
			_ = d.Set("project", idparts[0])
			_ = d.Set("repository", idparts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project|repository`")
		}
	}

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/pull-requests",
		d.Get("project").(string),
		d.Get("repository").(string),
	))

	if err != nil {
		return err
	}

	if req.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Workzone Reviewers object (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var settings PrSettings

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &settings)
	if decodeerr != nil {
		return decodeerr
	}

	d.Set("required_approvers", settings.RequiredApprovers)
	d.Set("required_successful_builds", settings.RequiredSuccessfulBuilds)
	d.Set("required_all_approvers", settings.RequiredAllApprovers)
	d.Set("required_all_tasks_complete", settings.RequiredAllTasksComplete)
	d.Set("no_needs_work_status", settings.NoNeedsWork)
	d.Set("merge_config", collapseMergeConfig(settings.MergeConfig))

	return nil
}

func resourcePrSettingsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)

	settings := &DeleteMergeConfig{}

	bytedata, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/settings/pull-requests",
		url.QueryEscape(project),
		url.QueryEscape(repository),
	), bytes.NewBuffer(bytedata))

	return err
}

func expandMergeConfig(l []interface{}) MergeConfig {
	mergeConfigMap := l[0].(map[string]interface{})
	mergeConfig := MergeConfig{
		DefaultStrategy: MergeStrategy{
			Id: mergeConfigMap["default_strategy"].(string),
		},
		CommitSummaries: mergeConfigMap["commit_summaries"].(int),
		Type:            "repository",
	}
	for _, item := range mergeConfigMap["enabled_strategies"].([]interface{}) {
		strategy := MergeStrategy{
			Id:      item.(string),
			Enabled: true,
		}
		mergeConfig.EnabledStrategies = append(mergeConfig.EnabledStrategies, strategy)
	}
	return mergeConfig
}

func collapseMergeConfig(rp MergeConfig) []interface{} {

	m := map[string]interface{}{
		"default_strategy":   rp.DefaultStrategy.Id,
		"commit_summaries":   rp.CommitSummaries,
		"enabled_strategies": rp.EnabledStrategies,
	}

	return []interface{}{m}
}
