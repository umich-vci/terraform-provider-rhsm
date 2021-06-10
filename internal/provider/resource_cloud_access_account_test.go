package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCloudAccessAccountAzure(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudAccessAccountAzure,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.azure_test_account", "nickname", "Terraform Acceptance Test Azure Account"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.azure_test_account", "account_id", "123e4567-e89b-12d3-a456-426614174000"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.azure_test_account", "provider_short_name", "MSAZ"),
				),
			},
		},
	})
}

const testAccResourceCloudAccessAccountAzure = `
resource "rhsm_cloud_access_account" "azure_test_account" {
	account_id          = "123e4567-e89b-12d3-a456-426614174000"
	provider_short_name = "MSAZ"
	nickname            = "Terraform Acceptance Test Azure Account"
  }
`
