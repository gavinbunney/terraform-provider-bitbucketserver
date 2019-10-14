---
id: bitbucketserver_group
title: bitbucketserver_group
---

Create a Bitbucket group.

## Example Usage

```hcl
resource "bitbucketserver_group" "browncoats" {
  name = "browncoats"
}
```

## Argument Reference

* `name` - Required. Group to create.

## Import

Import a group reference via the name.

```
terraform import bitbucketserver_group.test browncoats
```
