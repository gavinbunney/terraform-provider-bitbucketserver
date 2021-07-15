# Resource: bitbucketserver_user_group

Assign a User to an existing Bitbucket Group.

## Example Usage

```hcl
resource "bitbucketserver_user_group" "browncoat" {
  user  = "mreynolds"
  group = "browncoats"
}
```

## Argument Reference

* `user` - Required. User to assign group to.
* `group` - Required. Group to assign to the user.

## Import

Import a user group reference via the user and group keys:

```
terraform import bitbucketserver_user_group.browncoat mreynolds/browncoats
```
