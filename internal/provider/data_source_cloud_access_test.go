package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudAccess(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudAccess,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourceCloudAccess = `
data "rhsm_cloud_access" "ca" {}
`
