package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type BranchPermissionPayload struct {
	Type       string        `json:"type,omitempty"`
	Matcher    MatcherStruct `json:"matcher,omitempty"`
	Users      []string      `json:"users,omitempty"`
	Groups     []string      `json:"groups,omitempty"`
	AccessKeys []string      `json:"accessKeys,omitempty"`
}

type MatcherStruct struct {
	Id        string            `json:"id,omitempty"`
	DisplayId string            `json:"displayId,omitempty"`
	Type      MatcherStructType `json:"type,omitempty"`
	Active    bool              `json:"active,omitempty"`
}

type MatcherStructType struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BranchPermissionResponse struct {
	Id    int `json:"id"`
	Scope struct {
		ResourceID int    `json:"resourceId"`
		Type       string `json:"type"`
	} `json:"scope"`
	Type  string `json:"type"`
	Users []struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		ID           int    `json:"id"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
	} `json:"users"`
	Groups     []string `json:"groups"`
	AccessKeys []struct {
		Key struct {
			ID    int    `json:"id"`
			Text  string `json:"text"`
			Label string `json:"label"`
		} `json:"key"`
	} `json:"accessKeys"`
}

type AllRepositoryBranchPermissionsResponse struct {
	Size       int                        `json:"size"`
	Limit      int                        `json:"limit"`
	IsLastPage bool                       `json:"isLastPage"`
	Values     []BranchPermissionResponse `json:"values"`
}

func resourceBranchPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceBranchPermissionsCreate,
		Read:   resourceBranchPermissionsRead,
		Update: resourceBranchPermissionsCreate,
		Delete: resourceBranchPermissionsDelete,
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
			"ref_pattern": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"pull-request-only", "fast-forward-only", "no-deletes", "read-only"}, false),
			},
			"exception_users": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"exception_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"exception_access_keys": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"permission_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func newBranchPermissionPayloadFromResource(d *schema.ResourceData) *BranchPermissionPayload {
	branchPermissionPayload := &BranchPermissionPayload{
		Type: d.Get("type").(string),
	}

	for _, item := range d.Get("exception_users").([]interface{}) {
		branchPermissionPayload.Users = append(branchPermissionPayload.Users, item.(string))
	}

	for _, item := range d.Get("exception_groups").([]interface{}) {
		branchPermissionPayload.Groups = append(branchPermissionPayload.Groups, item.(string))
	}

	for _, item := range d.Get("exception_access_keys").([]interface{}) {
		branchPermissionPayload.AccessKeys = append(branchPermissionPayload.AccessKeys, item.(string))
	}

	matcherConfig := &MatcherStruct{
		Id:        d.Get("ref_pattern").(string),
		DisplayId: d.Get("ref_pattern").(string),
		Type: MatcherStructType{
			Id:   "PATTERN",
			Name: "Pattern",
		},
		Active: true,
	}

	branchPermissionPayload.Matcher = *matcherConfig

	return branchPermissionPayload
}

func resourceBranchPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	branchPermission := newBranchPermissionPayloadFromResource(d)

	request, err := json.Marshal(branchPermission)

	if err != nil {
		return err
	}

	res, err := client.Post(fmt.Sprintf("/rest/branch-permissions/2.0/projects/%s/repos/%s/restrictions",
		project,
		repository,
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	var branchPermissionResponse BranchPermissionResponse

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &branchPermissionResponse)

	if err != nil {
		return err
	}

	_ = d.Set("permission_id", branchPermissionResponse.Id)

	d.SetId(fmt.Sprintf("%s|%s|%s|%s",
		d.Get("project").(string),
		d.Get("repository").(string),
		d.Get("ref_pattern").(string),
		d.Get("type").(string)),
	)
	return resourceBranchPermissionsRead(d, m)
}

func resourceBranchPermissionsRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "|")
		if len(parts) == 4 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("ref_pattern", parts[2])
			_ = d.Set("type", parts[3])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project|repository|ref_pattern|type`")
		}
	}

	branchPermissionId := d.Get("permission_id")

	var err error

	if branchPermissionId == nil {
		err = getBranchPermissionFromList(d, m)
	} else {
		err = getBranchPermissionById(d, m)
	}

	if err != nil {
		return err
	}

	return nil
}

func getBranchPermissionById(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	id := d.Get("permission_id").(int)

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/branch-permissions/2.0/projects/%s/repos/%s/restrictions/%d",
		project,
		repository,
		id,
	))

	if err != nil {
		return err
	}

	var branchPermissionResponse BranchPermissionResponse

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&branchPermissionResponse)

	if err != nil {
		return err
	}

	_ = d.Set("permission_id", branchPermissionResponse.Id)
	_ = d.Set("type", branchPermissionResponse.Type)
	_ = d.Set("exception_groups", branchPermissionResponse.Groups)

	// Convert slice of structs back to slice object for branchPermissionResponse.Users
	exceptionUsers := make([]string, 0, len(branchPermissionResponse.Users))
	for _, item := range branchPermissionResponse.Users {
		exceptionUsers = append(exceptionUsers, item.Name)
	}
	_ = d.Set("exception_users", exceptionUsers)

	// Convert slice of structs back to slice object for branchPermissionResponse.Users
	exceptionAccessKeys := make([]int, 0, len(branchPermissionResponse.AccessKeys))
	for _, item := range branchPermissionResponse.AccessKeys {
		exceptionAccessKeys = append(exceptionAccessKeys, item.Key.ID)
	}
	_ = d.Set("exception_access_keys", exceptionAccessKeys)

	return nil
}

func getBranchPermissionFromList(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	restrictionType := d.Get("type").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/branch-permissions/2.0/projects/%s/repos/%s/restrictions",
		project,
		repository,
	))

	if err != nil {
		return err
	}

	var allRepositoryBranchPermissionsResponse AllRepositoryBranchPermissionsResponse

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&allRepositoryBranchPermissionsResponse)

	if err != nil {
		return err
	}

	for _, item := range allRepositoryBranchPermissionsResponse.Values {
		if strings.ToLower(strings.Replace(item.Type, "_", "-", -1)) == restrictionType {
			_ = d.Set("permission_id", item.Id)
			_ = d.Set("type", item.Type)
			_ = d.Set("exception_groups", item.Groups)

			// Convert slice of structs back to slice object for exception_users
			exceptionUsers := make([]string, 0, len(item.Users))
			for _, item := range item.Users {
				exceptionUsers = append(exceptionUsers, item.Name)
			}
			_ = d.Set("exception_users", exceptionUsers)

			// Convert slice of structs back to slice object for exception_access_keys
			exceptionAccessKeys := make([]int, 0, len(item.AccessKeys))
			for _, item := range item.AccessKeys {
				exceptionAccessKeys = append(exceptionAccessKeys, item.Key.ID)
			}
			_ = d.Set("exception_access_keys", exceptionAccessKeys)

			return nil
		}
	}

	return fmt.Errorf("incorrect ID format, should match `project|repository|ref_pattern`")
}

func resourceBranchPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/branch-permissions/2.0/projects/%s/repos/%s/restrictions/%d",
		d.Get("project").(string),
		d.Get("repository").(string),
		d.Get("permission_id").(int)))

	return err
}
