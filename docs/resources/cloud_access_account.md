# rhsm\_cloud\_access\_account Resource

Use this resource to create an entitlement for Red Hat Cloud Access for an account in a supported cloud provider.

## Example Usage

```hcl
resource "rhsm_cloud_access_account" "test_account" {
  account_id = "123e4567-e89b-12d3-a456-426614174000"
  provider_short_name = "MSAZ"
  nickname = "Test Account"
  gold_images = ["rhel-byos"]
}
```

## Argument Reference

* `account_id` - (Required) The ID of a cloud account that you would like to request Red Hat Cloud Access for.

* `provider_short_name` - (Required) The short name of the cloud provider that the `account_id` is in.
  This must be one of "AWS", "GCE", or "MSAZ".  Other cloud providers are supported but have not been tested
  so they are not in the list of valid options.

* `nickname` - (Optional) A nickname to help describe the account.  The default is an empty string.

* `gold_images` - (Optional) A list of gold images to request access to for the account.
  Images available to a cloud provider can be found with the `rhsm_cloud_access` data source.

## Attributes Reference

* `date_added` - The date the cloud account was added to Red Hat Cloud Access.

* `gold_image_status` - The status of requests for access to gold images.
