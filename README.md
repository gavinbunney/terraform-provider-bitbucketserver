# Bitbucket Server Terraform Provider

[![user guide](https://img.shields.io/badge/-user%20guide-blue)](https://registry.terraform.io/providers/gavinbunney/bitbucketserver/latest)

This terraform provider allows management of **Bitbucket Server** resources. The bundled terraform bitbucket provider works only for Bitbucket Cloud.

## Using the provider

Download a binary for your system from the release page and remove the `-os-arch` details so you're left with `terraform-provider-bitbucketserver`.
Use `chmod +x` to make it executable and then either place it at the root of your Terraform folder or in the Terraform plugin folder on your system.

See [User Guide](https://gavinbunney.github.io/terraform-provider-bitbucketserver) for details on all the provided data and resource types.

### Example

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

## Development Guide

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11+
    - correctly setup [GOPATH](http://golang.org/doc/code.html#GOPATH
    - add `$GOPATH/bin` to your `$PATH`
- clone this repository to `$GOPATH/src/github.com/gavinbunney/terraform-provider-bitbucketserver`

### Building the provider

To build the provider, run `make build`. This will also put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
```

### Testing

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc-bitbucket`.

```sh
$ make testacc-bitbucket
```

Alternatively, you can manually start Bitbucket Server docker container, run the acceptance tests and then shut down the docker.

```sh
$ scripts/start-docker-compose.sh
$ make testacc
$ scripts/stop-docker-compose.sh
```
