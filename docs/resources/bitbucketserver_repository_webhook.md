# Resource: bitbucketserver_repository_webhook

Manage a repository level webhook. Extends what Bitbucket does every time a repository changes, for example when code is pushed or a pull request is merged.

## Example Usage

```hcl
resource "bitbucketserver_project" "main" {
  key  = "MYPROJ"
  name = "my-project"
}

resource "bitbucketserver_repository" "main" {
  project = bitbucketserver_project.test.key
  name    = "repo"
}

resource "bitbucketserver_repository_webhook" "main" {
  project     = bitbucketserver_project.test.key
  repository  = bitbucketserver_repository.test.slug
  name        = "google"
  webhook_url = "https://www.google.com/"
  secret      = "abc"
  events      = ["repo:refs_changed"]
  active      = true
}
```

## Argument Reference

* `project` - Required. Project Key the repository is contained within.
* `repository` - Required. Repository slug to enable hook for.
* `name` - Required. Name of the webhook.
* `webhook_url` - Required. The URL of the webhook.
* `secret` - Optional. Secret used to authenticate the payload.
* `events` - Required. A list of events to trigger the webhook url.
* `active` - Optional. Enable or disable the webhook. Default: true

## Attribute Reference

* `webhook_id` - The webhook id.

## Import

Import a user reference using the project key, repository name and webhook name.

```
terraform import bitbucketserver_repository_webhook.main MYPROJ/repo/google
```
