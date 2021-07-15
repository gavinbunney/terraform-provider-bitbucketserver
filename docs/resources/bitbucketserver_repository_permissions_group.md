# Resource: bitbucketserver_repository_permissions_group

Set repository level permissions for a given group.

## Example Usage

```hcl
resource "bitbucketserver_repository_permissions_group" "my-repo" {
  project    = "MYPROJ"
  repository = "repo1"
  group      = "stash-users"
  permission = "REPO_WRITE"
}
```

## Argument Reference

* `project` - Required. Project key the repository is contained within.
* `repository` - Required. Repository slug to set the permissions for.
* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available permissions are:
                                                    
    * `REPO_READ`
    * `REPO_WRITE`
    * `REPO_ADMIN`

## Import

Import a group project level permissions via the project & group names:

```
terraform import bitbucketserver_repository_permissions_group.test MYPROJ/repo1/stash-users
```
