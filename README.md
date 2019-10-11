Bitbucket Server Terraform Provider
==================

[![Build Status](https://travis-ci.org/gavinbunney/terraform-provider-bitbucketserver.svg?branch=master)](https://travis-ci.org/gavinbunney/terraform-provider-bitbucketserver)

This terraform provider allows management of bitbucket server resources.

> Note: The bundled terraform bitbucket provider works for bitbucket cloud - this provider is for bitbucket server only!

## Using the provider

Download a binary for your system from the release page and remove the `-os-arch` details so you're left with `terraform-provider-bitbucketserver`.
Use `chmod +x` to make it executable and then either place it at the root of your Terraform folder or in the Terraform plugin folder on your system. 

### Quick Start

```hcl
provider "bitbucketserver" {
  server   = "https://mybitbucket.example.com"
  username = "admin"
  password = "password"
}

resource "bitbucketserver_project" "test" {
  key         = "TEST"
  name        = "test-01"
  description = "Test project"
}

resource "bitbucketserver_repository" "test" {
  project     = bitbucketserver_project.test.key
  name        = "test-01"
  description = "Test repository"
}
```

### Provider Configuration

The provider supports parameters to determine the bitbucket server and admin user/password to use.

```hcl
provider "bitbucketserver" {
  server   = "https://mybitbucket.example.com"
  username = "admin"
  password = "password"
}
```

You can also specify these parameters through the `BITBUCKET_SERVER`, `BITBUCKER_USERNAME` and `BITBUCKET_PASSWORD` environment variables.


### Create a Bitbucket Project

```hcl
resource "bitbucketserver_project" "test" {
  key         = "TEST"
  name        = "test-01"
  description = "Test project"
  avatar      = "data:(content type, e.g. image/png);base64,(data)"
}
```

* `key` - Required. Project key to set.
* `name` - Required. Name of the project.
* `description` - Optional. Description of the project.
* `avatar` - Optional. Avatar to use containing base64-encoded image data. Format: `data:(content type, e.g. image/png);base64,(data)`

#### Import Project

```bash
$ terraform import bitbucketserver_project.test TEST
```


### Assign Project Permissions for Group

```hcl
resource "bitbucketserver_project_permissions_group" "test" {
  project = "TEST"
  group = "stash-users"
  permission = "PROJECT_WRITE"
}
```

* `project` - Required. Project key to set permissions for.
* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are: `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`

#### Import Group Permission for Project

```bash
$ terraform import bitbucketserver_project_permissions_group.test TEST/stash-users
```


### Assign Project Permissions for User

```hcl
resource "bitbucketserver_project_permissions_user" "test" {
  project = "TEST"
  user = "admin"
  permission = "PROJECT_WRITE"
}
```

* `project` - Required. Project key to set permissions for.
* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are: `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`

#### Import Group Permission for Project

```bash
$ terraform import bitbucketserver_project_permissions_user.test TEST/admin
```


### Create a Bitbucket Repository

```hcl
resource "bitbucketserver_repository" "test" {
  name        = "test-01"
  description = "Test repository"
}
```

* `name` - Required. Name of the project.
* `slug` - Optional. Slug to use for the repository. Calculated if not defined.
* `description` - Optional. Description of the repository.
* `forkable` - Optional. Enable/disable forks of this repository. Default `true`
* `public` - Optional. Determine if this repository is public. Default `false`

#### Attributes

Additional to the above, the following attributes are emitted:

* `clone_ssh` - URL for SSH cloning of the repository.
* `clone_https` - URL for HTTPS cloning of the repository.

#### Import Repository

```bash
$ terraform import bitbucketserver_repository.test TEST/test-01
```


### Assign Repository Permissions for Group

```hcl
resource "bitbucketserver_repository_permissions_group" "test" {
  project    = "TEST"
  repository = "my-repo"
  group      = "stash-users"
  permission = "REPO_WRITE"
}
```

* `project` - Required. Project key to set permissions for.
* `repository` - Required. Repository slug to set permissions for.
* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available repository permissions are: `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`

#### Import Group Permission for Repository

```bash
$ terraform import bitbucketserver_repository_permissions_group.test TEST/my-repo/stash-users
```


### Assign Repository Permissions for User

```hcl
resource "bitbucketserver_repository_permissions_user" "test" {
  project    = "TEST"
  repository = "my-repo"
  user       = "admin"
  permission = "REPO_WRITE"
}
```

* `project` - Required. Project key to set permissions for.
* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available repository permissions are: `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`

#### Import Group Permission for Repository

```bash
$ terraform import bitbucketserver_repository_permissions_user.test TEST/my-repo/admin
```


### Create a Bitbucket User

```hcl
resource "bitbucketserver_user" "admin" {
  name          = "mreynolds"
  display_name  = "Malcolm Reynolds"
  email_address = "browncoat@example.com"
}
```

* `name` - Required. Username of the user.
* `display_name` - Required. User's name to display.
* `email_address` - Required. Email address of user.
* `password_length` - Optional. The length of the generated password on resource creation. Only applies on resource creation. Default `20`.

#### Attributes

Additional to the above, the following attributes are emitted:

* `initial_password` - The generated user password. Only available if password was handled on Terraform resource creation, not import.

#### Import User

```bash
$ terraform import bitbucketserver_user.test mreynolds
```


### Create a Bitbucket Group

```hcl
resource "bitbucketserver_group" "browncoats" {
  name          = "browncoats"
}
```

* `name` - Required. Group to create.

#### Import Group

```bash
$ terraform import bitbucketserver_group.test browncoats
```


### Assign a User to a Bitbucket Group

```hcl
resource "bitbucketserver_user_group" "browncoat" {
  user  = "mreynolds"
  group = "browncoats"
}
```

* `user` - Required. User to assign group to.
* `group` - Required. Group to assign to the user.

#### Import Group

```bash
$ terraform import bitbucketserver_user_group.browncoat mreynolds/browncoats
```


### Assign Global Permissions for Group

```hcl
resource "bitbucketserver_global_permissions_group" "test" {
  project = "TEST"
  group = "stash-users"
  permission = "ADMIN"
}
```

* `group` - Required. Name of the group permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`

#### Import Global Group Permissions

```bash
$ terraform import bitbucketserver_global_permissions_group.test my-group
```


### Assign Global Permissions for User

```hcl
resource "bitbucketserver_project_permissions_user" "test" {
  user = "admin"
  permission = "ADMIN"
}
```

* `user` - Required. Name of the user permissions are for.
* `permission` - Required. The permission to grant. Available project permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`

#### Import Global Group Permissions

```bash
$ terraform import bitbucketserver_global_permissions_user.test my-user
```


### Set Server License

```hcl
resource "bitbucketserver_license" "main" {
  license = "AAACLg0ODAoPeNqNVEtv4jAQvudXRNpbp"
}
```

* `license` - Required. License to apply.

#### Attributes

Additional to the above, the following attributes are emitted:

* `clone_ssh` - URL for SSH cloning of the repository.
* `clone_https` - URL for HTTPS cloning of the repository.
* `creation_date` - License creation date.
* `purchase_date` - License purchase date.
* `expiry_date` - Expiry date of the license.
* `maintenance_expiry_date` - Expiry date of the maintenance period.
* `grace_period_end_date` - Grace period beyond expiry date.
* `maximum_users` - Maximum number of users for license.
* `unlimited_users` - Boolean flag if this is an unlimited user license.
* `server_id` - Server ID.
* `support_entitlement_number` - Entitlement number for support requests.

#### Import License

```bash
$ terraform import bitbucketserver_license.main license
```


### Set Mail Server Configuration

```hcl
resource "bitbucketserver_mail_server" "mail" {
  hostname       = "mail.example.com"
  port           = 465
  protocol       = "SMTPS"
  sender_address = "bitbucket@example.com"
}
```

* `hostname` - Required. Hostname of the mail server.
* `port` - Required. Port number of the mail server. Typically port 25 or 587 for SMTP and 465 for SMTPS.
* `sender_address` - Required. Email address for notification emails.
* `protocol` - Optional. SMTP or SMTPS supported. Default `SMTP`
* `use_start_tls` - Optional. Use SSL/TLS if available. Default `true`
* `require_start_tls` - Optional. Require SSL to be used. Default `false`
* `username` - Optional. User to connect with.
* `password` - Optional. User to connect with.

#### Import Mail Configuration

```bash
$ terraform import bitbucketserver_mail_server.mail mail.example.com
```

---

## Data Sources


### Application Properties

Retrieve version information and other application properties.

```hcl
data "bitbucketserver_application_properties" "main" {}
```

#### Attributes

* `version` - Version of Bitbucket.
* `build_number` - Build number of the Bitbucket instance.
* `build_date` - Date the Bitbucket build was made,
* `display_name` - Name of the Bitbucket instance.


### Groups

Retrieve a list of groups, optionally matching the supplied `filter`.

```hcl
data "bitbucketserver_groups" "all" {}
```

* `filter` - Optional. If specified only group names containing the supplied string will be returned.

#### Attributes

* `groups` - List of maps containing a `name` key.


### Users

Retrieve a list of users for a group,optionally matching the supplied `filter`.

```hcl
data "bitbucketserver_group_users" "stash-users" {
  group = "stash-users"
}
```

* `group` - Required. Group to find the users for.
* `filter` - Optional. If specified only group names containing the supplied string will be returned.

#### Attributes

* `users` - List of users containing `name`, `email_address`, `display_name` and `active` keys.


### Global Permissions - Groups

Retrieve a list of groups that have been granted at least one global permission.

```hcl
data "bitbucketserver_global_permissions_groups" "all" { }
```

* `filter` - Optional. If specified only group names containing the supplied string will be returned.

#### Attributes

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`


### Global Permissions - Users

Retrieve a list of users that have been granted at least one global permission.

```hcl
data "bitbucketserver_global_permissions_users" "proj" { }
```

* `filter` - Optional. If specified only user names containing the supplied string will be returned.

#### Attributes

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are: `LICENSED_USER`, `PROJECT_CREATE`, `ADMIN`, `SYS_ADMIN`


### Project Permissions - Groups

Retrieve a list of groups that have been granted at least one permission for the specified project.

```hcl
data "bitbucketserver_project_permissions_groups" "proj" {
  project = "TEST"
}
```

* `project` - Required. Project Key to lookup permissions for.
* `filter` - Optional. If specified only group names containing the supplied string will be returned.

#### Attributes

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are: `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`


### Project Permissions - Users

Retrieve a list of users that have been granted at least one permission for the specified project.

```hcl
data "bitbucketserver_project_permissions_users" "proj" {
  project = "TEST"
}
```

* `project` - Required. Project Key to lookup permissions for.
* `filter` - Optional. If specified only user names containing the supplied string will be returned.

#### Attributes

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are: `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`


### Repository Permissions - Groups

Retrieve a list of groups that have been granted at least one permission for the specified repository.

```hcl
data "bitbucketserver_project_permissions_groups" "proj" {
  project    = "TEST"
  repository = "my-repo"
}
```

* `project` - Required. Project Key to lookup permissions for.
* `repository` - Required. Repository slug to lookup permissions for.
* `filter` - Optional. If specified only group names containing the supplied string will be returned.

#### Attributes

* `groups` - List of maps containing `name` and `permission` keys. Available permissions are: `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`


### Repository Permissions - Users

Retrieve a list of users that have been granted at least one permission for the specified repository.

```hcl
data "bitbucketserver_project_permissions_users" "proj" {
  project    = "TEST"
  repository = "my-repo"
}
```

* `project` - Required. Project Key to lookup permissions for.
* `repository` - Required. Repository slug to lookup permissions for.
* `filter` - Optional. If specified only user names containing the supplied string will be returned.

#### Attributes

* `users` - List of maps containing `name`, `email_address`, `display_name`, `active` and `permission` keys. Available permissions are: `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`


### Plugin Details

Retrieve details of an installed plugin.

```hcl
data "bitbucketserver_plugin" "myplugin" {
  key = "com.atlassian.upm.atlassian-universal-plugin-manager-plugin"
}
```

#### Attributes

* `key` - Unique key of the plugin.
* `enabled` - Set to `true` if the plugin is enabled.
* `enabled_by_default` - Set to `true` if the plugin is enabled by default (for system plugins). 
* `version` - Installed version of the plugin. 
* `name` - Name of the plugin.
* `description` - Plugin description.
* `user_installed` - Set to `true` if this is a user installed plugin vs a system bundled plugin.
* `optional` - Set to `true` if this is an optional plugin.
* `vendor_name` - Name of the vendor.
* `vendor_link` - Vendor homepage.
* `vendor_marketplace_link` - Plugin marketplace link.

---

## Development Guide

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

### Building The Provider

Clone repository to: `$GOPATH/src/github.com/gavinbunney/terraform-provider-bitbucketserver`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:gavinbunney/terraform-provider-bitbucketserver
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/gavinbunney/terraform-provider-bitbucketserver
$ make build
```

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-bitbucketserver
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
