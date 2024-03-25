package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceCloudAccessAccountAzure(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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

func TestAccResourceCloudAccessAccountAWS(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudAccessAccountAWS,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.aws_test_account", "nickname", "Terraform Acceptance Test AWS Account"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.aws_test_account", "account_id", "012345678912"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.aws_test_account", "provider_short_name", "AWS"),
				),
			},
		},
	})
}

func TestAccResourceCloudAccessAccountGCP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudAccessAccountGCP,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.gcp_test_account", "nickname", "Terraform Acceptance Test GCP Group"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.gcp_test_account", "account_id", "my.group@example.com"),
					resource.TestCheckResourceAttr(
						"rhsm_cloud_access_account.gcp_test_account", "provider_short_name", "GCE"),
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

const testAccResourceCloudAccessAccountAWS = `
resource "rhsm_cloud_access_account" "aws_test_account" {
	account_id          = "012345678912"
	provider_short_name = "AWS"
	nickname            = "Terraform Acceptance Test AWS Account"
  }
`

const testAccResourceCloudAccessAccountGCP = `
resource "rhsm_cloud_access_account" "gcp_test_account" {
	account_id          = "my.group@example.com"
	provider_short_name = "GCE"
	nickname            = "Terraform Acceptance Test GCP Group"
  }
`
