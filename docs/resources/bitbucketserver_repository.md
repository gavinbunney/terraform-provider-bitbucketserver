# Resource: bitbucketserver_repository

Create a Bitbucket Repository.

## Example Usage

```hcl
resource "bitbucketserver_repository" "test" {
  project     = "MYPROJ"
  name        = "test-01"
  description = "Test repository"
}
```

### Forking an existing repository

```hcl
resource "bitbucketserver_repository" "test" {
  project                 = "MYPROJ"
  name                    = "test-01"
  description             = "Test repository"
  fork_repository_project = "MY-ORIGIN-PROJ"
  fork_repository_slug    = "MY-ORIGIN-REPO"
}
```

> Note: Both `fork_repository_project` and `fork_repository_slug` are required to specified the origin repository to fork.

## Argument Reference

* `project` - Required. Name of the project to create the repository in.
* `name` - Required. Name of the repository.
* `slug` - Optional. Slug to use for the repository. Calculated if not defined.
* `description` - Optional. Description of the repository.
* `forkable` - Optional. Enable/disable forks of this repository. Default `true`
* `public` - Optional. Determine if this repository is public. Default `false`
* `enable_git_lfs` - Optional. Enable git-lfs for this repository. Default `false`
* `fork_repository_project` - Optional. Use this to fork an existing repository from the given project.
* `fork_repository_slug` - Optional. Use this to fork an existing repository from the given repository.

## Attribute Reference

Additional to the above, the following attributes are emitted:

* `clone_ssh` - URL for SSH cloning of the repository.
* `clone_https` - URL for HTTPS cloning of the repository.

## Import

Import a repository using the project key and repository slug:

```
terraform import bitbucketserver_repository.test MYPROJ/test-01
```
