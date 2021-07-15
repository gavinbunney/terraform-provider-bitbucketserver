---
id: data_bitbucketserver_global_permissions_users
title: bitbucketserver_global_permissions_users
---

Retrieve a list of users that have been granted at least one global permission.

## Example Usage

```hcl
data "bitbucketserver_global_permissions_users" "all" { }
```

### Applying a Custom Filter

Find any users starting with `malcolm`.
 
```hcl
data "bitbucketserver_global_permissions_users" "malcolms" {
  filter = "malcolm"
}
```

## Argument Reference

* `filter` - Optional. If specified only user's names/emails containing the supplied string will be returned.

## Attribute Reference

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`

    * `LICENSED_USER`
    * `PROJECT_CREATE`
    * `ADMIN`
    * `SYS_ADMIN`

