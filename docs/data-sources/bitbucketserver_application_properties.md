# Data Source: bitbucketserver_application_properties

This data source allows you to retrieve version information and other application properties of Bitbucket Server.

## Example Usage

```hcl
data "bitbucketserver_application_properties" "main" { }
```

## Attribute Reference

* `version` - Version of Bitbucket.
* `build_number` - Build number of the Bitbucket instance.
* `build_date` - Date the Bitbucket build was made,
* `display_name` - Name of the Bitbucket instance.
