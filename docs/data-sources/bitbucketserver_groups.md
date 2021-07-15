# Data Source: bitbucketserver_groups

This data source allows you to retrieve a list of groups, optionally matching the supplied `filter`.

## Example Usage

```hcl
data "bitbucketserver_groups" "all" { }
```

### Applying a Custom Filter

Find any groups starting with `dev`.
 
```hcl
data "bitbucketserver_groups" "dev-groups" {
  filter = "dev"
}
```

## Argument Reference

* `filter` - Optional. If specified only group names containing the supplied string will be returned.

## Attribute Reference

* `groups` - List of maps containing `name` key.
