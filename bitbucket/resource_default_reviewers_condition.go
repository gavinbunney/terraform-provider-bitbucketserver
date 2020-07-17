package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type Reviewer struct {
	ID int `json:"id,omitempty"`
}

type MatcherType struct {
	ID string `json:"id,omitempty"`
}

type Matcher struct {
	ID   string      `json:"id,omitempty"`
	Type MatcherType `json:"type,omitempty"`
}

type RefMatcher struct {
	ID   string `json:"id,omitempty"`
	Type struct {
		ID string `json:"id,omitempty"`
	} `json:"type,omitempty"`
}

type DefaultReviewersConditionPayload struct {
	SourceMatcher     Matcher    `json:"sourceMatcher,omitempty"`
	TargetMatcher     Matcher    `json:"targetMatcher,omitempty"`
	Reviewers         []Reviewer `json:"reviewers,omitempty"`
	RequiredApprovals int        `json:"requiredApprovals,omitempty"`
}

type DefaultReviewersConditionResp struct {
	ID                int        `json:"id,omitempty"`
	RequiredApprovals int        `json:"requiredApprovals,omitempty"`
	Reviewers         []Reviewer `json:"reviewers,omitempty"`
	SourceRefMatcher  RefMatcher `json:"sourceRefMatcher,omitempty"`
	TargetRefMatcher  RefMatcher `json:"targetRefMatcher,omitempty"`
}

var matcherDesc = `id can be either "any" to match all branches, "refs/heads/master" to match certain branch, "pattern" to match multiple branches or "development" to match branching model. type_id must be one of: "ANY_REF", "BRANCH", "PATTERN", "MODEL_BRANCH".`

var validMatcherTypeIDs = []string{
	"ANY_REF", "BRANCH", "PATTERN", "MODEL_BRANCH",
}

func resourceDefaultReviewersCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceDefaultReviewersConditionCreate,
		Read:   resourceDefaultReviewersConditionRead,
		Delete: resourceDefaultReviewersConditionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_slug": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_matcher": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				ForceNew:    true,
				Description: matcherDesc,
			},
			"target_matcher": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				ForceNew:    true,
				Description: matcherDesc,
			},
			"reviewers": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
				Set:         schema.HashInt,
				ForceNew:    true,
				MinItems:    1,
				Description: "IDs of users to become default reviewers when you create a pull request.",
			},
			"required_approvals": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The number of default reviewers that must approve a pull request.",
			},
		},
	}
}

func refMatcherToMatcher(refMatcher RefMatcher) Matcher {
	convertID := func(id string) string {
		if id == "ANY_REF_MATCHER_ID" {
			return "any"
		}

		return id
	}

	return Matcher{
		ID: convertID(refMatcher.ID),
		Type: MatcherType{
			ID: refMatcher.Type.ID,
		},
	}
}

func expandReviewers(set *schema.Set) []Reviewer {
	list := set.List()

	rs := make([]Reviewer, 0, len(list))

	for _, v := range list {
		rs = append(rs, Reviewer{ID: v.(int)})
	}

	return rs
}

func collapseReviewers(reviewers []Reviewer) *schema.Set {
	reviewerIDs := make([]interface{}, 0)

	for _, r := range reviewers {
		reviewerIDs = append(reviewerIDs, r.ID)
	}

	return schema.NewSet(schema.HashInt, reviewerIDs)
}

func expandMatcher(matcherMap map[string]interface{}) Matcher {
	return Matcher{
		ID: matcherMap["id"].(string),
		Type: MatcherType{
			ID: matcherMap["type_id"].(string),
		},
	}
}

func collapseMatcher(matcher Matcher) map[string]interface{} {
	return map[string]interface{}{
		"id":      matcher.ID,
		"type_id": matcher.Type.ID,
	}
}

func parseResourceID(id string) (string, string, string, error) {
	parts := strings.Split(id, ":")

	errMsg := fmt.Errorf("invalid format of ID (%s), expected condition_id:project_key or condition_id:project_key:repository_slug", id)

	if len(parts) != 2 && len(parts) != 3 {
		return "", "", "", errMsg
	}

	if len(parts) == 2 {
		if parts[0] == "" || parts[1] == "" {
			return "", "", "", errMsg
		}

		return parts[0], parts[1], "", nil
	}

	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", errMsg
	}

	return parts[0], parts[1], parts[2], nil
}

func createResourceID(conditionID int, projectKey string, repositorySlug string) string {
	if repositorySlug == "" {
		return fmt.Sprintf("%v:%s", conditionID, projectKey)
	}

	return fmt.Sprintf("%v:%s:%s", conditionID, projectKey, repositorySlug)
}

func getCreateConditionURI(projectKey string, repositorySlug string) string {
	if repositorySlug == "" {
		return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/condition",
			url.PathEscape(projectKey),
		)
	}

	return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/repos/%s/condition",
		url.PathEscape(projectKey),
		url.PathEscape(repositorySlug),
	)
}

func getReadConditionURI(projectKey string, repositorySlug string) string {
	if repositorySlug == "" {
		return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/conditions",
			url.PathEscape(projectKey),
		)
	}

	return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/repos/%s/conditions",
		url.PathEscape(projectKey),
		url.PathEscape(repositorySlug),
	)
}

func getDeleteConditionURI(conditionID string, projectKey string, repositorySlug string) string {
	if repositorySlug == "" {
		return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/condition/%s",
			url.PathEscape(projectKey),
			url.PathEscape(conditionID),
		)
	}

	return fmt.Sprintf("/rest/default-reviewers/1.0/projects/%s/repos/%s/condition/%s",
		url.PathEscape(projectKey),
		url.PathEscape(repositorySlug),
		url.PathEscape(conditionID),
	)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func resourceDefaultReviewersConditionCreate(d *schema.ResourceData, m interface{}) error {
	projectKey := d.Get("project_key").(string)
	repositorySlug := d.Get("repository_slug").(string)

	condition := &DefaultReviewersConditionPayload{
		SourceMatcher:     expandMatcher(d.Get("source_matcher").(map[string]interface{})),
		TargetMatcher:     expandMatcher(d.Get("target_matcher").(map[string]interface{})),
		Reviewers:         expandReviewers(d.Get("reviewers").(*schema.Set)),
		RequiredApprovals: d.Get("required_approvals").(int),
	}

	if contains(validMatcherTypeIDs, condition.SourceMatcher.Type.ID) == false {
		return fmt.Errorf("source_matcher.type_id %s must be one of %v", condition.SourceMatcher.Type.ID, validMatcherTypeIDs)
	}

	if contains(validMatcherTypeIDs, condition.TargetMatcher.Type.ID) == false {
		return fmt.Errorf("target_matcher.type_id %s must be one of %v", condition.TargetMatcher.Type.ID, validMatcherTypeIDs)
	}

	if condition.RequiredApprovals > len(condition.Reviewers) {
		return fmt.Errorf("required_approvals %d cannot be more than length of reviewers %d", condition.RequiredApprovals, len(condition.Reviewers))
	}

	bytedata, err := json.Marshal(condition)

	if err != nil {
		return err
	}

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Post(getCreateConditionURI(projectKey, repositorySlug), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	var newCondition DefaultReviewersConditionResp

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &newCondition)

	if err != nil {
		return err
	}

	d.SetId(createResourceID(newCondition.ID, projectKey, repositorySlug))

	return resourceDefaultReviewersConditionRead(d, m)
}

func resourceDefaultReviewersConditionRead(d *schema.ResourceData, m interface{}) error {
	conditionID, projectKey, repositorySlug, err := parseResourceID(d.Id())

	if err != nil {
		return err
	}

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(getReadConditionURI(projectKey, repositorySlug))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to find a matching default reviewers condition %s. API returned %d", d.Id(), resp.StatusCode)
	}

	var conditions []DefaultReviewersConditionResp

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &conditions)

	if err != nil {
		return err
	}

	for _, c := range conditions {
		cID := strconv.Itoa(c.ID)

		if cID == conditionID {
			d.Set("project_key", projectKey)
			d.Set("repository_slug", repositorySlug)
			d.Set("source_matcher", collapseMatcher(refMatcherToMatcher(c.SourceRefMatcher)))
			d.Set("target_matcher", collapseMatcher(refMatcherToMatcher(c.TargetRefMatcher)))
			d.Set("reviewers", collapseReviewers(c.Reviewers))
			d.Set("required_approvals", c.RequiredApprovals)

			return nil
		}
	}

	return fmt.Errorf("unable to find a matching default reviewers condition %s", d.Id())
}

func resourceDefaultReviewersConditionDelete(d *schema.ResourceData, m interface{}) error {
	conditionID, projectKey, repositorySlug, err := parseResourceID(d.Id())

	if err != nil {
		return err
	}

	client := m.(*BitbucketServerProvider).BitbucketClient

	_, err = client.Delete(getDeleteConditionURI(conditionID, projectKey, repositorySlug))

	return err
}
