---
id: bitbucketserver_repository
title: bitbucketserver_repository
---

Create a Bitbucket Repository.

## Example Usage

```hcl
resource "bitbucketserver_repository" "test" {
  project     = "MYPROJ"
  name        = "test-01"
  description = "Test repository"
}
```

## Argument Reference

* `project` - Required. Name of the project to create the repository in.
* `name` - Required. Name of the repository.
* `slug` - Optional. Slug to use for the repository. Calculated if not defined.
* `description` - Optional. Description of the repository.
* `forkable` - Optional. Enable/disable forks of this repository. Default `true`
* `public` - Optional. Determine if this repository is public. Default `false`

## Attribute Reference

Additional to the above, the following attributes are emitted:

* `clone_ssh` - URL for SSH cloning of the repository.
* `clone_https` - URL for HTTPS cloning of the repository.

## Import

Import a repository using the project key and repository slug:

```
terraform import bitbucketserver_repository.test MYPROJ/test-01
```
