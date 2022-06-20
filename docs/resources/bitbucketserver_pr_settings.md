# Resource: bitbucketserver_pr_settings

Provides the ability to manage pull request settings for the context repository.

## Example Usage

```hcl
resource "bitbucketserver_pr_settings" "test" {
  project                     = "MYPROJ"
  repository                  = "repo"
  no_needs_work_status        = true
  required_all_approvers      = true
  required_all_tasks_complete = true
  required_approvers          = 1
  required_successful_builds  = 1
  merge_config {
    default_strategy   = "no-ff"
    enabled_strategies = ["no-ff", "ff"]
    commit_summaries   = 30
  }
}
```

## Argument Reference

* `project` - Required. Project Key that contains target repository.
* `repository` - Required. Repository slug of target repository.
* `merge_config.default_strategy` - Required. Default [merge strategy](https://confluence.atlassian.com/bitbucketserver0717/pull-request-merge-strategies-1087535782.html?utm_campaign=in-app-help&amp%3Butm_source=stash&amp%3Butm_medium=in-app-help). Git merge strategies affect the way the Git history appears after merging a pull request. Must be one of `no-ff`, `ff`, `ff-only`, `rebase-no-ff`, `rebase-ff-only`, `squash`, `squash-ff-only`.
* `merge_config.enabled_strategies` - Required. List of enabled merge strategies. Must contain at least the strategy that you specify as the default one.
* `merge_config.commit_summaries` - Optional. Controls the number of commit summaries included in commit messages for pull requests. Default `20`.
* `required_approvers` - Optional. The number of approvals required on a pull request for it to be mergeable. Default `0`.
* `required_successful_builds` - Optional. The number of successful builds on a pull request for it to be mergeable. Default `0`.
* `required_all_approvers` - Optional. Whether or not all approvers must approve a pull request for it to be mergeable. Default `false`.
* `required_all_tasks_complete` - Optional. Whether or not all tasks on a pull request need to be completed for it to be mergeable. Default `false`.
* `no_needs_work_status` - Optional. Whether or not to block the merge if any reviewers have marked the pull request as 'needs work'. Default `false`.
