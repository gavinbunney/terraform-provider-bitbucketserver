Bitbucket Server Terraform Provider
==================

[![Build Status](https://travis-ci.org/gavinbunney/terraform-provider-bitbucketserver.svg?branch=master)](https://travis-ci.org/gavinbunney/terraform-provider-bitbucketserver) [![codecov](https://codecov.io/gh/gavinbunney/terraform-provider-bitbucketserver/branch/master/graph/badge.svg)](https://codecov.io/gh/gavinbunney/terraform-provider-bitbucketserver) [![user guide](https://img.shields.io/badge/-user%20guide-blue)](https://gavinbunney.github.io/terraform-provider-bitbucketserver)

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

See [User Guide](https://gavinbunney.github.io/terraform-provider-bitbucketserver) for details on all the provided data and resource types.

---

## Development Guide

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Then clone the repository [Terraform provider bitbucket](https://github.com/gavinbunney/terraform-provider-bitbucketserver)  

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc-bitbucket`. Alternatively, you can run the 
`scripts/start-docker-compose.sh` to start the Bitbucket docker container, run `make testacc` and then call 
`scripts/stop-docker-compose.sh` when done running tests.

*Note:* Acceptance tests create real resources, and often cost money to run if not running locally.
