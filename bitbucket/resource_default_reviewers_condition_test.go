package bitbucket

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketDefaultReviewersCondition_forProject(t *testing.T) {
	key := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDefaultReviewersConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDefaultReviewersConditionResourceForProject(key, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "project_key", "TEST"+key),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "reviewers.#", "1"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "required_approvals", "1"),
				),
			},
		},
	})
}

func TestAccBitbucketDefaultReviewersCondition_noRequiredApprovals(t *testing.T) {
	key := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDefaultReviewersConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDefaultReviewersConditionResourceForProject(key, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "project_key", "TEST"+key),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "reviewers.#", "1"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "required_approvals", "0"),
				),
			},
		},
	})
}

func TestAccBitbucketDefaultReviewersCondition_forRepository(t *testing.T) {
	key := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDefaultReviewersConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDefaultReviewersConditionResourceForRepository(key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "project_key", "TEST"+key),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "repository_slug", "test-repo-"+key),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "source_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.id", "any"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "target_matcher.type_id", "ANY_REF"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "reviewers.#", "1"),
					resource.TestCheckResourceAttr("bitbucketserver_default_reviewers_condition.test", "required_approvals", "1"),
				),
			},
			{
				ResourceName:      "bitbucketserver_default_reviewers_condition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBitbucketDefaultReviewersCondition_expectRequiredApprovalsError(t *testing.T) {
	key := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBitbucketDefaultReviewersConditionResourceForProject(key, 2),
				ExpectError: regexp.MustCompile("required_approvals 2 cannot be more than length of reviewers 1"),
			},
		},
	})
}

func testAccCheckBitbucketDefaultReviewersConditionDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucketserver_default_reviewers_condition.test"]

	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_default_reviewers_condition.test")
	}

	return nil
}

func testAccBitbucketDefaultReviewersConditionResourceForProject(key string, requiredApprovals int) string {
	return fmt.Sprintf(`
resource "bitbucketserver_project" "test" {
	key = "TEST%s"
	name = "test-project-%s"
}

data "bitbucketserver_user" "reviewer" {
	name = "admin"
}

resource "bitbucketserver_default_reviewers_condition" "test" {
	project_key			= bitbucketserver_project.test.key
	source_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	target_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	reviewers          = [data.bitbucketserver_user.reviewer.user_id]
	required_approvals = %d
}`, key, key, requiredApprovals)
}

func testAccBitbucketDefaultReviewersConditionResourceForRepository(key string) string {
	return fmt.Sprintf(`
resource "bitbucketserver_project" "test" {
	key = "TEST%s"
	name = "test-project-%s"
}

resource "bitbucketserver_repository" "test" {
	project = bitbucketserver_project.test.key
	name = "test-repo-%s"
}

data "bitbucketserver_user" "reviewer" {
	name = "admin"
}

resource "bitbucketserver_default_reviewers_condition" "test" {
	project_key			= bitbucketserver_project.test.key
	repository_slug		= bitbucketserver_repository.test.slug
	source_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	target_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	reviewers          = [data.bitbucketserver_user.reviewer.user_id]
	required_approvals = 1
}`, key, key, key)
}
