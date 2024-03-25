package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccDataSourceCloudAccess(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudAccess,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestCloudAccessDataSource_UpgradeFromVersion(t *testing.T) {
	/* ... */
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rhsm": {
						VersionConstraint: "0.6.1",
						Source:            "umich-vci/rhsm",
					},
				},
				Config: testAccDataSourceCloudAccess,
			},
			{
				ProtoV6ProviderFactories: protoV6ProviderFactories(),
				Config:                   testAccDataSourceCloudAccess,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

const testAccDataSourceCloudAccess = `
data "rhsm_cloud_access" "ca" {}
`
