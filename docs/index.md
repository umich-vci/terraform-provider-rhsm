---
page_title: "Provider: RHSM"
subcategory: ""
description: |-
  The Red Hat Subscription Manager (RHSM) provider is used to interact with Red Hat Subscription Manager.
---

# Red Hat Subscription Manager Provider

The Red Hat Subscription Manager (RHSM) provider is used to interact with Red Hat Subscription Manager.

## Example Usage

```terraform
provider "rhsm" {
  refresh_token = "tokenhere"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `refresh_token` (String) This is the [offline token](https://access.redhat.com/articles/3626371#bgenerating-a-new-offline-tokenb-3) used to generate access tokens for Red Hat Subscription Manager. This must be provided in the config or in the environment variable `RHSM_REFRESH_TOKEN`.
