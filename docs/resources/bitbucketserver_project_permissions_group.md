# Resource: bitbucketserver_project_permissions_group

Set project level permissions for a given group.

## Example Usage

```hcl
resource "bitbucketserver_project_permissions_group" "my-proj" {
  project    = "MYPROJ"
  group      = "stash-users"
  permission = "PROJECT_WRITE"
}
```

## Argument Reference

* `project` - Required. Project key to set permissions for.
* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are:

    * `PROJECT_READ`
    * `PROJECT_WRITE`
    * `PROJECT_ADMIN`

## Import

Import a group project level permissions via the project & group names:

```
terraform import bitbucketserver_project_permissions_group.test MYPROJ/stash-users
```
