package bitbucket

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testRepo string = "test-repo"

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"bitbucketserver": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("BITBUCKET_SERVER"); v == "" {
		t.Fatal("BITBUCKET_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("BITBUCKET_USERNAME"); v == "" {
		t.Fatal("BITBUCKET_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("BITBUCKET_PASSWORD"); v == "" {
		t.Fatal("BITBUCKET_PASSWORD must be set for acceptance tests")
	}
}
