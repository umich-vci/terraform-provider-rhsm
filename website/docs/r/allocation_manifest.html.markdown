---
layout: "rhsm"
page_title: "RHSM: rhsm_allocation"
sidebar_current: "docs-rhsm-resource-allocation"
description: |-
 Create a manifest from a Red Hat Subscription Manager allocation.
---

# rhsm\_allocation\_manifest

Use this resource to create a manifest from a RHSM subscription allocation that can be uploaded
to a Red Hat Satellite server.

## Example Usage

```hcl
resource "rhsm_allocation_manifest" "manifest" {
    allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
}
```

## Argument Reference

* `allocation_uuid` - (Required) The UUID of the subscription allocation to create the manifest from.

* `tainted` - (Optional) a boolean used to mark when the subscription allocation has been changed
  since the last time the manifest was generated.  This should be left unset.

## Attributes Reference

* `last_modified` - The date and time the subscription allocation was last modified.

* `manifest_last_modified` - The date and time the subscription allocation was last modified
  when the manifest was last generated.

* `manifest` - The manifest as downloaded from the RHSM portal.  This is a zip file which has been
  base64 encoded to a string.
