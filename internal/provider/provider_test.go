package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/umich-vci/terraform-provider-rhsm/internal/sdkprovider"
)

func testAccPreCheck(t *testing.T) {
	refreshToken := os.Getenv("RHSM_REFRESH_TOKEN")

	if refreshToken == "" {
		t.Fatalf("RHSM_REFRESH_TOKEN must be set for acceptance tests to run")
	}
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"rhsm": func() (tfprotov6.ProviderServer, error) {
			upgradedSdkProvider, err := tf5to6server.UpgradeServer(
				context.Background(),
				sdkprovider.New("test")().GRPCProvider,
			)

			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkProvider
				},

				providerserver.NewProtocol6(New("test")),
			}
			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudAccess,
			},
		},
	})
}

// func TestResource_UpgradeFromVersion(t *testing.T) {
// 	/* ... */
// 	resource.Test(t, resource.TestCase{
// 		Steps: []resource.TestStep{
// 			{
// 				ExternalProviders: map[string]resource.ExternalProvider{
// 					"rhsm": {
// 						VersionConstraint: "0.6.1",
// 						Source:            "hashicorp/<provider>",
// 					},
// 				},
// 				Config: `resource "rhsm_cloud_access_account" "aws_test_account" {
// 					account_id          = "012345678912"
// 					provider_short_name = "AWS"
// 					nickname            = "Terraform Acceptance Test AWS Account"
// 				  }`,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("provider_resource.aws_test_account", "<attribute>", "<value>"),
// 					/* ... */
// 				),
// 			},
// 			{
// 				ProtoV5ProviderFactories: protoV5ProviderFactories(),
// 				Config: `resource "provider_resource" "example" {
//                             /* ... */
//                         }`,
// 				// ConfigPlanChecks is a terraform-plugin-testing feature.
// 				// If acceptance testing is still using terraform-plugin-sdk/v2,
// 				// use `PlanOnly: true` instead. When migrating to
// 				// terraform-plugin-testing, switch to `ConfigPlanChecks` or you
// 				// will likely experience test failures.
// 				ConfigPlanChecks: resource.ConfigPlanChecks{
// 					PreApply: []plancheck.PlanCheck{
// 						plancheck.ExpectEmptyPlan(),
// 					},
// 				},
// 			},
// 		},
// 	})
// }
