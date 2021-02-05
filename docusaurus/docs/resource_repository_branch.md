---
id: bitbucketserver_repository_branch
title: bitbucketserver_repository_branch
---

Manage the lifecycle of new branches in a Bitbucket repository.

## Example Usage

```hcl
resource "bitbucketserver_repository_branch" "main" {
  project     = "MYPROJ"
  repository  = "repo1"
  branch_name = "ft/new-branch"
  source_ref  = "refs/head/master"
}
```

## Argument Reference

* `project` - Required. Project Key the repository is contained within.
* `repository` - Required. Repository slug to enable hook for.
* `branch_name` - Required. The display name of the new branch to be created.
* `source_ref` - Optional. The full reference name of the source branch to create from. Default is `refs/head/master`
