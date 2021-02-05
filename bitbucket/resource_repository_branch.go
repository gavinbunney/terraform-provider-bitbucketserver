package bitbucket

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRepositoryBranch() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryBranchCreate,
		Update: resourceRepositoryBranchUpdate,
		Read:   resourceRepositoryBranchRead,
		Delete: resourceRepositoryBranchDelete,
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
			"branch_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_ref": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: true,
				Default:  "refs/head/master",
			},
		},
	}
}

func resourceRepositoryBranchRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			d.Set("project", parts[0])
			d.Set("repository", parts[1])
			d.Set("branch_name", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/branch_name`")
		}
	}
	return nil
}

func resourceRepositoryBranchCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceBannerUpdate(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("branch_name").(string)))
	return resourceRepositoryBranchRead(d, m)
}

func resourceRepositoryBranchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	branchName := d.Get("branch_name").(string)
	sourceRef := d.Get("source_ref").(string)

	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{ "name": "%s", "startPoint": "%s" }`, branchName, sourceRef)))
	_, err := client.Post(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/branches", project, repository), body)
	if err != nil {
		return err
	}
	return resourceRepositoryBranchRead(d, m)
}

func resourceRepositoryBranchDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	branchName := d.Get("branch_name").(string)

	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{ "name": "%s" }`, branchName)))
	_, err := client.Do(http.MethodDelete, fmt.Sprintf("/rest/branch-utils/1.0/projects/%s/repos/%s/branches", project, repository), body, "application/json")
	return err
}
