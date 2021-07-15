# Resource: bitbucketserver_global_permissions_user

Set global permissions for a given user.

## Example Usage

```hcl
resource "bitbucketserver_global_permissions_user" "mreynolds" {
  user       = "mreynolds"
  permission = "ADMIN"
}
```

## Argument Reference

* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available global permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`

## Import

Import a user global permissions via the user's name:

```
terraform import bitbucketserver_global_permissions_user.test mreynolds
```
