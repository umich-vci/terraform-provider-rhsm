---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rhsm_cloud_access_account Resource - terraform-provider-rhsm"
subcategory: ""
description: |-
  Resource to manage entitlement for Red Hat Cloud Access for an account in a supported cloud provider.
---

# rhsm_cloud_access_account (Resource)

Resource to manage entitlement for Red Hat Cloud Access for an account in a supported cloud provider.

## Example Usage

```terraform
resource "rhsm_cloud_access_account" "test_account" {
  account_id          = "123e4567-e89b-12d3-a456-426614174000"
  provider_short_name = "MSAZ"
  nickname            = "Test Account"
  gold_images         = ["rhel-byos"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **account_id** (String) The ID of a cloud account that you would like to request Red Hat Cloud Access for.
- **provider_short_name** (String) (Required) The short name of the cloud provider that the `account_id` is in. This must be one of "AWS", "GCE", or "MSAZ".  Other cloud providers are supported but have not been tested so they are not in the list of valid options.

### Optional

- **gold_images** (Set of String) A list of gold images to request access to for the account. Images available to a cloud provider can be found with the `rhsm_cloud_access` data source.
- **id** (String) The ID of this resource.
- **nickname** (String) A nickname to help describe the account. Defaults to ``.

### Read-Only

- **date_added** (String) The date the cloud account was added to Red Hat Cloud Access.
- **gold_image_status** (String) The status of requests for access to gold images.

