---
layout: "rhsm"
page_title: "RHSM: rhsm_allocation_entitlement"
sidebar_current: "docs-rhsm-datasource-allocation-entitlement"
description: |-
 Get information about an entitlement associated with a Red Hat Subscription Manager allocation.
---

# rhsm\_allocation

Get information about an entitlement associated with a Red Hat Subscription Manager allocation.

## Example Usage

```hcl
data "rhsm_allocation_entitlement" "entitlement" {
    allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
    pool = "123e4567e89b12d3a456426655440000"
    quantity = 5
}
```

## Argument Reference

* `allocation_uuid` - (Required) The UUID of the subscription allocation to create the entitlement on.

* `entitlement_id` - (Required) The ID of the entitlement to look up in the specified allocation.

## Attributes Reference

* `contract_number` - The support contract associated with the entitlement.

* `sku` - The SKU of the entitlement.

* `quantity` - The number of entitlements you would like add to the allocation/use from the pool.
