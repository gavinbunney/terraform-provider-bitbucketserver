# Data Source: bitbucketserver_repository_hooks

Retrieve a list of repository hooks and their status' for the specified repo.

## Example Usage

```hcl
data "bitbucketserver_repository_hooks" "main" {
  project    = "TEST"
  repository = "repo1"
}

#  data.bitbucketserver_repository_hooks.main.hooks = [{
#     "key"         = "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook",
#     "name"        = "Reject Force Push",
#     "type"        = "PRE_RECEIVE",
#     "description" = "Reject all force pushes (git push --force) to this repository",
#     "version"     = "6.10.0",
#     "scope_types" = ["PROJECT", "REPOSITORY"],
#     "enabled"     = "false",
#     "configured"  = "false",
#     "scope_type"  = "REPOSITORY",
#  }]
```

### Applying a Custom Filter

Find specific types of repository hooks.
 
```hcl
data "bitbucketserver_project_hooks" "main" {
  project    = "TEST"
  repository = "repo1"
  type       = "PRE_RECEIVE"
}
```

## Argument Reference

* `project` - Required. Project Key the repository is contained within.
* `repository` - Required. Repository slug to lookup hooks for.
* `type` - Optional. Type of hook to find. Must be one of `PRE_RECEIVE`, `POST_RECEIVE`

## Attribute Reference

* `hooks` - List of maps containing:

    * `key` - Unique key identifying the hook e.g. `com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook`
    * `name` - Name of the hook e.g. `Reject Force Push`
    * `type` - Type of the hook e.g. `PRE_RECEIVE`
    * `description` - Detailed description e.g. `Reject all force pushes (git push --force) to this repository`
    * `version` - Version of the hook, for system hooks this is the bitbucket version e.g. `6.10.0`
    * `scope_types` - List of strings containing the scopes available for this hook, e.g. `["PROJECT", "REPOSITORY"]`
    * `enabled` - Set if this hook is enabled for this project
    * `configured` - Set if the hook is configured for this project 
    * `scope_type` - Type of scope applied for this hook, e.g. `REPOSITORY`
    * `scope_resource_id` - Reference ID of the applied scope, e.g. `1`
