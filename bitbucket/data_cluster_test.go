package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataCluster(t *testing.T) {
	config := `
		data "bitbucketserver_cluster" "main" {}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_cluster.main", "running", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_cluster.main", "local_node.#", "1"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "local_node.0.id"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "local_node.0.hostname"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "local_node.0.port"),
					resource.TestCheckResourceAttr("data.bitbucketserver_cluster.main", "local_node.0.local", "true"),
					resource.TestCheckResourceAttr("data.bitbucketserver_cluster.main", "nodes.#", "1"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "nodes.0.id"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "nodes.0.hostname"),
					resource.TestCheckResourceAttrSet("data.bitbucketserver_cluster.main", "nodes.0.port"),
					resource.TestCheckResourceAttr("data.bitbucketserver_cluster.main", "nodes.0.local", "true"),
				),
			},
		},
	})
}
