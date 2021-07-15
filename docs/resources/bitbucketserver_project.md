# Resource: bitbucketserver_project

Create a Bitbucket Project to hold repositories.

## Example Usage

```hcl
resource "bitbucketserver_project" "test" {
  key         = "TEST"
  name        = "test-01"
  description = "Test project"
  avatar      = "data:(content type, e.g. image/png);base64,(data)"
}
```

## Argument Reference

* `key` - Required. Project key to set.
* `name` - Required. Name of the project.
* `description` - Optional. Description of the project.
* `avatar` - Optional. Avatar to use containing base64-encoded image data. Format: `data:(content type, e.g. image/png);base64,(data)`
* `public` - Optional. Flag to make the project public or private. Default `false`.

## Import

Import a project reference via the key:

```
terraform import bitbucketserver_project.test TEST
```
