package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAllocationEntitlement(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAllocationEntitlement,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"scaffolding_resource.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceAllocationEntitlement = `
resource "scaffolding_resource" "foo" {
  sample_attribute = "bar"
}
`
