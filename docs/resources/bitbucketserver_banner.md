# Resource: bitbucketserver_banner

Manage the announcement banner, updating as required.

## Example Usage

```hcl
resource "bitbucketserver_banner" "main" {
  message = "Bitbucket is down for maintenance\n*Save your work*"
}
```

## Argument Reference

* `message` - Required. Information to display to the user. Markdown supported.
* `enabled` - Optional. Turn the announcement banner on/off. Default `true`.
* `audience` - Optional. Set the audience for the announcement. Must be one of `ALL` or `AUTHENTICATED`. Default `ALL`.

## Import

Import the banner:

```
terraform import bitbucketserver_banner.main banner
```
