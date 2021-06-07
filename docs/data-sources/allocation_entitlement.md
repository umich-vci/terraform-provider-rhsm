---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rhsm_allocation_entitlement Data Source - terraform-provider-rhsm"
subcategory: ""
description: |-
  Data source to get information about an entitlement associated with a Red Hat Subscription Manager allocation.
---

# rhsm_allocation_entitlement (Data Source)

Data source to get information about an entitlement associated with a Red Hat Subscription Manager allocation.

## Example Usage

```terraform
data "rhsm_allocation_entitlement" "entitlement" {
  allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
  entitlement_id  = "123e4567-e89b-12d3-a456-426655440000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **allocation_uuid** (String) The UUID of the subscription allocation to create the entitlement on.
- **entitlement_id** (String) The ID of the entitlement to look up in the specified allocation.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **contract_number** (String) The support contract associated with the entitlement.
- **quantity** (Number) The number of entitlements available in the pool.
- **sku** (String) The SKU of the entitlement.

