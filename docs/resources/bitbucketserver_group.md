# Resource: bitbucketserver_group

Create a Bitbucket group.

## Example Usage

```hcl
resource "bitbucketserver_group" "browncoats" {
  name = "browncoats"
}
```

## Argument Reference

* `name` - Required. Group to create.
* `import_if_exists` - Optional. Import groups that already exist in bitbucket into the terraform state file.

## Import

Import a group reference via the name.

```
terraform import bitbucketserver_group.test browncoats
```
