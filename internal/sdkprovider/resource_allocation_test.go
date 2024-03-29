package sdkprovider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAllocation(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAllocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rhsm_allocation.test", "name", "TerraformAcceptanceTestAllocation"),
					resource.TestMatchResourceAttr(
						"rhsm_allocation.test", "uuid", regexp.MustCompile(`[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}`)),
				),
			},
		},
	})
}

const testAccResourceAllocation = `
resource "rhsm_allocation" "test" {
	name = "TerraformAcceptanceTestAllocation"
  }
`
