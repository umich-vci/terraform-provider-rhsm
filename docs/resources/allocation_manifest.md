# rhsm\_allocation\_manifest Resource

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

## Attributes Reference

* `last_modified` - The date and time the subscription allocation was last modified.

* `manifest_last_modified` - The date and time the subscription allocation was last modified
  when the manifest was last generated.

* `manifest` - The manifest as downloaded from the RHSM portal.  This is a zip file which has been
  base64 encoded to a string.
