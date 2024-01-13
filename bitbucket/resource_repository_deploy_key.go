package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type SSHKey struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	CreatedDate   jsonTime             `json:"createdDate,omitempty"`
	UpdatedDate   jsonTime             `json:"updatedDate,omitempty"`
	URL           string               `json:"url,omitempty"`
	Active        bool                 `json:"active"`
	Events        []interface{}        `json:"events"`
	Configuration WebhookConfiguration `json:"configuration"`
}

func resourceRepositoryDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryDeployKeyCreate,
		Read:   resourceRepositoryDeployKeyRead,
		Delete: resourceRepositoryDeployKeyDelete,
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
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REPO_READ", "REPO_WRITE", "REPO_ADMIN"}, false),
			},
		},
	}
}

type KeyRequestKey struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

type KeyRequest struct {
	Key        KeyRequestKey `json:"key"`
	Permission string        `json:"permission"`
}

type KeyResponse struct {
	Key struct {
		Id                int    `json:"id"`
		ExpiryDays        int    `json:"expiryDays"`
		BitLength         int    `json:"bitLength"`
		AlgorithmType     string `json:"algorithmType"`
		CreatedDate       string `json:"createdDate"`
		Fingerprint       string `json:"fingerprint"`
		LastAuthenticated string `json:"lastAuthenticated"`
		Label             string `json:"label"`
		Text              string `json:"text"`
	} `json:"key"`
	Permission string `json:"permission"`
	Project    struct {
		Name        string `json:"name"`
		Key         string `json:"key"`
		Id          int    `json:"id"`
		Type        string `json:"type"`
		Public      bool   `json:"public"`
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		Namespace   string `json:"namespace"`
		Scope       string `json:"scope"`
	} `json:"project"`
	Repository struct {
		Name          string `json:"name"`
		Id            int    `json:"id"`
		State         string `json:"state"`
		Public        bool   `json:"public"`
		DefaultBranch string `json:"defaultBranch"`
		HierarchyId   string `json:"hierarchyId"`
		StatusMessage string `json:"statusMessage"`
		Archived      bool   `json:"archived"`
		Forkable      bool   `json:"forkable"`
		RelatedLinks  struct {
		} `json:"relatedLinks"`
		Partition int `json:"partition"`
		Origin    struct {
			Name          string `json:"name"`
			Id            int    `json:"id"`
			State         string `json:"state"`
			Public        bool   `json:"public"`
			DefaultBranch string `json:"defaultBranch"`
			HierarchyId   string `json:"hierarchyId"`
			StatusMessage string `json:"statusMessage"`
			Archived      bool   `json:"archived"`
			Forkable      bool   `json:"forkable"`
			RelatedLinks  struct {
			} `json:"relatedLinks"`
			Partition   int    `json:"partition"`
			Description string `json:"description"`
			Project     struct {
				Name        string `json:"name"`
				Key         string `json:"key"`
				Id          int    `json:"id"`
				Type        string `json:"type"`
				Public      bool   `json:"public"`
				Avatar      string `json:"avatar"`
				Description string `json:"description"`
				Namespace   string `json:"namespace"`
				Scope       string `json:"scope"`
			} `json:"project"`
			Scope string `json:"scope"`
			ScmId string `json:"scmId"`
			Slug  string `json:"slug"`
		} `json:"origin"`
		Description string `json:"description"`
		Project     struct {
			Name        string `json:"name"`
			Key         string `json:"key"`
			Id          int    `json:"id"`
			Type        string `json:"type"`
			Public      bool   `json:"public"`
			Avatar      string `json:"avatar"`
			Description string `json:"description"`
			Namespace   string `json:"namespace"`
			Scope       string `json:"scope"`
		} `json:"project"`
		Scope string `json:"scope"`
		ScmId string `json:"scmId"`
		Slug  string `json:"slug"`
	} `json:"repository"`
}

func resourceRepositoryDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	keyRequest := &KeyRequest{
		Key: KeyRequestKey{
			Label: d.Get("label").(string),
			Text:  d.Get("key").(string),
		},
		Permission: d.Get("permission").(string),
	}

	request, err := json.Marshal(keyRequest)

	client := m.(*BitbucketServerProvider).BitbucketClient
	resp, err := client.Post(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}
	return storeResponse(d, resp)
}

func resourceRepositoryDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh/%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Id()),
	))
	if err != nil {
		return err
	}
	return storeResponse(d, resp)
}

func resourceRepositoryDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh/%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Id()),
	))
	return err
}

func storeResponse(d *schema.ResourceData, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var keyResponse KeyResponse
	err = json.Unmarshal(body, &keyResponse)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(keyResponse.Key.Id))
	labelError := store(d, keyResponse, keyResponse.Key.Label, "label")
	permissionError := store(d, keyResponse, keyResponse.Permission, "permission")
	keyError := store(d, keyResponse, keyResponse.Key.Text, "key")
	projectError := store(d, keyResponse, keyResponse.Repository.Project.Name, "project")
	repositoryError := store(d, keyResponse, keyResponse.Repository.Name, "repository")

	return errors.Join(labelError, permissionError, keyError, projectError, repositoryError)
}

func store(d *schema.ResourceData, keyResponse KeyResponse, value string, name string) error {
	if value == "" {
		respAsJson, _ := json.Marshal(keyResponse)
		return errors.New(fmt.Sprintf("%s is nil in %s", name, string(respAsJson)))
	} else {
		return d.Set(name, value)
	}
}
