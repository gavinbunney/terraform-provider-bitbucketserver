# Resource: bitbucketserver_project_permissions_user

Set project permissions for a given user.

## Example Usage

```hcl
resource "bitbucketserver_project_permissions_user" "mreynolds" {
  project    = "MYPROJ"
  user       = "mreynolds"
  permission = "PROJECT_WRITE"
}
```

## Argument Reference

* `project` - Required. Name of the project to assign permissions to.
* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are:

    * `PROJECT_READ`
    * `PROJECT_WRITE`
    * `PROJECT_ADMIN`

## Import

Import a user project permissions via the project and user name:

```
terraform import bitbucketserver_global_permissions_user.test MYPROJ/mreynolds
```
