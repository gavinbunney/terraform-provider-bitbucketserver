package bitbucket

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type ClusterNode struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address struct {
		Hostname string `json:"hostName,omitempty"`
		Port     int    `json:"port,omitempty"`
	} `json:"address,omitempty"`
	Local bool `json:"local,omitempty"`
}

type Cluster struct {
	LocalNode ClusterNode   `json:"localNode,omitempty"`
	Nodes     []ClusterNode `json:"nodes,omitempty"`
	Running   bool          `json:"running,omitempty"`
}

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"local_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterNodeResourceSchema(),
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterNodeResourceSchema(),
			},
			"running": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func clusterNodeResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get("/rest/api/1.0/admin/cluster")

	if err != nil {
		return err
	}

	var cluster Cluster

	body, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		return readErr
	}

	decodeErr := json.Unmarshal(body, &cluster)
	if decodeErr != nil {
		return decodeErr
	}

	d.SetId("cluster")
	_ = d.Set("running", cluster.Running)

	var nodes []interface{}
	for _, node := range cluster.Nodes {
		n := make(map[string]interface{})
		n["id"] = node.ID
		n["name"] = node.Name
		n["hostname"] = node.Address.Hostname
		n["port"] = node.Address.Port
		n["local"] = node.Local
		nodes = append(nodes, n)
	}
	_ = d.Set("nodes", nodes)

	var localNode []interface{}
	n := make(map[string]interface{})
	n["id"] = cluster.LocalNode.ID
	n["name"] = cluster.LocalNode.Name
	n["hostname"] = cluster.LocalNode.Address.Hostname
	n["port"] = cluster.LocalNode.Address.Port
	n["local"] = cluster.LocalNode.Local
	localNode = append(localNode, n)
	_ = d.Set("local_node", localNode)

	return nil
}
