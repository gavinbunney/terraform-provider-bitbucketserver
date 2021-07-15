# Data Source: bitbucketserver_cluster

Gets information about the nodes that currently make up the Bitbucket cluster.

## Example Usage

```hcl
data "bitbucketserver_cluster" "main" { }

output "local_hostname" {
  value = "Bitbucket running on ${data.bitbucketserver_cluster.main.local_node.0.hostname}"
}
```

## Attribute Reference

* `local_node` - List with a single element, containing the local node details. See `node` schema below.
* `nodes` - List of nodes of the Bitbucket cluster.
* `running` - Flag is the cluster is running.

### Node Schema

Each node in the attributes above contains the following elements:

* `id` - Unique cluster identifier.
* `name` - Unique cluster identifier.
* `hostname` - Address hostname of the cluster node. Typically an IP address.
* `port` - Port of the cluster node. This is not the same as the Bitbucket UI port, rather the node cluster port.
* `local` - Flag if this is a local node.
