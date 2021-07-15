# Data Source: bitbucketserver_user

This data source allows you to retrieve Bitbucket user details.

## Example Usage

```hcl
data "bitbucketserver_user" "admin" {
  name = "admin"
}
```

## Argument Reference

* `name` - Unique name of the user.

## Attribute Reference

* `email_address` - User's email.
* `display_name` - User's display name.
* `user_id` - User's ID.
