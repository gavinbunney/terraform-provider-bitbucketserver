# Resource: bitbucketserver_project_branch_permissions

Provides the ability to apply project ref restrictions to enforce branch permissions. A restriction means preventing writes on the specified branch(es) by all except a set of users and/or groups, or preventing specific operations such as branch deletion.

## Example Usage

```hcl
resource "bitbucketserver_project_branch_permissions" "pr_only" {
  project          = "MYPROJ"
  repository       = "repo"
  ref_pattern      = "refs/heads/master"
  type             = "pull-request-only"
}

resource "bitbucketserver_project_branch_permissions" "no_deletes" {
  project     = "MYPROJ"
  repository  = "repo"
  ref_pattern = "heads/**/master"
  type        = "no-deletes"
  exception_users  = ["admin"]
  exception_groups = ["group_1", "group_2"]
}
```

## Argument Reference

* `project` - Required. Project Key that contains target repository.
* `repository` - Required. Repository slug of target repository.
* `ref_pattern` - Required. A wildcard pattern that may match multiple branches. You must specify a valid [Branch Permission Pattern](https://confluence.atlassian.com/bitbucketserver/branch-permission-patterns-776639814.html).
* `type` - Required. Type of the restriction. Must be one of `pull-request-only`, `fast-forward-only`, `no-deletes`, `read-only`.
* `exception_users` - Optional. List of usernames to whom restrictions do not apply.
* `exception_groups` - Optional. List of group names to which restrictions do not apply.
* `exception_access_keys` - Optional. List of access keys IDs to which restrictions do not apply.
