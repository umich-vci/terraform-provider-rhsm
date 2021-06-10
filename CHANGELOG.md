## 0.3.0 (unreleased)

BUG FIXES:

* `datasource/rhsm_cloud_access` Fix crash due to goldImageStatus value changing in API ([#5]())

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
