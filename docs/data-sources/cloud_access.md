# rhsm\_cloud\_access Data Source

Use this data source to look up information about cloud providers entitled to Red Hat Cloud Access.

## Example Usage

```hcl
data "rhsm_cloud_access" "ca" {}
```

## Attributes Reference

* `enabled_accounts` - A list where each entry is a single cloud provider

In each element of `enabled_accounts` are the following attributes:

* `accounts` - A list of cloud accounts that are enabled for cloud access in the cloud provider.

  * `id` - The id of the cloud

  * `date_added` - The date the account was added to cloud access.

  * `gold_image_status` - The status of any requests for gold image access for a cloud account.

  * `nickname` - A nickname associated with the cloud account.

  * `name` - The name of the cloud provider.

* `products` - A list of products that are entitled to the cloud provider.

  * `enabled_quantity` - The quantity of subscriptions allowed to be consumed by the cloud provider.

  * `image_groups` - A list of images associated with the cloud provider. These are used when
      requesting access to gold images for a cloud account.

  * `name` - The name of the product.

  * `next_renewal` - The renewal date of the subscription.

  * `sku` - The SKU of the product.

  * `total_quantity` - The total number of subscriptions of the product available.

* `short_name` - An abreviation of the cloud provider name. Used when adding or removing accounts.
