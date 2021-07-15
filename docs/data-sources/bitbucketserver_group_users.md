# Data Source: bitbucketserver_group_users

Retrieve a list of users for a specific group.

## Example Usage

```hcl
data "bitbucketserver_group_users" "stash-users" {
  group = "stash-users"
}
```

### Applying a Custom Filter

Find any users starting with `malcolm`.
 
```hcl
data "bitbucketserver_group_users" "malcolms" {
  group  = "stash-users"
  filter = "malcolm"
}
```

## Argument Reference

* `group` - Required. Group to find the users for.
* `filter` - Optional. If specified only users matching name/email for the supplied string will be returned.

## Attribute Reference

* `users` - List of maps containing `name`, `email_address`, `display_name` and `active` keys.
