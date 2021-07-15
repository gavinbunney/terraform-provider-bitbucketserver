# Data Source: bitbucketserver_repository_permissions_users

Retrieve a list of users that have been granted at least one permission for the specified repository.

## Example Usage

```hcl
data "bitbucketserver_repository_permissions_users" "my-repo-users" {
  project    = "TEST"
  repository = "my-repo"
}
```

### Applying a Custom Filter

Find repository users starting with `malcolm`.
 
```hcl
data "bitbucketserver_repository_permissions_users" "my-repo-malcolms" {
  project    = "TEST"
  repository = "my-repo"
  filter     = "malcolm"
}
```

## Argument Reference

* `project` - Required. Project Key to lookup permissions for.
* `repository` - Required. Repository slug to lookup permissions for.
* `filter` - Optional. If specified only user's names/emails containing the supplied string will be returned.

## Attribute Reference

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are:

    * `REPO_READ`
    * `REPO_WRITE`
    * `REPO_ADMIN`
