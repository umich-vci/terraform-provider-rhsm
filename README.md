<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Red Hat Subscription Manager

A Terraform provider for Red Hat Subscription Manager (RHSM)

This provider can be used to create and manage allocations, entitlements, and manifests for Red Hat Subscriptions
that can be used with Red Hat Satellite.

Support will be added in the future to manage Red Hat Cloud Access entitlements for cloud accounts.

The provider does not have working tests so it should probably be considered beta.

## Building/Installing

The provider is now available on the [Terraform Registry](https://registry.terraform.io/providers/umich-vci/rhsm/latest) so you probably don't need to build the provider.

That said, running `GO111MODULE=on go get -u github.com/umich-vci/terraform-provider-rhsm` should download
the code and result in a binary at `$GOPATH/bin/terraform-provider-rhsm`. You can then move the
binary to `~/.terraform.d/plugins` to use it with Terraform.

This has been tested with Terraform 0.12.x and Red Hat Satellite 6.6.x and 6.7.x.

## License

This project is licensed under the Mozilla Public License Version 2.0.
