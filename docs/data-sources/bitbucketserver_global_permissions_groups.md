# Data Source: bitbucketserver_global_permissions_groups

Retrieve a list of groups that have been granted at least one global permission.

## Example Usage

```hcl
data "bitbucketserver_global_permissions_groups" "all" { }
```

### Applying a Custom Filter

Find any groups starting with `dev`.
 
```hcl
data "bitbucketserver_global_permissions_groups" "dev-groups" {
  filter = "dev"
}
```

## Argument Reference

* `filter` - Optional. If specified only group names containing the supplied string will be returned.

## Attribute Reference

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are:

    * `LICENSED_USER`
    * `PROJECT_CREATE`
    * `ADMIN`
    * `SYS_ADMIN`

