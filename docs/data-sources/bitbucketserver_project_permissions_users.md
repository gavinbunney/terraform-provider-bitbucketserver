# Data Source: bitbucketserver_project_permissions_users

Retrieve a list of users that have been granted at least one permission for the specified project.

## Example Usage

```hcl
data "bitbucketserver_project_permissions_users" "test-users" {
  project = "TEST"
}
```

### Applying a Custom Filter

Find project users starting with `malcolm`.
 
```hcl
data "bitbucketserver_project_permissions_users" "malcolms" {
  project = "TEST"
  filter  = "malcolm"
}
```

## Argument Reference

* `project` - Required. Project Key to lookup permissions for.
* `filter` - Optional. If specified only user's names/emails containing the supplied string will be returned.

## Attribute Reference

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are:

    * `PROJECT_READ`
    * `PROJECT_WRITE`
    * `PROJECT_ADMIN`
