# Resource: bitbucketserver_mail_server

Setup mail server configuration.

## Example Usage

```hcl
resource "bitbucketserver_mail_server" "mail" {
  hostname       = "mail.example.com"
  port           = 465
  protocol       = "SMTPS"
  sender_address = "bitbucket@example.com"
}
```

## Argument Reference

* `hostname` - Required. Hostname of the mail server.
* `port` - Required. Port number of the mail server. Typically port 25 or 587 for SMTP and 465 for SMTPS.
* `sender_address` - Required. Email address for notification emails.
* `protocol` - Optional. SMTP or SMTPS supported. Default `SMTP`
* `use_start_tls` - Optional. Use SSL/TLS if available. Default `true`
* `require_start_tls` - Optional. Require SSL to be used. Default `false`
* `username` - Optional. User to connect with.
* `password` - Optional. User to connect with.

## Import

Import the existing mail server configuration with the hostname:

```
terraform import bitbucketserver_mail_server.mail mail.example.com
```
