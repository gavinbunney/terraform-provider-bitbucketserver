# Resource: bitbucketserver_plugin_config

Configure plugins.

## Example Usage

```hcl
resource "bitbucketserver_plugin_config" "mypluginconfig" {
  config_endpoint = "/rest/1.0/myplugin/config"
  values          = "{"\key\": \"value\"}"
}
```

## Argument Reference

* `config_endpoint` - Required. Path to the configuration endpoint. Relative to the bitbucket server url configured in the provider.
* `values` - Required. Plugin configuration in JSON format.

## Import

Import a plugin config reference via the key:

```
terraform import bitbucketserver_plugin_config.mypluginkey my-plugin-key
```
