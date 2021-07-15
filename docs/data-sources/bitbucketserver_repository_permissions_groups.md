# Data Source: bitbucketserver_repository_permissions_groups

Retrieve a list of groups that have been granted at least one repository level permission to the specified repo.

## Example Usage

```hcl
data "bitbucketserver_repository_permissions_groups" "my-repo-groups" {
  project    = "TEST"
  repository = "my-repo"
}
```

### Applying a Custom Filter

Find project groups starting with `dev` with project permissions.
 
```hcl
data "bitbucketserver_repository_permissions_groups" "my-repo-dev-groups" {
  project = "TEST"
  repository = "my-repo"
  filter  = "dev"
}
```

## Argument Reference

* `project` - Required. Project Key to lookup permissions for.
* `repository` - Required. Repository slug to lookup permissions for.
* `filter` - Optional. If specified only group names containing the supplied string will be returned.

## Attribute Reference

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are:

    * `REPO_READ`
    * `REPO_WRITE`
    * `REPO_ADMIN`
