# Resource: bitbucketserver_user

Create a Bitbucket user.

## Example Usage

```hcl
resource "bitbucketserver_user" "admin" {
  name          = "mreynolds"
  display_name  = "Malcolm Reynolds"
  email_address = "browncoat@example.com"
}
```

## Argument Reference

* `name` - Required. Username of the user.
* `display_name` - Required. User's name to display.
* `email_address` - Required. Email address of user.
* `password_length` - Optional. The length of the generated password on resource creation. Only applies on resource creation. Default `20`.

## Attribute Reference

* `initial_password` - The generated user password. Only available if password was handled on Terraform resource creation, not import.
* `user_id` - The user ID.

## Import

Import a user reference via the user's name.

```
terraform import bitbucketserver_user.test mreynolds
```
