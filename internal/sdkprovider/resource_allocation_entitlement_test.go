package sdkprovider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAllocationEntitlement(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
