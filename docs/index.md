# Red Hat Subscription Manager Provider

 The Red Hat Subscription Manager (RHSM) provider is used to interact with Red Hat Subscription Manager.

## Example Usage

```hcl
// Configure the Satellite Provider
provider "rhsm" {
    refresh_token = "tokenhere"
}

// Create a new Satellite Allocation
resource "rhsm_allocation" "foo" {
  name = "foo"
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `refresh_token` - (Optional) This is the [offline token](https://access.redhat.com/articles/3626371#bgenerating-a-new-offline-tokenb-3)
  used to generate access tokens for Red Hat Subscription Manager. This must be provided in
  the config or in the environment variable `RHSM_REFRESH_TOKEN`.
