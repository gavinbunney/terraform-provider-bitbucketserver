# Resource: bitbucketserver_user_access_token

Personal access tokens can be used to replace passwords over https, or to authenticate using the Bitbucket Server REST API over Basic Auth. 

For git operations, you can use your personal access token as a substitute for your password.
   
> Note: You can only create access tokens for your user account - i.e. the one that the provisioner has been configured to authenticate with!
> This is a restriction in the Bitbucket APIs.

## Example Usage

```hcl
resource "bitbucketserver_user_access_token" "token" {
  user        = "admin"
  name        = "my-token"
  permissions = ["REPO_READ", "PROJECT_ADMIN"]
}
```

## Argument Reference

* `user` - Required. Username of the user.
* `name` - Required. Name of the access token.
* `permissions` - Required. List of permissions to grant the access token.

     * `PROJECT_READ`
     * `PROJECT_WRITE`
     * `PROJECT_ADMIN`
     * `REPO_READ`
     * `REPO_WRITE`
     * `REPO_ADMIN`

## Attribute Reference

* `access_token` - The generated access token. Only available if token was generated on Terraform resource creation, not import/update.
* `created_date` - When the access token was generated.
* `last_authenticated` - When the access token was last used for authentication.

## Import

Import a user token reference via the token id.

```
terraform import bitbucketserver_user_access_token.test 413460754380
```
