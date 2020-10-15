# rhsm\_allocation Data Source

Use this data source to look up a RHSM Subscription allocation.

## Example Usage

```hcl
resource "rhsm_allocation" "foo" {
    uuid = "123e4567-e89b-12d3-a456-426655440000"
}
```

## Argument Reference

* `uuid` - The UUID of the subscription allocation to look up.

## Attributes Reference

* `name` - The name of the subscription allocation.

* `type` - The type of the subscription allocation.  The only one supported by this resource is `Satellite`.

* `version` - The version of the subscription allocation type.  This defaults to 6.5, but can be adjusted
  in the RHSM portal.

* `created_date` - The date and time the subscription allocation was created.

* `created_by` - The user account used to create the subscription allocation.

* `last_modified` - The date and time the subscription allocation was last modified.

* `entitlements_attached_quantity` - The number of entitlements associated with the subscription

* `entitlements_attached`

  * `reason` - The reason for the value of `valid`.

  * `valid` - If the entitlements associated with the subscription allocation are valid or not.
