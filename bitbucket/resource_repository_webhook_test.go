package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceRepositoryWebhook_simple(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key  = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name    = "repo"
		}

		resource "bitbucketserver_repository_webhook" "test" {
			project     = bitbucketserver_project.test.key
			repository  = bitbucketserver_repository.test.slug
			name        = "google"
			webhook_url = "https://www.google.com/"
			events      = ["repo:refs_changed"]
			active      = true
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "webhook_url", "https://www.google.com/"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "secret", ""),
				),
			},
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccBitbucketResourceRepositoryWebhook_complete(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	configString := `
		resource "bitbucketserver_project" "test" {
			key  = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name    = "repo"
		}

		resource "bitbucketserver_repository_webhook" "test" {
			project     = bitbucketserver_project.test.key
			repository  = bitbucketserver_repository.test.slug
			name        = "%v"
			webhook_url = "%v"
			secret      = "abc"
			events      = ["repo:refs_changed"]
			active      = true
		}
	`

	config := fmt.Sprintf(configString, projectKey, projectKey, "test", "https://www.oldurl.com/")
	newConfig := fmt.Sprintf(configString, projectKey, projectKey, "test2", "https://www.newurl.com/")

	// Create resource
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "name", "test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "webhook_url", "https://www.oldurl.com/"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "secret", "abc"),
				),
			},
			{
				Config: newConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "name", "test2"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "webhook_url", "https://www.newurl.com/"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_webhook.test", "secret", "abc"),
				),
			},
			{
				Config:  newConfig,
				Destroy: true,
			},
		},
	})
}
