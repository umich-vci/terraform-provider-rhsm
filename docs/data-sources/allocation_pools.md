# rhsm\_allocation\_pools Data Source

Get information about pools available to a Red Hat Subscription Manager allocation.

## Example Usage

```hcl
data "rhsm_allocation_pools" "pools" {
    allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
    pool = "123e4567e89b12d3a456426655440000"
    quantity = 5
}
```

## Argument Reference

* `allocation_uuid` - (Required) The UUID of the subscription allocation to create the entitlement on.

* `future` - (Optional) Should pools only valid in the future be listed.

## Attributes Reference

* `pools`

  * `contract_number` - The support contract associated with the entitlement.

  * `end_date` - The date the support contract ends.

  * `entitlements_available` - The number of entitlements available from the pool.

  * `id` - The ID of the pool.

  * `service_level` - The service level of the support contract.

  * `sku` - The SKU of the entitlement.

  * `start_date` - The date the support contract starts.

  * `subscription_name` - The friendly name of the sku.

  * `subscription_number` - The subscription number.
