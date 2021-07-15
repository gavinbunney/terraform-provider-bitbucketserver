# Data Source: bitbucketserver_plugin

This data source allows you to retrieve installed plugin information and license details.

## Example Usage

```hcl
data "bitbucketserver_plugin" "my-plugin" {
  key = "com.example.plugin-my-plugin"
}
```

## Argument Reference

* `key` - Unique key of the plugin.

## Attribute Reference

* `enabled` - Set to `true` if the plugin is enabled.
* `enabled_by_default` - Set to `true` if the plugin is enabled by default (for system plugins). 
* `version` - Installed version of the plugin. 
* `name` - Name of the plugin.
* `description` - Plugin description.
* `user_installed` - Set to `true` if this is a user installed plugin vs a system bundled plugin.
* `optional` - Set to `true` if this is an optional plugin.
* `vendor.name` - Name of the vendor.
* `vendor.link` - Vendor homepage.
* `vendor.marketplace_link` - Plugin marketplace link.
* `applied_license.0.valid` - Is the license valid. true/false.
* `applied_license.0.evaluation` - Is the license an evaluation. true/false.
* `applied_license.0.nearly_expired` - Is the license nearly expired. true/false.
* `applied_license.0.maintenance_expiry_date` - Date of maintenance expiry.
* `applied_license.0.maintenance_expired` - Is the maintenance expired. true/false.
* `applied_license.0.license_type` - Type of license.
* `applied_license.0.expiry_date` - Expiry date of the license.
* `applied_license.0.raw_license` - The raw license information.
* `applied_license.0.renewable` - Is the license renewabl. true/false.
* `applied_license.0.organization_name` - Name of the organization the license is for.
* `applied_license.0.enterprise` - Is the license for enterprise. true/false.
* `applied_license.0.data_center` - Is the license for data center. true/false.
* `applied_license.0.subscription` - Is the license a subscription. true/false.
* `applied_license.0.active` - Is the license active. true/false.
* `applied_license.0.auto_renewal` - Is the license renewed automatically. true/false.
* `applied_license.0.upgradable` - Is the license able to be upgraded. true/false.
* `applied_license.0.crossgradeable` - Can the license be crossgraded. true/false.
* `applied_license.0.purchase_past_server_cutoff_date` - The purchase date past the server cutoff date. true/false.
