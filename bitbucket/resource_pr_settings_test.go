package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourcePrSettings_requiredArgumentsOnly(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForRepositoryBasedTests(projectKey) + `
		resource "bitbucketserver_pr_settings" "test" {
			project    = bitbucketserver_project.test.key
			repository = bitbucketserver_repository.test.name
			merge_config {
				default_strategy   = "no-ff"
				enabled_strategies = ["no-ff"]
			}
		}
	`
	resourceName := "bitbucketserver_pr_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.default_strategy", "no-ff"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.enabled_strategies.0", "no-ff"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.commit_summaries", "20"),
					resource.TestCheckResourceAttr(resourceName, "no_needs_work_status", "false"),
					resource.TestCheckResourceAttr(resourceName, "required_all_approvers", "false"),
					resource.TestCheckResourceAttr(resourceName, "required_all_tasks_complete", "false"),
					resource.TestCheckResourceAttr(resourceName, "required_approvers", "0"),
					resource.TestCheckResourceAttr(resourceName, "required_successful_builds", "0"),
				),
			},
		},
	})
}

func TestAccBitbucketResourcePrSettings_allArguments(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForRepositoryBasedTests(projectKey) + `
		resource "bitbucketserver_pr_settings" "test" {
			project                     = bitbucketserver_project.test.key
			repository                  = bitbucketserver_repository.test.name
			no_needs_work_status        = true
			required_all_approvers      = true
			required_all_tasks_complete = true
			required_approvers          = 1
			required_successful_builds  = 1
			merge_config {
				default_strategy   = "no-ff"
				enabled_strategies = ["no-ff", "ff"]
				commit_summaries   = 30
			}
		}
	`

	resourceName := "bitbucketserver_pr_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.default_strategy", "no-ff"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.enabled_strategies.0", "no-ff"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.enabled_strategies.1", "ff"),
					resource.TestCheckResourceAttr(resourceName, "merge_config.0.commit_summaries", "30"),
					resource.TestCheckResourceAttr(resourceName, "no_needs_work_status", "true"),
					resource.TestCheckResourceAttr(resourceName, "required_all_approvers", "true"),
					resource.TestCheckResourceAttr(resourceName, "required_all_tasks_complete", "true"),
					resource.TestCheckResourceAttr(resourceName, "required_approvers", "1"),
					resource.TestCheckResourceAttr(resourceName, "required_successful_builds", "1"),
				),
			},
		},
	})
}
