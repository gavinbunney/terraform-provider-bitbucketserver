# Resource: bitbucketserver_repository_permissions_user

Set repository permissions for a given user.

## Example Usage

```hcl
resource "bitbucketserver_repository_permissions_user" "mreynolds-repo" {
  project    = "MYPROJ"
  repository = "repo1"
  user       = "mreynolds"
  permission = "REPO_WRITE"
}
```

## Argument Reference

* `project` - Required. Project key the repository is contained within.
* `repository` - Required. Repository slug to set the permissions for.
* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are:
                                                    
    * `REPO_READ`
    * `REPO_WRITE`
    * `REPO_ADMIN`

## Import

Import a group project level permissions via the project & group names:

```
terraform import bitbucketserver_repository_permissions_user.test MYPROJ/repo1/mreynolds
```
