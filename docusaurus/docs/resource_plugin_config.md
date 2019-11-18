---
id: bitbucketserver_plugin_config
title: bitbucketserver_plugin_config
---

Configure plugins.

## Example Usage

```hcl
resource "bitbucketserver_plugin_config" "mypluginconfig" {
  key    = "my-plugin-key"
  values = "{"\key\": \"value\"}"
}
```

## Argument Reference

* `key` - Required. Unique key of the plugin. This is not the same used by plugin install.
* `values` - Required. Plugin configuration in JSON format.

## Attribute Reference

* `key` - Unique key of the plugin. This is not the same used by plugin install.
* `validlicense` - Indicates if the plugins has a valid license.
* `values` - Plugin configuration in JSON format.

## Import

Import a plugin config reference via the key:

```
terraform import bitbucketserver_plugin_config.mypluginkey my-plugin-key
```
