## 0.6.0 (Unreleased)

ENHANCEMENTS:

* Updated [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk) to 2.23.0.
* Updated [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) to 0.13.0.
* Updated [gorhsm](https://github.com/umich-vci/gosatellite) to 1.300.0.
* Using go 1.18 instead of 1.17

## 0.5.0 (May 11, 2022)

ENHANCEMENTS:

* Updated [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk) to 2.16.0.
* Updated [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) to 0.8.1.
* Updated [gorhsm](https://github.com/umich-vci/gosatellite) to 1.264.0.
* `resource/rhsm_cloud_access_account` Updated examples and documentation.

## 0.4.0 (February 25, 2022)

ENHANCEMENTS:

* Updated [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk) to 2.10.1.
* Updated [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) to 0.5.0.
* Updated [gorhsm](https://github.com/umich-vci/gosatellite) to 1.196.0.
* Use StateContext instead of State as State is deprecated.
* Build for Darwin arm64

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

* Updated [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk) to 2.6.1.
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
