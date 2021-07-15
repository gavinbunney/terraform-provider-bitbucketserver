# Resource: bitbucketserver_global_permissions_group

Set global permissions for a given group.

## Example Usage

```hcl
resource "bitbucketserver_global_permissions_group" "test" {
  group      = "stash-users"
  permission = "ADMIN"
}
```

## Argument Reference

* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available global permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`

## Import

Import a group global permissions via the group name:

```
terraform import bitbucketserver_global_permissions_group.test my-group
```
