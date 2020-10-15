# rhsm\_allocation\_entitlement Resource

Use this resource to create and attach an entitlement to a RHSM Subscription Allocation
for a Red Hat Satellite server.

## Example Usage

```hcl
resource "rhsm_allocation_entitlement" "entitlement" {
    allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
    pool = "123e4567e89b12d3a456426655440000"
    quantity = 5
}
```

## Argument Reference

* `allocation_uuid` - (Required) The UUID of the subscription allocation to create the entitlement on.

* `pool` - (Required) The ID of the pool you would like to create the entitlement from.

* `quantity` - (Required) The number of entitlements you would like add to the allocation/use from the pool.

## Attributes Reference

* `contract_number` - The support contract associated with the entitlement.

* `sku` - The SKU of the entitlement.
