# Resource: bitbucketserver_default_reviewers_condition

Create a default reviewers condition for project or repository.

## Example Usage

```hcl
resource "bitbucketserver_default_reviewers_condition" "condition" {
	project_key			= "PRO"
	repository_slug		= "repository-1"
	source_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	target_matcher		= {
		id			= "any"
		type_id		= "ANY_REF"
	}
	reviewers          = [1]
	required_approvals = 1
}
```

## Argument Reference

* `project_key` - Required. Project key.
* `repository_slug` - Optional. Repository slug. If empty, default reviewers condition will be created for the whole project.
* `source_matcher.id` - Required. Source branch matcher id. It can be either `"any"` to match all branches, `"refs/heads/master"` to match certain branch, `"pattern"` to match multiple branches or `"development"` to match branching model.
* `source_matcher.type_id` - Required. Source branch matcher type.It must be one of: `"ANY_REF"`, `"BRANCH"`, `"PATTERN"`, `"MODEL_BRANCH"`.
* `target_matcher.id` - Required. Target branch matcher id. It can be either `"any"` to match all branches, `"refs/heads/master"` to match certain branch, `"pattern"` to match multiple branches or `"development"` to match branching model.
* `target_matcher.type_id` - Required. Target branch matcher type. It must be one of: `"ANY_REF"`, `"BRANCH"`, `"PATTERN"`, `"MODEL_BRANCH"`.
* `reviewers` - Required. IDs of Bitbucket users to become default reviewers when new pull request is created.
* `required_approvals` - Required. The number of default reviewers that must approve a pull request. Can't be higher than length of `reviewers`.

You can find more information about [how to use branch matchers here](https://confluence.atlassian.com/bitbucketserver/add-default-reviewers-to-pull-requests-834221295.html).

## Import

Import a default reviewers condition reference via the ID in this format `condition_id:project_key:repository_slug`.

```
terraform import bitbucketserver_default_reviewers_condition.test 1:pro:repo
```

When importing a default reviewers condition for the whole project omit the `repository_slug`.

```
terraform import bitbucketserver_default_reviewers_condition.test 1:pro
```
