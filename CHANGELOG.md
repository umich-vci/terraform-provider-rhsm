## 0.3.0 (June 15, 2021)

BREAKING CHANGES:

* `resource/rhsm_cloud_access_account` The ID is now stored as provider_short_name:account_id instead of just account_id
  for consistency. It also allows for import to be used with `rhsm_cloud_access_account` resources.

BUG FIXES:

* `datasource/rhsm_cloud_access` Fix error due to goldImageStatus type changing in API ([#5]())
* `resource/rhsm_cloud_access_account` Fix error due to goldImageStatus type changing in API ([#5]())
* `datasource/rhsm_allocation` Add missing attributes
* `datasource/rhsm_allocation_entitlement` Add missing attributes

ENHANCEMENTS:

* Updated [terraform-plugin-sdk](https://github.com/umich-vci/gosatellite) to 2.6.1.
* Reworked code to model the approach in
  [terraform-provider-scaffoling](https://github.com/hashicorp/terraform-provider-scaffolding).
* Added descriptions to resources and data sources to allow for usage in documentation
  generation and in the language server.
* Added a few basic acceptance tests.

## 0.2.0 (October 16, 2020)

FEATURES:

* **New Data Source:** `rhsm_cloud_access`
* **New Resource:** `rhsm_cloud_access_account`

## 0.1.0 (October 15, 2020)

Initial release of provider to Terraform Registry.
