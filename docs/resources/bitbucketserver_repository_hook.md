# Resource: bitbucketserver_repository_hook

Manage a repository level hook. Extends what Bitbucket does every time a repository changes, for example when code is pushed or a pull request is merged.

## Example Usage

```hcl
resource "bitbucketserver_repository_hook" "main" {
  project    = "MYPROJ"
  repository = "repo1"
  hook       = "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook"
}
```

## Argument Reference

* `project` - Required. Project Key the repository is contained within.
* `repository` - Required. Repository slug to enable hook for.
* `hook` - Required. The hook to enable on the repository.
* `settings` - Optional. Map of values to apply as settings for the hook. Contents dependant on the individual hook settings.
