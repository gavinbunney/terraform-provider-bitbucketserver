# Resource: bitbucketserver_project_hook

Manage a project level hook. Extends what Bitbucket does every time a repository changes, for example when code is pushed or a pull request is merged.

## Example Usage

```hcl
resource "bitbucketserver_project_hook" "main" {
  project = "MYPROJ"
  hook    = "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook"
}
```

## Argument Reference

* `project` - Required. Project Key the hook to enable is for.
* `hook` - Required. The hook to enable on the project.
* `settings` - Optional. Map of values to apply as settings for the hook. Contents dependant on the individual hook settings.
