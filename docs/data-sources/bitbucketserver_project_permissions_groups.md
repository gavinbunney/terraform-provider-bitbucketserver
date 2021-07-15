# Data Source: bitbucketserver_project_permissions_groups

Retrieve a list of groups that have been granted at least one project level permission to the specified project.

## Example Usage

```hcl
data "bitbucketserver_project_permissions_groups" "test-groups" {
  project = "TEST"
}
```

### Applying a Custom Filter

Find project groups starting with `dev` with project permissions.
 
```hcl
data "bitbucketserver_project_permissions_groups" "dev-groups" {
  project = "TEST"
  filter  = "dev"
}
```

## Argument Reference

* `project` - Required. Project Key to lookup permissions for.
* `filter` - Optional. If specified only group names containing the supplied string will be returned.

## Attribute Reference

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are:

    * `PROJECT_READ`
    * `PROJECT_WRITE`
    * `PROJECT_ADMIN`
