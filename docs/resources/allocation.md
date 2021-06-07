---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rhsm_allocation Resource - terraform-provider-rhsm"
subcategory: ""
description: |-
  Resource to manage a RHSM Subscription allocation for a Red Hat Satellite server.
---

# rhsm_allocation (Resource)

Resource to manage a RHSM Subscription allocation for a Red Hat Satellite server.

## Example Usage

```terraform
resource "rhsm_allocation" "foo" {
  name = "foo"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the subscription allocation.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **created_by** (String) The user account used to create the subscription allocation.
- **created_date** (String) The date and time the subscription allocation was created.
- **entitlements_attached** (List of Object) (see [below for nested schema](#nestedatt--entitlements_attached))
- **entitlements_attached_quantity** (Number) The number of entitlements associated with the subscription
- **last_modified** (String) The date and time the subscription allocation was last modified.
- **type** (String) The type of the subscription allocation.  The only one supported by this resource is `Satellite`.
- **uuid** (String) The UUID of the subscription allocation that was created.
- **version** (String) The version of the subscription allocation type.  This defaults in the API to 6.5 and cannot be set through the API. It can be adjusted in the RHSM portal.

<a id="nestedatt--entitlements_attached"></a>
### Nested Schema for `entitlements_attached`

Read-Only:

- **reason** (String)
- **valid** (Boolean)

