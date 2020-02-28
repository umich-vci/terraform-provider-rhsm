---
layout: "rhsm"
page_title: "RHSM: rhsm_allocation"
sidebar_current: "docs-rhsm-resource-allocation"
description: |-
 Create a Red Hat Subscription Manager allocation.
---

# rhsm\_allocation

Use this resource to create a RHSM Subscription allocation for a Red Hat Satellite server.

## Example Usage

```hcl
resource "rhsm_allocation" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The name of the subscription allocation.

## Attributes Reference

* `uuid` - The UUID of the subscription allocation that was created. This is also used as
  the ID of the resource.

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
