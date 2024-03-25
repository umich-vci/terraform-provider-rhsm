package sdkprovider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAllocationEntitlement(t *testing.T) {
	t.Skip("data source not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAllocationEntitlement,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.scaffolding_data_source.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccDataSourceAllocationEntitlement = `
data "scaffolding_data_source" "foo" {
  sample_attribute = "bar"
}
`
