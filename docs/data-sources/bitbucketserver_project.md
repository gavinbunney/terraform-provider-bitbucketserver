# Data Source: bitbucketserver_project_hooks

Retrieve a list of projects.

## Example Usage

```hcl
data "bitbucketserver_project" "main" {
  project = "TEST"
}

#  data.bitbucketserver_project.main = [{
#     "key"         = "TEST",
#     "name"        = "Test Project",
#     "description" = "A test project",
#     "public"      = "false",
#     "avatar"      = "avatar.png",
#  }]
```

### Applying a Custom Filter

Find a specific project.
 
```hcl
data "bitbucketserver_project" "test" {
  project = "TEST"
}
```

## Argument Reference

* `key` - Required. Project key to lookup for.

