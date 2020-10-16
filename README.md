<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Red Hat Subscription Manager

A Terraform provider for Red Hat Subscription Manager (RHSM)

This provider can be used to create and manage allocations, entitlements, and manifests for Red Hat Subscriptions
that can be used with Red Hat Satellite.

It can also be used to manage Red Hat Cloud Access for various cloud provider accounts.  Requesting access to
gold images is supported as well, but presently is only supported on Azure by Red Hat.

The provider does not have working tests so it should probably be considered beta.

## Building/Installing

The provider is now available on the [Terraform Registry](https://registry.terraform.io/providers/umich-vci/rhsm/latest) so you probably don't need to build the provider unless you want to contribute.

That said, running `GO111MODULE=on go get -u github.com/umich-vci/terraform-provider-rhsm` should download
the code and result in a binary at `$GOPATH/bin/terraform-provider-rhsm`. You can then move the
binary to `~/.terraform.d/plugins` to use it with Terraform.

This has been tested with Terraform 0.12.x.

## License

This project is licensed under the Mozilla Public License Version 2.0.
