# Resource: bitbucketserver_license

Set the license for the bitbucket server.

## Example Usage

```hcl
resource "bitbucketserver_license" "main" {
  license = "AAACLg0ODAoPeNqNVEtv4jAQvudXRNpbp"
}
```

## Argument Reference

* `license` - Required. License to apply.

## Attribute Reference

* `creation_date` - License creation date.
* `purchase_date` - License purchase date.
* `expiry_date` - Expiry date of the license.
* `maintenance_expiry_date` - Expiry date of the maintenance period.
* `grace_period_end_date` - Grace period beyond expiry date.
* `maximum_users` - Maximum number of users for license.
* `unlimited_users` - Boolean flag if this is an unlimited user license.
* `server_id` - Server ID.
* `support_entitlement_number` - Entitlement number for support requests.

## Import

Import the license details:

```
terraform import bitbucketserver_license.main license
```
