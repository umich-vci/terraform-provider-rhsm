# Terraform Provider for Red Hat Subscription Manager

A Terraform provider for Red Hat Subscription Manager (RHSM)

This provider can be used to manage Red Hat Cloud Access for various cloud provider accounts.  Requesting access to
gold images is supported as well for AWS, Azure, and GCP.

This provider can also be used to create and manage allocations, entitlements, and manifests for Red Hat Subscriptions
that can be used with Red Hat Satellite. This functionality is deprecated.

## Using the provider

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/umich-vci/rhsm/latest) so you
probably don't need to build the provider unless you want to contribute.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.1.16
-	[Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
